package party

import (
	"errors"
	"github.com/gin-gonic/gin"
	ricardoerr "gitlab.com/ricardo-public/errors/pkg/errors"
	tokens "gitlab.com/ricardo-public/jwt-tools/pkg"
	"gorm.io/gorm"
	"net/http"
	"ricardo/party-service/internal/core/app/party"
	"ricardo/party-service/internal/core/entities"
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
		_ = ricardoerr.GinErrorHandler(gtx, ricardoerr.New(ricardoerr.ErrBadRequest, ""))
		return
	}

	p := entities.Party{
		Name:   cpr.Name,
		UserID: cpr.UserID,
	}
	cParty, err := c.service.Save(gtx.Request.Context(), p)
	if err != nil {
		_ = ricardoerr.GinErrorHandler(gtx, err)
		return
	}

	gtx.JSON(http.StatusOK, cParty)
}

func (c controller) Update(gtx *gin.Context) {
	var upr entities.UpdatePartyRequest
	err := gtx.ShouldBindJSON(&upr)
	if err != nil {
		_ = ricardoerr.GinErrorHandler(gtx, ricardoerr.New(ricardoerr.ErrBadRequest, ""))
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
		_ = ricardoerr.GinErrorHandler(gtx, err)
		return
	}

	gtx.JSON(http.StatusOK, uParty)
}

func (c controller) Get(gtx *gin.Context) {
	parties, err := c.service.GetAll(gtx.Request.Context())
	if err != nil {
		_ = ricardoerr.GinErrorHandler(gtx, ricardoerr.New(ricardoerr.ErrBadRequest, ""))
		return
	}

	gtx.JSON(http.StatusOK, parties)
}

func (c controller) GetForUser(gtx *gin.Context) {
	userId, err := strconv.ParseUint(gtx.Param("user_id"), 10, 64)
	if err != nil {
		_ = ricardoerr.GinErrorHandler(gtx, ricardoerr.New(ricardoerr.ErrBadRequest, err.Error()))
	}

	parties, err := c.service.GetAllForUser(gtx.Request.Context(), uint(userId))
	if err != nil {
		_ = ricardoerr.GinErrorHandler(gtx, err)
		return
	}

	gtx.JSON(http.StatusOK, parties)
}

func (c controller) GetOne(gtx *gin.Context) {
	partyId, err := strconv.ParseUint(gtx.Param("party_id"), 10, 64)
	if err != nil {
		_ = ricardoerr.GinErrorHandler(gtx, ricardoerr.New(ricardoerr.ErrBadRequest, ""))
		return
	}

	parties, err := c.service.Get(gtx.Request.Context(), uint(partyId))
	if err != nil {
		_ = ricardoerr.GinErrorHandler(gtx, err)
		return
	}

	gtx.JSON(http.StatusOK, parties)
}

func (c controller) Delete(gtx *gin.Context) {
	var dpr entities.DeletePartyRequest
	err := gtx.ShouldBindJSON(&dpr)
	if err != nil {
		_ = ricardoerr.GinErrorHandler(gtx, ricardoerr.New(ricardoerr.ErrBadRequest, ""))
		return
	}

	_, err = c.canUpdateOrDelete(gtx, dpr.ID)
	if err != nil {
		return
	}

	err = c.service.Delete(gtx.Request.Context(), dpr.ID)
	if err != nil {
		_ = ricardoerr.GinErrorHandler(gtx, err)
		return
	}

	gtx.Status(http.StatusOK)
}

func (c controller) canUpdateOrDelete(gtx *gin.Context, partyID uint) (bool, error) {
	p, err := c.service.Get(gtx.Request.Context(), partyID)
	if err != nil {
		_ = ricardoerr.GinErrorHandler(gtx, ricardoerr.New(ricardoerr.ErrBadRequest, ""))
		return false, err
	}

	strToken, err := tokens.ExtractTokenFromHeader(gtx.GetHeader(tokens.AuthorizationHeader))
	if err != nil {
		_ = ricardoerr.GinErrorHandler(gtx, err)
		return false, err
	}

	pToken, err := tokens.Parse(strToken, c.accessSecret)
	claims := pToken.Claims.(tokens.RicardoClaims)
	userID, err := strconv.ParseUint(claims.Subject, 10, 64)

	if uint(userID) != p.UserID && claims.Role != tokens.AdminRole {
		err = errors.New("unauthorized to update or delete")
		_ = ricardoerr.GinErrorHandler(gtx, err)
		return false, err
	}

	return true, nil
}
