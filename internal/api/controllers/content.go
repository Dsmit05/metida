package controllers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/Dsmit05/metida/internal/models"

	"github.com/Dsmit05/metida/internal/api/response"
	"github.com/gin-gonic/gin"
)

type contentRepositoryI interface {
	CreatContent(email string, name string, description string) error
	ReadContent(email string, id int32) (*models.Content, error)
}

// UserContent defines the content controller methods
type UserContent struct {
	db contentRepositoryI
}

func NewWallEditorials(db contentRepositoryI) *UserContent {
	return &UserContent{db}
}

type CreateContentInput struct {
	Name        string `json:"name" binding:"required" example:"My First Content"`
	Description string `json:"description" binding:"required" example:"New content..."`
}

// @Summary CreateContent
// @Tags content
// @Description Create Content
// @ID protected-create-content
// @Accept json
// @Produce json
// @Param input body CreateContentInput true "credentials"
// @Success 200 {object} response.Success
// @Failure 400 {string} string "error"
// @Failure 417 {object} response.Error
// @Security ApiKeyAuth
// @Router /lk/content [POST]
func (o *UserContent) CreateContent(c *gin.Context) {
	var inputData CreateContentInput

	if err := c.ShouldBindJSON(&inputData); err != nil {
		response.GinError(c, http.StatusBadRequest, response.CodeInvalidParams, "bad data, try again", nil)
		return
	}

	email, ok := c.Get("email")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user without email"})
	}

	if err := o.db.CreatContent(email.(string), inputData.Name, inputData.Description); err != nil {
		response.GinError(c, http.StatusBadRequest, response.CodeBadRequest, err.Error(), err)
		return
	}

	response.GinSuccess(c, http.StatusOK, response.CodeOk, "", "Content created")
}

// @Summary Show Content
// @Tags content
// @Description Show Content
// @ID protected-show-content
// @Accept json
// @Produce json
// @Param id path int true "id"
// @Success 200 {object} models.Content
// @Failure 400 {object} response.Error
// @Failure 417 {object} response.Error
// @Security ApiKeyAuth
// @Router /lk/content/{id} [GET]
func (o *UserContent) ShowContent(c *gin.Context) {
	contentId := c.Param("id")
	if contentId == "" {
		err := fmt.Errorf("blog not selected")
		response.GinError(c, http.StatusBadRequest, response.CodeBadRequest, "you need to select the content", err)
	}

	email, ok := c.Get("email")
	if !ok {
		c.AbortWithStatusJSON(http.StatusBadRequest, gin.H{"error": "user without email"})
	}

	id, err := strconv.Atoi(contentId)
	if err != nil {
		err = fmt.Errorf("blog not selected")
		response.GinError(c, http.StatusBadRequest, response.CodeBadRequest, "input content number", err)
	}

	content, err := o.db.ReadContent(email.(string), int32(id))
	if err != nil {
		response.GinError(c, http.StatusBadRequest, response.CodeDBError, err.Error(), err)
		return
	}

	c.JSON(http.StatusOK, content)
}
