package seeders

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/activity_type"
	"gorm.io/gorm/clause"
)

// SeedActivityTypes seeds initial activity types
func SeedActivityTypes() error {
	activityTypes := []activity_type.ActivityType{
		{
			Name:        "Call",
			Code:        "call",
			Description: "Phone call with account or contact",
			Icon:        "phone",
			BadgeColor:  "secondary",
			Status:      "active",
			Order:       1,
		},
		{
			Name:        "Whatsapp / Chat",
			Code:        "whatsapp_chat",
			Description: "Instant messaging conversation with a prospect",
			Icon:        "message-circle",
			BadgeColor:  "secondary",
			Status:      "active",
			Order:       2,
		},
		{
			Name:        "Email",
			Code:        "email",
			Description: "Email communication",
			Icon:        "mail",
			BadgeColor:  "secondary",
			Status:      "active",
			Order:       3,
		},
		{
			Name:        "Note",
			Code:        "note",
			Description: "Internal note about the prospect",
			Icon:        "sticky-note",
			BadgeColor:  "outline",
			Status:      "active",
			Order:       4,
		},
		{
			Name:        "Follow_Up",
			Code:        "follow_up",
			Description: "Follow-up action planned for the prospect",
			Icon:        "rotate-cw",
			BadgeColor:  "default",
			Status:      "active",
			Order:       5,
		},
		{
			Name:        "Presentation / Demo / Meet",
			Code:        "presentation_demo_meet",
			Description: "Presentation, demo, or in-person meeting with the prospect",
			Icon:        "presentation",
			BadgeColor:  "default",
			Status:      "active",
			Order:       6,
		},
		{
			Name:        "Document Proposal Sent",
			Code:        "document_proposal_sent",
			Description: "Proposal or supporting document has been sent to the prospect",
			Icon:        "file-text",
			BadgeColor:  "default",
			Status:      "active",
			Order:       7,
		},
	}

	canonicalCodes := make([]string, 0, len(activityTypes))
	for _, at := range activityTypes {
		canonicalCodes = append(canonicalCodes, at.Code)
		if err := database.DB.Clauses(clause.OnConflict{
			Columns: []clause.Column{{Name: "code"}},
			DoUpdates: clause.AssignmentColumns([]string{
				"name",
				"description",
				"icon",
				"badge_color",
				"status",
				"order",
			}),
		}).Create(&at).Error; err != nil {
			return err
		}
	}

	var obsoleteTypes []activity_type.ActivityType
	if err := database.DB.Where("code NOT IN ?", canonicalCodes).Find(&obsoleteTypes).Error; err != nil {
		return err
	}

	for _, obsoleteType := range obsoleteTypes {
		if err := database.DB.Table("activities").
			Where("activity_type_id = ?", obsoleteType.ID).
			Update("activity_type_id", nil).Error; err != nil {
			log.Printf("Warning: Failed to nullify activity type references for %s: %v", obsoleteType.Code, err)
			continue
		}

		if err := database.DB.Unscoped().Delete(&obsoleteType).Error; err != nil {
			log.Printf("Warning: Failed to hard delete obsolete activity type %s: %v", obsoleteType.Code, err)
		}
	}

	log.Printf("Seeded %d canonical activity types", len(activityTypes))
	return nil
}
