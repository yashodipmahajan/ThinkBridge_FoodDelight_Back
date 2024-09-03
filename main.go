package main

import (
	db "FoodDelight/DB"
	"FoodDelight/component/admin/controller"
	restaurantController "FoodDelight/component/restaurant/controller"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {

	db := db.NewDBConnection()
	if db == nil {
		fmt.Println("DB connection failed...")
		return
	}
	router := mux.NewRouter().StrictSlash(true)
	originsOk := handlers.AllowedOrigins([]string{os.Getenv("ORIGIN_ALLOWED"), "http://localhost:4200"})
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})

	methodsOk := handlers.AllowedMethods([]string{"GET", "DELETE", "HEAD", "POST", "PUT", "OPTIONS"})

	router = router.PathPrefix("/api/v1/fooddelight").Subrouter()

	AdminController, err := controller.NewAdminController()
	if err != nil {
		return
	}
	AdminController.AdminRoutes(router)
	RestaurantController, err := restaurantController.NewRestaurantController()
	if err != nil {
		return
	}
	RestaurantController.RestaurantRoutes(router)

	fmt.Println(http.ListenAndServe(":4002", handlers.CORS(originsOk, headersOk, methodsOk)(router)))
}
