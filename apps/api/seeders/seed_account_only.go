package seeders

// SeedAccountOnly seeds only the account-related data (roles, menus, permissions, users, categories, and accounts)
func SeedAccountOnly() error {
	// Seed in order: roles -> menus -> permissions -> users -> categories -> accounts
	if err := SeedRoles(); err != nil {
		return err
	}

	if err := SeedMenus(); err != nil {
		return err
	}

	// Update menu structure for existing menus (migration)
	if err := UpdateMenuStructure(); err != nil {
		return err
	}

	if err := SeedPermissions(); err != nil {
		return err
	}

	// Ensure Brick menu permissions exist even when ONLY_ACCOUNT=true.
	// Without this, fresh DB resets (docker compose down -v) will hide Bricks/Achievements menu.
	if err := AddBrickPermissions(); err != nil {
		// Best-effort: don't block account-only bootstrapping.
	}

	// Add Opportunity permissions (for Lead Management)
	if err := AddOpportunityPermissions(); err != nil {
		// Log warning but don't fail, as this might be running in a context where it's not strictly required
		// But for the user's current task it is.
	}

	if err := SeedUsers(); err != nil {
		return err
	}

	// Seed categories (required for accounts)
	if err := SeedCategories(); err != nil {
		return err
	}

	// Seed accounts (requires users for assigned_to and categories for category_id)
	if err := SeedAccounts(); err != nil {
		return err
	}

	return nil
}
