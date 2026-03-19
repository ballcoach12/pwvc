package repository

import (
	"fmt"

	"pairwise/internal/domain"

	"gorm.io/gorm"
)

type ConsensusRepository interface {
	Create(consensus *domain.ConsensusScore) error
	Update(consensus *domain.ConsensusScore) error
	GetByID(id int) (*domain.ConsensusScore, error)
	GetByFeature(featureID int) (*domain.ConsensusScore, error)
	GetByProject(projectID int) ([]*domain.ConsensusScore, error)
	DeleteByFeature(featureID int) error
	GetConsensusScores(projectID int) (map[int]domain.ConsensusScore, error) // Backwards compatibility
	SaveConsensusScore(projectID int, score domain.ConsensusScore) error     // Backwards compatibility
	DeleteConsensusScore(projectID, featureID int) error                     // Backwards compatibility
}

type consensusRepository struct {
	db *gorm.DB
}

func NewConsensusRepository(db *gorm.DB) ConsensusRepository {
	return &consensusRepository{db: db}
}

// Create inserts a new consensus score (T034 - US5)
func (r *consensusRepository) Create(consensus *domain.ConsensusScore) error {
	if err := r.db.Create(consensus).Error; err != nil {
		return fmt.Errorf("failed to create consensus score: %w", err)
	}
	return nil
}

// Update modifies an existing consensus score (T034 - US5)
func (r *consensusRepository) Update(consensus *domain.ConsensusScore) error {
	if err := r.db.Save(consensus).Error; err != nil {
		return fmt.Errorf("failed to update consensus score: %w", err)
	}
	return nil
}

// GetByID retrieves a consensus score by its ID
func (r *consensusRepository) GetByID(id int) (*domain.ConsensusScore, error) {
	var consensus domain.ConsensusScore
	if err := r.db.First(&consensus, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get consensus score by ID: %w", err)
	}
	return &consensus, nil
}

// GetByFeature retrieves a consensus score by feature ID (T034 - US5)
func (r *consensusRepository) GetByFeature(featureID int) (*domain.ConsensusScore, error) {
	var consensus domain.ConsensusScore
	err := r.db.Where("feature_id = ?", featureID).First(&consensus).Error

	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to get consensus score by feature: %w", err)
	}
	return &consensus, nil
}

// GetByProject retrieves all consensus scores for a project (T034 - US5)
func (r *consensusRepository) GetByProject(projectID int) ([]*domain.ConsensusScore, error) {
	var consensus []*domain.ConsensusScore

	if err := r.db.Where("project_id = ?", projectID).Order("locked_at DESC").Find(&consensus).Error; err != nil {
		return nil, fmt.Errorf("failed to get consensus scores by project: %w", err)
	}

	return consensus, nil
}

// DeleteByFeature removes a consensus score by feature ID (T034 - US5)
func (r *consensusRepository) DeleteByFeature(featureID int) error {
	if err := r.db.Where("feature_id = ?", featureID).Delete(&domain.ConsensusScore{}).Error; err != nil {
		return fmt.Errorf("failed to delete consensus score: %w", err)
	}
	return nil
}

// GetConsensusScores retrieves consensus scores for all features in a project (backwards compatibility)
func (r *consensusRepository) GetConsensusScores(projectID int) (map[int]domain.ConsensusScore, error) {
	var consensusList []*domain.ConsensusScore
	if err := r.db.Where("project_id = ?", projectID).Find(&consensusList).Error; err != nil {
		return nil, fmt.Errorf("failed to get consensus scores: %w", err)
	}

	scores := make(map[int]domain.ConsensusScore)
	for _, consensus := range consensusList {
		scores[consensus.FeatureID] = *consensus
	}

	return scores, nil
}

// SaveConsensusScore saves or updates a consensus score for a feature (backwards compatibility)
func (r *consensusRepository) SaveConsensusScore(projectID int, score domain.ConsensusScore) error {
	score.ProjectID = projectID

	// Check if consensus exists
	existing, err := r.GetByFeature(score.FeatureID)
	if err != nil {
		return err
	}

	if existing == nil {
		return r.Create(&score)
	} else {
		existing.SValue = score.SValue
		existing.SComplexity = score.SComplexity
		existing.LockedBy = score.LockedBy
		existing.LockedAt = score.LockedAt
		existing.Rationale = score.Rationale
		return r.Update(existing)
	}
}

// DeleteConsensusScore removes a consensus score for a feature (backwards compatibility)
func (r *consensusRepository) DeleteConsensusScore(projectID, featureID int) error {
	if err := r.db.Where("project_id = ? AND feature_id = ?", projectID, featureID).Delete(&domain.ConsensusScore{}).Error; err != nil {
		return fmt.Errorf("failed to delete consensus score: %w", err)
	}
	return nil
}
