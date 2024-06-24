package app

import (
	pb "booking_service/genproto/booking_service"
	grpc_server "booking_service/internal/delivery/grpc/server"
	invest_grpc "booking_service/internal/delivery/grpc/services"
	"booking_service/internal/infrastructure/grpc_service_clients"
	repo "booking_service/internal/infrastructure/repository/postgresql"
	"booking_service/internal/pkg/config"
	"booking_service/internal/pkg/logger"
	"booking_service/internal/pkg/otlp"
	"booking_service/internal/pkg/postgres"
	"booking_service/internal/usecase"
	"fmt"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_zap "github.com/grpc-ecosystem/go-grpc-middleware/logging/zap"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type App struct {
	Config         *config.Config
	Logger         *zap.Logger
	DB             *postgres.PostgresDB
	GrpcServer     *grpc.Server
	ShutdownOTLP   func() error
	ServiceClients grpc_service_clients.ServiceClients
	//BrokerProducer event.BrokerProducer
}

func NewApp(cfg *config.Config) (*App, error) {
	// init logger
	logger, err := logger.New(cfg.LogLevel, cfg.Environment, cfg.APP+".log")
	if err != nil {
		return nil, err
	}

	//errkafkaProducer := kafka.NewProducer(cfg, logger)

	// otlp collector initialization
	shutdownOTLP, err := otlp.InitOTLPProvider(cfg)
	if err != nil {
		return nil, err
	}

	// init db
	db, err := postgres.New(cfg)
	if err != nil {
		return nil, err
	}

	// grpc server init
	grpcServer := grpc.NewServer(
		grpc.StreamInterceptor(grpc_middleware.ChainStreamServer(
			grpc_ctxtags.StreamServerInterceptor(),
			grpc_zap.StreamServerInterceptor(logger),
			grpc_recovery.StreamServerInterceptor(),
		)),
		grpc.UnaryInterceptor(grpc_server.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				grpc_ctxtags.UnaryServerInterceptor(),
				grpc_zap.UnaryServerInterceptor(logger),
				grpc_recovery.UnaryServerInterceptor(),
			),
			grpc_server.UnaryInterceptorData(logger),
		)),
	)

	return &App{
		Config:       cfg,
		Logger:       logger,
		DB:           db,
		GrpcServer:   grpcServer,
		ShutdownOTLP: shutdownOTLP,
		//BrokerProducer: kafkaProducer,
	}, nil
}

func (a *App) Run() error {

	// context timeout initialization
	contextTimeout, err := time.ParseDuration(a.Config.Context.Timeout)

	// Initialize Service Clients
	serviceClients, err := grpc_service_clients.New(a.Config)
	if err != nil {
		return fmt.Errorf("error during initialize service clients: %w", err)
	}
	a.ServiceClients = serviceClients

	// repositories initialization
	bookingAppointment := repo.NewBookingAppointment(a.DB)

	bookingPatients := repo.NewBookingPatients(a.DB)

	doctorNotes := repo.NewDoctorNotes(a.DB)

	bookingArchive := repo.NewBookingArchive(a.DB)

	doctorAvailability := repo.NewDoctorAvailability(a.DB)

	// usecase initialization

	appointmentsUseCase := usecase.NewBookedAppointments(bookingAppointment, contextTimeout)

	patientUseCase := usecase.NewBookedPatient(bookingPatients, contextTimeout)

	doctorNotesUseCase := usecase.NewBookedDoctorNotes(doctorNotes, contextTimeout)

	archiveUseCase := usecase.NewBookedArchive(bookingArchive, contextTimeout)

	doctorAvailabilityUseCase := usecase.NewBookedDoctorAvailability(doctorAvailability, contextTimeout)

	pb.RegisterBookedAppointmentsServiceServer(a.GrpcServer, invest_grpc.BookingAppointmentsNewRPC(a.Logger, appointmentsUseCase))

	pb.RegisterPatientsServiceServer(a.GrpcServer, invest_grpc.BookingPatientNewRPC(a.Logger, patientUseCase))

	pb.RegisterDoctorNotesServiceServer(a.GrpcServer, invest_grpc.BookingDoctorNotesNewRPC(a.Logger, doctorNotesUseCase))

	pb.RegisterArchiveServiceServer(a.GrpcServer, invest_grpc.BookingArchiveNewRPC(a.Logger, archiveUseCase))

	pb.RegisterDoctorTimeServiceServer(a.GrpcServer, invest_grpc.BookingDoctorAvailabilityNewRPC(a.Logger, doctorAvailabilityUseCase))
	a.Logger.Info("gRPC Server Listening", zap.String("url", a.Config.RPCPort))

	if err := grpc_server.Run(a.Config, a.GrpcServer); err != nil {
		return fmt.Errorf("gRPC fatal to serve grpc server over %s %w", a.Config.RPCPort, err)
	}

	return nil
}

func (a *App) Stop() {
	// close broker producer
	// closing client service connections
	a.ServiceClients.Close()
	// stop gRPC server
	a.GrpcServer.Stop()

	// database connection
	a.DB.Close()

	// shutdown otlp collector
	if err := a.ShutdownOTLP(); err != nil {
		a.Logger.Error("shutdown otlp collector", zap.Error(err))
	}

	// zap logger sync
	a.Logger.Sync()
}
