package handler

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/mohebul123/SpotSync/dto"
	"github.com/mohebul123/SpotSync/service"
)

type ZoneHandler struct {
	srv       service.ZoneService
	validator *validator.Validate
}

func NewZoneHandler(srv service.ZoneService, v *validator.Validate) *ZoneHandler {
	return &ZoneHandler{srv: srv, validator: v}
}

func (h *ZoneHandler) Create(c echo.Context) error {
	req := new(dto.CreateZoneRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"success": false, "message": "Invalid request payload"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"success": false, "message": "Validation failed", "errors": err.Error()})
	}

	res, err := h.srv.CreateZone(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"success": false, "message": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"success": true,
		"message": "Parking zone created successfully",
		"data":    res,
	})
}

func (h *ZoneHandler) GetAll(c echo.Context) error {
	res, err := h.srv.GetAllZones()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"success": false, "message": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"success": true,
		"message": "Parking zones retrieved successfully",
		"data":    res,
	})
}

func (h *ZoneHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"success": false, "message": "Invalid zone ID"})
	}

	res, err := h.srv.GetZoneByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, echo.Map{"success": false, "message": "Parking zone not found"})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"success": true,
		"message": "Parking zone retrieved successfully",
		"data":    res,
	})
}
