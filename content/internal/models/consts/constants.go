package consts

const (
	ContentConfigPath = "./configs/content-config.yaml"

	IdParam         = "id"
	CategoryIdParam = "category_id"

	OffsetParam   = "offset"
	FindPartParam = "find_part"

	LimitMaterials = 3
)

var (
	AllowedPhotoContentType = map[string]string{
		"image/png":  ".png",
		"image/jpeg": ".jpg",
	}

	AllowedVideoContentType = map[string]string{
		"video/mp4": ".mp4",
	}
)
