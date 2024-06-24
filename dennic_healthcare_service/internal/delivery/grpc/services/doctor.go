package services

import (
	pb "Healthcare_Evrone/genproto/healthcare-service"
	"Healthcare_Evrone/internal/entity"
	"Healthcare_Evrone/internal/pkg/minio"
	"Healthcare_Evrone/internal/pkg/otlp"
	"Healthcare_Evrone/internal/usecase"
	"context"
	"time"

	"go.opentelemetry.io/otel/attribute"

	"go.uber.org/zap"
)

type doctorRPC struct {
	logger *zap.Logger
	doctor usecase.DoctorUsecase
}

const (
	serviceNameDoctorDelivery           = "doctorDelivery"
	serviceNameDoctorDeliveryRepoPrefix = "doctorDelivery"
	dayOfWeek                           = "Monday"
)

func DoctorRPC(logget *zap.Logger, doctorUsecase usecase.DoctorUsecase) pb.DoctorServiceServer {
	return &doctorRPC{
		logget,
		doctorUsecase,
	}

}

func (r doctorRPC) CreateDoctor(ctx context.Context, doctor *pb.Doctor) (*pb.Doctor, error) {

	ctx, span := otlp.Start(ctx, serviceNameDoctorDelivery, serviceNameDoctorDeliveryRepoPrefix+"Create")
	span.SetAttributes(attribute.Key("CreateDoctor").String(doctor.Id))
	defer span.End()
	imageUrl := minio.RemoveImageUrl(doctor.ImageUrl)

	req := entity.Doctor{
		Id:            doctor.Id,
		Order:         doctor.Order,
		FirstName:     doctor.FirstName,
		LastName:      doctor.LastName,
		ImageUrl:      imageUrl,
		Gender:        doctor.Gender,
		BirthDate:     doctor.BirthDate,
		PhoneNumber:   doctor.PhoneNumber,
		Email:         doctor.Email,
		Address:       doctor.Address,
		City:          doctor.City,
		Country:       doctor.Country,
		Salary:        doctor.Salary,
		Bio:           doctor.Bio,
		StartWorkDate: doctor.StartWorkDate,
		EndWorkDate:   doctor.EndWorkDate,
		WorkYears:     doctor.WorkYears,
		DepartmentId:  doctor.DepartmentId,
		RoomNumber:    doctor.RoomNumber,
		Password:      doctor.Password,
	}

	resp, err := r.doctor.CreateDoctor(ctx, &req)
	if err != nil {
		r.logger.Error("Failed to create doctor", zap.Error(err))
		return nil, err
	}
	respImageUrl := minio.AddImageUrl(resp.ImageUrl, cfg.MinioService.Bucket.Doctor)

	return &pb.Doctor{
		Id:            resp.Id,
		Order:         resp.Order,
		FirstName:     resp.FirstName,
		LastName:      resp.LastName,
		ImageUrl:      respImageUrl,
		Gender:        resp.Gender,
		BirthDate:     resp.BirthDate,
		PhoneNumber:   resp.PhoneNumber,
		Email:         resp.Email,
		Password:      resp.Password,
		Address:       resp.Address,
		City:          resp.City,
		Country:       resp.Country,
		Salary:        resp.Salary,
		Bio:           resp.Bio,
		StartWorkDate: resp.StartWorkDate,
		EndWorkDate:   resp.EndWorkDate,
		WorkYears:     resp.WorkYears,
		DepartmentId:  resp.DepartmentId,
		RoomNumber:    resp.RoomNumber,
		CreatedAt:     resp.CreatedAt.String(),
		UpdatedAt:     resp.UpdatedAt.String(),
		DeletedAt:     resp.DeletedAt.String(),
	}, nil
}

