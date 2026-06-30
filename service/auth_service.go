package service

import (
	"spotsync-api/dto"
	"spotsync-api/models"
	"spotsync-api/repository"
	"spotsync-api/utils"

	"golang.org/x/crypto/bcrypt"
)

// AuthService contains the business logic for registration and login.
type AuthService struct {
	userRepo       *repository.UserRepository
	jwtSecret      string
	jwtExpiryHours int
}

func NewAuthService(userRepo *repository.UserRepository, jwtSecret string, jwtExpiryHours int) *AuthService {
	return &AuthService{
		userRepo:       userRepo,
		jwtSecret:      jwtSecret,
		jwtExpiryHours: jwtExpiryHours,
	}
}

// Register hashes the password and stores a new user.
func (s *AuthService) Register(req dto.RegisterRequest) (*dto.UserResponse, error) {
	// Ensure email is unique.
	exists, err := s.userRepo.ExistsByEmail(req.Email)
	if err != nil {
		return nil, utils.ErrInternal
	}
	if exists {
		return nil, utils.ErrEmailExists
	}

	// Hash the password with bcrypt (cost 12).
	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), 12)
	if err != nil {
		return nil, utils.ErrInternal
	}

	// Default role is "driver" if not supplied.
	role := req.Role
	if role == "" {
		role = "driver"
	}

	user := &models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashed),
		Role:     role,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, utils.ErrInternal
	}

	return &dto.UserResponse{
		ID:        user.ID,
		Name:      user.Name,
		Email:     user.Email,
		Role:      user.Role,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}, nil
}

// Login verifies credentials and returns a signed JWT.
func (s *AuthService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	user, err := s.userRepo.FindByEmail(req.Email)
	if err != nil {
		return nil, utils.ErrInternal
	}
	// Same generic error whether the email is missing or the password is wrong,
	// so attackers can't tell which emails are registered.
	if user == nil {
		return nil, utils.ErrInvalidCredentials
	}

	// Compare the stored bcrypt hash with the provided password.
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		return nil, utils.ErrInvalidCredentials
	}

	// Sign a JWT carrying the user id and role.
	token, err := utils.GenerateToken(user.ID, user.Role, s.jwtSecret, s.jwtExpiryHours)
	if err != nil {
		return nil, utils.ErrInternal
	}

	return &dto.LoginResponse{
		Token: token,
		User: dto.LoginUser{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Role:  user.Role,
		},
	}, nil
}
