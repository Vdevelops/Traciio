package seeders

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/account"
	"github.com/gilabs/crm-healthcare/api/internal/domain/activity"
	"github.com/gilabs/crm-healthcare/api/internal/domain/activity_type"
	"github.com/gilabs/crm-healthcare/api/internal/domain/contact"
	"github.com/gilabs/crm-healthcare/api/internal/domain/pipeline"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"github.com/gilabs/crm-healthcare/api/internal/domain/visit_report"
	"gorm.io/datatypes"
)

// SeedActivities seeds initial activities
func SeedActivities() error {
	// Check if activities already exist
	var count int64
	database.DB.Model(&activity.Activity{}).Count(&count)
	if count > 0 {
		log.Println("Activities already seeded, skipping...")
		return nil
	}

	// Get accounts
	var accounts []account.Account
	if err := database.DB.Find(&accounts).Error; err != nil {
		return err
	}
	if len(accounts) == 0 {
		log.Println("Warning: No accounts found, skipping activity seeding")
		return nil
	}

	// Get contacts
	var contacts []contact.Contact
	if err := database.DB.Find(&contacts).Error; err != nil {
		return err
	}

	// Get users
	var users []user.User
	if err := database.DB.Find(&users).Error; err != nil {
		return err
	}
	if len(users) == 0 {
		log.Println("Warning: No users found, skipping activity seeding")
		return nil
	}

	// Get activity types
	var activityTypes []activity_type.ActivityType
	if err := database.DB.Find(&activityTypes).Error; err != nil {
		return err
	}

	// Create map of activity types by code
	activityTypeMap := make(map[string]string)
	for _, at := range activityTypes {
		activityTypeMap[at.Code] = at.ID
	}

	// Get visit reports to create related activities
	var visitReports []visit_report.VisitReport
	if err := database.DB.Find(&visitReports).Error; err != nil {
		return err
	}

	// Helper function to marshal metadata
	marshalMetadata := func(data map[string]interface{}) datatypes.JSON {
		bytes, _ := json.Marshal(data)
		return bytes
	}

	now := time.Now()
	activities := []activity.Activity{}

	// Create activities from visit reports
	for _, vr := range visitReports {
		// Get visit activity type ID
		visitTypeID := activityTypeMap["visit"]
		
		// Activity for visit report creation
		activities = append(activities, activity.Activity{
			Type:           "visit",
			ActivityTypeID: &visitTypeID,
			AccountID:      vr.AccountID, // Already *string, no need for &
			ContactID:      vr.ContactID,
			DealID:         vr.DealID,    // Include DealID if exists
			LeadID:         vr.LeadID,    // Include LeadID if exists
			UserID:         vr.SalesRepID,
			Description:    "Visit report created: " + vr.Purpose,
			Timestamp:      vr.CreatedAt,
			Metadata: marshalMetadata(map[string]interface{}{
				"visit_report_id": vr.ID,
				"status":         vr.Status,
				"visit_date":     vr.VisitDate.Format("2006-01-02"),
			}),
		})

		// Activity for check-in if exists
		if vr.CheckInTime != nil {
			activities = append(activities, activity.Activity{
				Type:           "visit",
				ActivityTypeID: &visitTypeID,
				AccountID:      vr.AccountID, // Already *string, no need for &
				ContactID:      vr.ContactID,
				DealID:         vr.DealID,    // Include DealID if exists
				LeadID:         vr.LeadID,   // Include LeadID if exists
				UserID:         vr.SalesRepID,
				Description:    "Checked in for visit: " + vr.Purpose,
				Timestamp:      *vr.CheckInTime,
				Metadata: marshalMetadata(map[string]interface{}{
					"visit_report_id": vr.ID,
					"action":          "check_in",
					"location":        "GPS location recorded",
				}),
			})
		}

		// Activity for check-out if exists
		if vr.CheckOutTime != nil {
			activities = append(activities, activity.Activity{
				Type:           "visit",
				ActivityTypeID: &visitTypeID,
				AccountID:      vr.AccountID, // Already *string, no need for &
				ContactID:      vr.ContactID,
				DealID:         vr.DealID,    // Include DealID if exists
				LeadID:         vr.LeadID,    // Include LeadID if exists
				UserID:         vr.SalesRepID,
				Description:    "Checked out from visit: " + vr.Purpose,
				Timestamp:      *vr.CheckOutTime,
				Metadata: marshalMetadata(map[string]interface{}{
					"visit_report_id": vr.ID,
					"action":          "check_out",
					"location":        "GPS location recorded",
				}),
			})
		}

		// Activity for approval if exists
		if vr.Status == "approved" && vr.ApprovedAt != nil {
			activities = append(activities, activity.Activity{
				Type:           "visit",
				ActivityTypeID: &visitTypeID,
				AccountID:      vr.AccountID, // Already *string, no need for &
				ContactID:      vr.ContactID,
				DealID:         vr.DealID,     // Include DealID if exists
				LeadID:         vr.LeadID,     // Include LeadID if exists
				UserID:         vr.SalesRepID,
				Description:    "Visit report approved: " + vr.Purpose,
				Timestamp:      *vr.ApprovedAt,
				Metadata: marshalMetadata(map[string]interface{}{
					"visit_report_id": vr.ID,
					"action":          "approved",
					"approved_by":     vr.ApprovedBy,
				}),
			})
		}

		// Activity for rejection if exists
		if vr.Status == "rejected" && vr.ApprovedAt != nil {
			activities = append(activities, activity.Activity{
				Type:           "visit",
				ActivityTypeID: &visitTypeID,
				AccountID:      vr.AccountID, // Already *string, no need for &
				ContactID:      vr.ContactID,
				DealID:         vr.DealID,    // Include DealID if exists
				LeadID:         vr.LeadID,    // Include LeadID if exists
				UserID:         vr.SalesRepID,
				Description:    "Visit report rejected: " + vr.Purpose,
				Timestamp:      *vr.ApprovedAt,
				Metadata: marshalMetadata(map[string]interface{}{
					"visit_report_id": vr.ID,
					"action":          "rejected",
					"rejection_reason": vr.RejectionReason,
				}),
			})
		}
	}

	// Additional activities (calls, emails, tasks)
	if len(accounts) > 0 && len(users) > 0 {
		var contactID *string
		if len(contacts) > 0 {
			contactID = &contacts[0].ID
		}

		// Get activity type IDs
		callTypeID := activityTypeMap["call"]
		emailTypeID := activityTypeMap["email"]
		taskTypeID := activityTypeMap["task"]
		dealTypeID := activityTypeMap["deal"]

		// Call activity
		activities = append(activities, activity.Activity{
			Type:           "call",
			ActivityTypeID: &callTypeID,
			AccountID:      &accounts[0].ID,
			ContactID:      contactID,
			UserID:         users[0].ID,
			Description:    "Follow-up call regarding product inquiry",
			Timestamp:      now.Add(-3 * 24 * time.Hour), // 3 days ago
			Metadata: marshalMetadata(map[string]interface{}{
				"duration": "15 minutes",
				"outcome":  "Positive response, scheduled meeting",
			}),
		})

		// Email activity
		activities = append(activities, activity.Activity{
			Type:           "email",
			ActivityTypeID: &emailTypeID,
			AccountID:      &accounts[0].ID,
			ContactID:      contactID,
			UserID:         users[0].ID,
			Description:    "Sent product catalog and pricing information",
			Timestamp:      now.Add(-2 * 24 * time.Hour), // 2 days ago
			Metadata: marshalMetadata(map[string]interface{}{
				"subject": "Product Catalog - Q1 2024",
				"status":  "sent",
			}),
		})

		// Task activity
		if len(accounts) > 1 {
			var taskContactID *string
			if len(contacts) > 1 {
				taskContactID = &contacts[1].ID
			}
			activities = append(activities, activity.Activity{
				Type:           "task",
				ActivityTypeID: &taskTypeID,
				AccountID:      &accounts[1].ID,
				ContactID:      taskContactID,
				UserID:         users[0].ID,
				Description:    "Prepare proposal for new product line",
				Timestamp:      now.Add(-1 * 24 * time.Hour), // 1 day ago
				Metadata: marshalMetadata(map[string]interface{}{
					"priority": "high",
					"due_date": now.Add(2 * 24 * time.Hour).Format("2006-01-02"),
					"status":   "in_progress",
				}),
			})
		}

		// Deal activity
		if len(accounts) > 2 {
			var dealContactID *string
			if len(contacts) > 2 {
				dealContactID = &contacts[2].ID
			}
			activities = append(activities, activity.Activity{
				Type:           "deal",
				ActivityTypeID: &dealTypeID,
				AccountID:      &accounts[2].ID,
				ContactID:      dealContactID,
				UserID:         users[0].ID,
				Description:    "New deal opportunity: Annual supply contract",
				Timestamp:      now.Add(-5 * 24 * time.Hour), // 5 days ago
				Metadata: marshalMetadata(map[string]interface{}{
					"value":     500000000,
					"currency":  "IDR",
					"stage":     "negotiation",
					"probability": 0.7,
				}),
			})
		}
	}

	// Create additional activities specifically for deals
	var deals []pipeline.Deal
	if err := database.DB.Where("status IN (?)", []string{"open", "won"}).Limit(20).Find(&deals).Error; err != nil {
		log.Printf("Warning: Failed to get deals for activities: %v", err)
	} else if len(deals) > 0 {
		log.Printf("Creating activities linked to %d deals", len(deals))

		// Get activity type IDs
		callTypeID := activityTypeMap["call"]
		emailTypeID := activityTypeMap["email"]
		dealTypeID := activityTypeMap["deal"]

		for i, deal := range deals {
			// Get user assigned to deal or use first user
			var userID string
			if deal.AssignedTo != nil {
				userID = *deal.AssignedTo
			} else if len(users) > 0 {
				userID = users[0].ID
			} else {
				continue
			}

			// Create 2-4 activities per deal
			// Call activity
			activities = append(activities, activity.Activity{
				Type:           "call",
				ActivityTypeID: &callTypeID,
				AccountID:      &deal.AccountID,
				ContactID:      deal.ContactID,
				DealID:         &deal.ID, // Link to deal
				UserID:         userID,
				Description:    fmt.Sprintf("Follow-up call for deal: %s", deal.Title),
				Timestamp:      now.Add(-time.Duration(i*2+1) * 24 * time.Hour),
				Metadata: marshalMetadata(map[string]interface{}{
					"deal_id":  deal.ID,
					"duration": "20 minutes",
					"outcome":  "Positive discussion about deal progress",
				}),
			})

			// Email activity
			activities = append(activities, activity.Activity{
				Type:           "email",
				ActivityTypeID: &emailTypeID,
				AccountID:      &deal.AccountID,
				ContactID:      deal.ContactID,
				DealID:         &deal.ID, // Link to deal
				UserID:         userID,
				Description:    fmt.Sprintf("Sent proposal for deal: %s", deal.Title),
				Timestamp:      now.Add(-time.Duration(i*2) * 24 * time.Hour),
				Metadata: marshalMetadata(map[string]interface{}{
					"deal_id": deal.ID,
					"subject": fmt.Sprintf("Proposal for %s", deal.Title),
					"status":  "sent",
				}),
			})

			// Deal stage change activity (for some deals)
			if i%2 == 0 {
				activities = append(activities, activity.Activity{
					Type:           "deal",
					ActivityTypeID: &dealTypeID,
					AccountID:      &deal.AccountID,
					ContactID:      deal.ContactID,
					DealID:         &deal.ID, // Link to deal
					UserID:         userID,
					Description:    fmt.Sprintf("Deal stage updated: %s", deal.Title),
					Timestamp:      now.Add(-time.Duration(i*3) * 24 * time.Hour),
					Metadata: marshalMetadata(map[string]interface{}{
						"deal_id":     deal.ID,
						"stage":       "negotiation",
						"value":       deal.Value,
						"probability": deal.Probability,
					}),
				})
			}
		}

		log.Printf("Added %d activities for %d deals", len(deals)*2, len(deals))
	}

	// Create activities
	for _, act := range activities {
		if err := database.DB.Create(&act).Error; err != nil {
			return err
		}
		log.Printf("Created activity: %s (id: %s, type: %s, account_id: %v)", act.Description, act.ID, act.Type, act.AccountID)
	}

	log.Printf("Activities seeded successfully (%d activities created)", len(activities))
	return nil
}

