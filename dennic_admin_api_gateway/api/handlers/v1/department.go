package v1

import (
	"context"
	e "dennic_admin_api_gateway/api/handlers/regtool"
	"dennic_admin_api_gateway/api/models"
	"dennic_admin_api_gateway/api/models/model_healthcare_service"
	pb "dennic_admin_api_gateway/genproto/healthcare-service"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/spf13/cast"
	"google.golang.org/protobuf/encoding/protojson"
	"net/http"
	"time"
)

// CreateDepartment ...
// @Summary CreateDepartment
// @Description CreateDepartment - Api for crete department
// @Tags Department
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param DepartmentReq body model_healthcare_service.DepartmentReq true "DepartmentReq"
// @Success 200 {object} model_healthcare_service.DepartmentRes
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/department/create [post]
func (h *HandlerV1) CreateDepartment(c *gin.Context) {
	var (
		body        model_healthcare_service.DepartmentReq
		jsonMarshal protojson.MarshalOptions
	)
	jsonMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)

	if e.HandleError(c, err, h.log, http.StatusBadRequest, "CreateDepartment") {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	department, err := h.serviceManager.HealthcareService().DepartmentService().CreateDepartment(ctx, &pb.Department{
		Id:               uuid.NewString(),
		Name:             body.Name,
		Description:      body.Description,
		ImageUrl:         body.ImageUrl,
		FloorNumber:      body.FloorNumber,
		ShortDescription: body.ShortDescription,
	})

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, "CreateDepartment") {
		return
	}

	c.JSON(http.StatusOK, model_healthcare_service.DepartmentRes{
		Id:               department.Id,
		Order:            department.Order,
		Name:             department.Name,
		Description:      department.Description,
		ImageUrl:         department.ImageUrl,
		FloorNumber:      department.FloorNumber,
		ShortDescription: department.ShortDescription,
		CreatedAt:        department.CreatedAt,
		UpdatedAt:        e.UpdateTimeFilter(department.UpdatedAt),
		DeletedAt:        e.UpdateTimeFilter(department.DeletedAt),
	})
}

// GetDepartment ...
// @Summary GetDepartment
// @Description GetDepartment - Api for get department
// @Tags Department
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id query string true "id"
// @Param is_active query bool false "is_active"
// @Success 200 {object} model_healthcare_service.DepartmentRes
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/department/get [get]
func (h *HandlerV1) GetDepartment(c *gin.Context) {
	id := c.Query("id")
	isActive := c.Query("is_active")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	department, err := h.serviceManager.HealthcareService().DepartmentService().GetDepartmentById(ctx, &pb.GetReqStrDepartment{
		Field:    "id",
		Value:    id,
		IsActive: !cast.ToBool(isActive),
	})

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, "GetDepartment") {
		return
	}

	c.JSON(http.StatusOK, model_healthcare_service.DepartmentRes{
		Id:               department.Id,
		Order:            department.Order,
		Name:             department.Name,
		Description:      department.Description,
		ImageUrl:         department.ImageUrl,
		FloorNumber:      department.FloorNumber,
		ShortDescription: department.ShortDescription,
		CreatedAt:        department.CreatedAt,
		UpdatedAt:        e.UpdateTimeFilter(department.UpdatedAt),
		DeletedAt:        e.UpdateTimeFilter(department.DeletedAt),
	})
}

