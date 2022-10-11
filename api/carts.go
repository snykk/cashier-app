package api

import (
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
	"net/http"
	"path"
	"strconv"
	"strings"
	"text/template"
)

func (api *API) AddCart(w http.ResponseWriter, r *http.Request) {
	// Get username context to struct model.Cart.
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	// Check r.Form with key product, if not found then return response code 400 and message "Request Product Not Found".
	var name = r.FormValue("product")

	if name == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Request Product Not Found"})
		return
	}

	var list []model.Product
	var totalPrice float64
	for _, formList := range r.Form {
		for _, v := range formList {
			item := strings.Split(v, ",")
			p, _ := strconv.ParseFloat(item[2], 64)
			q, _ := strconv.ParseFloat(item[3], 64)
			total := p * q
			list = append(list, model.Product{
				Id:       item[0],
				Name:     item[1],
				Price:    item[2],
				Quantity: item[3],
				Total:    total,
			})
			totalPrice += total
		}
	}
	var username string
	cookie, _ := r.Cookie("session_token")
	listSession, _ := api.sessionsRepo.ReadSessions()
	for _, sessionItem := range listSession {
		if sessionItem.Token == cookie.Value {
			username = sessionItem.Username
		}
	}

	// Add data field Name, Cart and TotalPrice with struct model.Cart.
	cart := model.Cart{
		Name:       username,
		Cart:       list,
		TotalPrice: totalPrice,
	}

	api.cartsRepo.AddCart(cart)

	filepath := path.Join("views", "dashboard.html")
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	listProducts, err := api.products.ReadProducts()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
	}
	data := model.Dashboard{
		Product: listProducts,
		Cart:    cart,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
	}
}
