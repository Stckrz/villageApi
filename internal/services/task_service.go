package services

import (
	"github.com/Stckrz/villageApi/internal/db/models"
	"gorm.io/gorm"
)

type TaskService interface {
	// ListTasks() ([]models.Task, error)
	// ListTasksByBuildingId(buildingId uint) ([]models.Task, error)
	CreateTask(task models.Task) (models.Task, error)
	DeleteTask(id uint) error
	UpdateTask(task models.Task, id uint) (error)
}

type taskService struct {
	db *gorm.DB
}

func NewTaskService(db *gorm.DB) TaskService {
	return &taskService{db: db}
}

// func (s *taskService) ListTasks() ([]models.Task, error) {
// 	var tasks []models.Task
// 	if err := s.db.Preload("Categories").Find(&tasks).Error; err != nil {
// 		return nil, err
// 	}
//
// 	return tasks, nil
// }

// func (s *taskService) ListTasksByBuildingId(buildingId uint) ([]models.Task, error) {
// 	var tasks []models.Task
// 	if err := s.db.Find()
// if err := s.db.Preload("Categories").Find(&tasks).Error; err != nil {
// 	return nil, err
// }

// 	return tasks, nil
// }

func (s *taskService) CreateTask(task models.Task) (models.Task, error) {
	var count int64
	if err := s.db.Model(&models.Building{}).
		Where("id = ?", task.BuildingId).
		Count(&count).Error; err != nil {
		return models.Task{}, err
	}
	if count == 0 {
		return models.Task{}, gorm.ErrRecordNotFound
	}

	if err := s.db.Create(&task).Error; err != nil {
		return models.Task{}, err
	}
	if err := s.db.Preload("Building").First(&task, task.ID).Error; err != nil {
		return models.Task{}, err
	}

	return task, nil
}

func (s *taskService) DeleteTask(id uint) error {
	result := s.db.Delete(&models.Task{}, id)
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (s *taskService) UpdateTask(task models.Task, id uint) error {
	return s.db.Transaction(func(transaction *gorm.DB) error {
		result := transaction.
			Model(&models.Task{}).
			Where("id = ?", id).
			Updates(map[string]any{
				"name":         task.Name,
				"description":  task.Description,
				"building_id":  task.BuildingId,
				"is_completed": task.IsCompleted,
			})
		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}

		return nil

	})
}