// ListDepartments ...
// @Summary ListDepartments
// @Description ListDepartments - Api for list department
// @Tags Department
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param ListReq query models.ListReq false "ListReq"
// @Param search query string false "search" Enums(name, description) "search"
// @Success 200 {object} model_healthcare_service.ListDepartments
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/department [get]
func (h *HandlerV1) ListDepartments(c *gin.Context) {
	search := c.Query("search")
	value := c.Query("value")
	limit := c.Query("limit")
	page := c.Query("page")
	orderBy := c.Query("orderBy")
	isActive := c.Query("is_active")
	pageInt, limitInt, err := e.ParseQueryParams(page, limit)
	if e.HandleError(c, err, h.log, http.StatusBadRequest, "ListDepartments") {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	departments, err := h.serviceManager.HealthcareService().DepartmentService().GetAllDepartments(ctx, &pb.GetAllDepartment{
		Field:    search,
		Value:    value,
		IsActive: !cast.ToBool(isActive),
		Page:     int64(pageInt),
		Limit:    int64(limitInt),
		OrderBy:  orderBy,
	})
	if e.HandleError(c, err, h.log, http.StatusInternalServerError, "ListDepartments") {
		return
	}
	var departmentsRes model_healthcare_service.ListDepartments
	for _, departmentRes := range departments.Departments {
		departmentsRes.Departments = append(departmentsRes.Departments, &model_healthcare_service.DepartmentRes{
			Id:               departmentRes.Id,
			Order:            departmentRes.Order,
			Name:             departmentRes.Name,
			Description:      departmentRes.Description,
			ImageUrl:         departmentRes.ImageUrl,
			FloorNumber:      departmentRes.FloorNumber,
			ShortDescription: departmentRes.ShortDescription,
			CreatedAt:        departmentRes.CreatedAt,
			UpdatedAt:        e.UpdateTimeFilter(departmentRes.UpdatedAt),
			DeletedAt:        e.UpdateTimeFilter(departmentRes.DeletedAt),
		})
	}

	c.JSON(http.StatusOK, model_healthcare_service.ListDepartments{
		Count:       int32(departments.Count),
		Departments: departmentsRes.Departments,
	})
}

// UpdateDepartment ...
// @Summary UpdateDepartment
// @Description UpdateDepartment - Api for update department
// @Tags Department
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param UpdateDepartmentReq body model_healthcare_service.DepartmentReq true "UpdateDepartmentReq"
// @Success 200 {object} model_healthcare_service.DepartmentReq
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/department [put]
func (h *HandlerV1) UpdateDepartment(c *gin.Context) {
	var (
		body        model_healthcare_service.DepartmentReq
		jsonMarshal protojson.MarshalOptions
	)
	jsonMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)

	if e.HandleError(c, err, h.log, http.StatusBadRequest, "UpdateDepartment") {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	department, err := h.serviceManager.HealthcareService().DepartmentService().UpdateDepartment(ctx, &pb.Department{
		Id:               body.Id,
		Name:             body.Name,
		Description:      body.Description,
		ImageUrl:         body.ImageUrl,
		FloorNumber:      body.FloorNumber,
		ShortDescription: body.ShortDescription,
	})

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, "UpdateDepartment") {
		return
	}

	c.JSON(http.StatusOK, model_healthcare_service.DepartmentRes{
		Id:               department.Id,
		Order:            department.Order,
		Name:             department.Name,
		Description:      department.Description,
		ImageUrl:         department.ImageUrl,
		FloorNumber:      department.FloorNumber,
		ShortDescription: department.ShortDescription,
		CreatedAt:        department.CreatedAt,
		UpdatedAt:        e.UpdateTimeFilter(department.UpdatedAt),
		DeletedAt:        e.UpdateTimeFilter(department.DeletedAt),
	})
}

// DeleteDepartment ...
// @Summary DeleteDepartment
// @Description DeleteDepartment - Api for delete department
// @Tags Department
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param id query string true "id"
// @Param hard_delete query bool true "hard_delete"
// @Success 200 {object} models.StatusRes
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/department [delete]
func (h *HandlerV1) DeleteDepartment(c *gin.Context) {
	value := c.Query("id")
	isActive := c.Query("hard_delete")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	status, err := h.serviceManager.HealthcareService().DepartmentService().DeleteDepartment(ctx, &pb.GetReqStrDepartment{
		Field:    "id",
		Value:    value,
		IsActive: cast.ToBool(isActive),
	})

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, "DeleteDepartment") {
		return
	}

	c.JSON(http.StatusOK, models.StatusRes{Status: status.Status})
}
