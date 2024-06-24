package v1

import (
	"context"
	e "dennic_admin_api_gateway/api/handlers/regtool"
	"dennic_admin_api_gateway/api/models/model_user_service"
	a "dennic_admin_api_gateway/api/models/models_admin_service"
	pb "dennic_admin_api_gateway/genproto/user_service"
	email "dennic_admin_api_gateway/internal/pkg/email"
	jwt "dennic_admin_api_gateway/internal/pkg/tokens"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cast"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// CreateAdmin
// @Summary CreateAdmin
// @Description Api for CreateAdmin
// @Security ApiKeyAuth
// @Tags Admin
// @Accept json
// @Produce json
// @Param CreateAdmin body model_admin_service.CreateAdminReq true "CreateAdminReq"
// @Success 200 {object} model_admin_service.GetAdminResp
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/admin/create [POST]
func (h *HandlerV1) CreateAdmin(c *gin.Context) {

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	var (
		body        a.CreateAdminReq
		jsonMarshal protojson.MarshalOptions
	)
	jsonMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)

	if e.HandleError(c, err, h.log, http.StatusBadRequest, INVALID_REQUET_BODY) {
		return
	}

	err = body.Validate()

	if err != nil {
		e.HandleError(c, err, h.log, http.StatusBadRequest, err.Error())
		return
	}

	body.PhoneNumber = strings.TrimSpace(body.PhoneNumber)
	body.LastName = strings.TrimSpace(body.LastName)
	body.LastName = strings.ToLower(body.LastName)
	body.LastName = strings.Title(body.LastName)
	body.FirstName = strings.TrimSpace(body.FirstName)
	body.FirstName = strings.ToLower(body.FirstName)
	body.FirstName = strings.Title(body.FirstName)

	existsPhone, err := h.serviceManager.UserService().AdminService().CheckField(ctx, &pb.CheckAdminFieldReq{
		Field: "phone_number",
		Value: body.PhoneNumber,
	})

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, NOT_REGISTERED) {
		return
	}
	existsEmail, err := h.serviceManager.UserService().AdminService().CheckField(ctx, &pb.CheckAdminFieldReq{
		Field: "email",
		Value: body.Email,
	})
	if e.HandleError(c, err, h.log, http.StatusInternalServerError, NOT_REGISTERED) {
		return
	}

	if existsPhone.Status || existsEmail.Status {
		err = errors.New("you have already registered, try to login")
		_ = e.HandleError(c, err, h.log, http.StatusBadRequest, "Register")
		return
	}
	Email := body.Email
	if !a.IsValidEmail(Email) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid email address"})
		return
	}
	userInfo, err := e.GetUserInfo(c)
	if e.HandleError(c, err, h.log, http.StatusBadRequest, "cannot get admin by id") {
		return
	}
	userId := userInfo.UserId

	getAdmin, err := h.serviceManager.UserService().AdminService().Get(
		ctx, &pb.GetAdminReq{
			Field:    "id",
			Value:    userId,
			IsActive: false,
		})

	if e.HandleError(c, err, h.log, http.StatusBadRequest, "cannot get admin by id") {
		return
	}
	data := fmt.Sprintf("you have been appointed as an admin to Dennic.UZ by %s %s", getAdmin.FirstName, getAdmin.LastName)

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

	body.Id = uuid.New().String()
	access, refresh, err := h.jwthandler.GenerateAuthJWT(body.PhoneNumber, body.Id, "admin", "admin")

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, SERVICE_ERROR) {
		return
	}
	body.Password, err = e.HashPassword(body.Password)

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, SERVICE_ERROR) {
		return
	}
	response, err := h.serviceManager.UserService().AdminService().Create(
		ctx, &pb.Admin{
			Id:            body.Id,
			Role:          body.Role,
			FirstName:     body.FirstName,
			LastName:      body.LastName,
			BirthDate:     body.BrithDate,
			PhoneNumber:   body.PhoneNumber,
			Email:         body.Email,
			Password:      body.Password,
			Gender:        body.Gender,
			Salary:        body.Salary,
			Biography:     body.Biography,
			StartWorkYear: body.StartWorkYear,
			EndWorkYear:   body.EndWorkYear,
			WorkYears:     body.WorkYears,
			RefreshToken:  refresh,
		})

	if e.HandleError(c, err, h.log, http.StatusBadRequest, "cannot create user") {
		return
	}

	resp := a.CreateAdminResp{
		Id:            response.Id,
		AdminOrder:    response.AdminOrder,
		Role:          response.Role,
		FirstName:     response.FirstName,
		LastName:      response.LastName,
		BrithDate:     response.BirthDate,
		PhoneNumber:   response.PhoneNumber,
		Email:         response.Email,
		Password:      response.Password,
		Gender:        response.Gender,
		Salary:        response.Salary,
		Biography:     response.Biography,
		StartWorkYear: response.StartWorkYear,
		EndWorkYear:   response.EndWorkYear,
		WorkYears:     response.WorkYears,
		ImageUrl:      response.ImageUrl,
		AccessToken:   access,
		ReflashToken:  refresh,
		CreatedAt:     response.CreatedAt,
	}

	c.JSON(http.StatusOK, resp)
}

