package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/hekmon/plexhue/huecolors"
	"github.com/hekmon/plexwebhooks"
)

func mainHandler(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()
	// Create the multi part reader
	multiPartReader, err := r.MultipartReader()
	if err != nil {
		// Detect error type for the http answer
		if err == http.ErrNotMultipart || err == http.ErrMissingBoundary {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		// Try to write the error as http body
		_, wErr := w.Write([]byte(err.Error()))
		if wErr != nil {
			err = fmt.Errorf("request error: %v | write error: %v", err, wErr)
		}
		// Log the error
		fmt.Fprintf(os.Stderr, "can't create a multipart reader from request: %v\n", err)
		return
	}
	// Use the multipart reader to parse the request body
	payload, thumb, err := plexwebhooks.Extract(multiPartReader)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		// Try to write the error as http body
		_, wErr := w.Write([]byte(err.Error()))
		if wErr != nil {
			err = fmt.Errorf("request error: %v | write error: %v", err, wErr)
		}
		// Log the error
		fmt.Fprintf(os.Stderr, "can't create a multipart reader from request: %v\n", err)
		return
	}
	// Do something
	fmt.Println()
	fmt.Println(time.Now())
	if err == nil {
		fmt.Printf("%+v\n", *payload)
		if thumb != nil {
			fmt.Printf("Filename: %s | Size: %d\n", thumb.Filename, len(thumb.Data))
			colors, params, err := huecolors.GetHueColors(3, thumb.Data)
			if err == nil {
				fmt.Println(params)
				fmt.Println(colors)
			} else {
				fmt.Println(err)
			}
		}
	} else {
		fmt.Println(err)
	}
	fmt.Println()
}