func (r doctorRPC) GetDoctorById(ctx context.Context, str *pb.GetReqStrDoctor) (*pb.DoctorAndDoctorHours, error) {
	ctx, span := otlp.Start(ctx, serviceNameDoctorDelivery, serviceNameDoctorDeliveryRepoPrefix+"Get")
	span.SetAttributes(attribute.Key("GetDoctorById").String(str.Value))
	defer span.End()
	doctor, err := r.doctor.GetDoctorById(ctx, &entity.GetReqStr{Field: str.Field, Value: str.Value, IsActive: str.IsActive})
	if err != nil {
		r.logger.Error("Failed to get doctor", zap.Error(err))
		return nil, err
	}
	var doctorSpec = []*pb.DoctorSpec{}
	for _, specialization := range doctor.Specializations {
		doctorSpec = append(doctorSpec, &pb.DoctorSpec{
			Id:   specialization.Id,
			Name: specialization.Name,
		})
	}
	imageUrl := minio.AddImageUrl(doctor.ImageUrl, cfg.MinioService.Bucket.Doctor)
	return &pb.DoctorAndDoctorHours{
		Id:              doctor.Id,
		Order:           doctor.Order,
		FirstName:       doctor.FirstName,
		LastName:        doctor.LastName,
		ImageUrl:        imageUrl,
		Gender:          doctor.Gender,
		BirthDate:       doctor.BirthDate,
		PhoneNumber:     doctor.PhoneNumber,
		Email:           doctor.Email,
		Address:         doctor.Address,
		City:            doctor.City,
		Country:         doctor.Country,
		Salary:          doctor.Salary,
		StartTime:       doctor.StartTime,
		FinishTime:      doctor.FinishTime,
		DayOfWeek:       dayOfWeek,
		Bio:             doctor.Bio,
		StartWorkDate:   doctor.StartWorkDate,
		EndWorkDate:     doctor.EndWorkDate,
		WorkYears:       doctor.WorkYears,
		DepartmentId:    doctor.DepartmentId,
		RoomNumber:      doctor.RoomNumber,
		Password:        doctor.Password,
		Specializations: doctorSpec,
		CreatedAt:       doctor.CreatedAt.String(),
		UpdatedAt:       doctor.UpdatedAt.String(),
		DeletedAt:       doctor.DeletedAt.String(),
	}, nil
}

func (r doctorRPC) GetAllDoctors(ctx context.Context, all *pb.GetAllDoctorS) (*pb.ListDoctorsAndHours, error) {

	ctx, span := otlp.Start(ctx, serviceNameDoctorDelivery, serviceNameDoctorDeliveryRepoPrefix+"Get all")
	span.SetAttributes(attribute.Key("GetAllDoctors").String(all.Value))
	defer span.End()
	resp, err := r.doctor.GetAllDoctors(ctx, &entity.GetAll{
		Page:     all.Page,
		Limit:    all.Limit,
		Field:    all.Field,
		Value:    all.Value,
		OrderBy:  all.OrderBy,
		IsActive: all.IsActive,
	})
	if err != nil {
		r.logger.Error("Failed to get all doctors", zap.Error(err))
		return nil, err
	}

	var doctors pb.ListDoctorsAndHours
	for _, doctor := range resp.Doctors {
		imageUrl := minio.AddImageUrl(doctor.ImageUrl, cfg.MinioService.Bucket.Doctor)
		var doctorSpec = []*pb.DoctorSpec{}
		for _, specialization := range doctor.Specializations {
			doctorSpec = append(doctorSpec, &pb.DoctorSpec{
				Id:   specialization.Id,
				Name: specialization.Name,
			})
		}
		doctors.DoctorHours = append(doctors.DoctorHours, &pb.DoctorAndDoctorHours{
			Id:              doctor.Id,
			Order:           doctor.Order,
			FirstName:       doctor.FirstName,
			LastName:        doctor.LastName,
			ImageUrl:        imageUrl,
			Gender:          doctor.Gender,
			BirthDate:       doctor.BirthDate,
			PhoneNumber:     doctor.PhoneNumber,
			Email:           doctor.Email,
			Address:         doctor.Address,
			City:            doctor.City,
			Country:         doctor.Country,
			Salary:          doctor.Salary,
			StartTime:       doctor.StartTime,
			FinishTime:      doctor.FinishTime,
			DayOfWeek:       dayOfWeek,
			Bio:             doctor.Bio,
			StartWorkDate:   doctor.StartWorkDate,
			EndWorkDate:     doctor.EndWorkDate,
			WorkYears:       doctor.WorkYears,
			DepartmentId:    doctor.DepartmentId,
			RoomNumber:      doctor.RoomNumber,
			Password:        doctor.Password,
			Specializations: doctorSpec,
			CreatedAt:       doctor.CreatedAt.String(),
			UpdatedAt:       doctor.UpdatedAt.String(),
			DeletedAt:       doctor.DeletedAt.String(),
		})
	}
	doctors.Count = resp.Count

	return &doctors, nil
}

