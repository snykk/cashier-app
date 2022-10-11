package api

import (
	"a21hc3NpZ25tZW50/model"
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"text/template"
	"time"

	"github.com/google/uuid"
)

func (api *API) Register(w http.ResponseWriter, r *http.Request) {
	// Read username and password request with FormValue.
	creds := model.Credentials{}

	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	var username = r.FormValue("username")
	var password = r.FormValue("password")

	// Handle request if creds is empty send response code 400, and message "Username or Password empty"
	if username == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Username or Password empty"})
		return
	}

	creds.Username = username
	creds.Password = password

	err := api.usersRepo.AddUser(creds)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: err.Error()})
		return
	}

	filepath := path.Join("views", "status.html")
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	var data = map[string]string{
		"name":    creds.Username,
		"message": "register success!",
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
	}
}

func (api *API) Login(w http.ResponseWriter, r *http.Request) {
	// Read usernmae and password request with FormValue.
	if err := r.ParseForm(); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	var username = r.FormValue("username")
	var password = r.FormValue("password")

	// Handle request if creds is empty send response code 400, and message "Username or Password empty"
	if username == "" || password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Username or Password empty"})
		return
	}

	// Create creds object
	creds := model.Credentials{
		Username: username,
		Password: password,
	}

	listUser, err := api.usersRepo.ReadUser()
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	if !api.usersRepo.LoginValid(listUser, creds) {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Wrong User or Password!"})
		return
	}

	// Generate Cookie with Name "session_token", Path "/", Value "uuid generated with github.com/google/uuid", Expires time to 5 Hour.
	sessionToken := uuid.New()

	session := model.Session{
		Token:    sessionToken.String(),
		Username: creds.Username,
		Expiry:   time.Now().Add(time.Hour * 5),
	}
	err = api.sessionsRepo.AddSessions(session)

	cookie := &http.Cookie{
		Name:    "session_token",
		Value:   sessionToken.String(),
		Path:    "/",
		Expires: time.Now().Add(time.Hour * 5),
	}

	http.SetCookie(w, cookie)

	if err != nil {
		fmt.Println("masuk panik add")
		panic(err)
	}

	filepath := path.Join("views", "dashboard.html")
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		fmt.Println("Masuk panik tmpl")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	listProducts, err := api.products.ReadProducts()

	if err != nil {
		fmt.Println("masuk panik list")
		panic(err)
	}

	data := model.Dashboard{
		Product: listProducts,
		Cart: model.Cart{
			Name:       creds.Username,
			Cart:       []model.Product{},
			TotalPrice: 0,
		},
	}

	w.WriteHeader(http.StatusOK)
	err = tmpl.Execute(w, data)
	if err != nil {
		fmt.Println("masuk panik tmpee")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}
}

func (api *API) Logout(w http.ResponseWriter, r *http.Request) {
	//Read session_token and get Value:
	cookie, err := r.Cookie("session_token")
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "http: named cookie not present"})
	}

	sessionToken := cookie.Value

	// Delete session cookies
	api.sessionsRepo.DeleteSessions(sessionToken)

	filepath := path.Join("views", "login.html")
	tmpl, err := template.ParseFiles(filepath)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(model.ErrorResponse{Error: "Internal Server Error"})
	}
}

func (api *API) RegisterView(w http.ResponseWriter, r *http.Request) {
	filepath := path.Join("views", "register.html")
	var tmpl = template.Must(template.ParseFiles(filepath))
	var err = tmpl.Execute(w, nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func (api *API) LoginView(w http.ResponseWriter, r *http.Request) {
	filepath := path.Join("views", "login.html")
	var tmpl = template.Must(template.ParseFiles(filepath))
	var err = tmpl.Execute(w, nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
