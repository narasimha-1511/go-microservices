package handler

import (
	"fmt"
	"net/http"
)

type Order struct {}

func (o *Order) Create(w http.ResponseWriter, r *http.Request){
	fmt.Println("Creating an order");
	// w.WriteHeader(http.StatusCreated);
}

func (o *Order) List(w http.ResponseWriter, r *http.Request){
	fmt.Println("Listing orders");
}

func (o *Order) GetById(w http.ResponseWriter, r *http.Request){
	fmt.Println("Getting an order by its ID");
}

func (o *Order) UpdateById(w http.ResponseWriter, r *http.Request){
	fmt.Println("Updating an order by its ID");
}

func (o *Order) DeleteById(w http.ResponseWriter, r *http.Request){
	fmt.Println("Deleting an order by its ID");
}