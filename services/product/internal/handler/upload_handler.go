package handler

import (
	"context"
	"path/filepath"
	"strings"
	"time"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/product/internal/storage"
	"github.com/gofiber/fiber/v2"
)

type UploadHandler struct {
	storage *storage.MinioStorage
}

func NewUploadHandler(storage *storage.MinioStorage) *UploadHandler {
	return &UploadHandler{storage: storage}
}

func (h *UploadHandler) UploadImage(c *fiber.Ctx) error {
	// Try fetching from "image" form key
	fileHeader, err := c.FormFile("image")
	if err != nil {
		// Fallback to "file" form key
		fileHeader, err = c.FormFile("file")
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"success": false,
				"message": "Image file is required (form key 'image' or 'file')",
			})
		}
	}

	// Validate extension
	ext := strings.ToLower(filepath.Ext(fileHeader.Filename))
	allowed := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".webp": true,
		".gif":  true,
	}
	if !allowed[ext] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "Only image files are allowed (.jpg, .jpeg, .png, .webp, .gif)",
		})
	}

	// Open the file
	file, err := fileHeader.Open()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to open file",
		})
	}
	defer file.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Upload to Minio
	url, err := h.storage.UploadFile(ctx, fileHeader.Filename, file, fileHeader.Size, fileHeader.Header.Get("Content-Type"))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"success": false,
			"message": "Failed to upload to storage: " + err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"success":   true,
		"image_url": url,
	})
}
