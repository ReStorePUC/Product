package handler

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type File struct{}

func NewFile() *File {
	return &File{}
}

//// UploadFile is a function that handles the upload of a single file
//func (c *FileController) UploadFile(ctx *gin.Context) {
//	/*
//	  UploadFile function handles the upload of a single file.
//	  It gets the file from the form data, saves it to the defined path,
//	  generates a unique identifier for the file, saves the file metadata to the database,
//	  and returns a success message and the file metadata.
//	*/
//	// Get the file from the form data
//	file, err := ctx.FormFile("file")
//	if err != nil {
//		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
//		return
//	}
//	// Define the path where the file will be saved
//	filePath := filepath.Join("uploads", file.Filename)
//	// Save the file to the defined path
//	if err := ctx.SaveUploadedFile(file, filePath); err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
//		return
//	}
//	// Generate a unique identifier for the file
//	uuid := uuid.New().String()
//	// Save file metadata to database
//	fileMetadata := models.File{
//		Filename: file.Filename,
//		UUID:     uuid,
//	}
//	if err := c.DB.Create(&fileMetadata).Error; err != nil {
//		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file metadata"})
//		return
//	}
//	// Return a success message and the file metadata
//	ctx.JSON(http.StatusOK, gin.H{"message": "File uploaded successfully", "Details": fileMetadata})
//}

// UploadFiles is a function that handles the upload of multiple files
func (f *File) UploadFiles(c *gin.Context) {
	form, err := c.MultipartForm()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	files := form.File["files"]
	fileNames := []string{}
	for _, file := range files {
		fileName := uuid.New().String()
		filePath := filepath.Join("uploads", fileName+".png")
		fileNames = append(fileNames, fileName)

		if err := c.SaveUploadedFile(file, filePath); err != nil {
			c.IndentedJSON(http.StatusBadRequest, struct {
				Error string
			}{
				err.Error(),
			})
			return
		}
	}

	c.IndentedJSON(http.StatusCreated, struct {
		Files []string
	}{
		fileNames,
	})
}

// GetFile is a function that retrieves a file from the server
func (f *File) GetFile(c *gin.Context) {
	fileName := c.Param("file")

	// Define the path of the file to be retrieved
	filePath := filepath.Join("uploads", fileName+".png")
	// Open the file
	fileData, err := os.Open(filePath)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}
	defer fileData.Close()

	// Read the first 512 bytes of the file to determine its content type
	fileHeader := make([]byte, 512)
	_, err = fileData.Read(fileHeader)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}
	fileContentType := http.DetectContentType(fileHeader)
	// Get the file info
	fileInfo, err := fileData.Stat()
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}
	// Set the headers for the file transfer and return the file
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", fileName))
	c.Header("Content-Type", fileContentType)
	c.Header("Content-Length", fmt.Sprintf("%d", fileInfo.Size()))
	c.File(filePath)
}

// DeleteFile is a function that deletes a file from the server and its metadata from the database
func (f *File) DeleteFile(c *gin.Context) {
	fileName := c.Param("file")

	// Define the path of the file to be deleted
	filePath := filepath.Join("uploads", fileName+".png")
	// Delete the file from the server
	err := os.Remove(filePath)
	if err != nil {
		c.IndentedJSON(http.StatusBadRequest, struct {
			Error string
		}{
			err.Error(),
		})
		return
	}

	c.IndentedJSON(http.StatusOK, struct{}{})
}
