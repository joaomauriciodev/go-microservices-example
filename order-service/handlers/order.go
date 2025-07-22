package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/joaomauriciodev/go-microservices-example/db"
	"github.com/joaomauriciodev/go-microservices-example/models"
)

type User struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type Product struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Price string `json:"price"`
}

func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	var order models.Order

	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, "Requesição inválida", http.StatusBadRequest)
		return
	}

	userChan := make(chan string)
	productChan := make(chan string)

	go fetchService(ctx, fmt.Sprintf("http://user-service:8001/users/%d", order.UserID), userChan)
	go fetchService(ctx, fmt.Sprintf("http://product-service:8002/products/%d", order.ProductID), productChan)

	var userData, productData string
	for i := 0; i < 2; i++ {
		select {
		case u := <-userChan:
			userData = u
		case p := <-productChan:
			productData = p
		case <-ctx.Done():
			http.Error(w, "Timeout ao buscar dados", http.StatusGatewayTimeout)
			return
		}
	}

	var user User
	var product Product

	if err := json.Unmarshal([]byte(userData), &user); err != nil || user.ID == "" {
		http.Error(w, "Usuário inválido", http.StatusBadRequest)
		return
	}

	if err := json.Unmarshal([]byte(productData), &product); err != nil || product.ID == "" {
		http.Error(w, "Produto inválido", http.StatusBadRequest)
		return
	}

	_, err := db.DB.Exec("INSERT INTO orders (user_id, product_id, status) VALUES ($1, $2, $3)", order.UserID, order.ProductID, "CRIADO")
	if err != nil {
		http.Error(w, "Erro ao salvar pedido", http.StatusInternalServerError)
		return
	}

	order.Status = "CRIADO"
	json.NewEncoder(w).Encode(order)

}

func fetchService(ctx context.Context, url string, ch chan<- string) {
	req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		ch <- ""
		return
	}

	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	ch <- string(body)
}
