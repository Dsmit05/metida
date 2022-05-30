package controllers

import (
	"fmt"
	"net/http"
	"time"

	"github.com/Dsmit05/metida/internal/api/validation"
	"github.com/Dsmit05/metida/internal/cryptography"
	"github.com/Dsmit05/metida/internal/models"

	"github.com/Dsmit05/metida/internal/api/response"
	"github.com/Dsmit05/metida/internal/consts"
	"github.com/Dsmit05/metida/internal/logger"
	"github.com/gin-gonic/gin"
)

type userRepositoryI interface {
	CreateUser(name string, password string, email string, role string) error
	ReadUser(email string) (*models.User, error)
	CreateSession(email string, refreshToken string, userAgent string, ip string, expiresIn int64) error
	ReadSession(email string, userAgent string, ip string) (*models.Session, error)
	UpdateSessionTokenOnly(refreshToken string, newRefreshToken string, expiresIn int64) error
	ReadEmailRoleWithRefreshToken(refreshToken string) (*models.UserEmailRole, error)
}

type tokensI interface {
	CreateToken(email string, role string, ttl time.Duration) (string, error)
	ParseToken(inputToken string) (email string, role string, err error)
	CreateRefreshToken() (string, error)
}

// UserAuth defines the user controller methods
type UserAuth struct {
	db    userRepositoryI
	token tokensI
}

func NewUserAuth(db userRepositoryI, token tokensI) *UserAuth {
	return &UserAuth{db: db, token: token}
}

type CreateUserInput struct {
	Username string `json:"username" binding:"required" example:"Ivan"`
	Password string `json:"password" binding:"required" example:"Q@werty1_23"`
	Email    string `json:"email" binding:"required" example:"ivashka2015@gmail.com"`
}

// @Summary Sign Up
// @Tags auth
// @Description create account
// @ID create-account
// @Accept json
// @Produce json
// @Param input body CreateUserInput true "credentials"
// @Success 200 {object} response.Success
// @Failure 400 {object} response.Error
// @Router /auth/sign-up [post]
func (o *UserAuth) CreateUser(c *gin.Context) {
	// Validate input
	var inputData CreateUserInput
	if err := c.ShouldBindJSON(&inputData); err != nil {
		response.GinError(c, http.StatusBadRequest, response.CodeInvalidParams, "incorrect data, try again", nil)
		return
	}

	if !o.validateUserData(c, inputData) {
		return
	}

	// hashing input password
	passwordHash, err := cryptography.HashPassword(inputData.Password)
	if err != nil {
		response.GinError(c, http.StatusBadRequest, response.CodeInvalidParams, "incorrect password", nil)
		return
	}

	if err = o.db.CreateUser(inputData.Username, passwordHash, inputData.Email, consts.RoleUser); err != nil {
		response.GinError(c, http.StatusBadRequest, response.CodeBadRequest, err.Error(), err)
		return
	}

	// Здесь нужно возращать OK и предлагать подтвердить почту
	// Дальше пользователь должен подтвердить почту, после чего создаем ему сессию
	// Todo: данную реализацию можно сделать в следующих версиях апи
	// Для примера сделаем сквозную, без транзакций

	ip, agent := o.getIPandUserAgent(c)

	// Create refresh Token
	rToken, err := o.token.CreateRefreshToken()
	if err != nil {
		response.GinError(c, http.StatusInternalServerError, response.CodeCryptoError, "", nil)
		return
	}
	err = o.db.CreateSession(inputData.Email, rToken, agent, ip, time.Now().Add(consts.RefreshTokenTTL).Unix())
	if err != nil {
		response.GinError(c, http.StatusBadRequest, response.CodeBadRequest, err.Error(), err)
		return
	}

	// Create access token
	aToken, err := o.token.CreateToken(inputData.Email, consts.RoleUser, consts.AccessTokenTTL)
	if err != nil {
		response.GinError(c, http.StatusInternalServerError, response.CodeCryptoError, "", nil)
		return
	}

	response.GinSuccess(c,
		http.StatusOK, response.CodeOk,
		gin.H{"aToken": aToken, "rToken": rToken}, "Create New User")
}

