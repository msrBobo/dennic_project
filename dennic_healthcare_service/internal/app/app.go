package app

import (
	pb "Healthcare_Evrone/genproto/healthcare-service"
	grpc_server "Healthcare_Evrone/internal/delivery/grpc/server"
	invest_grpc "Healthcare_Evrone/internal/delivery/grpc/services"
	"Healthcare_Evrone/internal/infrastructure/grpc_service_clients"
	repo "Healthcare_Evrone/internal/infrastructure/repository/postgresql"
	"Healthcare_Evrone/internal/pkg/config"
	"Healthcare_Evrone/internal/pkg/logger"
	"Healthcare_Evrone/internal/pkg/otlp"
	"Healthcare_Evrone/internal/pkg/postgres"
	"Healthcare_Evrone/internal/usecase"
	"Healthcare_Evrone/internal/usecase/event"
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
	BrokerProducer event.BrokerProducer
}

func NewApp(cfg *config.Config) (*App, error) {
	// init logger
	logger, err := logger.New(cfg.LogLevel, cfg.Environment, cfg.APP+".log")
	if err != nil {
		return nil, err
	}

	//kafkaProducer := kafka.NewProducer(cfg, logger)

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
	}, nil
}

func (a *App) Run() error {
	var (
		contextTimeout time.Duration
	)

	// context timeout initialization
	contextTimeout, err := time.ParseDuration(a.Config.Context.Timeout)
	if err != nil {
		return fmt.Errorf("error during parse duration for context timeout : %w", err)
	}
	// Initialize Service Clients
	serviceClients, err := grpc_service_clients.New(a.Config)
	if err != nil {
		return fmt.Errorf("error during initialize service clients: %w", err)
	}
	a.ServiceClients = serviceClients

	// repositories initialization
	healthRepo := repo.NewDoctorRepo(a.DB)
	department := repo.NewDepartmentRepo(a.DB)
	specialization := repo.NewSpecializationRepo(a.DB)
	dwh := repo.NewDoctorWorkingHoursRepo(a.DB)
	ds := repo.NewDoctorServicesRepo(a.DB)
	reas := repo.NewReasonsRepo(a.DB)

	// usecase initialization
	doctorWorkingHours := usecase.NewDoctorWorkingHoursService(contextTimeout, dwh)
	pb.RegisterDoctorWorkingHoursServiceServer(a.GrpcServer, invest_grpc.DoctorWorkingHoursServiceRPC(a.Logger, doctorWorkingHours))

	doctorUsecase := usecase.NewDoctorService(contextTimeout, healthRepo)
	pb.RegisterDoctorServiceServer(a.GrpcServer, invest_grpc.DoctorRPC(a.Logger, doctorUsecase))

	departmentUsecase := usecase.NewDepartmentService(contextTimeout, department)
	pb.RegisterDepartmentServiceServer(a.GrpcServer, invest_grpc.DepartmentRPC(a.Logger, departmentUsecase))

	specializationUsecase := usecase.NewSpecializationService(contextTimeout, specialization)
	pb.RegisterSpecializationServiceServer(a.GrpcServer, invest_grpc.SpecializationRPC(a.Logger, specializationUsecase, a.BrokerProducer))

	doctorsServicesUsecase := usecase.NewDoctorServices(contextTimeout, ds)
	pb.RegisterDoctorsServiceServer(a.GrpcServer, invest_grpc.DoctorsServiceRPC(a.Logger, doctorsServicesUsecase))

	reasonsUsecase := usecase.NewReasons(contextTimeout, reas)
	pb.RegisterReasonsServiceServer(a.GrpcServer, invest_grpc.ReasonsServiceRPC(a.Logger, reasonsUsecase))

	a.Logger.Info("gRPC Server Listening", zap.String("url", a.Config.RPCPort))
	if err := grpc_server.Run(a.Config, a.GrpcServer); err != nil {
		return fmt.Errorf("gRPC fatal to serve grpc server over %s %w", a.Config.RPCPort, err)
	}
	return nil
}

func (a *App) Stop() {
	// close broker producer
	a.BrokerProducer.Close()
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
