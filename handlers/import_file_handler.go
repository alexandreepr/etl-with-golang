package handlers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"etl-with-golang/helpers"
	"io"
	"os"
)

// ImportFile handles the file upload for data import.
// It accepts a .txt file via multipart/form-data, processes its contents,
// and returns an import ID upon successful processing.
//
// POST /api/v1/file-import
//
// Request:
//   - Content-Type: multipart/form-data
//   - Body: form field "file" containing a .txt file
//
// Responses:
//   - 200 OK: File uploaded and processed successfully
//     {
//       "message": "File uploaded and data saved successfully.",
//       "importacao_id": "<uuid>"
//     }
//   - 400 Bad Request: Invalid file or request
//   - 500 Internal Server Error: Server-side processing error
func ImportFile(ctx *gin.Context) {
	file, header, err := ctx.Request.FormFile("file")
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Failed to retrieve file"})
		return
	}
	defer file.Close()

	if filepath.Ext(header.Filename) != ".txt" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "File extension not allowed"})
		return
	}

	tempFile, err := os.CreateTemp("", "upload-*.txt")
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create temporary file"})
		return
	}
	defer os.Remove(tempFile.Name())

	// Copy the uploaded file content to the temporary file
	if _, err := io.Copy(tempFile, file); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
		return
	}

	// Close the temporary file to ensure all data is flushed to disk
	tempFile.Close()

	tempFile, err = os.Open(tempFile.Name())
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reopen temporary file"})
		return
	}
	defer tempFile.Close()

	importId, err := helpers.ProcessFile(tempFile)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "File uploaded and data saved successfully.", "importacao_id": importId})
}
