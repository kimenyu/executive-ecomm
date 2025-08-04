package auth

import (
	"context"
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/kimenyu/executive/configs"
	"github.com/kimenyu/executive/types"
	"github.com/kimenyu/executive/utils"
	"log"
	"net/http"
	"time"
)

type contextKey string

const UserKey contextKey = "userID"

// Middleware
func WithJWTAuth(store types.UserStore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := utils.GetTokenFromRequest(r)

			token, err := validateJWT(tokenString)
			if err != nil || !token.Valid {
				log.Printf("invalid token: %v", err)
				permissionDenied(w)
				return
			}

			claims, ok := token.Claims.(jwt.MapClaims)
			if !ok {
				log.Println("invalid claims")
				permissionDenied(w)
				return
			}

			// Extract and parse UUID
			str, ok := claims["userID"].(string)
			if !ok {
				log.Println("userID claim not found or invalid")
				permissionDenied(w)
				return
			}

			userUUID, err := uuid.Parse(str)
			if err != nil {
				log.Printf("failed to parse UUID: %v", err)
				permissionDenied(w)
				return
			}

			// Optional DB lookup (for verification)
			_, err = store.GetUserByID(userUUID)
			if err != nil {
				log.Printf("user not found: %v", err)
				permissionDenied(w)
				return
			}

			// Add UUID to context
			ctx := context.WithValue(r.Context(), UserKey, userUUID)
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
	return token.SignedString(secret)
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

// Get user UUID from context
func GetUserIDFromContext(ctx context.Context) uuid.UUID {
	userID, ok := ctx.Value(UserKey).(uuid.UUID)
	if !ok {
		return uuid.Nil
	}
	return userID
}
