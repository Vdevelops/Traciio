package seeders

import (
	"log"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/schedule"
	"github.com/gilabs/crm-healthcare/api/internal/domain/task"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
)

// SeedSchedules seeds initial schedules data
// Note: Schedules are now linked to tasks, so tasks must be seeded first
func SeedSchedules() error {
	// Check if schedules already exist
	var count int64
	database.DB.Model(&schedule.Schedule{}).Count(&count)
	if count > 0 {
		log.Println("Schedules already seeded, skipping...")
		return nil
	}

	// Get tasks (schedules are linked to tasks)
	var tasks []task.Task
	if err := database.DB.Where("assigned_to IS NOT NULL").Find(&tasks).Error; err != nil {
		return err
	}
	if len(tasks) == 0 {
		log.Println("Warning: No tasks with assigned users found, skipping schedule seeding")
		return nil
	}

	// Get users for CreatedBy
	var users []user.User
	if err := database.DB.Find(&users).Error; err != nil {
		return err
	}
	if len(users) == 0 {
		log.Println("Warning: No users found, skipping schedule seeding")
		return nil
	}

	// Get admin user for CreatedBy
	var adminUser user.User
	if err := database.DB.Where("email = ?", "admin@example.com").First(&adminUser).Error; err != nil {
		log.Printf("Warning: Admin user not found, using first user: %v", err)
		adminUser = users[0]
	}

	// Time helpers
	now := time.Now()
	today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	tomorrow := today.Add(24 * time.Hour)
	nextWeek := today.Add(7 * 24 * time.Hour)
	yesterday := today.Add(-24 * time.Hour)

	// Helper function for string pointer
	strPtr := func(s string) *string {
		if s == "" {
			return nil
		}
		return &s
	}

	// Helper function for int pointer
	intPtr := func(i int) *int {
		return &i
	}

	schedules := []schedule.Schedule{}

	// Create schedules from tasks
	// Only use tasks that have assigned users and due dates
	for i, t := range tasks {
		if t.AssignedTo == nil {
			continue // Skip tasks without assigned user
		}

		// Calculate scheduled_at based on task due_date or use default
		var scheduledAt time.Time
		if t.DueDate != nil {
			// Schedule 1 hour before task due date as reminder
			scheduledAt = t.DueDate.Add(-1 * time.Hour)
		} else {
			// Use default times if no due date
			switch i % 4 {
			case 0:
				scheduledAt = tomorrow.Add(9 * time.Hour) // Tomorrow 09:00
			case 1:
				scheduledAt = today.Add(14 * time.Hour) // Today 14:00
			case 2:
				scheduledAt = yesterday.Add(10 * time.Hour) // Yesterday 10:00
			case 3:
				scheduledAt = nextWeek.Add(10 * time.Hour) // Next week 10:00
			}
		}

		// Determine status based on task status
		status := "pending"
		if t.Status == "completed" {
			status = "completed"
		} else if t.Status == "cancelled" {
			status = "cancelled"
		} else if scheduledAt.Before(now) {
			status = "confirmed" // Past schedules are confirmed
		}

		// Create schedule title from task
		title := "Reminder: " + t.Title
		if len(title) > 255 {
			title = title[:252] + "..."
		}

		// Create description from task description
		description := strPtr("")
		if t.Description != "" {
			desc := "Task: " + t.Description
			description = &desc
		}

		// Set reminder minutes before (default: 15 minutes)
		reminderMinutesBefore := intPtr(15)

		schedules = append(schedules, schedule.Schedule{
			TaskID:                   &t.ID,         // TaskID is now nullable (*string)
			UserID:                   *t.AssignedTo, // User from task.assigned_to
			Title:                    title,
			Description:              description,
			ScheduledAt:              scheduledAt,
			Status:                   status,
			ReminderMinutesBefore:    reminderMinutesBefore,
			GoogleCalendarSyncStatus: "not_synced",
			CreatedBy:                adminUser.ID,
		})

		// Limit to first 5 schedules to avoid too many
		if len(schedules) >= 5 {
			break
		}
	}

	// Insert schedules
	for _, s := range schedules {
		if err := database.DB.Create(&s).Error; err != nil {
			log.Printf("Error seeding schedule '%s': %v", s.Title, err)
			return err
		}
	}

	log.Printf("Successfully seeded %d schedules", len(schedules))
	return nil
}
