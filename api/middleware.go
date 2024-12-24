package api

import (
	"context"
	"log"
	"net"
	"net/http"
	"strings"
)

type Middleware func(http.Handler) http.HandlerFunc

func CorsMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		if origin != "" {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
		}

		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, Accept, X-Requested-With")
		w.Header().Set("Access-Control-Max-Age", "86400")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func LogMiddleware(next http.Handler) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Origin: %s, Method: %s, Endpoint: %s\n", r.Header.Get("Origin"), r.Method, r.URL.Path)
		next.ServeHTTP(w, r)
	})
}

func RateLimitMiddleware(next http.Handler) http.HandlerFunc {
	limiter := NewLimiter()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip, ok := r.Context().Value("ip").(string)

		if ip == "" || !ok {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if limiter.Allow(ip) == false {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getClientIP(r *http.Request) string {
	xForwardedFor := r.Header.Get("X-Forwarded-For")
	if xForwardedFor != "" {
		parts := strings.Split(xForwardedFor, ",")
		return strings.TrimSpace(parts[0])
	}

	xRealIP := r.Header.Get("X-Real-IP")
	if xRealIP != "" {
		return xRealIP
	}

	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return ""
	}
	return host
}

func IPMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ip := getClientIP(r)
		ctx := context.WithValue(r.Context(), "ip", ip)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func parseHeader(r *http.Request) string {
	split := strings.Split(r.Header.Get("Authorization"), " ")

	if len(split) != 2 || split[0] != "Bearer" {
		return ""
	}

	return split[1]
}

/*
	NOTE: Request should have the following structure

Authorization: Bearer <token>

	Body: {
	 event_type: string
	 payload?: Object
	}
*/
func AuthMiddleware(next http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenStr := parseHeader(r)

		if tokenStr == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		claims, err := parseToken(tokenStr)

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "claims", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	}
}

func ApplyMiddleware(next http.Handler, middleware ...Middleware) http.Handler {
	for i := len(middleware) - 1; i >= 0; i-- {
		next = middleware[i](next)
	}
	return next
}
