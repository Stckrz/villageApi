package services

import (

	"github.com/Stckrz/villageApi/internal/db/models"
	"gorm.io/gorm"
)

type BuildingService interface {
	GetBuildingByID(id uint) (models.Building, error)
	ListBuildings() ([]models.Building, error)
	CreateBuilding(building models.Building) (models.Building, error)
	DeleteBuilding(id uint) (error)
	UpdateBuilding(building models.Building, id uint) (error)
}

type buildingService struct {
	db *gorm.DB
}

func NewBuildingService(db *gorm.DB) BuildingService {
	return &buildingService{db: db}
} 

func (s *buildingService) ListBuildings() ([]models.Building, error) {
	var buildings []models.Building
	if err := s.db.Preload("Categories").Preload("Tasks").Preload("Tasks.Building").Find(&buildings).Error; err != nil {
		return nil, err
	}

	return buildings, nil
}

func (s *buildingService) CreateBuilding(building models.Building) (models.Building, error){
	if err := s.db.Create(&building).Error; err != nil{
		return models.Building{}, err
	}
	if err := s.db.Preload("Categories").First(&building, building.ID).Error; err != nil {
		return models.Building{}, err
	}

	return building, nil
}

func (s *buildingService) DeleteBuilding(id uint) error{
	result := s.db.Delete(&models.Building{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s *buildingService) UpdateBuilding(building models.Building, id uint) error{
	return s.db.Transaction(func(transaction *gorm.DB) error {
		result := transaction.
			Model(&models.Building{}).
			Where("id = ?", id).
			Updates(map[string]any{
				"name": building.Name,
				"description": building.Description,
				"thumbnail_path": building.ThumbnailPath,
				"image_path": building.ImagePath,
			})
			if result.Error != nil {
				return result.Error
			}

			if result.RowsAffected == 0 {
				return gorm.ErrRecordNotFound
			}
		if err := transaction.
			Where("building_id = ?", id).
			Delete(&models.BuildingCategory{}).Error; err != nil {
				return err
			}
		if len(building.Categories) > 0 {
			for categoryIndex := range building.Categories {
				building.Categories[categoryIndex].ID = 0
				building.Categories[categoryIndex].BuildingID = id
			}
			if err := transaction.Create(&building.Categories).Error; err != nil {
				return err
			}
		}
	return nil

	})
}

func (s *buildingService) GetBuildingByID(id uint) (models.Building, error) {
	var building models.Building
	err := s.db.
		Preload("Categories").
		Preload("Tasks").
		First(&building, id).Error
	if err != nil {
		return models.Building{}, err
	}
	return building, nil
}
