package handler

import (
	"bwastartup/helper"
	"bwastartup/user"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
}

func NewUserHandler(userService user.Service) *userHandler {
	return &userHandler{userService}
}

func (h *userHandler) RegisterUser(c *gin.Context) {

	var input user.RegisterUserInput

	err := c.ShouldBindJSON(&input)
	if err != nil {

		errors := helper.FormatValidationError(err)

		errorMassage := gin.H{"errors": errors}

		response := helper.APIRespone("Register Account Failed", http.StatusBadRequest, "error", errorMassage)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	newUser, err := h.userService.RegisterUser(input)

	if err != nil {
		response := helper.APIRespone("Register Account Failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(newUser, "token Token Roken")

	response := helper.APIRespone("Account han been created", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) Login(c *gin.Context) {
	var input user.LoginInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMassage := gin.H{"errors": errors}
		response := helper.APIRespone("Login Failed", http.StatusBadRequest, "error", errorMassage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	loogedinUser, err := h.userService.Login(input)

	if err != nil {
		errorMassage := gin.H{"errors": err.Error()}
		response := helper.APIRespone("Login Failed", http.StatusBadRequest, "error", errorMassage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(loogedinUser, "token token token")

	response := helper.APIRespone("Succesfully Loogedin", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}