func (r doctorRPC) UpdateDoctor(ctx context.Context, doctor *pb.Doctor) (*pb.Doctor, error) {

	ctx, span := otlp.Start(ctx, serviceNameDoctorDelivery, serviceNameDoctorDeliveryRepoPrefix+"Update")
	span.SetAttributes(attribute.Key("UpdateDoctor").String(doctor.Id))
	defer span.End()
	imageUrl := minio.RemoveImageUrl(doctor.ImageUrl)
	req := entity.Doctor{
		Id:            doctor.Id,
		Order:         doctor.Order,
		FirstName:     doctor.FirstName,
		LastName:      doctor.LastName,
		ImageUrl:      imageUrl,
		Gender:        doctor.Gender,
		BirthDate:     doctor.BirthDate,
		PhoneNumber:   doctor.PhoneNumber,
		Email:         doctor.Email,
		Address:       doctor.Address,
		City:          doctor.City,
		Country:       doctor.Country,
		Salary:        doctor.Salary,
		Bio:           doctor.Bio,
		StartWorkDate: doctor.StartWorkDate,
		EndWorkDate:   doctor.EndWorkDate,
		WorkYears:     doctor.WorkYears,
		DepartmentId:  doctor.DepartmentId,
		RoomNumber:    doctor.RoomNumber,
		Password:      doctor.Password,
		UpdatedAt:     time.Now(),
	}
	resp, err := r.doctor.UpdateDoctor(ctx, &req)
	if err != nil {
		r.logger.Error("Failed to update doctor", zap.Error(err))
		return nil, err
	}
	respImageUrl := minio.AddImageUrl(resp.ImageUrl, cfg.MinioService.Bucket.Doctor)
	return &pb.Doctor{
		Id:            resp.Id,
		Order:         resp.Order,
		FirstName:     resp.FirstName,
		LastName:      resp.LastName,
		ImageUrl:      respImageUrl,
		Gender:        resp.Gender,
		BirthDate:     resp.BirthDate,
		PhoneNumber:   resp.PhoneNumber,
		Email:         resp.Email,
		Password:      resp.Password,
		Address:       resp.Address,
		City:          resp.City,
		Country:       resp.Country,
		Salary:        resp.Salary,
		Bio:           resp.Bio,
		StartWorkDate: resp.StartWorkDate,
		EndWorkDate:   resp.EndWorkDate,
		WorkYears:     resp.WorkYears,
		DepartmentId:  resp.DepartmentId,
		RoomNumber:    resp.RoomNumber,
		CreatedAt:     resp.CreatedAt.String(),
		UpdatedAt:     resp.UpdatedAt.String(),
		DeletedAt:     resp.DeletedAt.String(),
	}, nil
}

func (r doctorRPC) DeleteDoctor(ctx context.Context, str *pb.GetReqStrDoctor) (*pb.StatusDoctor, error) {

	ctx, span := otlp.Start(ctx, serviceNameDoctorDelivery, serviceNameDoctorDeliveryRepoPrefix+"Delete")
	span.SetAttributes(attribute.Key("DeleteDoctor").String(str.Value))
	defer span.End()

	status, err := r.doctor.DeleteDoctor(ctx, &entity.GetReqStr{Field: str.Field, Value: str.Value, IsActive: str.IsActive})
	if err != nil {
		r.logger.Error("deleted doctor error", zap.Error(err))
		return nil, err
	}
	return &pb.StatusDoctor{Status: status}, nil
}

