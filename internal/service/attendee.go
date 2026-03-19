package service

import (
	"crypto/sha256"
	"fmt"

	"pairwise/internal/domain"
	"pairwise/internal/repository"
)

// AttendeeService handles business logic for attendees
type AttendeeService struct {
	attendeeRepo repository.AttendeeRepository
}

// NewAttendeeService creates a new attendee service
func NewAttendeeService(attendeeRepo repository.AttendeeRepository) *AttendeeService {
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

// SetPIN sets or updates the PIN for an attendee
func (s *AttendeeService) SetPIN(id int, req domain.SetPINRequest) error {
	if id <= 0 {
		return domain.NewAPIError(400, "Invalid attendee ID")
	}

	// Basic validation
	if req.PIN == "" {
		return domain.NewAPIError(400, "PIN is required")
	}

	if len(req.PIN) < 4 || len(req.PIN) > 20 {
		return domain.NewAPIError(400, "PIN must be between 4 and 20 characters")
	}

	// Verify attendee exists
	_, err := s.attendeeRepo.GetByID(id)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.NewAPIError(404, "Attendee not found")
		}
		return domain.NewAPIError(500, "Failed to verify attendee", err.Error())
	}

	// Set the PIN
	err = s.attendeeRepo.SetPIN(id, req.PIN)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.NewAPIError(404, "Attendee not found")
		}
		return domain.NewAPIError(500, "Failed to set PIN", err.Error())
	}

	return nil
}

// CreateAttendeeWithoutPIN creates a new attendee without requiring a PIN (for invite workflow)
func (s *AttendeeService) CreateAttendeeWithoutPIN(projectID int, req domain.CreateAttendeeWithoutPINRequest) (*domain.Attendee, error) {
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
	attendee, err := s.attendeeRepo.CreateWithoutPIN(projectID, req)
	if err != nil {
		return nil, domain.NewAPIError(500, "Failed to create attendee", err.Error())
	}

	return attendee, nil
}

// GenerateInviteToken generates an invite token for an attendee to set up their PIN
func (s *AttendeeService) GenerateInviteToken(attendeeID int) (*domain.ResetPINResponse, error) {
	if attendeeID <= 0 {
		return nil, domain.NewAPIError(400, "Invalid attendee ID")
	}

	// Verify attendee exists
	_, err := s.attendeeRepo.GetByID(attendeeID)
	if err != nil {
		if err == domain.ErrNotFound {
			return nil, domain.NewAPIError(404, "Attendee not found")
		}
		return nil, domain.NewAPIError(500, "Failed to verify attendee", err.Error())
	}

	// Generate the token
	token, err := s.attendeeRepo.GenerateInviteToken(attendeeID)
	if err != nil {
		return nil, domain.NewAPIError(500, "Failed to generate invite token", err.Error())
	}

	return &domain.ResetPINResponse{
		InviteToken: token,
		ExpiresAt:   "24 hours",
	}, nil
}

// SetupPIN allows an attendee to set their PIN using an invite token
func (s *AttendeeService) SetupPIN(req domain.SetupPINRequest) error {
	// Validate the invite token
	attendee, err := s.attendeeRepo.GetByInviteToken(req.InviteToken)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.NewAPIError(400, "Invalid or expired invite token")
		}
		return domain.NewAPIError(500, "Failed to validate invite token", err.Error())
	}

	// Verify the attendee ID matches
	if attendee.ID != req.AttendeeID {
		return domain.NewAPIError(400, "Attendee ID does not match invite token")
	}

	// Set the PIN
	err = s.attendeeRepo.SetPIN(req.AttendeeID, req.PIN)
	if err != nil {
		return domain.NewAPIError(500, "Failed to set PIN", err.Error())
	}

	// Clear the invite token
	err = s.attendeeRepo.ClearInviteToken(req.AttendeeID)
	if err != nil {
		// Log the error but don't fail the request since PIN was set successfully
		// TODO: Add proper logging
	}

	return nil
}

// ChangePIN allows an authenticated attendee to change their own PIN
func (s *AttendeeService) ChangePIN(attendeeID int, req domain.ChangePINRequest) error {
	if attendeeID <= 0 {
		return domain.NewAPIError(400, "Invalid attendee ID")
	}

	// Get current attendee data
	attendee, err := s.attendeeRepo.GetByID(attendeeID)
	if err != nil {
		if err == domain.ErrNotFound {
			return domain.NewAPIError(404, "Attendee not found")
		}
		return domain.NewAPIError(500, "Failed to retrieve attendee", err.Error())
	}

	// Verify current PIN (simple hash comparison)
	currentPINHash := hashPINInService(req.CurrentPIN)
	if attendee.PinHash == nil || *attendee.PinHash != currentPINHash {
		return domain.NewAPIError(401, "Current PIN is incorrect")
	}

	// Set the new PIN
	err = s.attendeeRepo.SetPIN(attendeeID, req.NewPIN)
	if err != nil {
		return domain.NewAPIError(500, "Failed to update PIN", err.Error())
	}

	return nil
}

// hashPINInService creates a simple hash of the PIN (service layer implementation)
func hashPINInService(pin string) string {
	hash := sha256.Sum256([]byte(pin))
	return fmt.Sprintf("%x", hash)
}
