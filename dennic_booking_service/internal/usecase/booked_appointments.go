package usecase

import (
	_ "booking_service/genproto/booking_service"
	appointment "booking_service/internal/entity/booked_appointments"
	"booking_service/internal/pkg/otlp"
	"context"
	"time"
)

const (
	serviceNameAppointments = "AppointmentsService"
	spanNameAppointments    = "AppointmentsUsecase"
)

// BookedAppointmentsUseCase -.
type BookedAppointmentsUseCase struct {
	repo       BookedAppointments
	ctxTimeout time.Duration
}

// NewBookedAppointments -.
func NewBookedAppointments(r BookedAppointments, ctxTimeout time.Duration) *BookedAppointmentsUseCase {
	return &BookedAppointmentsUseCase{
		repo:       r,
		ctxTimeout: ctxTimeout,
	}
}

func (r *BookedAppointmentsUseCase) CreateAppointment(ctx context.Context, req *appointment.CreateAppointment) (*appointment.Appointment, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameAppointments, spanNameAppointments+"Create")
	span.End()

	return r.repo.CreateAppointment(ctx, req)
}

func (r *BookedAppointmentsUseCase) GetAppointment(ctx context.Context, req *appointment.FieldValueReq) (*appointment.Appointment, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameAppointments, spanNameAppointments+"Get")
	span.End()

	return r.repo.GetAppointment(ctx, req)
}

func (r *BookedAppointmentsUseCase) GetAllAppointment(ctx context.Context, req *appointment.GetAllAppointment) (*appointment.AppointmentsType, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameAppointments, spanNameAppointments+"List")
	span.End()

	return r.repo.GetAllAppointment(ctx, req)
}

func (r *BookedAppointmentsUseCase) GetFilteredAppointments(ctx context.Context, req *appointment.GetFilteredRequest) (*appointment.AppointmentsType, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameAppointments, spanNameAppointments+"List")
	span.End()

	return r.repo.GetFilteredAppointments(ctx, req)
}

func (r *BookedAppointmentsUseCase) UpdateAppointment(ctx context.Context, req *appointment.UpdateAppointment) (*appointment.Appointment, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameAppointments, spanNameAppointments+"Update")
	span.End()

	return r.repo.UpdateAppointment(ctx, req)
}

func (r *BookedAppointmentsUseCase) DeleteAppointment(ctx context.Context, req *appointment.FieldValueReq) (*appointment.StatusRes, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameAppointments, spanNameAppointments+"Delete")
	span.End()

	return r.repo.DeleteAppointment(ctx, req)
}
