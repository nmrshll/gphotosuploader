package main

import (
	"fmt"
	"os"

	"github.com/simonedegiacomi/gphotosuploader/api"
	"github.com/simonedegiacomi/gphotosuploader/auth"
)

// Simple example which consist in the upload of a single image
func main() {
	// Load cookie for credentials from a json file
	credentials, err := auth.NewCookieCredentialsFromFile("auth.json")
	if err != nil {
		panic(err)
	}

	// Open the file to upload
	file, err := os.Open("path/to/image.png")
	if err != nil {
		panic(err)
	}

	// Create an UploadOptions object that describes the upload.
	options, err := api.NewUploadOptionsFromFile(file)
	if err != nil {
		panic(err)
	}

	// Create an upload using the NewUpload method from the api package
	upload, err := api.NewUpload(options, *credentials)
	if err != nil {
		panic(err)
	}

	// Finally upload the image
	uploadRes, err := upload.Upload()
	if err != nil {
		panic(err)
	}
	fmt.Println("success! image uploaded: ", uploadRes.URLString())

	// Image uploaded!
}
