package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"
	"golang.org/x/crypto/sha3"
)

var users = map[string]string{
	"user1": hashAndSalt([]byte("password1")),
	"user2": "password2",
}

// signin takes a username and password and gives back a session token
func signin(w http.ResponseWriter, r *http.Request) {
	var body credentials

	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		log.Print(err)
	}

	if bcrypt.CompareHashAndPassword([]byte(users[body.Username]), []byte(body.Password)) == nil {
		p, s, salt, err := generateStrings([]int{64, 64, 32})

		if err != nil {
			log.Fatal(err)
		}

		buf := []byte(s + salt)
		secret := make([]byte, 64)
		sha3.ShakeSum256(secret, buf)

		userToken := makeToken(Token{
			Version: "v1",
			Public:  p,
			Secret:  base64.URLEncoding.EncodeToString([]byte(s)),
			Salt:    salt,
		})

		serverToken := makeToken(Token{
			Version: "v1",
			Public:  p,
			Secret:  fmt.Sprintf("%x", secret),
			Salt:    salt,
		})

		if err != nil {
			log.Fatal(err)
		}

		if err := rdb.Set(ctx, p, serverToken, 0).Err(); err != nil {
			log.Fatal(err)
		}

		w.Write([]byte(userToken))
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("password is not in db"))
	}
}

// welcome is an example of a function requiring an authenticated user
func welcome(w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.Header.Get("Auth-Token"), "v1") {
		token, err := parseToken(r.Header.Get("Auth-Token"))

		if err != nil {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("parsing auth token errored. is it valid?"))
		}

		entry := rdb.Get(ctx, token.Public).String()

		if entry != "" {
			entryToken, err := parseToken(entry)

			if err != nil {
				w.WriteHeader(http.StatusBadGateway)
				w.Write([]byte("parsing auth token in database errored. is it valid?"))
			}

			fmt.Println("E", entryToken.Secret)

			unb64dSecret, err := base64.URLEncoding.DecodeString(entryToken.Secret)

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("problem decoding secret from base64"))
				w.Write([]byte(err.Error()))
			}

			buf := []byte(string(unb64dSecret) + token.Salt)
			secret := make([]byte, 64)
			sha3.ShakeSum256(secret, buf)

			fmt.Printf("%x\n", secret)
		} else {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("entry is nil"))
		}
	}
}
