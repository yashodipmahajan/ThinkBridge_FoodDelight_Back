package controller

import (
	"FoodDelight/component/restaurant/service"
	"FoodDelight/model"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type RestaurantController struct {
	RestaurantService *service.RestaurantService
}

func NewRestaurantController() (*RestaurantController, error) {
	var restaurantService = service.NewRestaurantService()
	return &RestaurantController{
		RestaurantService: restaurantService,
	}, nil
}

func (c *RestaurantController) RestaurantRoutes(router *mux.Router) {

	adminRouter := router.PathPrefix("/admin/restaurant").Subrouter()
	// contactguardedrouter := userRouter.PathPrefix("/").Subrouter()

	adminRouter.HandleFunc("/add-restaurant", c.AddRestaurant).Methods(http.MethodPost)
	adminRouter.HandleFunc("/all-restaurants", c.AllRestaurants).Methods(http.MethodGet)
	adminRouter.HandleFunc("/{id}/update-restaurant", c.UpdateRestaurant).Methods(http.MethodPut)
	adminRouter.HandleFunc("/{id}/delete-restaurant", c.DeleteRestaurant).Methods(http.MethodDelete)

	// contactguardedrouter.Use(authorization.ContactMiddleware)

	fmt.Println("==========================Restaurant-Routes-Registered==========================")

}

func (controller *RestaurantController) AddRestaurant(w http.ResponseWriter, r *http.Request) {
	fmt.Println("_____________Add-Restaurant____________")
	var newRestaurant = &model.Restaurant{}

	json.NewDecoder(r.Body).Decode(&newRestaurant)

	addedrestaurant, err := controller.RestaurantService.AddRestaurant(newRestaurant)

	if err != nil {
		json.NewEncoder(w).Encode(err)
		w.WriteHeader(http.StatusBadGateway)
		return
	}

	json.NewEncoder(w).Encode("restaurant added successfully.. ")
	json.NewEncoder(w).Encode(&addedrestaurant)
	w.WriteHeader(200)

}

func (controller *RestaurantController) UpdateRestaurant(w http.ResponseWriter, r *http.Request) {
	fmt.Println("_________Update-Restaurant____________")
	var restaurantToUpdate = &model.Restaurant{}

	json.NewDecoder(r.Body).Decode(&restaurantToUpdate)

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

	updatedRestaurant, err := controller.RestaurantService.UpdateRestaurant(restaurantToUpdate, restaurantId)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	// json.NewEncoder(w).Encode("restaurant updated successfully.. ")
	json.NewEncoder(w).Encode(&updatedRestaurant)
	w.WriteHeader(200)

}

func (controller *RestaurantController) DeleteRestaurant(w http.ResponseWriter, r *http.Request) {
	fmt.Println("__________Delete-Restaurant____________")
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

	err = controller.RestaurantService.DeleteRestaurant(restaurantId)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}

	json.NewEncoder(w).Encode("restaurant deleted successfully.. ")
	w.WriteHeader(200)
}

func (controller *RestaurantController) AllRestaurants(w http.ResponseWriter, r *http.Request) {

	fmt.Println("________________ALL-RESTAURANTS______________")
	allrestaurants, err := controller.RestaurantService.AllRestaurants()
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadGateway)
		return
	}
	json.NewEncoder(w).Encode(allrestaurants)
	w.WriteHeader(200)
}
