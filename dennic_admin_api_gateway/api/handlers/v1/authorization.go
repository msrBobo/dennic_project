package v1

import (
	"context"
	e "dennic_admin_api_gateway/api/handlers/regtool"
	"dennic_admin_api_gateway/api/models"
	"dennic_admin_api_gateway/api/models/model_user_service"
	a "dennic_admin_api_gateway/api/models/models_admin_service"
	pb "dennic_admin_api_gateway/genproto/user_service"
	email "dennic_admin_api_gateway/internal/pkg/email"
	r "dennic_admin_api_gateway/internal/pkg/rand"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"google.golang.org/protobuf/encoding/protojson"
)

// Login ...
// @Summary Login
// @Description Login - Api for registering users
// @Tags Authorization
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param Login body model_user_service.LoginReq true "Login Req"
// @Success 200 {object} model_user_service.Response
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/admin/login [post]
func (h *HandlerV1) Login(c *gin.Context) {
	var (
		body        model_user_service.LoginReq
		jsonMarshal protojson.MarshalOptions
	)

	jsonMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)

	if e.HandleError(c, err, h.log, http.StatusBadRequest, "Login") {
		return
	}
	Email := body.Email
	if !a.IsValidEmail(Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email address"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	if !e.ValidatePassword(body.Password) {
		err := errors.New("invalid password")
		_ = e.HandleError(c, err, h.log, http.StatusBadRequest, err.Error())
		return
	}
	var (
		admin *pb.Admin
	)

	admin, err = h.serviceManager.UserService().AdminService().Get(ctx, &pb.GetAdminReq{
		Field:    "email",
		Value:    body.Email,
		IsActive: false,
	})
	if e.HandleError(c, err, h.log, http.StatusBadRequest, NOT_REGISTERED) {
		return
	}
	if !e.CheckHashPassword(admin.Password, body.Password) {
		err = errors.New("incorrect password")
		_ = e.HandleError(c, err, h.log, http.StatusBadRequest, err.Error())
		return
	}

	access, refresh, err := h.jwthandler.GenerateAuthJWT(admin.PhoneNumber, admin.Id, "admin", "admin")
	if e.HandleError(c, err, h.log, http.StatusInternalServerError, SERVICE_ERROR) {
		return
	}

	_, err = h.serviceManager.UserService().AdminService().UpdateRefreshToken(ctx, &pb.UpdateRefreshTokenAdminReq{
		Id:           admin.Id,
		RefreshToken: refresh,
	})

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, SERVICE_ERROR) {
		return
	}

	c.JSON(http.StatusOK, a.CreateAdminResp{
		Id:            admin.Id,
		AdminOrder:    admin.AdminOrder,
		Role:          admin.Role,
		FirstName:     admin.FirstName,
		LastName:      admin.LastName,
		BrithDate:     admin.BirthDate,
		PhoneNumber:   admin.PhoneNumber,
		Email:         admin.Email,
		Password:      admin.Password,
		Gender:        admin.Gender,
		Salary:        admin.Salary,
		Biography:     admin.Biography,
		StartWorkYear: admin.StartWorkYear,
		EndWorkYear:   admin.EndWorkYear,
		WorkYears:     admin.WorkYears,
		ImageUrl:      admin.ImageUrl,
		AccessToken:   access,
		ReflashToken:  refresh,
		CreatedAt:     admin.CreatedAt,
	})
}

// ForgetPassword ...
// @Summary ForgetPassword
// @Description ForgetPassword - Api for registering users
// @Tags Authorization
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param ForgetPassword body model_user_service.EmailRerReq true "RegisterModelReq"
// @Success 200 {object} model_user_service.MessageRes
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/admin/forget-password [post]
func (h *HandlerV1) ForgetPassword(c *gin.Context) {
	var (
		body        model_user_service.EmailRerReq
		jsonMarshal protojson.MarshalOptions
	)

	jsonMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)

	if e.HandleError(c, err, h.log, http.StatusBadRequest, INVALID_REQUET_BODY) {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	Email := body.Email
	if !a.IsValidEmail(Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email address"})
		return
	}

	existsPhone, err := h.serviceManager.UserService().AdminService().CheckField(ctx, &pb.CheckAdminFieldReq{
		Field: "email",
		Value: body.Email,
	})

	if e.HandleError(c, err, h.log, http.StatusBadRequest, NOT_REGISTERED) {
		return
	}
	if !existsPhone.Status {
		err = errors.New(NOT_REGISTERED)
		_ = e.HandleError(c, err, h.log, http.StatusBadRequest, NOT_REGISTERED)
		return
	}

	codeRed, _ := h.redis.Client.Get(ctx, body.Email).Result()

	if codeRed != "" {
		err = errors.New(CODE_EXPIRATION_NOT_OVER)
		if e.HandleError(c, err, h.log, http.StatusBadRequest, CODE_EXPIRATION_NOT_OVER) {
			return
		}
	}

	// TODO A method that Cosessends a code to a number
	code := r.GenerateRandomNumbers()
	fmt.Println(code)
	data := fmt.Sprintf("Your verification code: %s", code)

	emailData := email.EmailData{
		Message: data,
	}
	// send otp email
	go func() {
		err := email.SendEmail([]string{body.Email}, "Dennic.UZ\n", *h.cfg, "internal/pkg/email/emailotp.html", emailData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "While Send Email"})
			return
		}
	}()
	err = h.redis.Client.Set(ctx, body.Email, code, h.cfg.Redis.Time).Err()
	if e.HandleError(c, err, h.log, http.StatusInternalServerError, SERVICE_ERROR) {
		return
	}

	c.JSON(http.StatusOK, model_user_service.MessageRes{
		Message: "Code has been sent to your email, please check.",
	})
}

