package repository

import (
	"fmt"

	"pairwise/internal/domain"

	"gorm.io/gorm"
)

type ScoringRepository interface {
	Create(score *domain.FibonacciScore) error
	Update(score *domain.FibonacciScore) error
	GetByID(id int) (*domain.FibonacciScore, error)
	GetByFeatureAndAttendee(featureID, attendeeID int, criterionType string) (*domain.FibonacciScore, error)
	GetByProject(projectID int, criterionType string) ([]*domain.FibonacciScore, error)
	GetByFeature(featureID int, criterionType string) ([]*domain.FibonacciScore, error)
	DeleteByID(id int) error
}

type scoringRepository struct {
	db *gorm.DB
}

func NewScoringRepository(db *gorm.DB) ScoringRepository {
	return &scoringRepository{db: db}
}

// Create inserts a new Fibonacci score (T030 - US4)
func (r *scoringRepository) Create(score *domain.FibonacciScore) error {
	if err := r.db.Create(score).Error; err != nil {
		return fmt.Errorf("failed to create fibonacci score: %w", err)
	}
	return nil
}

// Update modifies an existing Fibonacci score (T030 - US4)
func (r *scoringRepository) Update(score *domain.FibonacciScore) error {
	if err := r.db.Save(score).Error; err != nil {
		return fmt.Errorf("failed to update fibonacci score: %w", err)
	}
	return nil
}

// GetByID retrieves a score by its ID
func (r *scoringRepository) GetByID(id int) (*domain.FibonacciScore, error) {
	var score domain.FibonacciScore
	if err := r.db.First(&score, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get fibonacci score by ID: %w", err)
	}
	return &score, nil
}

// GetByFeatureAndAttendee finds a score for a specific feature/attendee/criterion combination (T030 - US4)
func (r *scoringRepository) GetByFeatureAndAttendee(featureID, attendeeID int, criterionType string) (*domain.FibonacciScore, error) {
	var score domain.FibonacciScore
	err := r.db.Where("feature_id = ? AND attendee_id = ? AND criterion_type = ?", featureID, attendeeID, criterionType).
		First(&score).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get fibonacci score by feature and attendee: %w", err)
	}
	return &score, nil
}

// GetByProject retrieves all scores for a project, optionally filtered by criterion type (T030 - US4)
func (r *scoringRepository) GetByProject(projectID int, criterionType string) ([]*domain.FibonacciScore, error) {
	var scores []*domain.FibonacciScore

	query := r.db.Joins("JOIN features ON fibonacci_scores.feature_id = features.id").
		Where("features.project_id = ?", projectID)

	if criterionType != "" {
		query = query.Where("fibonacci_scores.criterion_type = ?", criterionType)
	}

	query = query.Order("fibonacci_scores.submitted_at DESC")

	if err := query.Find(&scores).Error; err != nil {
		return nil, fmt.Errorf("failed to get fibonacci scores by project: %w", err)
	}

	return scores, nil
}

// GetByFeature retrieves all scores for a specific feature, optionally filtered by criterion type (T030 - US4)
func (r *scoringRepository) GetByFeature(featureID int, criterionType string) ([]*domain.FibonacciScore, error) {
	var scores []*domain.FibonacciScore

	query := r.db.Where("feature_id = ?", featureID)

	if criterionType != "" {
		query = query.Where("criterion_type = ?", criterionType)
	}

	query = query.Order("submitted_at DESC")

	if err := query.Find(&scores).Error; err != nil {
		return nil, fmt.Errorf("failed to get fibonacci scores by feature: %w", err)
	}

	return scores, nil
}

// DeleteByID removes a score by its ID
func (r *scoringRepository) DeleteByID(id int) error {
	if err := r.db.Delete(&domain.FibonacciScore{}, id).Error; err != nil {
		return fmt.Errorf("failed to delete fibonacci score: %w", err)
	}
	return nil
}
