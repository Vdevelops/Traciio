package seeders

import (
	"errors"
	"fmt"
	"hash/fnv"
	"math"
	"strings"

	"github.com/gilabs/crm-healthcare/api/internal/config"
	"github.com/gilabs/crm-healthcare/api/internal/domain/account"
	"github.com/gilabs/crm-healthcare/api/internal/domain/contact"
	contactrole "github.com/gilabs/crm-healthcare/api/internal/domain/contact_role"
	leaddomain "github.com/gilabs/crm-healthcare/api/internal/domain/lead"
	"github.com/gilabs/crm-healthcare/api/pkg/geocoding"
	"gorm.io/gorm"
)

func syncConvertedLeadCustomer(db *gorm.DB, sourceLead *leaddomain.Lead, assignedTo *string) (*account.Account, *contact.Contact, error) {
	if sourceLead == nil {
		return nil, nil, fmt.Errorf("source lead is required for converted customer sync")
	}

	accountEntity, err := upsertConvertedLeadAccount(db, sourceLead, assignedTo)
	if err != nil {
		return nil, nil, err
	}

	contactEntity, err := upsertConvertedLeadContact(db, sourceLead, accountEntity.ID)
	if err != nil {
		return nil, nil, err
	}

	return accountEntity, contactEntity, nil
}

