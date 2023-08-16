package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

	"github.com/adhupraba/go-jwt-auth/controllers"
	"github.com/adhupraba/go-jwt-auth/helpers"
	"github.com/adhupraba/go-jwt-auth/initializers"
	"github.com/adhupraba/go-jwt-auth/middlewares"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.ConnectDb()
	initializers.MigrateDatabase()
}

func main() {
	router := chi.NewRouter()

	router.Use(cors.AllowAll().Handler)

	router.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		helpers.RespondWithJson(w, http.StatusOK, struct{ message string }{message: "pong"})
	})

	router.Post("/signup", controllers.Signup)
	router.Post("/signin", controllers.Signin)
	router.Get("/validate", middlewares.RequireAuth(controllers.Validate))

	port := os.Getenv("PORT")

	serve := http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Println("Server started on port", port)

	err := serve.ListenAndServe()

	if err != nil {
		log.Fatal(err)
	}
}
