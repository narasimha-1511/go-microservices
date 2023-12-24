package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)
func main(){	
	// TODO: Implement
	fmt.Println("Orders API");

	router := chi.NewRouter()

	router.Use(middleware.Logger)
	router.Get("/hello", basicHandler)

	server := &http.Server{
		Addr: ":3000",
		Handler: router,
	}

	err:= server.ListenAndServe()
	if err != nil {
		fmt.Println("Failed to listen tot he server",err)
	}

}

func basicHandler(w http.ResponseWriter, r *http.Request) {
	// fmt.Fprintf(w, "Hello World")
	w.Write([]byte("Hello World"))
}	