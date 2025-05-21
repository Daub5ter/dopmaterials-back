package consts

const (
	FileServerConfigPath = "./configs/file-server-config.yaml"

	IdParam         = "id"
	CategoryIdParam = "category_id"

	FileNameParam              = "filename"
	FileNameWithExtensionParam = "filename_with_extension"
	FileNamePartParam          = "filename_part"

	HttpRequestBodyTooLarge = "http: request body too large"

	StoragePhotosPath = "storage/photos/"
	MultiFormPhotoKey = "photo"
	MaxPhotosSizeMB   = 50

	StorageVideosPath = "storage/videos/"
	MultiFormVideoKey = "video"
	MaxVideoSizeMB    = 500

	OffsetParam   = "offset"
	FindPartParam = "find_part"
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
