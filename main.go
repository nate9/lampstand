package main // import "github.com/nate9/lampstand"

import (
	"fmt"
	"github.com/julienschmidt/httprouter"
	"log"
	"net/http"
	"os"
)

func main() {
	service := NewPassageService("./bible.db")
	bind := fmt.Sprintf("%s:%s", os.Getenv("OPENSHIFT_GO_IP"), os.Getenv("OPENSHIFT_GO_PORT"))
	if bind == ":" {
		bind = "localhost:8080"
	}
	fmt.Println("Serving lampstand on : " + bind)
	log.Fatal(RunPassageService(service, bind))
}

func checkErr(err error) {
	if err != nil {
		fmt.Println(err)
		log.Println(err)
	}
}

func RunPassageService(s *PassageService, host string) error {
	router := httprouter.New()
	router.GET("/api/:version/verses", s.findVerses)
	router.NotFound = http.FileServer(http.Dir("static"))
	return http.ListenAndServe(host, router)
}
