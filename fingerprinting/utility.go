package fingerprinting

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"gitlab.com/opennota/screengen"
)

//DownloadFile downloads a file using a GET http request
func downloadFile(url, filePath string) (err error) {

	// Create the file in the folder
	out, err := os.Create(filePath)
	if err != nil {
		return
	}
	defer closeSafely(out)

	// Get the data
	resp, err := http.Get(url) // nolint: gosec
	if err != nil {
		return
	}
	defer closeSafely(resp.Body)

	// Check server response
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status (%s) when downloading the file %s", resp.Status, url)
	}

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return
}

// closeSafely closes an entity and logs in case of errors
func closeSafely(toClose io.Closer) {
	err := toClose.Close()
	if err != nil {
		log.Println(err)
	}
}

// removeSafely calls os.Remove and logs in case of errors
func removeSafely(pathToFile string) {
	err := os.Remove(pathToFile)
	if err != nil {
		log.Println(err)
	}
}

// closeGeneratorSafely closes a screengen.Generator and logs in case of errors
func closeGeneratorSafely(g *screengen.Generator) {
	err := g.Close()
	if err != nil {
		log.Println(err)
	}
}
