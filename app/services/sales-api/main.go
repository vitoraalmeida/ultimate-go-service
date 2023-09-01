package main

import (
	"errors"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
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
	// Configuration

	// usar tipos literais (structs sem nomes) garantimos que quem vai receber
	// a configuração deve ser preciso sobre o que vai chegar na sua API. Assim
	// removemos o risco de criar abstrações que podem ou não abarcar elementos
	// que queremos utilizar. Precisão!
	// Preciso de uma banana, mas como gorilas de vez enquando possuem bananas,
	// vou esperar um gorila inteiro, sendo que pode ser que ele não esteja com
	// uma banana? Não, vou esperar especificamente bananas.
	cfg := struct {
		conf.Version
		// Podemos  definir variávies de ambiente que conf irá procurar let
		// nesse caso: WEB_RED_TIMEOUT e assim sucessivamente
		Web struct {
			ReadTimeout     time.Duration `conf:"default:5s"` // define configs padrão
			WriteTimeout    time.Duration `conf:"default:10s"`
			IdleTimeout     time.Duration `conf:"default:120s"`
			ShutdownTimeout time.Duration `conf:"default:20s"`
			APIHost         string        `conf:"default:0.0.0.0:3000"`
			DebugHost       string        `conf:"default:0.0.0.0:4000"`
		}
	}{
		Version: conf.Version{
			Build: build,
			Desc:  "copyright information here",
		},
	}

	const prefix = "SALES" // prefixo para as variáveis de ambiente -> SALES_WEB_READ_TIMEOUT
	// Parse recebe o objeto que criamos para a configuração e procura por
	// variáveis de ambiente ou flags de linha de comando que sobrescrevam o
	// default.
	help, err := conf.Parse(prefix, &cfg)
	if err != nil {
		// na tentativa de fazer o parsing de command line flags, se o usuário
		// digitou "help", conf.Parse retorna um erro específico para indicar
		// que devemos mostrar a mensagem de ajuda
		if errors.Is(err, conf.ErrHelpWanted) {
			fmt.Println(help)
			return nil
		}
		return fmt.Errorf("parsing config: %w", err)
	}

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
