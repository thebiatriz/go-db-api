package handlers

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/thebiatriz/go-db-api/internal/models"
	"github.com/thebiatriz/go-db-api/internal/usecases"
)

type UserHandler struct {
	userUsecase usecases.UserUsecase
}

func NewUserHandler(userUsecase usecases.UserUsecase) UserHandler {
	return UserHandler{
		userUsecase: userUsecase,
	}
}

func (u UserHandler) Login(c *gin.Context) {
	var req models.LoginRequest
	err := c.BindJSON(&req)

	if err != nil {
		handlePayloadErrors(c, err)
		return
	}

	token, err := u.userUsecase.Login(req.Email, req.Password)

	if err != nil {
		handleDomainError(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, gin.H{"token": token})
}

func (u UserHandler) GetUsers(c *gin.Context) {
	pageStr := c.DefaultQuery("page", "1")
	limitStr := c.DefaultQuery("limit", "10")
	queryName := c.DefaultQuery("name", "")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	if limit > 50 {
		limit = 50
	}

	users, err := u.userUsecase.GetUsers(page, limit, queryName)

	if err != nil {
		handleDomainError(c, err)
		return
	}
	c.IndentedJSON(http.StatusOK, users)
}

func (u UserHandler) GetUserById(c *gin.Context) {
	id := c.Param("id")

	userId, err := strconv.Atoi(id)

	if err != nil {
		response := models.Response{
			Message: msgInvalidIDNumeric,
		}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	user, err := u.userUsecase.GetUserById(userId)

	if err != nil {
		handleDomainError(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, user)
}

func (u UserHandler) CreateUser(c *gin.Context) {
	var req models.CreateUserRequest

	err := c.BindJSON(&req)

	if err != nil {
		handlePayloadErrors(c, err)
		return
	}

	req.Name = strings.TrimSpace(req.Name)
	req.Email = strings.TrimSpace(req.Email)

	userToCreate := models.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
	}

	newUser, err := u.userUsecase.CreateUser(userToCreate)

	if err != nil {
		handleDomainError(c, err)
		return
	}

	c.IndentedJSON(http.StatusCreated, newUser)
}

func (u UserHandler) DeleteUser(c *gin.Context) {
	targetId := c.Param("id")

	targetUserId, err := strconv.Atoi(targetId)

	if err != nil {
		response := models.Response{
			Message: msgInvalidIDNumeric,
		}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	requesterIdStr, exists := c.Get("userId")
	if !exists {
		response := models.Response{
			Message: msgUserIDMissing,
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	requesterId, ok := requesterIdStr.(int)
	if !ok {
		response := models.Response{
			Message: msgInvalidIDFormat,
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	requesterRoleStr, exists := c.Get("userRole")
	if !exists {
		response := models.Response{
			Message: msgUserRoleMissing,
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	requesterRole, ok := requesterRoleStr.(string)
	if !ok {
		response := models.Response{
			Message: msgInvalidRoleFormat,
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	err = u.userUsecase.DeleteUser(targetUserId, requesterId, requesterRole)

	if err != nil {
		handleDomainError(c, err)
		return
	}

	c.Status(http.StatusNoContent)

}

func (u UserHandler) UpdateUser(c *gin.Context) {
	var req models.UpdateUserRequest
	targetId := c.Param("id")

	targetUserId, err := strconv.Atoi(targetId)

	if err != nil {
		response := models.Response{
			Message: msgInvalidIDNumeric,
		}
		c.IndentedJSON(http.StatusBadRequest, response)
		return
	}

	requesterIdStr, exists := c.Get("userId")
	if !exists {
		response := models.Response{
			Message: msgUserIDMissing,
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	requesterId, ok := requesterIdStr.(int)
	if !ok {
		response := models.Response{
			Message: msgInvalidIDFormat,
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	requesterRoleStr, exists := c.Get("userRole")
	if !exists {
		response := models.Response{
			Message: msgUserRoleMissing,
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	requesterRole, ok := requesterRoleStr.(string)
	if !ok {
		response := models.Response{
			Message: msgInvalidRoleFormat,
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	err = c.BindJSON(&req)

	if err != nil {
		handlePayloadErrors(c, err)
		return
	}

	if req.Name != nil {
		*req.Name = strings.TrimSpace(*req.Name)
	}

	if req.Email != nil {
		*req.Email = strings.TrimSpace(*req.Email)
	}

	updatedUser, err := u.userUsecase.UpdateUser(targetUserId, requesterId, requesterRole, req)

	if err != nil {
		handleDomainError(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, updatedUser)
}

func (u UserHandler) GetMe(c *gin.Context) {
	userIdStr, exists := c.Get("userId")

	if !exists {
		response := models.Response{
			Message: msgUserIDMissing,
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	userId, ok := userIdStr.(int)

	if !ok {
		response := models.Response{
			Message: msgInvalidIDFormat,
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	currentUser, err := u.userUsecase.GetUserById(userId)

	if err != nil {
		handleDomainError(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, currentUser)
}

func (u UserHandler) DeleteMe(c *gin.Context) {
	userIdStr, exists := c.Get("userId")

	if !exists {
		response := models.Response{
			Message: msgUserIDMissing,
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	userId, ok := userIdStr.(int)

	if !ok {
		response := models.Response{
			Message: msgInvalidIDFormat,
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	userRoleStr, exists := c.Get("userRole")
	if !exists {
		response := models.Response{
			Message: msgUserRoleMissing,
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	userRole, ok := userRoleStr.(string)
	if !ok {
		response := models.Response{
			Message: msgInvalidRoleFormat,
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	err := u.userUsecase.DeleteUser(userId, userId, userRole)

	if err != nil {
		handleDomainError(c, err)
		return
	}

	c.Status(http.StatusNoContent)
}

func (u UserHandler) UpdateMe(c *gin.Context) {
	var req models.UpdateUserRequest
	userIdStr, exists := c.Get("userId")

	if !exists {
		response := models.Response{
			Message: msgUserIDMissing,
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	userId, ok := userIdStr.(int)
	if !ok {
		response := models.Response{
			Message: msgInvalidIDFormat,
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	userRoleStr, exists := c.Get("userRole")
	if !exists {
		response := models.Response{
			Message: msgUserRoleMissing,
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	userRole, ok := userRoleStr.(string)
	if !ok {
		response := models.Response{
			Message: msgInvalidRoleFormat,
		}
		c.IndentedJSON(http.StatusInternalServerError, response)
		return
	}

	err := c.BindJSON(&req)

	if err != nil {
		handlePayloadErrors(c, err)
		return
	}

	if req.Name != nil {
		*req.Name = strings.TrimSpace(*req.Name)
	}

	if req.Email != nil {
		*req.Email = strings.TrimSpace(*req.Email)
	}

	updatedUser, err := u.userUsecase.UpdateUser(userId, userId, userRole, req)

	if err != nil {
		handleDomainError(c, err)
		return
	}

	c.IndentedJSON(http.StatusOK, updatedUser)
}
