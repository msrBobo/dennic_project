package services

import (
	pb "booking_service/genproto/booking_service"
	"booking_service/internal/entity/patients"
	"booking_service/internal/pkg/otlp"
	"booking_service/internal/usecase"
	"context"
	"github.com/rickb777/date"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
	//"google.golang.org/grpc/codes"
)

const (
	serviceNamePatient     = "patientsService"
	sapmNamePatientService = "patientsService"
)

type BookingPatient struct {
	logger               *zap.Logger
	bookedPatientUseCase usecase.Patient
}

func BookingPatientNewRPC(logger *zap.Logger, PatientUsaCase usecase.Patient) *BookingPatient {
	return &BookingPatient{
		logger:               logger,
		bookedPatientUseCase: PatientUsaCase,
	}
}

func (r *BookingPatient) CreatePatient(ctx context.Context, req *pb.CreatePatientReq) (*pb.Patient, error) {
	ctx, span := otlp.Start(ctx, serviceNamePatient, sapmNamePatientService+"Create")
	span.SetAttributes(
		attribute.Key("id").String(req.Id),
	)
	defer span.End()
	reqDate, err := date.AutoParse(req.BirthDate)
	if err != nil {
		return nil, err
	}

	res, err := r.bookedPatientUseCase.CreatePatient(ctx, &patients.CreatedPatient{
		Id:             req.Id,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		BirthDate:      reqDate,
		Gender:         req.Gender,
		BloodGroup:     req.BloodGroup,
		PhoneNumber:    req.PhoneNumber,
		City:           req.City,
		Country:        req.Country,
		Address:        req.Address,
		PatientProblem: req.PatientProblem,
	})

	if err != nil {
		return nil, err
	}

	return &pb.Patient{
		Id:             res.Id,
		FirstName:      res.FirstName,
		LastName:       res.LastName,
		BirthDate:      res.BirthDate.String(),
		Gender:         res.Gender,
		BloodGroup:     res.BloodGroup,
		PhoneNumber:    res.PhoneNumber,
		City:           res.City,
		Country:        res.Country,
		Address:        res.Address,
		PatientProblem: res.PatientProblem,
		CreatedAt:      res.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      res.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt:      res.DeletedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (r *BookingPatient) GetPatient(ctx context.Context, req *pb.PatientFieldValueReq) (*pb.Patient, error) {
	ctx, span := otlp.Start(ctx, serviceNamePatient, sapmNamePatientService+"Get")
	span.SetAttributes(
		attribute.Key(req.Field).String(req.Value),
	)
	defer span.End()
	res, err := r.bookedPatientUseCase.GetPatient(ctx, &patients.FieldValueReq{
		Field:        req.Field,
		Value:        req.Value,
		DeleteStatus: req.IsActive,
	})
	if err != nil {
		return nil, err
	}

	return &pb.Patient{
		Id:             res.Id,
		FirstName:      res.FirstName,
		LastName:       res.LastName,
		BirthDate:      res.BirthDate.String(),
		Gender:         res.Gender,
		BloodGroup:     res.BloodGroup,
		PhoneNumber:    res.PhoneNumber,
		City:           res.City,
		Country:        res.Country,
		Address:        res.Address,
		PatientProblem: res.PatientProblem,
		CreatedAt:      res.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      res.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt:      res.DeletedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (r *BookingPatient) GetAllPatients(ctx context.Context, req *pb.GetAllPatientsReq) (*pb.Patients, error) {
	ctx, span := otlp.Start(ctx, serviceNamePatient, sapmNamePatientService+"List")
	defer span.End()

	var patentsRes pb.Patients
	allPatients, err := r.bookedPatientUseCase.GetAllPatiens(ctx, &patients.GetAllPatients{
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

	for _, patient := range allPatients.Patients {
		var patientRes pb.Patient
		patientRes.Id = patient.Id
		patientRes.FirstName = patient.FirstName
		patientRes.LastName = patient.LastName
		patientRes.BirthDate = patient.BirthDate.String()
		patientRes.Gender = patient.Gender
		patientRes.BloodGroup = patient.BloodGroup
		patientRes.PhoneNumber = patient.PhoneNumber
		patientRes.City = patient.City
		patientRes.Country = patient.Country
		patientRes.Address = patient.Address
		patientRes.PatientProblem = patient.PatientProblem
		patientRes.CreatedAt = patient.CreatedAt.Format("2006-01-02 15:04:05")
		patientRes.UpdatedAt = patient.UpdatedAt.Format("2006-01-02 15:04:05")
		patientRes.DeletedAt = patient.DeletedAt.Format("2006-01-02 15:04:05")
		patentsRes.Patients = append(patentsRes.Patients, &patientRes)
	}
	patentsRes.Count = allPatients.Count

	return &patentsRes, nil
}

func (r *BookingPatient) UpdatePatient(ctx context.Context, req *pb.UpdatePatientReq) (*pb.Patient, error) {
	ctx, span := otlp.Start(ctx, serviceNamePatient, sapmNamePatientService+"Update")
	span.SetAttributes(
		attribute.Key(req.Field).String(req.Value),
	)
	defer span.End()
	reqData, err := date.AutoParse(req.BirthDate)
	if err != nil {
		return nil, err
	}

	res, err := r.bookedPatientUseCase.UpdatePatient(ctx, &patients.UpdatePatient{
		Field:          req.Field,
		Value:          req.Value,
		FirstName:      req.FirstName,
		LastName:       req.LastName,
		BirthDate:      reqData,
		Gender:         req.Gender,
		BloodGroup:     req.BloodGroup,
		City:           req.City,
		Country:        req.Country,
		Address:        req.Address,
		PatientProblem: req.PatientProblem,
	})
	if err != nil {
		return nil, err
	}

	return &pb.Patient{
		Id:             res.Id,
		FirstName:      res.FirstName,
		LastName:       res.LastName,
		BirthDate:      res.BirthDate.String(),
		Gender:         res.Gender,
		BloodGroup:     res.BloodGroup,
		PhoneNumber:    res.PhoneNumber,
		City:           res.City,
		Country:        res.Country,
		Address:        res.Address,
		PatientProblem: res.PatientProblem,
		CreatedAt:      res.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:      res.UpdatedAt.Format("2006-01-02 15:04:05"),
		DeletedAt:      res.DeletedAt.Format("2006-01-02 15:04:05"),
	}, nil
}

func (r *BookingPatient) UpdatePhonePatient(ctx context.Context, req *pb.UpdatePhoneNumber) (*pb.PatientStatus, error) {
	ctx, span := otlp.Start(ctx, serviceNamePatient, sapmNamePatientService+"UpdatePhone")
	span.SetAttributes(
		attribute.Key(req.Field).String(req.Value),
	)
	defer span.End()
	res, err := r.bookedPatientUseCase.UpdatePhonePatient(ctx, &patients.UpdatePhoneNumber{
		Field:       req.Field,
		Value:       req.Value,
		PhoneNumber: req.PhoneNumber,
	})

	if err != nil {
		return &pb.PatientStatus{Status: res.Status}, err
	}

	return &pb.PatientStatus{Status: res.Status}, nil
}

func (r *BookingPatient) DeletePatient(ctx context.Context, req *pb.PatientFieldValueReq) (*pb.PatientStatus, error) {
	ctx, span := otlp.Start(ctx, serviceNamePatient, sapmNamePatientService+"Delete")
	span.SetAttributes(
		attribute.Key(req.Field).String(req.Value),
	)
	defer span.End()
	res, err := r.bookedPatientUseCase.DeletePatient(ctx, &patients.FieldValueReq{
		Field:        req.Field,
		Value:        req.Value,
		DeleteStatus: req.IsActive,
	})

	if err != nil {
		return &pb.PatientStatus{Status: res.Status}, err
	}

	return &pb.PatientStatus{Status: res.Status}, nil
}
