package handlers

import (
	"fmt"
	"net/http"
	"time"

	"backend/models"
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

type UserauthHandler interface {
	UserSignUp(c echo.Context) error           // An example JSON response handler function
	UserSignIn(c echo.Context) error           // An example JSON response handler function
	GetCurrentUserViaJWT(c echo.Context) error // An example JSON response handler function
}

type userauthHandler struct {
	userauthService services.UserauthService
}

func NewUserauthHandler(userauthSvc services.UserauthService) *userauthHandler {
	return &userauthHandler{userauthSvc}
}

func (h *userauthHandler) UserSignUp(c echo.Context) error {
	var user models.User

	if err := c.Bind(&user); err != nil {
		return berror.NewAPIError(http.StatusBadRequest, berror.PROBLEM_PARSING_JSON, err)
	}

	if err := c.Validate(user); err != nil {
		return err
	}

	err := h.userauthService.SignUp(c.Request().Context(), user)
	if err != nil {
		if err == global.ErrUserAlreadyExists {
			return c.JSON(http.StatusConflict, map[string]string{
				"data": "user already exists",
			})
		}
		return err
	}

	return c.JSON(http.StatusOK, map[string]string{
		"data": fmt.Sprintf("user %s created successfully", user.Username),
	})
}

func (h *userauthHandler) UserSignIn(c echo.Context) error {

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

	user, err := h.userauthService.SignIn(c.Request().Context(), params.UserEmail, params.Password)
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

	userClaims := &jwtdata.UserJWTTokenData{
		UserId:    uint64(user.ID),
		UserEmail: user.Email,
		UserName:  user.Username,
		UserRole:  user.Role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(viper.GetDuration("jwt.expiration")).Unix(),
		},
	}

	secret := viper.GetString("jwt.secret")
	tokenString, err := helpers.EncodeJWT(userClaims, secret)
	if err != nil {
		return errors.WithStack(err)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"accessToken": tokenString,
			"userRole":    user.Role,
		},
	})
}

func (h *userauthHandler) GetCurrentUserViaJWT(c echo.Context) error {
	tokenString := helpers.ExtractJWTFromHeader(c)
	fmt.Println("tokenString", tokenString)
	if tokenString == "" {
		return c.JSON(http.StatusUnauthorized, map[string]interface{}{
			"data": "invalid jwt provided",
		})
	}

	secret := viper.GetString("jwt.secret")
	var userClaims jwtdata.UserJWTTokenData

	err := helpers.DecodeJWT(c, &userClaims, secret)
	if err != nil {
		if err.Error() == "token is expired" || err.Error() == "token is invalid" {
			return c.JSON(http.StatusUnauthorized, map[string]interface{}{
				"data": "invalid jwt provided",
			})
		} else {
			panic(err)
		}
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"data": map[string]interface{}{
			"user": userClaims,
		},
	})
}
