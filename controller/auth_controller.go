package controller

import (
	"errors"
	"log"

	"github.com/Valeron93/myfiles/controller/schema"
	"github.com/Valeron93/myfiles/middleware"
	"github.com/Valeron93/myfiles/service/auth"
	"github.com/Valeron93/myfiles/validation"
	"github.com/gofiber/fiber/v2"
)

type AuthController struct {
	auth auth.Auth
}

func NewAuth(auth auth.Auth) *AuthController {
	return &AuthController{
		auth,
	}
}

func (a *AuthController) RegisterPage(c *fiber.Ctx) error {
	return c.Render("register", nil, "layout")
}

func (a *AuthController) Register(c *fiber.Ctx) error {
	registerSchema := schema.RegisterUser{}
	if err := c.BodyParser(&registerSchema); err != nil {
		return err
	}

	_, err := a.auth.RegisterUser(
		c.Context(),
		auth.UserRegistrationInfo(registerSchema),
	)

	if err != nil {
		return c.Render("register-form", fiber.Map{
			"Form":   registerSchema,
			"Errors": translateAuthErrors(err),
		})
	}

	session, err := a.auth.CreateSession(
		c.Context(),
		registerSchema.Username,
		registerSchema.Password,
	)

	if err != nil {
		return c.Render("register-form", fiber.Map{
			"Form":   registerSchema,
			"Errors": translateAuthErrors(err),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     middleware.SessionCookieKey,
		Value:    session.ID.String(),
		Path:     "/",
		HTTPOnly: true,
	})

	c.Set("HX-Redirect", "/")
	return nil
}

func (a *AuthController) LoginPage(c *fiber.Ctx) error {
	return c.Render("login", nil, "layout")
}

func (a *AuthController) Login(c *fiber.Ctx) error {
	loginSchema := schema.LoginUser{}
	if err := c.BodyParser(&loginSchema); err != nil {
		return err
	}

	session, err :=
		a.auth.CreateSession(c.Context(), loginSchema.Username, loginSchema.Password)

	if err != nil {
		return c.Render("login-form", fiber.Map{
			"Form":   loginSchema,
			"Errors": translateAuthErrors(err),
		})
	}

	c.Cookie(&fiber.Cookie{
		Name:     middleware.SessionCookieKey,
		Value:    session.ID.String(),
		Path:     "/",
		HTTPOnly: true,
	})

	c.Set("HX-Redirect", "/")
	return nil
}

// translates known errors to string that can be displayed in html
// if error is unknown, returns internal server error
func translateAuthErrors(err error) []string {
	var validateError validation.Error
	if errors.As(err, &validateError) {
		return validateError.Messages
	}

	if errors.Is(err, auth.ErrUserAlreadyExists) {
		return []string{"User with this username already exists"}
	}

	if errors.Is(err, auth.ErrInvalidCredentials) {
		return []string{"Invalid credentials"}
	}
	log.Printf("unhandled error: %v", err)
	return []string{"Internal Server Error"}
}
