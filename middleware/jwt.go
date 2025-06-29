package middleware

import (
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt"
	"github.com/pdrm26/hotel-reservation/core"
	"github.com/pdrm26/hotel-reservation/db"
)

type Claim struct {
	Id        string `json:"id"`
	Email     string `json:"email"`
	ValidTill string `json:"validtill"`
}

func JWTAuthentication(userStore db.UserStore) fiber.Handler {
	return func(c *fiber.Ctx) error {
		token := c.Get("X-Api-Token")
		if token == "" {
			return core.TokenMissingError()
		}

		claims, err := validateToken(token)
		if err != nil {
			return core.TokenInvalidError()
		}

		expires, exists := claims["expires"]
		if !exists {
			return core.TokenInvalidError()
		}

		expiresStr, ok := expires.(string)
		if !ok {
			return core.TokenInvalidError()
		}

		expirationTime, err := time.Parse(time.RFC3339, expiresStr)
		if err != nil {
			return core.TokenInvalidError()
		}

		if time.Now().After(expirationTime) {
			return core.TokenExpiredError()
		}

		userID := claims["id"].(string)
		user, err := userStore.GetUserByID(c.Context(), userID)
		if err != nil {
			return core.NotFoundError("user")
		}

		c.Context().SetUserValue("user", user)

		return c.Next()
	}
}

func validateToken(tokenStr string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			fmt.Println("invalid signin method", t.Header["alg"])
			return nil, core.UnAuthorizedError()
		}
		secret := os.Getenv("JWT_SECRET")
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("failed to parse JWT token:", err)
		return nil, core.TokenInvalidError()
	}

	if !token.Valid {
		fmt.Println("invalid token")
		return nil, core.TokenInvalidError()
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, core.UnAuthorizedError()
	}

	return claims, nil
}
