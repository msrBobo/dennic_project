package v1

import (
	"context"
	e "dennic_api_gateway/api/handlers/regtool"
	"dennic_api_gateway/api/models"
	"dennic_api_gateway/api/models/model_booking_service"
	"dennic_api_gateway/genproto/booking_service"
	pb "dennic_api_gateway/genproto/booking_service"
	"dennic_api_gateway/genproto/healthcare-service"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"google.golang.org/protobuf/encoding/protojson"
)

// CreateBookedAppointment ...
// @Summary CreateBookedAppointment
// @Description CreateBookedAppointment - Api for create booked appointment
// @Tags Appointment
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param CreateAppointmentReq body model_booking_service.CreateAppointmentReq true "CreateAppointmentReq"
// @Success 200 {object} model_booking_service.Appointment
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/appointment/create [post]
func (h *HandlerV1) CreateBookedAppointment(c *gin.Context) {
	var (
		body        model_booking_service.CreateAppointmentReq
		jsonMarshal protojson.MarshalOptions
	)
	jsonMarshal.UseProtoNames = true

	key := e.GenerateRandomNumber(8)

	userInfo, err := e.GetUserInfo(c)

	if e.HandleError(c, err, h.log, http.StatusUnauthorized, "missing token in the header") {
		return
	}

	err = c.ShouldBindJSON(&body)

	if e.HandleError(c, err, h.log, http.StatusBadRequest, "CreateBookedAppointment") {
		return
	}

	if body.PatientId == "" {
		c.JSON(http.StatusOK, gin.H{
			"message": "BRINCHI IMKONIYAT KETTI",
		})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	doctor, err := h.serviceManager.HealthcareService().DoctorService().GetDoctorById(ctx, &healthcare.GetReqStrDoctor{
		Field:    "id",
		Value:    body.DoctorId,
		IsActive: false,
	})

	if e.HandleError(c, err, h.log, http.StatusBadRequest, "CreateBookedAppointment") {
		return
	}

	patient, err := h.serviceManager.BookingService().PatientService().GetPatient(ctx, &pb.PatientFieldValueReq{
		Field:    "id",
		Value:    body.PatientId,
		IsActive: false,
	})

	if e.HandleError(c, err, h.log, http.StatusBadRequest, "CreateBookedAppointment") {
		return
	}

	service, err := h.serviceManager.HealthcareService().DoctorsService().GetDoctorServiceByID(ctx, &healthcare.GetReqStr{
		Field:    "id",
		Value:    body.DoctorServiceId,
		IsActive: false,
	})

	ExpiresAt, err := time.Parse("15:04:05", body.AppointmentTime)
	if e.HandleError(c, err, h.log, http.StatusBadRequest, "CreateBookedAppointment") {
		return
	}

	ExpiresAt = ExpiresAt.Add(15 * time.Minute)

	res, err := h.serviceManager.BookingService().BookedAppointment().CreateAppointment(ctx, &pb.CreateAppointmentReq{
		DepartmentId:    doctor.DepartmentId,
		DoctorId:        body.DoctorId,
		PatientId:       body.PatientId,
		DoctorServiceId: body.DoctorServiceId,
		PatientProblem:  patient.PatientProblem,
		PaymentType:     "cash",
		PaymentAmount:   service.OfflinePrice,
		AppointmentDate: body.AppointmentDate,
		AppointmentTime: body.AppointmentTime,
		Duration:        body.Duration,
		Key:             key,
		ExpiresAt:       ExpiresAt.Format("15:04:05"),
		Status:          "waiting",
		UserId:          userInfo.UserId,
	})

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, "CreateBookedAppointment") {
		return
	}

	_, err = h.serviceManager.BookingService().DoctorTimes().UpdateTimeStatus(ctx, &pb.DoctorTimeUpdateStatusReq{
		DepartmentId: doctor.DepartmentId,
		DoctorId:     body.DoctorId,
		DoctorDate:   body.AppointmentDate,
		StartTime:    body.AppointmentTime,
		Status:       "unavailable",
	})

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, "CreateBookedAppointment") {
		return
	}

	c.JSON(http.StatusOK, model_booking_service.Appointment{
		Id:                   res.Id,
		DepartmentId:         res.DepartmentId,
		DoctorId:             res.DoctorId,
		PatientId:            res.PatientId,
		AppointmentDate:      res.AppointmentDate,
		AppointmentStartTime: res.AppointmentTime,
		Duration:             res.Duration,
		Key:                  res.Key,
		ExpiresAt:            res.ExpiresAt,
		PatientStatus:        res.Status,
		PatientProblem:       res.PatientProblem,
		DoctorServiceId:      res.DoctorServiceId,
		PaymentType:          res.PaymentType,
		PaymentAmount:        float64(res.PaymentAmount),
		UserId:               res.UserId,
		CreatedAt:            res.CreatedAt,
		UpdatedAt:            e.UpdateTimeFilter(res.UpdatedAt),
	})
}

