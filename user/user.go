package user

import (
	"encoding/json"
	"fmt"

	"github.com/team-evian-fiicode25/business-logic/data"
	"github.com/team-evian-fiicode25/business-logic/database"
	"gorm.io/datatypes"
)

const (
	CitizenProfileType = "citizen"
	AuthorityAdminType = "authority"
)

func CreateCitizenProfile(userID string, firstName string, lastName string) (*data.CitizenProfile, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID cannot be empty")
	}

	citizenProfile := &data.CitizenProfile{
		UserID:           userID,
		FirstName:        firstName,
		LastName:         lastName,
		NotificationsOn:  true,
		ExperiencePoints: 0,
	}

	db := database.GetDB()
	err := db.Create(&citizenProfile).Error
	if err != nil {
		return nil, err
	}

	return citizenProfile, nil
}

func CreateAuthorityAdmin(userID string, organization string, privileges map[string]interface{}) (*data.AuthorityAdmin, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID cannot be empty")
	}

	if organization == "" {
		return nil, fmt.Errorf("organization cannot be empty")
	}

	if privileges == nil {
		privileges = map[string]interface{}{
			"canModifyRoutes": true,
		}
	}

	privilegesBytes, err := json.Marshal(privileges)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal privileges to JSON: %w", err)
	}

	privilegesJSON := datatypes.JSON(privilegesBytes)

	authorityAdmin := &data.AuthorityAdmin{
		UserID:       userID,
		Organization: organization,
		Privileges:   privilegesJSON,
	}

	db := database.GetDB()
	err = db.Create(&authorityAdmin).Error
	if err != nil {
		return nil, err
	}

	return authorityAdmin, nil
}

func GetCitizenProfile(userID string) (*data.CitizenProfile, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID cannot be empty")
	}

	db := database.GetDB()
	var profile data.CitizenProfile
	result := db.Where("user_id = ?", userID).First(&profile)

	if result.Error != nil {
		return nil, result.Error
	}

	return &profile, nil
}

func GetAuthorityAdmin(userID string) (*data.AuthorityAdmin, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID cannot be empty")
	}

	db := database.GetDB()
	var admin data.AuthorityAdmin
	result := db.Where("user_id = ?", userID).First(&admin)

	if result.Error != nil {
		return nil, result.Error
	}

	return &admin, nil
}
