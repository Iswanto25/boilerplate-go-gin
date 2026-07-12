package audit

import (
	"context"
	"log/slog"
	"time"

	"github.com/edustack/go-boilerplate/pkg/response"
	"gorm.io/gorm"
)

const (
	auditWorkerCount  = 10
	auditChannelSize  = 1000
	auditWriteTimeout = 5 * time.Second
)

type QueryParams struct {
	Page     int
	PageSize int
	Date     string
	Method   string
	Status   string
	Search   string
}

type Repository struct {
	db     *gorm.DB
	logCh  chan Logs
}

func NewRepository(db *gorm.DB) *Repository {
	r := &Repository{
		db:    db,
		logCh: make(chan Logs, auditChannelSize),
	}
	for i := 0; i < auditWorkerCount; i++ {
		go r.logWorker()
	}
	return r
}

func (r *Repository) logWorker() {
	for entry := range r.logCh {
		ctx, cancel := context.WithTimeout(context.Background(), auditWriteTimeout)
		if err := r.db.WithContext(ctx).Create(&entry).Error; err != nil {
			slog.Error("audit: failed to save log", "error", err)
		}
		cancel()
	}
}

func (r *Repository) SaveAsync(entry response.AuditEntry) {
	logEntry := Logs{
		Date:      &entry.Date,
		Name:      entry.Name,
		Role:      entry.Role,
		Host:      &entry.Host,
		Status:    &entry.Status,
		Data:      entry.Data,
		UserID:    entry.UserID,
		IP:        &entry.IP,
		Method:    &entry.Method,
		CreatedAt: entry.CreatedAt,
	}

	select {
	case r.logCh <- logEntry:
	default:
		slog.Warn("audit: log channel full, dropping log entry")
	}
}

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
	if params.Status != "" {
		query = query.Where("status = ?", params.Status)
	}
	if params.Search != "" {
		like := "%" + params.Search + "%"
		query = query.Where("host ILIKE ?", like)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (params.Page - 1) * params.PageSize
	if err := query.
		Order("createdAt DESC").
		Limit(params.PageSize).
		Offset(offset).
		Find(&logs).Error; err != nil {
		return nil, 0, err
	}

	return logs, total, nil
}

func (r *Repository) FindByID(ctx context.Context, id int64) (*Logs, error) {
	var logEntry Logs
	if err := r.db.WithContext(ctx).First(&logEntry, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &logEntry, nil
}
