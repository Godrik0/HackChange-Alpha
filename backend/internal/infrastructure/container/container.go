package container

import (
	"fmt"

	"github.com/Godrik0/HackChange-Alpha/backend/internal/application/services"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/config"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/domain/interfaces"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/infrastructure/http"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/infrastructure/http/handlers"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/infrastructure/promo"
	"github.com/Godrik0/HackChange-Alpha/backend/internal/infrastructure/providers"
	"gorm.io/gorm"
)

type Container struct {
	Config *config.Config
	Logger interfaces.Logger

	LoggerProvider     providers.LoggerProvider
	DatabaseProvider   providers.DatabaseProvider
	RepositoryProvider providers.RepositoryProvider
	MLServiceProvider  providers.MLServiceProvider

	DB       *gorm.DB
	MLClient interfaces.MLService

	ClientRepo interfaces.ClientRepository

	ClientService  interfaces.ClientService
	ScoringService interfaces.ScoringService

	ClientHandler *handlers.ClientHandler

	HTTPServer *http.Server
}

type Options struct {
	LoggerProvider     providers.LoggerProvider
	DatabaseProvider   providers.DatabaseProvider
	RepositoryProvider providers.RepositoryProvider
	MLServiceProvider  providers.MLServiceProvider
}

func New(cfg *config.Config, opts *Options) (*Container, error) {
	c := &Container{
		Config: cfg,
	}

	c.setupProviders(opts)

	if err := c.initLogger(); err != nil {
		return nil, fmt.Errorf("failed to initialize logger: %w", err)
	}

	if err := c.initDatabase(); err != nil {
		return nil, fmt.Errorf("failed to initialize database: %w", err)
	}

	if err := c.initRepositories(); err != nil {
		return nil, fmt.Errorf("failed to initialize repositories: %w", err)
	}

	if err := c.initMLClient(); err != nil {
		return nil, fmt.Errorf("failed to initialize ML client: %w", err)
	}

	if err := c.initServices(); err != nil {
		return nil, fmt.Errorf("failed to initialize services: %w", err)
	}

	if err := c.initHandlers(); err != nil {
		return nil, fmt.Errorf("failed to initialize handlers: %w", err)
	}

	if err := c.initHTTPServer(); err != nil {
		return nil, fmt.Errorf("failed to initialize HTTP server: %w", err)
	}

	return c, nil
}

func (c *Container) setupProviders(opts *Options) {
	if opts == nil {
		opts = &Options{}
	}

	if opts.LoggerProvider != nil {
		c.LoggerProvider = opts.LoggerProvider
	} else {
		c.LoggerProvider = &providers.DefaultLoggerProvider{}
	}

	if opts.DatabaseProvider != nil {
		c.DatabaseProvider = opts.DatabaseProvider
	} else {
		c.DatabaseProvider = &providers.PostgresProvider{}
	}

	if opts.RepositoryProvider != nil {
		c.RepositoryProvider = opts.RepositoryProvider
	} else {
		c.RepositoryProvider = &providers.DefaultRepositoryProvider{}
	}

	if opts.MLServiceProvider != nil {
		c.MLServiceProvider = opts.MLServiceProvider
	} else {
		c.MLServiceProvider = &providers.DefaultMLServiceProvider{}
	}
}

func (c *Container) initLogger() error {
	logger, err := c.LoggerProvider.ProvideLogger(c.Config)
	if err != nil {
		return err
	}

	c.Logger = logger
	return nil
}

func (c *Container) initDatabase() error {
	db, err := c.DatabaseProvider.ProvideDatabase(c.Config, c.Logger)
	if err != nil {
		return err
	}

	if err := c.DatabaseProvider.RunMigrations(db, c.Logger); err != nil {
		return err
	}

	c.DB = db
	return nil
}

func (c *Container) initRepositories() error {
	c.ClientRepo = c.RepositoryProvider.ProvideClientRepository(c.DB, c.Logger)
	return nil
}

func (c *Container) initMLClient() error {
	c.MLClient = c.MLServiceProvider.ProvideMLService(c.Config, c.Logger)
	return nil
}

func (c *Container) initServices() error {
	promoProvider := promo.NewStaticPromoProvider()

	c.ClientService = services.NewClientService(
		c.ClientRepo,
		c.MLClient,
		c.Logger,
	)

	c.ScoringService = services.NewScoringService(
		c.ClientRepo,
		c.MLClient,
		promoProvider,
		c.Logger,
	)
	return nil
}

func (c *Container) initHandlers() error {
	c.ClientHandler = handlers.NewClientHandler(
		c.ClientService,
		c.ScoringService,
		c.Logger,
	)
	return nil
}

func (c *Container) initHTTPServer() error {
	c.HTTPServer = http.NewServer(
		c.Config,
		c.ClientHandler,
		c.Logger,
	)
	return nil
}

func (c *Container) Close() error {
	c.Logger.Info("Closing container resources")

	if c.DB != nil {
		sqlDB, err := c.DB.DB()
		if err != nil {
			c.Logger.Error("Failed to get SQL DB", "error", err)
			return err
		}
		if err := sqlDB.Close(); err != nil {
			c.Logger.Error("Failed to close database", "error", err)
			return err
		}
	}

	c.Logger.Info("Container resources closed")
	return nil
}
