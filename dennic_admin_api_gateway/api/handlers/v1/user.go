package v1

import (
	"context"
	e "dennic_admin_api_gateway/api/handlers/regtool"
	"dennic_admin_api_gateway/api/models/model_user_service"
	pb "dennic_admin_api_gateway/genproto/user_service"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/spf13/cast"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// CreateUser
// @Summary CreateUser
// @Description Api for CreateUser
// @Security ApiKeyAuth
// @Tags User
// @Accept json
// @Produce json
// @Param CreateAdmin body model_user_service.User true "User"
// @Success 200 {object} model_user_service.GetUserResp
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/user/create [POST]
func (h *HandlerV1) CreateUser(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	var (
		body        model_user_service.User
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

	existsPhone, err := h.serviceManager.UserService().UserService().CheckField(ctx, &pb.CheckFieldUserReq{
		Field: "phone_number",
		Value: body.PhoneNumber,
	})

	body.Id = uuid.New().String()

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, NOT_REGISTERED) {
		return
	}
	if existsPhone.Status {
		err = errors.New("you have already registered, try to login")
		_ = e.HandleError(c, err, h.log, http.StatusBadRequest, "Register")
		return
	}

	access, refresh, err := h.jwthandler.GenerateAuthJWT(body.PhoneNumber, body.Id, "admin", "user")

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, SERVICE_ERROR) {
		return
	}

	request := &pb.User{
		Id:           body.Id,
		FirstName:    body.FirstName,
		LastName:     body.LastName,
		BirthDate:    body.BrithDate,
		PhoneNumber:  body.PhoneNumber,
		Password:     body.Password,
		Gender:       body.Gender,
		RefreshToken: refresh,
	}

	if request.Gender == "female" {
		request.ImageUrl = "https://minio.dennic.uz/user/154b1855-67fa-48dc-8dd5-1694d8160b81.JPG"
	} else {
		request.ImageUrl = "https://minio.dennic.uz/user/249af563-90cc-47e5-bbe4-5f1127611220.JPG"
	}

	_, err = h.serviceManager.UserService().UserService().Create(ctx, request)

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, SERVICE_ERROR) {
		return
	}

	c.JSON(http.StatusOK, &model_user_service.Response{
		Id:           body.Id,
		FirstName:    body.FirstName,
		LastName:     body.LastName,
		BrithDate:    body.BrithDate,
		PhoneNumber:  body.PhoneNumber,
		Gender:       body.Gender,
		AccessToken:  access,
		RefreshToken: refresh,
	})
}

// GetUserByID
// @Summary GetUserByID
// @Description Api for GetUserByID
// @Security ApiKeyAuth
// @Tags User
// @Accept json
// @Produce json
// @Param id  query string true "id"
// @Param is_active query bool false "is_active"
// @Success 200 {object} model_user_service.GetUserResp
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/user/get [GET]
func (h *HandlerV1) GetUserByID(c *gin.Context) {

	id := c.Query("id")
	isActive := c.Query("is_active")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	response, err := h.serviceManager.UserService().UserService().Get(
		ctx, &pb.GetUserReq{
			Field:    "id",
			Value:    id,
			IsActive: !cast.ToBool(isActive),
		})

	if e.HandleError(c, err, h.log, http.StatusBadRequest, "cannot get user by id") {
		return
	}
	// Parse the birthdate from the service response using the correct layout
	birthDate, err := time.Parse("2006-01-02 15:04:05 -0700 MST", response.BirthDate)
	if err != nil {
		e.HandleError(c, err, h.log, http.StatusInternalServerError, "invalid birthdate format")
		return
	}

	// Format the birthdate as "2006-01-02"
	formattedBirthDate := birthDate.Format("2006-01-02")

	resp := model_user_service.GetUserResp{
		Id:          response.Id,
		UserOrder:   response.UserOrder,
		FirstName:   response.FirstName,
		LastName:    response.LastName,
		BrithDate:   formattedBirthDate, // Use the formatted birthdate
		PhoneNumber: response.PhoneNumber,
		Password:    response.Password,
		Gender:      response.Gender,
		ImageUrl:    response.ImageUrl,
		CreatedAt:   response.CreatedAt,
		UpdatedAt:   response.UpdatedAt,
	}

	c.JSON(http.StatusOK, resp)
}

