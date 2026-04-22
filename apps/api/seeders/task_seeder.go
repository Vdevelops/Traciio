package seeders

import (
	"log"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/account"
	"github.com/gilabs/crm-healthcare/api/internal/domain/contact"
	"github.com/gilabs/crm-healthcare/api/internal/domain/schedule"
	"github.com/gilabs/crm-healthcare/api/internal/domain/task"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
)

// SeedTasks seeds initial tasks data
func SeedTasks() error {
	// Check if tasks already exist
	var count int64
	database.DB.Model(&task.Task{}).Count(&count)
	if count > 0 {
		log.Println("Tasks already seeded, skipping...")
		return nil
	}

	// Get users for AssignedTo and CreatedBy
	var users []user.User
	if err := database.DB.Find(&users).Error; err != nil {
		return err
	}
	if len(users) == 0 {
		log.Println("Warning: No users found, skipping task seeding")
		return nil
	}
	defaultUser := users[0]

	// Get admin user for AssignedFrom
	var adminUser user.User
	if err := database.DB.Where("email = ?", "admin@example.com").First(&adminUser).Error; err != nil {
		log.Printf("Warning: Admin user not found, using first user for assigned_from: %v", err)
		adminUser = defaultUser
	}

	// Get sales user for AssignedTo
	var salesUser user.User
	if err := database.DB.Where("email = ?", "sales@example.com").First(&salesUser).Error; err != nil {
		log.Printf("Warning: Sales user not found, using default user for assigned_to: %v", err)
		salesUser = defaultUser
	}

	// Get accounts
	var accounts []account.Account
	if err := database.DB.Find(&accounts).Error; err != nil {
		return err
	}

	// Get contacts
	var contacts []contact.Contact
	if err := database.DB.Find(&contacts).Error; err != nil {
		return err
	}

	now := time.Now()
	tomorrow := now.Add(24 * time.Hour)
	dayAfterTomorrow := now.Add(48 * time.Hour)
	nextWeek := now.Add(7 * 24 * time.Hour)
	nextWeekPlus2 := now.Add(9 * 24 * time.Hour)
	nextWeekPlus3 := now.Add(10 * 24 * time.Hour)

	// Task templates with variety
	taskTemplates := []struct {
		title       string
		description string
		taskType    string
		priority    string
		status      string
	}{
		{
			title:       "Follow up kunjungan RS untuk penawaran produk baru",
			description: "Hubungi kembali PIC di rumah sakit untuk menindaklanjuti presentasi produk cardiovascular.",
			taskType:    "call",
			priority:    "high",
			status:      "pending",
		},
		{
			title:       "Siapkan proposal kerja sama tahunan",
			description: "Susun proposal lengkap untuk kontrak suplai tahunan termasuk skema diskon.",
			taskType:    "meeting",
			priority:    "urgent",
			status:      "in_progress",
		},
		{
			title:       "Kirim email materi produk ke apotek",
			description: "Kirim katalog produk dan brosur promo ke apotek target.",
			taskType:    "email",
			priority:    "medium",
			status:      "pending",
		},
		{
			title:       "Kunjungi RS untuk presentasi produk baru",
			description: "Lakukan kunjungan ke rumah sakit untuk presentasi produk cardiovascular terbaru. Pastikan membawa sample dan brosur.",
			taskType:    "meeting",
			priority:    "high",
			status:      "pending",
		},
		{
			title:       "Follow up dengan klien untuk closing deal",
			description: "Hubungi klien untuk follow up proposal yang sudah dikirim minggu lalu. Target closing deal bulan ini.",
			taskType:    "call",
			priority:    "urgent",
			status:      "pending",
		},
		{
			title:       "Submit visit report kunjungan kemarin",
			description: "Lengkapi dan submit visit report untuk kunjungan ke apotek kemarin. Pastikan semua foto dan dokumentasi sudah lengkap.",
			taskType:    "general",
			priority:    "medium",
			status:      "in_progress",
		},
		{
			title:       "Persiapan meeting dengan direktur rumah sakit",
			description: "Siapkan presentasi lengkap untuk meeting dengan direktur rumah sakit. Fokus pada produk premium dan value proposition.",
			taskType:    "meeting",
			priority:    "high",
			status:      "pending",
		},
		{
			title:       "Koordinasi dengan tim untuk event produk",
			description: "Koordinasi dengan tim marketing dan technical support untuk persiapan event launching produk baru.",
			taskType:    "meeting",
			priority:    "medium",
			status:      "pending",
		},
		{
			title:       "Review performa sales bulan ini",
			description: "Analisis performa sales bulan ini, identifikasi gap dan buat action plan untuk improvement.",
			taskType:    "general",
			priority:    "low",
			status:      "pending",
		},
		{
			title:       "Update database klien terbaru",
			description: "Update informasi kontak dan kebutuhan produk klien di database CRM.",
			taskType:    "general",
			priority:    "medium",
			status:      "in_progress",
		},
		{
			title:       "Follow up email untuk follow-up meeting",
			description: "Kirim email follow-up setelah meeting dengan klien untuk menindaklanjuti poin-poin yang dibahas.",
			taskType:    "email",
			priority:    "high",
			status:      "pending",
		},
		{
			title:       "Persiapan demo produk untuk klien baru",
			description: "Siapkan demo produk dan presentasi untuk klien potensial. Pastikan semua equipment dan sample ready.",
			taskType:    "meeting",
			priority:    "urgent",
			status:      "pending",
		},
	}

	// Due dates for future tasks (varied)
	futureDueDates := []*time.Time{
		&tomorrow,
		&dayAfterTomorrow,
		&nextWeek,
		&nextWeekPlus2,
		&nextWeekPlus3,
	}

	tasks := []task.Task{}

	// Create variatif tasks for each user (minimal 5 tasks total)
	for i, u := range users {
		if i > 0 {
			break // Only create tasks for the first user to limit exactly 5 tasks
		}
		userID := u.ID
		accountIndex := i % len(accounts)
		if accountIndex >= len(accounts) {
			continue
		}

		// Get accounts and contacts for this user
		var userAccounts []account.Account
		var userContacts []contact.Contact
		if len(accounts) > 0 {
			// Rotate accounts for variety
			for j := 0; j < 3 && (accountIndex+j) < len(accounts); j++ {
				userAccounts = append(userAccounts, accounts[(accountIndex+j)%len(accounts)])
			}
		}
		if len(contacts) > 0 {
			for j := 0; j < 3 && (accountIndex+j) < len(contacts); j++ {
				userContacts = append(userContacts, contacts[(accountIndex+j)%len(contacts)])
			}
		}

		// Create 5 tasks per user (minimal 5 as requested)
		numTasks := 5
		for j := 0; j < numTasks; j++ {
			templateIndex := (i*numTasks + j) % len(taskTemplates)
			template := taskTemplates[templateIndex]

			// Select account and contact
			var accountID *string
			var contactID *string
			if len(userAccounts) > 0 {
				accID := userAccounts[j%len(userAccounts)].ID
				accountID = &accID
			}
			if len(userContacts) > 0 {
				contID := userContacts[j%len(userContacts)].ID
				contactID = &contID
			}

			// Skip if no account (required)
			if accountID == nil {
				continue
			}

			// Select due date (varied future dates)
			dueDateIndex := (i*numTasks + j) % len(futureDueDates)
			dueDate := futureDueDates[dueDateIndex]

			// Determine assigned from (admin or self)
			var assignedFrom *string
			if i%2 == 0 && adminUser.ID != "" {
				assignedFrom = &adminUser.ID
			}

			tasks = append(tasks, task.Task{
				Title:        template.title,
				Description:  template.description,
				Type:         template.taskType,
				Status:       template.status,
				Priority:     template.priority,
				DueDate:      dueDate,
				AssignedTo:   &userID,
				AssignedFrom: assignedFrom,
				AccountID:    accountID,
				ContactID:    contactID,
				CreatedBy:    userID,
			})
		}
	}

	// Removed additional completed tasks loop to keep exact count to 5

	for _, t := range tasks {
		// Skip tasks that don't have minimum required foreign keys
		if t.AccountID == nil {
			continue
		}

		if err := database.DB.Create(&t).Error; err != nil {
			return err
		}
		log.Printf("Created task: %s (id: %s, status: %s, priority: %s)", t.Title, t.ID, t.Status, t.Priority)

		// Auto-create schedule for task if it has due_date and assigned_to
		if t.DueDate != nil && t.AssignedTo != nil {
			// Calculate scheduled_at (1 hour before due date)
			scheduledAt := t.DueDate.Add(-1 * time.Hour)
			reminderMinutes := 60 // Default: 1 hour before

			// Determine status based on task status
			status := "pending"
			if t.Status == "completed" {
				status = "completed"
			} else if t.Status == "cancelled" {
				status = "cancelled"
			} else if scheduledAt.Before(time.Now()) {
				status = "confirmed" // Past schedules are confirmed
			}

			// Create schedule title from task
			title := "Reminder: " + t.Title
			if len(title) > 255 {
				title = title[:252] + "..."
			}

			// Create description from task description
			var descriptionPtr *string
			if t.Description != "" {
				desc := "Task: " + t.Description
				descriptionPtr = &desc
			}

			schedule := schedule.Schedule{
				TaskID:                   &t.ID, // TaskID is now nullable (*string)
				UserID:                   *t.AssignedTo,
				Title:                    title,
				Description:              descriptionPtr,
				ScheduledAt:              scheduledAt,
				Status:                   status,
				ReminderMinutesBefore:    &reminderMinutes,
				GoogleCalendarSyncStatus: "not_synced", // Default: not synced (toggle false)
				CreatedBy:                t.CreatedBy,
			}

			if err := database.DB.Create(&schedule).Error; err != nil {
				log.Printf("Warning: Failed to create schedule for task %s: %v", t.ID, err)
				// Don't fail task seeding if schedule creation fails
			} else {
				log.Printf("Created schedule for task: %s (schedule id: %s)", t.Title, schedule.ID)
			}
		}
	}

	log.Println("Tasks seeded successfully")
	return nil
}
