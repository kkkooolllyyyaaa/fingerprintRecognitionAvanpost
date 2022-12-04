package handler

import (
	"context"
	"fingerprintRecognitionAvanpost/internal/algorithm"
	"fingerprintRecognitionAvanpost/internal/file"
	"fingerprintRecognitionAvanpost/internal/myimage"
	"fingerprintRecognitionAvanpost/internal/preprocess"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/image/bmp"
	"image"
	"mime/multipart"
	"net/http"
)

type PredictHandler struct {
	algorithm algorithm.Algorithm
}

func NewPredictHandler(algorithm algorithm.Algorithm) *PredictHandler {
	return &PredictHandler{algorithm: algorithm}
}

func (h *PredictHandler) HandlePredictPost(ctx *gin.Context) {
	fl, err := extractFile(ctx)
	if err != nil {
		APIResponse(ctx, "Can't extract image", http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	img, err := bmp.Decode(fl)
	if err != nil {
		APIResponse(ctx, "Can't decode bmp image", http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	preprocessedData, err := preprocessToData(ctx, img)
	if err != nil {
		APIResponse(ctx, "Can't preprocess data", http.StatusInternalServerError, http.MethodPost, nil)
		return
	}

	foundFilename, err := h.algorithm.Predict(preprocessedData)
	if err != nil {
		APIResponse(ctx, "Can't predict for image", http.StatusInternalServerError, http.MethodPost, nil)
		return
	}

	APIResponse(
		ctx,
		"Success",
		http.StatusOK,
		http.MethodPost,
		PredictPostResponse{
			Person: foundFilename,
		},
	)
}

func extractFile(ctx *gin.Context) (multipart.File, error) {
	formImage, err := ctx.FormFile("image")
	if err != nil {
		return nil, errors.Wrap(err, "Extract form file")
	}

	file, err := formImage.Open()
	if err != nil {
		return nil, errors.Wrap(err, "Open image from form file")
	}

	return file, nil
}

func preprocessToData(ctx context.Context, img image.Image) (*preprocess.Data, error) {
	gray := file.ToGray(img)
	myImg := myimage.NewMyImage(gray, "filename")
	preprocessedData, err := preprocess.PreprocessOne(ctx, myImg, false)
	if err != nil {
		return nil, err
	}
	return preprocessedData, nil
}
