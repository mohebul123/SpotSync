package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// JWTMiddleware checks for a valid bearer token and extracts user data
func JWTMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if authHeader == "" {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"success": false,
				"message": "Missing authorization token",
			})
		}

		// Extract token from "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"success": false,
				"message": "Invalid authorization format",
			})
		}

		tokenString := parts[1]
		secret := []byte(os.Getenv("JWT_SECRET"))

		// Parse and validate token
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, echo.ErrUnauthorized
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"success": false,
				"message": "Invalid or expired token",
			})
		}

		// Extract claims and inject into Echo Context
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.JSON(http.StatusUnauthorized, echo.Map{
				"success": false,
				"message": "Invalid token claims",
			})
		}

		// Store userID and role as float64/string in context safely
		c.Set("userID", uint(claims["user_id"].(float64)))
		c.Set("role", claims["role"].(string))

		return next(c)
	}
}

// RequireAdmin middleware ensures the requester has the 'admin' role
func RequireAdmin(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		role, ok := c.Get("role").(string)
		if !ok || role != "admin" {
			return c.JSON(http.StatusForbidden, echo.Map{
				"success": false,
				"message": "Access denied: Admin role required",
			})
		}
		return next(c)
	}
}
