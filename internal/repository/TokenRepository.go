package repository

import (
	"test_go_project/internal/entity"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type TokenRepo struct {
	db *gorm.DB
}

func NewTokenRepo(db *gorm.DB) *TokenRepo {
	return &TokenRepo{db}
}

func (r *TokenRepo) SaveRefreshPair(refreshRecord *entity.RefreshRecord) error {
	result := r.db.Clauses(clause.OnConflict{
		UpdateAll: true,
	}).Create(refreshRecord)

	return result.Error
}

func (r *TokenRepo) ExistsById(id uuid.UUID) (bool, error) {
	var exists bool
	err := r.db.Model(&entity.RefreshRecord{}).
		Select("1").
		Where("user_id = ?", id).
		Limit(1).
		Find(&exists).Error

	if err != nil {
		return false, err
	}

	return exists, nil
}

func (r *TokenRepo) GetById(id uuid.UUID) (*entity.RefreshRecord, error) {
	var refreshRecord entity.RefreshRecord
	err := r.db.First(&refreshRecord, "user_id = ?", id).Error

	if err != nil {
		return nil, err
	}

	return &refreshRecord, nil
}

func (r *TokenRepo) DeleteById(id uuid.UUID) error {
	result := r.db.Delete(&entity.RefreshRecord{}, id)

	return result.Error
}