func upsertConvertedLeadAccount(db *gorm.DB, sourceLead *leaddomain.Lead, assignedTo *string) (*account.Account, error) {
	companyName := strings.TrimSpace(sourceLead.CompanyName)
	if companyName == "" {
		companyName = strings.TrimSpace(sourceLead.FirstName + " " + sourceLead.LastName)
	}
	if companyName == "" {
		return nil, fmt.Errorf("converted lead has no company name to sync")
	}

	categoryID, err := resolveConvertedLeadAccountCategoryID(db, sourceLead)
	if err != nil {
		return nil, err
	}

	lookup := strings.TrimSpace(sourceLead.Email)
	var entity account.Account
	query := db.Where("email = ?", lookup)
	if lookup == "" {
		lookup = strings.TrimSpace(sourceLead.Phone)
		query = db.Where("phone = ?", lookup)
	}
	if lookup != "" {
		if err := query.First(&entity).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	entity.Name = companyName
	entity.CategoryID = categoryID
	entity.Address = sourceLead.Address
	entity.City = sourceLead.City
	entity.Province = sourceLead.Province
	entity.Phone = sourceLead.Phone
	entity.Email = sourceLead.Email
	entity.PostalCode = sourceLead.PostalCode
	entity.Country = "Indonesia"
	entity.Website = sourceLead.Website
	entity.Industry = sourceLead.Industry
	entity.Status = "active"
	entity.AssignedTo = assignedTo
	if entity.Latitude == nil || entity.Longitude == nil {
		entity.Latitude, entity.Longitude = resolveConvertedLeadCoordinates(sourceLead)
	}

	if entity.ID == "" {
		if err := db.Create(&entity).Error; err != nil {
			return nil, err
		}
	} else if err := db.Save(&entity).Error; err != nil {
		return nil, err
	}

	return &entity, nil
}

func resolveConvertedLeadCoordinates(sourceLead *leaddomain.Lead) (*float64, *float64) {
	if sourceLead == nil {
		return nil, nil
	}

	if config.AppConfig != nil && config.AppConfig.Geocoding.Enabled {
		service := geocoding.NewGeocodingService(config.AppConfig.Geocoding.Provider, config.AppConfig.Geocoding.APIKey)
		if result, err := service.GeocodeAddressWithFallback(sourceLead.Address, sourceLead.City, sourceLead.Province); err == nil && result != nil {
			return &result.Latitude, &result.Longitude
		}
	}

	return fallbackConvertedLeadCoordinates(sourceLead)
}

func fallbackConvertedLeadCoordinates(sourceLead *leaddomain.Lead) (*float64, *float64) {
	baseLat, baseLng := -6.9915, 110.4180
	locationKey := strings.ToLower(strings.TrimSpace(sourceLead.City + "|" + sourceLead.Province + "|" + sourceLead.CompanyName + "|" + sourceLead.Email))
	if locationKey == "" {
		locationKey = strings.ToLower(strings.TrimSpace(sourceLead.FirstName + "|" + sourceLead.LastName + "|" + sourceLead.Phone))
	}

	hasher := fnv.New32a()
	_, _ = hasher.Write([]byte(locationKey))
	seed := hasher.Sum32()

	latOffset := (float64(seed%4000)/100000.0 - 0.02)
	lngOffset := (float64((seed/4000)%4000)/100000.0 - 0.02)
	lat := baseLat + latOffset
	lng := baseLng + lngOffset

	if math.IsNaN(lat) || math.IsNaN(lng) {
		return nil, nil
	}

	return &lat, &lng
}

func upsertConvertedLeadContact(db *gorm.DB, sourceLead *leaddomain.Lead, accountID string) (*contact.Contact, error) {
	roleID, err := resolveConvertedLeadContactRoleID(db, sourceLead)
	if err != nil {
		return nil, err
	}

	contactName := strings.TrimSpace(sourceLead.AuthorityPerson)
	if contactName == "" {
		contactName = strings.TrimSpace(sourceLead.FirstName + " " + sourceLead.LastName)
	}
	if contactName == "" {
		contactName = strings.TrimSpace(sourceLead.CompanyName)
	}
	if contactName == "" {
		return nil, fmt.Errorf("converted lead has no contact name to sync")
	}

	lookup := strings.TrimSpace(sourceLead.Email)
	var entity contact.Contact
	query := db.Where("email = ? AND account_id = ?", lookup, accountID)
	if lookup == "" {
		lookup = strings.TrimSpace(sourceLead.Phone)
		query = db.Where("phone = ? AND account_id = ?", lookup, accountID)
	}
	if lookup != "" {
		if err := query.First(&entity).Error; err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
	}

	entity.AccountID = accountID
	entity.RoleID = roleID
	entity.Name = contactName
	entity.Phone = sourceLead.Phone
	entity.Email = sourceLead.Email
	entity.Position = sourceLead.JobTitle
	entity.Notes = strings.TrimSpace(sourceLead.Notes)

	if entity.ID == "" {
		if err := db.Create(&entity).Error; err != nil {
			return nil, err
		}
	} else if err := db.Save(&entity).Error; err != nil {
		return nil, err
	}

	return &entity, nil
}

func resolveConvertedLeadAccountCategoryID(db *gorm.DB, sourceLead *leaddomain.Lead) (string, error) {
	companyName := strings.ToLower(strings.TrimSpace(sourceLead.CompanyName))
	industry := strings.ToLower(strings.TrimSpace(sourceLead.Industry))
	categoryCode := "HOSPITAL"

	switch {
	case strings.Contains(companyName, "apotek") || strings.Contains(industry, "pharmacy"):
		categoryCode = "PHARMACY"
	case strings.Contains(companyName, "klinik") || strings.Contains(industry, "clinic"):
		categoryCode = "CLINIC"
	case strings.Contains(companyName, "rs") || strings.Contains(companyName, "rumah sakit") || strings.Contains(industry, "hospital"):
		categoryCode = "HOSPITAL"
	}

	var category account.Category
	if err := db.Where("code = ?", categoryCode).First(&category).Error; err != nil {
		return "", err
	}
	return category.ID, nil
}

func resolveConvertedLeadContactRoleID(db *gorm.DB, sourceLead *leaddomain.Lead) (string, error) {
	roleCode := "PIC"
	title := strings.ToLower(strings.TrimSpace(sourceLead.JobTitle))
	if strings.Contains(title, "owner") || strings.Contains(title, "director") || strings.Contains(title, "manager") || strings.Contains(title, "head") {
		roleCode = "MANAGER"
	}

	var role contactrole.ContactRole
	if err := db.Where("code = ?", roleCode).First(&role).Error; err != nil {
		return "", err
	}
	return role.ID, nil
}