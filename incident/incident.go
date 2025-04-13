package incident

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/team-evian-fiicode25/business-logic/data"
	"github.com/team-evian-fiicode25/business-logic/database"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type LatLng struct {
	Lat float64 `json:"latitude"`
	Lng float64 `json:"longitude"`
}

func ReportTrafficIncident(userID, locationWKT, description, incidentType string) (*data.TrafficIncident, error) {
	var reportedBy *string
	if userID != "" {
		reportedBy = &userID
	}

	locationJSON, err := json.Marshal(locationWKT)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal location WKT: %w", err)
	}

	incident := &data.TrafficIncident{
		IncidentID:   uuid.New().String(),
		ReportedBy:   reportedBy,
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

func pointsToLineString(points []LatLng) string {
	if len(points) == 0 {
		return "LINESTRING EMPTY"
	}
	coords := make([]string, 0, len(points))
	for _, p := range points {
		coords = append(coords, fmt.Sprintf("%f %f", p.Lng, p.Lat))
	}
	return "LINESTRING(" + strings.Join(coords, ", ") + ")"
}

func GetOpenTrafficIncidentsByRoute(points []LatLng, tolerance float64) ([]data.TrafficIncident, error) {
	db := database.GetDB()
	var incidents []data.TrafficIncident

	lineStringWKT := pointsToLineString(points)
	if lineStringWKT == "LINESTRING EMPTY" {
		return nil, fmt.Errorf("no points provided for route")
	}

	query := `
		SELECT *
		FROM traffic_incidents
		WHERE status = ?
		  AND ST_DWithin(
		      ST_GeomFromText(location::jsonb #>> '{}', 4326),
		      ST_GeomFromText(?, 4326),
		      ?
		  )
	`

	if err := db.Raw(query, "Open", lineStringWKT, tolerance).Scan(&incidents).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("error fetching incidents: %w", err)
	}

	return incidents, nil
}
