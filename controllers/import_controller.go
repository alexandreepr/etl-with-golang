package controllers

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"path/filepath"
	"etl-with-golang/helpers"
	"io"
	"os"
)

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

	batchChan, wg := helpers.SetupWorkers()

	if err := helpers.ProcessFile(tempFile, batchChan); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	close(batchChan)
	wg.Wait()

	ctx.JSON(http.StatusOK, gin.H{"message": "File uploaded and data saved successfully"})
}
