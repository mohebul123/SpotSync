package handler

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/mohebul123/SpotSync/dto"
	"github.com/mohebul123/SpotSync/service"
)

type ReservationHandler struct {
	srv       service.ReservationService
	validator *validator.Validate
}

func NewReservationHandler(srv service.ReservationService, v *validator.Validate) *ReservationHandler {
	return &ReservationHandler{srv: srv, validator: v}
}

// Spec #6: Reserve Parking Spot
func (h *ReservationHandler) Book(c echo.Context) error {
	userID := c.Get("userID").(uint)

	req := new(dto.CreateReservationRequest)
	if err := c.Bind(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"success": false, "message": "Invalid request payload"})
	}

	if err := h.validator.Struct(req); err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"success": false, "message": "Validation failed", "errors": err.Error()})
	}

	res, err := h.srv.BookSpot(userID, req)
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"success": false, "message": err.Error()})
	}

	return c.JSON(http.StatusCreated, echo.Map{
		"success": true,
		"message": "Reservation confirmed successfully",
		"data":    res,
	})
}

func (h *ReservationHandler) GetMyReservations(c echo.Context) error {
	userID := c.Get("userID").(uint)

	res, err := h.srv.GetDriverReservations(userID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"success": false, "message": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"success": true,
		"message": "My reservations retrieved successfully",
		"data":    res,
	})
}

func (h *ReservationHandler) Cancel(c echo.Context) error {
	userID := c.Get("userID").(uint)

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"success": false, "message": "Invalid reservation ID"})
	}

	err = h.srv.CancelReservation(userID, uint(id))
	if err != nil {
		if strings.Contains(err.Error(), "forbidden") {
			return c.JSON(http.StatusForbidden, echo.Map{"success": false, "message": err.Error()})
		}
		return c.JSON(http.StatusBadRequest, echo.Map{"success": false, "message": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"success": true,
		"message": "Reservation cancelled successfully",
	})
}

func (h *ReservationHandler) GetAllReservations(c echo.Context) error {
	res, err := h.srv.GetAllReservations()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, echo.Map{"success": false, "message": err.Error()})
	}

	return c.JSON(http.StatusOK, echo.Map{
		"success": true,
		"message": "All reservations retrieved successfully",
		"data":    res,
	})
}