// GetBookedAppointment ...
// @Summary GetBookedAppointment
// @Description GetBookedAppointment - API to get Booked appointment by ID
// @Tags Appointment
// @Accept json
// @Produce json
// @Param id query string true "id"
// @Success 200 {object} model_booking_service.Appointment
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/appointment/get [get]
func (h *HandlerV1) GetBookedAppointment(c *gin.Context) {
	id := c.Query("id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	res, err := h.serviceManager.BookingService().BookedAppointment().GetAppointment(ctx, &pb.AppointmentFieldValueReq{
		Field:    "id",
		Value:    id,
		IsActive: false,
	})

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, "GetBookedAppointment") {
		return
	}

	appointment := model_booking_service.Appointment{
		Id:                   res.Id,
		DepartmentId:         res.DepartmentId,
		DoctorId:             res.DoctorId,
		PatientId:            res.PatientId,
		AppointmentDate:      res.AppointmentDate,
		AppointmentStartTime: res.AppointmentTime,
		Duration:             res.Duration,
		Key:                  res.Key,
		ExpiresAt:            res.ExpiresAt,
		PatientStatus:        res.Status,
		PatientProblem:       res.PatientProblem,
		DoctorServiceId:      res.DoctorServiceId,
		PaymentType:          res.PaymentType,
		PaymentAmount:        float64(res.PaymentAmount),
		UserId:               res.UserId,
		CreatedAt:            res.CreatedAt,
		UpdatedAt:            e.UpdateTimeFilter(res.UpdatedAt),
	}

	doctor, err := h.serviceManager.HealthcareService().DoctorService().GetDoctorById(ctx, &healthcare.GetReqStrDoctor{
		Field:    "id",
		Value:    appointment.DoctorId,
		IsActive: false,
	})
	if e.HandleError(c, err, h.log, http.StatusInternalServerError, "GetDoctorById") {
		return
	}

	appointment.DoctorFirstName = doctor.FirstName
	appointment.DoctorLastName = doctor.LastName
	appointment.DoctorWorkingYears = int64(doctor.WorkYears)

	patient, err := h.serviceManager.BookingService().PatientService().GetPatient(ctx, &pb.PatientFieldValueReq{
		Field:    "id",
		Value:    res.PatientId,
		IsActive: false,
	})
	if e.HandleError(c, err, h.log, http.StatusInternalServerError, "GetPatient") {
		return
	}

	appointment.PatientFullName = patient.FirstName + " " + patient.LastName
	appointment.PatientPhoneNumber = patient.PhoneNumber

	appointments, err := h.serviceManager.BookingService().BookedAppointment().GetFilteredAppointments(ctx, &pb.GetFilteredRequest{
		Field:    "doctor_id",
		Value:    doctor.Id,
		IsActive: false,
		Status:   "attended",
	})

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, "GetBookedAppointment -> GetFilteredAppointments") {
		return
	}

	appointment.DoctorPatientsCount = appointments.Count

	startTime, err := time.Parse("15:04:05", appointment.AppointmentStartTime)
	if e.HandleError(c, err, h.log, http.StatusInternalServerError, "parsing appointment start time back to time") {
		return
	}
	duration := time.Duration(appointment.Duration) * time.Minute
	finishTime := startTime.Add(duration)
	appointment.AppointmentFinishTime = finishTime.Format("15:04:05")

	c.JSON(http.StatusOK, appointment)
}

