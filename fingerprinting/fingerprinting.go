package fingerprinting

import (
	"errors"
)

//nolint
const (
	PHOTO            = "AgA"
	VIDEO            = "BAA"
	ANIMATION        = "CgA"
	STICKER          = "CAA"
	VOICE            = "AwA"
	DOCUMENT         = "BQA"
	AUDIO            = "CQA"
	VIDEONOTE        = "DQA"
)

func GetFingerprint(fileID, url, folderPath string) (filePath, aHash, pHash string, histogram []float64, err error) {

	fileIDPrefix := fileID[:3]

	switch fileIDPrefix {
	case STICKER:
		return getStickerFingerprint(url, folderPath)
	case PHOTO:
		return getPhotoFingerprint(url, folderPath)
	case VIDEO:
		return getVideoFingerprint(url, folderPath)
	case ANIMATION:
		return getVideoFingerprint(url, folderPath)
	}

	return "", "", "", []float64{}, errors.New("media type not supported")

}