func (r doctorRPC) ListDoctorsByDepartmentId(ctx context.Context, dep *pb.GetReqStrDep) (*pb.ListDoctors, error) {
	ctx, span := otlp.Start(ctx, serviceNameDoctorDelivery, serviceNameDoctorDeliveryRepoPrefix+"Get all")
	span.SetAttributes(attribute.Key("ListDoctorsByDepartmentId").String(dep.DepartmentId))
	defer span.End()
	resp, err := r.doctor.ListDoctorsByDepartmentId(ctx, &entity.GetReqStrDep{
		DepartmentId: dep.DepartmentId,
		IsActive:     dep.IsActive,
		Page:         dep.Page,
		Limit:        dep.Limit,
		Field:        dep.Field,
		Value:        dep.Value,
		OrderBy:      dep.OrderBy,
	})
	if err != nil {
		r.logger.Error("Failed to get all doctors", zap.Error(err))
		return nil, err
	}

	var doctors pb.ListDoctors
	for _, doctor := range resp.Doctors {
		imageUrl := minio.AddImageUrl(doctor.ImageUrl, cfg.MinioService.Bucket.Doctor)
		var doctorSpec = []*pb.DoctorSpec{}
		for _, specialization := range doctor.Specializations {
			doctorSpec = append(doctorSpec, &pb.DoctorSpec{
				Id:   specialization.Id,
				Name: specialization.Name,
			})
		}
		doctors.Doctors = append(doctors.Doctors, &pb.Doctor{
			Id:              doctor.Id,
			Order:           doctor.Order,
			FirstName:       doctor.FirstName,
			LastName:        doctor.LastName,
			ImageUrl:        imageUrl,
			Gender:          doctor.Gender,
			BirthDate:       doctor.BirthDate,
			PhoneNumber:     doctor.PhoneNumber,
			Email:           doctor.Email,
			Address:         doctor.Address,
			City:            doctor.City,
			Country:         doctor.Country,
			Salary:          doctor.Salary,
			Bio:             doctor.Bio,
			StartWorkDate:   doctor.StartWorkDate,
			EndWorkDate:     doctor.EndWorkDate,
			WorkYears:       doctor.WorkYears,
			DepartmentId:    doctor.DepartmentId,
			RoomNumber:      doctor.RoomNumber,
			Password:        doctor.Password,
			Specializations: doctorSpec,
			CreatedAt:       doctor.CreatedAt.String(),
			UpdatedAt:       doctor.UpdatedAt.String(),
			DeletedAt:       doctor.DeletedAt.String(),
		})
	}
	doctors.Count = resp.Count

	return &doctors, nil
}

func (r doctorRPC) ListDoctorBySpecializationId(ctx context.Context, spec *pb.GetReqStrSpec) (*pb.ListDoctorsAndHours, error) {
	ctx, span := otlp.Start(ctx, serviceNameDoctorDelivery, serviceNameDoctorDeliveryRepoPrefix+"Get all")
	span.SetAttributes(attribute.Key("ListDoctorBySpecializationId").String(spec.SpecializationId))
	defer span.End()
	resp, err := r.doctor.ListDoctorBySpecializationId(ctx, &entity.GetReqStrSpec{
		SpecializationId: spec.SpecializationId,
		IsActive:         spec.IsActive,
		Page:             spec.Page,
		Limit:            spec.Limit,
		Field:            spec.Field,
		Value:            spec.Value,
		OrderBy:          spec.OrderBy,
	})
	if err != nil {
		r.logger.Error("Failed to get all doctors", zap.Error(err))
		return nil, err
	}

	var doctors pb.ListDoctorsAndHours
	for _, doctor := range resp.Doctors {
		imageUrl := minio.AddImageUrl(doctor.ImageUrl, cfg.MinioService.Bucket.Doctor)
		var doctorSpec = []*pb.DoctorSpec{}
		for _, specialization := range doctor.Specializations {
			doctorSpec = append(doctorSpec, &pb.DoctorSpec{
				Id:   specialization.Id,
				Name: specialization.Name,
			})
		}
		doctors.DoctorHours = append(doctors.DoctorHours, &pb.DoctorAndDoctorHours{
			Id:              doctor.Id,
			Order:           doctor.Order,
			FirstName:       doctor.FirstName,
			LastName:        doctor.LastName,
			ImageUrl:        imageUrl,
			Gender:          doctor.Gender,
			BirthDate:       doctor.BirthDate,
			PhoneNumber:     doctor.PhoneNumber,
			Email:           doctor.Email,
			Address:         doctor.Address,
			City:            doctor.City,
			Country:         doctor.Country,
			Salary:          doctor.Salary,
			Bio:             doctor.Bio,
			StartWorkDate:   doctor.StartWorkDate,
			EndWorkDate:     doctor.EndWorkDate,
			WorkYears:       doctor.WorkYears,
			DepartmentId:    doctor.DepartmentId,
			RoomNumber:      doctor.RoomNumber,
			Password:        doctor.Password,
			Specializations: doctorSpec,
			CreatedAt:       doctor.CreatedAt.String(),
			UpdatedAt:       doctor.UpdatedAt.String(),
			DeletedAt:       doctor.DeletedAt.String(),
		})
	}
	doctors.Count += resp.Count
	return &doctors, nil
}
