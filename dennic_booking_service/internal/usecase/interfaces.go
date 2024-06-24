// Package usecase implements application business logic. Each logic group in own file.
package usecase

import (
	"booking_service/internal/entity/archive"
	appointment "booking_service/internal/entity/booked_appointments"
	"booking_service/internal/entity/doctor_availability"
	"booking_service/internal/entity/doctor_notes"
	"booking_service/internal/entity/patients"
	"context"
)

//go:generate mockgen -source=interfaces.go -destination=./mocks_test.go -package=usecase_test

type (
	//datient -.
	Patient interface {
		CreatePatient(ctx context.Context, req *patients.CreatedPatient) (*patients.Patient, error)
		GetPatient(ctx context.Context, req *patients.FieldValueReq) (*patients.Patient, error)
		GetAllPatiens(ctx context.Context, req *patients.GetAllPatients) (*patients.PatientsType, error)
		UpdatePatient(ctx context.Context, req *patients.UpdatePatient) (*patients.Patient, error)
		UpdatePhonePatient(ctx context.Context, req *patients.UpdatePhoneNumber) (*patients.StatusRes, error)
		DeletePatient(ctx context.Context, req *patients.FieldValueReq) (*patients.StatusRes, error)
	}

	// BookedAppointments -.
	BookedAppointments interface {
		CreateAppointment(ctx context.Context, req *appointment.CreateAppointment) (*appointment.Appointment, error)
		GetAppointment(ctx context.Context, req *appointment.FieldValueReq) (*appointment.Appointment, error)
		GetAllAppointment(ctx context.Context, req *appointment.GetAllAppointment) (*appointment.AppointmentsType, error)
		GetFilteredAppointments(ctx context.Context, req *appointment.GetFilteredRequest) (*appointment.AppointmentsType, error)
		UpdateAppointment(ctx context.Context, req *appointment.UpdateAppointment) (*appointment.Appointment, error)
		DeleteAppointment(ctx context.Context, req *appointment.FieldValueReq) (*appointment.StatusRes, error)
	}
	// DoctorNotes -.
	DoctorNotes interface {
		CreateDoctorNotes(ctx context.Context, req *doctor_notes.CreatedDoctorNote) (*doctor_notes.DoctorNote, error)
		GetDoctorNotes(ctx context.Context, req *doctor_notes.FieldValueReq) (*doctor_notes.DoctorNote, error)
		GetAllDoctorNotes(ctx context.Context, req *doctor_notes.GetAllNotes) (*doctor_notes.DoctorNotesType, error)
		UpdateDoctorNotes(ctx context.Context, req *doctor_notes.UpdateDoctorNoteReq) (*doctor_notes.DoctorNote, error)
		DeleteDoctorNotes(ctx context.Context, req *doctor_notes.FieldValueReq) (*doctor_notes.StatusRes, error)
	}

	// Archive -.
	Archive interface {
		CreateArchive(ctx context.Context, req *archive.CreatedArchive) (*archive.Archive, error)
		GetArchive(ctx context.Context, req *archive.FieldValueReq) (*archive.Archive, error)
		GetAllArchive(ctx context.Context, req *archive.GetAllArchives) (*archive.ArchivesType, error)
		UpdateArchive(ctx context.Context, req *archive.UpdateArchive) (*archive.Archive, error)
		DeleteArchive(ctx context.Context, req *archive.FieldValueReq) (*archive.StatusRes, error)
	}

	// DoctorAvailability -.
	DoctorAvailability interface {
		CreateDoctorAvailability(ctx context.Context, req *doctor_availability.CreateDoctorAvailability) (*doctor_availability.DoctorAvailability, error)
		GetDoctorAvailability(ctx context.Context, req *doctor_availability.FieldValueReq) (*doctor_availability.DoctorAvailability, error)
		GetAllDoctorAvailability(ctx context.Context, req *doctor_availability.GetAllReq) (*doctor_availability.DoctorAvailabilityType, error)
		UpdateDoctorAvailability(ctx context.Context, req *doctor_availability.UpdateDoctorAvailability) (*doctor_availability.DoctorAvailability, error)
		DeleteDoctorAvailability(ctx context.Context, req *doctor_availability.FieldValueReq) (*doctor_availability.StatusRes, error)
		GeDoctorTimeStatus(ctx context.Context, req *doctor_availability.DoctorTimeStatusReq) (*doctor_availability.TimeStatus, error)
		UpdateTimeStatus(ctx context.Context, req *doctor_availability.DoctorTimeUpdateStatusReq) (*doctor_availability.TimeStatus, error)
		GetDoctorTimeByDoctorID(ctx context.Context, req *doctor_availability.DoctorTimeFieldValueByIdReq) (*doctor_availability.Hours, error)
	}
)
