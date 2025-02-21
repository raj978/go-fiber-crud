package user

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

// UserHandler handles user-related HTTP requests.
type UserHandler struct {
	userService UserService
}

// NewUserHandler registers the user routes and attaches the JWT middleware to protect them.
func NewUserHandler(userRoute fiber.Router, userService UserService) {
	handler := &UserHandler{
		userService: userService,
	}

	// Unprotected routes
	userRoute.Post("", handler.createUser)

	// Wrap all routes in a JWT-protected group
	// protected := userRoute.Group("/", JWTMiddleware)

	// protected.Get("", handler.getUsers)
	// protected.Get("/:userID", handler.JWTMiddleware, handler.getUser)
	// protected.Put("/:userID", handler.checkIfUserExistsMiddleware, handler.JWTMiddleware, handler.updateUser)
	// protected.Delete("/:userID", handler.checkIfUserExistsMiddleware, handler.JWTMiddleware, handler.deleteUser)

	userRoute.Get("", handler.getUsers)
	userRoute.Get("/:userID", JWTMiddleware, handler.getUser)
	userRoute.Put("/:userID", handler.checkIfUserExistsMiddleware, JWTMiddleware, handler.updateUser)
	userRoute.Delete("/:userID", handler.checkIfUserExistsMiddleware, JWTMiddleware, handler.deleteUser)
}

// Gets all users.
func (handler *UserHandler) getUsers(ctx *fiber.Ctx) error {
	// Create cancellable context.
	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Get all users.
	users, err := handler.userService.GetUsers(customContext)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status": "success",
		"data":   users,
	})
}

// Gets a single user.
func (handler *UserHandler) getUser(ctx *fiber.Ctx) error {
	// Create cancellable context.
	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Fetch parameter.
	targetedUserID := ctx.Params("userID")

	// Get one user.
	user, err := handler.userService.GetUser(customContext, targetedUserID)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	// Return results.
	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status": "success",
		"data":   user,
	})
}

// Creates a single user.
func (handler *UserHandler) createUser(ctx *fiber.Ctx) error {
	// Initialize variables.
	user := &User{}

	// Create cancellable context.
	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Parse request body.
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	// Create one user.
	createdUser, err := handler.userService.CreateUser(customContext, user)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	// Return result.
	return ctx.Status(fiber.StatusCreated).JSON(&fiber.Map{
		"status":  "success",
		"message": "User has been created successfully!",
		"data":    createdUser,
	})
}

// Updates a single user.
func (handler *UserHandler) updateUser(ctx *fiber.Ctx) error {
	// Initialize variables.
	user := &User{}
	targetedUserID := ctx.Params("userID")

	// Create cancellable context.
	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Parse request body.
	if err := ctx.BodyParser(user); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	// Update one user.
	updatedUser, err := handler.userService.UpdateUser(customContext, targetedUserID, user)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	// Return result.
	return ctx.Status(fiber.StatusOK).JSON(&fiber.Map{
		"status":  "success",
		"message": "User has been updated successfully!",
		"data":    updatedUser,
	})
}

// Deletes a single user.
func (handler *UserHandler) deleteUser(ctx *fiber.Ctx) error {
	// Initialize previous user ID.
	targetedUserID := ctx.Params("userID")

	// Create cancellable context.
	customContext, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Delete one user.
	err := handler.userService.DeleteUser(customContext, targetedUserID)

	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(&fiber.Map{
			"status":  "fail",
			"message": err.Error(),
		})
	}

	// Return 204 No Content.
	return ctx.SendStatus(fiber.StatusNoContent)
}
