package v1

import (
	"dennic_api_gateway/api/models"
	"dennic_api_gateway/internal/pkg/email"
	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
)

// Support
// @Summary Support
// @Description Api  Support
// @Tags Support
// @Accept json
// @Produce json
// @Param CreateAdmin body models.SupportReq true "CreateAdminReq"
// @Success 200 {object} models.SupportRes
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/support [POST]
func (h *HandlerV1) Support(c *gin.Context) {
	var (
		body        models.SupportReq
		jsonMarshal protojson.MarshalOptions
	)
	jsonMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	emailData := email.EmailData{
		FullName:    body.FullName,
		PhoneNumber: body.PhoneNumber,
		Email:       body.Email,
		Message:     body.Message,
	}

	go func() {
		err := email.SendEmail([]string{"dennic.uz@gmail.com"}, "Support Request", *h.cfg, "internal/pkg/email/emailotp.html", emailData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "While Send Email"})
			return
		}
	}()

	c.JSON(http.StatusOK, &models.SupportRes{Message: "Your message was sent to our admins!"})
}
