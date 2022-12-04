package app

import (
	"context"
	"fingerprintRecognitionAvanpost/internal/algorithm"
	"fingerprintRecognitionAvanpost/internal/preprocess"
	"fingerprintRecognitionAvanpost/internal/route"
	"fingerprintRecognitionAvanpost/pkg/logger"
	helmet "github.com/danielkov/gin-helmet"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	"golang.org/x/sync/errgroup"
	"net/http"
)

func RunTest(ctx context.Context, erg *errgroup.Group, server *http.Server, data []*preprocess.Data) error {
	alg := algorithm.NewKeyPointsAlgorithm(data)

	logger.Info(ctx).Msg("Setuping routers...")
	router := SetupRouter(alg)
	server.Handler = router.Handler()

	logger.Info(ctx).Msg("Listening http connects...")
	erg.Go(func() error {
		return server.ListenAndServe()
	})

	return nil
}

func SetupRouter(algorithm algorithm.Algorithm) *gin.Engine {
	/**
	@description Init Router
	*/
	router := gin.Default()
	/**
	@description Setup Mode Application
	*/
	gin.SetMode(gin.ReleaseMode)

	/**
	@description Setup Middleware
	*/
	const RequestBodyMaxSize = 1024 * 1024
	router.Use(
		cors.New(cors.Config{
			AllowOrigins:  []string{"*"},
			AllowMethods:  []string{"*"},
			AllowHeaders:  []string{"*"},
			AllowWildcard: true,
		}),
		limits.RequestSizeLimiter(RequestBodyMaxSize),
	)
	router.Use(helmet.Default())
	router.Use(gzip.Gzip(gzip.BestCompression))

	/**
	@description Init All Route
	*/
	route.InitPredictRoutes(algorithm, router)

	return router
}
