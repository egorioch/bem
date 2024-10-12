package service

import (
	"bem/internal/domain/doc/models"
	"bem/internal/domain/doc/repository"
	"context"
	"fmt"
	"golang.org/x/exp/slog"
)

type DocumentService struct {
	DocumentRepo *repository.DocumentRepository
	Logger       *slog.Logger
	Cache        *Cache
}

func CreateNewDocumentService(dr *repository.DocumentRepository, logger *slog.Logger, cache *Cache) *DocumentService {
	return &DocumentService{
		DocumentRepo: dr,
		Logger:       logger,
		Cache:        cache,
	}
}

func (ds *DocumentService) GetDocument(ctx context.Context, documentID string) (*models.Document, error) {
	doc, err := ds.DocumentRepo.GetDocument(ctx, documentID)

	return doc, err
}

func (ds *DocumentService) GetAllDocuments(ctx context.Context) (*[]models.Document, error) {
	docs, err := ds.DocumentRepo.GetAllDocuments(ctx)

	return docs, err
}

func (ds *DocumentService) DeleteDocumentByID(ctx context.Context, documentID string) error {
	err := ds.DocumentRepo.DeleteDocument(ctx, documentID)
	if err != nil {
		return fmt.Errorf("document not found or could not be deleted")
	}

	return nil
}

func (ds *DocumentService) SaveDocument(ctx context.Context, doc *models.Document) error {
	err := ds.SaveDocument(ctx, doc)

	return err
}
