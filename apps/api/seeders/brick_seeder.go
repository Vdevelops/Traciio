package seeders

import (
	"log"
	"time"

	"github.com/gilabs/crm-healthcare/api/internal/database"
	brickdomain "github.com/gilabs/crm-healthcare/api/internal/domain/brick"
	"github.com/gilabs/crm-healthcare/api/internal/domain/group"
	monthlytargetdomain "github.com/gilabs/crm-healthcare/api/internal/domain/monthly_target"
	"github.com/gilabs/crm-healthcare/api/internal/domain/role"
	"github.com/gilabs/crm-healthcare/api/internal/domain/user"
	"golang.org/x/crypto/bcrypt"
)

// brickManagerMapping defines which bricks belong to which manager
// and which sales reps are assigned to each brick.
type brickManagerMapping struct {
ManagerEmail string
ManagerName  string
Bricks       []brickDefinition
}

type brickDefinition struct {
Name        string
Code        string
Description string
Province    string
Regency     string
// SalesEmails are the sales rep emails assigned to this brick
SalesEmails []string
}

// getSeedMappings returns the deterministic brick-to-manager-to-salesrep mapping.
// Each manager has distinct bricks and distinct sales reps, ensuring proper RBAC isolation.
func getSeedMappings() []brickManagerMapping {
return []brickManagerMapping{
{
ManagerEmail: "salesmanager@example.com",
ManagerName:  "Sales Manager",
Bricks: []brickDefinition{
{
Name:        "Solo, Jawa Tengah",
Code:        "SOLO-JT",
Description: "Brick untuk wilayah Solo, Jawa Tengah",
Province:    "Jawa Tengah",
Regency:     "Surakarta",
SalesEmails: []string{"sales@example.com", "sales1@example.com"},
},
{
Name:        "Yogyakarta, DI Yogyakarta",
Code:        "YOGYA-DIY",
Description: "Brick untuk wilayah Yogyakarta, DI Yogyakarta",
Province:    "DI Yogyakarta",
Regency:     "Yogyakarta",
SalesEmails: []string{"sales2@example.com", "sales3@example.com"},
},
{
Name:        "Medan, Sumatera Utara",
Code:        "MDN-SU",
Description: "Brick untuk wilayah Medan, Sumatera Utara",
Province:    "Sumatera Utara",
Regency:     "Medan",
SalesEmails: nil,
},
},
},
{
ManagerEmail: "salesmanager2@example.com",
ManagerName:  "Sales Manager 2",
Bricks: []brickDefinition{
{
Name:        "Semarang, Jawa Tengah",
Code:        "SMG-JT",
Description: "Brick untuk wilayah Semarang, Jawa Tengah",
Province:    "Jawa Tengah",
Regency:     "Semarang",
SalesEmails: []string{"sales4@example.com", "sales5@example.com"},
},
{
Name:        "Jakarta Pusat, DKI Jakarta",
Code:        "JKT-PST",
Description: "Brick untuk wilayah Jakarta Pusat, DKI Jakarta",
Province:    "DKI Jakarta",
Regency:     "Jakarta Pusat",
SalesEmails: []string{"sales6@example.com"},
},
{
Name:        "Jakarta Selatan, DKI Jakarta",
Code:        "JKT-SEL",
Description: "Brick untuk wilayah Jakarta Selatan, DKI Jakarta",
Province:    "DKI Jakarta",
Regency:     "Jakarta Selatan",
SalesEmails: nil,
},
{
Name:        "Denpasar, Bali",
Code:        "DPS-BALI",
Description: "Brick untuk wilayah Denpasar, Bali",
Province:    "Bali",
Regency:     "Denpasar",
SalesEmails: nil,
},
},
},
{
ManagerEmail: "salesmanager3@example.com",
ManagerName:  "Sales Manager 3",
Bricks: []brickDefinition{
{
Name:        "Surabaya, Jawa Timur",
Code:        "SBY-JT",
Description: "Brick untuk wilayah Surabaya, Jawa Timur",
Province:    "Jawa Timur",
Regency:     "Surabaya",
SalesEmails: []string{"sales7@example.com", "sales8@example.com"},
},
{
Name:        "Malang, Jawa Timur",
Code:        "MLG-JT",
Description: "Brick untuk wilayah Malang, Jawa Timur",
Province:    "Jawa Timur",
Regency:     "Malang",
SalesEmails: []string{"sales9@example.com"},
},
{
Name:        "Sidoarjo, Jawa Timur",
Code:        "SDA-JT",
Description: "Brick untuk wilayah Sidoarjo, Jawa Timur",
Province:    "Jawa Timur",
Regency:     "Sidoarjo",
SalesEmails: nil,
},
{
Name:        "Jakarta Barat, DKI Jakarta",
Code:        "JKT-BAR",
Description: "Brick untuk wilayah Jakarta Barat, DKI Jakarta",
Province:    "DKI Jakarta",
Regency:     "Jakarta Barat",
SalesEmails: nil,
},
{
Name:        "Jakarta Timur, DKI Jakarta",
Code:        "JKT-TIM",
Description: "Brick untuk wilayah Jakarta Timur, DKI Jakarta",
Province:    "DKI Jakarta",
Regency:     "Jakarta Timur",
SalesEmails: nil,
},
{
Name:        "Bandung, Jawa Barat",
Code:        "BDG-JB",
Description: "Brick untuk wilayah Bandung, Jawa Barat",
Province:    "Jawa Barat",
Regency:     "Bandung",
SalesEmails: []string{"sales10@example.com"},
},
{
Name:        "Bekasi, Jawa Barat",
Code:        "BKS-JB",
Description: "Brick untuk wilayah Bekasi, Jawa Barat",
Province:    "Jawa Barat",
Regency:     "Bekasi",
SalesEmails: nil,
},
},
},
}
}

