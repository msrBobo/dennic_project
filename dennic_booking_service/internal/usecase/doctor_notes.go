package usecase

import (
	"booking_service/internal/entity/doctor_notes"
	"booking_service/internal/pkg/otlp"
	"context"
	"time"
)

const (
	serviceNameDoctorNotes = "DoctorNotesService"
	spanNameDoctorNotes    = "DoctorNotesUsecase"
)

// BookedDoctorNotesUseCase -.
type BookedDoctorNotesUseCase struct {
	Repo       DoctorNotes
	ctxTimeout time.Duration
}

// NewBookedDoctorNotes -.
func NewBookedDoctorNotes(r DoctorNotes, ctxTimeout time.Duration) *BookedDoctorNotesUseCase {
	return &BookedDoctorNotesUseCase{
		Repo:       r,
		ctxTimeout: ctxTimeout,
	}
}

func (r *BookedDoctorNotesUseCase) CreateDoctorNotes(ctx context.Context, req *doctor_notes.CreatedDoctorNote) (*doctor_notes.DoctorNote, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameDoctorNotes, spanNameDoctorNotes+"Create")
	span.End()

	return r.Repo.CreateDoctorNotes(ctx, req)
}

func (r *BookedDoctorNotesUseCase) GetDoctorNotes(ctx context.Context, req *doctor_notes.FieldValueReq) (*doctor_notes.DoctorNote, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameDoctorNotes, spanNameDoctorNotes+"Get")
	span.End()

	return r.Repo.GetDoctorNotes(ctx, req)
}

func (r *BookedDoctorNotesUseCase) GetAllDoctorNotes(ctx context.Context, req *doctor_notes.GetAllNotes) (*doctor_notes.DoctorNotesType, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameDoctorNotes, spanNameDoctorNotes+"List")
	span.End()

	return r.Repo.GetAllDoctorNotes(ctx, req)
}

func (r *BookedDoctorNotesUseCase) UpdateDoctorNotes(ctx context.Context, req *doctor_notes.UpdateDoctorNoteReq) (*doctor_notes.DoctorNote, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameDoctorNotes, spanNameDoctorNotes+"Update")
	span.End()

	return r.Repo.UpdateDoctorNotes(ctx, req)
}

func (r *BookedDoctorNotesUseCase) DeleteDoctorNotes(ctx context.Context, req *doctor_notes.FieldValueReq) (*doctor_notes.StatusRes, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameDoctorNotes, spanNameDoctorNotes+"Delete")
	span.End()

	return r.Repo.DeleteDoctorNotes(ctx, req)
}
