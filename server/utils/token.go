package token

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var JWT_SECRET string
var JWT_EXPIRATION int

func Generate(username string) (string, error) {
	claims := jwt.MapClaims{
		"authorized": true,
		"username":   username,
		"exp":        time.Now().Add(time.Hour * time.Duration(JWT_EXPIRATION)).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString([]byte(JWT_SECRET))
}

func Validate(c *gin.Context) error {
	tokenString := ExtractToken(c)

	_, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(JWT_SECRET), nil
	})

	if err != nil {
		return err
	}

	return nil
}

func ExtractToken(c *gin.Context) string {
	token := c.Query("token")
	if token != "" {
		return token
	}

	authHeader := c.Request.Header.Get("Authorization")
	var bearerToken = strings.Split(authHeader, " ")

	if len(bearerToken) == 2 {
		return strings.Trim(bearerToken[1], "\" ")
	}

	return ""
}

func ExtractTokenUsername(c *gin.Context) (string, error) {
	tokenString := ExtractToken(c)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(JWT_SECRET), nil
	})

	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)

	if ok && token.Valid {
		username := claims["username"].(string)
		if err != nil {
			return "", err
		}
		return username, nil
	}

	return "", nil
}

func JwtMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := Validate(c)
		if err != nil {
			c.String(http.StatusUnauthorized, err.Error())
			c.Abort()
			return
		}
		c.Next()
	}
}

func SetEnvs() {
	// get values from .env
	env_secret := os.Getenv("JWT_SECRET")
	env_expiration := os.Getenv("JWT_EXPIRATION_HOURS")

	// convert expiration to int
	result, err := strconv.Atoi(env_expiration)

	if err != nil {
		panic(err)
	}

	JWT_EXPIRATION = result
	JWT_SECRET = env_secret
}
