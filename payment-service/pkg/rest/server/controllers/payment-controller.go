package controllers

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/mahendraintelops/intelops-procurement-solution-v10/payment-service/pkg/rest/server/daos/clients/sqls"
	"github.com/mahendraintelops/intelops-procurement-solution-v10/payment-service/pkg/rest/server/models"
	"github.com/mahendraintelops/intelops-procurement-solution-v10/payment-service/pkg/rest/server/services"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"net/http"
	"os"
	"strconv"
)

type PaymentController struct {
	paymentService *services.PaymentService
}

func NewPaymentController() (*PaymentController, error) {
	paymentService, err := services.NewPaymentService()
	if err != nil {
		return nil, err
	}
	return &PaymentController{
		paymentService: paymentService,
	}, nil
}

func (paymentController *PaymentController) CreatePayment(context *gin.Context) {
	// validate input
	var input models.Payment
	if err := context.ShouldBindJSON(&input); err != nil {
		log.Error(err)
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// trigger payment creation
	if _, err := paymentController.paymentService.CreatePayment(&input); err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Payment created successfully"})
}

func (paymentController *PaymentController) UpdatePayment(context *gin.Context) {
	// validate input
	var input models.Payment
	if err := context.ShouldBindJSON(&input); err != nil {
		log.Error(err)
		context.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// trigger payment update
	if _, err := paymentController.paymentService.UpdatePayment(id, &input); err != nil {
		log.Error(err)
		if errors.Is(err, sqls.ErrNotExists) {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{"message": "Payment updated successfully"})
}

func (paymentController *PaymentController) FetchPayment(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// trigger payment fetching
	payment, err := paymentController.paymentService.GetPayment(id)
	if err != nil {
		log.Error(err)
		if errors.Is(err, sqls.ErrNotExists) {
			context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	serviceName := os.Getenv("SERVICE_NAME")
	collectorURL := os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	if len(serviceName) > 0 && len(collectorURL) > 0 {
		// get the current span by the request context
		currentSpan := trace.SpanFromContext(context.Request.Context())
		currentSpan.SetAttributes(attribute.String("payment.id", strconv.FormatInt(payment.Id, 10)))
	}

	context.JSON(http.StatusOK, payment)
}

func (paymentController *PaymentController) DeletePayment(context *gin.Context) {
	id, err := strconv.ParseInt(context.Param("id"), 10, 64)
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// trigger payment deletion
	if err := paymentController.paymentService.DeletePayment(id); err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, gin.H{
		"message": "Payment deleted successfully",
	})
}

func (paymentController *PaymentController) ListPayments(context *gin.Context) {
	// trigger all payments fetching
	payments, err := paymentController.paymentService.ListPayments()
	if err != nil {
		log.Error(err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	context.JSON(http.StatusOK, payments)
}

func (*PaymentController) PatchPayment(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "PATCH",
	})
}

func (*PaymentController) OptionsPayment(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "OPTIONS",
	})
}

func (*PaymentController) HeadPayment(context *gin.Context) {
	context.JSON(http.StatusOK, gin.H{
		"message": "HEAD",
	})
}
