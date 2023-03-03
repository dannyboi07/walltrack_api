package main

import (
	"net/http"
	"os"
	"walltrack/common"
	"walltrack/controller"
	"walltrack/db"
	"walltrack/middleware"
	"walltrack/util"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func main() {
	if err := util.ReadEnv(); err != nil {
		util.Log.Fatalln("Failed to load env vars, err:", err)
	}

	if err := common.InitKeys(); err != nil {
		util.Log.Fatalln("Failed to load RSA keys, err:", err)
	}

	if err := db.InitDb(); err != nil {
		util.Log.Fatalln("Failed to initialize db, err:", err)
	} else if err := db.RunConfig(); err != nil {
		util.Log.Fatalln("Failed to configure db, err:", err)
	}

	var r *chi.Mux = chi.NewRouter()

	r.Use(chiMiddleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowOriginFunc: func(r *http.Request, origin string) bool {
			return origin == "http://localhost:3000"
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	}))

	r.Route("/walltrack_api", func(r chi.Router) {
		r.Get("/ping", controller.Ping)

		r.Route("/v1", func(r chi.Router) {

			r.Route("/auth", func(r chi.Router) {
				r.Get("/refresh", controller.RefreshToken)

				r.Group(func(r chi.Router) {
					r.Use(middleware.JsonRoute)
					r.Post("/login", controller.Login)
					r.Post("/register", controller.Register)
				})
			})
		})

	})

	util.Log.Println("Starting server on port:", os.Getenv("PORT"))
	if err := http.ListenAndServe("0.0.0.0:"+os.Getenv("PORT"), r); err != nil {
		util.Log.Fatalf("Failed to start server on port: %s, err: %s", os.Getenv("PORT"), err.Error())
	}
}
