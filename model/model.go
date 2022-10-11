package model

import "time"

type Dashboard struct {
	Product []Product
	Cart    Cart
}

type Cart struct {
	Name       string `json:"name"`
	Cart       []Product
	TotalPrice float64 `json:"total_price"`
}

type Product struct {
	Id       string `json:"id"`
	Name     string `json:"name"`
	Price    string `json:"price"`
	Quantity string
	Total    float64
}

type Session struct {
	Token    string    `json:"token"`
	Username string    `json:"username"`
	Expiry   time.Time `json:"expiry"`
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
