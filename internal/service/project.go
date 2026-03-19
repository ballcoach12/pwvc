package service

import (
	"pairwise/internal/domain"
	"pairwise/internal/repository"
)

// ProjectService handles business logic for projects
type ProjectService struct {
	projectRepo repository.ProjectRepository
}

// NewProjectService creates a new project service
func NewProjectService(projectRepo repository.ProjectRepository) *ProjectService {
	return &ProjectService{
		projectRepo: projectRepo,
	}
}

// CreateProject creates a new project with validation
func (s *ProjectService) CreateProject(req domain.CreateProjectRequest) (*domain.Project, error) {
	// Basic validation
	if req.Name == "" {
		return nil, domain.NewAPIError(400, "Project name is required")
	}

	if len(req.Name) > 255 {
		return nil, domain.NewAPIError(400, "Project name must be less than 255 characters")
	}

	// Create the project
	project, err := s.projectRepo.Create(req)
	if err != nil {
		return nil, domain.NewAPIError(500, "Failed to create project", err.Error())
	}

	return project, nil
}

// GetProject retrieves a project by ID
func (s *ProjectService) GetProject(id int) (*domain.Project, error) {
	if id <= 0 {
		return nil, domain.NewAPIError(400, "Invalid project ID")
	}

	project, err := s.projectRepo.GetByID(id)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.NewAPIError(404, "Project not found")
		}
		return nil, domain.NewAPIError(500, "Failed to retrieve project", err.Error())
	}

	return project, nil
}

// UpdateProject updates an existing project
func (s *ProjectService) UpdateProject(id int, req domain.UpdateProjectRequest) (*domain.Project, error) {
	if id <= 0 {
		return nil, domain.NewAPIError(400, "Invalid project ID")
	}

	// Basic validation
	if req.Name == "" {
		return nil, domain.NewAPIError(400, "Project name is required")
	}

	if len(req.Name) > 255 {
		return nil, domain.NewAPIError(400, "Project name must be less than 255 characters")
	}

	// Validate status if provided
	if req.Status != "" {
		validStatuses := map[string]bool{
			"active":    true,
			"inactive":  true,
			"completed": true,
		}
		if !validStatuses[req.Status] {
			return nil, domain.NewAPIError(400, "Invalid status. Must be one of: active, inactive, completed")
		}
	}

	project, err := s.projectRepo.Update(id, req)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.NewAPIError(404, "Project not found")
		}
		return nil, domain.NewAPIError(500, "Failed to update project", err.Error())
	}

	return project, nil
}

// DeleteProject deletes a project
func (s *ProjectService) DeleteProject(id int) error {
	if id <= 0 {
		return domain.NewAPIError(400, "Invalid project ID")
	}

	err := s.projectRepo.Delete(id)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.NewAPIError(404, "Project not found")
		}
		return domain.NewAPIError(500, "Failed to delete project", err.Error())
	}

	return nil
}

// ListProjects retrieves all projects
func (s *ProjectService) ListProjects() ([]domain.Project, error) {
	projects, err := s.projectRepo.List()
	if err != nil {
		return nil, domain.NewAPIError(500, "Failed to retrieve projects", err.Error())
	}

	// Return empty slice instead of nil if no projects found
	if projects == nil {
		projects = []domain.Project{}
	}

	return projects, nil
}

// UpdateProjectInviteCode updates a project's invite code (T016 - US1)
func (s *ProjectService) UpdateProjectInviteCode(id int, inviteCode string) (*domain.Project, error) {
	if id <= 0 {
		return nil, domain.NewAPIError(400, "Invalid project ID")
	}

	project, err := s.projectRepo.UpdateInviteCode(id, inviteCode)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.NewAPIError(404, "Project not found")
		}
		return nil, domain.NewAPIError(500, "Failed to update invite code", err.Error())
	}

	return project, nil
}

// GetProjectByInviteCode retrieves a project by its invite code (T016 - US1)
func (s *ProjectService) GetProjectByInviteCode(inviteCode string) (*domain.Project, error) {
	if inviteCode == "" {
		return nil, domain.NewAPIError(400, "Invite code is required")
	}

	project, err := s.projectRepo.GetByInviteCode(inviteCode)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.NewAPIError(404, "Invalid invite code")
		}
		return nil, domain.NewAPIError(500, "Failed to retrieve project", err.Error())
	}

	return project, nil
}
