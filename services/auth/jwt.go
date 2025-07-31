package auth

import (
	"context"
	"fmt"
	"github.com/kimenyu/executive/configs"
	"github.com/kimenyu/executive/types"
	"github.com/kimenyu/executive/utils"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserKey contextKey = "userID"

func WithJWTAuth(store types.UserStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := utils.GetTokenFromRequest(r)

			token, err := validateJWT(tokenString)
			if err != nil {
				log.Printf("failed to validate token: %v", err)
				permissionDenied(w)
				return
			}

			if !token.Valid {
				log.Println("invalid token")
				permissionDenied(w)
				return
			}

			claims := token.Claims.(jwt.MapClaims)
			str := claims["userID"].(string)

			userID, err := strconv.Atoi(str)
			if err != nil {
				log.Printf("failed to convert userID to int: %v", err)
				permissionDenied(w)
				return
			}

			u, err := store.GetUserByID(userID)
			if err != nil {
				log.Printf("failed to get user by id: %v", err)
				permissionDenied(w)
				return
			}

			// Inject userID into context
			ctx := context.WithValue(r.Context(), UserKey, u.ID)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func CreateJWT(secret []byte, userID string) (string, error) {
	expiration := time.Second * time.Duration(configs.Envs.JWTExpirationInSeconds)

	claims := jwt.MapClaims{
		"userID":    userID,
		"expiresAt": time.Now().Add(expiration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func validateJWT(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(configs.Envs.JWTSecret), nil
	})
}

func permissionDenied(w http.ResponseWriter) {
	utils.WriteError(w, http.StatusForbidden, fmt.Errorf("permission denied"))
}

func GetUserIDFromContext(ctx context.Context) int {
	userID, ok := ctx.Value(UserKey).(int)
	if !ok {
		return -1
	}

	return userID
}
