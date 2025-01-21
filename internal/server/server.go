package server

import (
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/mikhaylov123ty/go-diploma-5.6/internal/utils"
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

	"github.com/go-chi/chi/v5"
	"github.com/golang-jwt/jwt/v4"
)

const (
	userHandlerPath = "/api/user"
)

type Server struct {
	address      string
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
	usersRepo users.Storage,
	ordersRepo orders.Storage,
	balanceRepo balance.Storage,
	witdrawRepo withdrawals.Storage,
	secretKey string) *Server {
	return &Server{
		address:      address,
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
			router.Post("/", s.signToken(api.NewRegisterHandler(s.usersRepo).Handle))
		})

		router.Route("/login", func(router chi.Router) {
			router.Post("/", s.signToken(api.NewAuthHandler(s.usersRepo).Handle))
		})

		router.Route("/orders", func(router chi.Router) {
			router.Post("/", s.authHandler(api.NewPostOrdersHandler(s.ordersRepo, s.usersRepo).Handle))
			router.Get("/", s.authHandler(api.NewGetOrdersHandler(s.ordersRepo, s.usersRepo).Handle))
		})

		router.Route("/balance", func(router chi.Router) {
			router.Get("/", s.authHandler(api.NewGetBalanceHandler(s.balanceRepo).Handle))

			router.Route("/withdraw", func(router chi.Router) {
				router.Post("/", s.authHandler(api.NewWithdrawHandler(s.balanceRepo, s.ordersRepo, s.withdrawRepo).Handle))
			})
		})

		router.Route("/withdrawals", func(router chi.Router) {
			router.Get("/", s.authHandler(api.NewWithdrawalsHandler(s.withdrawRepo).Handle))
		})

	})

	return router
}

func (s *Server) authHandler(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, err := s.parseToken(authHeader)
		if err != nil {
			log.Println(err)
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		ctx := context.WithValue(r.Context(), utils.ContextKey("login"), claims.Login)
		r = r.WithContext(ctx)

		next(w, r)
	}
}

func (s *Server) parseToken(tokenString string) (*Claims, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	if claims.ExpiresAt.Before(time.Now()) {
		return nil, errors.New("token expired")
	}

	fmt.Println("Claims", claims)
	return claims, nil
}

func (s *Server) signToken(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		data := &models.UserData{}
		body, err := io.ReadAll(r.Body)
		if err != nil {
			log.Println("read body error", err)
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
