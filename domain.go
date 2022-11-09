package main

import (
	"crypto/rand"
	"fmt"
	"log"
	"math/big"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type Token struct {
	Version string
	Public  string
	Secret  string
	Salt    string
}

type credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

func hashAndSalt(pwd []byte) string {
	// Use GenerateFromPassword to hash & salt pwd.
	// MinCost is just an integer constant provided by the bcrypt
	// package along with DefaultCost & MaxCost.
	// The cost can be any value you want provided it isn't lower
	// than the MinCost (4)
	hash, err := bcrypt.GenerateFromPassword(pwd, bcrypt.MinCost)
	if err != nil {
		log.Println(err)
	} // GenerateFromPassword returns a byte slice so we need to
	// convert the bytes to a string and return it
	return string(hash)
}

func generateStrings(bits []int) (a, b, c *big.Int, err error) {
	if a, err = rand.Prime(rand.Reader, bits[0]); err != nil {
		return nil, nil, nil, err
	}

	if b, err = rand.Prime(rand.Reader, bits[1]); err != nil {
		return nil, nil, nil, err
	}

	if c, err = rand.Prime(rand.Reader, bits[2]); err != nil {
		return nil, nil, nil, err
	}

	return a, b, c, err
}

func parseToken(token string) (Token, error) {
	var tok Token
	toks := strings.Split(token, ".")

	fmt.Println(toks[0])

	tok.Version = toks[0]
	tok.Public = toks[1]
	tok.Secret = toks[2]

	if len(toks) == 4 {
		tok.Salt = toks[3]
	}

	return tok, nil
}
