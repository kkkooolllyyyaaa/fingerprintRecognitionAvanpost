package route

import (
	"fingerprintRecognitionAvanpost/internal/algorithm"
	"fingerprintRecognitionAvanpost/internal/handler"
	"github.com/gin-gonic/gin"
)

func InitPredictRoutes(algorithm algorithm.Algorithm, route *gin.Engine) {

	currencyRateHandler := handler.NewPredictHandler(algorithm)

	groupRoute := route.Group("api/")
	groupRoute.POST("/predict", currencyRateHandler.HandlePredictPost)
}
