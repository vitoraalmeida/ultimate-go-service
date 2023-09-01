// debug provê um handler para os endpoints de debug
package debug

import (
	"expvar"
	"net/http"
	"net/http/pprof"
)

// StandardLibraryMux registra todas as rotas de debug da stdlib num novo mux
// evitando usar o DefaultServerMux, pois ele pode ser acessado/modificado por
// qualquer lib terceira que importarmos, expondo ao risco de alguma delas
// adicionar algum handler no nosso serviço sem sabermos.
func StandardLibraryMux() *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("/debug/pprof/", pprof.Index)
	mux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	mux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	mux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	mux.HandleFunc("/debug/pprof/trace", pprof.Trace)
	mux.Handle("/debug/vars", expvar.Handler())

	return mux
}
