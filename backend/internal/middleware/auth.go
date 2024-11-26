package middleware

import (
    "github.com/gofiber/fiber/v2"
    "github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(secret string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        tokenString := c.Get("JWT")
        if tokenString == "" {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Missing or invalid token",
            })
        }

        token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
            return []byte(secret), nil
        })
				
        if err != nil || !token.Valid {
            return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
                "error": "Invalid token",
            })
        }

        // Add claims to context if needed
        claims := token.Claims.(jwt.MapClaims)
        c.Locals("user", claims)

        return c.Next()
    }
}
