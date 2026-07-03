package audit

import (
	"context"
	"log/slog"
	"time"

	"github.com/edustack/go-boilerplate/pkg/response"
	"gorm.io/gorm"
)

// QueryParams berisi filter dan paginasi untuk query list Logs.
type QueryParams struct {
	Page     int
	PageSize int
	Date     string // filter by date, format: "2006-01-02"
	Method   string // filter by HTTP method
	Status   int    // filter by HTTP status code
	Search   string // search by path atau host
}

// Repository menangani penulisan dan pembacaan audit Logs ke database.
type Repository struct {
	db *gorm.DB
}

// NewRepository membuat instance audit repository.
func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// SaveAsync mengimplementasikan response.AuditLogger.
func (r *Repository) SaveAsync(entry response.AuditEntry) {
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		logEntry := Logs{
			Date:      entry.Date,
			UserID:    entry.UserID,
			Name:      entry.Name,
			Role:      entry.Role,
			Host:      entry.Host,
			Path:      entry.Path,
			Method:    entry.Method,
			Status:    entry.Status,
			Data:      entry.Data,
			CreatedAt: entry.CreatedAt,
		}

		if err := r.db.WithContext(ctx).Create(&logEntry).Error; err != nil {
			slog.Error("audit: failed to save log", "error", err)
		}
	}()
}

// FindAll mengambil daftar Logs dengan filter dan paginasi.
func (r *Repository) FindAll(ctx context.Context, params QueryParams) ([]Logs, int64, error) {
	var logs []Logs
	var total int64

	query := r.db.WithContext(ctx).Model(&Logs{})

	if params.Date != "" {
		query = query.Where("date = ?", params.Date)
	}
	if params.Method != "" {
		query = query.Where("method = ?", params.Method)
	}
	if params.Status != 0 {
		query = query.Where("status = ?", params.Status)
	}
	if params.Search != "" {
		like := "%" + params.Search + "%"
		query = query.Where("path ILIKE ? OR host ILIKE ?", like, like)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.
		Order("created_at DESC").
		Limit(params.PageSize).
		Offset(offset).
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

// FindByID mengambil satu Logs berdasarkan ID-nya.
func (r *Repository) FindByID(ctx context.Context, id int64) (*Logs, error) {
	var logEntry Logs
	if err := r.db.WithContext(ctx).First(&logEntry, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &logEntry, nil
}
