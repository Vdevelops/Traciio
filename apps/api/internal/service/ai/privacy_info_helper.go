package ai

import (
	ai_settings "github.com/gilabs/crm-healthcare/api/internal/domain/ai_settings"
)

// buildPrivacyInfoMessage builds a formatted privacy information message
func buildPrivacyInfoMessage(dataPrivacy ai_settings.DataPrivacySettings, isDefault bool) string {
	privacyInfo := "**PENGATURAN DATA PRIVACY:**\n\n"
	
	if isDefault {
		privacyInfo += "Pengaturan data privacy belum dikonfigurasi. Secara default, semua jenis data dapat diakses:\n\n"
	} else {
		privacyInfo += "Akses ke data berikut diizinkan:\n\n"
	}

	// Route Optimization
	privacyInfo += "**ROUTE OPTIMIZATION:**\n"
	privacyInfo += formatPrivacyToggle(dataPrivacy.AllowRouteOptimization, "**Route Optimization** (Optimasi Rute Kunjungan)")

	// Sales domain
	privacyInfo += "\n**SALES:**\n"
	privacyInfo += formatPrivacyToggle(dataPrivacy.AllowLeads, "**Leads** (Prospek)")
	privacyInfo += formatPrivacyToggle(dataPrivacy.AllowPipelines, "**Pipeline** (Alur Penjualan)")
	privacyInfo += formatPrivacyToggle(dataPrivacy.AllowDeals, "**Deals** (Kesempatan Penjualan)")
	privacyInfo += formatPrivacyToggle(dataPrivacy.AllowSchedule, "**Schedules** (Jadwal Kunjungan)")
	privacyInfo += formatPrivacyToggle(dataPrivacy.AllowVisitReports, "**Visit Reports** (Laporan Kunjungan)")
	privacyInfo += formatPrivacyToggle(dataPrivacy.AllowTasks, "**Tasks** (Tugas)")
	privacyInfo += formatPrivacyToggle(dataPrivacy.AllowActivities, "**Activities** (Aktivitas)")

	// Inventory domain
	privacyInfo += "\n**INVENTORY:**\n"
	privacyInfo += formatPrivacyToggle(dataPrivacy.AllowProducts, "**Products** (Produk)")

	// Customer domain
	privacyInfo += "\n**CUSTOMERS:**\n"
	privacyInfo += formatPrivacyToggle(dataPrivacy.AllowAccounts, "**Accounts** (Akun/Fasilitas Kesehatan)")
	privacyInfo += formatPrivacyToggle(dataPrivacy.AllowContacts, "**Contacts** (Kontak/Dokter/Apoteker)")

	// Analytics domain
	privacyInfo += "\n**ANALYTICS:**\n"
	privacyInfo += formatPrivacyToggle(dataPrivacy.AllowSalesPerformance, "**Sales Performance** (Analisis Performa Penjualan)")
	privacyInfo += formatPrivacyToggle(dataPrivacy.AllowProductAnalysis, "**Product Analytics** (Analisis Produk)")
	privacyInfo += formatPrivacyToggle(dataPrivacy.AllowReports, "**Reports** (Laporan)")

	// Management domain
	privacyInfo += "\n**MANAGEMENT:**\n"
	privacyInfo += formatPrivacyToggle(dataPrivacy.AllowUsers, "**Users** (Pengguna)")
	privacyInfo += formatPrivacyToggle(dataPrivacy.AllowRoles, "**Roles** (Peran & Izin)")
	privacyInfo += formatPrivacyToggle(dataPrivacy.AllowGroups, "**Groups** (Segmentasi Grup)")
	privacyInfo += formatPrivacyToggle(dataPrivacy.AllowBrickManagement, "**Bricks** (Manajemen Wilayah/Teritori)")
	privacyInfo += formatPrivacyToggle(dataPrivacy.AllowTarget, "**Targets** (Manajemen Target/Quota)")
	
	return privacyInfo
}

// formatPrivacyToggle formats a single privacy toggle line
func formatPrivacyToggle(allowed bool, label string) string {
	if allowed {
		return "✓ " + label + "\n"
	}
	return "✗ " + label + " - DISABLED\n"
}
