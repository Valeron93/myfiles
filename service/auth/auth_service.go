package auth

import (
	"context"
	"log"
	"time"

	"github.com/Valeron93/myfiles/model"
	"github.com/Valeron93/myfiles/validation"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserRegistrationInfo struct {
	Username        string
	Password        string
	ConfirmPassword string
}

type Auth interface {
	// Returns validation.Error if input data is invalid
	RegisterUser(ctx context.Context, opts UserRegistrationInfo) (model.User, error)
	DeleteUser(ctx context.Context, user model.User) error
	GetUser(ctx context.Context, username string) (model.User, error)
	CreateSession(ctx context.Context, username, password string) (model.UserSession, error)
	GetSession(ctx context.Context, sessionId uuid.UUID) (model.UserSession, error)
	RevokeSession(ctx context.Context, session model.UserSession) error
}

func NewSQL(db *gorm.DB) (Auth, error) {
	return &authSQL{db}, nil
}

type authSQL struct {
	db *gorm.DB
}

func (a *authSQL) RegisterUser(ctx context.Context, opts UserRegistrationInfo) (model.User, error) {
	if err := validation.Struct(opts); err != nil {
		return model.User{}, err
	}

	// generate password early to prevent timing-based attacks
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(opts.Password), bcrypt.DefaultCost)
	if err != nil {
		return model.User{}, err
	}

	user := model.User{
		ID:             uuid.New(),
		Username:       opts.Username,
		HashedPassword: hashedPassword,
		CreatedAt:      time.Now(),
	}

	err = gorm.G[model.User](a.db).
		Create(ctx, &user)

	return user, wrapGormError(err)
}

func (a *authSQL) GetUser(ctx context.Context, username string) (model.User, error) {
	var user model.User

	user, err := gorm.G[model.User](a.db).
		Where("username = ?", username).
		First(ctx)

	return user, wrapGormError(err)
}

func (a *authSQL) DeleteUser(ctx context.Context, user model.User) error {
	log.Printf("unimplemented: %s", "DeleteUser")
	return nil
}

func (a *authSQL) CreateSession(ctx context.Context, username, password string) (model.UserSession, error) {

	user, err := a.GetUser(ctx, username)
	if err != nil {
		// compare hash here to prevent timing attacks
		_ = bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password))
		return model.UserSession{}, ErrInvalidCredentials
	}

	if err := bcrypt.CompareHashAndPassword(user.HashedPassword, []byte(password)); err != nil {
		return model.UserSession{}, ErrInvalidCredentials
	}

	session := model.UserSession{
		ID:        uuid.New(),
		User:      user,
		CreatedAt: time.Now(),
	}

	err = gorm.G[model.UserSession](a.db).Create(ctx, &session)
	return session, wrapGormError(err)
}

func (a *authSQL) GetSession(ctx context.Context, sessionId uuid.UUID) (model.UserSession, error) {

	session := model.UserSession{
		ID: sessionId,
	}

	session, err := gorm.G[model.UserSession](a.db).
		Preload("User", nil).
		Where("id = ?", sessionId).
		First(ctx)

	return session, wrapGormError(err)
}

func (a *authSQL) RevokeSession(ctx context.Context, session model.UserSession) error {
	log.Printf("unimplemented: %s", "RevokeSession")
	return nil
}
