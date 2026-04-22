package interfaces

import (
	"github.com/gilabs/crm-healthcare/api/internal/domain/contact"
)

// ContactRepository defines the interface for contact repository
type ContactRepository interface {
	// FindByID finds a contact by ID
	FindByID(id string) (*contact.Contact, error)
	
	// List returns a list of contacts with pagination
	List(req *contact.ListContactsRequest) ([]contact.Contact, int64, error)
	
	// Create creates a new contact
	Create(contact *contact.Contact) error
	
	// Update updates a contact
	Update(contact *contact.Contact) error
	
	// Delete soft deletes a contact
	Delete(id string) error
	
	// FindByAccountID finds contacts by account ID
	FindByAccountID(accountID string) ([]contact.Contact, error)
}

