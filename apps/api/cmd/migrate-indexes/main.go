package main

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/config"
	"github.com/gilabs/crm-healthcare/api/internal/database"
	"gorm.io/gorm"
)

func main() {
	// Initialize logger
	log.Println("🚀 Starting Performance Index Migration (A+ Upgrade)...")

	// Load configuration
	if err := config.Load(); err != nil {
		log.Fatal("Failed to load config:", err)
	}

	// Connect to database
	if err := database.Connect(); err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer database.Close()

	// DROP broken FTS index explicitly - REMOVED for Optimal fix
	// We want these indexes to exist!
	
	// Add indexes
	addIndexes(database.DB)

	log.Println("✅ A+ Performance Indexes Applied Successfully!")
}

func addIndexes(db *gorm.DB) {
	indexes := []struct {
		Table   string
		Column  string
		Name    string
		Unique  bool
		Message string
	}{
		// Users Table
		{"users", "email", "idx_users_email_perf", false, "Optimizing user lookup by email"},
		{"users", "role_id", "idx_users_role_id_perf", false, "Optimizing user filtering by role"},
		{"users", "group_id", "idx_users_group_id_perf", false, "Optimizing user filtering by group"},
		{"users", "brick_id", "idx_users_brick_id_perf", false, "Optimizing user filtering by brick"},

		// Leads Table
		{"leads", "status_id", "idx_leads_status_id_perf", false, "Optimizing lead filtering by status"},
		{"leads", "assigned_to", "idx_leads_assigned_to_perf", false, "Optimizing lead filtering by owner"},
		{"leads", "created_at", "idx_leads_created_at_perf", false, "Optimizing lead sorting/filtering by date"},

		// Activities Table
		{"activities", "user_id", "idx_activities_user_id_perf", false, "Optimizing activity log filtering by user"},
		{"activities", "created_at", "idx_activities_created_at_perf", false, "Optimizing activity log sorting"},
		{"activities", "deal_id", "idx_activities_deal_id_perf", false, "Optimizing activity log filtering by deal"},

		// Deals Table
		{"deals", "stage_id", "idx_deals_stage_id_perf", false, "Optimizing deal pipeline view"}, // Fixed column name if it was pipeline_stage_id? Check entity. It is StageID -> stage_id.
		{"deals", "assigned_to", "idx_deals_assigned_to_perf", false, "Optimizing deal filtering by owner"}, // Entity uses AssignedTo -> assigned_to (user_id in logs was likely wrong guess or alias)

		// Pipeline Stages Table
		// {"pipeline_stages", "pipeline_id", ...} // PipelineStage doesn't have pipeline_id in Entity? Check Entity. It only has ID, Name, Code... No pipeline_id. Removing.

		// Visit Reports
		{"visit_reports", "sales_rep_id", "idx_visit_reports_sales_rep_id_perf", false, "Optimizing visit report filtering by user"}, // Fixed: sales_rep_id
		{"visit_reports", "visit_date", "idx_visit_reports_visit_date_perf", false, "Optimizing visit report filtering by date"}, // Fixed: visit_date

		// Accounts - Spatial index for BBOX map queries
		{"accounts", "latitude, longitude", "idx_accounts_lat_lng_spatial", false, "Optimizing account map viewport (BBOX) queries"},

		// Deal Histories - Optimizing pipeline movement score aggregation
		{"deal_histories", "changed_at", "idx_deal_histories_changed_at_perf", false, "Optimizing deal histories by changed_at timestamp"},
	}

	for _, idx := range indexes {
		log.Printf("Processing: %s...", idx.Message)
		
		// Check if index exists first to avoid errors
		if db.Migrator().HasIndex(idx.Table, idx.Name) {
			log.Printf("  -> Index %s already exists. Skipping.", idx.Name)
			continue
		}

		// Create index
		query := "CREATE INDEX CONCURRENTLY IF NOT EXISTS " + idx.Name + " ON " + idx.Table + " (" + idx.Column + ")"
		if idx.Unique {
			query = "CREATE UNIQUE INDEX CONCURRENTLY IF NOT EXISTS " + idx.Name + " ON " + idx.Table + " (" + idx.Column + ")"
		}

		// Note: GORM's Exec doesn't support CONCURRENTLY inside a transaction, so we use standard CREATE INDEX if basic Exec fails? 
		// Actually best to just use non-concurrent for safety in this script unless we are sure.
		// Retrying simple create if concurrent fails.
		if err := db.Exec(query).Error; err != nil {
			log.Printf("  -> Failed concurrent creation, trying standard: %v", err)
			query = "CREATE INDEX IF NOT EXISTS " + idx.Name + " ON " + idx.Table + " (" + idx.Column + ")"
			if err := db.Exec(query).Error; err != nil {
				log.Printf("  ❌ Failed to create index %s: %v", idx.Name, err)
			} else {
				log.Printf("  ✅ Created index %s", idx.Name)
			}
		} else {
			log.Printf("  ✅ Created index %s", idx.Name)
		}
	}

	// SPECIAL: FTS Indexes (GIN)
	ftsIndexes := []struct {
		Table      string
		Name       string
		Expression string
	}{
		{
			"deals", 
			"idx_deals_fts", 
			"USING GIN(to_tsvector('english', title || ' ' || COALESCE(description, '') || ' ' || COALESCE(notes, '')))",
		},
		{
			"tasks",
			"idx_tasks_fts",
			"USING GIN(to_tsvector('english', title || ' ' || COALESCE(description, '')))",
		},
	}

	for _, idx := range ftsIndexes {
		log.Printf("Processing FTS Index: %s...", idx.Name)
		if db.Migrator().HasIndex(idx.Table, idx.Name) {
			log.Printf("  -> FTS Index %s already exists. Skipping.", idx.Name)
			continue
		}
		
		query := "CREATE INDEX IF NOT EXISTS " + idx.Name + " ON " + idx.Table + " " + idx.Expression
		if err := db.Exec(query).Error; err != nil {
			log.Printf("  ❌ Failed to create FTS index %s: %v", idx.Name, err)
		} else {
			log.Printf("  ✅ Created FTS index %s", idx.Name)
		}
	}
}
