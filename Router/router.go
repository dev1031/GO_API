package Router

import (
	"clone_project/Middleware"
	"github.com/gorilla/mux"
)

func Router() *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/",Middleware.GetLandingPage).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/users",Middleware.GetAllUsers).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/register_user", Middleware.CreateUser).Methods("POST","OPTIONS")
	router.HandleFunc("/api/get_user" , Middleware.GetUser).Methods("GET", "OPTIONS")
	router.HandleFunc("/api/get_near_points", Middleware.GetNearPoints).Methods("GET","OPTIONS")
	return router
}