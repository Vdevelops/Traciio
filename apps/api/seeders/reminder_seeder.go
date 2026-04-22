package seeders

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/notification"
	"github.com/gilabs/crm-healthcare/api/internal/domain/reminder"
	"github.com/gilabs/crm-healthcare/api/internal/domain/task"
)

// SeedReminders seeds initial reminders data
func SeedReminders() error {
	// Check if reminders already exist
	var count int64
	database.DB.Model(&reminder.Reminder{}).Count(&count)
	if count > 0 {
		log.Println("Reminders already seeded, skipping...")
		return nil
	}

	// Get tasks
	var tasks []task.Task
	if err := database.DB.Find(&tasks).Error; err != nil {
		return err
	}
	if len(tasks) == 0 {
		log.Println("Warning: No tasks found, skipping reminder seeding")
		return nil
	}

	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)
	nextWeek := now.Add(7 * 24 * time.Hour)
	nextMonth := now.Add(30 * 24 * time.Hour)

	reminders := []reminder.Reminder{}

	// Create reminders for each task
	for i, t := range tasks {
		if i >= 3 {
			break // Only create reminders for first 3 tasks
		}

		var remindAt time.Time
		var message string
		switch i {
		case 0:
			remindAt = tomorrow.Add(-2 * time.Hour) // 2 hours before due date
			message = "Jangan lupa untuk " + t.Title
		case 1:
			remindAt = nextWeek.Add(-1 * 24 * time.Hour) // 1 day before due date
			message = "Persiapan untuk " + t.Title
		case 2:
			remindAt = nextMonth.Add(-3 * 24 * time.Hour) // 3 days before due date
			message = "Reminder: " + t.Title
		}

		reminders = append(reminders, reminder.Reminder{
			TaskID:      t.ID,
			RemindAt:    remindAt,
			ReminderType: "in_app",
			Message:     message,
			IsSent:      false,
			CreatedBy:   t.CreatedBy,
		})
	}

	// Create reminders
	for _, rem := range reminders {
		if err := database.DB.Create(&rem).Error; err != nil {
			return err
		}
		log.Printf("Created reminder for task %s (remind_at: %s)", rem.TaskID, rem.RemindAt.Format("2006-01-02 15:04:05"))
	}

	log.Printf("Seeded %d reminders", len(reminders))
	return nil
}

// SeedNotifications seeds initial notifications data (connected to reminders)
func SeedNotifications() error {
	// Check if notifications already exist
	var count int64
	database.DB.Model(&notification.Notification{}).Count(&count)
	if count > 0 {
		log.Println("Notifications already seeded, skipping...")
		return nil
	}

	// Get reminders that are already sent (or create some test notifications)
	var reminders []reminder.Reminder
	if err := database.DB.Preload("Task").Find(&reminders).Error; err != nil {
		return err
	}

	if len(reminders) == 0 {
		log.Println("Warning: No reminders found, skipping notification seeding")
		return nil
	}

	notifications := []notification.Notification{}

	// Create notifications for some reminders (simulate already sent reminders)
	for i, rem := range reminders {
		if i >= 2 {
			break // Only create notifications for first 2 reminders
		}

		// Mark reminder as sent
		rem.IsSent = true
		sentAt := rem.RemindAt.Add(5 * time.Minute) // Simulate sent 5 minutes after remind_at
		rem.SentAt = &sentAt
		if err := database.DB.Model(&rem).Updates(map[string]interface{}{
			"is_sent": true,
			"sent_at": sentAt,
		}).Error; err != nil {
			log.Printf("Warning: Failed to mark reminder as sent: %v", err)
		}

		// Create notification
		title := "Reminder"
		message := rem.Message
		if message == "" {
			message = "You have a reminder"
		}

		if rem.Task != nil {
			title = "Task Reminder: " + rem.Task.Title
			if message == "" {
				message = "Task reminder: " + rem.Task.Title
			}
		}

		// Prepare notification data
		notificationData := map[string]interface{}{
			"reminder_id": rem.ID,
			"task_id":     rem.TaskID,
			"remind_at":   rem.RemindAt.Format(time.RFC3339),
		}
		if rem.Task != nil {
			notificationData["task_title"] = rem.Task.Title
		}

		dataJSON, _ := json.Marshal(notificationData)

		notifications = append(notifications, notification.Notification{
			UserID:  rem.CreatedBy,
			Title:   title,
			Message: message,
			Type:    "reminder",
			Data:    string(dataJSON),
			IsRead:  i == 0, // First notification is read, second is unread
		})
	}

	// Create notifications
	for _, notif := range notifications {
		if err := database.DB.Create(&notif).Error; err != nil {
			return err
		}
		log.Printf("Created notification for user %s: %s", notif.UserID, notif.Title)
	}

	log.Printf("Seeded %d notifications", len(notifications))
	return nil
}

