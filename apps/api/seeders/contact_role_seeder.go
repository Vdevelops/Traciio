package seeders

import (
	"log"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	"github.com/gilabs/crm-healthcare/api/internal/domain/contact_role"
)

// SeedContactRoles seeds the contact_roles table
func SeedContactRoles() error {
	log.Println("Seeding contact roles...")

	contactRoles := []contact_role.ContactRole{
		{
			Name:        "Doctor",
			Code:        "DOCTOR",
			Description: "Medical doctor role",
			BadgeColor:  "default",
			Status:      "active",
		},
		{
			Name:        "PIC",
			Code:        "PIC",
			Description: "Person in Charge role",
			BadgeColor:  "secondary",
			Status:      "active",
		},
		{
			Name:        "Manager",
			Code:        "MANAGER",
			Description: "Manager role",
			BadgeColor:  "outline",
			Status:      "active",
		},
		{
			Name:        "Other",
			Code:        "OTHER",
			Description: "Other contact role",
			BadgeColor:  "success",
			Status:      "active",
		},
	}

	for _, cr := range contactRoles {
		var existing contact_role.ContactRole
		if err := database.DB.Where("code = ?", cr.Code).First(&existing).Error; err != nil {
			if err.Error() == "record not found" {
				if err := database.DB.Create(&cr).Error; err != nil {
					return err
				}
				log.Printf("Created contact role: %s", cr.Name)
			} else {
				return err
			}
		} else {
			log.Printf("Contact role %s already exists, skipping", cr.Code)
		}
	}

	log.Println("Contact roles seeded successfully")
	return nil
}

