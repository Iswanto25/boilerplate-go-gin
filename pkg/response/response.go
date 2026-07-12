package response

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"regexp"
	"strings"
	"time"

	appErr "github.com/edustack/go-boilerplate/pkg/errors"
	"github.com/edustack/go-boilerplate/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

type PaginatedResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	Meta    Pagination  `json:"meta"`
}

type Pagination struct {
	Page      int   `json:"page"`
	PageSize  int   `json:"page_size"`
	TotalData int64 `json:"total_data"`
	TotalPage int   `json:"total_page"`
}

type AuditEntry struct {
	Date      string
	UserID    *string
	Name      *string
	Role      *string
	Host      string
	IP        string
	Method    string
	Status    string
	Data      json.RawMessage
	CreatedAt time.Time
}

type AuditLogger interface {
	SaveAsync(entry AuditEntry)
}

var auditLogger AuditLogger

func SetAuditLogger(l AuditLogger) {
	auditLogger = l
}

var sensitiveKeys = regexp.MustCompile(`(?i)"(password|accessToken|refreshToken|token|secret)"\s*:\s*"[^"]*"`)

func maskSensitive(raw []byte) []byte {
	return sensitiveKeys.ReplaceAllFunc(raw, func(match []byte) []byte {
		prefix := regexp.MustCompile(`("[^"]*"\s*:\s*)`).Find(match)
		return append(prefix, []byte(`"***"`)...)
	})
}

func captureBody(c *gin.Context) json.RawMessage {
	if c.Request.Body == nil || c.Request.ContentLength == 0 {
		return nil
	}
	raw, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil
	}
	c.Request.Body = io.NopCloser(bytes.NewBuffer(raw))

	var parsed map[string]interface{}
	if json.Unmarshal(raw, &parsed) != nil {
		return nil
	}
	masked := maskSensitive(raw)

	var result json.RawMessage
	json.Unmarshal(masked, &result)
	return result
}

func CaptureRequestBody() gin.HandlerFunc {
	return func(c *gin.Context) {
		body := captureBody(c)
		if body != nil {
			c.Set("requestBody", body)
		}
		c.Next()
	}
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

	path := c.Request.URL.Path
	if c.Request.URL.RawQuery != "" {
		path = path + "?" + c.Request.URL.RawQuery
	}

	userName := "Guest"
	if name != nil {
		userName = *name
	}
	startTime, _ := c.Get("startTime")
	if st, ok := startTime.(int64); ok {
		responseTime := time.Now().UnixMilli() - st
		logger.LogConsole(c.Request.Method, path, statusCode, userName, responseTime)
	} else {
		slog.Info(c.Request.Method + " " + path)
	}

	if auditLogger == nil {
		return
	}

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

	reqID := uuid.New().String()

	source := "Success"
	if statusCode >= 400 {
		source = "Error"
	}

	var reqBody json.RawMessage
	if rb, exists := c.Get("requestBody"); exists {
		reqBody, _ = rb.(json.RawMessage)
	}
	if reqBody == nil {
		reqBody = json.RawMessage("{}")
	}

	query := make(map[string]string)
	for k, v := range c.Request.URL.Query() {
		query[k] = strings.Join(v, ",")
	}
	queryJSON, _ := json.Marshal(query)

	params := make(map[string]string)
	for _, p := range c.Params {
		params[p.Key] = p.Value
	}
	paramsJSON, _ := json.Marshal(params)

	requestPayload := map[string]interface{}{
		"body":   json.RawMessage(reqBody),
		"query":  json.RawMessage(queryJSON),
		"reqId":  reqID,
		"action": "AuditLog",
		"params": json.RawMessage(paramsJSON),
		"userId": userID,
	}

	var respData interface{}
	if body != nil {
		b, _ := json.Marshal(body)
		masked := maskSensitive(b)
		json.Unmarshal(masked, &respData)
	}

	dateTimeStr := time.Now().Format("2006-01-02 15:04:05")

	responsePayload := map[string]interface{}{
		"data":      respData,
		"reqId":     reqID,
		"source":    source,
		"userId":    userID,
		"message":   "",
		"timestamp": dateTimeStr,
		"userAgent": userAgent,
	}

	logData := map[string]interface{}{
		"request":  requestPayload,
		"response": responsePayload,
	}

	rawData, _ := json.Marshal(logData)

	host := fmt.Sprintf("%s://%s%s", guessScheme(c), c.Request.Host, c.Request.URL.Path)

	entry := AuditEntry{
		Date:      dateTimeStr,
		UserID:    userID,
		Name:      name,
		Role:      role,
		Host:      host,
		IP:        ip,
		Method:    c.Request.Method,
		Status:    fmt.Sprintf("%d", statusCode),
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

func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	body := APIResponse{Success: true, Message: message, Data: data}
	c.JSON(statusCode, body)
	saveLog(c, statusCode, body)
}

func Error(c *gin.Context, statusCode int, message string) {
	body := APIResponse{Success: false, Message: message}
	c.JSON(statusCode, body)
	saveLog(c, statusCode, body)
}

func WriteError(c *gin.Context, err error) {
	var appError *appErr.AppError
	if errors.As(err, &appError) {
		Error(c, appError.StatusCode, appError.Message)
		return
	}
	Error(c, http.StatusInternalServerError, "internal server error")
}

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
