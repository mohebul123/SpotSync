package service

import (
	"errors"

	"github.com/mohebul123/SpotSync/dto"
	"github.com/mohebul123/SpotSync/models"
	"github.com/mohebul123/SpotSync/repository"
)

type ReservationService interface {
	BookSpot(userID uint, req *dto.CreateReservationRequest) (*dto.ReservationWithZoneResponse, error)
	CancelReservation(userID uint, resID uint) error
	GetDriverReservations(userID uint) ([]dto.ReservationWithZoneResponse, error)
	GetAllReservations() ([]dto.ReservationWithZoneResponse, error)
}

type reservationService struct {
	repo     repository.ReservationRepository
	zoneRepo repository.ZoneRepository
}

func NewReservationService(repo repository.ReservationRepository, zoneRepo repository.ZoneRepository) ReservationService {
	return &reservationService{repo: repo, zoneRepo: zoneRepo}
}

func (s *reservationService) BookSpot(userID uint, req *dto.CreateReservationRequest) (*dto.ReservationWithZoneResponse, error) {
	var newRes models.Reservation

	err := s.repo.WithTransaction(func(txRepo repository.ReservationRepository) error {
		zone, err := txRepo.GetZoneForUpdate(req.ZoneID)
		if err != nil {
			return errors.New("parking zone not found")
		}

		exists, err := txRepo.ExistsActiveByLicensePlate(req.LicensePlate)
		if err != nil {
			return err
		}
		if exists {
			return errors.New("this vehicle already has an active reservation")
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

	zone, _ := s.zoneRepo.FindByID(newRes.ZoneID)
	activeCount, _ := s.zoneRepo.GetActiveReservationsCount(zone.ID)
	available := zone.TotalCapacity - int(activeCount)

	return &dto.ReservationWithZoneResponse{
		ID:     newRes.ID,
		UserID: newRes.UserID,
		ZoneID: newRes.ZoneID,
		Zone: dto.ZoneResponse{
			ID:             zone.ID,
			Name:           zone.Name,
			Type:           zone.Type,
			TotalCapacity:  zone.TotalCapacity,
			AvailableSpots: available,
			PricePerHour:   zone.PricePerHour,
			CreatedAt:      zone.CreatedAt,
			UpdatedAt:      zone.UpdatedAt,
		},
		LicensePlate: newRes.LicensePlate,
		Status:       newRes.Status,
		CreatedAt:    newRes.CreatedAt,
		UpdatedAt:    newRes.UpdatedAt,
	}, nil
}

func (s *reservationService) CancelReservation(userID uint, resID uint) error {
	return s.repo.WithTransaction(func(txRepo repository.ReservationRepository) error {
		res, err := txRepo.FindByID(resID)
		if err != nil {
			return errors.New("reservation not found")
		}

		if res.UserID != userID {
			return errors.New("forbidden: you can only cancel your own reservation")
		}

		if res.Status != "active" {
			return errors.New("reservation is already completed or cancelled")
		}

		res.Status = "cancelled"
		return txRepo.Update(res)
	})
}

func (s *reservationService) GetDriverReservations(userID uint) ([]dto.ReservationWithZoneResponse, error) {
	reservations, err := s.repo.FindAllByUserID(userID)
	if err != nil {
		return nil, err
	}

	var res []dto.ReservationWithZoneResponse
	for _, r := range reservations {
		zone, _ := s.zoneRepo.FindByID(r.ZoneID)
		activeCount, _ := s.zoneRepo.GetActiveReservationsCount(zone.ID)
		available := zone.TotalCapacity - int(activeCount)

		res = append(res, dto.ReservationWithZoneResponse{
			ID:     r.ID,
			UserID: r.UserID,
			ZoneID: r.ZoneID,
			Zone: dto.ZoneResponse{
				ID:             zone.ID,
				Name:           zone.Name,
				Type:           zone.Type,
				TotalCapacity:  zone.TotalCapacity,
				AvailableSpots: available,
				PricePerHour:   zone.PricePerHour,
				CreatedAt:      zone.CreatedAt,
				UpdatedAt:      zone.UpdatedAt,
			},
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			CreatedAt:    r.CreatedAt,
			UpdatedAt:    r.UpdatedAt,
		})
	}
	return res, nil
}

func (s *reservationService) GetAllReservations() ([]dto.ReservationWithZoneResponse, error) {
	reservations, err := s.repo.GetAll()
	if err != nil {
		return nil, err
	}

	var res []dto.ReservationWithZoneResponse
	for _, r := range reservations {
		zone, _ := s.zoneRepo.FindByID(r.ZoneID)
		activeCount, _ := s.zoneRepo.GetActiveReservationsCount(zone.ID)
		available := zone.TotalCapacity - int(activeCount)

		res = append(res, dto.ReservationWithZoneResponse{
			ID:     r.ID,
			UserID: r.UserID,
			ZoneID: r.ZoneID,
			Zone: dto.ZoneResponse{
				ID:             zone.ID,
				Name:           zone.Name,
				Type:           zone.Type,
				TotalCapacity:  zone.TotalCapacity,
				AvailableSpots: available,
				PricePerHour:   zone.PricePerHour,
				CreatedAt:      zone.CreatedAt,
				UpdatedAt:      zone.UpdatedAt,
			},
			LicensePlate: r.LicensePlate,
			Status:       r.Status,
			CreatedAt:    r.CreatedAt,
			UpdatedAt:    r.UpdatedAt,
		})
	}
	return res, nil
}
