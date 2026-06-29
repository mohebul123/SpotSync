package repository

import (
	"github.com/mohebul123/SpotSync/models"
	"gorm.io/gorm"
)

type ReservationRepository interface {
	WithTransaction(txFunc func(txRepo ReservationRepository) error) error
	GetZoneForUpdate(zoneID uint) (*models.ParkingZone, error)
	CountActiveReservations(zoneID uint) (int64, error)
	Create(res *models.Reservation) error
	FindByID(id uint) (*models.Reservation, error)
	Update(res *models.Reservation) error
	FindAllByUserID(userID uint) ([]models.Reservation, error)
	GetAll() ([]models.Reservation, error)
	ExistsActiveByLicensePlate(licensePlate string) (bool, error)
}

type reservationRepository struct {
	db *gorm.DB
}

func NewReservationRepository(db *gorm.DB) ReservationRepository {
	return &reservationRepository{db: db}
}

func (r *reservationRepository) WithTransaction(txFunc func(txRepo ReservationRepository) error) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		txRepo := &reservationRepository{db: tx}
		return txFunc(txRepo)
	})
}

func (r *reservationRepository) GetZoneForUpdate(zoneID uint) (*models.ParkingZone, error) {
	var zone models.ParkingZone
	err := r.db.Raw("SELECT * FROM parking_zones WHERE id = ? FOR UPDATE", zoneID).Scan(&zone).Error
	if err != nil {
		return nil, err
	}
	if zone.ID == 0 {
		return nil, gorm.ErrRecordNotFound
	}
	return &zone, nil
}

func (r *reservationRepository) CountActiveReservations(zoneID uint) (int64, error) {
	var count int64
	err := r.db.Model(&models.Reservation{}).Where("zone_id = ? AND status = ?", zoneID, "active").Count(&count).Error
	return count, err
}

func (r *reservationRepository) Create(res *models.Reservation) error {
	return r.db.Create(res).Error
}

func (r *reservationRepository) FindByID(id uint) (*models.Reservation, error) {
	var res models.Reservation
	err := r.db.Preload("Zone").First(&res, id).Error
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *reservationRepository) Update(res *models.Reservation) error {
	return r.db.Save(res).Error
}

func (r *reservationRepository) FindAllByUserID(userID uint) ([]models.Reservation, error) {
	var reservations []models.Reservation

	err := r.db.Preload("Zone").Where("user_id = ?", userID).Order("created_at desc").Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepository) GetAll() ([]models.Reservation, error) {
	var reservations []models.Reservation

	err := r.db.Preload("Zone").Order("id desc").Find(&reservations).Error
	return reservations, err
}

func (r *reservationRepository) ExistsActiveByLicensePlate(licensePlate string) (bool, error) {
	var count int64
	err := r.db.Model(&models.Reservation{}).Where("license_plate = ? AND status = ?", licensePlate, "active").Count(&count).Error
	return count > 0, err
}
