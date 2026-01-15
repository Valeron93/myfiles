package controller

import (
	"context"
	"errors"
	"log"

	"github.com/Valeron93/myfiles/controller/schema"
	"github.com/Valeron93/myfiles/service"
	"github.com/Valeron93/myfiles/validation"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const SessionCookieKey = "session_id"

type AuthController struct {
	auth service.Auth
}

type SessionCtx struct{}

func NewAuth(auth service.Auth) *AuthController {
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
		service.UserRegistrationInfo(registerSchema),
	)

	if err != nil {
		var validationError validation.Error
		if errors.As(err, &validationError) {
			return c.Render("register-form", fiber.Map{
				"Form":   registerSchema,
				"Errors": validationError.Messages,
			})
		}
		log.Printf("%#+v", err)
		if errors.Is(err, service.ErrAlreadyExists) {
			return c.Render("register-form", fiber.Map{
				"Form":   registerSchema,
				"Errors": []string{"User with this username already exists."},
			})
		}

		return err
	}

	session, err := a.auth.CreateSession(
		c.Context(),
		registerSchema.Username,
		registerSchema.Password,
	)

	if err != nil {
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     SessionCookieKey,
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
		if errors.Is(err, service.ErrInvalidCredentials) {
			return c.Render("login-form", fiber.Map{
				"Form":   loginSchema,
				"Errors": []string{"Invalid credentials."},
			})
		}
		return err
	}

	c.Cookie(&fiber.Cookie{
		Name:     SessionCookieKey,
		Value:    session.ID.String(),
		Path:     "/",
		HTTPOnly: true,
	})

	c.Set("HX-Redirect", "/")
	return nil
}

func (a *AuthController) InjectSession() func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		sessionID := c.Cookies(SessionCookieKey)
		if sessionID == "" {
			return c.Next()
		}

		sessionUUID, err := uuid.Parse(sessionID)
		if err != nil {
			return c.Next()
		}

		session, err := a.auth.GetSession(
			c.Context(),
			sessionUUID,
		)

		if err != nil {
			if errors.Is(err, service.ErrNotFound) {
				c.ClearCookie(SessionCookieKey)
			}
			return c.Next()
		}

		c.SetUserContext(context.WithValue(
			context.Background(),
			SessionCtx{},
			session,
		))

		return c.Next()
	}
}
