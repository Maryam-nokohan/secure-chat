package fileupload

import (
	"errors"
	"fmt"
	"net/http"
)

type Rule struct {

	AllowedTypes map[string]string
	MaxSizeBytes int64
}

var (
	ErrTooLarge       = errors.New("file exceeds the maximum allowed size")
	ErrEmptyFile      = errors.New("file is empty")
	ErrDisallowedType = errors.New("file type is not allowed")
)


var AvatarRule = Rule{
	MaxSizeBytes: 5 << 20, // 5MB
	AllowedTypes: map[string]string{
		"image/png":  "png",
		"image/jpeg": "jpg",
		"image/webp": "webp",
	},
}

var ChatAttachmentRule = Rule{
	MaxSizeBytes: 25 << 20, // 25MB
	AllowedTypes: map[string]string{
		"image/png":       "png",
		"image/jpeg":      "jpg",
		"image/webp":      "webp",
		"image/gif":       "gif",
		"application/pdf": "pdf",
		"text/plain":      "txt",
	},
}

func Validate(data []byte, rule Rule) (mimeType, ext string, err error) {
	if len(data) == 0 {
		return "", "", ErrEmptyFile
	}
	if int64(len(data)) > rule.MaxSizeBytes {
		return "", "", fmt.Errorf("%w (max %d bytes)", ErrTooLarge, rule.MaxSizeBytes)
	}
	sniff := data
	if len(sniff) > 512 {
		sniff = sniff[:512]
	}
	detected := http.DetectContentType(sniff)
	ext, ok := rule.AllowedTypes[detected]
	if !ok {
		return "", "", fmt.Errorf("%w: %s", ErrDisallowedType, detected)
	}
	return detected, ext, nil
}