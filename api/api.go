package api

import (
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	httpRequests "github.com/Iknite-Space/sqlc-example-api/campay_api/Payment"
	"github.com/Iknite-Space/sqlc-example-api/db/repo"
	"github.com/gin-gonic/gin"
	//"google.golang.org/protobuf/internal/errors"
	//"google.golang.org/protobuf/internal/errors"
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

	r.POST("/customer", h.handleCreateCustomer)
	// r.GET("/customer/:phone", h.handleGetCustomerByPhoneNo)

	// r.POST("/product", h.handleCreateProduct)
	//r.GET("/product/:id", h.handleGetProductById)
	r.POST("/order/:customer_id/placeorder", h.handleCreateOrder)

	return r
}

func (h *MessageHandler) handleCreateCustomer(c *gin.Context) {
	var req repo.CreateCustomerParams
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	customer, err := h.querier.CreateCustomer(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *MessageHandler) handleGetCustomerByPhoneNo(c *gin.Context) {
	var req = c.Param("phone")
	if req == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Phone number is required"})
		return
	}
	customer, err := h.querier.GetCustomerByPhoneNo(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, customer)
}

func (h *MessageHandler) handleCreateProduct(c *gin.Context) {
	var req repo.CreateProductParams
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	product, err := h.querier.CreateProduct(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *MessageHandler) handleGetProductById(c *gin.Context) {
	var req = c.Param("id")
	if req == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id number is required"})
		return
	}
	product, err := h.querier.GetProductById(c, req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, product)
}

func (h *MessageHandler) handleCreateOrder(c *gin.Context) {
	id := c.Param("customer_id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id not available"})
		return
	}
	var req repo.CreateOrderParams
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var reg repo.CreateCustomerParams
	errw := c.ShouldBindJSON(&reg)
	if errw != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": errw.Error()})
		return
	}

	// Check if the customer exists
	_, errs := h.querier.GetCustomerById(c, id)
	if errs != nil {
		if errors.Is(errs, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer ID Not found"})
		}
		// else {
		// 	c.JSON(http.StatusInternalServerError, gin.H{"error": errs.Error()})
		// }
		// return

	}

	order, err := h.querier.CreateOrder(c, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	apikey := os.Getenv("API_KEY")

	trans := httpRequests.RequestPayment(apikey, reg.Phoneno, fmt.Sprintf("%v", req.TotalPrice), "description", "ex_ref")

	time.Sleep(10 * time.Second)

	state := httpRequests.CheckPaymentStatus(apikey, trans.Reference)

	var st = state.Status
	var od = order.ID

	orders, errq := h.querier.UpdateOrderById(c, od, st)
	if errq != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errq.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Success": orders, "Payment_details": trans, "status": state})

}

/*
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

*/
