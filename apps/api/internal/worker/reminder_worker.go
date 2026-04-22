package worker

import (
	"encoding/json"
	"log"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/domain/notification"
	"github.com/gilabs/crm-healthcare/api/internal/domain/reminder"
	"github.com/gilabs/crm-healthcare/api/internal/hub"
	"github.com/gilabs/crm-healthcare/api/internal/repository/interfaces"
	notificationservice "github.com/gilabs/crm-healthcare/api/internal/service/notification"
)

type ReminderWorker struct {
	reminderRepo        interfaces.ReminderRepository
	notificationService *notificationservice.Service
	notificationHub     *hub.NotificationHub
	ticker              *time.Ticker
	stopChan            chan bool
}

func NewReminderWorker(
	reminderRepo interfaces.ReminderRepository,
	notificationService *notificationservice.Service,
	notificationHub *hub.NotificationHub,
	interval time.Duration,
) *ReminderWorker {
	return &ReminderWorker{
		reminderRepo:        reminderRepo,
		notificationService: notificationService,
		notificationHub:     notificationHub,
		ticker:              time.NewTicker(interval),
		stopChan:            make(chan bool),
	}
}

// Start starts the reminder worker with panic recovery
func (w *ReminderWorker) Start() {
	log.Println("Reminder worker started")

	go func() {
		// CRITICAL: Panic recovery to prevent worker crash
		defer func() {
			if r := recover(); r != nil {
				log.Printf("PANIC in reminder worker: %v. Restarting worker...", r)
				// Restart the worker after panic
				go w.Start()
			}
		}()

		for {
			select {
			case <-w.ticker.C:
				// Wrap processReminders with panic recovery for each tick
				func() {
					defer func() {
						if r := recover(); r != nil {
							log.Printf("PANIC in processReminders: %v", r)
						}
					}()
					w.processReminders()
				}()
			case <-w.stopChan:
				w.ticker.Stop()
				log.Println("Reminder worker stopped")
				return
			}
		}
	}()
}

// Stop stops the reminder worker
func (w *ReminderWorker) Stop() {
	w.stopChan <- true
}

// processReminders processes pending reminders
func (w *ReminderWorker) processReminders() {
	now := time.Now()

	// Find pending reminders
	reminders, err := w.reminderRepo.FindPendingReminders(now)
	if err != nil {
		log.Printf("Error finding pending reminders: %v", err)
		return
	}

	if len(reminders) == 0 {
		return
	}

	log.Printf("Processing %d pending reminders", len(reminders))

	for _, rem := range reminders {
		if err := w.processReminder(&rem); err != nil {
			log.Printf("Error processing reminder %s: %v", rem.ID, err)
			// Continue with next reminder
			continue
		}
	}
}

// processReminder processes a single reminder
func (w *ReminderWorker) processReminder(rem *reminder.Reminder) error {
	// Only process in_app reminders
	if rem.ReminderType != "in_app" {
		// Mark as sent even if we skip it
		return w.reminderRepo.MarkAsSent(rem.ID, time.Now())
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

	// Create notification
	createReq := &notification.CreateNotificationRequest{
		UserID:  rem.CreatedBy,
		Title:   title,
		Message: message,
		Type:    "reminder",
		Data:    string(dataJSON),
	}

	notif, err := w.notificationService.CreateNotification(createReq)
	if err != nil {
		return err
	}

	// Broadcast notification via WebSocket
	notifMap := map[string]interface{}{
		"id":         notif.ID,
		"user_id":    notif.UserID,
		"title":      notif.Title,
		"message":    notif.Message,
		"type":       notif.Type,
		"is_read":    notif.IsRead,
		"created_at": notif.CreatedAt.Format(time.RFC3339),
	}
	w.notificationHub.BroadcastNotification(rem.CreatedBy, notifMap)

	// Mark reminder as sent
	return w.reminderRepo.MarkAsSent(rem.ID, time.Now())
}
