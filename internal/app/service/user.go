package service

import (
	"errors"

	"github.com/ervinismu/purplestore/internal/app/model"
	"github.com/ervinismu/purplestore/internal/app/repository"
	"github.com/ervinismu/purplestore/internal/app/schema"
	"github.com/ervinismu/purplestore/internal/app/util"
)

type userService struct {
	UserRepository repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) userService {
	return userService{
		UserRepository: userRepo,
	}
}

func (userService *userService) Register(req schema.RegisterReq) (*schema.RegisterRes, error) {

	//check if email is not exist in database
	_, err := userService.UserRepository.GetByEmail(req.Email)
	if err == nil {
		return nil, errors.New("user with email " + req.Email + " already exist")
	}

	//compose hashed password using BCrypt
	hashedPass := util.HashPasword(req.Password)

	newUser := model.User{
		Email:    req.Email,
		Password: hashedPass,
	}

	savedUser, err := userService.UserRepository.Create(newUser)
	if err != nil {
		return nil, err
	}

	registerRes := schema.RegisterRes{
		Message: "registration success for email " + savedUser.Email,
	}

	return &registerRes, nil
}

func (userService *userService) Login(req schema.LoginReq) (*schema.LoginRes, error) {

	userData, err := userService.UserRepository.GetByEmail(req.Email)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !util.PasswordVerified(req.Password, userData.Password) {
		return nil, errors.New("password does not match")
	}

	//compose JWT
	accessToken := util.BuildJWT(req.Email)

	loginRes := schema.LoginRes{
		Message:     "login success",
		AccessToken: accessToken,
	}

	return &loginRes, nil
}
