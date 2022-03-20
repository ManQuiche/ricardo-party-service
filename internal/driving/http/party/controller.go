package party

import (
	"errors"
	"github.com/gin-gonic/gin"
	tokens "gitlab.com/ricardo-public/jwt-tools/pkg"
	"gorm.io/gorm"
	"net/http"
	"ricardo/party-service/internal/core/app/party"
	"ricardo/party-service/internal/core/entities"
	errors2 "ricardo/party-service/pkg/errors"
	"strconv"
)

type Controller interface {
	Create(gtx *gin.Context)
	Update(gtx *gin.Context)
	Get(gtx *gin.Context)
	GetForUser(gtx *gin.Context)
	GetOne(gtx *gin.Context)
	Delete(gtx *gin.Context)
}

type controller struct {
	service      party.Service
	accessSecret []byte
}

func NewController(service party.Service, accessSecret []byte) Controller {
	return controller{service: service, accessSecret: accessSecret}
}

func (c controller) Create(gtx *gin.Context) {
	var cpr entities.CreatePartyRequest
	err := gtx.ShouldBindJSON(&cpr)
	if err != nil {
		_ = errors2.GinErrorHandler(gtx, err, http.StatusBadRequest)
		return
	}

	p := entities.Party{
		Name:   cpr.Name,
		UserID: cpr.UserID,
	}
	cParty, err := c.service.Save(gtx.Request.Context(), p)
	if err != nil {
		_ = errors2.GinErrorHandler(gtx, err, http.StatusInternalServerError)
		return
	}

	gtx.JSON(http.StatusOK, cParty)
}

func (c controller) Update(gtx *gin.Context) {
	var upr entities.UpdatePartyRequest
	err := gtx.ShouldBindJSON(&upr)
	if err != nil {
		_ = errors2.GinErrorHandler(gtx, err, http.StatusBadRequest)
		return
	}

	_, err = c.canUpdateOrDelete(gtx, upr.ID)
	if err != nil {
		return
	}

	p := entities.Party{
		Model: gorm.Model{
			ID: upr.ID,
		},
		Name:   upr.Name,
		UserID: upr.UserID,
	}

	uParty, err := c.service.Save(gtx.Request.Context(), p)
	if err != nil {
		_ = errors2.GinErrorHandler(gtx, err, http.StatusInternalServerError)
		return
	}

	gtx.JSON(http.StatusOK, uParty)
}

func (c controller) Get(gtx *gin.Context) {
	parties, err := c.service.GetAll(gtx.Request.Context())
	if err != nil {
		_ = errors2.GinErrorHandler(gtx, err, http.StatusInternalServerError)
		return
	}

	gtx.JSON(http.StatusOK, parties)
}

func (c controller) GetForUser(gtx *gin.Context) {
	userId, err := strconv.ParseUint(gtx.Param("user_id"), 10, 64)
	if err != nil {
		_ = errors2.GinErrorHandler(gtx, err, http.StatusBadRequest)
	}

	parties, err := c.service.GetAllForUser(gtx.Request.Context(), uint(userId))
	if err != nil {
		_ = errors2.GinErrorHandler(gtx, err, http.StatusInternalServerError)
		return
	}

	gtx.JSON(http.StatusOK, parties)
}

func (c controller) GetOne(gtx *gin.Context) {
	partyId, err := strconv.ParseUint(gtx.Param("party_id"), 10, 64)
	if err != nil {
		_ = errors2.GinErrorHandler(gtx, err, http.StatusBadRequest)
	}

	parties, err := c.service.Get(gtx.Request.Context(), uint(partyId))
	if err != nil {
		_ = errors2.GinErrorHandler(gtx, err, http.StatusInternalServerError)
		return
	}

	gtx.JSON(http.StatusOK, parties)
}

func (c controller) Delete(gtx *gin.Context) {
	var dpr entities.DeletePartyRequest
	err := gtx.ShouldBindJSON(&dpr)
	if err != nil {
		_ = errors2.GinErrorHandler(gtx, err, http.StatusBadRequest)
		return
	}

	_, err = c.canUpdateOrDelete(gtx, dpr.ID)
	if err != nil {
		return
	}

	err = c.service.Delete(gtx.Request.Context(), dpr.ID)
	if err != nil {
		_ = errors2.GinErrorHandler(gtx, err, http.StatusInternalServerError)
		return
	}

	gtx.Status(http.StatusOK)
}

func (c controller) canUpdateOrDelete(gtx *gin.Context, partyID uint) (bool, error) {
	p, err := c.service.Get(gtx.Request.Context(), partyID)
	if err != nil {
		_ = errors2.GinErrorHandler(gtx, err, http.StatusNotFound)
		return false, err
	}

	strToken, err := tokens.ExtractTokenFromHeader(gtx.GetHeader(tokens.AuthorizationHeader))
	if err != nil {
		_ = errors2.GinErrorHandler(gtx, err, http.StatusInternalServerError)
		return false, err
	}

	pToken, err := tokens.Parse(strToken, c.accessSecret)
	claims := pToken.Claims.(tokens.RicardoClaims)
	userID, err := strconv.ParseUint(claims.Subject, 10, 64)

	if uint(userID) != p.UserID && claims.Role != tokens.AdminRole {
		err = errors.New("unauthorized to update or delete")
		_ = errors2.GinErrorHandler(gtx, err, http.StatusUnauthorized)
		return false, err
	}

	return true, nil
}
