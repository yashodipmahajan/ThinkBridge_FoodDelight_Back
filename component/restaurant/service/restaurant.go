package service

import (
	db "FoodDelight/DB"
	"FoodDelight/model"
	"errors"
	"fmt"

	"github.com/google/uuid"
)

type RestaurantService struct {
}

func NewRestaurantService() *RestaurantService {

	return &RestaurantService{}
}

// ADD-RESTAURANT
func (service *RestaurantService) AddRestaurant(restaurant *model.Restaurant) (*model.Restaurant, error) {
	if restaurant.Name == "" {
		return nil, errors.New("restaurant name should not be empty")
	}
	if restaurant.Description == "" {
		return nil, errors.New("restaurant description should not be empty")
	}
	if restaurant.Location == "" {
		return nil, errors.New("restaurant location should not be empty")
	}

	db := db.NewDBConnection()
	defer db.Close()

	// Auto-migrate the Admin model to create the table if it doesn't exist
	err := db.AutoMigrate(&model.Restaurant{}).Error
	if err != nil {
		return nil, errors.New("failed to auto-migrate the admins table")
	}
	// Create the new restaurant
	if err := db.Create(&restaurant).Error; err != nil {
		return nil, err
	}

	return restaurant, nil

}

// ALL-RESTAURANTS
func (service *RestaurantService) AllRestaurants() ([]*model.Restaurant, error) {

	// fmt.Println("function called.")

	db := db.NewDBConnection()
	defer db.Close()

	var allrestaurants []*model.Restaurant

	err := db.Find(&allrestaurants).Error
	if err != nil {
		fmt.Println("Error while getting all restaurants:", err)
		return nil, errors.New("error while getting all restaurants")
	}

	return allrestaurants, nil
}

// UPDATE-RESTAURANT
func (service *RestaurantService) UpdateRestaurant(restaurantToUpdate *model.Restaurant, id uuid.UUID) (*model.Restaurant, error) {

	if restaurantToUpdate.Name == "" {
		return nil, errors.New("restaurant name should not be empty")
	}
	if restaurantToUpdate.Description == "" {
		return nil, errors.New("restaurant description should not be empty")
	}
	if restaurantToUpdate.Location == "" {
		return nil, errors.New("restaurant location should not be empty")
	}

	db := db.NewDBConnection()
	defer db.Close()

	restaurant := &model.Restaurant{}

	err := db.Where("id = ?", id).First(restaurant).Error
	if err != nil {
		return nil, errors.New("restaurant with this id not found")
	}

	var err2 = db.Model(restaurant).Where("id=?", restaurant.ID).Updates(restaurantToUpdate)
	if err2 != nil {
		return nil, err2.Error
	}

	return restaurantToUpdate, nil
}

// DELETE-RESTAURANT
func (service *RestaurantService) DeleteRestaurant(id uuid.UUID) error {

	db := db.NewDBConnection()
	defer db.Close()

	restaurant := &model.Restaurant{}

	err := db.Where("id = ?", id).First(restaurant).Error
	if err != nil {
		return errors.New("restaurant with this id not found")
	}

	var err2 = db.Model(restaurant).Where("id=?", restaurant.ID).Delete(restaurant)
	if err2 != nil {
		return err2.Error
	}

	return nil
}
