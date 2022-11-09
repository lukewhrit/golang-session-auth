package main

import (
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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

		secret := sha3.Sum256(append(s.Bytes()[:], salt.Bytes()[:]...))

		if err != nil {
			log.Fatal(err)
		}

		if err := rdb.Set(ctx, p.String(), fmt.Sprintf("v1.%s.%s.%s", p.String(),
			base64.URLEncoding.EncodeToString(secret[:]), salt.String()), 0).Err(); err != nil {
			log.Fatal(err)
		}

		w.Write([]byte(fmt.Sprintf("v1.%s.%s", p.String(), base64.URLEncoding.EncodeToString(s.Bytes()))))
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("password is not in db"))
	}
}

// welcome is an example of a function requiring an authenticated user
func welcome(w http.ResponseWriter, r *http.Request) {
	token, err := parseToken(r.Header.Get("Auth-Token"))

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("problem parsing token. is it valid?"))
	}

	resp, err := rdb.Get(ctx, token.Public).Bytes()

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	dbToken, err := parseToken(string(resp))

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	}

	// todo UN BASE 64

	saltedIncomingToken := sha3.Sum256(append([]byte(token.Secret), []byte(dbToken.Salt)...))

	fmt.Println(saltedIncomingToken)

	// compare authHeader with a token in the database
	fmt.Println(subtle.ConstantTimeCompare(saltedIncomingToken[:], []byte(dbToken.Secret)))
}
