package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/vitoraalmeida/service/foundation/logger"
	"go.uber.org/zap"
)

// é alterado por ldflags
var build = "develop"

func main() {
	// construímos nosso logger e passaremos ele concretamente para os componentes
	// que precisarmos. Não devemos adicionar loggers em contexts, pois assim
	// acabamos passando para todos os lugares de forma desnecessária
	log, err := logger.New("SALES-API")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Zap doc
	// Sync calls the underlying Core's Sync method, flushing any buffered log entries.
	// Applications should take care to call Sync before exiting.
	defer log.Sync()

	if err := run(log); err != nil {
		log.Errorw("startup", "ERROR", err)
		log.Sync()
		os.Exit(1)
	}
}

// coordena a inicialização e desligamento do sistema
func run(log *zap.SugaredLogger) error {

	// -------------------------------------------------------------------------
	// GOMAXPROCS

	log.Infow("startup", "GOMAXPROCS", runtime.GOMAXPROCS(0), "BUILD", build)

	// -------------------------------------------------------------------------

	// Canal para onde poderá ser enviado sinais de SO para encerrar o programa
	shutdown := make(chan os.Signal, 1)
	// Notify registra que determinados sinais serão direcionados pelo channel
	// passado. SIGTERM = kubernetes envia quando finaliza o container.
	// SIGINT = Ctrl + c
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// Fica aguardando por sinais de shutdown e encerra caso chegue
	sig := <-shutdown

	log.Infow("shutdown", "status", "shutdown started", "signal", sig)
	defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)
	// executa ações necessárias para terminar de forma segura

	return nil
}
