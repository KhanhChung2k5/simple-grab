package repository

import (
	"context"
	"errors"
	"fmt"

	"github.com/KhanhChung2k5/simple-grab/internal/model"
	"gorm.io/gorm"
)

// Error constants for user repository
var (
	ErrUserNotFound = errors.New("user not found")
	ErrEmailExists  = errors.New("email already exists")
)

// UserRepository is the repository for the user model
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetByEmail retrieves a user by their email
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

// GetByID retrieves a user by their ID
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

// CreateWithRoleProfile creates a new user with a role profile
func (r *UserRepository) CreateWithRoleProfile(
	ctx context.Context,
	email string,
	passwordHash string,
	role string,
	phone *string,
) (*model.User, error) {
	var createdUser *model.User

	// create the user
	err := r.db.WithContext(ctx).Transaction(
		func(tx *gorm.DB) error {
			user := &model.User{
				Email:        email,
				PasswordHash: passwordHash,
				Role:         role,
				Phone:        phone,
			}

			// create the user
			if err := tx.Create(user).Error; err != nil {
				if errors.Is(err, gorm.ErrDuplicatedKey) {
					return ErrEmailExists
				}
				return fmt.Errorf("create user: %w", err)
			}
			// check if the role is driver
			if role == model.RoleDriver {
				// create the driver profile
				driver := &model.Driver{
					UserID:      user.ID,
					VehicleType: "car",
				}
				// create the driver profile
				if err := tx.Create(driver).Error; err != nil {
					return fmt.Errorf("create driver profile: %w", err)
				}
			}
			// set the created user
			createdUser = user
			return nil
		},
	)
	if err != nil {
		return nil, err
	}

	return createdUser, nil
}
