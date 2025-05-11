package handlers

import (
	"net/http"
	"strings"
	"surfe/internal/services"

	"github.com/labstack/echo/v4"
)

type ActionHandler struct {
	actionService services.ActionService
}

func NewActionHandler(actionService services.ActionService) *ActionHandler {
	return &ActionHandler{
		actionService: actionService,
	}
}

func (h *ActionHandler) GetNextActionProbabilities(c echo.Context) error {
	actionType := strings.ToUpper(c.Param("type"))

	probabilities, err := h.actionService.GetNextActionProbabilities(actionType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, probabilities)
}

func (h *ActionHandler) GetReferralIndex(c echo.Context) error {
	referralIndex, err := h.actionService.GetReferralIndex()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, referralIndex)
}
