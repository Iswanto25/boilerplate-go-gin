package response

import (
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	appErr "github.com/edustack/go-boilerplate/pkg/errors"
	"github.com/edustack/go-boilerplate/pkg/logger"
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

// AuditEntry matches the Express.js respons.ts log payload.
type AuditEntry struct {
	Date      string
	UserID    *string
	Name      *string
	Role      *string
	Host      string // full URL: proto://host/path
	IP        string
	Method    string
	Status    string // stored as string (Express.js: code.toString())
	Data      json.RawMessage
	CreatedAt time.Time
}

// AuditLogger adalah interface yang harus diimplementasikan oleh audit.Repository.
type AuditLogger interface {
	SaveAsync(entry AuditEntry)
}

var auditLogger AuditLogger

func SetAuditLogger(logger AuditLogger) {
	auditLogger = logger
}

func saveLog(c *gin.Context, statusCode int, body interface{}) {
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

	// Compute request path
	path := c.Request.URL.Path
	if c.Request.URL.RawQuery != "" {
		path = path + "?" + c.Request.URL.RawQuery
	}

	// Console log: {METHOD} {PATH} {STATUS} | {userName} | {responseTime}ms
	userName := "Guest"
	if name != nil {
		userName = *name
	}
	startTime, _ := c.Get("startTime")
	if st, ok := startTime.(int64); ok {
		responseTime := time.Now().UnixMilli() - st
		logger.LogConsole(c.Request.Method, path, statusCode, userName, responseTime)
	} else {
		slog.Info(fmt.Sprintf("%s %s %d | %s", c.Request.Method, path, statusCode, userName))
	}

	// Save to database via audit logger
	if auditLogger == nil {
		return
	}

	// Build data payload matching Express.js format
	dateTimeNow := time.Now().Format("2006-01-02 15:04:05")
	forwardedFor := c.GetHeader("X-Forwarded-For")
	ip := forwardedFor
	if ip == "" {
		ip = c.ClientIP()
	}
	if ip == "" {
		ip = c.Request.RemoteAddr
	}
	userAgent := c.GetHeader("User-Agent")
	if userAgent == "" {
		userAgent = "Unknown"
	}

	statusStr := fmt.Sprintf("%d", statusCode)

	source := "Success"
	if statusCode >= 400 {
		source = "Error"
	}

	// Build data JSON matching Express.js log payload
	var dataPayload map[string]interface{}
	if body != nil {
		// The body is already the API response struct, extract the actual data
		if b, err := json.Marshal(body); err == nil {
			_ = json.Unmarshal(b, &dataPayload)
		}
	}

	logData := map[string]interface{}{
		"userAgent": userAgent,
		"timestamp": dateTimeNow,
		"source":    source,
		"message":   "",
		"data":      dataPayload,
	}

	rawData, _ := json.Marshal(logData)

	host := fmt.Sprintf("%s://%s%s", guessScheme(c), c.Request.Host, c.Request.URL.Path)

	dateTimeNowFull := time.Now().Format("2006-01-02 15:04:05")

	entry := AuditEntry{
		Date:      dateTimeNowFull,
		UserID:    userID,
		Name:      name,
		Role:      role,
		Host:      host,
		IP:        ip,
		Method:    c.Request.Method,
		Status:    statusStr,
		Data:      rawData,
		CreatedAt: time.Now(),
	}

	auditLogger.SaveAsync(entry)
}

func guessScheme(c *gin.Context) string {
	if proto := c.GetHeader("X-Forwarded-Proto"); proto != "" {
		return proto
	}
	if c.Request.TLS != nil {
		return "https"
	}
	return "http"
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
