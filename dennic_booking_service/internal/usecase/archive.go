package usecase

import (
	"booking_service/internal/entity/archive"
	"booking_service/internal/pkg/otlp"
	"context"
	"time"
)

const (
	serviceNameArchive = "ArchiveService"
	spanNameArchive    = "ArchiveUsecase"
)

// BookedArchiveUseCase -.
type BookedArchiveUseCase struct {
	Repo       Archive
	ctxTimeout time.Duration
}

// NewBookedArchive -.
func NewBookedArchive(r Archive, ctxTimeout time.Duration) *BookedArchiveUseCase {
	return &BookedArchiveUseCase{
		Repo:       r,
		ctxTimeout: ctxTimeout,
	}
}

func (r *BookedArchiveUseCase) CreateArchive(ctx context.Context, req *archive.CreatedArchive) (*archive.Archive, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameArchive, spanNameArchive+"Create")
	defer span.End()

	return r.Repo.CreateArchive(ctx, req)
}

func (r *BookedArchiveUseCase) GetArchive(ctx context.Context, req *archive.FieldValueReq) (*archive.Archive, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameArchive, spanNameArchive+"Get")
	defer span.End()

	return r.Repo.GetArchive(ctx, req)
}

func (r *BookedArchiveUseCase) GetAllArchive(ctx context.Context, req *archive.GetAllArchives) (*archive.ArchivesType, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameArchive, spanNameArchive+"List")
	defer span.End()

	return r.Repo.GetAllArchive(ctx, req)
}

func (r *BookedArchiveUseCase) UpdateArchive(ctx context.Context, req *archive.UpdateArchive) (*archive.Archive, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameArchive, spanNameArchive+"Update")
	defer span.End()

	return r.Repo.UpdateArchive(ctx, req)
}

func (r *BookedArchiveUseCase) DeleteArchive(ctx context.Context, req *archive.FieldValueReq) (*archive.StatusRes, error) {
	ctx, cancel := context.WithTimeout(ctx, r.ctxTimeout)
	defer cancel()

	ctx, span := otlp.Start(ctx, serviceNameArchive, spanNameArchive+"Delete")
	defer span.End()

	return r.Repo.DeleteArchive(ctx, req)
}
