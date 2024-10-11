package repository

import (
	"bem/internal/domain/doc/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/exp/slog"
)

type DocumentRepository struct {
	db     *sql.DB
	logger *slog.Logger
}

func NewDocumentRepository(db *sql.DB, logger *slog.Logger) *DocumentRepository {
	return &DocumentRepository{
		db:     db,
		logger: logger,
	}
}

func (r *DocumentRepository) GetDocument(ctx context.Context, documentID string) (*models.Document, error) {
	query := `SELECT id, name, mime_type, file, public, location, token, owner_email, created_at
              FROM documents WHERE id = $1`

	var doc models.Document
	err := r.db.QueryRowContext(ctx, query, documentID).Scan(
		&doc.MD.ID, &doc.MD.Name, &doc.MD.MimeType, &doc.MD.File, &doc.MD.Public, &doc.Location, &doc.MD.Token, &doc.MD.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("document not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching document: %w", err)
	}

	return &doc, nil
}

// Функция сохранения документа
func (dr *DocumentRepository) SaveDocument(ctx context.Context, doc *models.Document) error {
	query := `INSERT INTO documents (name, mime_type, file, public, location, token, owner_email)
              VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	err := dr.db.QueryRowContext(ctx, query,
		doc.MD.Name, doc.MD.MimeType, doc.MD.File, doc.MD.Public, doc.Location, doc.MD.Token, doc.MD.Grant).Scan(&doc.MD.ID)
	if err != nil {
		return fmt.Errorf("error saving document: %w", err)
	}
	return nil
}
