package handlers

import "content/internal/models"

type database interface {
	GetMaterials(categoryId, limit *uint32, offset uint32) ([]*models.Material, error)
	GetSearchedMaterials(ids []string) ([]*models.Material, error)
	GetMaterial(id uint64) (*models.Material, error)
	GetMaterialsIdsByCategory(categoryId uint32) ([]uint64, error)
	InsertMaterial(animation *models.Material) (operationTransaction any, id uint64, err error)
	UpdateMaterial(animation *models.Material) (operationTransaction any, err error)
	DeleteMaterial(id uint64) (operationTransaction any, err error)
	RestoreMaterial(id uint64) (operationTransaction any, err error)

	GetNullParentCategories() ([]*models.Category, error)
	GetSubsidiariesCategories(categoryId uint32) ([]*models.Category, error)
	GetCategory(id uint32) (*models.Category, error)
	InsertCategory(category *models.Category) (id uint32, err error)
	UpdateCategory(category *models.Category) error
	DeleteCategory(id uint32) (operationTransaction any, err error)

	ConfirmOperation(operationTransaction any) error
	CancelOperation(operationTransaction any) error
}

type search interface {
	PutMaterial(material models.MaterialSearch) error
	DeleteMaterial(materialId string) error
	SearchMaterials(findPart string, categoryId *uint32, offset uint32) (ids []string, err error)
}
