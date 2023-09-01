// Recebe log estruturado e devolve em formato legível para humanos
package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

var service string

func init() {
	// vai filtrar os logs estruturados que tiverem o campo "service"
	flag.StringVar(&service, "service", "", "filtra qual serviço deve ser convertido para texto legível")
}

func main() {
	flag.Parse()
	var b strings.Builder

	service := strings.ToLower(service)

	scanner := bufio.NewScanner(os.Stdin)
	// scanner.Scan lê uma linha por vez
	for scanner.Scan() {
		s := scanner.Text()

		// o log que estamos utilizando está estruturado em json, então serializamos
		// para um mapa com chaves do tipo string
		m := make(map[string]any)
		err := json.Unmarshal([]byte(s), &m)
		if err != nil {
			if service == "" {
				fmt.Println(s)
			}
			continue
		}

		// checa se foi adicionado um filtro de serviço e para o serviço atual
		// se ele for diferente do que queremos filtrar, ignora e busca o próximo
		if service != "" && strings.ToLower(m["service"].(string)) != service {
			continue
		}

		// adiciona trace
		traceID := "00000000-0000-0000-0000-000000000000"
		if v, ok := m["trace_id"]; ok {
			traceID = fmt.Sprintf("%v", v)
		}

		// constrói a string que será mostrada no log com os campos desejados
		// na ordem desejada
		b.Reset()
		b.WriteString(fmt.Sprintf("%s: %s: %s: %s: %s: %s: ",
			m["service"],
			m["ts"],
			m["level"],
			traceID,
			m["caller"],
			m["msg"],
		))

		// Adiciona o restante dos campos do log estruturado
		for k, v := range m {
			switch k {
			case "service", "ts", "level", "trace_id", "caller", "msg":
				continue
			}

			b.WriteString(fmt.Sprintf("%s[%v]: ", k, v))
		}

		// Write the new log format, removing the last :
		out := b.String()
		fmt.Println(out[:len(out)-2])
	}

	if err := scanner.Err(); err != nil {
		log.Println(err)
	}
}
