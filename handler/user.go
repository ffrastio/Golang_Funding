package handler

import (
	"bwastartup/auth"
	"bwastartup/helper"
	"bwastartup/user"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService user.Service
	authService auth.Service
}

func NewUserHandler(userService user.Service, authService auth.Service) *userHandler {
	return &userHandler{userService, authService}
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

	token, err := h.authService.GenerateToken(newUser.ID)
	if err != nil {
		response := helper.APIRespone("Register Account Failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(newUser, token)

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

	token, err := h.authService.GenerateToken(loogedinUser.ID)
	if err != nil {
		response := helper.APIRespone("Login Failed", http.StatusBadRequest, "error", nil)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	formatter := user.FormatUser(loogedinUser, token)

	response := helper.APIRespone("Succesfully Loogedin", http.StatusOK, "success", formatter)

	c.JSON(http.StatusOK, response)
}

func (h *userHandler) CheckEmailAvaibility(c *gin.Context) {
	// ada input dari user
	// input email di mapping ke struct input

	var input user.CheckEmailInput

	err := c.ShouldBindJSON(&input)
	if err != nil {
		errors := helper.FormatValidationError(err)
		errorMassage := gin.H{"errors": errors}
		response := helper.APIRespone("Email Checking Failed", http.StatusBadRequest, "error", errorMassage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	isEmailAvailable, err := h.userService.IsEmailAvailable(input)
	if err != nil {
		errorMassage := gin.H{"errors": "Server Error"}
		response := helper.APIRespone("Email Checking Failed", http.StatusBadRequest, "error", errorMassage)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{
		"is_available": isEmailAvailable,
	}

	metaMessage := "Email has been Registered"

	if isEmailAvailable {
		metaMessage = "Email is Available"
	}

	response := helper.APIRespone(metaMessage, http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)

}

func (h *userHandler) UploadAvatar(c *gin.Context) {
	// input dari user
	// simpan gambar dalam folder "/images"
	// service memanggil repository
	// JWT (sementara hardcode seakan-akan akun sudah login)
	// repo mengambil ID user ID = 1
	// repo update data user simpan lokasi file
	file, err := c.FormFile("avatar")
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIRespone("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}
	userID := 1

	path := fmt.Sprintf("images/%d-%s", userID, file.Filename)

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIRespone("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	_, err = h.userService.SaveAvatar(userID, path)
	if err != nil {
		data := gin.H{"is_uploaded": false}
		response := helper.APIRespone("Failed to upload avatar image", http.StatusBadRequest, "error", data)
		c.JSON(http.StatusBadRequest, response)
		return
	}

	data := gin.H{"is_uploaded": true}
	response := helper.APIRespone("Success to upload avatar image", http.StatusOK, "success", data)
	c.JSON(http.StatusOK, response)

}
