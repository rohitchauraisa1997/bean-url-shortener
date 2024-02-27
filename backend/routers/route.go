// MIT License

// Copyright (c) The RAI Authors

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.
package routers

import (
	"net/http"

	"backend/handlers"
	"backend/middlewares"
	"backend/repositories"
	"backend/services"

	"github.com/labstack/echo/v4"
	bean "github.com/retail-ai-inc/bean/v2"
	"github.com/retail-ai-inc/bean/v2/store/redis"
	"github.com/spf13/viper"
)

type Repositories struct {
	userAuthRepo  repositories.UserauthRepository // added by bean
	urlRepo       repositories.UrlRepository      // added by bean
	userRepo      repositories.UserRepository
	adminAuthRepo repositories.AdminAuthRepository
}

type Services struct {
	userauthSvc  services.UserauthService // added by bean
	urlSvc       services.UrlService      // added by bean
	userSvc      services.UserService
	adminAuthSvc services.AdminAuthService
}

type Handlers struct {
	userauthHdlr  handlers.UserauthHandler // added by bean
	urlHdlr       handlers.UrlHandler      // added by bean
	adminAuthHdlr handlers.AdminauthHandler
}

func Init(b *bean.Bean) {

	e := b.Echo

	repos := &Repositories{
		userAuthRepo:  repositories.NewUserauthRepository(b.DBConn.MasterMySQLDB), // added by bean
		urlRepo:       repositories.NewUrlRepository(redis.NewMasterCache(b.DBConn.MasterRedisDB, viper.GetString("database.redis.prefix"))),
		userRepo:      repositories.NewUserRepository(redis.NewMasterCache(b.DBConn.MasterRedisDB, viper.GetString("database.redis.prefix"))),
		adminAuthRepo: repositories.NewAdminAuthRepository(b.DBConn.MasterMySQLDB),
	}

	svcs := &Services{
		userauthSvc:  services.NewUserauthService(repos.userAuthRepo), // added by bean
		urlSvc:       services.NewUrlService(repos.urlRepo),           // added by bean
		userSvc:      services.NewUserService(repos.userRepo),
		adminAuthSvc: services.NewAdminAuthService(repos.adminAuthRepo),
	}

	hdlrs := &Handlers{
		userauthHdlr:  handlers.NewUserauthHandler(svcs.userauthSvc),                       // added by bean
		urlHdlr:       handlers.NewUrlHandler(svcs.urlSvc, svcs.userSvc, svcs.userauthSvc), // added by bean
		adminAuthHdlr: handlers.NewAdminAuthHandler(svcs.adminAuthSvc),
	}

	// Default index page goes to above JSON (/json) index page.
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": `backend 🚀`,
		})
	})

	// accessible to everyone
	e.GET("/:url", hdlrs.urlHdlr.ResolveUrl)

	// provisioning apis for users
	user := e.Group("/user")
	user.POST("/signup", hdlrs.userauthHdlr.UserSignUp)
	user.POST("/signin", hdlrs.userauthHdlr.UserSignIn)
	user.GET("/me", hdlrs.userauthHdlr.GetCurrentUserViaJWT)

	// accessible to users
	urlShortener := e.Group("/url-shortener")
	urlShortener.Use(middlewares.JWTMiddleware())
	urlShortener.GET("/resolutions/analytics", hdlrs.urlHdlr.GetUrlResolutionAnalyticsForUser)
	urlShortener.POST("/api/shorten", hdlrs.urlHdlr.ShortenUrl, middlewares.UserOnlyMiddleware)

	// accessible only to admin
	// admin can only go through the analytics but cant use the application to create urls
	// admin := e.Group("/admin")
	// admin.POST("/signin", hdlrs.adminAuthHdlr.AdminSignIn)
	// admin.Use(middlewares.JWTMiddleware())
	urlShortener.GET("/resolutions/analytics/all", hdlrs.urlHdlr.GetUrlResolutionAnalyticsForAllUsers, middlewares.AdminOnlyMiddleware)
}
