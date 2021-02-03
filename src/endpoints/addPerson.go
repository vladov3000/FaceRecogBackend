package endpoints

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/vladov3000/FaceRecogBackend/src/database"
)

func AddPersonHandler(db database.Database) ReqHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("Add Person Endpoint Hit")

		if r.Method != http.MethodPost {
			text := "405 - expected POST request"
			log.Print(text)

			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w, text)
			return
		}

		// parse person from request body
		var person database.Person

		decoder := json.NewDecoder(r.Body)
		err := decoder.Decode(&person)
		if err != nil {
			text := "400 - couldn't parse JSON in request body"
			log.Printf("Error: %s Response: %s", err, text)

			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, text)
			return
		}

		// add person to database
		if err = db.AddPerson(person); err != nil {
			text := fmt.Sprintf("500 - failed to add person %+v to database", person)
			log.Printf("Error: %s Response: %s", err, text)

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, text)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}
