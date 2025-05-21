package consts

const (
	ApiGatewayConfigPath = "./configs/api-gateway-config.yaml"
	ApiGatewayCertPem    = "./build/api-gateway-cert.pem"
	ApiGatewayKeyPem     = "./build/api-gateway-key.pem"

	FileNameParam              = "filename"
	FileNameWithExtensionParam = "filename_with_extension"
	FileNamePartParam          = "file_name_part"

	ContentUrl    = "http://content:44301/"
	FileServerUrl = "http://file-server:44302/"

	IdParam         = "id"
	CategoryIdParam = "category_id"

	OffsetParam   = "offset"
	FindPartParam = "find_part"

	// TODO DELETE
	StorageVideosPath = "storage/videos/"
)
