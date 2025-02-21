package user

import (
	"context"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

// If user does not exist, do not allow one to access the API.
func (handler *UserHandler) checkIfUserExistsMiddleware(ctx *fiber.Ctx) error {
	// Create a new customized context.
	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Fetch parameter.
	targetedUserID := ctx.Params("userID")

	if targetedUserID == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": "User not found!",
		})
	}

	// Check if user exists.
	searchedUser, err := handler.userService.GetUser(customContext, targetedUserID)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	if searchedUser == nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": "User not found!",
		})
	}

	return ctx.Next()
}

// JWTMiddleware verifies JWT tokens provided in the Authorization header using HS256 algorithm.
func JWTMiddleware(c *fiber.Ctx) error {
	// Extract the Authorization header
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
	}

	// Expected header format: "Bearer <token>"
	tokenParts := strings.Split(authHeader, " ")
	if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Authorization header format"})
	}

	tokenString := tokenParts[1]

	// Retrieve the JWT secret from environment variables
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "JWT secret not configured"})
	}

	// Parse and validate the JWT token using HS256 method
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure that the token method conforms to HS256
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fiber.NewError(fiber.StatusUnauthorized, "Unexpected signing method")
		}
		return []byte(secret), nil
	})

	if err != nil || token == nil || !token.Valid {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
	}

	// Token is valid; proceed to the next handler
	return c.Next()
}

// JWTMiddlewareWithUserService verifies JWT tokens and ensures that the token matches the one stored in the corresponding user record.
// This middleware requires a UserService to fetch the user from the database using an identifier from the token claims.
// func JWTMiddlewareWithUserService(userService UserService) fiber.Handler {
// 	return func(c *fiber.Ctx) error {
// 		// Extract the Authorization header
// 		authHeader := c.Get("Authorization")
// 		if authHeader == "" {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Missing Authorization header"})
// 		}

// 		// Expected header format: "Bearer <token>"
// 		tokenParts := strings.Split(authHeader, " ")
// 		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid Authorization header format"})
// 		}

// 		tokenString := tokenParts[1]

// 		// Retrieve the JWT secret from environment variables
// 		secret := os.Getenv("JWT_SECRET")
// 		if secret == "" {
// 			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "JWT secret not configured"})
// 		}

// 		// Parse and validate the JWT token using HS256 method
// 		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 				return nil, fiber.NewError(fiber.StatusUnauthorized, "Unexpected signing method")
// 			}
// 			return []byte(secret), nil
// 		})

// 		if err != nil || token == nil || !token.Valid {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid or expired token"})
// 		}

// 		// Extract claims and get user identifier (assuming the claim "user" holds a unique identifier, e.g., email or user id)
// 		claims, ok := token.Claims.(jwt.MapClaims)
// 		if !ok {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token claims"})
// 		}

// 		userIdentifier, ok := claims["user"].(string)
// 		if !ok || userIdentifier == "" {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User identifier not found in token"})
// 		}

// 		// Fetch the user record from the database using the user identifier
// 		dbUser, err := userService.GetUser(c.Context(), userIdentifier)
// 		if err != nil || dbUser == nil {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
// 		}

// 		// Check if the stored auth token for the user matches the token provided
// 		if dbUser.AuthToken != tokenString {
// 			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token does not match stored user token"})
// 		}

// 		// Optionally, store the user in the context for later use
// 		c.Locals("user", dbUser)
// 		return c.Next()
// 	}
// }
