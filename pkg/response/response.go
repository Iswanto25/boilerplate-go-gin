package response

import (
	"errors"
	"net/http"

	appErr "github.com/edustack/go-boilerplate/pkg/errors"
	"github.com/gin-gonic/gin"
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

func Success(c *gin.Context, statusCode int, message string, data interface{}) {
	c.JSON(statusCode, APIResponse{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(c *gin.Context, statusCode int, message string) {
	c.JSON(statusCode, APIResponse{
		Success: false,
		Message: message,
	})
}

func WriteError(c *gin.Context, err error) {
	var appErr *appErr.AppError
	if errors.As(err, &appErr) {
		Error(c, appErr.StatusCode, appErr.Message)
		return
	}
	Error(c, http.StatusInternalServerError, "internal server error")
}

func Paginated(c *gin.Context, statusCode int, message string, data interface{}, page, pageSize int, totalData int64) {
	totalPage := int(totalData) / pageSize
	if int(totalData)%pageSize > 0 {
		totalPage++
	}

	c.JSON(statusCode, PaginatedResponse{
		Success: true,
		Message: message,
		Data:    data,
		Meta: Pagination{
			Page:      page,
			PageSize:  pageSize,
			TotalData: totalData,
			TotalPage: totalPage,
		},
	})
}