// ListBookedAppointments ...
// @Summary ListBookedAppointments
// @Description ListBookedAppointments - API to list doctor notes
// @Tags Appointment
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Success 200 {object} model_booking_service.AppointmentsType
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/appointment [get]
func (h *HandlerV1) ListBookedAppointments(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	userInfo, err := e.GetUserInfo(c)

	if e.HandleError(c, err, h.log, http.StatusUnauthorized, "missing token in the header") {
		return
	}

	res, err := h.serviceManager.BookingService().BookedAppointment().GetFilteredAppointments(ctx, &pb.GetFilteredRequest{
		Field:    "user_id",
		Value:    userInfo.UserId,
		IsActive: false,
	})

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, "ListBookedAppointments") {
		return
	}

	var response model_booking_service.AppointmentsType
	for _, appointment := range res.Appointments {
		doctor, err := h.serviceManager.HealthcareService().DoctorService().GetDoctorById(ctx, &healthcare.GetReqStrDoctor{
			Field:    "id",
			Value:    appointment.DoctorId,
			IsActive: false,
		})
		if e.HandleError(c, err, h.log, http.StatusInternalServerError, "GetDoctorById") {
			return
		}

		patient, err := h.serviceManager.BookingService().PatientService().GetPatient(ctx, &booking_service.PatientFieldValueReq{
			Field:    "id",
			Value:    appointment.PatientId,
			IsActive: false,
		})
		if e.HandleError(c, err, h.log, http.StatusInternalServerError, "GetPatient") {
			return
		}

		var app model_booking_service.Appointment
		app.Id = appointment.Id
		app.DepartmentId = appointment.DepartmentId
		app.DoctorId = appointment.DoctorId
		app.PatientId = appointment.PatientId
		app.AppointmentDate = appointment.AppointmentDate
		app.AppointmentStartTime = appointment.AppointmentTime
		app.Duration = appointment.Duration
		app.Key = appointment.Key
		app.UserId = appointment.UserId
		app.ExpiresAt = appointment.ExpiresAt
		app.PatientStatus = appointment.Status
		app.PatientProblem = appointment.PatientProblem
		app.PaymentType = appointment.PaymentType
		app.PaymentAmount = float64(appointment.PaymentAmount)
		app.DoctorServiceId = appointment.DoctorServiceId
		app.CreatedAt = appointment.CreatedAt
		app.UpdatedAt = e.UpdateTimeFilter(appointment.UpdatedAt)
		app.DoctorFirstName = doctor.FirstName
		app.DoctorLastName = doctor.LastName
		app.PatientFullName = patient.FirstName + " " + patient.LastName
		app.PatientPhoneNumber = patient.PhoneNumber
		app.PatientBirthDate = patient.BirthDate

		startTime, err := time.Parse("15:04:05", app.AppointmentStartTime)
		if e.HandleError(c, err, h.log, http.StatusInternalServerError, "parsing appointment start time back to time") {
			return
		}
		duration := time.Duration(app.Duration) * time.Minute
		finishTime := startTime.Add(duration)
		app.AppointmentFinishTime = finishTime.Format("15:04:05")

		response.Appointments = append(response.Appointments, &app)
	}

	c.JSON(http.StatusOK, &model_booking_service.AppointmentsType{
		Appointments: response.Appointments,
		Count:        res.Count,
	})
}

// UpdateBookedAppointment ...
// @Summary UpdateBookedAppointment
// @Description UpdateDoctorNote - API to update appointment
// @Tags Appointment
// @Accept json
// @Produce json
// @Param UpdateAppointmentReq body model_booking_service.UpdateAppointmentReq true "UpdateAppointmentReq"
// @Success 200 {object} model_booking_service.Appointment
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/appointment [put]
func (h *HandlerV1) UpdateBookedAppointment(c *gin.Context) {
	var (
		body        model_booking_service.UpdateAppointmentReq
		jsonMarshal protojson.MarshalOptions
	)
	jsonMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	if e.HandleError(c, err, h.log, http.StatusBadRequest, "UpdateBookedAppointment") {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	res, err := h.serviceManager.BookingService().BookedAppointment().UpdateAppointment(ctx, &pb.UpdateAppointmentReq{
		AppointmentDate: body.AppointmentDate,
		AppointmentTime: body.AppointmentTime,
		Duration:        body.Duration,
		Key:             body.Key,
		ExpiresAt:       body.ExpiresAt,
		Status:          body.PatientStatus,
		Field:           "id",
		Value:           body.BookedAppointmentId,
		DoctorServiceId: body.DoctorServiceId,
		PatientProblem:  body.PatientProblem,
		PaymentType:     body.PaymentType,
		PaymentAmount:   float32(body.PaymentAmount),
		UserId:          body.UserId,
	})

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, "UpdateBookedAppointment") {
		return
	}

	c.JSON(http.StatusOK, model_booking_service.Appointment{
		Id:                   res.Id,
		DepartmentId:         res.DepartmentId,
		DoctorId:             res.DoctorId,
		PatientId:            res.PatientId,
		AppointmentDate:      res.AppointmentDate,
		AppointmentStartTime: res.AppointmentTime,
		Duration:             res.Duration,
		Key:                  res.Key,
		ExpiresAt:            res.ExpiresAt,
		PatientStatus:        res.Status,
		DoctorServiceId:      res.DoctorServiceId,
		PatientProblem:       res.PatientProblem,
		PaymentType:          res.PaymentType,
		PaymentAmount:        float64(res.PaymentAmount),
		UserId:               res.UserId,
		CreatedAt:            res.CreatedAt,
		UpdatedAt:            e.UpdateTimeFilter(res.UpdatedAt),
	})
}

