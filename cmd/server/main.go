package main

import (
	"net/http"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/mcaubrey/go_rest_api/internal/database"
	"github.com/mcaubrey/go_rest_api/internal/services/comment"
	"github.com/mcaubrey/go_rest_api/internal/services/user"
	transportHTTP "github.com/mcaubrey/go_rest_api/internal/transport/http"
)

// App - contain application information
type App struct {
	Version string
	Name    string
}

// Run - sets up our application.
func (app *App) Run() error {
	logrus.SetFormatter(&logrus.JSONFormatter{})

	env := godotenv.Load()
	if env != nil {
		logrus.Fatal("Could not read env file")
	}

	logrus.WithFields(logrus.Fields{
		"AppName":    app.Name,
		"AppVersion": app.Version,
	}).Info("Setting up application")

	db, err := database.NewDatabase()
	if err != nil {
		return err
	}

	commentService := comment.NewService(db)
	userService := user.NewService(db)

	err = database.MigrateDB(db)
	if err != nil {
		return err
	}

	handler := transportHTTP.NewHandler(commentService, userService)
	handler.SetupRoutes()

	if err := http.ListenAndServe(":8080", handler.Router); err != nil {
		return err
	}

	return nil
}

func main() {
	logrus.Info("Starting up...")
	app := App{
		Name:    "Commenting Service",
		Version: "0.0.1",
	}

	if err := app.Run(); err != nil {
		logrus.Error("Error starting up our REST API")
		logrus.Fatal(err)
	}
}
