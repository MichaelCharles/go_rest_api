package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/joho/godotenv"

	"github.com/mcaubrey/go_rest_api/internal/database"
	"github.com/mcaubrey/go_rest_api/internal/services/comment"
	transportHTTP "github.com/mcaubrey/go_rest_api/internal/transport/http"
)

// App - the struct which contains things like pointers to database connections
type App struct {
}

// Run - sets up our application.
func (app *App) Run() error {
	env := godotenv.Load()
	if env != nil {
		log.Fatal("Error loading .env file")
	}

	fmt.Println("Setting up our app...")
	db, err := database.NewDatabase()
	if err != nil {
		return err
	}
	err = database.MigrateDB(db)
	if err != nil {
		return err
	}

	commentService := comment.NewService(db)

	handler := transportHTTP.NewHandler(commentService)
	handler.SetupRoutes()

	if err := http.ListenAndServe(":8080", handler.Router); err != nil {
		return err
	}

	return nil
}

func main() {
	fmt.Println("Starting up...")
	app := App{}
	if err := app.Run(); err != nil {
		fmt.Println("Error starting up our REST API")
		fmt.Println(err)
	}
}
