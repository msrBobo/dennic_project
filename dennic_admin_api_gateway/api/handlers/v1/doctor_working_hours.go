package v1

import (
	"context"
	e "dennic_admin_api_gateway/api/handlers/regtool"
	"dennic_admin_api_gateway/api/models"
	"dennic_admin_api_gateway/api/models/model_healthcare_service"
	pb "dennic_admin_api_gateway/genproto/healthcare-service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/spf13/cast"
	"google.golang.org/protobuf/encoding/protojson"
)

// CreateDoctorWorkingHours ...
// @Summary CreateDoctorWorkingHours
// @Description CreateDoctorWorkingHours - Api for crete doctor_working_hours
// @Tags Doctor Working Hours
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param DoctorWorkingHoursReq body model_healthcare_service.DoctorWorkingHoursReq true "DoctorServiceReq"
// @Success 200 {object} model_healthcare_service.DoctorWorkingHoursRes
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/doctor-working-hours/create [post]
func (h *HandlerV1) CreateDoctorWorkingHours(c *gin.Context) {
	var (
		body        model_healthcare_service.DoctorWorkingHoursReq
		jsonMarshal protojson.MarshalOptions
	)
	jsonMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)

	if e.HandleError(c, err, h.log, http.StatusBadRequest, "CreateDoctorWorkingHours") {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	dwh, err := h.serviceManager.HealthcareService().DoctorWorkingHoursService().CreateDoctorWorkingHours(ctx, &pb.DoctorWorkingHours{
		DoctorId:   body.DoctorId,
		DayOfWeek:  body.DayOfWeek,
		StartTime:  body.StartTime,
		FinishTime: body.FinishTime,
	})

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, "CreateDoctorWorkingHours") {
		return
	}

	c.JSON(http.StatusOK, model_healthcare_service.DoctorWorkingHoursRes{
		Id:         dwh.Id,
		DoctorId:   dwh.DoctorId,
		DayOfWeek:  dwh.DayOfWeek,
		StartTime:  dwh.StartTime,
		FinishTime: dwh.FinishTime,
		CreatedAt:  dwh.CreatedAt,
		UpdatedAt:  e.UpdateTimeFilter(dwh.UpdatedAt),
		DeletedAt:  e.UpdateTimeFilter(dwh.DeletedAt),
	})
}

// GetDoctorWorkingHours ...
// @Summary GetDoctorWorkingHours
// @Description GetDoctorWorkingHours - Api for get doctor_working_hours
// @Tags Doctor Working Hours
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id query string true "id"
// @Param is_active query bool false "is_active"
// @Success 200 {object} model_healthcare_service.DoctorWorkingHoursRes
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/doctor-working-hours/get [get]
func (h *HandlerV1) GetDoctorWorkingHours(c *gin.Context) {
	id := c.Query("id")
	isActive := c.Query("is_active")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	dwh, err := h.serviceManager.HealthcareService().DoctorWorkingHoursService().GetDoctorWorkingHoursById(ctx, &pb.GetReqInt{
		Field:    "id",
		Value:    id,
		IsActive: !cast.ToBool(isActive),
	})

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, "GetDoctorWorkingHours") {
		return
	}

	c.JSON(http.StatusOK, model_healthcare_service.DoctorWorkingHoursRes{
		Id:         dwh.Id,
		DoctorId:   dwh.DoctorId,
		DayOfWeek:  dwh.DayOfWeek,
		StartTime:  dwh.StartTime,
		FinishTime: dwh.FinishTime,
		CreatedAt:  dwh.CreatedAt,
		UpdatedAt:  e.UpdateTimeFilter(dwh.UpdatedAt),
		DeletedAt:  e.UpdateTimeFilter(dwh.DeletedAt),
	})
}

