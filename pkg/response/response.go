package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	appErr "github.com/edustack/go-boilerplate/pkg/errors"
	"github.com/gin-gonic/gin"
)

// APIResponse adalah struktur standar untuk semua response API.
type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// PaginatedResponse adalah struktur response untuk data yang dipaginasi.
type PaginatedResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    Pagination  `json:"meta"`
}

// Pagination berisi metadata paginasi.
type Pagination struct {
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	TotalData int64 `json:"total_data"`
	TotalPage int   `json:"total_page"`
}

// AuditEntry adalah data log yang dikirim ke AuditLogger.
// Struct ini didefinisikan di pkg/response agar tidak ada import cycle:
// pkg/response TIDAK mengimport internal/audit — cukup mendefinisikan kontrak datanya.
type AuditEntry struct {
	Date      string
	UserID    *string
	Name      *string
	Role      *string
	Host      string
	Path      string
	Method    string
	Status    int
	Data      json.RawMessage
	CreatedAt time.Time
}

// AuditLogger adalah interface yang harus diimplementasikan oleh audit.Repository.
// Dengan interface ini, pkg/response tidak perlu mengimport internal/audit secara langsung.
type AuditLogger interface {
	SaveAsync(entry AuditEntry)
}

// auditLogger adalah dependency global yang di-inject saat aplikasi startup.
var auditLogger AuditLogger

// SetAuditLogger meng-inject implementasi AuditLogger ke dalam package response.
// Dipanggil sekali saat aplikasi startup di main.go.
func SetAuditLogger(logger AuditLogger) {
	auditLogger = logger
}

// saveLog adalah helper internal yang membaca konteks Gin dan menyimpan log.
func saveLog(c *gin.Context, statusCode int, body interface{}) {
	if auditLogger == nil {
		return
	}

	var userID *string
	var name *string
	var role *string

	if uid, exists := c.Get("user_id"); exists {
		s := fmt.Sprintf("%v", uid)
		userID = &s
	}
	if n, exists := c.Get("name"); exists {
		s, _ := n.(string)
		name = &s
	}
	if r, exists := c.Get("role"); exists {
		s, _ := r.(string)
		role = &s
	}

	var rawData json.RawMessage
	if body != nil {
		if b, err := json.Marshal(body); err == nil {
			rawData = b
		}
	}

	entry := AuditEntry{
		Date:      time.Now().Format("2006-01-02"),
		UserID:    userID,
		Name:      name,
		Role:      role,
		Host:      c.Request.Host,
		Path:      c.Request.URL.Path,
		Method:    c.Request.Method,
		Status:    statusCode,
		Data:      rawData,
		CreatedAt: time.Now(),
	}

	auditLogger.SaveAsync(entry)
}

// Success mengirim response sukses dan menyimpan audit log.
func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	body := APIResponse{Success: true, Message: message, Data: data}
	c.JSON(statusCode, body)
	saveLog(c, statusCode, body)
}

// Error mengirim response error dan menyimpan audit log.
func Error(c *gin.Context, statusCode int, message string) {
	body := APIResponse{Success: false, Message: message}
	c.JSON(statusCode, body)
	saveLog(c, statusCode, body)
}

// WriteError mengkonversi AppError dan mengirim response error.
func WriteError(c *gin.Context, err error) {
	var appError *appErr.AppError
	if errors.As(err, &appError) {
		Error(c, appError.StatusCode, appError.Message)
		return
	}
	Error(c, http.StatusInternalServerError, "internal server error")
}

// Paginated mengirim response terpaginasi dan menyimpan audit log.
func Paginated(c *gin.Context, statusCode int, message string, data interface{}, page, pageSize int, totalData int64) {
	totalPage := int(totalData) / pageSize
	if int(totalData)%pageSize > 0 {
		totalPage++
	}

	body := PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta: Pagination{
			Page:      page,
			PageSize:  pageSize,
			TotalData: totalData,
			TotalPage: totalPage,
		},
	}
	c.JSON(statusCode, body)
	saveLog(c, statusCode, body)
}
