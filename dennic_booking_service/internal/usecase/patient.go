package usecase

import (
	"booking_service/internal/entity/patients"
	"booking_service/internal/pkg/otlp"
	"context"
	"time"
)

const (
	serviceNamePatient = "PatientUsecase"
	spanNamePatient    = "PatientUsecase"
)

// BookedPatientUseCase -.
type BookedPatientUseCase struct {
	Repo       Patient
	ctxTimeout time.Duration
}

// NewBookedPatient -.
func NewBookedPatient(r Patient, ctxTimeout time.Duration) *BookedPatientUseCase {
	return &BookedPatientUseCase{
		Repo:       r,
		ctxTimeout: ctxTimeout,
	}
}

func (r *BookedPatientUseCase) CreatePatient(ctx context.Context, req *patients.CreatedPatient) (*patients.Patient, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNamePatient, spanNamePatient+"Create")
	defer span.End()

	return r.Repo.CreatePatient(ctx, req)
}

func (r *BookedPatientUseCase) GetPatient(ctx context.Context, req *patients.FieldValueReq) (*patients.Patient, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNamePatient, spanNamePatient+"Get")
	span.End()

	return r.Repo.GetPatient(ctx, req)
}

func (r *BookedPatientUseCase) GetAllPatiens(ctx context.Context, req *patients.GetAllPatients) (*patients.PatientsType, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNamePatient, spanNamePatient+"List")
	span.End()

	return r.Repo.GetAllPatiens(ctx, req)
}

func (r *BookedPatientUseCase) UpdatePatient(ctx context.Context, req *patients.UpdatePatient) (*patients.Patient, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNamePatient, spanNamePatient+"Update")
	span.End()

	return r.Repo.UpdatePatient(ctx, req)
}

func (r *BookedPatientUseCase) UpdatePhonePatient(ctx context.Context, req *patients.UpdatePhoneNumber) (*patients.StatusRes, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNamePatient, spanNamePatient+"UpdatePhone")
	span.End()

	return r.Repo.UpdatePhonePatient(ctx, req)
}

func (r *BookedPatientUseCase) DeletePatient(ctx context.Context, req *patients.FieldValueReq) (*patients.StatusRes, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNamePatient, spanNamePatient+"Delete")
	span.End()

	return r.Repo.DeletePatient(ctx, req)
}
