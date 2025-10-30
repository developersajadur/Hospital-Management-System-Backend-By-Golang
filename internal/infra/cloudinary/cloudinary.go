package cloudinary

import (
	"hospital_management_system/config"
	"log"

	"github.com/cloudinary/cloudinary-go/v2"
)

func NewCloudinary() (*cloudinary.Cloudinary, error) {
	cld, err := cloudinary.NewFromParams(
		config.ENV.CloudinaryCloudName,
		config.ENV.CloudinaryApiKey,
		config.ENV.CloudinaryApiSecret,
	)
	if err != nil {
		log.Println("Failed to initialize Cloudinary:", err)
		return nil, err
	}
	return cld, nil
}
