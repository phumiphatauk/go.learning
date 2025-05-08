package user

import "go.learning/models"

type service struct {
	Repository
}

type Service interface {
	GetUserByID(id uint) (*User, error)
	CreateUser(user CreateUser) error
	UpdateUser(user UpdateUser) error
	DeleteUser(id uint) error
}

func NewService(repository Repository) Service {
	return service{repository}
}

func (s service) CreateUser(user CreateUser) error {
	// Convert user to models.User
	newUser := &models.User{
		Email:          user.Email,
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		HashedPassword: user.Password, // Hash the password here
		Active:         true,
	}

	// Call the repository to create the user
	err := s.Repository.CreateUser(newUser)
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
		DeletedAt: user.DeletedAt.Time.Format("2006-01-02 15:04:05"),
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
