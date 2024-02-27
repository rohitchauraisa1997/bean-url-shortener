package middlewares

import (
	"net/http"

	jwtdata "backend/packages/jwt"

	"github.com/labstack/echo/v4"
	"github.com/retail-ai-inc/bean/v2/helpers"
	"github.com/spf13/viper"
)

// JWTMiddleware is the JWT validation middleware, it also set the data in context.
func JWTMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tokenString := helpers.ExtractJWTFromHeader(c)
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

			c.Set("userId", userClaims.UserId)
			c.Set("user", userClaims)
			return next(c)
		}
	}
}

func AdminOnlyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(jwtdata.UserJWTTokenData)
		if user.UserRole != "admin" {
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"error": "Access forbidden",
			})
		}
		return next(c)
	}
}

func UserOnlyMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Get("user").(jwtdata.UserJWTTokenData)
		if user.UserRole != "user" {
			return c.JSON(http.StatusForbidden, map[string]interface{}{
				"error": "Access forbidden",
			})
		}
		return next(c)
	}
}
