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
	db = database.NewMongoDB("FaceRecogApp")
	defer db.Disconnect()

	// result, err := inferer.GetResults("test-images/obama.jpg")
	// if err != nil {
	// 	log.Fatal(err)
	// }
	//log.Printf("%+v", result)

	// person := database.NewPerson("1", "1", result.Encodings)
	//log.Printf("%+v", person)

	// err = db.AddPerson(person)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// personFromDB, err := db.GetPerson(bson.M{"encoding": person.Encoding})
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// log.Printf("%+v", personFromDB)

	// setup route handlers
	resultsHandler := endpoints.GetResultsHandler(tempImgFolder, inferer, db)

	http.HandleFunc("/status", endpoints.StatusHandler)
	http.HandleFunc("/results", resultsHandler)

	// start server
	log.Printf("Listening on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
