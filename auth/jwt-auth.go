package auth

import (
	"net/http"
	"strings"

	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// JWT
var jwtKey = []byte("my_secret_key")
var tokens []string

type Claims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func JwtTokenGen(ctx *gin.Context) {
	token, _ := generateJWT()
	tokens = append(tokens, token)

	ctx.JSON(http.StatusOK, gin.H{
		"token": token,
	})
}

func JwtTokenCheck(ctx *gin.Context) {

	bearerToken := ctx.Request.Header.Get("Authorization")
	reqToken := strings.Split(bearerToken, " ")[1]
	claims := &Claims{}
	tkn, err := jwt.ParseWithClaims(reqToken, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			ctx.JSON(http.StatusUnauthorized, gin.H{
				"message": "Unauthorized: Access Denied",
			})
			return
		}
		ctx.JSON(http.StatusBadRequest, gin.H{
			"message": "Bad Request",
		})
		return
	}
	if !tkn.Valid {
		ctx.JSON(http.StatusUnauthorized, gin.H{
			"message": "Unauthorized: Access Denied",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

func generateJWT() (string, error) {
	expirationTime := time.Now().Add(5 * time.Minute)
	claims := &Claims{
		Username: "username",
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(jwtKey)
}
