package server

import (
	"github.com/ipaqsa/netcom/logger"
	"log"
	"net/http"
	"node/pkg/service"
)

var infoLogger = logger.NewLogger("INFO")
var errorLogger = logger.NewLogger("ERROR")

func Run(port string) {
	service.NewNode()
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("web/static"))))
	http.HandleFunc("/register", registerHandler)
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/remove", removeHandler)
	http.HandleFunc("/update", updateHandler)
	http.HandleFunc("/file", fileHandler)
	http.HandleFunc("/profile", profileHandler)
	http.HandleFunc("/send", sendHandler)
	http.HandleFunc("/", homeHandler)
	infoLogger.Printf("node start at %s", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
