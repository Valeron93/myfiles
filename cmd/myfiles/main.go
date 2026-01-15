package main

import (
	"log"

	"github.com/Valeron93/myfiles/controller"
	"github.com/Valeron93/myfiles/model"
	"github.com/Valeron93/myfiles/service"
	"github.com/Valeron93/myfiles/views"
	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"gorm.io/gorm"
)

func main() {

	db, err := gorm.Open(sqlite.Open("./tmp/db.sqlite"), &gorm.Config{
		TranslateError: true,
	})

	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}

	if err := db.AutoMigrate(
		&model.User{},
		&model.UserSession{},
	); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}
	auth, err := service.NewAuthSQL(db)
	if err != nil {
		log.Fatalf("failed to create auth service: %v", err)
	}

	authController := controller.NewAuth(auth)

	app := fiber.New(fiber.Config{
		Views: views.Engine,
	})

	app.Use(logger.New(), authController.InjectSession())
	app.Get("/", func(c *fiber.Ctx) error {
		title := "index title!"
		session, ok := c.UserContext().Value(controller.SessionCtx{}).(model.UserSession)
		if ok {
			title = session.User.Username
		}
		return c.Render("index", fiber.Map{
			"Title": title,
		}, "layout")
	})

	app.Get("/register", authController.RegisterPage)
	app.Post("/api/register", authController.Register)

	app.Get("/login", authController.LoginPage)
	app.Post("/api/login", authController.Login)

	log.Fatal(app.Listen(":3000"))
}
