package handlers

import (
	"bytes"
	"encoding/base64"
	"image"
	"image/jpeg"
	"image/png"
	"img-compress-demo/backend/compression"
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	timeFormat        = "2006-01-02 15:04:05"
	maxFormSize       = 32 << 20
	defaultComponents = 64
)

type CompressRequest struct {
	Image         string `json:"image"`
	Method        string `json:"method"`
	NumComponents int    `json:"numComponents"`
}

type ComponentLevelResponse struct {
	NumComponents int    `json:"numComponents"`
	DataSize      int    `json:"dataSize"`
	ImageData     string `json:"imageData"`
}

type CompressResponse struct {
	Method          string                   `json:"method"`
	OriginalSize    int                      `json:"originalSize"`
	ComponentLevels []ComponentLevelResponse `json:"componentLevels"`
}

func logf(requestID, format string, args ...interface{}) {
	log.Printf("["+timeFormat+"] [%s] "+format, append([]interface{}{time.Now(), requestID}, args...)...)
}

func errorResponse(c *gin.Context, requestID string, status int, message string, err error) {
	logf(requestID, "ERROR: %s: %v", message, err)
	c.JSON(status, gin.H{"error": message})
}

func CompressImage(c *gin.Context) {
	requestID := time.Now().Format("20060102-150405")
	logf(requestID, "Starting compression request")

	if err := c.Request.ParseMultipartForm(maxFormSize); err != nil {
		errorResponse(c, requestID, http.StatusBadRequest, "Failed to parse form", err)
		return
	}

	file, header, err := c.Request.FormFile("image")
	if err != nil {
		errorResponse(c, requestID, http.StatusBadRequest, "No image file provided", err)
		return
	}
	defer file.Close()

	logf(requestID, "Received file: %s (size: %d bytes)", header.Filename, header.Size)

	method := c.Request.FormValue("method")
	numComponents := defaultComponents
	if numComponentsStr := c.Request.FormValue("numComponents"); numComponentsStr != "" {
		if parsed, err := strconv.Atoi(numComponentsStr); err == nil {
			numComponents = parsed
		}
	}

	logf(requestID, "Compression parameters: method=%s, components=%d", method, numComponents)

	var img image.Image
	switch filepath.Ext(header.Filename) {
	case ".jpg", ".jpeg":
		img, err = jpeg.Decode(file)
	case ".png":
		img, err = png.Decode(file)
	default:
		errorResponse(c, requestID, http.StatusBadRequest, "Unsupported image format. Use JPG or PNG", nil)
		return
	}

	if err != nil {
		errorResponse(c, requestID, http.StatusBadRequest, "Failed to decode image: "+err.Error(), err)
		return
	}

	logf(requestID, "Image decoded successfully: %dx%d pixels", img.Bounds().Dx(), img.Bounds().Dy())

	var compressor compression.Compressor
	switch method {
	case "DCT":
		compressor = compression.NewDCTCompressor()
	case "SVD":
		compressor = compression.NewSVDCompressor()
	default:
		errorResponse(c, requestID, http.StatusBadRequest, "Invalid method. Use DCT or SVD", nil)
		return
	}

	logf(requestID, "Starting %s compression...", method)
	startTime := time.Now()
	result, err := compressor.Compress(img, numComponents)
	if err != nil {
		errorResponse(c, requestID, http.StatusInternalServerError, "Compression failed: "+err.Error(), err)
		return
	}
	logf(requestID, "Compression completed in %v (%d component levels)", time.Since(startTime), len(result.ComponentLevels))

	logf(requestID, "Encoding %d component images to base64...", len(result.ComponentLevels))
	encodeStart := time.Now()
	componentResponses := make([]ComponentLevelResponse, len(result.ComponentLevels))
	for i, level := range result.ComponentLevels {
		imageData, err := imageToBase64(level.Image)
		if err != nil {
			errorResponse(c, requestID, http.StatusInternalServerError, "Failed to encode image", err)
			return
		}

		componentResponses[i] = ComponentLevelResponse{
			NumComponents: level.NumComponents,
			DataSize:      level.DataSize,
			ImageData:     imageData,
		}
	}
	logf(requestID, "Encoding completed in %v", time.Since(encodeStart))

	logf(requestID, "Request completed successfully (total: %v)", time.Since(startTime))

	c.JSON(http.StatusOK, CompressResponse{
		Method:          result.Method,
		OriginalSize:    result.OriginalSize,
		ComponentLevels: componentResponses,
	})
}

func imageToBase64(img image.Image) (string, error) {
	var buf bytes.Buffer
	if err := png.Encode(&buf, compression.ImageToRGBA(img)); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(buf.Bytes()), nil
}