// SeedBricks seeds initial brick data with deterministic manager and sales rep assignments.
// Each sales manager gets distinct bricks and distinct sales reps, ensuring proper RBAC isolation:
//   - salesmanager@example.com  -> Solo, Yogyakarta, Medan              -> [sales@, sales1@, sales2@, sales3@]
//   - salesmanager2@example.com -> Semarang, Jakarta PST/SEL, Denpasar  -> [sales4@, sales5@, sales6@]
//   - salesmanager3@example.com -> Surabaya, Malang, Sidoarjo, JKT BAR/TIM, Bandung, Bekasi -> [sales7@..sales10@]
func SeedBricks() error {
// Check if bricks already exist
var count int64
database.DB.Model(&brickdomain.Brick{}).Count(&count)
if count > 0 {
log.Println("Bricks already seeded, skipping...")
return nil
}

// Get roles
var salesManagerRole role.Role
if err := database.DB.Where("code = ?", "sales_manager").First(&salesManagerRole).Error; err != nil {
log.Printf("Warning: Sales manager role not found: %v", err)
return err
}

var salesRole role.Role
if err := database.DB.Where("code = ?", "sales").First(&salesRole).Error; err != nil {
log.Printf("Warning: Sales role not found: %v", err)
return err
}

// Get groups for user creation fallback
var salesGroup group.Group
var salesGroupID *string
if err := database.DB.Where("code = ?", "SALES").First(&salesGroup).Error; err == nil {
salesGroupID = &salesGroup.ID
}

// Hash password for any new users
hashedPassword, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
if err != nil {
return err
}

now := time.Now()
currentYear := now.Year()
currentMonth := int(now.Month())
mappings := getSeedMappings()
brickCount := 0

for _, mapping := range mappings {
// Resolve or create the manager user
manager, mErr := resolveOrCreateUser(mapping.ManagerEmail, mapping.ManagerName, salesManagerRole.ID, salesGroupID, string(hashedPassword))
if mErr != nil {
log.Printf("Error resolving manager %s: %v", mapping.ManagerEmail, mErr)
return mErr
}
log.Printf("Manager resolved: %s (ID: %s)", manager.Email, manager.ID)

for i, bDef := range mapping.Bricks {
// Create the brick
brick := brickdomain.Brick{
Name:        bDef.Name,
Code:        bDef.Code,
Description: bDef.Description,
Province:    bDef.Province,
Regency:     bDef.Regency,
ManagerID:   &manager.ID,
Status:      "active",
}
if cErr := database.DB.Create(&brick).Error; cErr != nil {
log.Printf("Error creating brick %s: %v", brick.Name, cErr)
continue
}
log.Printf("Created brick: %s (ID: %s, Manager: %s)", brick.Name, brick.ID, manager.Email)
brickCount++

// Create monthly target for the brick
baseTarget := int64(30000000000)
targetVariation := int64(10000000000 * (i % 6))
targetAmount := baseTarget + targetVariation

monthlyTarget := monthlytargetdomain.MonthlyTarget{
BrickID:      &brick.ID,
UserID:       nil,
GroupID:      nil,
Year:         currentYear,
Month:        currentMonth,
TargetAmount: targetAmount,
}
if tErr := database.DB.Create(&monthlyTarget).Error; tErr != nil {
log.Printf("Warning: Failed to create monthly target for brick %s: %v", brick.Name, tErr)
}

// Assign sales reps to this brick deterministically
for _, salesEmail := range bDef.SalesEmails {
salesUser, sErr := resolveOrCreateUser(salesEmail, "", salesRole.ID, salesGroupID, string(hashedPassword))
if sErr != nil {
log.Printf("Warning: Failed to resolve sales user %s: %v", salesEmail, sErr)
continue
}
if uErr := database.DB.Model(&user.User{}).Where("id = ?", salesUser.ID).Update("brick_id", brick.ID).Error; uErr != nil {
log.Printf("Warning: Failed to assign %s to brick %s: %v", salesEmail, brick.Name, uErr)
} else {
log.Printf("Assigned sales user %s -> brick %s (manager: %s)", salesEmail, brick.Name, manager.Email)
}
}

// Assign the manager to their first brick (so manager has a brick_id too)
if i == 0 {
if uErr := database.DB.Model(&user.User{}).Where("id = ?", manager.ID).Update("brick_id", brick.ID).Error; uErr != nil {
log.Printf("Warning: Failed to assign manager %s to brick %s: %v", manager.Email, brick.Name, uErr)
} else {
log.Printf("Assigned manager %s -> brick %s (primary brick)", manager.Email, brick.Name)
}
}
}
}

log.Printf("Seeded %d bricks with deterministic manager and sales rep assignments", brickCount)
return nil
}

