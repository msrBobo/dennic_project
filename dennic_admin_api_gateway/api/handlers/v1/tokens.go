package v1

import (
	e "dennic_admin_api_gateway/api/handlers/regtool"
	"dennic_admin_api_gateway/api/models/model_user_service"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetTokens
// @Summary GetTokens
// @Description GetTokens
// @Tags Token
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} model_user_service.Tokens "Successful response"
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/token/get-token [GET]
func (h *HandlerV1) GetTokens(c *gin.Context) {

	access, refresh, err := h.jwthandler.GenerateAuthJWT("phone_number1", "123e4567-e89b-12d3-a456-426614174001", "1", "admin")
	if e.HandleError(c, err, h.log, http.StatusInternalServerError, "GetTokens") {
		return
	}
	c.JSON(http.StatusOK, model_user_service.Tokens{
		AccessToken:  access,
		RefreshToken: refresh,
	})
}
