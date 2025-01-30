package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

const secret = "testsecret"

func generateJWT(valid bool) string {
	claims := jwt.MapClaims{
		"id":   1,
		"name": "Test User",
		"exp":  time.Now().Add(time.Hour).Unix(),
	}

	if !valid {
		claims["exp"] = time.Now().Add(-time.Hour).Unix()
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := token.SignedString([]byte(secret))
	return signedToken
}

func TestJWTMiddleware(t *testing.T) {
	app := fiber.New()
	app.Use(JWTMiddleware(secret))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Success"})
	})

	t.Run("should return unauthorized if token is missing", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return unauthorized if token format is incorrect", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "InvalidTokenFormat")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should return unauthorized if token is invalid", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer invalid.token.value")

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should pass if token is valid", func(t *testing.T) {
		validToken := generateJWT(true)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		req.Header.Set("Authorization", "Bearer "+validToken)

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})
}

func TestJWTMiddlewareParam(t *testing.T) {
	app := fiber.New()
	app.Use(JWTMiddlewareParam(secret))
	app.Get("/", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"message": "Success"})
	})

	t.Run("should return unauthorized if token is missing in query", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusAlreadyReported, resp.StatusCode)
	})

	t.Run("should return unauthorized if token is invalid in query", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/?token=invalid.token.value", nil)

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("should pass if token is valid in query", func(t *testing.T) {
		validToken := generateJWT(true)
		req := httptest.NewRequest(http.MethodGet, "/?token="+validToken, nil)

		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
	})
}
