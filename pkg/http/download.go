package http

import (
	"io"
	"net/http"
	"os"
)

func Download(f *os.File, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	_, err = io.Copy(f, resp.Body)
	return err
}
