package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

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
		log.Printf("Failed to create inferer: %s", err)
		return
	}

	// setup route handlers
	resultsHandler := endpoints.GetResultsHandler(tempImgFolder, inferer)

	http.HandleFunc("/status", endpoints.StatusHandler)
	http.HandleFunc("/results", resultsHandler)

	// start server
	log.Printf("Listening on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
