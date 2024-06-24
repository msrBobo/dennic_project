package v1

import (
	"context"
	"database/sql"
	"dennic_api_gateway/genproto/booking_service"
	"dennic_api_gateway/genproto/healthcare-service"
	bk "dennic_api_gateway/internal/infrastructure/grpc_service_client"
	hs "dennic_api_gateway/internal/infrastructure/grpc_service_client"
	"errors"
	"time"
)

//// CronTime Expire promo codes every 24 hours
//func CronTime() {
//	c := cron.New()
//	_, err := c.AddFunc("0 0 * * *", func() {
//		start := time.Now().Add(time.Hour * 5)
//		if err := DoctorsTime; err != nil {
//		}
//		elapsed := time.Since(start)
//		fmt.Println("Function execution time:", elapsed)
//	})
//	if err != nil {
//		return
//	}
//}

type DoctorTime struct {
	hs hs.HealthcareServiceI
	bk bk.BookingServiceI
}

func (r *DoctorTime) DoctorsTime() error {
	doctors, err := r.hs.DoctorService().GetAllDoctors(context.Background(), &healthcare.GetAllDoctorS{})
	if err != nil {
		return err
	}
	//date := time.Now().AddDate(0, 0, 8)
	date := time.Now().Add(time.Hour * 5)
	for i := 0; i < 8; i++ {
		if date.Weekday() == time.Sunday {
			date = date.AddDate(0, 0, 1)
		}

		for _, doctor := range doctors.DoctorHours {

			StartTime, err := time.Parse("15:04:05", doctor.StartTime)
			if err != nil {
				return err
			}
			FinishTime, err := time.Parse("15:04:05", doctor.FinishTime)
			if err != nil {
				return err
			}
			for t := StartTime; t.Before(FinishTime); t = t.Add(time.Hour) {
				if t.Format("15:04") != "12:00" {
					res, err := r.bk.DoctorTimes().GeDoctorTimeStatus(context.Background(), &booking_service.DoctorTimeStatusReq{
						DepartmentId: doctor.DepartmentId,
						DoctorId:     doctor.Id,
						DoctorDate:   date.Format("2006-01-02"),
						StartTime:    StartTime.String(),
					})
					if res.Status == "" || errors.Is(err, sql.ErrNoRows) {
						_, err = r.bk.DoctorTimes().CreateDoctorTime(context.Background(), &booking_service.CreateDoctorTimeReq{
							DepartmentId: doctor.DepartmentId,
							DoctorId:     doctor.Id,
							DoctorDate:   date.Format("2006-01-02"),
							StartTime:    doctor.StartTime,
							EndTime:      doctor.FinishTime,
							Status:       "available",
						})
						if err != nil {
							return err
						}
					}

				}
			}
		}
		date = date.AddDate(0, 0, 1)
	}
	return nil
}
