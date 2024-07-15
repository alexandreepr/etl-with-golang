package handlers

import (
	"bytes"
	"etl-with-golang/helpers"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"io"
)

func TestImportFile(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("Successful file import", func(t *testing.T) {
		content := []byte("test content")
		tmpfile, err := os.CreateTemp("", "test*.txt")
		assert.NoError(t, err)
		defer os.Remove(tmpfile.Name())
		_, err = tmpfile.Write(content)
		assert.NoError(t, err)
		tmpfile.Close()

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, err := writer.CreateFormFile("file", "test.txt")
		assert.NoError(t, err)
		file, err := os.Open(tmpfile.Name())
		assert.NoError(t, err)
		_, err = io.Copy(part, file)
		assert.NoError(t, err)
		writer.Close()

		responseRecorder := httptest.NewRecorder()
		context, _ := gin.CreateTestContext(responseRecorder)
		context.Request, _ = http.NewRequest("POST", "/import-file", body)
		context.Request.Header.Set("Content-Type", writer.FormDataContentType())

		helpers.SetProcessFileWrapper(func(file multipart.File) (string, error) {
			return "mock-import-id", nil
		})

		ImportFile(context)

		assert.Equal(t, http.StatusOK, responseRecorder.Code)
		assert.Contains(t, responseRecorder.Body.String(), "File uploaded and data saved successfully")
	})

	t.Run("Invalid file extension", func(t *testing.T) {
		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		_, err := writer.CreateFormFile("file", "test.pdf")
		assert.NoError(t, err)
		writer.Close()

		responseRecorder := httptest.NewRecorder()
		context, _ := gin.CreateTestContext(responseRecorder)
		context.Request, _ = http.NewRequest("POST", "/import-file", body)
		context.Request.Header.Set("Content-Type", writer.FormDataContentType())

		ImportFile(context)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		assert.Contains(t, responseRecorder.Body.String(), "File extension not allowed")
	})

	t.Run("Failed to retrieve file", func(t *testing.T) {
		responseRecorder := httptest.NewRecorder()
		context, _ := gin.CreateTestContext(responseRecorder)
		context.Request, _ = http.NewRequest("POST", "/import-file", nil)

		ImportFile(context)

		assert.Equal(t, http.StatusBadRequest, responseRecorder.Code)
		assert.Contains(t, responseRecorder.Body.String(), "Failed to retrieve file")
	})
}
