package service

import (
	"context"
	"crypto/rsa"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/MFCaballero/pangea/appdomain"
	"github.com/dgrijalva/jwt-go"
)

const (
	privKeyPath = "keys/app.rsa"     // openssl genrsa -out app.rsa keysize
	pubKeyPath  = "keys/app.rsa.pub" // openssl rsa -in app.rsa -pubout > app.rsa.pub
	TokenName   = "x-access-token"
)

var (
	verifyKey *rsa.PublicKey
	signKey   *rsa.PrivateKey
)

type Token struct {
	UserID uint
	Name   string
	Email  string
	*jwt.StandardClaims
}

type JwtAuthenticator struct {
}

func fatal(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func initKeys() {
	dir, _ := os.Getwd()
	dir += "/configs"

	signBytes, err := ioutil.ReadFile(dir + privKeyPath)
	fatal(err)

	signKey, err = jwt.ParseRSAPrivateKeyFromPEM(signBytes)
	fatal(err)

	verifyBytes, err := ioutil.ReadFile(dir + pubKeyPath)
	fatal(err)

	verifyKey, err = jwt.ParseRSAPublicKeyFromPEM(verifyBytes)
	fatal(err)
}

var authenticator *JwtAuthenticator

func GetAuthenticator() *JwtAuthenticator {
	if authenticator == nil {
		authenticator = &JwtAuthenticator{}
		initKeys()
	}
	return authenticator
}

func (jwtAuth JwtAuthenticator) IsTokenExists(r *http.Request) (bool, string) {

	var token string = ""

	if token = r.Header.Get(TokenName); token != "" {
		return true, token

	} else if cookie, err := r.Cookie(TokenName); err == nil {
		token = cookie.Value
		return true, token

	} else if keys, ok := r.URL.Query()[TokenName]; ok {

		token = keys[0]
		if token != "" {
			return true, token
		}
	}
	return false, token

}

func (jwtAuth JwtAuthenticator) IsUserTokenValid(token string) bool {

	tk := &Token{}

	_, err := jwt.ParseWithClaims(token, tk, func(token *jwt.Token) (interface{}, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		return verifyKey, nil
	})
	return err == nil
}

func (jwtAuth JwtAuthenticator) UserFromToken(tokenString string) (*appdomain.User, error) {

	tk := &Token{}

	_, err := jwt.ParseWithClaims(tokenString, tk, func(token *jwt.Token) (interface{}, error) {
		// since we only use the one private key to sign the tokens,
		// we also only use its public counter part to verify
		return verifyKey, nil
	})
	if err != nil {
		return nil, err
	}

	//this is for simplicity, we could also just have the id of the user in the token and fetch the rest from the db
	var usr = appdomain.User{
		Email: tk.Email,
		Name:  tk.Name,
		ID:    tk.UserID,
	}
	return &usr, err
}

func (jwtAuth JwtAuthenticator) JwtVerify(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		exist, token := jwtAuth.IsTokenExists(r) //Grab the token from the header
		if !exist {
			w.WriteHeader(http.StatusForbidden)
			return
		}
		usr, err := jwtAuth.UserFromToken(token)

		if err != nil {
			w.WriteHeader(http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), "user", usr)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (jwtAuth JwtAuthenticator) GetTokenForUser(user *appdomain.User) (string, error) {

	expiresAt := time.Now().Add(time.Minute * 100000).Unix()

	tk := &Token{
		UserID: user.ID,
		Name:   user.Name,
		Email:  user.Email,
		StandardClaims: &jwt.StandardClaims{
			ExpiresAt: expiresAt,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, tk)

	// create a signer for rsa 256
	tokenString, err := token.SignedString(signKey)
	if err != nil {
		return "", errors.New("message Sorry, error while Signing Token")
	}

	return tokenString, nil
}
