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
	return string(hashedPassword), nil
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
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