// UpdatePassword
// @Summary UpdatePassword
// @Description Api for UpdatePassword
// @Tags Authorization
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param NewPassword  query string true "NewPassword"
// @Success 200 {object} model_user_service.GetUserResp
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/admin/update-password [PUT]
func (h *HandlerV1) UpdatePassword(c *gin.Context) {
	newPassword := c.Query("NewPassword")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	userInfo, err := e.GetUserInfo(c)

	if e.HandleError(c, err, h.log, http.StatusUnauthorized, "missing token in the header") {
		return
	}

	admin, err := h.serviceManager.UserService().AdminService().Get(ctx, &pb.GetAdminReq{
		Field:    "id",
		Value:    userInfo.UserId,
		IsActive: false,
	})

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, SERVICE_ERROR) {
		return
	}

	hashPass, err := e.HashPassword(newPassword)

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, SERVICE_ERROR) {
		return
	}

	response, err := h.serviceManager.UserService().AdminService().ChangePassword(ctx, &pb.ChangeAdminPasswordReq{
		PhoneNumber: admin.PhoneNumber,
		Password:    hashPass,
	})

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, SERVICE_ERROR) {
		return
	}

	c.JSON(http.StatusOK, &models.StatusRes{Status: response.Status})
}

// VerifyOtpCode ...
// @Summary VerifyOtpCode
// @Description VerifyOtpCode - Api for Verify Otp Code users
// @Tags Authorization
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param VerifyOtpCode query model_user_service.VerifyOtpCodeReq true "VerifyOtpCode"
// @Failure 200 {object} model_user_service.MessageRes
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/admin/verify-otp-code [post]
func (h *HandlerV1) VerifyOtpCode(c *gin.Context) {
	email := c.Query("email")
	code := c.Query("code")

	reqCode := cast.ToInt64(code)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	if !a.IsValidEmail(email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email address"})
		return
	}

	redisRes, err := h.redis.Client.Get(ctx, email).Result()

	if e.HandleError(c, err, h.log, http.StatusBadRequest, "code is expired") {
		return
	}

	var redisCode int64

	err = json.Unmarshal([]byte(redisRes), &redisCode)

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, SERVICE_ERROR) {
		return
	}

	if reqCode != redisCode {
		err = errors.New(INVALID_CODE)
		_ = e.HandleError(c, err, h.log, http.StatusBadRequest, INVALID_CODE)
		return
	}

	err = h.redis.Client.Del(ctx, email).Err()

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, SERVICE_ERROR) {
		return
	}

	admin, err := h.serviceManager.UserService().AdminService().Get(ctx, &pb.GetAdminReq{
		Field:    "email",
		Value:    email,
		IsActive: false,
	})

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, SERVICE_ERROR) {
		return
	}

	access, err := h.jwthandler.GenerateJWT(admin.PhoneNumber, admin.Id, "admin")
	if e.HandleError(c, err, h.log, http.StatusInternalServerError, SERVICE_ERROR) {
		return
	}

	c.JSON(http.StatusOK, &models.AccessToken{Token: access})

}

// SenOtpCode ...
// @Summary SenOtpCode
// @Description SenOtpCode - Api for sen otp code users
// @Tags Authorization
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param SenOtpCode body model_user_service.EmailRerReq true "RegisterModelReq"
// @Success 200 {object} model_user_service.MessageRes
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/admin/send-otp [post]
func (h *HandlerV1) SenOtpCode(c *gin.Context) {
	var (
		body        model_user_service.EmailRerReq
		jsonMarshal protojson.MarshalOptions
	)

	jsonMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)

	if e.HandleError(c, err, h.log, http.StatusBadRequest, INVALID_REQUET_BODY) {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	existsPhone, err := h.serviceManager.UserService().AdminService().CheckField(ctx, &pb.CheckAdminFieldReq{
		Field: "email",
		Value: body.Email,
	})
	if e.HandleError(c, err, h.log, http.StatusBadRequest, INVALID_REQUET_BODY) {
		return
	}

	codeRed, _ := h.redis.Client.Get(ctx, body.Email).Result()
	if codeRed != "" {
		err = errors.New(CODE_EXPIRATION_NOT_OVER)
		if e.HandleError(c, err, h.log, http.StatusBadRequest, CODE_EXPIRATION_NOT_OVER) {
			return
		}
	}

	if !existsPhone.Status {
		err = errors.New(NOT_REGISTERED)
		_ = e.HandleError(c, err, h.log, http.StatusBadRequest, NOT_REGISTERED)
		return
	}

	// TODO A method that Cosessends a code to a number
	code := r.GenerateRandomNumbers()
	fmt.Println(code)
	data := fmt.Sprintf("Your verification code: %s", code)

	emailData := email.EmailData{
		Message: data,
	}
	// send otp email
	go func() {
		err := email.SendEmail([]string{body.Email}, "Dennic.UZ\n", *h.cfg, "internal/pkg/email/emailotp.html", emailData)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "While Send Email"})
			return
		}
	}()

	err = h.redis.Client.Set(ctx, body.Email, code, h.cfg.Redis.Time).Err()
	if e.HandleError(c, err, h.log, http.StatusInternalServerError, SERVICE_ERROR) {
		return
	}

	c.JSON(http.StatusOK, model_user_service.MessageRes{
		Message: "Code has been sent to your email, please check.",
	})
}
