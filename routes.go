package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"net/http"

	"crypto/rand"

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
		var p, s, salt *big.Int
		var err error

		// maybe make the public string a uuid?
		if p, err = rand.Prime(rand.Reader, 64); err != nil {
			log.Fatal(err)
		}

		if s, err = rand.Prime(rand.Reader, 64); err != nil {
			log.Fatal(err)
		}

		if salt, err = rand.Prime(rand.Reader, 32); err != nil {
			log.Fatal(err)
		}

		secret := sha3.Sum256(append(s.Bytes()[:], salt.Bytes()[:]...))

		if err != nil {
			log.Fatal(err)
		}

		// goes to user, this should be compounded into a token format like: v1.public.secret.salt
		x := map[string]interface{}{
			"public": p.String(),                                   // public string
			"secret": base64.URLEncoding.EncodeToString(s.Bytes()), // secret before hashing & salting
		}

		// goes to database
		y := map[string]interface{}{
			"public": p.String(),                                   // public string, primary key for session entries
			"secret": base64.URLEncoding.EncodeToString(secret[:]), // secret after hashing & salting
		}

		fmt.Println(x)
		fmt.Println(y)
		w.Write([]byte(fmt.Sprintf("v1.%s.%s", p.String(), base64.URLEncoding.EncodeToString(s.Bytes()))))
	} else {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("password is not in db"))
	}
}

// welcome is an example of a function requiring an authenticated user
func welcome(w http.ResponseWriter, r *http.Request) {
	secret := r.Header.Get("Auth-Token")
	public := r.Header.Get("Auth-Token-Identifier")

	w.Write([]byte(authHeader))

	// compare authHeader with a token in the database
}
