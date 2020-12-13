package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/gofrs/uuid"
)

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
}

func signin(w http.ResponseWriter, r *http.Request) {
	var creds credentials

	// Decode JSON request body
	if err := json.NewDecoder(r.Body).Decode(&creds); err != nil {
		// If the body is invalid, error
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	// Get the expected password from the database
	expectedPassword, ok := users[creds.Username]

	if !ok || expectedPassword != creds.Password {
		w.WriteHeader(http.StatusUnauthorized)
	}

	// Create a new session token
	uuid, _ := uuid.NewV4()
	sessionToken := uuid.String()

	q := cache.SetEX(ctx, sessionToken, creds.Username, 120*time.Second)

	if err := q.Err(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Now, set the client cookie
	http.SetCookie(w, &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken,
		Expires: time.Now().Add(120 * time.Second),
	})

	w.Write([]byte("welcome"))
}

func welcome(w http.ResponseWriter, r *http.Request) {
	c, err := r.Cookie("session_token")

	if err != nil {
		if err == http.ErrNoCookie {
			// If the cookie doesn't exist, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		// For any other error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	sessionToken := c.Value

	// Get name of user from cache
	get, err := cache.Get(ctx, sessionToken).Result()

	if err != nil {
		if err == redis.Nil {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Write([]byte(
		fmt.Sprintf("Welcome %s!", get),
	))
}
