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

// @Summary Get next action probabilities
// @Description Get probabilities of next actions based on current action type
// @Tags actions
// @Accept json
// @Produce json
// @Param type path string true "Action Type"
// @Success 200 {object} models.ActionProbability
// @Failure 500 {object} error
// @Router /actions/{type}/next [get]
func (h *ActionHandler) GetNextActionProbabilities(c echo.Context) error {
	actionType := strings.ToUpper(c.Param("type"))

	probabilities, err := h.actionService.GetNextActionProbabilities(actionType)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, probabilities)
}

// @Summary Get referral index
// @Description Get the referral index showing how many users each user has referred
// @Tags actions
// @Accept json
// @Produce json
// @Success 200 {object} map[int]int
// @Failure 500 {object} error
// @Router /actions/referral [get]
func (h *ActionHandler) GetReferralIndex(c echo.Context) error {
	referralIndex, err := h.actionService.GetReferralIndex()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "Internal server error"})
	}

	return c.JSON(http.StatusOK, referralIndex)
}
