package repo

import (
	"booking_service/internal/infrastructure/repository"
	"booking_service/internal/pkg/postgres"
)

type IBookingStorage interface {
	Archive() repository.Archive
	BookedAppointments() repository.BookedAppointments
	DoctorAvailability() repository.DoctorAvailability
	DoctorNotes() repository.DoctorNotes
	Patients() repository.Patient
}

type BookingStoragePg struct {
	archive            repository.Archive
	bookedAppointments repository.BookedAppointments
	doctorAvailability repository.DoctorAvailability
	doctorNotes        repository.DoctorNotes
	patients           repository.Patient
}

func NewStorage(db *postgres.PostgresDB) IBookingStorage {
	return &BookingStoragePg{
		archive:            NewBookingArchive(db),
		bookedAppointments: NewBookingAppointment(db),
		doctorAvailability: NewDoctorAvailability(db),
		doctorNotes:        NewDoctorNotes(db),
		patients:           NewBookingPatients(db),
	}
}

func (s *BookingStoragePg) Archive() repository.Archive {
	return s.archive
}

func (s *BookingStoragePg) BookedAppointments() repository.BookedAppointments {
	return s.bookedAppointments
}

func (s *BookingStoragePg) DoctorAvailability() repository.DoctorAvailability {
	return s.doctorAvailability
}

func (s *BookingStoragePg) DoctorNotes() repository.DoctorNotes {
	return s.doctorNotes
}

func (s *BookingStoragePg) Patients() repository.Patient {
	return s.patients
}
