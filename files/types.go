package files

import "github.com/google/uuid"

const (
	// ContentTypeTar is the MIME type for tar archives.
	ContentTypeTar = "application/x-tar"
	// ContentTypeZip is the MIME type for zip archives.
	ContentTypeZip = "application/zip"

	// FormatZip requests zip format conversion on download.
	FormatZip = "zip"
)

// UploadResponse contains the hash identifying an uploaded file.
type UploadResponse struct {
	ID uuid.UUID `json:"hash" format:"uuid"`
}
