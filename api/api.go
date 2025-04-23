package api

import (
	"crypto/rand"
	"database/sql"
	"errors"
	"net/http"
	"os"
	"time"

	httpRequests "github.com/Iknite-Space/sqlc-example-api/campay"
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

	r.POST("/customer", h.handleCreateCustomer)
	// r.GET("/customer/:phone", h.handleGetCustomerByPhoneNo)

	r.POST("/product", h.handleCreateProduct)
	r.GET("/product/:id", h.handleGetProductById)
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

/*
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
*/

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

	// Check if the customer exists
	customer, errs := h.querier.GetCustomerById(c, id)
	if errs != nil {
		if errors.Is(errs, sql.ErrNoRows) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Customer ID Not found"})
		}

	}

	order, err := h.querier.CreateOrder(c, id, req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "here"})
		return
	}

	orde, err := h.querier.UpdateOrderTotalPriceById(c, order.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	apikey := os.Getenv("API_KEY")

	description := "e-commerce order payment"
	ref := generateRandomReference(12)

	// making payment request via campay api
	trans := httpRequests.RequestPayment(apikey, customer[0].Phoneno, orde.TotalPrice, description, ref)

	time.Sleep(25 * time.Second)
	// checking the payment status from campay api
	state := httpRequests.CheckPaymentStatus(apikey, trans.Reference)

	// updating the payment status from campay api to the database
	var updateParams = repo.UpdateOrderByIdParams{
		ID:          order.ID,
		OrderStatus: &state.Status,
	}
	orders, errq := h.querier.UpdateOrderById(c, updateParams)
	if errq != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": errq.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Payment_details": trans, "Successfully created order": orders})

	c.JSON(http.StatusOK, orders)

}

// function to generate random string to be used as reference of payments
func generateRandomReference(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		randomByte, _ := rand.Prime(rand.Reader, 8)
		result[i] = charset[randomByte.Int64()%int64(len(charset))]
	}
	return string(result)
}
