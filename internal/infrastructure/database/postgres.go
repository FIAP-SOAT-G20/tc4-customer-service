package database

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/internal/infrastructure/config"
	"github.com/FIAP-SOAT-G20/fiap-tech-challenge-3-lambda-auth-tf/internal/infrastructure/logger"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
)

type Database struct {
	*gorm.DB
	dbURL string
}

type GormLogger struct {
	*slog.Logger
}

func (l *GormLogger) LogMode(gormlogger.LogLevel) gormlogger.Interface {
	return l
}

func (l *GormLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.Logger.InfoContext(ctx, msg, "data", data)
}

func (l *GormLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.Logger.WarnContext(ctx, msg, "data", data)
}

func (l *GormLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.Logger.ErrorContext(ctx, msg, "data", data)
}

func (l *GormLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()

	if err != nil {
		l.Logger.ErrorContext(ctx, "database query failed",
			"error", err,
			"elapsed", elapsed,
			"sql", sql,
			"rows", rows,
		)
		return
	}

	l.Logger.DebugContext(ctx, "database query",
		"elapsed", elapsed,
		"sql", sql,
		"rows", rows,
	)
}

func NewPostgresConnection(cfg *config.Config, logger *logger.Logger) (*Database, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%d sslmode=disable",
		cfg.DBHost,
		cfg.DBUser,
		cfg.DBPassword,
		cfg.DBName,
		cfg.DBPort,
	)

	dbURL := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", cfg.DBUser, cfg.DBPassword, cfg.DBHost, cfg.DBPort, cfg.DBName)

	// Configure GORM with slog logger
	gormConfig := &gorm.Config{
		Logger:      &GormLogger{Logger: logger.Logger},
		PrepareStmt: true,
	}

	db, err := gorm.Open(postgres.Open(dsn), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	// Get generic database object sql.DB to use its functions
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get database instance: %w", err)
	}

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)

	// SetMaxOpenConns sets the maximum number of open connections to the database
	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)

	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused
	sqlDB.SetConnMaxLifetime(cfg.DBMaxLifetime)

	return &Database{db, dbURL}, nil
}
