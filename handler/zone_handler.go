package handler

import (
	"net/http"
	"strconv"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/mohebul123/SpotSync/dto"
	"github.com/mohebul123/SpotSync/service"
)

type ZoneAPIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Errors  interface{} `json:"errors,omitempty"`
}

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
		return c.JSON(http.StatusBadRequest, ZoneAPIResponse{Success: false, Message: "Invalid request payload"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, ZoneAPIResponse{Success: false, Message: "Validation failed", Errors: err.Error()})
	}

	res, err := h.srv.CreateZone(req)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ZoneAPIResponse{Success: false, Message: err.Error()})
	}

	return c.JSON(http.StatusCreated, ZoneAPIResponse{
		Success: true,
		Message: "Parking zone created successfully",
		Data:    res,
	})
}

func (h *ZoneHandler) GetAll(c echo.Context) error {
	res, err := h.srv.GetAllZones()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, ZoneAPIResponse{Success: false, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, ZoneAPIResponse{
		Success: true,
		Message: "Parking zones retrieved successfully",
		Data:    res,
	})
}

func (h *ZoneHandler) GetByID(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ZoneAPIResponse{Success: false, Message: "Invalid zone ID"})
	}

	res, err := h.srv.GetZoneByID(uint(id))
	if err != nil {
		return c.JSON(http.StatusNotFound, ZoneAPIResponse{Success: false, Message: "Parking zone not found"})
	}

	return c.JSON(http.StatusOK, ZoneAPIResponse{
		Success: true,
		Message: "Parking zone retrieved successfully",
		Data:    res,
	})
}

func (h *ZoneHandler) Update(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ZoneAPIResponse{Success: false, Message: "Invalid zone ID"})
	}

	req := new(dto.UpdateZoneRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, ZoneAPIResponse{Success: false, Message: "Invalid request payload"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, ZoneAPIResponse{Success: false, Message: "Validation failed", Errors: err.Error()})
	}

	res, err := h.srv.UpdateZone(uint(id), req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, ZoneAPIResponse{Success: false, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, ZoneAPIResponse{
		Success: true,
		Message: "Parking zone updated successfully",
		Data:    res,
	})
}

func (h *ZoneHandler) Delete(c echo.Context) error {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ZoneAPIResponse{Success: false, Message: "Invalid zone ID"})
	}

	err = h.srv.DeleteZone(uint(id))
	if err != nil {
		return c.JSON(http.StatusBadRequest, ZoneAPIResponse{Success: false, Message: err.Error()})
	}

	return c.JSON(http.StatusOK, ZoneAPIResponse{
		Success: true,
		Message: "Parking zone deleted successfully",
	})
}