// resolveOrCreateUser fetches an existing user by email or creates a new one.
func resolveOrCreateUser(email, name, roleID string, groupID *string, hashedPassword string) (*user.User, error) {
var u user.User
if err := database.DB.Where("email = ?", email).First(&u).Error; err == nil {
return &u, nil
}

// Derive name from email if not provided
if name == "" {
name = email
}

u = user.User{
Email:     email,
Password:  hashedPassword,
Name:      name,
AvatarURL: "https://api.dicebear.com/7.x/lorelei/svg?seed=" + email,
RoleID:    roleID,
GroupID:   groupID,
Status:    "active",
}
if err := database.DB.Create(&u).Error; err != nil {
return nil, err
}
log.Printf("Created user: %s", u.Email)
return &u, nil
}

// AssignBricksToUsers re-assigns all active sales users to their designated bricks
// using the same deterministic mapping as SeedBricks. This function overwrites any
// previous brick_id assignment to ensure consistency with the seed mappings.
func AssignBricksToUsers() error {
mappings := getSeedMappings()
assigned := 0

for _, mapping := range mappings {
// Get the manager
var manager user.User
if err := database.DB.Where("email = ?", mapping.ManagerEmail).First(&manager).Error; err != nil {
log.Printf("Warning: Manager %s not found, skipping assignment: %v", mapping.ManagerEmail, err)
continue
}

// Assign manager to their first brick
if len(mapping.Bricks) > 0 {
var firstBrick brickdomain.Brick
if err := database.DB.Where("code = ?", mapping.Bricks[0].Code).First(&firstBrick).Error; err == nil {
result := database.DB.Model(&user.User{}).Where("email = ?", mapping.ManagerEmail).Update("brick_id", firstBrick.ID)
if result.Error == nil && result.RowsAffected > 0 {
log.Printf("Assigned manager %s -> brick %s (primary brick)", mapping.ManagerEmail, firstBrick.Name)
assigned++
}
}
}

for _, bDef := range mapping.Bricks {
if len(bDef.SalesEmails) == 0 {
continue
}

// Find the brick by code
var brick brickdomain.Brick
if err := database.DB.Where("code = ?", bDef.Code).First(&brick).Error; err != nil {
log.Printf("Warning: Brick %s not found, skipping: %v", bDef.Code, err)
continue
}

// Assign each designated sales rep to this brick
for _, salesEmail := range bDef.SalesEmails {
result := database.DB.Model(&user.User{}).Where("email = ?", salesEmail).Update("brick_id", brick.ID)
if result.Error != nil {
log.Printf("Warning: Failed to assign %s to brick %s: %v", salesEmail, brick.Name, result.Error)
} else if result.RowsAffected > 0 {
log.Printf("Assigned %s -> brick %s (manager: %s)", salesEmail, brick.Name, manager.Email)
assigned++
}
}
}
}

log.Printf("Deterministically assigned %d sales users to bricks", assigned)
return nil
}
