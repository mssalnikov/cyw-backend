package utils

import (
	"net/http"
	"log"
	"github.com/go-redis/redis"
	"fmt"
)

func AuthenticationMiddleware(next http.Handler) http.Handler {
	log.Println(next)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Session-Token")
		log.Println(r.URL.String())
		switch url := r.URL.String(); url {
		case "/":
			next.ServeHTTP(w, r)
		case "/auth":
			next.ServeHTTP(w, r)
		default:
			userId, err := RedisCon.Get(fmt.Sprintf("TOKEN:%s", token)).Int64()
			if err == redis.Nil {
				log.Println("Unauthorized request at: ", r.URL)
				http.Error(w, "Forbidden", http.StatusForbidden)
			} else if err != nil {
				log.Println(err)
				http.Error(w, "Bad token", http.StatusBadRequest)
			} else {
				log.Printf("Authenticated user %d\n", userId)
				// Pass down the request to the next middleware (or final handler)
				next.ServeHTTP(w, r)
			}
		}
	})
	return nil
}
