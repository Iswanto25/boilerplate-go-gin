package audit

import (
	"net/http"
	"strconv"

	"github.com/edustack/go-boilerplate/pkg"
	"github.com/gin-gonic/gin"
)

// Handler menangani HTTP request untuk audit logs.
type Handler struct {
	service Service
}

// NewHandler membuat instance audit handler.
func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

// GetLogs godoc
// @Summary  Get list of audit logs
// @Tags     Audit
// @Produce  json
// @Param    page      query  int     false  "Page number (default: 1)"
// @Param    page_size query  int     false  "Items per page (default: 20, max: 100)"
// @Param    date      query  string  false  "Filter by date (format: 2006-01-02)"
// @Param    method    query  string  false  "Filter by HTTP method (GET, POST, etc.)"
// @Param    status    query  int     false  "Filter by HTTP status code"
// @Param    search    query  string  false  "Search by path or host"
// @Success  200  {object}  pkg.PaginatedResponse
// @Router   /api/v1/audit/logs [get]
func (h *Handler) GetLogs(c *gin.Context) {
	params := QueryParams{
		Page:     parseIntQuery(c, "page", 1),
		PageSize: parseIntQuery(c, "page_size", 20),
		Date:     c.Query("date"),
		Method:   c.Query("method"),
		Status:   c.Query("status"),
		Search:   c.Query("search"),
	}

	logs, total, err := h.service.GetLogs(c.Request.Context(), params)
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, "failed to retrieve audit logs")
		return
	}

	pkg.Paginated(c, http.StatusOK, "Audit logs retrieved successfully",
		logs, params.Page, params.PageSize, total)
}

// GetLogDetail godoc
// @Summary  Get audit log detail by ID
// @Tags     Audit
// @Produce  json
// @Param    id  path  int  true  "Log ID"
// @Success  200  {object}  pkg.APIResponse
// @Router   /api/v1/audit/logs/{id} [get]
func (h *Handler) GetLogDetail(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil || id <= 0 {
		pkg.Error(c, http.StatusBadRequest, "Invalid log ID, must be a positive integer")
		return
	}

	log, err := h.service.GetLogByID(c.Request.Context(), id)
	if err != nil {
		pkg.WriteError(c, err)
		return
	}

	pkg.Success(c, http.StatusOK, "Log retrieved successfully", log)
}

// parseIntQuery mengambil query param integer dengan nilai default.
func parseIntQuery(c *gin.Context, key string, defaultVal int) int {
	val := c.Query(key)
	if val == "" {
		return defaultVal
	}
	n, err := strconv.Atoi(val)
	if err != nil || n <= 0 {
		return defaultVal
	}
	return n
}
