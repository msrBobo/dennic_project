package services

import (
	pb "booking_service/genproto/booking_service"
	"booking_service/internal/entity/doctor_availability"
	"booking_service/internal/pkg/otlp"
	"booking_service/internal/usecase"
	"context"
	"github.com/rickb777/date"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	"time"
)

const (
	serviceNameDoctorAvailability     = "DoctorAvailabilityService"
	spanNameDoctorAvailabilityService = "DoctorAvailabilityService"
)

type BookingDoctorAvailability struct {
	logger                          *zap.Logger
	bookedDoctorAvailabilityUseCase usecase.DoctorAvailability
}

func BookingDoctorAvailabilityNewRPC(
	logger *zap.Logger, DoctorAvailability usecase.DoctorAvailability) *BookingDoctorAvailability {
	return &BookingDoctorAvailability{
		logger:                          logger,
		bookedDoctorAvailabilityUseCase: DoctorAvailability,
	}
}

func (r *BookingDoctorAvailability) CreateDoctorTime(ctx context.Context, req *pb.CreateDoctorTimeReq) (*pb.DoctorTime, error) {
	ctx, span := otlp.Start(ctx, serviceNameDoctorAvailability, spanNameDoctorAvailabilityService+"Create")
	span.SetAttributes(
		attribute.Key("doctor_id").String(req.DoctorId),
	)
	defer span.End()

	reqDate, err := date.AutoParse(req.DoctorDate)
	if err != nil {
		return nil, err
	}
	reqStartTime, err := time.Parse("15:04:05", req.StartTime)
	if err != nil {
		return nil, err
	}
	reqEndTime, err := time.Parse("15:04:05", req.EndTime)
	if err != nil {
		return nil, err
	}

	res, err := r.bookedDoctorAvailabilityUseCase.CreateDoctorAvailability(ctx, &doctor_availability.CreateDoctorAvailability{
		DepartmentId: req.DepartmentId,
		DoctorId:     req.DoctorId,
		DoctorDate:   reqDate,
		StartTime:    reqStartTime,
		EndTime:      reqEndTime,
		Status:       req.Status,
	})

	if err != nil {
		return nil, err
	}

	return &pb.DoctorTime{
		Id:           res.Id,
		DepartmentId: res.DepartmentId,
		DoctorId:     res.DoctorId,
		DoctorDate:   res.DoctorDate.String(),
		StartTime:    res.StartTime.Format("15:04:05"),
		EndTime:      res.EndTime.Format("15:04:05"),
		Status:       res.Status,
		CreatedAt:    res.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    res.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt:    res.DeletedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (r *BookingDoctorAvailability) GetDoctorTime(ctx context.Context, req *pb.DoctorTimeFieldValueReq) (*pb.DoctorTime, error) {
	ctx, span := otlp.Start(ctx, serviceNameDoctorAvailability, spanNameDoctorAvailabilityService+"Get")
	span.SetAttributes(
		attribute.Key(req.Field).String(req.Value),
	)
	defer span.End()

	res, err := r.bookedDoctorAvailabilityUseCase.GetDoctorAvailability(ctx, &doctor_availability.FieldValueReq{
		Field:        req.Field,
		Value:        req.Value,
		DeleteStatus: req.IsActive,
	})

	if err != nil {
		return nil, err
	}

	return &pb.DoctorTime{
		Id:           res.Id,
		DepartmentId: res.DepartmentId,
		DoctorId:     res.DoctorId,
		DoctorDate:   res.DoctorDate.String(),
		StartTime:    res.StartTime.Format("15:04:05"),
		EndTime:      res.EndTime.Format("15:04:05"),
		Status:       res.Status,
		CreatedAt:    res.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    res.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt:    res.DeletedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (r *BookingDoctorAvailability) GetAllDoctorTimes(ctx context.Context, req *pb.GetAllDoctorTimesReq) (*pb.DoctorTimes, error) {
	ctx, span := otlp.Start(ctx, serviceNameDoctorAvailability, spanNameDoctorAvailabilityService+"List")
	defer span.End()

	var docAvails pb.DoctorTimes

	allDocAvails, err := r.bookedDoctorAvailabilityUseCase.GetAllDoctorAvailability(ctx, &doctor_availability.GetAllReq{
		Page:         req.Page,
		Limit:        req.Limit,
		DeleteStatus: req.IsActive,
		Field:        req.Field,
		Value:        req.Value,
		OrderBy:      req.OrderBy,
	})

	if err != nil {
		return nil, err
	}

	for _, availability := range allDocAvails.DoctorAvailabilitys {
		var docAvail pb.DoctorTime
		docAvail.Id = availability.Id
		docAvail.DepartmentId = availability.DepartmentId
		docAvail.DoctorId = availability.DoctorId
		docAvail.DoctorDate = availability.DoctorDate.String()
		docAvail.StartTime = availability.StartTime.Format("15:04:05")
		docAvail.EndTime = availability.EndTime.Format("15:04:05")
		docAvail.Status = availability.Status
		docAvail.CreatedAt = availability.CreatedAt.Format("2006-01-02 15:04:05")
		docAvail.UpdatedAt = availability.UpdatedAt.Format("2006-01-02 15:04:05")
		docAvail.DeletedAt = availability.DeletedAt.Format("2006-01-02 15:04:05")
		docAvails.DoctorTimes = append(docAvails.DoctorTimes, &docAvail)
	}

	docAvails.Count = allDocAvails.Count

	return &docAvails, nil
}

func (r *BookingDoctorAvailability) UpdateDoctorTime(ctx context.Context, req *pb.UpdateDoctorTimeReq) (*pb.DoctorTime, error) {
	ctx, span := otlp.Start(ctx, serviceNameDoctorAvailability, spanNameDoctorAvailabilityService+"Update")
	span.SetAttributes(
		attribute.Key(req.Field).String(req.Value),
	)
	defer span.End()

	reqDate, err := date.AutoParse(req.DoctorDate)
	if err != nil {
		return nil, err
	}
	reqStartTime, err := time.Parse("15:04:05", req.StartTime)
	if err != nil {
		return nil, err
	}
	reqEndTime, err := time.Parse("15:04:05", req.EndTime)
	if err != nil {
		return nil, err
	}

	res, err := r.bookedDoctorAvailabilityUseCase.UpdateDoctorAvailability(ctx, &doctor_availability.UpdateDoctorAvailability{
		Field:        req.Field,
		Value:        req.Value,
		DepartmentId: req.DepartmentId,
		DoctorId:     req.DoctorId,
		DoctorDate:   reqDate,
		StartTime:    reqStartTime,
		EndTime:      reqEndTime,
		Status:       req.Status,
	})

	if err != nil {
		return nil, err
	}

	return &pb.DoctorTime{
		Id:           res.Id,
		DepartmentId: res.DepartmentId,
		DoctorId:     res.DoctorId,
		DoctorDate:   res.DoctorDate.String(),
		StartTime:    res.StartTime.Format("15:04:05"),
		EndTime:      res.EndTime.Format("15:04:05"),
		Status:       res.Status,
		CreatedAt:    res.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    res.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt:    res.DeletedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (r *BookingDoctorAvailability) DeleteDoctorTime(ctx context.Context, req *pb.DoctorTimeFieldValueReq) (*pb.DoctorTimeDeleteStatus, error) {
	ctx, span := otlp.Start(ctx, serviceNameDoctorAvailability, spanNameDoctorAvailabilityService+"Delete")
	span.SetAttributes(
		attribute.Key(req.Field).String(req.Value),
	)
	defer span.End()

	res, err := r.bookedDoctorAvailabilityUseCase.DeleteDoctorAvailability(ctx, &doctor_availability.FieldValueReq{
		Field:        req.Field,
		Value:        req.Value,
		DeleteStatus: req.IsActive,
	})

	if err != nil {
		return &pb.DoctorTimeDeleteStatus{Status: res.Status}, err
	}

	return &pb.DoctorTimeDeleteStatus{Status: res.Status}, err
}

func (r *BookingDoctorAvailability) GeDoctorTimeStatus(ctx context.Context, req *pb.DoctorTimeStatusReq) (*pb.DoctorTimeStatus, error) {
	ctx, span := otlp.Start(ctx, serviceNameDoctorAvailability, spanNameDoctorAvailabilityService+"GetStatus")
	span.SetAttributes(
		attribute.Key("doctor_id").String(req.DoctorId),
	)
	defer span.End()

	res, err := r.bookedDoctorAvailabilityUseCase.GeDoctorTimeStatus(ctx, &doctor_availability.DoctorTimeStatusReq{
		DepartmentID: req.DepartmentId,
		DoctorID:     req.DoctorId,
		DoctorDate:   req.DoctorDate,
		StartTime:    req.StartTime,
	})

	if err != nil {
		return nil, err
	}

	return &pb.DoctorTimeStatus{
		Status: res.Status,
	}, nil
}

func (r *BookingDoctorAvailability) UpdateTimeStatus(ctx context.Context, req *pb.DoctorTimeUpdateStatusReq) (*pb.DoctorTimeStatus, error) {
	ctx, span := otlp.Start(ctx, serviceNameDoctorAvailability, spanNameDoctorAvailabilityService+"UpdateStatus")
	span.SetAttributes(
		attribute.Key("doctor_id").String(req.DoctorId),
	)
	defer span.End()

	res, err := r.bookedDoctorAvailabilityUseCase.UpdateTimeStatus(ctx, &doctor_availability.DoctorTimeUpdateStatusReq{
		DepartmentID: req.DepartmentId,
		DoctorID:     req.DoctorId,
		DoctorDate:   req.DoctorDate,
		StartTime:    req.StartTime,
		Status:       req.Status,
	})

	if err != nil {
		return nil, err
	}

	return &pb.DoctorTimeStatus{
		Status: res.Status,
	}, nil
}

func (r *BookingDoctorAvailability) GetDoctorTimeByDoctorID(ctx context.Context, req *pb.DoctorTimeFieldValueByIdReq) (*pb.Hours, error) {
	ctx, span := otlp.Start(ctx, serviceNameDoctorAvailability, spanNameDoctorAvailabilityService+"UpdateStatus")
	span.SetAttributes(
		attribute.Key("doctor_id").String(req.DoctorId),
	)
	defer span.End()

	var hours pb.Hours

	res, err := r.bookedDoctorAvailabilityUseCase.GetDoctorTimeByDoctorID(ctx, &doctor_availability.DoctorTimeFieldValueByIdReq{
		DoctorID: req.DoctorId,
		Date:     req.DoctorDate,
	})

	if err != nil {
		return nil, err
	}

	for _, hour := range res.Hours {
		hours.Hours = append(hours.Hours, hour)
	}

	return &pb.Hours{
		Hours: res.Hours,
	}, nil
}
