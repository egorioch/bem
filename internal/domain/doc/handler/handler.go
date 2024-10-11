package handler

import (
	"bem/internal/domain/doc/models"
	"bem/internal/domain/doc/service"
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/slog"
	"io/ioutil"
	"net/http"
)

type DocumentHandler struct {
	documentService *service.DocumentService
	logger          *slog.Logger
}

func NewDocumentHandler(ds *service.DocumentService, logger *slog.Logger) *DocumentHandler {
	return &DocumentHandler{
		documentService: ds,
		logger:          logger,
	}
}

func (dh *DocumentHandler) GetDocumentHandler(c *gin.Context) {
	documentID := c.Param("id") // Получаем ID документа из URL

	doc, err := dh.documentService.GetDocument(c.Request.Context(), documentID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Document not found"})
		return
	}

	c.JSON(http.StatusOK, doc.MD)
}

func (dh *DocumentHandler) SaveDocumentHandler(c *gin.Context) {
	var doc models.Document
	form, err := c.MultipartForm()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid form data"})
		return
	}

	metaData := form.Value["meta"]
	if err := json.Unmarshal([]byte(metaData[0]), &doc.MD); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid meta data"})
		return
	}

	file, _, err := c.Request.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File read error"})
		return
	}

	// Сохранение файла на диск
	filePath := fmt.Sprintf("static/uploads/%s", doc.MD.Name)
	if err := ioutil.WriteFile(filePath, fileBytes, 0644); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File write error"})
		return
	}
	doc.Location = filePath

	if err = dh.documentService.SaveDocument(c.Request.Context(), &doc); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "File save error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "File uploaded successfully"})
}
