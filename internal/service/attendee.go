package service

import (
	"pairwise/internal/domain"
	"pairwise/internal/repository"
)

// AttendeeService handles business logic for attendees
type AttendeeService struct {
	attendeeRepo *repository.AttendeeRepository
}

// NewAttendeeService creates a new attendee service
func NewAttendeeService(attendeeRepo *repository.AttendeeRepository) *AttendeeService {
	return &AttendeeService{
		attendeeRepo: attendeeRepo,
	}
}

// CreateAttendee creates a new attendee for a project
func (s *AttendeeService) CreateAttendee(projectID int, req domain.CreateAttendeeRequest) (*domain.Attendee, error) {
	if projectID <= 0 {
		return nil, domain.NewAPIError(400, "Invalid project ID")
	}

	// Basic validation
	if req.Name == "" {
		return nil, domain.NewAPIError(400, "Attendee name is required")
	}

	if len(req.Name) > 255 {
		return nil, domain.NewAPIError(400, "Attendee name must be less than 255 characters")
	}

	// Create the attendee
	attendee, err := s.attendeeRepo.Create(projectID, req)
	if err != nil {
		return nil, domain.NewAPIError(500, "Failed to create attendee", err.Error())
	}

	return attendee, nil
}

// GetAttendee retrieves an attendee by ID
func (s *AttendeeService) GetAttendee(id int) (*domain.Attendee, error) {
	if id <= 0 {
		return nil, domain.NewAPIError(400, "Invalid attendee ID")
	}

	attendee, err := s.attendeeRepo.GetByID(id)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.NewAPIError(404, "Attendee not found")
		}
		return nil, domain.NewAPIError(500, "Failed to retrieve attendee", err.Error())
	}

	return attendee, nil
}

// GetProjectAttendees retrieves all attendees for a project
func (s *AttendeeService) GetProjectAttendees(projectID int) ([]domain.Attendee, error) {
	if projectID <= 0 {
		return nil, domain.NewAPIError(400, "Invalid project ID")
	}

	attendees, err := s.attendeeRepo.GetByProjectID(projectID)
	if err != nil {
		return nil, domain.NewAPIError(500, "Failed to retrieve attendees", err.Error())
	}

	// Return empty slice instead of nil if no attendees found
	if attendees == nil {
		attendees = []domain.Attendee{}
	}

	return attendees, nil
}

// DeleteAttendee deletes an attendee
func (s *AttendeeService) DeleteAttendee(id int) error {
	if id <= 0 {
		return domain.NewAPIError(400, "Invalid attendee ID")
	}

	err := s.attendeeRepo.Delete(id)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.NewAPIError(404, "Attendee not found")
		}
		return domain.NewAPIError(500, "Failed to delete attendee", err.Error())
	}

	return nil
}
