package server

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"fmt"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/transactions"
	"io"
	"log"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/mikhaylov123ty/go-diploma-5.6/internal/models"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/server/api"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/balance"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/orders"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/users"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/storage/withdrawals"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/utils"

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const (
	userHandlerPath = "/api/user"
)

type Server struct {
	address      string
	transactions transactions.Handler
	usersRepo    users.Storage
	ordersRepo   orders.Storage
	balanceRepo  balance.Storage
	withdrawRepo withdrawals.Storage
	secretKey    string
}

type Claims struct {
	jwt.RegisteredClaims
	Login string
}

func New(
	address string,
	transactions transactions.Handler,
	usersRepo users.Storage,
	ordersRepo orders.Storage,
	balanceRepo balance.Storage,
	witdrawRepo withdrawals.Storage,
	secretKey string) *Server {
	return &Server{
		address:      address,
		transactions: transactions,
		usersRepo:    usersRepo,
		ordersRepo:   ordersRepo,
		balanceRepo:  balanceRepo,
		withdrawRepo: witdrawRepo,
		secretKey:    secretKey,
	}
}

func (s *Server) Start() error {
	router := s.newRouter()

	slog.Info("starting server", slog.String("address", s.address))
	return http.ListenAndServe(s.address, router)
}

func (s *Server) newRouter() *chi.Mux {
	router := chi.NewRouter()
	router.Use(s.withlogger, s.withGZipEncode)

	router.Route(userHandlerPath, func(router chi.Router) {

		router.Route("/register", func(router chi.Router) {
			router.Post("/", s.signToken(api.NewRegisterHandler(s.usersRepo, s.transactions).Handle))
		})

		router.Route("/login", func(router chi.Router) {
			router.Post("/", s.signToken(api.NewAuthHandler(s.usersRepo, s.transactions).Handle))
		})

		router.Route("/orders", func(router chi.Router) {
			router.Post("/", s.authHandler(api.NewPostOrdersHandler(s.ordersRepo, s.usersRepo, s.transactions).Handle))
			router.Get("/", s.authHandler(api.NewGetOrdersHandler(s.ordersRepo, s.usersRepo, s.transactions).Handle))
		})

		router.Route("/balance", func(router chi.Router) {
			router.Get("/", s.authHandler(api.NewGetBalanceHandler(s.balanceRepo, s.transactions).Handle))

			router.Route("/withdraw", func(router chi.Router) {
				router.Post("/", s.authHandler(api.NewWithdrawHandler(s.balanceRepo, s.ordersRepo, s.withdrawRepo, s.transactions).Handle))
			})
		})

		router.Route("/withdrawals", func(router chi.Router) {
			router.Get("/", s.authHandler(api.NewWithdrawalsHandler(s.withdrawRepo, s.transactions).Handle))
		})

	})

	return router
}

func (s *Server) authHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			slog.WarnContext(r.Context(), "auth handler empty")
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, err := s.parseToken(r.Context(), authHeader)
		if err != nil {
			slog.ErrorContext(r.Context(), "auth handler",
				slog.String("method", "parse tokent"),
				slog.String("error", err.Error()))
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), utils.ContextKey("login"), claims.Login)
		r = r.WithContext(ctx)

		next(w, r)
	}
}

func (s *Server) parseToken(ctx context.Context, tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %w", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return nil, fmt.Errorf("token is expired")
	}

	slog.DebugContext(ctx, "parse token", slog.Any("claims", claims))

	return claims, nil
}

func (s *Server) signToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &models.UserData{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			slog.ErrorContext(r.Context(), "sign token",
				slog.String("method", "read body"),
				slog.String("error", err.Error()))
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err = json.Unmarshal(body, &data); err != nil {
			log.Println("unmarshal body", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		if data.Login == "" {
			log.Println("empty login")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		slog.DebugContext(r.Context(), "signing token", slog.Any("data", data))

		r.Body.Close()

		newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(time.Now().Add(30 * time.Minute)),
			},
			Login: data.Login,
		})
		tokenString, err := newToken.SignedString([]byte(s.secretKey))
		if err != nil {
			return
		}

		ctx := context.WithValue(r.Context(), utils.ContextKey("token"), tokenString)
		ctx = context.WithValue(ctx, utils.ContextKey("login"), data.Login)
		ctx = context.WithValue(ctx, utils.ContextKey("pass"), data.Pass)
		r = r.WithContext(ctx)

		slog.DebugContext(r.Context(), "sign token", slog.String("token string", tokenString))
		next(w, r)
	}
}

func (s *Server) withlogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		newUUID := uuid.New()
		ctx := context.WithValue(r.Context(), utils.ContextKey("processID"), newUUID)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

// middleware эндпоинтов для компрессии
func (s *Server) withGZipEncode(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Проверка хедеров
		headers := strings.Split(r.Header.Get("Accept-Encoding"), ",")
		if !utils.ArrayContains(headers, "gzip") {
			next.ServeHTTP(w, r)
			return
		}

		gz, err := gzip.NewWriterLevel(w, gzip.BestCompression)
		if err != nil {
			slog.ErrorContext(r.Context(), "gZip encode", slog.String("error", err.Error()))
		}
		defer gz.Close()

		slog.DebugContext(r.Context(), "compressing request with gzip")

		w.Header().Set("Content-Encoding", "gzip")

		next.ServeHTTP(utils.GzipWriter{ResponseWriter: w, Writer: gz}, r)
	})
}
