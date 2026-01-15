package middleware

import (
	"context"
	"errors"

	"github.com/Valeron93/myfiles/service/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

const SessionCookieKey = "session_id"

type SessionCtx struct{}

func InjectSession(a auth.Auth) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		sessionID := c.Cookies(SessionCookieKey)
		if sessionID == "" {
			return c.Next()
		}

		sessionUUID, err := uuid.Parse(sessionID)
		if err != nil {
			return c.Next()
		}

		session, err := a.GetSession(
			c.Context(),
			sessionUUID,
		)

		if err != nil {
			if errors.Is(err, auth.ErrUserNotFound) {
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
