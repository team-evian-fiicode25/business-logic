package apis

import (
	"business_logic/internal/handlers"
	"business_logic/internal/services"
	"net/http"
)

func SetupRoutes() {
	mux := http.NewServeMux()

	consumerService := services.NewConsumerService("http://localhost:8080/graphql")
	// authoritiesService := services.NewAuthoritiesService("https://your-graphql-api.com/graphql")

	registerAuthRoutes(mux, "/api/consumer", handlers.NewConsumerAuthHandler(consumerService))
	// registerAuthRoutes(mux, "/api/authorities", handlers.NewAuthoritiesAuthHandler(authoritiesService))

	http.ListenAndServe(":8000", mux)
}

func registerAuthRoutes(mux *http.ServeMux, basePath string, handler handlers.AuthHandler) {
	mux.HandleFunc(basePath+"/auth/register", handler.Register)
	mux.HandleFunc(basePath+"/auth/login", handler.Login)
}
