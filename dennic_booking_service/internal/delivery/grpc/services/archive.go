package services

import (
	pb "booking_service/genproto/booking_service"
	"booking_service/internal/entity/archive"
	"booking_service/internal/pkg/otlp"
	"booking_service/internal/usecase"
	"context"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"

	"time"
)

const (
	serviceNameArchive     = "ArchiveService"
	spanNameArchiveService = "ArchiveService"
)

type BookingArchive struct {
	logger               *zap.Logger
	bookedArchiveUseCase usecase.Archive
}

func BookingArchiveNewRPC(logger *zap.Logger, ArchiveUsaCase usecase.Archive) *BookingArchive {

	return &BookingArchive{
		logger:               logger,
		bookedArchiveUseCase: ArchiveUsaCase,
	}
}

func (r *BookingArchive) CreateArchive(ctx context.Context, req *pb.CreateArchiveReq) (*pb.Archive, error) {
	ctx, span := otlp.Start(ctx, serviceNameArchive, spanNameArchiveService+"Create")
	span.SetAttributes(
		attribute.Key("id").Int64(req.DoctorAvailabilityId),
	)
	defer span.End()

	startTime, err := time.Parse("15:04:05", req.StartTime)
	if err != nil {
		return nil, err
	}

	endTime, err := time.Parse("15:04:05", req.EndTime)
	if err != nil {
		return nil, err
	}

	res, err := r.bookedArchiveUseCase.CreateArchive(ctx, &archive.CreatedArchive{
		DoctorAvailabilityId: req.DoctorAvailabilityId,
		StartTime:            startTime,
		EndTime:              endTime,
		PatientProblem:       req.PatientProblem,
		Status:               req.Status,
		PaymentType:          req.PaymentType,
		PaymentAmount:        float64(req.PaymentAmount),
	})

	if err != nil {
		return nil, err
	}

	return &pb.Archive{
		Id:                   res.Id,
		DoctorAvailabilityId: res.DoctorAvailabilityId,
		StartTime:            res.StartTime.Format("15:04:05"),
		EndTime:              res.EndTime.Format("15:04:05"),
		PatientProblem:       res.PatientProblem,
		Status:               res.Status,
		PaymentType:          res.PaymentType,
		PaymentAmount:        float32(res.PaymentAmount),
		CreatedAt:            res.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:            res.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt:            res.DeletedAt.Format("2006-01-02 15:04:05"),
	}, nil

}

func (r *BookingArchive) GetArchive(ctx context.Context, req *pb.ArchiveFieldValueReq) (*pb.Archive, error) {
	ctx, span := otlp.Start(ctx, serviceNameArchive, spanNameArchiveService+"Get")
	span.SetAttributes(
		attribute.Key(req.Field).String(req.Value),
	)
	defer span.End()

	res, err := r.bookedArchiveUseCase.GetArchive(ctx, &archive.FieldValueReq{
		Field:        req.Field,
		Value:        req.Value,
		DeleteStatus: req.IsActive,
	})
	if err != nil {
		return nil, err
	}

	return &pb.Archive{
		Id:                   res.Id,
		DoctorAvailabilityId: res.DoctorAvailabilityId,
		StartTime:            res.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:              res.EndTime.Format("2006-01-02 15:04:05"),
		PatientProblem:       res.PatientProblem,
		Status:               res.Status,
		PaymentType:          res.PaymentType,
		PaymentAmount:        float32(res.PaymentAmount),
		CreatedAt:            res.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:            res.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt:            res.DeletedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (r *BookingArchive) GetAllArchives(ctx context.Context, req *pb.GetAllArchivesReq) (*pb.Archives, error) {
	ctx, span := otlp.Start(ctx, serviceNameArchive, spanNameArchiveService+"List")
	defer span.End()

	var archives pb.Archives

	archivesRes, err := r.bookedArchiveUseCase.GetAllArchive(ctx, &archive.GetAllArchives{
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

	for _, archiveRes := range archivesRes.Archives {
		var archiv pb.Archive
		archiv.Id = archiveRes.Id
		archiv.DoctorAvailabilityId = archiveRes.DoctorAvailabilityId
		archiv.StartTime = archiveRes.StartTime.Format("2006-01-02 15:04:05")
		archiv.EndTime = archiveRes.EndTime.Format("2006-01-02 15:04:05")
		archiv.PatientProblem = archiveRes.PatientProblem
		archiv.Status = archiveRes.Status
		archiv.PaymentType = archiveRes.PaymentType
		archiv.PaymentAmount = float32(archiveRes.PaymentAmount)
		archiv.CreatedAt = archiveRes.CreatedAt.Format("2006-01-02 15:04:05")
		archiv.UpdatedAt = archiveRes.UpdatedAt.Format("2006-01-02 15:04:05")
		archiv.DeletedAt = archiveRes.DeletedAt.Format("2006-01-02 15:04:05")
		archives.Archives = append(archives.Archives, &archiv)
	}
	archives.Count = archivesRes.Count

	return &archives, nil
}

func (r *BookingArchive) UpdateArchive(ctx context.Context, req *pb.UpdateArchiveReq) (*pb.Archive, error) {
	ctx, span := otlp.Start(ctx, serviceNameArchive, spanNameArchiveService+"Update")
	span.SetAttributes(
		attribute.Key(req.Field).String(req.Value),
	)
	defer span.End()

	startTime, err := time.Parse("15:04:05", req.StartTime)
	if err != nil {
		return nil, err
	}

	endTime, err := time.Parse("15:04:05", req.EndTime)
	if err != nil {
		return nil, err
	}

	res, err := r.bookedArchiveUseCase.UpdateArchive(ctx, &archive.UpdateArchive{
		Field:                req.Field,
		Value:                req.Value,
		DoctorAvailabilityId: req.DoctorAvailabilityId,
		StartTime:            startTime,
		EndTime:              endTime,
		PatientProblem:       req.PatientProblem,
		Status:               req.Status,
		PaymentType:          req.PaymentType,
		PaymentAmount:        float64(req.PaymentAmount),
	})

	if err != nil {
		return nil, err
	}

	return &pb.Archive{
		Id:                   res.Id,
		DoctorAvailabilityId: res.DoctorAvailabilityId,
		StartTime:            res.StartTime.Format("2006-01-02 15:04:05"),
		EndTime:              res.EndTime.Format("2006-01-02 15:04:05"),
		PatientProblem:       res.PatientProblem,
		Status:               res.Status,
		PaymentType:          res.PaymentType,
		PaymentAmount:        float32(res.PaymentAmount),
		CreatedAt:            res.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:            res.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt:            res.DeletedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (r *BookingArchive) DeleteArchive(ctx context.Context, req *pb.ArchiveFieldValueReq) (*pb.DeleteArchiveStatus, error) {
	ctx, span := otlp.Start(ctx, serviceNameArchive, spanNameArchiveService+"Delete")
	span.SetAttributes(
		attribute.Key(req.Field).String(req.Value),
	)
	defer span.End()

	res, err := r.bookedArchiveUseCase.DeleteArchive(ctx, &archive.FieldValueReq{
		Field:        req.Field,
		Value:        req.Value,
		DeleteStatus: req.IsActive,
	})

	if err != nil {
		return &pb.DeleteArchiveStatus{Status: res.Status}, err
	}

	return &pb.DeleteArchiveStatus{Status: res.Status}, nil
}
