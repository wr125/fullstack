package main

import (
	"log"

	"github.com/joho/godotenv"
	"github.com/wr125/fullstack/views/db"
	"github.com/wr125/fullstack/views/handlers"
	"github.com/wr125/fullstack/views/services"

	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"os"
)

// In production, the secret key of the CookieStore
// and database name would be obtained from a .env file
// const (
// 	SECRET_KEY string = "secret"
// 	DB_NAME    string = "app_data.db"
// )

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(".env file could not load")
	}
	secretKey := os.Getenv("SECRET_KEY")
	dbname := os.Getenv("DB_NAME")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	e := echo.New()

	e.Static("/", "assets")

	e.HTTPErrorHandler = handlers.CustomHTTPErrorHandler

	// Helpers Middleware
	// e.Use(middleware.Recover())
	e.Use(middleware.Logger())

	// Session Middleware
	e.Use(session.Middleware(sessions.NewCookieStore([]byte(secretKey))))

	store, err := db.NewStore(dbname)
	if err != nil {
		e.Logger.Fatalf("failed to create store: %s", err)
	}

	us := services.NewUserServices(services.User{}, store)
	ah := handlers.NewAuthHandler(us)

	ts := services.NewTodoServices(services.Todo{}, store)
	th := handlers.NewTaskHandler(ts)

	// Setting Routes
	handlers.SetupRoutes(e, ah, th)

	// Start Server
	e.Logger.Fatal(e.Start("0.0.0.0:" + port))
}

/*
https://gist.github.com/taforyou/544c60ffd072c9573971cf447c9fea44
https://gist.github.com/mhewedy/4e45e04186ed9d4e3c8c86e6acff0b17

https://github.com/CurtisVermeeren/gorilla-sessions-tutorial
*/
