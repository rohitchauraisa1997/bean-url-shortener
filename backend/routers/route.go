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
	userauthRepo repositories.UserauthRepository // added by bean
	urlRepo      repositories.UrlRepository      // added by bean
	userRepo     repositories.UserRepository
}

type Services struct {
	userauthSvc services.UserauthService // added by bean
	urlSvc      services.UrlService      // added by bean
	userSvc     services.UserService
}

type Handlers struct {
	userauthHdlr handlers.UserauthHandler // added by bean
	urlHdlr      handlers.UrlHandler      // added by bean
}

func Init(b *bean.Bean) {

	e := b.Echo

	repos := &Repositories{
		userauthRepo: repositories.NewUserauthRepository(b.DBConn.MasterMySQLDB), // added by bean
		urlRepo:      repositories.NewUrlRepository(redis.NewMasterCache(b.DBConn.MasterRedisDB, viper.GetString("database.redis.prefix"))),
		userRepo:     repositories.NewUserRepository(redis.NewMasterCache(b.DBConn.MasterRedisDB, viper.GetString("database.redis.prefix"))),
	}

	svcs := &Services{
		userauthSvc: services.NewUserauthService(repos.userauthRepo), // added by bean
		urlSvc:      services.NewUrlService(repos.urlRepo),           // added by bean
		userSvc:     services.NewUserService(repos.userRepo),
	}

	hdlrs := &Handlers{
		userauthHdlr: handlers.NewUserauthHandler(svcs.userauthSvc),     // added by bean
		urlHdlr:      handlers.NewUrlHandler(svcs.urlSvc, svcs.userSvc), // added by bean
	}

	// Default index page goes to above JSON (/json) index page.
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"message": `backend 🚀`,
		})
	})

	user := e.Group("/user")
	user.POST("/signup", hdlrs.userauthHdlr.UserSignUp)
	user.POST("/signin", hdlrs.userauthHdlr.UserSignIn)
	user.GET("/me", hdlrs.userauthHdlr.GetCurrentUserViaJWT)

	urlShortener := e.Group("/url-shortener")
	urlShortener.Use(middlewares.JWTMiddleware())
	urlShortener.GET("/resolutions/analytics", hdlrs.urlHdlr.GetUrlResolutionAnalyticsForUser)
	urlShortener.GET("/client-ip", func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{
			"clientIp": c.RealIP(),
		})
	})
	urlShortener.POST("/api/shorten", hdlrs.urlHdlr.ShortenUrl)
	e.GET("/:url", hdlrs.urlHdlr.ResolveUrl)
	urlShortener.GET("/admin/resolutions/analytics", hdlrs.urlHdlr.GetUrlResolutionAnalyticsForAllUsers)
}
