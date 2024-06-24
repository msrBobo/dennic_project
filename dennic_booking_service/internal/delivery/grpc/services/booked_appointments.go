package services

import (
	pb "booking_service/genproto/booking_service"
	_ "booking_service/internal/delivery/grpc"
	appointment "booking_service/internal/entity/booked_appointments"
	"booking_service/internal/pkg/otlp"
	"booking_service/internal/usecase"
	"context"
	"time"

	"github.com/rickb777/date"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

const (
	serviceNameAppointments     = "AppointmentsService"
	spanNameAppointmentsService = "AppointmentsService"
)

type BookingAppointments struct {
	logger                   *zap.Logger
	bookedAppointmentUseCase usecase.BookedAppointments
}

func BookingAppointmentsNewRPC(logger *zap.Logger, AppointmentUsaCase usecase.BookedAppointments) *BookingAppointments {

	return &BookingAppointments{
		logger:                   logger,
		bookedAppointmentUseCase: AppointmentUsaCase,
	}
}

func (r *BookingAppointments) CreateAppointment(ctx context.Context, req *pb.CreateAppointmentReq) (*pb.Appointment, error) {
	ctx, span := otlp.Start(ctx, serviceNameAppointments, spanNameAppointmentsService+"Create")
	span.SetAttributes(
		attribute.Key("doctor_id").String(req.DoctorId),
	)
	defer span.End()

	Date, err := date.AutoParse(req.AppointmentDate)
	if err != nil {
		return nil, err
	}
	Time, err := time.Parse("15:04:05", req.AppointmentTime)
	if err != nil {
		return nil, err
	}
	expTime, err := time.Parse("15:04:05", req.ExpiresAt)
	if err != nil {
		return nil, err
	}
	res, err := r.bookedAppointmentUseCase.CreateAppointment(ctx, &appointment.CreateAppointment{
		DepartmentId:    req.DepartmentId,
		DoctorId:        req.DoctorId,
		PatientId:       req.PatientId,
		ServiceId:       req.DoctorServiceId,
		AppointmentDate: Date,
		AppointmentTime: Time,
		Duration:        req.Duration,
		Key:             req.Key,
		ExpiresAt:       expTime,
		PatientProblem:  req.PatientProblem,
		Status:          req.Status,
		PaymentType:     req.PaymentType,
		PaymentAmount:   float64(req.PaymentAmount),
		UserId:          req.UserId,
	})

	if err != nil {
		return nil, err
	}

	return &pb.Appointment{
		Id:              res.Id,
		DepartmentId:    res.DepartmentId,
		DoctorId:        res.DoctorId,
		PatientId:       res.PatientId,
		DoctorServiceId: res.ServiceId,
		UserId:          res.UserId,
		AppointmentDate: res.AppointmentDate.String(),
		AppointmentTime: res.AppointmentTime.Format("15:04:05"),
		Duration:        res.Duration,
		Key:             res.Key,
		ExpiresAt:       res.ExpiresAt.Format("2006-01-02 15:04:05"),
		PatientProblem:  req.PatientProblem,
		Status:          res.Status,
		PaymentType:     res.PaymentType,
		PaymentAmount:   float32(res.PaymentAmount),
		CreatedAt:       res.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       res.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt:       res.DeletedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (r *BookingAppointments) GetAppointment(ctx context.Context, req *pb.AppointmentFieldValueReq) (*pb.Appointment, error) {
	ctx, span := otlp.Start(ctx, serviceNameAppointments, spanNameAppointmentsService+"Get")
	span.SetAttributes(
		attribute.Key(req.Field).String(req.Value),
	)
	defer span.End()

	res, err := r.bookedAppointmentUseCase.GetAppointment(ctx, &appointment.FieldValueReq{
		Field:        req.Field,
		Value:        req.Value,
		DeleteStatus: req.IsActive,
	})

	if err != nil {
		return nil, err
	}

	return &pb.Appointment{
		Id:              res.Id,
		DepartmentId:    res.DepartmentId,
		DoctorId:        res.DoctorId,
		PatientId:       res.PatientId,
		DoctorServiceId: res.ServiceId,
		AppointmentDate: res.AppointmentDate.String(),
		AppointmentTime: res.AppointmentTime.Format("15:04:05"),
		Duration:        res.Duration,
		Key:             res.Key,
		ExpiresAt:       res.ExpiresAt.Format("2006-01-02 15:04:05"),
		PatientProblem:  res.PatientProblem,
		Status:          res.Status,
		PaymentType:     res.PaymentType,
		PaymentAmount:   float32(res.PaymentAmount),
		UserId:          res.UserId,
		CreatedAt:       res.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       res.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt:       res.DeletedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (r *BookingAppointments) GetAllAppointment(ctx context.Context, req *pb.GetAllAppointmentsReq) (*pb.Appointments, error) {
	ctx, span := otlp.Start(ctx, serviceNameAppointments, spanNameAppointmentsService+"List")
	defer span.End()

	var appointmentsRes pb.Appointments

	allAppointment, err := r.bookedAppointmentUseCase.GetAllAppointment(ctx, &appointment.GetAllAppointment{
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

	for _, appoint := range allAppointment.Appointments {
		var appointmentRes pb.Appointment
		appointmentRes.Id = appoint.Id
		appointmentRes.DepartmentId = appoint.DepartmentId
		appointmentRes.UserId = appoint.UserId
		appointmentRes.DoctorId = appoint.DoctorId
		appointmentRes.PatientId = appoint.PatientId
		appointmentRes.DoctorServiceId = appoint.ServiceId
		appointmentRes.AppointmentDate = appoint.AppointmentDate.String()
		appointmentRes.AppointmentTime = appoint.AppointmentTime.Format("15:04:05")
		appointmentRes.Duration = appoint.Duration
		appointmentRes.Key = appoint.Key
		appointmentRes.PatientProblem = appoint.PatientProblem
		appointmentRes.Status = appoint.Status
		appointmentRes.PaymentType = appoint.PaymentType
		appointmentRes.PaymentAmount = float32(appoint.PaymentAmount)
		appointmentRes.ExpiresAt = appoint.ExpiresAt.Format("2006-01-02 15:04:05")
		appointmentRes.CreatedAt = appoint.CreatedAt.Format("2006-01-02 15:04:05")
		appointmentRes.UpdatedAt = appoint.UpdatedAt.Format("2006-01-02 15:04:05")
		appointmentRes.DeletedAt = appoint.DeletedAt.Format("2006-01-02 15:04:05")

		appointmentsRes.Appointments = append(appointmentsRes.Appointments, &appointmentRes)
	}
	appointmentsRes.Count = allAppointment.Count

	return &appointmentsRes, nil
}

func (r *BookingAppointments) GetFilteredAppointments(ctx context.Context, req *pb.GetFilteredRequest) (*pb.Appointments, error) {
	ctx, span := otlp.Start(ctx, serviceNameAppointments, spanNameAppointmentsService+"Filtered List")
	defer span.End()

	var appointmentsRes pb.Appointments

	allAppointment, err := r.bookedAppointmentUseCase.GetFilteredAppointments(ctx, &appointment.GetFilteredRequest{
		Page:         req.Page,
		Limit:        req.Limit,
		DeleteStatus: req.IsActive,
		Field:        req.Field,
		Value:        req.Value,
		OrderBy:      req.OrderBy,
		Status:       req.Status,
	})
	if err != nil {
		return nil, err
	}

	for _, appoint := range allAppointment.Appointments {
		var appointmentRes pb.Appointment
		appointmentRes.Id = appoint.Id
		appointmentRes.DepartmentId = appoint.DepartmentId
		appointmentRes.DoctorId = appoint.DoctorId
		appointmentRes.PatientId = appoint.PatientId
		appointmentRes.UserId = appoint.UserId
		appointmentRes.DoctorServiceId = appoint.ServiceId
		appointmentRes.AppointmentDate = appoint.AppointmentDate.String()
		appointmentRes.AppointmentTime = appoint.AppointmentTime.Format("15:04:05")
		appointmentRes.Duration = appoint.Duration
		appointmentRes.Key = appoint.Key
		appointmentRes.PatientProblem = appoint.PatientProblem
		appointmentRes.Status = appoint.Status
		appointmentRes.PaymentType = appoint.PaymentType
		appointmentRes.PaymentAmount = float32(appoint.PaymentAmount)
		appointmentRes.ExpiresAt = appoint.ExpiresAt.Format("2006-01-02 15:04:05")
		appointmentRes.CreatedAt = appoint.CreatedAt.Format("2006-01-02 15:04:05")
		appointmentRes.UpdatedAt = appoint.UpdatedAt.Format("2006-01-02 15:04:05")
		appointmentRes.DeletedAt = appoint.DeletedAt.Format("2006-01-02 15:04:05")

		appointmentsRes.Appointments = append(appointmentsRes.Appointments, &appointmentRes)
	}
	appointmentsRes.Count = allAppointment.Count

	return &appointmentsRes, nil
}

func (r *BookingAppointments) UpdateAppointment(ctx context.Context, req *pb.UpdateAppointmentReq) (*pb.Appointment, error) {
	ctx, span := otlp.Start(ctx, serviceNameAppointments, spanNameAppointmentsService+"Update")
	span.SetAttributes(
		attribute.Key(req.Field).String(req.Value),
	)
	defer span.End()

	reqDate, err := date.AutoParse(req.AppointmentDate)
	if err != nil {
		return nil, err
	}

	reqTime, err := time.Parse("15:04:05", req.AppointmentTime)
	if err != nil {
		return nil, err
	}
	expreqTime, err := time.Parse("2006-01-02 15:04:05", req.ExpiresAt)
	if err != nil {
		return nil, err
	}

	res, err := r.bookedAppointmentUseCase.UpdateAppointment(ctx, &appointment.UpdateAppointment{
		Field:           req.Field,
		Value:           req.Value,
		DepartmentId:    req.DepartmentId,
		DoctorId:        req.DoctorId,
		PatientId:       req.PatientId,
		ServiceId:       req.DoctorServiceId,
		AppointmentDate: reqDate,
		AppointmentTime: reqTime,
		Duration:        req.Duration,
		Key:             req.Key,
		ExpiresAt:       expreqTime,
		Status:          req.Status,
		PatientProblem:  req.PatientProblem,
		PaymentType:     req.PaymentType,
		PaymentAmount:   float64(req.PaymentAmount),
		UserId:          req.UserId,
	})
	if err != nil {
		return nil, err
	}

	return &pb.Appointment{
		Id:              res.Id,
		DepartmentId:    res.DepartmentId,
		DoctorId:        res.DoctorId,
		PatientId:       res.PatientId,
		DoctorServiceId: res.ServiceId,
		AppointmentDate: res.AppointmentDate.String(),
		AppointmentTime: res.AppointmentTime.Format("15:04:05"),
		Duration:        res.Duration,
		Key:             res.Key,
		ExpiresAt:       res.ExpiresAt.Format("2006-01-02 15:04:05"),
		PatientProblem:  res.PatientProblem,
		Status:          res.Status,
		PaymentType:     res.PaymentType,
		PaymentAmount:   float32(res.PaymentAmount),
		UserId:          res.UserId,
		CreatedAt:       res.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:       res.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt:       res.DeletedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (r *BookingAppointments) DeleteAppointment(ctx context.Context, req *pb.AppointmentFieldValueReq) (*pb.DeleteAppointmentStatus, error) {
	ctx, span := otlp.Start(ctx, serviceNameAppointments, spanNameAppointmentsService+"Delete")
	span.SetAttributes(
		attribute.Key(req.Field).String(req.Value),
	)
	defer span.End()

	res, err := r.bookedAppointmentUseCase.DeleteAppointment(ctx, &appointment.FieldValueReq{
		Field:        req.Field,
		Value:        req.Value,
		DeleteStatus: req.IsActive,
	})

	if err != nil {
		return &pb.DeleteAppointmentStatus{Status: res.Status}, err
	}

	return &pb.DeleteAppointmentStatus{Status: res.Status}, err
}
