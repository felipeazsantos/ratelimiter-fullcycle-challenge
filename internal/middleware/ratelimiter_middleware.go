package middleware

import (
	"net/http"

	usecase "github.com/felipeazsantos/ratelimiter-fullcycle-challenge/internal/usecase/limiter"
)

func RateLimiterMiddleware(useCase *usecase.RateLimiterUseCase, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allow, err := useCase.Execute(w, r)

		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		if !allow {
			http.Error(w, "you have reached the maximum number of requests or actions allowed within a certain time frame", http.StatusMethodNotAllowed)
			return
		}

		next.ServeHTTP(w, r)
	})
}
