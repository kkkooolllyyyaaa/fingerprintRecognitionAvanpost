package app

import (
	"context"
	"fingerprintRecognitionAvanpost/internal/algorithm"
	"fingerprintRecognitionAvanpost/internal/file"
	"fingerprintRecognitionAvanpost/internal/route"
	helmet "github.com/danielkov/gin-helmet"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/gzip"
	limits "github.com/gin-contrib/size"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"net/http"
)

func RunTest(ctx context.Context, erg *errgroup.Group, server *http.Server) error {
	fileRoot := "files/preprocessed/SOCOFing/Real/"

	dataWorker := file.NewDataWorker(fileRoot)
	data, err := dataWorker.InitByReadingAll()
	if err != nil {
		return errors.Wrap(err, "Error while init by reading all")
	}

	alg := algorithm.NewKeyPointsAlgorithm(data)

	router := SetupRouter(alg)
	server.Handler = router.Handler()

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
