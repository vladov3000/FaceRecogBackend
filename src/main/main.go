package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/vladov3000/FaceRecogBackend/src/database"
	"github.com/vladov3000/FaceRecogBackend/src/endpoints"
	"github.com/vladov3000/FaceRecogBackend/src/infer"
)

func main() {
	// parse command line arguments
	portPtr := flag.Int("p", 8000, "port to listen on")
	tempImgFolderPtr := flag.String("i", "image-uploads", "folder for temporary image uploads")
	flag.Parse()

	port := *portPtr
	tempImgFolder := *tempImgFolderPtr

	// setup folder
	createTempImgFolder(tempImgFolder)

	// create inferer
	model := "build/model/model.dat"
	sp := "build/model/shape_predictor.dat"
	inferer, err := infer.NewInferer(model, sp)
	if err != nil {
		log.Fatalf("Failed to create inferer: %s", err)
	}
	defer inferer.Free()

	// create database
	var db database.Database
	db = database.NewMongoDB("FaceRecogDB")
	defer db.Disconnect()

	// result, err := inferer.GetResults("test-images/obama.jpg")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Printf("%+v", result)

	// person := database.NewPerson("b", "obama", result.Encodings, []string{"hello world", "former president"})
	// log.Printf("%+v", person)

	// err = db.AddPerson(person)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// personFromDB, _, err := db.GetPerson(bson.M{"encoding": person.Encoding})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Printf("%+v", personFromDB)

	// setup route handlers
	getPeopleHandler := endpoints.GetPeopleHandler(tempImgFolder, inferer, db)
	addPersonHandler := endpoints.AddPersonHandler(db)

	http.Handle("/", http.FileServer(http.Dir("./static")))
	http.HandleFunc("/status", endpoints.StatusHandler)
	http.HandleFunc("/getPeople", getPeopleHandler)
	http.HandleFunc("/addPerson", addPersonHandler)

	// start server
	log.Printf("Listening on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
