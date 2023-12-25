package handler

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi"
	"github.com/google/uuid"
	"github.com/narasimha-1511/go-microservices/model"
	"github.com/narasimha-1511/go-microservices/repository/order"
)

type Order struct {
	Repo *order.RedisRepo;
}

func (o *Order) Create(w http.ResponseWriter, r *http.Request){
	// fmt.Println("Creating an order");
	var body struct{
	 CustomerId uuid.UUID 			`json:"customer_id"`
	 LineItems []model.LineItem `json:"line_items"`
	}

	if err:= json.NewDecoder(r.Body).Decode(&body); err!=nil{
		fmt.Println("Error decoding the request body",err);
		w.WriteHeader(http.StatusBadRequest);
		return;
	}
	
	
	now := time.Now().UTC()

	order := model.Order{
		OrderID: rand.Uint64(),
		CustomerID: body.CustomerId,
		LineItems: body.LineItems,
		CreatedAt: &now,
	}

	err:= o.Repo.Insert(r.Context(), order)

	if err != nil {
		fmt.Println("Error inserting the order", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	res, err := json.Marshal(order)
	
	if err != nil {
		fmt.Println("Error marshalling the order", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	w.Write(res)


}

func (o *Order) List(w http.ResponseWriter, r *http.Request){
	// fmt.Println("Listing orders");

	cursorStr := r.URL.Query().Get("cursor")

	if cursorStr == "" {
		cursorStr = "0"
	}

	const decimal = 10
	const bitSize = 64

	cursor, err := strconv.ParseUint(cursorStr, decimal, bitSize)

	if err != nil {
		fmt.Println("Error parsing the cursor", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	const size = 50

	res ,err := o.Repo.FindAll(r.Context(), order.FindAllPage{
		Size: size,
		Offset: cursor,
	})

	if err != nil {
		fmt.Println("Error fetching the orders", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var Response struct{
		Items []model.Order `json:"items"`
		Next uint64 `json:"next,omitempty"`
	}

	Response.Items = res.Orders
	Response.Next = res.Cursor

	data, err := json.Marshal(Response)

	if err != nil {
		fmt.Println("Error marshalling the response", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)

}

func (o *Order) GetById(w http.ResponseWriter, r *http.Request){
	// fmt.Println("Getting an order by its ID");

	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	id, err := strconv.ParseUint(idParam, base, bitSize)

	if err != nil {
		fmt.Println("Error parsing the order ID", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	or ,err := o.Repo.FindById(r.Context(), id)

	if err != nil {
		fmt.Println("Error fetching the order", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err:= json.NewEncoder(w).Encode(or); err!=nil{
		fmt.Println("Error encoding the order", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func (o *Order) UpdateById(w http.ResponseWriter, r *http.Request){
	// fmt.Println("Updating an order by its ID");

	var body struct{
	  Status string `json:"status"`
	}

	if err:= json.NewDecoder(r.Body).Decode(&body); err!=nil{
		fmt.Println("Error decoding the request body",err);
		w.WriteHeader(http.StatusBadRequest);
		return;
	}

	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	id, err := strconv.ParseUint(idParam, base, bitSize)

	if err != nil {
		fmt.Println("Error parsing the order ID", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	or ,err := o.Repo.FindById(r.Context(), id)

	if err != nil {
		fmt.Println("Error fetching the order", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	const completedstatus = "completed"
	const shippedstatus = "shipped"
	now := time.Now().UTC()

	switch body.Status {
	case completedstatus:
		if or.CompletedAt != nil {
			fmt.Println("Order already completed")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		or.CompletedAt = &now
	case shippedstatus:
		if or.ShippedAt != nil {
			fmt.Println("Order already shipped")
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		or.ShippedAt = &now
	default:
		fmt.Println("Invalid status", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = o.Repo.UpdateById(r.Context(), or)

	if err != nil {
		fmt.Println("Error updating the order", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	if err:= json.NewEncoder(w).Encode(or); err!=nil{
		fmt.Println("Error encoding the order", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	

}

func (o *Order) DeleteById(w http.ResponseWriter, r *http.Request){
	// fmt.Println("Deleting an order by its ID");

	idParam := chi.URLParam(r, "id")

	const base = 10
	const bitSize = 64

	id, err := strconv.ParseUint(idParam, base, bitSize)

	if err != nil {
		fmt.Println("Error parsing the order ID", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = o.Repo.DeleteByID(r.Context(), id)

	if err != nil {
		fmt.Println("Error deleting the order", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
	

}