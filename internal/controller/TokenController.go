package controller

import (
	"test_go_project/internal/customErrors"
	"test_go_project/internal/service"
	"test_go_project/internal/dto"

	"net/http"
	"errors"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type TokenController struct {
    service *service.TokenService
}

func NewTokenController(s *service.TokenService) *TokenController {
    return &TokenController{service: s}
}

// @Summary Get a new token pair for given GUID
// @Description Generates JWT access and base64 refresh tokens for given GUID
// @Tags tokens
// @Produce json
// @Param guid path string true "User ID which tokens are generated for" example:"a739620f-89de-4a12-976c-37a9661c1b91"
// @Success 200 {object} dto.TokenPair
// @Failure 400 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /tokens/{guid} [get]
func (c *TokenController) GetNewTokenPair(context *gin.Context) {
	userID, err := uuid.Parse(context.Param("guid"))

	if err != nil {
		context.AbortWithStatusJSON(http.StatusBadRequest, dto.ErrorResponse{Error: customErrors.ErrUUIDParsing.Error()})
		return
	}

	tokenPair, err := c.service.GetNewTokenPair(userID, 
		context.GetHeader("User-Agent"),
		context.ClientIP())

	if err != nil {
		if customErrors.IsCustomError(err) {
			context.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		} else {
			context.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		}
		return
	}

	context.JSON(http.StatusOK, tokenPair)
}

// @Summary Update a token pair
// @Description Updates both tokens after performing all required checks
// @Tags tokens
// @Produce json
// @Param refresh_token path string true "base64 refresh token" example:"fFqZUrwRFSLSLpw0eYXHRGYKlkDKs1TboXmlOy6S0pg"
// @Success 200 {object} dto.TokenPair
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /tokens/refresh/{refresh_token} [get]
// @Security BearerAuth
func (c *TokenController) UpdateTokenPair(context *gin.Context) {
    userID, err := getUserIDFromAccessToken(context)

	if err != nil {
		return
	}

	accessJTI := context.GetString("accessJTI")

	tokenPair, err := c.service.UpdateTokenPair(*userID, 
		context.GetHeader("User-Agent"), 
		context.ClientIP(),
		context.Param("refresh_token"),
		accessJTI)

	if err != nil {
		switch {
		case errors.Is(err, customErrors.ErrRecordNotFound):
			context.AbortWithStatusJSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})

		case errors.Is(err, customErrors.ErrUserAgentMismatch),
			errors.Is(err, customErrors.ErrUserIDMismatch),
			errors.Is(err, customErrors.ErrRefreshMismatch),
			errors.Is(err, customErrors.ErrTokenPairing):
			context.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})

		case errors.Is(err, customErrors.ErrDB),
			errors.Is(err, customErrors.ErrUnexpectedState):
			context.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		
		default:
			context.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		}
		return
	}

	context.JSON(http.StatusOK, tokenPair)
}

// @Summary Get GUID
// @Description Returns GUID extracted from token and checks if record with such GUID exists in DB
// @Tags user
// @Produce json
// @Success 200 {object} dto.UserIdDto
// @Failure 404 {object} dto.ErrorResponse
// @Failure 500 {object} dto.ErrorResponse
// @Router /user [get]
// @Security BearerAuth
func (c *TokenController) GetGUID(context *gin.Context) {
    userID, err := getUserIDFromAccessToken(context)

	if err != nil {
		return
	}

	err = c.service.ExistsById(*userID)
	
	if err != nil {
		switch {
		case errors.Is(err, customErrors.ErrRecordNotFound):
			context.AbortWithStatusJSON(http.StatusNotFound, dto.ErrorResponse{Error: err.Error()})

		case errors.Is(err, customErrors.ErrDB):
			context.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		
		default:
			context.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		}
		return
	}

	context.JSON(http.StatusOK,  dto.UserIdDto{
		UserID: userID.String(),
	})
}

// @Summary Clear tokens
// @Description Removes the record about access/refresh token pair
// @Tags tokens
// @Success 200
// @Failure 500 {object} dto.ErrorResponse
// @Router /tokens [delete]
// @Security BearerAuth
func (c *TokenController) ClearTokens(context *gin.Context) {
    userID, err := getUserIDFromAccessToken(context)

	if err != nil {
		return
	}

	err = c.service.DeleteById(*userID)
	
	if err != nil {
		if customErrors.IsCustomError(err) {
			context.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		} else {
			context.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		}
		return
	}

	context.Status(http.StatusOK)
}

func getUserIDFromAccessToken(context *gin.Context) (*uuid.UUID, error) {
	rawUserID, exists := context.Get("userID")
	if !exists {
		err := customErrors.ErrMissingUserID
		context.AbortWithStatusJSON(http.StatusUnauthorized, dto.ErrorResponse{Error: err.Error()})
		return nil, err
	}

	userId, ok := rawUserID.(uuid.UUID)
	if !ok {
		err := customErrors.ErrUUIDParsing
		context.AbortWithStatusJSON(http.StatusInternalServerError, dto.ErrorResponse{Error: err.Error()})
		return nil, err
	}

	return &userId, nil
}