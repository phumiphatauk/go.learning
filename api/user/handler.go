package user

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

type Handler interface {
	Register(c echo.Context) (err error)
	GetList(c echo.Context) (err error)
	Get(c echo.Context) (err error)
	Update(c echo.Context) (err error)
	Delete(c echo.Context) (err error)
}

type handler struct {
	service Service
}

func NewHandler(service Service) Handler {
	return handler{service}
}

func (h handler) Register(c echo.Context) (err error) {
	var req CreateUser
	if err = c.Bind(&req); err != nil {
		return
	}

	err = h.service.CreateUser(req)
	if err != nil {
		return
	}

	return c.JSON(http.StatusCreated, nil)
}

func (h handler) GetList(c echo.Context) (err error) {
	var queryParams GetUserList
	// Bind query parameters to the struct
	if err = c.Bind(&queryParams); err != nil {
		return
	}

	// Validate query parameters
	if queryParams.Page < 1 {
		queryParams.Page = 1
	}

	if queryParams.Limit < 1 {
		queryParams.Limit = 10
	}

	if queryParams.Sort == "" {
		queryParams.Sort = "first_name"
	}

	if queryParams.SortDirection == "" {
		queryParams.SortDirection = "asc"
	}

	// Call the service to get the user list
	users, err := h.service.GetUserList(queryParams)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, "Failed to get user list")
	}

	return c.JSON(http.StatusOK, users)
}

func (h handler) Get(c echo.Context) (err error) {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "ID is required")
	}

	// Convert id to uint
	var userID uint
	if _, err := fmt.Sscanf(id, "%d", &userID); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID format")
	}

	user, err := h.service.GetUserByID(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	return c.JSON(http.StatusOK, user)
}

func (h handler) Update(c echo.Context) (err error) {
	var req UpdateUser
	if err = c.Bind(&req); err != nil {
		return
	}

	err = h.service.UpdateUser(req)
	if err != nil {
		return
	}

	return c.JSON(http.StatusOK, nil)
}

func (h handler) Delete(c echo.Context) (err error) {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, "ID is required")
	}

	// Convert id to uint
	var userID uint
	if _, err := fmt.Sscanf(id, "%d", &userID); err != nil {
		return c.JSON(http.StatusBadRequest, "Invalid ID format")
	}

	err = h.service.DeleteUser(userID)
	if err != nil {
		return c.JSON(http.StatusNotFound, "User not found")
	}

	return c.JSON(http.StatusNoContent, nil)
}
