package repository

import (
	"a21hc3NpZ25tZW50/db"
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
	"fmt"
)

type CartRepository struct {
	db db.DB
}

func NewCartRepository(db db.DB) CartRepository {
	return CartRepository{db}
}

func (u *CartRepository) ReadCart() (model.Cart, error) {
	records, err := u.db.Load("carts")
	if err != nil {
		return model.Cart{}, err
	}

	if len(records) == 0 {
		return model.Cart{}, fmt.Errorf("Cart not found!")
	}

	var cart model.Cart
	err = json.Unmarshal([]byte(records), &cart)
	if err != nil {
		return model.Cart{}, err
	}

	return cart, nil
}

func (u *CartRepository) AddCart(cart model.Cart) error {
	jsonData, err := json.Marshal(cart)
	if err != nil {
		return err
	}

	u.db.Save("carts", jsonData)

	return nil
}

func (u *CartRepository) ResetCarts() error {
	err := u.db.Reset("carts", []byte(""))
	if err != nil {
		return err
	}

	return nil
}
