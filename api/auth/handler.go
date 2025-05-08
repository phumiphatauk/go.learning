package auth

import (
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler interface {
	Login(c echo.Context) (err error)
	RefreshToken(c echo.Context) (err error)
	Logout(c echo.Context) (err error)
}
type handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return handler{service}
}

func (h handler) Login(c echo.Context) (err error) {
	var req Login
	if err = c.Bind(&req); err != nil {
		return
	}

	loginResponse, err := h.service.Login(req.Email, req.Password)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, loginResponse)
}

func (h handler) RefreshToken(c echo.Context) (err error) {
	var req RefreshToken
	if err = c.Bind(&req); err != nil {
		return
	}

	refreshTokenResponse, err := h.service.RefreshToken(req.RefreshToken)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, refreshTokenResponse)
}

func (h handler) Logout(c echo.Context) (err error) {
	var req Logout
	if err = c.Bind(&req); err != nil {
		return
	}

	logoutResponse, err := h.service.Logout(req.RefreshToken)
	if err != nil {
		return
	}

	return c.JSON(http.StatusNoContent, logoutResponse)
}
