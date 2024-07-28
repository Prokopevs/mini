package handler

import (
	"net/http"

	"github.com/Prokopevs/mini/server/internal/model"
	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// @Summary  	 Create message
// @Description  Add message to DB
// @Accept 	 	 json
// @Produce 	 json
// @Param        request  body  model.MessageCreate  true  "body"
// @Success 	 200  {array}   model.Message
// @Failure      400  {object}  errorResponse
// @Failure      500  {object}  errorResponse
// @Router       /api/v1/create [post]“
func (h *HTTP) CreateMessage(c *gin.Context) {
	m := &model.MessageCreate{}
	if err := c.ShouldBindJSON(m); err != nil {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "Bad request", Message: "wrond data provided"})
		return
	}

	if m.Message == "" {
		c.JSON(http.StatusBadRequest, errorResponse{Error: "Bad request", Message: "message cannot be empty"})
		return
	}

	err := h.service.CreateMessage(c.Request.Context(), m)
	if err != nil {
		h.log.Errorw("Create message.", "err", err)
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "Internal server error", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, "Ok")
}

// @Summary  	 Get all messages
// @Description  Get messages
// @Accept 	 	 json
// @Produce 	 json
// @Success 	 200  {array}   model.Message
// @Failure      500  {object}  errorResponse
// @Router       /api/v1/getMessages [get]“
func (h *HTTP) GetMessages(c *gin.Context) {
	res, err := h.service.GetMessages(c.Request.Context())
	if err != nil {
		h.log.Errorw("Get messages.", "err", err)
		c.JSON(http.StatusInternalServerError, errorResponse{Error: "Internal server error", Message: err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