// DeleteBookedAppointment ...
// @Summary DeleteBookedAppointment
// @Description DeleteBookedAppointment - API to delete an appointment
// @Tags Appointment
// @Accept json
// @Produce json
// @Param id query integer true "id"
// @Success 200 {object} models.StatusRes
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/appointment [delete]
func (h *HandlerV1) DeleteBookedAppointment(c *gin.Context) {
	id := c.Query("id")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	status, err := h.serviceManager.BookingService().BookedAppointment().DeleteAppointment(ctx, &pb.AppointmentFieldValueReq{
		Field:    "id",
		Value:    id,
		IsActive: false,
	})

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, "DeleteBookedAppointment") {
		return
	}

	c.JSON(http.StatusOK, models.StatusRes{Status: status.Status})
}

// GetAppointmentDates ...
// @Summary GetAppointmentDates
// @Description AppointmentDates - API to get Appointment Dates
// @Tags Appointment
// @Accept json
// @Produce json
// @Success 200 {object} []models.Date
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/appointment/dates [get]
func (h *HandlerV1) GetAppointmentDates(c *gin.Context) {
	dayNow := time.Now().Add(time.Hour * 5)
	var date []models.Date
	for i := 0; i < 8; i++ {
		day := dayNow.AddDate(0, 0, i)
		if day.Weekday() != time.Sunday {
			date = append(date, models.Date{
				Date: day.Format("2006-01-02"),
				Week: day.Format("Monday"),
			})
		} else if day.Weekday() != time.Sunday {
			i = i - 1
		}
	}
	c.JSON(http.StatusOK, date)
}

// Booking ...
// @Summary Booking
// @Description Booking - Api for create booked appointment
// @Tags Appointment
// @Security ApiKeyAuth
// @Accept json
// @Produce json
// @Param Booking body model_booking_service.BookingReq true "Booking"
// @Success 200 {object} []model_booking_service.BookingResp
// @Failure 400 {object} model_common.StandardErrorModel
// @Failure 500 {object} model_common.StandardErrorModel
// @Router /v1/appointment/booking [post]
func (h *HandlerV1) Booking(c *gin.Context) {
	var (
		body        model_booking_service.BookingReq
		times       []model_booking_service.BookingResp
		status      bool
		timeOfDay   bool
		jsonMarshal protojson.MarshalOptions
	)
	jsonMarshal.UseProtoNames = true

	err := c.ShouldBindJSON(&body)
	date, err := time.Parse("2006-01-02", body.Date)
	if e.HandleError(c, err, h.log, http.StatusBadRequest, "Booking") {
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(h.cfg.Context.Timeout))
	defer cancel()

	hours, err := h.serviceManager.BookingService().DoctorTimes().GetDoctorTimeByDoctorID(ctx, &pb.DoctorTimeFieldValueByIdReq{
		DoctorId:   body.DoctorID,
		DoctorDate: body.Date,
	})

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, "Booking") {
		return
	}

	doctor, err := h.serviceManager.HealthcareService().DoctorService().GetDoctorById(ctx, &healthcare.GetReqStrDoctor{
		Field:    "id",
		Value:    body.DoctorID,
		IsActive: false,
	})

	if e.HandleError(c, err, h.log, http.StatusInternalServerError, "Booking") {
		return
	}

	for _, hour := range hours.Hours {
		timeStatus, err := h.serviceManager.BookingService().DoctorTimes().GeDoctorTimeStatus(ctx, &pb.DoctorTimeStatusReq{
			DepartmentId: doctor.DepartmentId,
			DoctorId:     body.DoctorID,
			DoctorDate:   date.Format("2006-01-02"),
			StartTime:    hour,
		})
		if e.HandleError(c, err, h.log, http.StatusInternalServerError, "Booking") {
			return
		}
		timeOfDayStatus, err := time.Parse("15:04", hour)
		if e.HandleError(c, err, h.log, http.StatusInternalServerError, "Booking") {
			return
		}
		if timeStatus.Status == "available" {
			status = true
		} else {
			status = false
		}
		if timeOfDayStatus.Hour() <= 12 {
			timeOfDay = true
		} else {
			timeOfDay = false
		}

		times = append(times, model_booking_service.BookingResp{
			Time:      hour,
			Status:    status,
			TimeOfDay: timeOfDay,
		})
	}

	c.JSON(http.StatusOK, times)

}
