package service

import (
	db "FoodDelight/DB"
	"FoodDelight/model"
	"errors"
	"net/mail"

	"github.com/google/uuid"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type AdminService struct {
}

func NewAdminService() *AdminService {

	return &AdminService{}
}

const cost = 10

func doesAdminExistWithUsername(username string) error {
	db := db.NewDBConnection()

	defer db.Close()

	admin := &model.Admin{}

	err := db.Where("username = ?", username).First(admin).Error

	if err == nil {
		return errors.New("admin exists with this username")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// If there's an error other than record not found
		return err
	}

	return nil

}
func doesAdminExistWithEmail(email string) error {
	db := db.NewDBConnection()

	defer db.Close()

	admin := &model.Admin{}

	err := db.Where("email = ?", email).First(admin).Error

	if err == nil {
		return errors.New("admin exists with this email")
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		// If there's an error other than record not found
		return err
	}

	return nil
}

//REG ADMIN

func (service *AdminService) RegisterAdmin(admin *model.Admin) (*model.Admin, error) {

	db := db.NewDBConnection()
	defer db.Close()

	// Auto-migrate the Admin model to create the table if it doesn't exist
	err := db.AutoMigrate(&model.Admin{}).Error
	if err != nil {
		return nil, errors.New("failed to auto-migrate the admins table")
	}

	// Check if the username already exists
	err = doesAdminExistWithUsername(admin.Username)
	if err != nil {
		return nil, err
	}

	// Validate the email address
	_, err = mail.ParseAddress(admin.Email)
	if err != nil {
		return nil, errors.New("invalid email address")
	}

	// Check if the email already exists
	err = doesAdminExistWithEmail(admin.Email)
	if err != nil {
		return nil, err
	}

	// Hash the password
	hashedPass, err := bcrypt.GenerateFromPassword([]byte(admin.Password), cost)
	if err != nil {
		return nil, err
	}
	admin.Password = string(hashedPass)

	// Create the new admin
	if err := db.Create(&admin).Error; err != nil {
		return nil, err
	}

	return admin, nil

}

// LOGIN-ADMIN
func (service *AdminService) LoginAdmin(credentials *model.Admin) (error, *model.Admin) {

	db := db.NewDBConnection()
	defer db.Close()

	err := doesAdminExistWithUsername(credentials.Username)
	if err == nil {
		return errors.New("admin not exist with this username"), nil
	}

	var foundAdmin = &model.Admin{}
	err = db.Where("username = ?", credentials.Username).First(foundAdmin).Error
	if err != nil {
		return errors.New("admin not exists with this username"), nil
	}

	// Validate password
	err = bcrypt.CompareHashAndPassword([]byte(foundAdmin.Password), []byte(credentials.Password))
	if err != nil {

		return errors.New("invalid password"), nil
	}

	return nil, foundAdmin
}

// LOGOUT-ADMIN
func (service *AdminService) LogoutAdmin(id uuid.UUID) (*model.Admin, error) {

	err, foundAdmin := model.FindAdminByID(id)
	if err != nil {
		return nil, err
	}

	return foundAdmin, nil

}

//UPDATE ADMIN

func (service *AdminService) UpdateAdmin(adminToUpdate *model.Admin, id uuid.UUID) (*model.Admin, error) {
	if adminToUpdate.Fullname == "" {
		return nil, errors.New("restaurant name should not be empty")
	}
	if adminToUpdate.Username == "" {
		return nil, errors.New("restaurant description should not be empty")
	}
	if adminToUpdate.Email == "" {
		return nil, errors.New("email description should not be empty")
	}

	db := db.NewDBConnection()
	defer db.Close()

	admin := &model.Admin{}

	err := db.Where("id = ?", id).First(admin).Error
	if err != nil {
		return nil, errors.New("admin with this id not found")
	}

	var err2 = db.Model(admin).Where("id=?", admin.ID).Updates(adminToUpdate)
	if err2 != nil {
		return nil, err2.Error
	}

	return adminToUpdate, nil
}

// DELETE ADMIN
func (service *AdminService) DeleteAdmin(id uuid.UUID) error {

	db := db.NewDBConnection()
	defer db.Close()

	admin := &model.Admin{}

	err := db.Where("id = ?", id).First(admin).Error
	if err != nil {
		return errors.New("admin with this id not found")
	}

	var err2 = db.Model(admin).Where("id=?", admin.ID).Delete(admin)
	if err2 != nil {
		return err2.Error
	}

	return nil
}
