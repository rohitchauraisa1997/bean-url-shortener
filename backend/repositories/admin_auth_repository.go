package repositories

import (
	"backend/models"
	"backend/packages/global"
	"context"
	"fmt"

	// "github.com/retail-ai-inc/bean/trace"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type adminAuthRepo struct {
	MasterMySqlDB *gorm.DB
}

type AdminAuthRepository interface {
	AuthenticateAdmin(ctx context.Context, userEmailId string, userPassword string) (*models.Admin, error)
}

func NewAdminAuthRepository(MasterMySqlDB *gorm.DB) *adminAuthRepo {
	return &adminAuthRepo{MasterMySqlDB}
}

func (r *adminAuthRepo) AuthenticateAdmin(ctx context.Context, userEmail string, userPassword string) (*models.Admin, error) {
	var admin *models.Admin
	getHashedPassword("beanAdminPassword")

	err := r.MasterMySqlDB.Where("email = ?", userEmail).First(&admin).Error
	fmt.Println("admin", admin)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			fmt.Println("ErrRecordNotFound")
			fmt.Println("ErrRecordNotFound")
			return nil, global.ErrUserNotExists
		}
		return nil, errors.WithStack(err)
	}
	if admin == nil {
		return nil, global.ErrInvalidUserCredentials
	}
	// verify the signin password
	if checkPasswordHash(userPassword, admin.Password) {
		return admin, nil
	} else {
		return nil, global.ErrInvalidUserCredentials
	}
}
