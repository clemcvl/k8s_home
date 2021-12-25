package pkg

import (
	"fmt"
	"html/template"
	"net/http"

	log "github.com/sirupsen/logrus"
)

type Component struct {
	Name string `json:"name"`
}

type DataPasser struct {
	state Component
}

func StartServer() {
	passer := &DataPasser{state: Component{Name: "toto"}}
	log.Println("setuping routes")
	setupRoutes(passer)
	log.Println("starting server")
	err := http.ListenAndServe(":8080", logRequest(http.DefaultServeMux))
	if err != nil {
		log.Warnln("Error creating webserver", err)
	}
}

func setupRoutes(passer *DataPasser) {

	// Setup Common routes
	http.HandleFunc("/ping", pingRoute)
	http.HandleFunc("/data", passer.getTemplate)

}

func pingRoute(w http.ResponseWriter, req *http.Request) {
	log.Println("ping")
	fmt.Fprintf(w, "pong")
}

func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.EscapedPath() != "/ping" {
			log.WithFields(log.Fields{
				"url":           r.URL,
				"method":        r.Method,
				"remote_addr":   r.RemoteAddr,
				"forwarded_for": r.Header.Get("X-Forwarded-For"),
			}).Info("received request")
		}
		handler.ServeHTTP(w, r)
	})
}

func (p *DataPasser) getTemplate(w http.ResponseWriter, req *http.Request) {
	state := p.state
	temp, err := template.ParseGlob("templates/*")
	if err != nil {
		log.Errorln("Error parsing template", err)
		fmt.Fprintf(w, "Error parsing template")
		return
	}

	templateName := fmt.Sprintf("%s.html", req.URL.Path[1:])
	err = temp.ExecuteTemplate(w, templateName, state)
	if err != nil {
		log.Warnln("Error executing template", err)
		fmt.Fprintf(w, "Error executing template")
		return
	}
}
