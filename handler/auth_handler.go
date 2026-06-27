package handler

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/mohebul123/SpotSync/dto"
	"github.com/mohebul123/SpotSync/service"
)

type AuthHandler struct {
	srv       service.AuthService
	validator *validator.Validate
}

func NewAuthHandler(srv service.AuthService, v *validator.Validate) *AuthHandler {
	return &AuthHandler{srv: srv, validator: v}
}

func (h *AuthHandler) Register(c echo.Context) error {
	req := new(dto.RegisterRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"message": "Invalid request payload",
		})
	}
	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  err.Error(),
		})
	}

	res, err := h.srv.RegisterUser(req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"message": err.Error(),
		})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"success": true,
		"message": "User registered successfully",
		"data":    res,
	})
}

func (h *AuthHandler) Login(c echo.Context) error {
	req := new(dto.LoginRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"message": "Invalid request payload",
		})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{
			"success": false,
			"message": "Validation failed",
			"errors":  err.Error(),
		})
	}
	res, err := h.srv.LoginUser(req)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, echo.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, echo.Map{
		"success": true,
		"message": "Login successful",
		"data":    res,
	})
}
