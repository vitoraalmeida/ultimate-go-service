// Package logger provê uma função de conveniência para contruir um logger para
// se utilizado em outros lugares
package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// New constrói um Sugared Logger que escreve no stdout e provê timestamps
// legíveis para humanos
func New(service string, outputPaths ...string) (*zap.SugaredLogger, error) {
	config := zap.NewProductionConfig()

	// formato = "2023-08-30T11:26:10.471-0300"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.DisableStacktrace = true
	config.InitialFields = map[string]any{
		"service": service,
	}

	config.OutputPaths = []string{"stdout"}
	if outputPaths != nil {
		config.OutputPaths = outputPaths
	}

	log, err := config.Build(zap.WithCaller(true))
	if err != nil {
		return nil, err
	}

	// SugaredLogger permite criar logs sem ser necessário criar de forma estruturada
	// sem construir definindo os tipos, chaves e valores manualmente
	return log.Sugar(), nil
}
