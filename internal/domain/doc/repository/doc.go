package repository

import (
	"bem/internal/domain/doc/models"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"golang.org/x/exp/slog"
	"os"
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
	query := `SELECT id, name, mime_type, file, public, location, created_at, token FROM documents WHERE id = $1`

	var doc models.Document
	err := r.db.QueryRowContext(ctx, query, documentID).Scan(
		&doc.MD.ID, &doc.MD.Name, &doc.MD.MimeType, &doc.MD.File, &doc.MD.Public, &doc.Location, &doc.MD.CreatedAt, &doc.MD.Token)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("document not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching document: %w", err)
	}

	return &doc, nil
}

func (r *DocumentRepository) GetAllDocuments(ctx context.Context) (*[]models.Document, error) {
	query := `SELECT id, name, mime_type, file, public, location, created_at, token FROM documents`

	var docs []models.Document
	rows, err := r.db.QueryContext(ctx, query)

	if err != nil {
		return nil, fmt.Errorf("error fetching documents: %w", err)
	}

	for rows.Next() {
		var doc models.Document
		if err := rows.Scan(&doc.MD.ID, &doc.MD.Name, &doc.MD.MimeType,
			&doc.MD.File, &doc.MD.Public, &doc.Location, &doc.MD.CreatedAt, &doc.MD.Token); err != nil {
			return nil, fmt.Errorf("error scanning document: %w", err)
		}
		docs = append(docs, doc)
	}
	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("document not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching document: %w", err)
	}

	return &docs, nil
}

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

func (r *DocumentRepository) DeleteDocument(ctx context.Context, documentID string) error {
	// Удаляем документ из базы данных и получаем путь к файлу
	var location string
	deleteQuery := `DELETE FROM documents WHERE id = $1 RETURNING location`
	err := r.db.QueryRowContext(ctx, deleteQuery, documentID).Scan(&location)
	fmt.Printf("location %s \n", location)
	if err == sql.ErrNoRows {
		return fmt.Errorf("document not found")
	}
	if err != nil {
		return fmt.Errorf("error deleting document from database: %w", err)
	}

	// Удаляем файл с диска
	if err = os.Remove(location); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error deleting file: %w", err)
	}

	return nil
}
