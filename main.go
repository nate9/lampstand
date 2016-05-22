package main // import "github.com/nate9/lampstand"

import (
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/julienschmidt/httprouter"
	"net/http"
	"os"
)

func main() {
	log.SetFormatter(&log.TextFormatter{DisableColors:true})
	log.SetOutput(os.Stdout)
	service := NewPassageService("./bible.db")
	bind := fmt.Sprintf("%s:%s", os.Getenv("OPENSHIFT_GO_IP"), os.Getenv("OPENSHIFT_GO_PORT"))
	if bind == ":" {
		bind = "localhost:8080"
	}
	log.Info("Serving lampstand on : " + bind)
	log.Fatal(RunPassageService(service, bind))
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		log.Error(err)
	}
}

func RunPassageService(s *PassageService, host string) error {
	router := httprouter.New()
	router.GET("/api/:version/verses", s.findVerses)
	router.NotFound = http.FileServer(http.Dir("static"))
	return http.ListenAndServe(host, router)
}