// ListDoctorWorkingHours ...
// @Summary ListDoctorWorkingHours
// @Description ListDoctorWorkingHours - Api for list doctor_working_hours
// @Tags Doctor Working Hours
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param ListReq query models.ListReq false "ListReq"
// @Param search query string false "search" Enums(day_of_week) "search"
// @Success 200 {object} model_healthcare_service.ListDoctorWorkingHours
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/doctor-working-hours [get]
func (h *HandlerV1) ListDoctorWorkingHours(c *gin.Context) {
	search := c.Query("search")
	value := c.Query("value")
	limit := c.Query("limit")
	page := c.Query("page")
	orderBy := c.Query("orderBy")
	isActive := c.Query("is_active")

	pageInt, limitInt, err := e.ParseQueryParams(page, limit)
	if e.HandleError(c, err, h.log, http.StatusBadRequest, "ListDoctorWorkingHours") {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	dwhs, err := h.serviceManager.HealthcareService().DoctorWorkingHoursService().GetAllDoctorWorkingHours(ctx, &pb.GetAllDoctorWorkingHoursReq{
		Field:    search,
		Value:    value,
		IsActive: !cast.ToBool(isActive),
		Page:     int64(pageInt),
		Limit:    int64(limitInt),
		OrderBy:  orderBy,
	})

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, "ListDoctorWorkingHours") {
		return
	}
	var dwhsRes model_healthcare_service.ListDoctorWorkingHours
	for _, dwhRes := range dwhs.Dwh {
		dwhsRes.ListDWH = append(dwhsRes.ListDWH, &model_healthcare_service.DoctorWorkingHoursRes{
			Id:         dwhRes.Id,
			DoctorId:   dwhRes.DoctorId,
			DayOfWeek:  dwhRes.DayOfWeek,
			StartTime:  dwhRes.StartTime,
			FinishTime: dwhRes.FinishTime,
			CreatedAt:  dwhRes.CreatedAt,
			UpdatedAt:  e.UpdateTimeFilter(dwhRes.UpdatedAt),
			DeletedAt:  e.UpdateTimeFilter(dwhRes.DeletedAt),
		})
	}

	c.JSON(http.StatusOK, model_healthcare_service.ListDoctorWorkingHours{
		Count:   dwhs.Count,
		ListDWH: dwhsRes.ListDWH,
	})
}

// UpdateDoctorWorkingHours ...
// @Summary UpdateDoctorWorkingHours
// @Description UpdateDoctorWorkingHours - Api for update doctor_working_hours
// @Tags Doctor Working Hours
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param UpdateDoctorWorkingHoursReq body model_healthcare_service.DoctorWorkingHoursReq true "UpdateDoctorWorkingHoursReq"
// @Success 200 {object} model_healthcare_service.DoctorWorkingHoursRes
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/doctor-working-hours [put]
func (h *HandlerV1) UpdateDoctorWorkingHours(c *gin.Context) {
	var (
		body        model_healthcare_service.DoctorWorkingHoursReq
		jsonMarshal protojson.MarshalOptions
	)
	jsonMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)

	id := c.Query("id")

	if e.HandleError(c, err, h.log, http.StatusBadRequest, "UpdateDoctorWorkingHours") {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	dwh, err := h.serviceManager.HealthcareService().DoctorWorkingHoursService().UpdateDoctorWorkingHours(ctx, &pb.DoctorWorkingHours{
		Id:         cast.ToInt32(id),
		DoctorId:   body.DoctorId,
		DayOfWeek:  body.DayOfWeek,
		StartTime:  body.StartTime,
		FinishTime: body.FinishTime,
	})

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, "UpdateDoctorWorkingHours") {
		return
	}

	c.JSON(http.StatusOK, model_healthcare_service.DoctorWorkingHoursRes{
		Id:         dwh.Id,
		DoctorId:   dwh.DoctorId,
		DayOfWeek:  dwh.DayOfWeek,
		StartTime:  dwh.StartTime,
		FinishTime: dwh.FinishTime,
		CreatedAt:  dwh.CreatedAt,
		UpdatedAt:  e.UpdateTimeFilter(dwh.UpdatedAt),
		DeletedAt:  e.UpdateTimeFilter(dwh.DeletedAt),
	})
}

// DeleteDoctorWorkingHours ...
// @Summary DeleteDoctorWorkingHours
// @Description DeleteDoctorWorkingHours - Api for delete doctor_working_hours
// @Tags Doctor Working Hours
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id query string true "id"
// @Param hard_delete query bool true "hard_delete"
// @Success 200 {object} models.StatusRes
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/doctor-working-hours [delete]
func (h *HandlerV1) DeleteDoctorWorkingHours(c *gin.Context) {
	value := c.Query("id")
	isActive := c.Query("hard_delete")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	status, err := h.serviceManager.HealthcareService().DoctorWorkingHoursService().DeleteDoctorWorkingHours(ctx, &pb.GetReqInt{
		Field:    "id",
		Value:    value,
		IsActive: cast.ToBool(isActive),
	})

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, "DeleteDoctorWorkingHours") {
		return
	}

	c.JSON(http.StatusOK, models.StatusRes{Status: status.Status})
}
