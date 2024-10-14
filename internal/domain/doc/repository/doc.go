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
	query := `SELECT id, name, mime_type, file, public, location, created_at FROM documents WHERE id = $1`

	var doc models.Document
	err := r.db.QueryRowContext(ctx, query, documentID).Scan(
		&doc.MD.ID, &doc.MD.Name, &doc.MD.MimeType, &doc.MD.File, &doc.MD.Public, &doc.Location, &doc.MD.CreatedAt)

	if errors.Is(err, sql.ErrNoRows) {
		return nil, fmt.Errorf("document not found")
	}
	if err != nil {
		return nil, fmt.Errorf("error fetching document: %w", err)
	}

	return &doc, nil
}

func (r *DocumentRepository) GetAllDocuments(ctx context.Context, email string) (*[]models.Document, error) {
	query := `
		SELECT d.id, d.name, d.mime_type, d.file, d.public, d.location, d.created_at
		FROM documents d
		JOIN document_grants dg ON d.id = dg.document_id
		WHERE dg.user_email = $1
		ORDER BY d.created_at DESC;
	`
	rows, err := r.db.QueryContext(ctx, query, email)

	if err != nil {
		return nil, fmt.Errorf("error fetching documents: %w", err)
	}

	var docs []models.Document
	for rows.Next() {
		var doc models.Document
		if err := rows.Scan(&doc.MD.ID, &doc.MD.Name, &doc.MD.MimeType, &doc.MD.File, &doc.MD.Public, &doc.Location, &doc.MD.CreatedAt); err != nil {
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
	query := `INSERT INTO documents (name, mime_type, file, public, location, owner_email)
              VALUES ($1, $2, $3, $4, $5, $6) RETURNING id`
	err := dr.db.QueryRowContext(ctx, query,
		doc.MD.Name, doc.MD.MimeType, doc.MD.File, doc.MD.Public, doc.Location, doc.MD.Grant).Scan(&doc.MD.ID)
	if err != nil {
		return fmt.Errorf("error saving document: %w", err)
	}
	return nil
}

func (r *DocumentRepository) SaveDocumentGrants(ctx context.Context, documentID string, users []string) error {
	query := `INSERT INTO document_grants (document_id, user_email) VALUES ($1, $2)`
	for _, userEmail := range users {
		_, err := r.db.ExecContext(ctx, query, documentID, userEmail)
		if err != nil {
			return fmt.Errorf("error saving document grant: %w", err)
		}
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
