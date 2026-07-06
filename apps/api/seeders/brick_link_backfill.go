package seeders

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/database"
)

// BackfillBrickLinks populates brick_id for records that can infer it from related users.
//
// Idempotent:
// - Aligns records to the current account territory when account linkage exists
// - Falls back to assigned user brick only when no account-derived brick is available
//
// This is primarily intended to keep analytics (brick performance/targets) from appearing empty
// in dev environments where seeders created data before brick assignments were present.
func BackfillBrickLinks() error {
	// Accounts: derive from assigned_to user's brick_id only when still empty.
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

	// Deals: sync from linked account territory first.
	if err := database.DB.Exec(`
		UPDATE deals d
		SET brick_id = a.brick_id
		FROM accounts a
		WHERE d.account_id = a.id
		  AND a.brick_id IS NOT NULL
		  AND d.deleted_at IS NULL
		  AND a.deleted_at IS NULL
		  AND (d.brick_id IS NULL OR d.brick_id <> a.brick_id)
	`).Error; err != nil {
		return err
	}

	// Deals: fallback to assigned_to user's brick_id only if still empty.
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

	// Visit reports: sync from linked account territory first.
	if err := database.DB.Exec(`
		UPDATE visit_reports vr
		SET brick_id = a.brick_id
		FROM accounts a
		WHERE vr.account_id = a.id
		  AND a.brick_id IS NOT NULL
		  AND vr.deleted_at IS NULL
		  AND a.deleted_at IS NULL
		  AND (vr.brick_id IS NULL OR vr.brick_id <> a.brick_id)
	`).Error; err != nil {
		return err
	}

	// Visit reports: fallback to sales_rep_id user's brick_id only if still empty.
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
