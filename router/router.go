package router 

import (
    "github.com/gorilla/mux"
    "square-8-challenge/middleware"
)

func Router() *mux.Router {
    
    // initialize router
    r := mux.NewRouter()

    // endpoints
    r.HandleFunc("/api/containers", middleware.GetContainers).Methods("GET")
    r.HandleFunc("/api/containers/{id}", middleware.GetContainer).Methods("GET")
    r.HandleFunc("/api/containers", middleware.CreateContainer).Methods("POST")
    r.HandleFunc("/api/containers/{id}", middleware.DeleteContainer).Methods("DELETE")
    r.HandleFunc("/api/containers/{id}", middleware.UpdateContainer).Methods("PUT")

    return r
}
