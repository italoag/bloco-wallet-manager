package logger

import (
	"go.uber.org/zap"
)

// Logger é a interface que define os métodos de logging que serão utilizados na aplicação.
type Logger interface {
	Info(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Debug(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
}

// zapLogger implementa a interface Logger utilizando o Uber Zap.
type zapLogger struct {
	logger *zap.Logger
}

// NewLogger inicializa um novo logger com base no nível de log fornecido.
func NewLogger(level string) Logger {
	var cfg zap.Config
	if level == "debug" {
		cfg = zap.NewDevelopmentConfig()
	} else {
		cfg = zap.NewProductionConfig()
	}

	logger, err := cfg.Build()
	if err != nil {
		panic("Falha ao inicializar o logger: " + err.Error())
	}

	return &zapLogger{logger: logger}
}

// Info registra uma mensagem de informação.
func (z *zapLogger) Info(msg string, fields ...zap.Field) {
	z.logger.Info(msg, fields...)
}

// Error registra uma mensagem de erro.
func (z *zapLogger) Error(msg string, fields ...zap.Field) {
	z.logger.Error(msg, fields...)
}

// Debug registra uma mensagem de debug.
func (z *zapLogger) Debug(msg string, fields ...zap.Field) {
	z.logger.Debug(msg, fields...)
}

// Warn registra uma mensagem de aviso.
func (z *zapLogger) Warn(msg string, fields ...zap.Field) {
	z.logger.Warn(msg, fields...)
}