type AuthenticationUserInput struct {
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// @Summary Sign In
// @Tags auth
// @Description log in account
// @ID login-account
// @Accept json
// @Produce json
// @Param input body AuthenticationUserInput true "credentials"
// @Success 200 {object} response.Success
// @Failure 400 {object} response.Error
// @Router /auth/sign-in [post]
func (o *UserAuth) AuthenticationUser(c *gin.Context) {
	var inputData AuthenticationUserInput

	if err := c.ShouldBindJSON(&inputData); err != nil {
		response.GinError(c, http.StatusBadRequest, response.CodeInvalidParams, "bad data, try again", nil)
		return
	}

	user, err := o.db.ReadUser(inputData.Email)
	if err != nil {
		response.GinError(c, http.StatusBadRequest, response.CodeDBError, err.Error(), err)
		return
	}

	// check password with hash
	if !cryptography.CheckPassword(user.Password, inputData.Password) {
		err = fmt.Errorf("wrong password")
		response.GinError(c, http.StatusUnauthorized, response.CodeBadRequest, err.Error(), err)
		return
	}

	newRefreshToken, err := o.token.CreateRefreshToken()
	if err != nil {
		response.GinError(c, http.StatusInternalServerError, response.CodeCryptoError, "", nil)
		return
	}

	ip, agent := o.getIPandUserAgent(c)

	// при каждом логине создаем новую сессию
	err = o.db.CreateSession(
		inputData.Email, newRefreshToken, agent, ip, time.Now().Add(consts.RefreshTokenTTL).Unix())
	if err != nil {
		response.GinError(c, http.StatusBadRequest, response.CodeBadRequest, err.Error(), err)
		return
	}

	// Create access token
	aToken, err := o.token.CreateToken(inputData.Email, user.Role, consts.AccessTokenTTL)
	if err != nil {
		response.GinError(c, http.StatusInternalServerError, response.CodeCryptoError, "", nil)
		return
	}

	response.GinSuccess(c, http.StatusOK, response.CodeOk,
		gin.H{"aToken": aToken, "rToken": newRefreshToken}, "Authentication well")

}

type RefreshTokenInput struct {
	RefreshToken string `json:"rtoken" binding:"required"`
}

// @Summary Refresh token
// @Tags auth
// @Description refresh access token
// @ID refresh-token
// @Accept json
// @Produce json
// @Param input body RefreshTokenInput true "credentials"
// @Success 200 {object} response.Success
// @Failure 400 {object} response.Error
// @Router /auth/refresh [post]
func (o *UserAuth) RefreshTokenUser(c *gin.Context) {
	var inputData RefreshTokenInput

	if err := c.ShouldBindJSON(&inputData); err != nil {
		response.GinError(c, http.StatusBadRequest, response.CodeInvalidParams, "bad data, try again", err)
		return
	}

	userData, err := o.db.ReadEmailRoleWithRefreshToken(inputData.RefreshToken)
	if err != nil {
		response.GinError(c, http.StatusUnauthorized, response.CodeBadRequest, err.Error(), err)
		return
	}

	// Check ttl refresh token
	if userData.ExpiresIn <= time.Now().Unix() {
		response.GinError(c, http.StatusExpectationFailed, response.CodeUnknownUser, "Please log in", err)
		return
	}

	// Create access token
	aToken, err := o.token.CreateToken(userData.Email, userData.Role, consts.AccessTokenTTL)
	if err != nil {
		response.GinError(c, http.StatusInternalServerError, response.CodeCryptoError, "", nil)
		return
	}

	rToken, err := o.token.CreateRefreshToken()
	if err != nil {
		response.GinError(c, http.StatusInternalServerError, response.CodeCryptoError, "", nil)
		return
	}

	err = o.db.UpdateSessionTokenOnly(
		inputData.RefreshToken, rToken, time.Now().Add(consts.RefreshTokenTTL).Unix())
	if err != nil {
		response.GinError(c, http.StatusBadRequest, response.CodeBadRequest, "", nil)
		return
	}

	response.GinSuccess(c,
		http.StatusOK, response.CodeOk,
		gin.H{"aToken": aToken, "rToken": rToken}, "Token refresh")
}

func (o *UserAuth) getIPandUserAgent(c *gin.Context) (IP, UserAgent string) {
	val, ok := c.Request.Header["User-Agent"]

	if !ok {
		logger.Debug("not have User-Agent Header", c.Request.Header)
	} else {
		// Todo: Здесь можно использовать готовые библиотеки для парсинга UserAgent
		UserAgent = val[0]
	}

	IP = c.ClientIP()

	return
}

// validateUserData check valid CreateUserInput.
func (o *UserAuth) validateUserData(c *gin.Context, userData CreateUserInput) bool {
	if err := validation.IsEmailValid(userData.Email); err != nil {
		response.GinError(c, http.StatusBadRequest, response.CodeInvalidParams, err.Error(), err)
		return false
	}

	if err := validation.IsUserNameValid(userData.Username); err != nil {
		response.GinError(c, http.StatusBadRequest, response.CodeInvalidParams, err.Error(), err)
		return false
	}

	if err := validation.IsPasswordValid(userData.Password); err != nil {
		response.GinError(c, http.StatusBadRequest, response.CodeInvalidParams, err.Error(), err)
		return false
	}

	return true
}
