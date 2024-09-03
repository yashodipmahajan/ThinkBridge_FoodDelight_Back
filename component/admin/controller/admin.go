package controller

import (
	"FoodDelight/component/admin/service"
	"FoodDelight/component/security/middleware/authorization"
	"FoodDelight/model"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type AdminController struct {
	Adminservice *service.AdminService
}

func NewAdminController() (*AdminController, error) {
	var AdminService = service.NewAdminService()
	return &AdminController{
		Adminservice: AdminService,
	}, nil
}

func (c *AdminController) AdminRoutes(router *mux.Router) {

	adminRouter := router.PathPrefix("/admin").Subrouter()
	// contactguardedrouter := userRouter.PathPrefix("/").Subrouter()

	adminRouter.HandleFunc("/register", c.RegisterAdmin).Methods(http.MethodPost)
	adminRouter.HandleFunc("/login", c.LoginAdmin).Methods(http.MethodPost)
	adminRouter.HandleFunc("/{id}/logout", c.LogoutAdmin).Methods(http.MethodPost)
	adminRouter.HandleFunc("/{id}/update-admin", c.UpdateAdmin).Methods(http.MethodPut)
	adminRouter.HandleFunc("/{id}/delete-admin", c.DeleteAdmin).Methods(http.MethodDelete)

	// contactguardedrouter.Use(authorization.ContactMiddleware)

	fmt.Println("==========================Admin-Routes-Registered==========================")

}

// REG ADMIN
func (controller *AdminController) RegisterAdmin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("ADMIN REGISTRATION STARTED.....")
	var newAdmin = &model.Admin{}
	json.NewDecoder(r.Body).Decode(&newAdmin)

	regAdmin, err := controller.Adminservice.RegisterAdmin(newAdmin)
	if err != nil {
		fmt.Println(err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(&regAdmin)
}

// LOGIN ADMIN
func (controller *AdminController) LoginAdmin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOGIN ROUTE STARTED..")

	var credentials *model.Admin

	json.NewDecoder(r.Body).Decode(&credentials)

	err, loggedinAdmin := controller.Adminservice.LoginAdmin(credentials)
	if err != nil {
		json.NewEncoder(w).Encode("USER NOT FOUND")
		w.WriteHeader(401)
		return
	}

	var admin = loggedinAdmin

	claim := &authorization.Claims{
		Id: admin.ID,
	}

	token, err := claim.Coder()
	if err != nil {
		json.NewEncoder(w).Encode(err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	fmt.Println("token", token)

	http.SetCookie(w, &http.Cookie{
		Name:    "auth",
		Value:   token,
		Expires: time.Now().Add(time.Minute + 100),
	})

	// json.NewEncoder(w).Encode(&token)

	var reponse = model.ResponseDto{
		LoggedInAdmin: *loggedinAdmin,
		Token:         token,
	}

	json.NewEncoder(w).Encode(&reponse)
	w.WriteHeader(200)

}

// LOGOUT-ADMIN
func (controller *AdminController) LogoutAdmin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("LOGIN ROUTE STARTED..")
	params := mux.Vars(r)
	err := uuid.Validate(params["id"])
	if err != nil {
		http.Error(w, errors.New("invalid userId").Error(), http.StatusBadGateway)
		return
	}
	adminId, err := uuid.Parse(params[("id")])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	_, err = controller.Adminservice.LogoutAdmin(adminId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	cookie, err := r.Cookie("auth")

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "auth",
		Expires: time.Now(), // Set the expiration to a past time
		MaxAge:  -1,         // Negative value to delete the cookie immediately
	})

	currentcookie, err := r.Cookie("auth")
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("cookie not found")
	} else {
		fmt.Println("currentcookie..", currentcookie)
	}
	json.NewEncoder(w).Encode(cookie.Value)
	json.NewEncoder(w).Encode("LogOut Successfully..")

	w.WriteHeader(200)

}

// UPDATE-ADMIN
func (controller *AdminController) UpdateAdmin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("_________Update-Admin________-")
	var adminToUpdate = &model.Admin{}

	json.NewDecoder(r.Body).Decode(&adminToUpdate)

	params := mux.Vars(r)
	err := uuid.Validate(params["id"])
	if err != nil {
		http.Error(w, errors.New("invalid userId").Error(), http.StatusBadGateway)
		return
	}
	restaurantId, err := uuid.Parse(params[("id")])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	updatedAdmin, err := controller.Adminservice.UpdateAdmin(adminToUpdate, restaurantId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	json.NewEncoder(w).Encode("admin updated successfully.. ")
	json.NewEncoder(w).Encode(&updatedAdmin)
	w.WriteHeader(200)
}

// DELETE-ADMIN
func (controller *AdminController) DeleteAdmin(w http.ResponseWriter, r *http.Request) {
	fmt.Println("__________Delete-Admin__________")
	params := mux.Vars(r)
	err := uuid.Validate(params["id"])
	if err != nil {
		http.Error(w, errors.New("invalid userId").Error(), http.StatusBadGateway)
		return
	}
	adminId, err := uuid.Parse(params[("id")])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	err = controller.Adminservice.DeleteAdmin(adminId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	json.NewEncoder(w).Encode("admin deleted successfully.. ")
	w.WriteHeader(200)
}
