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
	MaxSizeBytes: 100 << 20, // 100MB (audio/video need more headroom than the old 25MB)
	AllowedTypes: map[string]string{
		"image/png":       "png",
		"image/jpeg":      "jpg",
		"image/webp":      "webp",
		"image/gif":       "gif",
		"application/pdf": "pdf",
		"text/plain":      "txt",
		"audio/mpeg":      "mp3",
		"audio/wav":       "wav",
		"audio/ogg":       "ogg",
		"audio/webm":      "weba",
		"video/mp4":       "mp4",
		"video/webm":      "webm",
		"video/ogg":       "ogv",
		"video/quicktime": "mov",
		// http.DetectContentType can't tell an audio-only WebM/Ogg container
		// apart from a video one (same container signature), so a recorded
		// voice note (MediaRecorder's audio/webm or audio/ogg output) sniffs
		// as these instead. video/webm is already allowed above; Ogg needs
		// its own entry since the sniffer reports it as application/ogg.
		"application/ogg": "ogg",
	},
}

// ValidateOpaque checks size/emptiness only, for payloads that are
// ciphertext rather than a real media file (e.g. an encrypted voice note).
// Sniffing the content type of AES-GCM output is meaningless - it's
// indistinguishable from random bytes.
func ValidateOpaque(data []byte, maxSize int64) error {
	if len(data) == 0 {
		return ErrEmptyFile
	}
	if int64(len(data)) > maxSize {
		return fmt.Errorf("%w (max %d bytes)", ErrTooLarge, maxSize)
	}
	return nil
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