// ListUsers
// @Summary ListUsers
// @Description Api for ListUsers
// @Security ApiKeyAuth
// @Tags User
// @Accept json
// @Produce json
// @Param Page  query string false "Page"
// @Param Limit query string false "Limit"
// @Param searchField query string false "searchField" Enums(first_name, last_name, gender, phone_number, created_at)
// @Param Value query string false "Value"
// @Param OrderBy query string false "OrderBy"
// @Param is_active query bool false "is_active"
// @Success 200 {object} model_user_service.GetUserResp
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/user [GET]
func (h *HandlerV1) ListUsers(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	page := c.Query("Page")
	limit := c.Query("Limit")
	field := c.Query("searchField")
	value := c.Query("Value")
	orderBy := c.Query("OrderBy")
	isActive := c.Query("is_active")
	pageInt, limitInt, err := e.ParseQueryParams(page, limit)
	if e.HandleError(c, err, h.log, http.StatusBadRequest, "failed to list users") {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	response, err := h.serviceManager.UserService().UserService().ListUsers(
		ctx, &pb.ListUsersReq{
			Page:     pageInt,
			Limit:    limitInt,
			IsActive: !cast.ToBool(isActive),
			Value:    value,
			Field:    field,
			OrderBy:  orderBy,
		})
	if e.HandleError(c, err, h.log, http.StatusInternalServerError, SERVICE_ERROR) {
		return
	}

	var users model_user_service.ListUserResp

	for _, in := range response.Users {
		user := model_user_service.GetUserResp{
			Id:          in.Id,
			UserOrder:   in.UserOrder,
			FirstName:   in.FirstName,
			LastName:    in.LastName,
			BrithDate:   in.BirthDate,
			PhoneNumber: in.PhoneNumber,
			Password:    in.Password,
			Gender:      in.Gender,
			ImageUrl:    in.ImageUrl,
			CreatedAt:   in.CreatedAt,
			UpdatedAt:   in.UpdatedAt,
		}

		users.Users = append(users.Users, user)
		users.Count = response.Count
	}

	c.JSON(http.StatusOK, users)
}

// UpdateUser
// @Summary UpdateUser
// @Description Api for UpdateUser
// @Security ApiKeyAuth
// @Tags User
// @Accept json
// @Produce json
// @Param UpdUserReq body model_user_service.UpdUserReq true "UpdUserReq"
// @Success 200 {object} model_user_service.GetUserResp
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/user [PUT]
func (h *HandlerV1) UpdateUser(c *gin.Context) {
	var (
		body        model_user_service.UpdUserResp
		jspbMarshal protojson.MarshalOptions
	)
	jspbMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if e.HandleError(c, err, h.log, http.StatusBadRequest, INVALID_REQUET_BODY) {
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	req := &pb.User{
		Id:        body.Id,
		FirstName: body.FirstName,
		LastName:  body.LastName,
		BirthDate: body.BrithDate,
		Gender:    body.Gender,
		ImageUrl:  body.ImageUrl,
	}

	response, err := h.serviceManager.UserService().UserService().Update(ctx, req)
	if e.HandleError(c, err, h.log, http.StatusInternalServerError, SERVICE_ERROR) {
		return
	}

	resp := model_user_service.UpdUserResp{
		Id:        response.Id,
		FirstName: response.FirstName,
		LastName:  response.LastName,
		BrithDate: response.BirthDate,
		Gender:    response.Gender,
		ImageUrl:  response.ImageUrl,
		UpdatedAt: response.UpdatedAt,
	}

	c.JSON(http.StatusOK, resp)
}

// DeleteUser
// @Summary DeleteUser
// @Description Api for DeleteUser
// @Security ApiKeyAuth
// @Tags User
// @Accept json
// @Produce json
// @Param id  query string true "id"
// @Param hard_delete query bool false "hard_delete"
// @Success 200 {object} model_user_service.CheckUserFieldResp
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/user [DELETE]
func (h *HandlerV1) DeleteUser(c *gin.Context) {
	var jspbMarshal protojson.MarshalOptions
	jspbMarshal.UseProtoNames = true

	id := c.Query("id")
	hard_delete := c.Query("hard_delete")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	response, err := h.serviceManager.UserService().UserService().Delete(
		ctx, &pb.DeleteUserReq{
			Field:    "id",
			Value:    id,
			IsActive: cast.ToBool(hard_delete),
		})

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, SERVICE_ERROR) {
		return
	}

	c.JSON(http.StatusOK, response)
}
