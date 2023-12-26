package handlers

import (
	"backend/models"
	"backend/services"
	"net/http"
	"time"

	ferror "github.com/retail-ai-inc/bean/error"
	"github.com/spf13/viper"

	"github.com/labstack/echo/v4"
)

type UrlHandler interface {
	GetUrlResolutionAnalytics(c echo.Context) error
	ShortenUrl(c echo.Context) error
	ResolveUrl(c echo.Context) error
}

type urlHandler struct {
	urlService  services.UrlService
	userService services.UserService
}

func NewUrlHandler(urlSvc services.UrlService, userSvc services.UserService) *urlHandler {
	return &urlHandler{
		urlService:  urlSvc,
		userService: userSvc,
	}
}

func (handler *urlHandler) GetUrlResolutionAnalytics(c echo.Context) error {
	tctx := c.Request().Context()
	resp, err := handler.urlService.GetAnalytics(tctx)
	if err != nil {
		panic(err)
	}

	return c.JSON(http.StatusOK, resp)
}

func (handler *urlHandler) ShortenUrl(c echo.Context) error {
	tctx := c.Request().Context()
	var bodyParams struct {
		Url    string        `json:"url" validate:"required,url"`
		Expiry time.Duration `json:"expiry"`
	}

	if err := c.Bind(&bodyParams); err != nil {
		return ferror.NewAPIError(http.StatusBadRequest, ferror.PROBLEM_PARSING_JSON, err)
	}

	if err := c.Validate(bodyParams); err != nil {
		return err
	}
	userIp := c.RealIP()

	var shortenedUrl string
	quotaRemaining, err := handler.userService.IsUserApiQuotaRemaining(tctx, userIp)
	if err != nil {
		panic(err)
	}
	if !quotaRemaining {
		ttl, err := handler.userService.GetUserTTL(tctx, userIp)
		if err != nil {
			panic(err)
		}
		return c.JSON(http.StatusServiceUnavailable, map[string]interface{}{
			"error":                    "Rate limit exceeded",
			"rate_limit_reset_in_mins": ttl / time.Nanosecond / time.Minute,
		})
	}
	ttlInMins := bodyParams.Expiry * time.Minute
	shortenedUrl, err = handler.urlService.ShortenUrl(tctx, bodyParams.Url, userIp, ttlInMins)
	if err != nil {
		panic(err)
	}
	handler.userService.DecrementApiQouta(tctx, userIp)

	res := models.Response{
		CustomShort: viper.GetString("backend.domain") + "/" + shortenedUrl,
		URL:         bodyParams.Url,
		Expiry:      (ttlInMins / time.Nanosecond / time.Minute) * 60,
		CreatedBy:   userIp,
	}

	return c.JSON(http.StatusOK, res)
}

func (handler *urlHandler) ResolveUrl(c echo.Context) error {
	tctx := c.Request().Context()
	shortenedUrl := c.Param("url")
	c.Response().Header().Set("Cache-Control", "no-store, max-age=0")
	fullUrl, err := handler.urlService.ResolveUrl(tctx, shortenedUrl)
	if err != nil {
		panic(err)
	}
	return c.Redirect(http.StatusPermanentRedirect, fullUrl)
}
