package services

import (
	pb "booking_service/genproto/booking_service"
	"booking_service/internal/entity/doctor_notes"
	"booking_service/internal/pkg/otlp"
	"booking_service/internal/usecase"
	"context"

	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

const (
	serviceNameDoctorNotes     = "DoctorNotesService"
	spanNameDoctorNotesService = "DoctorNotesService"
)

type BookingDoctorNotes struct {
	logger                   *zap.Logger
	bookedDoctorNotesUseCase usecase.DoctorNotes
}

func BookingDoctorNotesNewRPC(
	logger *zap.Logger, DoctorNotesUseCase usecase.DoctorNotes) *BookingDoctorNotes {
	return &BookingDoctorNotes{
		logger:                   logger,
		bookedDoctorNotesUseCase: DoctorNotesUseCase,
	}
}

func (r *BookingDoctorNotes) CreateDoctorNote(ctx context.Context, req *pb.CreateDoctorNoteReq) (*pb.DoctorNote, error) {
	ctx, span := otlp.Start(ctx, serviceNameDoctorNotes, spanNameDoctorNotesService+"Create")
	span.SetAttributes(
		attribute.Key("appointment_id").Int64(req.AppointmentId),
	)
	defer span.End()

	res, err := r.bookedDoctorNotesUseCase.CreateDoctorNotes(ctx, &doctor_notes.CreatedDoctorNote{
		AppointmentId: req.AppointmentId,
		DoctorId:      req.DoctorId,
		PatientId:     req.PatientId,
		Prescription:  req.Prescription,
	})

	if err != nil {
		return nil, err
	}

	return &pb.DoctorNote{
		Id:            res.Id,
		AppointmentId: res.AppointmentId,
		DoctorId:      res.DoctorId,
		PatientId:     res.PatientId,
		Prescription:  res.Prescription,
		CreatedAt:     res.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     res.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt:     res.DeletedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (r *BookingDoctorNotes) GetDoctorNote(ctx context.Context, req *pb.FieldValueReq) (*pb.DoctorNote, error) {
	ctx, span := otlp.Start(ctx, serviceNameDoctorNotes, spanNameDoctorNotesService+"Get")
	span.SetAttributes(
		attribute.Key(req.Field).String(req.Value),
	)
	defer span.End()

	res, err := r.bookedDoctorNotesUseCase.GetDoctorNotes(ctx, &doctor_notes.FieldValueReq{
		Field:        req.Field,
		Value:        req.Value,
		DeleteStatus: req.IsActive,
	})

	if err != nil {
		return nil, err
	}

	return &pb.DoctorNote{
		Id:            res.Id,
		AppointmentId: res.AppointmentId,
		DoctorId:      res.DoctorId,
		PatientId:     res.PatientId,
		Prescription:  res.Prescription,
		CreatedAt:     res.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     res.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt:     res.DeletedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (r *BookingDoctorNotes) GetAllNotes(ctx context.Context, req *pb.GetAllReq) (*pb.DoctorNotes, error) {
	ctx, span := otlp.Start(ctx, serviceNameDoctorNotes, spanNameDoctorNotesService+"List")
	defer span.End()

	var notesRes pb.DoctorNotes

	allNotes, err := r.bookedDoctorNotesUseCase.GetAllDoctorNotes(ctx, &doctor_notes.GetAllNotes{
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

	for _, note := range allNotes.DoctorNotes {
		var noteRes pb.DoctorNote
		noteRes.Id = note.Id
		noteRes.AppointmentId = note.AppointmentId
		noteRes.DoctorId = note.DoctorId
		noteRes.PatientId = note.PatientId
		noteRes.Prescription = note.Prescription
		noteRes.CreatedAt = note.CreatedAt.Format("2006-01-02 15:04:05")
		noteRes.UpdatedAt = note.UpdatedAt.Format("2006-01-02 15:04:05")
		noteRes.DeletedAt = note.DeletedAt.Format("2006-01-02 15:04:05")
		notesRes.DoctorNotes = append(notesRes.DoctorNotes, &noteRes)
	}
	notesRes.Count = allNotes.Count

	return &notesRes, nil
}

func (r *BookingDoctorNotes) UpdateDoctorNote(ctx context.Context, req *pb.UpdateDoctorNoteReq) (*pb.DoctorNote, error) {
	ctx, span := otlp.Start(ctx, serviceNameDoctorNotes, spanNameDoctorNotesService+"Update")
	span.SetAttributes(
		attribute.Key(req.Field).String(req.Value),
	)
	defer span.End()

	res, err := r.bookedDoctorNotesUseCase.UpdateDoctorNotes(ctx, &doctor_notes.UpdateDoctorNoteReq{
		Field:         req.Field,
		Value:         req.Value,
		AppointmentId: req.AppointmentId,
		DoctorId:      req.DoctorId,
		PatientId:     req.PatientId,
		Prescription:  req.Prescription,
	})

	if err != nil {
		return nil, err
	}

	return &pb.DoctorNote{
		Id:            res.Id,
		AppointmentId: res.AppointmentId,
		DoctorId:      res.DoctorId,
		PatientId:     res.PatientId,
		Prescription:  res.Prescription,
		CreatedAt:     res.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     res.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt:     res.DeletedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (r *BookingDoctorNotes) DeleteDoctorNote(ctx context.Context, req *pb.FieldValueReq) (*pb.StatusRes, error) {
	ctx, span := otlp.Start(ctx, serviceNameDoctorNotes, spanNameDoctorNotesService+"Delete")
	span.SetAttributes(
		attribute.Key(req.Field).String(req.Value),
	)
	defer span.End()

	res, err := r.bookedDoctorNotesUseCase.DeleteDoctorNotes(ctx, &doctor_notes.FieldValueReq{
		Field:        req.Field,
		Value:        req.Value,
		DeleteStatus: req.IsActive,
	})

	if err != nil {
		return &pb.StatusRes{Status: res.Status}, err
	}

	return &pb.StatusRes{Status: res.Status}, err
}
