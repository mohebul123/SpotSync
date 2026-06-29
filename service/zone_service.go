package service

import (
	"errors"

	"github.com/mohebul123/SpotSync/dto"
	"github.com/mohebul123/SpotSync/models"
	"github.com/mohebul123/SpotSync/repository"
)

type ZoneService interface {
	CreateZone(req *dto.CreateZoneRequest) (*dto.ZoneResponse, error)
	GetAllZones() ([]dto.ZoneResponse, error)
	GetZoneByID(id uint) (*dto.ZoneResponse, error)
	UpdateZone(id uint, req *dto.UpdateZoneRequest) (*dto.ZoneResponse, error)
	DeleteZone(id uint) error
}

type zoneService struct {
	repo repository.ZoneRepository
}

func NewZoneService(repo repository.ZoneRepository) ZoneService {
	return &zoneService{repo: repo}
}

func (s *zoneService) CreateZone(req *dto.CreateZoneRequest) (*dto.ZoneResponse, error) {
	zone := &models.ParkingZone{
		Name:          req.Name,
		Type:          req.Type,
		TotalCapacity: req.TotalCapacity,
		PricePerHour:  req.PricePerHour,
	}

	if err := s.repo.Create(zone); err != nil {
		return nil, err
	}

	return &dto.ZoneResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: zone.TotalCapacity,
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt,
		UpdatedAt:      zone.UpdatedAt,
	}, nil
}

func (s *zoneService) GetAllZones() ([]dto.ZoneResponse, error) {
	zones, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var res []dto.ZoneResponse
	for _, zone := range zones {
		activeCount, _ := s.repo.GetActiveReservationsCount(zone.ID)
		available := zone.TotalCapacity - int(activeCount)

		res = append(res, dto.ZoneResponse{
			ID:             zone.ID,
			Name:           zone.Name,
			Type:           zone.Type,
			TotalCapacity:  zone.TotalCapacity,
			AvailableSpots: available,
			PricePerHour:   zone.PricePerHour,
			CreatedAt:      zone.CreatedAt,
			UpdatedAt:      zone.UpdatedAt,
		})
	}
	return res, nil
}

func (s *zoneService) GetZoneByID(id uint) (*dto.ZoneResponse, error) {
	zone, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}

	activeCount, _ := s.repo.GetActiveReservationsCount(zone.ID)
	available := zone.TotalCapacity - int(activeCount)

	return &dto.ZoneResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: available,
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt,
		UpdatedAt:      zone.UpdatedAt,
	}, nil
}

func (s *zoneService) UpdateZone(id uint, req *dto.UpdateZoneRequest) (*dto.ZoneResponse, error) {
	zone, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.New("parking zone not found")
	}

	zone.Name = req.Name
	zone.Type = req.Type
	zone.TotalCapacity = req.TotalCapacity
	zone.PricePerHour = req.PricePerHour

	if err := s.repo.Update(zone); err != nil {
		return nil, err
	}

	activeCount, _ := s.repo.GetActiveReservationsCount(zone.ID)
	available := zone.TotalCapacity - int(activeCount)

	return &dto.ZoneResponse{
		ID:             zone.ID,
		Name:           zone.Name,
		Type:           zone.Type,
		TotalCapacity:  zone.TotalCapacity,
		AvailableSpots: available,
		PricePerHour:   zone.PricePerHour,
		CreatedAt:      zone.CreatedAt,
		UpdatedAt:      zone.UpdatedAt,
	}, nil
}

func (s *zoneService) DeleteZone(id uint) error {
	_, err := s.repo.FindByID(id)
	if err != nil {
		return errors.New("parking zone not found")
	}
	return s.repo.Delete(id)
}
