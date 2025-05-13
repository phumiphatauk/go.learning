package user

import (
	"go.learning/models"
	"go.learning/utils"
)

type service struct {
	Repository
}

type Service interface {
	GetUserList(queryParams GetUserList) (*GetUserListResponse, error)
	GetUserByID(id uint) (*User, error)
	CreateUser(user CreateUser) error
	UpdateUser(user UpdateUser) error
	DeleteUser(id uint) error
}

func NewService(repository Repository) Service {
	return service{repository}
}

func (s service) GetUserList(queryParams GetUserList) (*GetUserListResponse, error) {
	// Call the repository to get the user list
	users, total, err := s.Repository.GetUserList(queryParams)
	if err != nil {
		return nil, err
	}

	var userList []User
	for _, user := range users {
		userList = append(userList, User{
			ID:        user.ID,
			Email:     user.Email,
			FirstName: user.FirstName,
			LastName:  user.LastName,
			Active:    user.Active,
			CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
		})
	}

	// Create the response
	response := &GetUserListResponse{
		PaginationResponse: PaginationResponse{
			Total: total,
		},
		Data: userList,
	}

	return response, nil
}

func (s service) CreateUser(user CreateUser) error {

	// Generate a hashed password
	hashedPassword, err := utils.GenerateHashedPassword(user.Password)
	if err != nil {
		return err
	}

	// Convert user to models.User
	newUser := &models.User{
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		HashedPassword: hashedPassword,
		Active:         true,
	}

	// Call the repository to create the user
	err = s.Repository.CreateUser(newUser)
	if err != nil {
		return err
	}

	return nil
}

func (s service) GetUserByID(id uint) (*User, error) {
	// Call the repository to get the user by ID
	user, err := s.Repository.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	// Convert models.User to User
	return &User{
		ID:        user.ID,
		Email:     user.Email,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Active:    user.Active,
		CreatedAt: user.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: user.UpdatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (s service) UpdateUser(user UpdateUser) error {
	// Convert user to models.User
	updatedUser := &models.User{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Active:    user.Active,
	}

	// Call the repository to update the user
	err := s.Repository.UpdateUser(updatedUser)
	if err != nil {
		return err
	}

	return nil
}

func (s service) DeleteUser(id uint) error {
	// Call the repository to delete the user
	err := s.Repository.DeleteUser(id)
	if err != nil {
		return err
	}

	return nil
}
