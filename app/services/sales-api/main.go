package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"

	"github.com/ardanlabs/conf/v3"
	"github.com/vitoraalmeida/service/app/services/sales-api/handlers"
	"github.com/vitoraalmeida/service/business/web/auth"
	"github.com/vitoraalmeida/service/business/web/v1/debug"
	"github.com/vitoraalmeida/service/foundation/keystore"
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
			// definir o host do serviço de debug em outro ip/rota para impossbilitar o acesso externo
			DebugHost string `conf:"default:0.0.0.0:4000"`
			//adicionar noprint no fim da tag de configuração caso não queira que essa info vá para o log
			//DebugHost       string        `conf:"default:0.0.0.0:4000,noprint"`
			//adicionar mask no fim da tag de configuração caso queira que apareça, mas mascarado
			//DebugHost       string        `conf:"default:0.0.0.0:4000,mask"`
		}
		// informações para lidar com autenticação
		Auth struct {
			KeysFolder string `conf:"default:zarf/keys/"`                           // informações com keys definidas a priori
			ActiveKID  string `conf:"default:54bb2165-71e1-41a6-af3e-7da4a0e1e2c1"` // nome da key PEM que será pré-definida
			Issuer     string `conf:"default:service project"`                      // define quem é o criador do token
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
	// App Starting

	log.Infow("starting service", "version", build)
	defer log.Infow("shutdown complete")

	// gera string que contém as informações de configução que foram usadas
	out, err := conf.String(&cfg)
	if err != nil {
		return fmt.Errorf("generating config for output: %w", err)
	}
	log.Infow("startup", "config", out)

	// -------------------------------------------------------------------------
	// Inicia suporte à autenticação

	log.Infow("startup", "status", "initializing authentication support")

	// Criação do armazenamento de chaves em memória usando chaves criadas anteriormente
	ks, err := keystore.NewFS(os.DirFS(cfg.Auth.KeysFolder))
	if err != nil {
		return fmt.Errorf("reading keys: %w", err)
	}

	authCfg := auth.Config{
		Log:       log,
		KeyLookup: ks,
	}

	// objeto que armazena informações para lidar com autenticação/autorização
	auth, err := auth.New(authCfg)
	if err != nil {
		return fmt.Errorf("constructing auth: %w", err)
	}

	// -------------------------------------------------------------------------
	// Inicia serviço de debug

	log.Infow("startup", "status", "debug v1 router started", "host", cfg.Web.DebugHost)

	// iniciamos uma goroutine separada para servir os endpoints de debug
	// caso a goroutine principal morra, não tem problema esta fica orfã, pois
	// ela apenas realiza leitura
	go func() {
		if err := http.ListenAndServe(cfg.Web.DebugHost, debug.Mux(build, log)); err != nil {
			log.Errorw("shutdown", "status", "debug v1 router closed", "host", cfg.Web.DebugHost, "ERROR", err)
		}
	}()

	// -------------------------------------------------------------------------
	// Inicia o serviço da API
	log.Infow("startup", "status", "initializing V1 API support")

	// Canal para onde poderá ser enviado sinais de SO para encerrar o programa
	shutdown := make(chan os.Signal, 1)
	// Notify registra que determinados sinais serão direcionados pelo channel
	// passado. SIGTERM = kubernetes envia quando finaliza o container.
	// SIGINT = Ctrl + c
	signal.Notify(shutdown, syscall.SIGINT, syscall.SIGTERM)

	// cria uma instâcia do nosso mux
	apiMux := handlers.APIMux(handlers.APIMuxConfig{
		Shutdown: shutdown,
		Log:      log,
		Auth:     auth,
	})

	// cria uma instância de http.Server customizada com os valores de configuração
	api := http.Server{
		Addr:         cfg.Web.APIHost,
		Handler:      apiMux,
		ReadTimeout:  cfg.Web.ReadTimeout,
		WriteTimeout: cfg.Web.WriteTimeout,
		IdleTimeout:  cfg.Web.IdleTimeout,
		ErrorLog:     zap.NewStdLog(log.Desugar()),
	}

	// cria uma channel para sinalizar erros na API
	// channels sem buffer (tamanho 1) garantem a quem enviou que o receptor
	// recebeu, com o custo de que o remetente fica esperando até que o
	// destinatário receba (aumento de latencia)
	serverErrors := make(chan error, 1)

	go func() {
		log.Infow("startup", "status", "api router started", "host", api.Addr)
		// executa o servidor e caso erros sejam retornados, envia pelo canal
		// assim garantimos que go routines que sejam iniciadas pelo servidor
		// não fiquem orfãs caso o servidor pare, pois a goroutine main assume
		// a responsabilidade
		serverErrors <- api.ListenAndServe()
	}()

	// -------------------------------------------------------------------------
	// Shutdown
	// select fica aguardando (blocking) enquanto algum dos casos ocorrer,
	// ou no nosso contexto enquanto algum sinal não chega pelas channels.
	// Serão sinais de shutdown ou de erro. O que ocorrer primeiro será executado
	select {
	case err := <-serverErrors:
		return fmt.Errorf("server error: %w", err)
		// Fica aguardando por sinais de shutdown e encerra caso chegue
	case sig := <-shutdown:
		log.Infow("shutdown", "status", "shutdown started", "signal", sig)
		defer log.Infow("shutdown", "status", "shutdown complete", "signal", sig)

		// cria um timer que será usado para, caso o servidor não consiga finalizar
		// gracefully (api.Shutdown), o programa é interrompido para que não
		// fique em executação eternamente.
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Web.ShutdownTimeout)
		defer cancel()

		// api.Shutdown fecham todos os listeners, as conexões que estão paradas
		// e então aguarda por trabalhos que já começaram terminaram (evitar
		// dados corrompidos). Se o contexto
		// passado finalizar antes disso acontecer, retorna um erro de contexto
		if err := api.Shutdown(ctx); err != nil {
			api.Close()
			return fmt.Errorf("could not stop server gracefully: %w", err)
		}
	}
	return nil
}
