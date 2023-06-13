package controller

import (
	"net/http"

	"github.com/ervinismu/purplestore/internal/app/schema"
	"github.com/gin-gonic/gin"
)

type UserService interface {
	Register(req schema.RegisterReq) (*schema.RegisterRes, error)
	Login(req schema.LoginReq) (*schema.LoginRes, error)
}

type userController struct {
	UserService UserService
}

func NewUserController(userService UserService) userController {
	return userController{
		UserService: userService,
	}
}

func (ctrl *userController) Register(c *gin.Context) {
	var req schema.RegisterReq

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	userData, err := ctrl.UserService.Register(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": userData,
	})
}

func (ctrl *userController) Login(c *gin.Context) {
	var req schema.LoginReq

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	loginData, err := ctrl.UserService.Login(req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": loginData,
	})
}
