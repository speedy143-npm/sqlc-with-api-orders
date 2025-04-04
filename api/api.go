package api

import (
	"net/http"

	"github.com/Iknite-Space/sqlc-example-api/db/repo"
	"github.com/gin-gonic/gin"
)

type MessageHandler struct {
	querier repo.Querier
}

func NewMessageHandler(querier repo.Querier) *MessageHandler {
	return &MessageHandler{
		querier: querier,
	}
}

func (h *MessageHandler) WireHttpHandler() http.Handler {

	r := gin.Default()
	r.Use(gin.CustomRecovery(func(c *gin.Context, _ any) {
		c.String(http.StatusInternalServerError, "Internal Server Error: panic")
		c.AbortWithStatus(http.StatusInternalServerError)
	}))

	r.POST("/message", h.handleCreateMessage)
	r.GET("/message/:id", h.handleGetMessage)
	r.GET("/thread/:id/messages", h.handleGetThreadMessages)
	r.DELETE("/message/:id", h.handleDeleteMessage)
	r.PATCH("/message/:id", h.handleUpdateMessage)
	//creating a thread
	r.POST("/thread", h.handleCreateThread)

	return r
}

func (h *MessageHandler) handleCreateMessage(c *gin.Context) {
	var req repo.CreateMessageParams
	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Check if the thread exists
	_, err = h.querier.GetThreadByID(c, req.Thread)
	if err != nil {
		// Handle the case where the thread does not exist
		c.JSON(http.StatusBadRequest, gin.H{"error": "Thread does not exist"})
		return
	}

	message, err := h.querier.CreateMessage(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, message)
}

func (h *MessageHandler) handleGetMessage(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	message, err := h.querier.GetMessageByID(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, message)
}

func (h *MessageHandler) handleGetThreadMessages(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id is required"})
		return
	}

	messages, err := h.querier.GetMessagesByThread(c, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"thread":   id,
		"topic":    "example",
		"messages": messages,
	})
}

// My code to handle deletion of messages
// using message id

func (h *MessageHandler) handleDeleteMessage(c *gin.Context) {
	// Extract the message ID from the request URI or query parameters
	messageID := c.Param("id") // Assuming the ID is part of the URL, like /messages/:id

	if messageID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "message ID is required"})
		return
	}

	// Perform the deletion
	err := h.querier.DeleteMessage(c, messageID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with a success message
	c.JSON(http.StatusOK, gin.H{"message": "message deleted successfully"})
}

// My code to handle update of messages
// using message id

func (h *MessageHandler) handleUpdateMessage(c *gin.Context) {
	// Get the message ID from the URL parameters
	messageID := c.Param("id")

	// Parse the update details from the request body
	var req repo.PatchMessageParams
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Attach the message ID to the update parameters
	//req.ID = messageID

	// Update the message in the database
	message, err := h.querier.PatchMessage(c, messageID, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Respond with the updated message
	c.JSON(http.StatusOK, message)
}

// creating a create thread handler
func (h *MessageHandler) handleCreateThread(c *gin.Context) {
	var req repo.CreateThreadParams
	err := c.ShouldBindBodyWithJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	thread, err := h.querier.CreateThread(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, thread)
}

//creating a get thread by id handler
