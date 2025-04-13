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

func pointsToLineString(points []LatLng) string {
	if len(points) == 0 {
		return "LINESTRING EMPTY"
	}
	coords := make([]string, len(points))
	for i, p := range points {
		coords[i] = fmt.Sprintf("%f %f", p.Lng, p.Lat)
	}
	return "LINESTRING(" + strings.Join(coords, ", ") + ")"
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
        WHERE status = 'Open'
          AND ST_DWithin(
              ST_GeomFromText(location #>> '{}')::geography,
              ST_GeomFromText(?)::geography,
              ?
          )
    `
	if err := db.Raw(query, lineStringWKT, tolerance).Scan(&incidents).Error; err != nil && err != gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("error fetching incidents: %w", err)
	}

	return incidents, nil
}
