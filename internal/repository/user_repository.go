package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/KhanhChung2k5/simple-grab/internal/model"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailExists  = errors.New("email already exists")
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(ctx context.Context, email, passwordHash, role string, phone *string) (*model.User, error) {
	user := &model.User{
		Email:        email,
		Phone:        phone,
		PasswordHash: passwordHash,
		Role:         role,
	}

	err := r.db.WithContext(ctx).Create(user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			return nil, ErrEmailExists
		}
		return nil, fmt.Errorf("create user: %w", err)
	}

	return user, nil
}

func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by email: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("get user by id: %w", err)
	}

	return &user, nil
}

func (r *UserRepository) CreateDriverProfile(ctx context.Context, userID string) (*model.Driver, error) {
	driver := &model.Driver{
		UserID:      userID,
		VehicleType: "car",
	}

	err := r.db.WithContext(ctx).Create(driver).Error
	if err != nil {
		return nil, fmt.Errorf("create driver profile: %w", err)
	}

	return driver, nil
}
