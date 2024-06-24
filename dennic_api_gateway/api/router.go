package api

import (
	// "github.com/casbin/casbin/v2"
	_ "dennic_api_gateway/api/docs"
	"dennic_api_gateway/api/middleware/casbin"
	"dennic_api_gateway/internal/pkg/redis"
	"fmt"
	"time"

	v1 "dennic_api_gateway/api/handlers/v1"
	"dennic_api_gateway/api/middleware"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"

	grpcClients "dennic_api_gateway/internal/infrastructure/grpc_service_client"
	"dennic_api_gateway/internal/pkg/config"
)

type RouteOption struct {
	Config         *config.Config
	Logger         *zap.Logger
	ContextTimeout time.Duration
	Service        grpcClients.ServiceClient
	Redis          *redis.RedisDB
	//BrokerProducer event.BrokerProducer

}

// NewRoute
// @title Dennic Project
// @version 1.7
// @host swag.dennic.uz
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func NewRoute(option RouteOption) *gin.Engine {
	router := gin.New()

	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	HandlerV1 := v1.New(&v1.HandlerV1Config{
		Config:         option.Config,
		Logger:         option.Logger,
		ContextTimeout: option.ContextTimeout,
		Service:        option.Service,
		Redis:          option.Redis,

		//BrokerProducer: option.BrokerProducer,
	})

	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AllowCredentials = true
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowBrowserExtensions = true
	corsConfig.AllowMethods = []string{"*"}
	router.Use(cors.New(corsConfig))

	err := HandlerV1.CronTime()
	if err != nil {
		fmt.Println(err)
	}

	router.Use(middleware.GinTracing())

	router.Use(casbin.NewAuthorizer())

	api := router.Group("/v1")

	api.POST("/file-upload", HandlerV1.UploadFile)
	api.POST("/support", HandlerV1.Support)

	// customer
	customer := api.Group("/customer")
	customer.POST("/register", HandlerV1.Register)
	customer.POST("/verify", HandlerV1.Verify)
	customer.POST("/forget-password", HandlerV1.ForgetPassword)
	customer.PUT("/update-password", HandlerV1.UpdatePassword)
	customer.POST("/verify-otp-code", HandlerV1.VerifyOtpCode)
	customer.POST("/send-otp", HandlerV1.SenOtpCode)
	customer.POST("/login", HandlerV1.Login)
	customer.POST("/logout", HandlerV1.LogOut)

	// user
	user := api.Group("/user")
	user.GET("/get", HandlerV1.GetUserByID)
	user.GET("/", HandlerV1.ListUsers)
	user.PUT("/update", HandlerV1.UpdateUser)
	user.PUT("/update-refresh-token", HandlerV1.UpdateRefreshToken)
	user.DELETE("/", HandlerV1.DeleteUser)

	token := api.Group("/token")
	token.GET("/get-token", HandlerV1.GetTokens)

	// archive
	// archive := api.Group("/archive")
	// archive.POST("/", HandlerV1.CreateArchive)
	// archive.GET("/get", HandlerV1.GetArchive)
	// archive.GET("/", HandlerV1.ListArchive)
	// archive.PUT("/", HandlerV1.UpdateArchive)
	// archive.DELETE("/", HandlerV1.DeleteArchive)

	// doctor notes
	doctorNote := api.Group("/doctor-notes")
	doctorNote.GET("/get", HandlerV1.GetDoctorNote)
	doctorNote.GET("/", HandlerV1.ListDoctorNotes)

	// appointment
	appointment := api.Group("/appointment")
	appointment.POST("/create", HandlerV1.CreateBookedAppointment)
	appointment.POST("/booking", HandlerV1.Booking)
	appointment.GET("/get", HandlerV1.GetBookedAppointment)
	appointment.GET("/dates", HandlerV1.GetAppointmentDates)
	appointment.GET("/", HandlerV1.ListBookedAppointments)
	appointment.PUT("/", HandlerV1.UpdateBookedAppointment)
	appointment.DELETE("/", HandlerV1.DeleteBookedAppointment)

	// doctorTime
	doctorTime := api.Group("/doctor-time")
	doctorTime.GET("/get", HandlerV1.GetDoctorTimes)
	doctorTime.GET("/", HandlerV1.ListDoctorTimes)

	// patient
	patient := api.Group("/patient")
	patient.POST("/create", HandlerV1.CreatePatient)
	patient.GET("/get", HandlerV1.GetPatient)
	patient.GET("/", HandlerV1.ListPatient)
	patient.PUT("/", HandlerV1.UpdatePatient)
	patient.PUT("/phone", HandlerV1.UpdatePhonePatient)
	patient.DELETE("/", HandlerV1.DeletePatient)

	// department
	department := api.Group("/department")
	department.POST("/", HandlerV1.CreateDepartment)
	department.GET("/get", HandlerV1.GetDepartment)
	department.GET("/", HandlerV1.ListDepartments)
	department.PUT("/", HandlerV1.UpdateDepartment)
	department.DELETE("/", HandlerV1.DeleteDepartment)

	// doctor
	doctor := api.Group("/doctor")
	doctor.POST("/", HandlerV1.CreateDoctor)
	doctor.GET("/get", HandlerV1.GetDoctor)
	doctor.GET("/", HandlerV1.ListDoctors)
	doctor.GET("/top", HandlerV1.ListTopDoctors)
	doctor.GET("/spec", HandlerV1.ListDoctorsBySpecializationId)
	doctor.PUT("/", HandlerV1.UpdateDoctor)
	doctor.DELETE("/", HandlerV1.DeleteDoctor)

	// specialization
	specialization := api.Group("/specialization")
	specialization.GET("/get", HandlerV1.GetSpecialization)
	specialization.GET("/", HandlerV1.ListSpecializations)

	// doctorServices
	doctorServices := api.Group("/doctor-services")
	doctorServices.GET("/get", HandlerV1.GetDoctorService)
	doctorServices.GET("/", HandlerV1.ListDoctorServices)

	// doctorWorkingHours

	doctorWorkingHours := api.Group("/doctor-working-hours")
	doctorWorkingHours.GET("/get", HandlerV1.GetDoctorWorkingHours)
	doctorWorkingHours.GET("/", HandlerV1.ListDoctorWorkingHours)

	// reasons
	reasons := api.Group("/reasons")
	reasons.GET("/get", HandlerV1.GetReasons)
	reasons.GET("/", HandlerV1.ListReasons)

	// session
	session := api.Group("session")
	session.GET("/", HandlerV1.GetUserSessions)
	session.DELETE("/", HandlerV1.DeleteSessionById)

	url := ginSwagger.URL("swagger/doc.json")
	api.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

	return router
}
