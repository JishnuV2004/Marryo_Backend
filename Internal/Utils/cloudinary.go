package utils

import (
	"context"
	"os"

	"github.com/cloudinary/cloudinary-go/v2"
)

var CloudCtx = context.Background()

func Cloudinary() (*cloudinary.Cloudinary, error) {
	return cloudinary.NewFromParams(
		os.Getenv("CLOUDINARY_CLOUD_NAME"),
		os.Getenv("CLOUDINARY_API_KEY"),
		os.Getenv("CLOUDINARY_API_SECRET"),
	)
}