// GetAdminByID
// @Summary GetAdminByID
// @Description Api for GetAdminByID
// @Security ApiKeyAuth
// @Tags Admin
// @Accept json
// @Produce json
// @Param id query string true "id"
// @Param is_active query bool false "is_active"
// @Success 200 {object} model_admin_service.GetAdminResp
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/admin/get [GET]
func (h *HandlerV1) GetAdminByID(c *gin.Context) {
	id := c.Query("id")
	isActive := c.Query("is_active")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	response, err := h.serviceManager.UserService().AdminService().Get(
		ctx, &pb.GetAdminReq{
			Field:    "id",
			Value:    id,
			IsActive: !cast.ToBool(isActive),
		})

	if e.HandleError(c, err, h.log, http.StatusBadRequest, "cannot get user by id") {
		return
	}

	resp := a.GetAdminResp{
		Id:            response.Id,
		AdminOrder:    uint64(response.AdminOrder),
		FirstName:     response.FirstName,
		LastName:      response.LastName,
		Role:          response.Role,
		BrithDate:     response.BirthDate,
		PhoneNumber:   response.PhoneNumber,
		Email:         response.Email,
		Password:      response.Password,
		Gender:        response.Gender,
		Salary:        response.Salary,
		Biography:     response.Biography,
		StartWorkYear: response.StartWorkYear,
		EndWorkYear:   response.EndWorkYear,
		WorkYears:     response.WorkYears,
		RefreshToken:  response.RefreshToken,
		ImageUrl:      response.ImageUrl,
		CreatedAt:     response.CreatedAt,
		UpdatedAt:     response.UpdatedAt,
		DeletedAt:     response.DeletedAt,
	}

	c.JSON(http.StatusOK, resp)

}

