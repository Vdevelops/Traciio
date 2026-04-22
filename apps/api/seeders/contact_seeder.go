package seeders

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/account"
	"github.com/gilabs/crm-healthcare/api/internal/domain/category"
	"github.com/gilabs/crm-healthcare/api/internal/domain/contact"
	"github.com/gilabs/crm-healthcare/api/internal/domain/contact_role"
)

// SeedContacts seeds initial contacts
func SeedContacts() error {
	// Check if contacts already exist
	var count int64
	database.DB.Model(&contact.Contact{}).Count(&count)
	if count > 0 {
		log.Println("Contacts already seeded, skipping...")
		return nil
	}

	// Get accounts for relationships (with category preload)
	var accounts []account.Account
	if err := database.DB.Preload("Category").Find(&accounts).Error; err != nil {
		return err
	}

	if len(accounts) == 0 {
		log.Println("Warning: No accounts found, skipping contact seeding")
		return nil
	}

	// Get contact roles by code
	var doctorRole, picRole, managerRole, otherRole contact_role.ContactRole
	if err := database.DB.Where("code = ?", "DOCTOR").First(&doctorRole).Error; err != nil {
		return err
	}
	if err := database.DB.Where("code = ?", "PIC").First(&picRole).Error; err != nil {
		return err
	}
	if err := database.DB.Where("code = ?", "MANAGER").First(&managerRole).Error; err != nil {
		return err
	}
	if err := database.DB.Where("code = ?", "OTHER").First(&otherRole).Error; err != nil {
		return err
	}

	// Get categories by code for mapping
	var hospitalCategory, clinicCategory, pharmacyCategory category.Category
	if err := database.DB.Where("code = ?", "HOSPITAL").First(&hospitalCategory).Error; err != nil {
		return err
	}
	if err := database.DB.Where("code = ?", "CLINIC").First(&clinicCategory).Error; err != nil {
		return err
	}
	if err := database.DB.Where("code = ?", "PHARMACY").First(&pharmacyCategory).Error; err != nil {
		return err
	}

	// Map accounts by category for easier assignment
	hospitalAccounts := []account.Account{}
	clinicAccounts := []account.Account{}
	pharmacyAccounts := []account.Account{}

	for _, acc := range accounts {
		if acc.Category != nil {
			switch acc.Category.Code {
			case "HOSPITAL":
				hospitalAccounts = append(hospitalAccounts, acc)
			case "CLINIC":
				clinicAccounts = append(clinicAccounts, acc)
			case "PHARMACY":
				pharmacyAccounts = append(pharmacyAccounts, acc)
			}
		}
	}

	contacts := []contact.Contact{}

	// Contacts for hospitals
	if len(hospitalAccounts) > 0 {
		// RSUD Jakarta contacts
		contacts = append(contacts, contact.Contact{
			AccountID: hospitalAccounts[0].ID,
			Name:      "Dr. Ahmad Wijaya, Sp.PD",
			RoleID:    doctorRole.ID,
			Phone:     "+6281234568001",
			Email:     "ahmad.wijaya@rsud-jakarta.go.id",
			Position:  "Kepala Bagian Internal Medicine",
			Notes:     "Spesialis penyakit dalam, kontak utama untuk produk kardiovaskular",
		})
		contacts = append(contacts, contact.Contact{
			AccountID: hospitalAccounts[0].ID,
			Name:      "Budi Santoso",
			RoleID:    picRole.ID,
			Phone:     "+6281234568002",
			Email:     "budi.santoso@rsud-jakarta.go.id",
			Position:  "Manager Procurement",
			Notes:     "Menangani pengadaan obat dan alat medis",
		})

		// RSCM contacts
		if len(hospitalAccounts) > 1 {
			contacts = append(contacts, contact.Contact{
				AccountID: hospitalAccounts[1].ID,
				Name:      "Dr. Siti Nurhaliza, Sp.JP",
				RoleID:    doctorRole.ID,
				Phone:     "+6281234568003",
				Email:     "siti.nurhaliza@rscm.go.id",
				Position:  "Kepala Bagian Kardiologi",
				Notes:     "Spesialis jantung dan pembuluh darah",
			})
			contacts = append(contacts, contact.Contact{
				AccountID: hospitalAccounts[1].ID,
				Name:      "Rina Kartika",
				RoleID:    managerRole.ID,
				Phone:     "+6281234568004",
				Email:     "rina.kartika@rscm.go.id",
				Position:  "Direktur Operasional",
				Notes:     "Keputusan strategis untuk pengadaan",
			})
		}

		// RS Pondok Indah contacts
		if len(hospitalAccounts) > 2 {
			contacts = append(contacts, contact.Contact{
				AccountID: hospitalAccounts[2].ID,
				Name:      "Dr. Michael Chen, Sp.OG",
				RoleID:    doctorRole.ID,
				Phone:     "+6281234568005",
				Email:     "michael.chen@rspondokindah.co.id",
				Position:  "Kepala Bagian Obstetri & Ginekologi",
				Notes:     "Spesialis kandungan dan kebidanan",
			})
		}
	}

	// Contacts for clinics
	if len(clinicAccounts) > 0 {
		// Klinik Sehat Sentosa contacts
		contacts = append(contacts, contact.Contact{
			AccountID: clinicAccounts[0].ID,
			Name:      "Dr. Indra Gunawan",
			RoleID:    doctorRole.ID,
			Phone:     "+6281234568006",
			Email:     "indra.gunawan@kliniksehatsentosa.com",
			Position:  "Dokter Umum",
			Notes:     "Dokter praktik umum, menerima pasien walk-in",
		})
		contacts = append(contacts, contact.Contact{
			AccountID: clinicAccounts[0].ID,
			Name:      "Sari Dewi",
			RoleID:    picRole.ID,
			Phone:     "+6281234568007",
			Email:     "sari.dewi@kliniksehatsentosa.com",
			Position:  "Administrator",
			Notes:     "Menangani administrasi dan pengadaan",
		})

		// Klinik Medika Pratama contacts
		if len(clinicAccounts) > 1 {
			contacts = append(contacts, contact.Contact{
				AccountID: clinicAccounts[1].ID,
				Name:      "Dr. Lisa Permata",
				RoleID:    doctorRole.ID,
				Phone:     "+6281234568008",
				Email:     "lisa.permata@medikapratama.com",
				Position:  "Dokter Spesialis Anak",
				Notes:     "Spesialis kesehatan anak",
			})
		}

		// Klinik Bunda Sejahtera contacts
		if len(clinicAccounts) > 2 {
			contacts = append(contacts, contact.Contact{
				AccountID: clinicAccounts[2].ID,
				Name:      "Dr. Maria Sari",
				RoleID:    doctorRole.ID,
				Phone:     "+6281234568009",
				Email:     "maria.sari@bundasejahtera.com",
				Position:  "Dokter Spesialis Kandungan",
				Notes:     "Spesialis kandungan dan keluarga berencana",
			})
		}
	}

	// Contacts for pharmacies
	if len(pharmacyAccounts) > 0 {
		// Apotek Kimia Farma contacts
		contacts = append(contacts, contact.Contact{
			AccountID: pharmacyAccounts[0].ID,
			Name:      "Apt. Dedi Kurniawan, S.Farm",
			RoleID:    managerRole.ID,
			Phone:     "+6281234568010",
			Email:     "dedi.kurniawan@kimiafarma.co.id",
			Position:  "Apoteker Kepala",
			Notes:     "Menangani pengadaan dan manajemen stok obat",
		})
		contacts = append(contacts, contact.Contact{
			AccountID: pharmacyAccounts[0].ID,
			Name:      "Rizki Pratama",
			RoleID:    picRole.ID,
			Phone:     "+6281234568011",
			Email:     "rizki.pratama@kimiafarma.co.id",
			Position:  "Supervisor",
			Notes:     "Koordinator operasional harian",
		})

		// Apotek Guardian contacts
		if len(pharmacyAccounts) > 1 {
			contacts = append(contacts, contact.Contact{
				AccountID: pharmacyAccounts[1].ID,
				Name:      "Apt. Sarah Putri, S.Farm",
				RoleID:    managerRole.ID,
				Phone:     "+6281234568012",
				Email:     "sarah.putri@guardian.co.id",
				Position:  "Branch Manager",
				Notes:     "Manager cabang, keputusan pengadaan",
			})
		}
	}

	// Create contacts
	for _, cont := range contacts {
		if err := database.DB.Create(&cont).Error; err != nil {
			return err
		}
		log.Printf("Created contact: %s (id: %s, account_id: %s, role_id: %s)", cont.Name, cont.ID, cont.AccountID, cont.RoleID)
	}

	log.Printf("Contacts seeded successfully (%d contacts created)", len(contacts))
	return nil
}

