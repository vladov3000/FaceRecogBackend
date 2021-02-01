package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/vladov3000/FaceRecogBackend/src/infer"
)

// maximum image that can be uploaded 10 << 20
// specifies a maximum upload of 10 MB files.
const MAX_IMG_SIZE = 10 << 20

type ReqHandler = func(w http.ResponseWriter, r *http.Request)

func getResultsHandler(tempImgFolder string, inferer infer.Inferer) ReqHandler {
	return func(w http.ResponseWriter, r *http.Request) {
		log.Print("Results Endpoint Hit")

		if r.Method != http.MethodPost {
			text := "405 - expected POST request"
			log.Print(text)

			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprint(w, text)
			return
		}

		// We are expecting an image (binary data) so multipart/form-data
		// is the most appropriate content type to parse.
		// Read more here: https://stackoverflow.com/questions/4007969/application-x-www-form-urlencoded-or-multipart-form-data

		contentType := r.Header.Get("content-type")
		if !strings.HasPrefix(contentType, "multipart/form-data") {
			text := "415 - expected multipart/form-data as prefix of content type, but got " + contentType
			log.Print(text)

			w.WriteHeader(http.StatusUnsupportedMediaType)
			fmt.Fprint(w, text)
			return
		}

		// Ensure content length is correct

		contentLength, err := strconv.Atoi(r.Header.Get("content-length"))
		if err != nil {
			text := "400 - failed to parse content length to int"
			log.Printf("Error: %s Response: %s", err, text)

			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, text)
			return
		}
		if contentLength > MAX_IMG_SIZE {
			text := fmt.Sprintf("413 - max content size of %d, got size of %d", MAX_IMG_SIZE, contentLength)
			log.Print(text)

			w.WriteHeader(http.StatusRequestEntityTooLarge)
			fmt.Fprint(w, text)
			return
		}

		// Get file from request

		r.ParseMultipartForm(MAX_IMG_SIZE)

		fileName := "imgFile"
		file, fileHeader, err := r.FormFile(fileName)
		if err != nil {
			text := fmt.Sprintf("404 - failed to get file %s from request", fileName)
			log.Printf("Error: %s Response: %s", err, text)

			w.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(w, text)
			return
		}
		defer file.Close()

		// Get bytes from file

		fileBytes, err := ioutil.ReadAll(file)
		if err != nil {
			text := "500 - failed to read file"
			log.Printf("Error: %s Response: %s", err, text)

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, text)
			return
		}

		// Get file extension

		fileExt := filepath.Ext(fileHeader.Filename)
		if fileExt != ".jpg" && fileExt != ".png" && fileExt != ".jpeg" {
			text := "422 - Expected file extension of .jpg/.png/.jpeg, got " + fileExt
			log.Print(text)

			w.WriteHeader(http.StatusUnprocessableEntity)
			fmt.Fprint(w, text)
			return
		}

		// Create temporary file

		tempFileName := "upload-*" + fileExt
		tempFile, err := ioutil.TempFile(tempImgFolder, tempFileName)
		if err != nil {
			text := "500 - failed to create temp file"
			log.Printf("Error for %s: %s Response: %s", tempFileName, err, text)

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, text)
			return
		}
		tempFileName = tempFile.Name()
		defer os.Remove(tempFileName)

		// write file bytes to temp file

		_, err = tempFile.Write(fileBytes)
		if err != nil {
			text := "500 - failed to write to temp file"
			log.Printf("Error: %s Response: %s", err, text)

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, text)
			return
		}

		// run inferer to get results
		result, err := inferer.GetResults(tempFileName)
		if err != nil {
			text := "500 - failed to run inferer"
			log.Printf("Error: %s Response: %s", err, text)

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, text)
			return
		}

		toSend, err := json.Marshal(result)
		if err != nil {
			text := "500 - failed to marshal into json"
			log.Printf("Error: %s Response: %s", err, text)

			w.WriteHeader(http.StatusInternalServerError)
			fmt.Fprint(w, text)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(toSend)
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	log.Print("Status Endpoint Hit")
	fmt.Fprint(w, "Online\n")
}

func createTempImgFolder(tempImgFolder string) {
	// check if folder exists
	_, err := os.Stat(tempImgFolder)
	if err == nil || !os.IsNotExist(err) {
		return
	}

	log.Printf("Creating temporary file folder at %s", tempImgFolder)

	err = os.Mkdir(tempImgFolder, 0777)
	if err != nil {
		log.Printf("Error when creating %s: %s", tempImgFolder, err)
	}
}

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
	resultsHandler := getResultsHandler(tempImgFolder, inferer)

	http.HandleFunc("/status", statusHandler)
	http.HandleFunc("/results", resultsHandler)

	// start server
	log.Printf("Listening on port %d", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}
