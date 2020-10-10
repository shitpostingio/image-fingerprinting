package fingerprinting

import (
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"os"
	"strings"

	"github.com/AlessandroPomponio/hsv/histogram"

	"github.com/corona10/goimagehash"
)

// getPhotoFingerprint downloads and fingerprints a photo.
func getPhotoFingerprint(url, folderPath string) (filePath, aHash, pHash string, histogram []float64, err error) {

	filePath, err = downloadPhoto(url, folderPath)
	if err != nil {
		return
	}

	aHash, pHash, histogram, err = FingerprintPhoto(filePath)
	return

}

// FingerprintPhoto fingerprints a photo given its path.
func FingerprintPhoto(filePath string) (aHash string, pHash string, retHistogram []float64, err error) {

	file, err := os.Open(filePath) // nolint: gosec
	if err != nil {
		log.Println("Unable to open file ", filePath, err)
		return
	}
	defer closeSafely(file)

	//Telegram doesn't tell us what the content-type is when
	//downloading a picture, so we'll try to decode it both
	//as a jpeg and a png file
	var img image.Image
	if strings.HasSuffix(filePath, ".jpg") {
		img, err = jpeg.Decode(file)
		if err != nil {
			log.Println("Unable to decode file ", filePath, err)
			return
		}
	} else if strings.HasSuffix(filePath, ".png") {
		img, err = png.Decode(file)
		if err != nil {
			log.Println("Unable to decode file ", filePath, err)
			return
		}
	} else {
		log.Println("Unsupported file extension")
		return
	}

	//We want both AverageHash and PerceptionHash
	hash, err := goimagehash.AverageHash(img)
	if err != nil {
		log.Println("Unable to get average hash for file: ", filePath, err)
		return
	}
	aHash = hash.ToString()

	hash, err = goimagehash.PerceptionHash(img)
	if err != nil {
		log.Println("Unable to get perception hash for file: ", filePath, err)
		return
	}
	pHash = hash.ToString()

	retHistogram = histogram.With32BinsConcurrent(img, histogram.RoundClosest)

	return

}

// downloadPhoto gets the appropriate name for the photo and downloads it
// in the folderPath provided.
func downloadPhoto(url, folderPath string) (filePath string, err error) {

	filePath = fmt.Sprintf("%s/%s", folderPath, getPhotoName(url))
	err = downloadFile(url, filePath)
	return

}

// getPhotoName returns the name of the photo, making sure
// it has an extension.
func getPhotoName(url string) string {

	fileNameStart := strings.LastIndex(url, "/") + 1
	fileName := url[fileNameStart:]

	if !strings.HasSuffix(fileName, ".jpg") && !strings.HasSuffix(fileName, ".png") {
		fileName = fileName + ".jpg"
	}

	return fileName

}
