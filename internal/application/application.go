package application

import (
	"log/slog"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/nathanjms/go-api-template/internal/awsHelper"
	"github.com/nathanjms/go-api-template/internal/database"
	"github.com/nathanjms/go-api-template/internal/env"
	"github.com/nathanjms/go-api-template/internal/jwtHelper"
)

type Config struct {
	BaseURL  string
	HTTPPort int
	JWT      struct {
		SecretKey string
	}
	AWS struct {
		Bucket string
	}
}

type Application struct {
	Config            Config
	DB                *database.DB
	sentryInitialized bool
	Logger            *slog.Logger
	S3                *awsHelper.S3Helper
	JWTService        *jwtHelper.JWTService
}

func New(logger *slog.Logger) (*Application, error) {
	app := &Application{}
	// --- Initialize Sentry Error reporting ---
	if err := initSentry(logger, app); err != nil {
		return nil, err
	}

	// --- Config ---
	cfg := initConfig()

	// --- DB ---
	db, err := database.New(env.GetString("DB_DSN", "root:password@tcp(localhost:3306)/api-db"))
	if err != nil {
		return nil, err
	}

	// --- AWS ---
	s3 := initS3(cfg.AWS.Bucket)

	// --- JWT ---
	jwtService, err := jwtHelper.NewJWTService(cfg.JWT.SecretKey)
	if err != nil {
		return nil, err
	}

	app.Config = cfg
	app.DB = db
	app.Logger = logger
	app.S3 = s3
	app.JWTService = jwtService

	return app, nil
}

func initSentry(logger *slog.Logger, app *Application) error {
	// --- Initialize Sentry Error reporting ---
	sentryDsn := env.GetString("SENTRY_DSN", "")
	if sentryDsn != "" {
		if err := sentry.Init(sentry.ClientOptions{
			Dsn:              env.GetString("SENTRY_DSN", ""),
			EnableTracing:    true,
			TracesSampleRate: 1.0,
		}); err != nil {
			return err
		}
		app.sentryInitialized = true // Set the flag after successful initialization
	} else {
		logger.Info("Sentry disabled")
	}

	return nil
}

func initConfig() Config {
	var cfg Config

	cfg.BaseURL = env.GetString("BASE_URL", "http://localhost")
	cfg.HTTPPort = env.GetInt("PORT", 3000)
	cfg.JWT.SecretKey = env.GetString("RSA_PRIVATE_KEY", "secret")
	cfg.AWS.Bucket = env.GetString("AWS_BUCKET", "bucket")

	return cfg
}

func initS3(bucket string) *awsHelper.S3Helper {
	accessKeyId := env.GetString("AWS_ACCESS_KEY", "secret")
	accessKeySecret := env.GetString("AWS_SECRET_KEY", "secret")
	awsAccountId := env.GetString("AWS_ACCOUNT_ID", "123456789")
	return awsHelper.New(accessKeyId, accessKeySecret, awsAccountId, bucket)
}

// Add a new Close method to Application
func (app *Application) Close() {
	if app.DB != nil {
		app.DB.Close()
	}
	if app.sentryInitialized {
		sentry.Flush(2 * time.Second)
	}
}

// reportError reports the error to Sentry and logs it
func (app *Application) ReportError(err error) {
	// 1. Log the error
	app.Logger.Error(err.Error())

	// 2. Capture the error in Sentry
	sentry.CaptureException(err)
}
