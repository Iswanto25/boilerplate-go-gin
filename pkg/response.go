package pkg

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type APIResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Hint    string      `json:"hint,omitempty"`
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
		LogConsole(c.Request.Method, path, statusCode, userName, responseTime)
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

	var reqBody json.RawMessage
	if rb, exists := c.Get("requestBody"); exists {
		reqBody, _ = rb.(json.RawMessage)
	}

	reqID := uuid.New().String()
	now := time.Now()
	timeISO := now.Format("2006-01-02T15:04:05.000-07:00")

	durationMs := int64(0)
	if st, ok := startTime.(int64); ok {
		durationMs = now.UnixMilli() - st
	}

	reqTime := timeISO
	resTime := timeISO

	var requestData interface{}
	if c.Request.Method == "GET" || c.Request.Method == "DELETE" {
		queryMap := make(map[string]string)
		for k, v := range c.Request.URL.Query() {
			queryMap[k] = strings.Join(v, ",")
		}
		queryMap["reqTime"] = reqTime
		requestData = map[string]interface{}{
			"queryParams": queryMap,
		}
	} else {
		reqMap := make(map[string]interface{})
		if len(reqBody) > 0 {
			json.Unmarshal(reqBody, &reqMap)
		}
		if reqMap == nil {
			reqMap = make(map[string]interface{})
		}
		reqMap["reqTime"] = reqTime
		requestData = reqMap
	}

	var responseData interface{}
	if body != nil {
		b, _ := json.Marshal(body)
		masked := maskSensitive(b)

		var respMap map[string]interface{}
		json.Unmarshal(masked, &respMap)
		if respMap != nil {
			if data, ok := respMap["data"]; ok && data != nil {
				if dm, ok := data.(map[string]interface{}); ok {
					responseData = dm
				} else if dmArr, ok := data.([]interface{}); ok {
					responseData = dmArr
				}
			}
			if meta, ok := respMap["meta"]; ok && meta != nil {
				if respMapData, ok := responseData.(map[string]interface{}); ok {
					if mm, ok := meta.(map[string]interface{}); ok {
						if tc, ok := mm["total_data"]; ok {
							respMapData["totalData"] = tc
						}
						if ps, ok := mm["page_size"]; ok {
							respMapData["itemCount"] = ps
						}
					}
				}
			}
		}
	}
	if responseData == nil {
		responseData = make(map[string]interface{})
	}
	respMap, _ := responseData.(map[string]interface{})
	if respMap != nil {
		respMap["resTime"] = resTime
	}

	pathFull := c.Request.URL.Path
	if c.Request.URL.RawQuery != "" {
		pathFull = pathFull + "?" + c.Request.URL.RawQuery
	}

	logData := map[string]interface{}{
		"time":       timeISO,
		"level":      "INFO",
		"msg":        "HTTP Transaction completed",
		"reqId":      reqID,
		"userId":     userID,
		"method":     c.Request.Method,
		"path":       pathFull,
		"status":     statusCode,
		"durationMs": durationMs,
		"userAgent":  userAgent,
		"request":    requestData,
		"response":   responseData,
	}

	rawData, _ := json.Marshal(logData)
	host := fmt.Sprintf("%s://%s%s", guessScheme(c), c.Request.Host, c.Request.URL.Path)

	entry := AuditEntry{
		Date:      now.Format("2006-01-02 15:04:05"),
		UserID:    userID,
		Name:      name,
		Role:      role,
		Host:      host,
		IP:        ip,
		Method:    c.Request.Method,
		Status:    fmt.Sprintf("%d", statusCode),
		Data:      rawData,
		CreatedAt: now,
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
	var appError *AppError
	if errors.As(err, &appError) {
		msg := appError.Error()

		hint := appError.Hint
		if appError.StatusCode == http.StatusInternalServerError && os.Getenv("APP_ENV") == "production" {
			hint = ""
		}

		if appError.Err != nil {
			slog.Error("error cause", "cause", appError.Err.Error(), "code", appError.Code)
		}

		body := APIResponse{Success: false, Message: msg, Hint: hint}
		c.JSON(appError.StatusCode, body)
		saveLog(c, appError.StatusCode, body)
		return
	}

	msg := "internal server error"
	if os.Getenv("APP_ENV") != "production" {
		msg = err.Error()
	}
	slog.Error("unexpected error type", "error", err.Error())
	body := APIResponse{Success: false, Message: msg}
	c.JSON(http.StatusInternalServerError, body)
	saveLog(c, http.StatusInternalServerError, body)
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
