package seeders

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/database"
)

// BackfillBrickLinks populates brick_id for records that can infer it from related users.
//
// Idempotent:
// - Only fills rows where brick_id is currently NULL
// - Never overwrites an existing brick_id
//
// This is primarily intended to keep analytics (brick performance/targets) from appearing empty
// in dev environments where seeders created data before brick assignments were present.
func BackfillBrickLinks() error {
	// Accounts: derive from assigned_to user's brick_id
	if err := database.DB.Exec(`
		UPDATE accounts a
		SET brick_id = u.brick_id
		FROM users u
		WHERE a.brick_id IS NULL
		  AND a.assigned_to = u.id
		  AND u.brick_id IS NOT NULL
		  AND a.deleted_at IS NULL
	`).Error; err != nil {
		return err
	}

	// Deals: derive from assigned_to user's brick_id
	if err := database.DB.Exec(`
		UPDATE deals d
		SET brick_id = u.brick_id
		FROM users u
		WHERE d.brick_id IS NULL
		  AND d.assigned_to = u.id
		  AND u.brick_id IS NOT NULL
		  AND d.deleted_at IS NULL
	`).Error; err != nil {
		return err
	}

	// Visit reports: derive from sales_rep_id user's brick_id
	if err := database.DB.Exec(`
		UPDATE visit_reports vr
		SET brick_id = u.brick_id
		FROM users u
		WHERE vr.brick_id IS NULL
		  AND vr.sales_rep_id = u.id
		  AND u.brick_id IS NOT NULL
		  AND vr.deleted_at IS NULL
	`).Error; err != nil {
		return err
	}

	log.Println("✅ Backfilled brick_id links for accounts/deals/visit_reports")
	return nil
}
