package endpoints

import (
	"fmt"
	"log"
	"net/http"
)

func statusHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Status Endpoint Hit")
	fmt.Fprint(w, "Online\n")
}
