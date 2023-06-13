package repository

import (
	"errors"

	"github.com/ervinismu/purplestore/internal/app/model"
	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	Create(user model.User) (*model.User, error)
	GetByEmail(email string) (*model.User, error)
}

type userRepo struct {
	DB *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) UserRepository {
	return &userRepo{
		DB: db,
	}
}

// Create implements UserRepository
func (repo *userRepo) Create(user model.User) (*model.User, error) {
	sqlStatement := `
		INSERT INTO users (email, password)
		VALUES ($1, $2) 
	`

	_, err := repo.DB.Exec(sqlStatement, user.Email, user.Password)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

// GetByEmail implements UserRepository
func (repo *userRepo) GetByEmail(email string) (*model.User, error) {
	var (
		data         model.User
		sqlStatement = `
			SELECT id, email, password
			FROM users
			WHERE email = $1
			LIMIT 1
		`
	)

	err := repo.DB.QueryRowx(sqlStatement, email).StructScan(&data)
	if err != nil {
		return nil, err
	}

	if data.Id == 0{
		return nil, errors.New("data not found")
	}

	return &data, nil
}
