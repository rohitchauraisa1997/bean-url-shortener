package repositories

import (
	"backend/models"
	"backend/packages/global"
	"context"
	"fmt"

	// "github.com/retail-ai-inc/bean/trace"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type userAuthRepo struct {
	MasterMySqlDB *gorm.DB
}

type UserauthRepository interface {
	GetAllUsers(ctx context.Context) ([]models.User, error)
	UserSignUp(ctx context.Context, userToAdd models.User) error
	AuthenticateUser(ctx context.Context, userEmailId string, userPassword string) (*models.User, error)
}

func NewUserauthRepository(MasterMySqlDB *gorm.DB) *userAuthRepo {
	return &userAuthRepo{MasterMySqlDB}
}

func getHashedPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	fmt.Println("password", password)
	fmt.Println("hashedPassword", string(hashedPassword))
	return string(hashedPassword), nil
}

func checkPasswordHash(password, hash string) bool {
	fmt.Println("password", password)
	fmt.Println("hash", hash)
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	fmt.Println("err", err)
	return err == nil
}

func (r *userAuthRepo) GetAllUsers(ctx context.Context) ([]models.User, error) {
	var allUsers []models.User

	if err := r.MasterMySqlDB.Select("id,username,email").Find(&allUsers).Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return allUsers, nil
}

func (r *userAuthRepo) UserSignUp(ctx context.Context, userToAdd models.User) error {
	var existingUser models.User

	if err := r.MasterMySqlDB.Where("username = ? OR email = ?", userToAdd.Username, userToAdd.Email).First(&existingUser).Error; err == nil {
		return global.ErrUserAlreadyExists
	}

	hashedPassword, err := getHashedPassword(userToAdd.Password)
	if err != nil {
		return errors.WithStack(err)
	}
	userToAdd.Password = string(hashedPassword)
	userToAdd.Role = "user"

	// Create a new user
	// Insert the user into the database
	if err := r.MasterMySqlDB.Create(&userToAdd).Error; err != nil {
		return fmt.Errorf("error creating user: %s", err.Error())
	}

	return nil
}

func (r *userAuthRepo) AuthenticateUser(ctx context.Context, userEmail string, userPassword string) (*models.User, error) {
	var user *models.User

	err := r.MasterMySqlDB.Where("email = ?", userEmail).First(&user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, global.ErrUserNotExists
		}
		return nil, errors.WithStack(err)
	}
	if user == nil {
		return nil, global.ErrInvalidUserCredentials
	}
	// verify the signin password
	if checkPasswordHash(userPassword, user.Password) {
		return user, nil
	} else {
		return nil, global.ErrInvalidUserCredentials
	}

}
