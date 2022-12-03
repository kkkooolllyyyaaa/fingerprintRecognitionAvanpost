package handler

import (
	"fingerprintRecognitionAvanpost/internal/algorithm"
	"fingerprintRecognitionAvanpost/internal/preprocess"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"golang.org/x/image/bmp"
	"mime/multipart"
	"net/http"
	"strconv"
)

type PredictHandler struct {
	algorithm algorithm.Algorithm
}

func NewPredictHandler(algorithm algorithm.Algorithm) *PredictHandler {
	return &PredictHandler{algorithm: algorithm}
}

func (h *PredictHandler) HandlePredictPost(ctx *gin.Context) {
	file, err := extractFile(ctx)
	if err != nil {
		APIResponse(ctx, "Can't extract image", http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	_, err = bmp.Decode(file)

	if err != nil {
		APIResponse(ctx, "Can't decode bmp image", http.StatusBadRequest, http.MethodPost, nil)
		return
	}

	personIndex, err := h.algorithm.Predict(preprocess.Data{})
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
			Person: strconv.Itoa(personIndex),
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
