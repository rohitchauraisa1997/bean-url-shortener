package handlers

import (
	"net/http"
	"time"

	"backend/packages/global"
	"backend/services"

	jwtdata "backend/packages/jwt"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/pkg/errors"
	berror "github.com/retail-ai-inc/bean/v2/error"
	"github.com/retail-ai-inc/bean/v2/helpers"
	"github.com/spf13/viper"
)

type AdminauthHandler interface {
	AdminSignIn(c echo.Context) error // An example JSON response handler function
}

type adminauthHandler struct {
	adminauthService services.AdminAuthService
}

func NewAdminAuthHandler(adminauthSvc services.AdminAuthService) *adminauthHandler {
	return &adminauthHandler{adminauthSvc}
}

func (h *adminauthHandler) AdminSignIn(c echo.Context) error {

	var params struct {
		UserEmail string `json:"email" validate:"required,email"`
		Password  string `json:"password" validate:"required"`
	}

	if err := c.Bind(&params); err != nil {
		return berror.NewAPIError(http.StatusBadRequest, berror.PROBLEM_PARSING_JSON, err)
	}

	if err := c.Validate(params); err != nil {
		return err
	}

	admin, err := h.adminauthService.SignIn(c.Request().Context(), params.UserEmail, params.Password)
	if err != nil {
		if err == global.ErrInvalidUserCredentials {
			return c.JSON(http.StatusUnauthorized, map[string]string{
				"data": "invalid credentials",
			})
		}

		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}

	adminClaims := &jwtdata.UserJWTTokenData{
		UserId:    uint64(admin.ID),
		UserEmail: admin.Email,
		UserName:  admin.Username,
		UserRole:  "admin",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(viper.GetDuration("jwt.expiration")).Unix(),
		},
	}

	secret := viper.GetString("jwt.secret")
	tokenString, err := helpers.EncodeJWT(adminClaims, secret)
	if err != nil {
		return errors.WithStack(err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"accessToken": tokenString,
		},
	})
}
