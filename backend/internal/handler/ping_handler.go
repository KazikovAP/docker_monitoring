package handler

import (
	"fmt"
	"net"
	"net/http"

	"github.com/KazikovAP/docker_monitoring/backend/internal/model"
	"github.com/KazikovAP/docker_monitoring/backend/internal/service"
	"github.com/gin-gonic/gin"
)

type PingHandler struct {
	service service.PingService
}

func NewPingHandler(pingService service.PingService) *PingHandler {
	return &PingHandler{service: pingService}
}

func (h *PingHandler) GetAllPings(ctx *gin.Context) {
	pings, err := h.service.GetAllPings()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, pings)
}

func (h *PingHandler) AddPing(ctx *gin.Context) {
	var ping model.Ping
	if err := ctx.ShouldBindJSON(&ping); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if net.ParseIP(ping.IPAddress) == nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid IP address"})
		return
	}

	if err := h.service.AddPing(&ping); err != nil {
		if err.Error() == fmt.Sprintf("IP address %s already exists", ping.IPAddress) {
			ctx.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		} else {
			ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{"message": "Ping added successfully"})
}
