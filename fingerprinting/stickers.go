package fingerprinting

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/AlessandroPomponio/hsv/histogram"

	"github.com/corona10/goimagehash"
	"golang.org/x/image/webp"
)

// getStickerFingerprint downloads and fingerprints a sticker.
func getStickerFingerprint(url, folderPath string) (filePath, aHash, pHash string, histogram []float64, err error) {

	if stickerIsAnimated(url) {
		err = errors.New("animated stickers are not supported")
		return
	}

	filePath, err = downloadSticker(url, folderPath)
	if err != nil {
		return
	}

	pHash, histogram, err = FingerprintSticker(filePath)
	return

}

// FingerprintSticker fingerprints a sticker given its path.
func FingerprintSticker(filePath string) (pHash string, retHistogram []float64, err error) {

	if !strings.HasSuffix(filePath, ".webp") {
		log.Println("Unsupported file extension")
		return
	}

	file, err := os.Open(filePath) // nolint: gosec
	if err != nil {
		log.Println("Unable to open file ", filePath, err)
		return
	}
	defer closeSafely(file)

	img, err := webp.Decode(file)
	if err != nil {
		log.Println("Unable to decode file ", filePath, err)
		return
	}

	hash, err := goimagehash.PerceptionHash(img)
	if err != nil {
		log.Println("Unable to get perception hash for file: ", filePath, err)
		return
	}
	pHash = hash.ToString()

	retHistogram = histogram.With32BinsConcurrent(img, histogram.RoundClosest)

	return

}

func stickerIsAnimated(url string) bool {
	return strings.HasSuffix(url, "tgs")
}

// downloadSticker gets the appropriate name for the sticker and downloads it
// in the folderPath provided.
func downloadSticker(url, folderPath string) (filePath string, err error) {

	filePath = fmt.Sprintf("%s/%s", folderPath, getStickerName(url))
	err = downloadFile(url, filePath)
	return

}

// getStickerName returns the name of the sticker.
func getStickerName(url string) string {

	fileNameStart := strings.LastIndex(url, "/") + 1
	fileName := url[fileNameStart:]

	if !strings.HasSuffix(fileName, ".webp") && !strings.HasSuffix(fileName, ".tgs") {
		fileName = fileName + ".webp"
	}

	return fileName

}

