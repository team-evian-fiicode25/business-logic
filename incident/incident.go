package user

import (
	"encoding/json"
	"fmt"

	"github.com/team-evian-fiicode25/business-logic/data"
	"github.com/team-evian-fiicode25/business-logic/database"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

func ReportTrafficIncident(userID, locationWKT, description, incidentType string) (*data.TrafficIncident, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID cannot be empty")
	}

	locationJSON, err := json.Marshal(locationWKT)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal location WKT: %w", err)
	}

	incident := &data.TrafficIncident{
		IncidentID:   "",
		ReportedBy:   userID,
		Location:     datatypes.JSON(locationJSON),
		Description:  description,
		IncidentType: incidentType,
		Status:       "Open",
	}

	db := database.GetDB()
	if err := db.Create(incident).Error; err != nil {
		return nil, err
	}

	return incident, nil
}

func GetOpenTrafficIncidentsByRoute(polylineWKT string, tolerance float64) ([]data.TrafficIncident, error) {
	db := database.GetDB()
	var incidents []data.TrafficIncident

	query := `
		SELECT *
		FROM traffic_incidents
		WHERE status = ?
		  AND ST_DWithin(
		      location::geometry,
		      ST_GeomFromText(?, 4326),
		      ?
		  )
	`
	if err := db.Raw(query, "Open", polylineWKT, tolerance).Scan(&incidents).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("error fetching incidents: %w", err)
	}

	return incidents, nil
}
