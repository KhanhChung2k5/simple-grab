package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/KhanhChung2k5/simple-grab/internal/model"
	"github.com/KhanhChung2k5/simple-grab/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)


//Error constants (invalid credentials, invalid role)
var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrInvalidRole        = errors.New("invalid role")
)

//AuthService is a service for the auth model
type AuthService struct {
	users     *repository.UserRepository
	jwtSecret []byte
	jwtExpiry time.Duration
}

//NewAuthService creates a new auth service
func NewAuthService(users *repository.UserRepository, jwtSecret string, jwtExpiryHours int) *AuthService {
	return &AuthService{
		users:     users,
		jwtSecret: []byte(jwtSecret),
		jwtExpiry: time.Duration(jwtExpiryHours) * time.Hour,
	}
}

//Claims is a struct for the claims
type Claims struct {
	UserID string `json:"user_id"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

//Register registers a new user
func (s *AuthService) Register(ctx context.Context, req model.RegisterRequest) (*model.RegisterResponse, error) {
	//Check if the role is valid
	if req.Role != model.RoleRider && req.Role != model.RoleDriver {
		return nil, ErrInvalidRole
	}

	//Hash the password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	//Create the user
	user, err := s.users.Create(ctx, req.Email, string(hash), req.Role, req.Phone)
	if err != nil {
		return nil, err
	}

	//Create the driver profile if the role is driver
	if req.Role == model.RoleDriver {
		if _, err := s.users.CreateDriverProfile(ctx, user.ID); err != nil {
			return nil, err
		}
	}

	//Return the user
	return &model.RegisterResponse{
		ID:        user.ID,
		Email:     user.Email,
		Role:      user.Role,
		Phone:     user.Phone,
		CreatedAt: user.CreatedAt,
	}, nil
}

//Login logs in a user
func (s *AuthService) Login(ctx context.Context, req model.LoginRequest) (*model.LoginResponse, error) {
	//Get the user by email
	user, err := s.users.GetByEmail(ctx, req.Email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, err
	}

	//Compare the password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	//Generate the token
	token, err := s.generateToken(user)
	if err != nil {
		return nil, err
	}

	user.PasswordHash = ""
	//Return the token
	return &model.LoginResponse{
		Token: token,
		ExpiresIn:   int(s.jwtExpiry.Seconds()),
	}, nil
}

//GetMe gets the user by id
func (s *AuthService) GetMe(ctx context.Context, userID string) (*model.User, error) {
	user, err := s.users.GetByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	//Clear the password hash
	user.PasswordHash = ""
	return user, nil
}

//ParseToken parses a token string
func (s *AuthService) ParseToken(tokenString string) (*Claims, error) {
	//Parse the token
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}

	//Get the claims
	claims, ok := token.Claims.(*Claims)
	//Check if the token is valid
	if !ok || !token.Valid {
		//Return an error if the token is invalid
		return nil, errors.New("invalid token")
	}

	return claims, nil
}

//GenerateToken generates a token for a user
func (s *AuthService) generateToken(user *model.User) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: user.ID,
		Role:   user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   user.ID,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(s.jwtExpiry)),
		},
	}

	//Create the token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtSecret)
}
