package birthdayservice

import (
	"context"
	"fmt"
	"time"

	"github.com/kianooshaz/skeleton/foundation/derror"
	"github.com/kianooshaz/skeleton/foundation/pagination"
	"github.com/kianooshaz/skeleton/services/user/birthday/persistence"
	birthdayproto "github.com/kianooshaz/skeleton/services/user/birthday/proto"
)

// calculateAge calculates age based on the date of birth.
func calculateAge(dateOfBirth time.Time) int {
	now := time.Now()
	age := now.Year() - dateOfBirth.Year()

	// Adjust if birthday hasn't occurred this year yet.
	if now.YearDay() < dateOfBirth.YearDay() {
		age--
	}

	return age
}

// validateAge validates that the age is within acceptable bounds.
func (s *Service) validateAge(age int) error {
	if age < s.config.MinAge {
		return fmt.Errorf("age %d is below minimum allowed age %d", age, s.config.MinAge)
	}
	if age > s.config.MaxAge {
		return fmt.Errorf("age %d is above maximum allowed age %d", age, s.config.MaxAge)
	}
	return nil
}

// Create creates a new birthday record.
func (s *Service) Create(ctx context.Context, req birthdayproto.CreateRequest) (birthdayproto.CreateResponse, error) {
	s.logger.Info("Creating birthday record", "user_id", req.UserID)

	// Check if birthday already exists for this user.
	exists, err := s.persister.ExistsByUserID(ctx, req.UserID)
	if err != nil {
		return birthdayproto.CreateResponse{}, fmt.Errorf("checking existing birthday: %w", err)
	}
	if exists {
		return birthdayproto.CreateResponse{}, derror.ErrUserAlreadyExists
	}

	// Calculate age.
	age := calculateAge(req.DateOfBirth)

	// Validate age.
	if err := s.validateAge(age); err != nil {
		return birthdayproto.CreateResponse{}, fmt.Errorf("validating age: %w", err)
	}

	// Create birthday record.
	birthday := birthdayproto.Birthday{
		ID:          birthdayproto.NewBirthdayID(),
		UserID:      req.UserID,
		DateOfBirth: req.DateOfBirth,
		Age:         age,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.persister.Create(ctx, birthday); err != nil {
		return birthdayproto.CreateResponse{}, fmt.Errorf("creating birthday record: %w", err)
	}

	s.logger.Info("Birthday record created successfully", "birthday_id", birthday.ID, "user_id", req.UserID)

	return birthdayproto.CreateResponse{Data: birthday}, nil
}

// Get retrieves a birthday record by ID.
func (s *Service) Get(ctx context.Context, req birthdayproto.GetRequest) (birthdayproto.GetResponse, error) {
	s.logger.Info("Getting birthday record", "birthday_id", req.ID)

	birthday, err := s.persister.Get(ctx, req.ID)
	if err != nil {
		return birthdayproto.GetResponse{}, fmt.Errorf("getting birthday record: %w", err)
	}

	return birthdayproto.GetResponse{Data: birthday}, nil
}

// GetByUserID retrieves a birthday record by user ID.
func (s *Service) GetByUserID(ctx context.Context, req birthdayproto.GetByUserIDRequest) (birthdayproto.GetByUserIDResponse, error) {
	s.logger.Info("Getting birthday record by user ID", "user_id", req.UserID)

	birthday, err := s.persister.GetByUserID(ctx, req.UserID)
	if err != nil {
		return birthdayproto.GetByUserIDResponse{}, fmt.Errorf("getting birthday record by user ID: %w", err)
	}

	return birthdayproto.GetByUserIDResponse{Data: birthday}, nil
}

// Update updates an existing birthday record.
func (s *Service) Update(ctx context.Context, req birthdayproto.UpdateRequest) (birthdayproto.UpdateResponse, error) {
	s.logger.Info("Updating birthday record", "birthday_id", req.ID)

	// Get existing birthday record.
	existingBirthday, err := s.persister.Get(ctx, req.ID)
	if err != nil {
		return birthdayproto.UpdateResponse{}, fmt.Errorf("getting existing birthday record: %w", err)
	}

	// Calculate new age.
	age := calculateAge(req.DateOfBirth)

	// Validate age.
	if err := s.validateAge(age); err != nil {
		return birthdayproto.UpdateResponse{}, fmt.Errorf("validating age: %w", err)
	}

	// Update birthday record.
	updatedBirthday := birthdayproto.Birthday{
		ID:          existingBirthday.ID,
		UserID:      existingBirthday.UserID,
		DateOfBirth: req.DateOfBirth,
		Age:         age,
		CreatedAt:   existingBirthday.CreatedAt,
		UpdatedAt:   time.Now(),
	}

	if err := s.persister.Update(ctx, updatedBirthday); err != nil {
		return birthdayproto.UpdateResponse{}, fmt.Errorf("updating birthday record: %w", err)
	}

	s.logger.Info("Birthday record updated successfully", "birthday_id", req.ID)

	return birthdayproto.UpdateResponse{Data: updatedBirthday}, nil
}

// Delete removes a birthday record.
func (s *Service) Delete(ctx context.Context, req birthdayproto.DeleteRequest) error {
	s.logger.Info("Deleting birthday record", "birthday_id", req.ID)

	if err := s.persister.Delete(ctx, req.ID); err != nil {
		return fmt.Errorf("deleting birthday record: %w", err)
	}

	s.logger.Info("Birthday record deleted successfully", "birthday_id", req.ID)

	return nil
}

// List retrieves a paginated list of birthday records with optional filters.
func (s *Service) List(ctx context.Context, req birthdayproto.ListRequest) (birthdayproto.ListResponse, error) {
	s.logger.Info("Listing birthday records", "page_number", req.PageNumber, "page_rows", req.PageRows)

	// Convert service filters to persistence filters.
	filters := persistence.ListFilters{
		UserID:     req.UserID,
		MinAge:     req.MinAge,
		MaxAge:     req.MaxAge,
		BirthMonth: req.BirthMonth,
	}

	// Get birthday records.
	birthdays, err := s.persister.List(ctx, req.Page, req.OrderBy, filters)
	if err != nil {
		return birthdayproto.ListResponse{}, fmt.Errorf("listing birthday records: %w", err)
	}

	// Get total count.
	totalCount, err := s.persister.Count(ctx, filters)
	if err != nil {
		return birthdayproto.ListResponse{}, fmt.Errorf("counting birthday records: %w", err)
	}

	// Build response.
	response := pagination.NewResponse(req.Page, totalCount, birthdays)

	return birthdayproto.ListResponse(response), nil
}