// ListAdmins
// @Summary ListAdmins
// @Description Api for ListAdmins
// @Security ApiKeyAuth
// @Tags Admin
// @Accept json
// @Produce json
// @Param Page  query string false "Page"
// @Param Limit query string false "Limit"
// @Param searchField query string false "searchField" Enums(first_name, last_name, gender, phone_number, email, biography, created_at)
// @Param Value query string false "Value"
// @Param OrderBy query string false "OrderBy"
// @Param is_active query bool false "is_active"
// @Success 200 {object} model_admin_service.GetAdminResp
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/admin [GET]
func (h *HandlerV1) ListAdmins(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	page := c.Query("Page")
	limit := c.Query("Limit")
	field := c.Query("searchField")
	value := c.Query("Value")
	orderBy := c.Query("OrderBy")
	is_active := c.Query("is_active")

	pageInt, limitInt, err := e.ParseQueryParams(page, limit)
	if e.HandleError(c, err, h.log, http.StatusBadRequest, "failed to list users") {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	response, err := h.serviceManager.UserService().AdminService().ListAdmins(
		ctx, &pb.ListAdminsReq{
			Page:     pageInt,
			Limit:    limitInt,
			IsActive: !cast.ToBool(is_active),
			Value:    value,
			Field:    field,
			OrderBy:  orderBy,
		})
	if e.HandleError(c, err, h.log, http.StatusInternalServerError, SERVICE_ERROR) {
		return
	}

	var admins a.ListAdminrResp

	for _, in := range response.Admins {
		admin := a.GetAdminResp{
			Id:            in.Id,
			AdminOrder:    uint64(in.AdminOrder),
			FirstName:     in.FirstName,
			LastName:      in.LastName,
			Role:          in.Role,
			BrithDate:     in.BirthDate,
			PhoneNumber:   in.PhoneNumber,
			Email:         in.Email,
			Password:      in.Password,
			Gender:        in.Gender,
			Salary:        in.Salary,
			Biography:     in.Biography,
			StartWorkYear: in.StartWorkYear,
			EndWorkYear:   in.EndWorkYear,
			WorkYears:     in.WorkYears,
			ImageUrl:      in.ImageUrl,
			CreatedAt:     in.CreatedAt,
			UpdatedAt:     in.UpdatedAt,
			DeletedAt:     in.DeletedAt,
		}

		admins.Admins = append(admins.Admins, admin)
		admins.Count = response.Count
	}

	c.JSON(http.StatusOK, admins)
}

// UpdateAdmin
// @Summary UpdateAdmin
// @Description Api for UpdateAdmin
// @Security ApiKeyAuth
// @Tags Admin
// @Accept json
// @Produce json
// @Param UpdAdminReq body model_admin_service.UpdAdminReq true "UpdAdminReq"
// @Success 200 {object} model_admin_service.UpdAdminResp
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/admin [PUT]
func (h *HandlerV1) UpdateAdmin(c *gin.Context) {
	var (
		body        a.UpdAdminReq
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if e.HandleError(c, err, h.log, http.StatusBadRequest, INVALID_REQUET_BODY) {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	req := &pb.Admin{
		Id:            body.Id,
		FirstName:     body.FirstName,
		LastName:      body.LastName,
		BirthDate:     body.BrithDate,
		Gender:        body.Gender,
		Salary:        body.Salary,
		Biography:     body.Biography,
		StartWorkYear: body.StartWorkYear,
		EndWorkYear:   body.EndWorkYear,
		WorkYears:     body.WorkYears,
		ImageUrl:      body.ImageUrl,
	}

	response, err := h.serviceManager.UserService().AdminService().Update(ctx, req)
	if e.HandleError(c, err, h.log, http.StatusInternalServerError, SERVICE_ERROR) {
		return
	}

	resp := a.UpdAdminResp{
		Id:            response.Id,
		FirstName:     response.FirstName,
		LastName:      response.LastName,
		BrithDate:     response.BirthDate,
		Gender:        response.Gender,
		Salary:        response.Salary,
		Biography:     response.Biography,
		StartWorkYear: response.StartWorkYear,
		EndWorkYear:   response.EndWorkYear,
		WorkYears:     response.WorkYears,
		ImageUrl:      response.ImageUrl,
		UpdatedAt:     response.UpdatedAt,
	}

	c.JSON(http.StatusOK, resp)
}

// UpdateAdminRefreshToken
// @Summary Update Admin Refresh Token
// @Description Update the refresh token of the admin
// @Security ApiKeyAuth
// @Tags Admin
// @Accept json
// @Produce json
// @Param RefreshToken body model_user_service.RefreshToken true "RefreshToken"
// @Success 200 {object} model_user_service.UpdateRefreshTokenUserResp "Successful response"
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/admin/update-admin-refresh-token [PUT]
func (h *HandlerV1) UpdateAdminRefreshToken(c *gin.Context) {

	var (
		RefreshToken model_user_service.RefreshToken
		jspbMarshal  protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&RefreshToken)
	if e.HandleError(c, err, h.log, http.StatusBadRequest, INVALID_REQUET_BODY) {
		return
	}

	claims, err := jwt.ExtractClaim(RefreshToken.RefreshToken)
	if e.HandleError(c, err, h.log, http.StatusUnauthorized, "missing token in the header") {
		return
	}

	access, refresh, err := h.jwthandler.GenerateAuthJWT(cast.ToString(claims["phone"]), cast.ToString(claims["id"]), cast.ToString(claims["session_id"]), "admin")
	if e.HandleError(c, err, h.log, http.StatusInternalServerError, SERVICE_ERROR) {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	_, err = h.serviceManager.UserService().AdminService().UpdateRefreshToken(ctx, &pb.UpdateRefreshTokenAdminReq{
		Id:           cast.ToString(claims["id"]),
		RefreshToken: refresh,
	})

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, SERVICE_ERROR) {
		return
	}
	resp := model_user_service.UpdateRefreshTokenUserResp{
		AccessToken:  access,
		RefreshToken: refresh,
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteAdmin
// @Summary DeleteAdmin
// @Description Api for DeleteAdmin
// @Security ApiKeyAuth
// @Tags Admin
// @Accept json
// @Produce json
// @Param id  query string true "id"
// @Param hard_delete query bool false "hard_delete"
// @Success 200 {object} model_user_service.CheckUserFieldResp
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/admin [DELETE]
func (h *HandlerV1) DeleteAdmin(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Query("id")
	hard_delete := c.Query("hard_delete")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	response, err := h.serviceManager.UserService().AdminService().Delete(
		ctx, &pb.DeleteAdminReq{
			Field:    "id",
			Value:    id,
			IsActive: cast.ToBool(hard_delete),
		})

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, SERVICE_ERROR) {
		return
	}

	c.JSON(http.StatusOK, response)
}
