package main

import (
	"FxService/controllers"
	"FxService/model/dto/response"
	"FxService/service/common"
	"FxService/service/dal"
	"context"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/gin-gonic/gin"

	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	sdkresource "go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

var (
	resource          *sdkresource.Resource
	initResourcesOnce sync.Once
)

func main() {
	_ = initTracerProvider()
	dal.Init(os.Getenv("MONGODB_URI"))
	common.InitLog()

	router := gin.Default()
	router.Use(otelgin.Middleware(os.Getenv("OTEL_SERVICE_NAME")))
	router.Use(GlobalErrorHandler)
	controllers.AddRoutes(router)
	log.Fatal(router.Run("0.0.0.0:8080"))

}

func GlobalErrorHandler(c *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			// Handle the error here, you can log it or send an error response
			c.IndentedJSON(http.StatusOK, common.GetSimpleResponse[response.ForexDataResponse](nil, response.InternalError, &[]response.Error{
				{Code: "INTERNAL_ERROR", Message: "Please try again later", Details: "Something went wrong. Please try again later"},
			}))
			common.Logger.Errorf("Unhandled Error in %s %s. Exception:%v", c.Request.Method, c.Request.URL, err)
		}
	}()
	c.Next()
}

func initResource() *sdkresource.Resource {
	initResourcesOnce.Do(func() {
		extraResources, _ := sdkresource.New(
			context.Background(),
			sdkresource.WithOS(),
			sdkresource.WithProcess(),
			sdkresource.WithContainer(),
			sdkresource.WithHost(),
		)
		resource, _ = sdkresource.Merge(
			sdkresource.Default(),
			extraResources,
		)
	})
	return resource
}

func initTracerProvider() *sdktrace.TracerProvider {
	ctx := context.Background()

	exporter, err := otlptracegrpc.New(ctx)
	if err != nil {
		log.Fatalf("OTLP Trace gRPC Creation: %v", err)
	}
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(initResource()),
	)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))
	return tp
}
