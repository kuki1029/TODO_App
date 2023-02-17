package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html"

	"todo/app/controller"
	"todo/app/repo"
	"todo/routes"
)

// This function will create all the needed routes for our different pages

// Stop the Fiber application
func exit(app *fiber.App) {
	_ = app.Shutdown()
}

func main() {
	// Setup the database
	repo.ConnectToDB()
	// Setup repo and controller interfaces
	ur := repo.NewUserRepo(repo.DB.DbConn)
	uc := controller.NewUserController(*ur)
	tr := repo.NewTaskRepo(repo.DB.DbConn)
	tc := controller.NewTaskController(*tr)
	// Create a new engine
	engine := html.New("./resources/views", ".html")

	// Pass the engine to the Views
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	// This serves the css files so it the HTML can render it
	app.Static("/static", "./static")
	app.Static("/resources/JS", "./resources/JS")
	// Serves all the HTML filse
	app.Static("/", "./resources/views", fiber.Static{
		Index: "login.html",
	})

	routes.SetupTaskRoutes(app, tc)
	routes.SetupUserRoutes(app, uc)

	// Close any connections on interrupt signal
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		<-c
		exit(app)
	}()

	// Start listening on the specified address
	if err := app.Listen("127.0.0.1:8080"); err != nil {
		log.Panic(err)
	}
}
