package party

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
	"ricardo/party-service/internal/core/app/party"
	"ricardo/party-service/internal/core/entities"
	errors2 "ricardo/party-service/pkg/errors"
)

type Controller interface {
	Create(gtx *gin.Context)
	Update(gtx *gin.Context)
	Get(gtx *gin.Context)
	GetForUser(gtx *gin.Context)
	Delete(gtx *gin.Context)
}

type controller struct {
	service party.Service
}

func NewController(service party.Service) Controller {
	return controller{service: service}
}

func (c controller) Create(gtx *gin.Context) {
	var cpr entities.CreatePartyRequest
	_ = gtx.ShouldBindJSON(&cpr)

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
	_ = gtx.ShouldBindJSON(&upr)

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
	var gpr entities.GetPartyRequest
	_ = gtx.ShouldBindJSON(&gpr)

	parties, err := c.service.GetAll(gtx.Request.Context())
	if err != nil {
		_ = errors2.GinErrorHandler(gtx, err, http.StatusInternalServerError)
		return
	}

	gtx.JSON(http.StatusOK, parties)
}

func (c controller) GetForUser(gtx *gin.Context) {
	var gpr entities.GetPartyRequest
	_ = gtx.ShouldBindJSON(&gpr)

	userID := gtx.GetUint("user_id")
	parties, err := c.service.GetAllForUser(gtx.Request.Context(), userID)
	if err != nil {
		_ = errors2.GinErrorHandler(gtx, err, http.StatusInternalServerError)
		return
	}

	gtx.JSON(http.StatusOK, parties)
}

func (c controller) Delete(gtx *gin.Context) {
	var dpr entities.DeletePartyRequest
	_ = gtx.ShouldBindJSON(&dpr)

	err := c.service.Delete(gtx.Request.Context(), dpr.ID)
	if err != nil {
		_ = errors2.GinErrorHandler(gtx, err, http.StatusInternalServerError)
		return
	}

	gtx.Status(http.StatusOK)
}
