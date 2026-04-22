package user

import (
	"errors"

	"github.com/gilabs/crm-healthcare/api/internal/domain/activity"
	"github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	"github.com/gilabs/crm-healthcare/api/internal/domain/task"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"github.com/gilabs/crm-healthcare/api/internal/domain/visit_report"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)


type ProfileService struct {
	userRepo       interfaces.UserRepository
	activityRepo   interfaces.ActivityRepository
	dealRepo       interfaces.DealRepository
	visitReportRepo interfaces.VisitReportRepository
	taskRepo       interfaces.TaskRepository
}

func NewProfileService(
	userRepo interfaces.UserRepository,
	activityRepo interfaces.ActivityRepository,
	dealRepo interfaces.DealRepository,
	visitReportRepo interfaces.VisitReportRepository,
	taskRepo interfaces.TaskRepository,
) *ProfileService {
	return &ProfileService{
		userRepo:        userRepo,
		activityRepo:    activityRepo,
		dealRepo:        dealRepo,
		visitReportRepo: visitReportRepo,
		taskRepo:        taskRepo,
	}
}

// GetProfile returns complete profile data for a user
func (s *ProfileService) GetProfile(userID string) (*user.ProfileResponse, error) {
	// Get user
	u, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Get stats
	stats, err := s.getStats(userID)
	if err != nil {
		return nil, err
	}

	// Get activities
	activities, err := s.getActivities(userID)
	if err != nil {
		return nil, err
	}

	// Get transactions (from deals)
	transactions, err := s.getTransactions(userID)
	if err != nil {
		return nil, err
	}

	return &user.ProfileResponse{
		User:         u.ToUserResponse(),
		Stats:        stats,
		Activities:   activities,
		Transactions: transactions,
	}, nil
}

// getStats calculates profile statistics
func (s *ProfileService) getStats(userID string) (*user.ProfileStats, error) {
	// Count visit reports created by user
	visitReq := &visit_report.ListVisitReportsRequest{
		SalesRepID: userID,
		Page:       1,
		PerPage:    1,
	}
	_, totalVisits, err := s.visitReportRepo.List(visitReq)
	if err != nil {
		return nil, err
	}

	// Count deals assigned to user
	dealReq := &pipeline.ListDealsRequest{
		AssignedTo: userID,
		Page:       1,
		PerPage:    1,
	}
	_, totalDeals, err := s.dealRepo.List(dealReq)
	if err != nil {
		return nil, err
	}

	// Count tasks assigned to user
	taskReq := &task.ListTasksRequest{
		AssignedTo: userID,
		Page:       1,
		PerPage:    1,
	}
	_, totalTasks, err := s.taskRepo.List(taskReq)
	if err != nil {
		return nil, err
	}

	return &user.ProfileStats{
		Visits: int(totalVisits),
		Deals:  int(totalDeals),
		Tasks:  int(totalTasks),
	}, nil
}

// getActivities returns recent activities for user
func (s *ProfileService) getActivities(userID string) ([]user.ProfileActivity, error) {
	req := &activity.ListActivitiesRequest{
		UserID:  userID,
		Page:    1,
		PerPage: 3,
	}

	activities, _, err := s.activityRepo.List(req)
	if err != nil {
		return nil, err
	}

	result := make([]user.ProfileActivity, len(activities))
	for i, a := range activities {
		// Generate activity title based on type
		title := s.generateActivityTitle(&a)
		
		result[i] = user.ProfileActivity{
			ID:          a.ID,
			Title:       title,
			Description: a.Description,
			Type:        a.Type,
			Date:        a.Timestamp,
		}
	}

	return result, nil
}

// generateActivityTitle generates a title for activity based on type
func (s *ProfileService) generateActivityTitle(a *activity.Activity) string {
	switch a.Type {
	case "visit":
		return "Visit Report"
	case "call":
		return "Phone Call"
	case "email":
		return "Email Sent"
	case "task":
		return "Task Created"
	case "deal":
		return "Deal Updated"
	default:
		return "Activity"
	}
}

// getTransactions returns transactions from deals (as placeholder)
func (s *ProfileService) getTransactions(userID string) ([]user.ProfileTransaction, error) {
	req := &pipeline.ListDealsRequest{
		AssignedTo: userID,
		Page:       1,
		PerPage:    10,
	}

	deals, _, err := s.dealRepo.List(req)
	if err != nil {
		return nil, err
	}

	result := make([]user.ProfileTransaction, 0, len(deals))
	for _, deal := range deals {
		// Map deal status to transaction status
		status := "pending"
		if deal.Status == "won" {
			status = "paid"
		} else if deal.Status == "lost" {
			status = "failed"
		}

		result = append(result, user.ProfileTransaction{
			ID:      deal.ID,
			Product: deal.Title,
			Status:  status,
			Date:    deal.CreatedAt,
			Amount:  deal.Value,
		})
	}

	return result, nil
}

// UpdateProfile updates user profile information
func (s *ProfileService) UpdateProfile(userID string, req *user.UpdateProfileRequest) (*user.UserResponse, error) {
	// Find user
	u, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound
		}
		return nil, err
	}

	// Update name if provided
	if req.Name != "" {
		u.Name = req.Name
	}

	if err := s.userRepo.Update(u); err != nil {
		return nil, err
	}

	// Reload with role
	updatedUser, err := s.userRepo.FindByID(u.ID)
	if err != nil {
		return nil, err
	}

	return updatedUser.ToUserResponse(), nil
}

// ChangePassword changes user password
func (s *ProfileService) ChangePassword(userID string, req *user.ChangePasswordRequest) error {
	// Validate password confirmation
	if req.Password != req.ConfirmPassword {
		return errors.New("passwords do not match")
	}

	// Find user
	u, err := s.userRepo.FindByID(userID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrUserNotFound
		}
		return err
	}

	// Verify current password
	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.CurrentPassword))
	if err != nil {
		return ErrIncorrectPassword
	}
	
	// Hash new password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	u.Password = string(hashedPassword)

	return s.userRepo.Update(u)
}

