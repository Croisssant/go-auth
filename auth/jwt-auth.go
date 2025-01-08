package auth

import (
	"fmt"
	"net/http"
	"strings"

	"time"

	"croissant.com/go/auth/models"
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

func generateToken(userId uint) (string, error) {
	claims := jwt.MapClaims{}
	claims["user_id"] = userId
	claims["exp"] = time.Now().Add(time.Hour * 1).Unix()

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

func verifyToken(tokenString string) (jwt.MapClaims, error) {

	// Parsing tokenString
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {

		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return jwtKey, nil
	})

	if err != nil {
		return nil, err
	}

	// Validating token
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, fmt.Errorf("invalid token")
}

// Auth Middleware to check if user has valid JWT Token
func AuthenticationMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		tokenString := ctx.GetHeader("Authorization")
		if tokenString == "" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "missing authentication token"})
			ctx.Abort()
			return
		}

		// Token should be prefixed with "Bearer "
		tokenParts := strings.Split(tokenString, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication token, missing 'Bearer' key"})
			ctx.Abort()
			return
		}

		tokenString = tokenParts[1]
		claims, err := verifyToken(tokenString)
		if err != nil {
			ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authentication token"})
			ctx.Abort()
			return
		}

		ctx.Set("user_id", claims["user_id"])
		ctx.Next()
	}
}

// Authentication Handlers
func Login(ctx *gin.Context) {
	var user models.User

	// Check if user credentials are valid
	if err := ctx.ShouldBindJSON(&user); err != nil {
		fmt.Printf("Error: %v \n", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		ctx.Abort()
		return
	}

	if user.Username == "user" && user.Password == "password" {
		// Generate JWT Token when all checks pass
		token, err := generateToken(user.ID)
		if err != nil {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Error generating token"})
			ctx.Abort()
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"token": token})
	} else {
		fmt.Printf("User: %v, Password: %v \n", user.Username, user.Password)
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "invalid credentials"})
	}
}

func Register(ctx *gin.Context) {
	var user models.User

	if err := ctx.ShouldBindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid data"})
		ctx.Abort()
		return
	}

	user.ID = 1
	ctx.JSON(http.StatusCreated, gin.H{"message": "user registered successfully"})
}
