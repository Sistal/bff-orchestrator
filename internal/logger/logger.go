package logger

import (
	"os"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	instance *zap.Logger
	once     sync.Once
)

// Get retorna la instancia singleton del logger.
// En producción usa formato JSON; en desarrollo usa formato legible por consola.
func Get() *zap.Logger {
	once.Do(func() {
		env := os.Getenv("ENVIRONMENT")

		var cfg zap.Config
		if env == "production" {
			cfg = zap.NewProductionConfig()
			cfg.EncoderConfig.TimeKey = "timestamp"
			cfg.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
		} else {
			cfg = zap.NewDevelopmentConfig()
			cfg.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		}

		var err error
		instance, err = cfg.Build(zap.AddCallerSkip(0))
		if err != nil {
			// Fallback al logger de emergencia (escribe en stderr)
			instance = zap.NewExample()
		}
	})
	return instance
}

// Sync flushea los buffers pendientes. Llamar en el shutdown del servidor.
func Sync() {
	if instance != nil {
		_ = instance.Sync()
	}
}
