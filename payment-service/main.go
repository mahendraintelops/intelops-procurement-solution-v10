package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/mahendraintelops/intelops-procurement-solution-v10/payment-service/config"
	restcontrollers "github.com/mahendraintelops/intelops-procurement-solution-v10/payment-service/pkg/rest/server/controllers"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/sinhashubham95/go-actuator"
	log "github.com/sirupsen/logrus"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"os"

	"github.com/mahendraintelops/intelops-procurement-solution-v10/payment-service/pkg/rest/client"
)

var (
	serviceName  = os.Getenv("SERVICE_NAME")
	collectorURL = os.Getenv("OTEL_EXPORTER_OTLP_ENDPOINT")
	insecure     = os.Getenv("INSECURE_MODE")
)

func main() {

	// rest server configuration
	router := gin.Default()
	var restTraceProvider *sdktrace.TracerProvider
	if len(serviceName) > 0 && len(collectorURL) > 0 {
		// add opentel
		restTraceProvider = config.InitRestTracer(serviceName, collectorURL, insecure)
		router.Use(otelgin.Middleware(serviceName))
	}
	defer func() {
		if restTraceProvider != nil {
			if err := restTraceProvider.Shutdown(context.Background()); err != nil {
				log.Printf("Error shutting down tracer provider: %v", err)
			}
		}
	}()
	// add actuator
	addActuator(router)
	// add prometheus
	addPrometheus(router)

	paymentController, err := restcontrollers.NewPaymentController()
	if err != nil {
		log.Errorf("error occurred: %v", err)
		os.Exit(1)
	}

	v1 := router.Group("/v1")
	{

		v1.GET("/payments/:id", paymentController.FetchPayment)
		v1.POST("/payments", paymentController.CreatePayment)
		v1.PUT("/payments/:id", paymentController.UpdatePayment)
		v1.DELETE("/payments/:id", paymentController.DeletePayment)
		v1.GET("/payments", paymentController.ListPayments)
		v1.PATCH("/payments/:id", paymentController.PatchPayment)
		v1.HEAD("/payments", paymentController.HeadPayment)
		v1.OPTIONS("/payments", paymentController.OptionsPayment)

	}

	Port := ":4565"
	log.Println("Server started")
	if err = router.Run(Port); err != nil {
		log.Errorf("error occurred: %v", err)
		os.Exit(1)
	}

	// this will not be called as the control won't reach here.
	// call external services here if the HasRestClients value is true
	// (that means this repo is a client to external service(s)
	var err0 error

	bNodeC2, err0 := client.ExecuteNodeC2()
	if err0 != nil {
		log.Printf("error occurred: %v", err0)
		return
	}
	log.Printf("response received: %s", string(bNodeC2))

}

func prometheusHandler() gin.HandlerFunc {
	h := promhttp.Handler()

	return func(c *gin.Context) {
		h.ServeHTTP(c.Writer, c.Request)
	}
}

func addPrometheus(router *gin.Engine) {
	router.GET("/metrics", prometheusHandler())
}

func addActuator(router *gin.Engine) {
	actuatorHandler := actuator.GetActuatorHandler(&actuator.Config{Endpoints: []int{
		actuator.Env,
		actuator.Info,
		actuator.Metrics,
		actuator.Ping,
		// actuator.Shutdown,
		actuator.ThreadDump,
	},
		Env:     "dev",
		Name:    "payment-service",
		Port:    4565,
		Version: "0.0.1",
	})
	ginActuatorHandler := func(ctx *gin.Context) {
		actuatorHandler(ctx.Writer, ctx.Request)
	}
	router.GET("/actuator/*endpoint", ginActuatorHandler)
}

func init() {
	// Log as JSON instead of the default ASCII formatter.
	// log.SetFormatter(&log.JSONFormatter{})
	log.SetFormatter(&log.TextFormatter{
		DisableColors: false,
		FullTimestamp: true,
	})
	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	log.SetOutput(os.Stdout)
	// Only log the warning severity or above.
	log.SetLevel(log.InfoLevel)
}
