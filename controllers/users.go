package controllers

import (
	"log/slog"
	"net/http"
	"strings"
	"unicode"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/preguntame/preguntame-backend/auth"
	"github.com/preguntame/preguntame-backend/models"
)

type loginDTO struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type registerDTO struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

/**
 * Endpoint for handling the login
 * It checks that the credentials are valid and then returns a JWT with the user claims
 */
func Login(e echo.Context) error {
	params := loginDTO{}

	// Check if the body contains a loginDTO
	if err := e.Bind(&params); err != nil {
		slog.WarnContext(e.Request().Context(), "Error binding the request", "error", err)
		return err
	}

	// Search the user by the email in the database
	user, err := models.FindUserByEmail(params.Email)
	if err != nil {
		slog.Error("Error searching user in database", "error", err)
		return err
	}

	// Return an error in case that the user doesn't exists or that the password doesn't match
	if user == nil || params.Password != user.Password {
		return e.String(http.StatusBadRequest, "User and password not match")
	}

	// Generates the JWT
	jwt, err := auth.NewTokenForUser(*user)
	if err != nil {
		slog.Error("Error generating jwt", "error", err)
		return err
	}

	// If we haven't returend yet, it means that the loggin was successful
	return e.String(http.StatusOK, jwt)
}

/**
 * Endpoint for handliong new user registers
 */
func Register(e echo.Context) error {
	params := registerDTO{}

	// Check if the body contains a registerDTO
	if err := e.Bind(&params); err != nil {
		slog.WarnContext(e.Request().Context(), "Error binding the request", "error", err)
		return err
	}

	// Generate a unique ID for the new user
	uuid, err := uuid.NewUUID()
	if err != nil {
		slog.Error("Error generating id for user", "error", err)
		return err
	}

	if !strings.Contains(params.Email, "@") || len(params.Email) > 256 || !strings.Contains(params.Email, ".") {
		slog.Info("The email must contains an @, an . and must be less or equal than 256", "email", params.Email)
		return e.String(http.StatusBadRequest, "The email must contains an @, an . and must be less or equal than 256")
	}

	mayus := false
	minus := false
	num := false

	for _, r := range params.Password {
		if unicode.IsUpper(r) {
			mayus = true
		}
		if unicode.IsLower(r) {
			minus = true
		}
		if unicode.IsDigit(r) {
			num = true
		}
	}

	if !mayus || !minus || len(params.Password) < 8 || !num {
		slog.Info("The password must contains at least one uppercase character, one lowcase character, a number and the length must be greather than 8 characters")
		return e.String(http.StatusBadRequest, "The password must contains at least one uppercase character, one lowcase character a number and the length must be greather than 8 characters")
	}

	user := models.User{
		Id:       uuid.String(),
		Name:     params.Name,
		Password: params.Password,
		Email:    params.Email,
	}

	// Saves the user in the database
	err = models.InsertUser(user)
	if err != nil {
		slog.Error("Error inserting user in database", "error", err)
		return err
	}

	return e.String(http.StatusOK, "Successful register")
}
