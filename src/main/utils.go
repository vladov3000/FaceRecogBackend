package main

import (
	"log"
	"os"
)

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
