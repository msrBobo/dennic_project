package model_booking_service

type Appointment struct {
	Id                    int64   `json:"id"`
	DepartmentId          string  `json:"department_id"`
	DoctorId              string  `json:"doctor_id"`
	DoctorFirstName       string  `json:"doctor_first_name"`
	DoctorLastName        string  `json:"doctor_last_name"`
	DoctorPatientsCount   int64   `json:"doctor_patients_count"`
	DoctorWorkingYears    int64   `json:"doctor_working_years"`
	PatientId             string  `json:"patient_id"`
	PatientFullName       string  `json:"patient_full_name"`
	PatientPhoneNumber    string  `json:"patient_phone_number"`
	AppointmentDate       string  `json:"appointment_date"`
	AppointmentStartTime  string  `json:"appointment_start_time"`
	AppointmentFinishTime string  `json:"appointment_finish_time"`
	Duration              int64   `json:"duration"`
	Key                   string  `json:"key"`
	ExpiresAt             string  `json:"expires_at"`
	PatientStatus         string  `json:"patient_status"`
	PatientProblem        string  `json:"patient_problem"`
	DoctorServiceId       string  `json:"doctor_service_id"`
	PaymentType           string  `json:"payment_type"`
	PaymentAmount         float64 `json:"payment_amount"`
	UserId                string  `json:"user_id"`
	CreatedAt             string  `json:"created_at"`
	UpdatedAt             string  `json:"updated_at"`
}

type AppointmentsType struct {
	Count        int64          `json:"count"`
	Appointments []*Appointment `json:"appointments"`
}

type CreateAppointmentReq struct {
	DepartmentId    string  `json:"-"`
	DoctorId        string  `json:"doctor_id"`
	PatientId       string  `json:"patient_id"`
	AppointmentDate string  `json:"appointment_date"`
	AppointmentTime string  `json:"appointment_time"`
	Duration        int64   `json:"-"`
	Key             string  `json:"-"`
	ExpiresAt       string  `json:"-"`
	PatientStatus   string  `json:"-"`
	PatientProblem  string  `json:"-"`
	DoctorServiceId string  `json:"doctor_service_id"`
	PaymentType     string  `json:"-"`
	PaymentAmount   float64 `json:"-"`
}

type UpdateAppointmentReq struct {
	BookedAppointmentId string  `json:"booked_appointment_id"`
	AppointmentDate     string  `json:"appointment_date"`
	AppointmentTime     string  `json:"appointment_time"`
	Duration            int64   `json:"duration"`
	Key                 string  `json:"key"`
	ExpiresAt           string  `json:"expires_at"`
	PatientStatus       string  `json:"patient_status"`
	PatientProblem      string  `json:"patient_problem"`
	DoctorServiceId     string  `json:"doctor_service_id"`
	PaymentType         string  `json:"payment_type"`
	PaymentAmount       float64 `json:"payment_amount"`
	UserId              string  `json:"user_id"`
}

type BookingReq struct {
	DoctorID string `json:"doctor_id"`
	Date     string `json:"date"`
}

type BookingResp struct {
	Time      string `json:"time"`
	Status    bool   `json:"status"`
	TimeOfDay bool   `json:"time_of_day"`
}
