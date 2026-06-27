package service

import (
	"errors"
	"math"

	"github.com/mohebul123/SpotSync/dto"
	"github.com/mohebul123/SpotSync/models"
	"github.com/mohebul123/SpotSync/repository"
)

type ReservationService interface {
	BookSpot(userID uint, req *dto.CreateReservationRequest) (*dto.ReservationResponse, error)
	CancelReservation(userID uint, resID uint) (*dto.ReservationResponse, error)
	GetDriverReservations(userID uint) ([]dto.ReservationResponse, error)
}

type reservationService struct {
	repo     repository.ReservationRepository
	zoneRepo repository.ZoneRepository // To fetch price per hour safely
}

func NewReservationService(repo repository.ReservationRepository, zoneRepo repository.ZoneRepository) ReservationService {
	return &reservationService{repo: repo, zoneRepo: zoneRepo}
}

func (s *reservationService) BookSpot(userID uint, req *dto.CreateReservationRequest) (*dto.ReservationResponse, error) {
	var newRes models.Reservation

	err := s.repo.WithTransaction(func(txRepo repository.ReservationRepository) error {
		zone, err := txRepo.GetZoneForUpdate(req.ZoneID)
		if err != nil {
			return errors.New("parking zone not found")
		}

		activeCount, err := txRepo.CountActiveReservations(req.ZoneID)
		if err != nil {
			return err
		}

		if int(activeCount) >= zone.TotalCapacity {
			return errors.New("no available spots in this parking zone")
		}

		newRes = models.Reservation{
			UserID:       userID,
			ZoneID:       req.ZoneID,
			LicensePlate: req.LicensePlate,
			Status:       "active",
		}

		return txRepo.Create(&newRes)
	})

	if err != nil {
		return nil, err
	}

	return &dto.ReservationResponse{
		ID:           newRes.ID,
		UserID:       newRes.UserID,
		ZoneID:       newRes.ZoneID,
		LicensePlate: newRes.LicensePlate,
		Status:       newRes.Status,
		StartTime:    newRes.CreatedAt,
		EndTime:      newRes.UpdatedAt,
		TotalCost:    0.0,
	}, nil
}

func (s *reservationService) CancelReservation(userID uint, resID uint) (*dto.ReservationResponse, error) {
	var updatedRes *models.Reservation
	var totalCost float64

	err := s.repo.WithTransaction(func(txRepo repository.ReservationRepository) error {
		res, err := txRepo.FindByIDAndUserID(resID, userID)
		if err != nil {
			return errors.New("active reservation not found")
		}

		if res.Status != "active" {
			return errors.New("reservation is already completed or cancelled")
		}

		zone, err := txRepo.GetZoneForUpdate(res.ZoneID)
		if err != nil {
			return err
		}

		res.Status = "cancelled"
		if err := txRepo.Update(res); err != nil {
			return err
		}

		// Calculate cost: EndTime (UpdatedAt) - StartTime (CreatedAt)
		duration := res.UpdatedAt.Sub(res.CreatedAt)
		hoursSpent := duration.Hours()

		if hoursSpent < 0.02 { // Grace period under 1 min
			totalCost = 0.0
		} else {
			totalCost = math.Ceil(hoursSpent) * zone.PricePerHour
		}

		updatedRes = res
		return nil
	})

	if err != nil {
		return nil, err
	}

	return &dto.ReservationResponse{
		ID:           updatedRes.ID,
		UserID:       updatedRes.UserID,
		ZoneID:       updatedRes.ZoneID,
		LicensePlate: updatedRes.LicensePlate,
		Status:       updatedRes.Status,
		StartTime:    updatedRes.CreatedAt,
		EndTime:      updatedRes.UpdatedAt,
		TotalCost:    totalCost,
	}, nil
}

func (s *reservationService) GetDriverReservations(userID uint) ([]dto.ReservationResponse, error) {
	reservations, err := s.repo.FindAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	var res []dto.ReservationResponse
	for _, r := range reservations {
		var totalCost float64
		// Fetch zone details to get the rate for dynamic calculation
		zone, _ := s.zoneRepo.FindByID(r.ZoneID)

		if r.Status != "active" && zone != nil {
			duration := r.UpdatedAt.Sub(r.CreatedAt)
			hours := duration.Hours()
			if hours > 0.02 {
				totalCost = math.Ceil(hours) * zone.PricePerHour
			}
		}

		res = append(res, dto.ReservationResponse{
			ID:           r.ID,
			UserID:       r.UserID,
			ZoneID:       r.ZoneID,
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			StartTime:    r.CreatedAt,
			EndTime:      r.UpdatedAt,
			TotalCost:    totalCost,
		})
	}
	return res, nil
}
