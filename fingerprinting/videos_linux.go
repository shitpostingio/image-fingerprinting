package fingerprinting

import (
	"fmt"
	"image/png"
	"log"
	"os"
	"strings"

	"gitlab.com/opennota/screengen"
)

// getPhotoFingerprint downloads and fingerprints a video.
func getVideoFingerprint(url, folderPath string) (filePath, aHash, pHash string, histogram []float64, err error) {

	filePath, err = downloadVideo(url, folderPath)
	if err != nil {
		return
	}

	aHash, pHash, histogram, err = FingerprintVideo(filePath)
	return

}

// FingerprintVideo fingerprints a video given its path.
func FingerprintVideo(fileName string) (aHash string, pHash string, histogram []float64, err error) {

	//Initialization of the screen generator
	generator, err := screengen.NewGenerator(fileName)
	if err != nil {
		log.Println(fmt.Sprintf("Unable to create a screenshot generator for file %s", fileName))
		return
	}
	defer closeGeneratorSafely(generator)

	//We want to grab a screenshot in the middle of the video
	framePath := fileName + ".png"
	err = extractFrame(framePath, generator.Duration/2, generator)
	if err != nil {
		return
	}

	aHash, pHash, histogram, err = FingerprintPhoto(framePath)
	return

}

//extractFrame extracts a frame from a video file
func extractFrame(frameOutputName string, timestamp int64, generator *screengen.Generator) (err error) {

	frame, err := generator.Image(timestamp)
	if err != nil {
		log.Println("Unable to screengrab from file")
		return
	}

	outputFile, err := os.Create(frameOutputName)
	if err != nil {
		log.Println("Unable to save frame on disk")
		return
	}
	defer closeSafely(outputFile)

	err = png.Encode(outputFile, frame)
	if err != nil {
		log.Println("Unable to encode frame")
	}

	return
}

// downloadVideo gets the appropriate name for the video and downloads it
// in the folderPath provided.
func downloadVideo(url, folderPath string) (filePath string, err error) {

	filePath = fmt.Sprintf("%s/%s", folderPath, getVideoName(url))
	err = downloadFile(url, filePath)
	return

}

// getVideoName adds the MP4 extension to telegram videos, so we
// can proceed to the fingerprinting.
func getVideoName(url string) string {

	fileNameStart := strings.LastIndex(url, "/") + 1
	fileName := url[fileNameStart:]
	return fileName + ".mp4"

}
