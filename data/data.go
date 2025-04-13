package data

import (
	"database/sql/driver"
	"fmt"
	"time"

	"gorm.io/datatypes"
)

type Geography struct {
	WKT string
}

func (g Geography) Value() (driver.Value, error) {
	return g.WKT, nil
}

func (g *Geography) Scan(value interface{}) error {
	if str, ok := value.(string); ok {
		g.WKT = str
		return nil
	}
	return fmt.Errorf("failed to scan Geography: %v", value)
}

type CitizenProfile struct {
	UserID           string         `gorm:"type:uuid;primaryKey"`
	FirstName        string         `gorm:"not null"`
	LastName         string         `gorm:"not null"`
	PreferredModes   datatypes.JSON `gorm:"type:jsonb"`
	NotificationsOn  bool           `gorm:"not null;default:true"`
	SavedRoutes      datatypes.JSON `gorm:"type:jsonb"`
	ExperiencePoints int            `gorm:"default:0"`
}

type AuthorityAdmin struct {
	UserID       string         `gorm:"type:uuid;primaryKey"`
	Organization string         `gorm:"not null"`
	Privileges   datatypes.JSON `gorm:"type:jsonb;not null"`
}

type Route struct {
	RouteID       string    `gorm:"type:uuid;primaryKey"`
	Name          string    `gorm:"not null"`
	RouteType     string    `gorm:"not null;check:route_type IN ('Running', 'Bike', 'PublicTransport')"`
	PartnerID     string    `gorm:"type:uuid"`
	Geometry      Geography `gorm:"type:geography(LINESTRING,4326)"`
	EstimateTime  time.Duration
	IsEcoFriendly bool      `gorm:"not null;default:false"`
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`

	Partner    *Partner    `gorm:"foreignKey:PartnerID"`
	RouteNodes []RouteNode `gorm:"foreignKey:RouteID"`
}

type RouteNode struct {
	NodeID      string    `gorm:"type:uuid;primaryKey"`
	RouteID     string    `gorm:"type:uuid"`
	Sequence    int       `gorm:"not null"`
	Location    Geography `gorm:"type:geography(POINT,4326)"`
	Description string
	CreatedAt   time.Time `gorm:"autoCreateTime"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime"`

	Route Route `gorm:"constraint:OnDelete:CASCADE"`
}

type Achievement struct {
	AchievementID     string `gorm:"type:uuid;primaryKey"`
	PhotoURL          string
	Name              string `gorm:"not null"`
	Description       string
	ExperienceAwarded int `gorm:"not null"`
}

type CitizenAchievement struct {
	UserID        string    `gorm:"type:uuid;primaryKey"`
	AchievementID string    `gorm:"type:uuid;primaryKey"`
	AwardedAt     time.Time `gorm:"autoCreateTime"`

	Citizen     CitizenProfile `gorm:"foreignKey:UserID;references:UserID;constraint:OnDelete:CASCADE;belongsTo:CitizenProfile"`
	Achievement Achievement    `gorm:"foreignKey:AchievementID;references:AchievementID;constraint:OnDelete:CASCADE"`
}

type Notification struct {
	NotificationID string    `gorm:"type:uuid;primaryKey"`
	UserID         string    `gorm:"type:uuid"`
	Title          string    `gorm:"not null"`
	Message        string    `gorm:"not null"`
	Seen           bool      `gorm:"not null;default:false"`
	CreatedAt      time.Time `gorm:"autoCreateTime"`

	Citizen CitizenProfile `gorm:"foreignKey:UserID;references:UserID;constraint:OnDelete:CASCADE;belongsTo:CitizenProfile"`
}

type Suggestion struct {
	SuggestionID string    `gorm:"type:uuid;primaryKey"`
	SuggestedBy  string    `gorm:"type:uuid"`
	Message      string    `gorm:"not null"`
	Status       string    `gorm:"not null;check:status IN ('Pending', 'Approved', 'Rejected')"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`

	Citizen CitizenProfile `gorm:"foreignKey:SuggestedBy;references:UserID;belongsTo:CitizenProfile"`
}

type TrafficIncident struct {
	IncidentID   string         `gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ReportedBy   *string        `gorm:"type:uuid"`
	Location     datatypes.JSON `gorm:"type:jsonb"`
	Description  string
	IncidentType string    `gorm:"not null;check:incident_type IN ('Accident', 'Roadblock', 'BadWeather', 'Hazard', 'Traffic', 'Other')"`
	Status       string    `gorm:"not null;check:status IN ('Open', 'Resolved', 'Closed')"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`

	Citizen *CitizenProfile `gorm:"foreignKey:ReportedBy;references:UserID;constraint:OnDelete:CASCADE;"`
}

type Partner struct {
	PartnerID    string `gorm:"type:uuid;primaryKey"`
	UserID       string `gorm:"type:uuid;uniqueIndex;nullable"`
	CompanyName  string `gorm:"not null"`
	Type         string `gorm:"not null;check:type IN ('RideSharing', 'PublicTransport', 'Other')"`
	WebsiteURL   string
	ContactEmail string    `gorm:"not null"`
	Rating       float32   `gorm:"check:rating >= 0 AND rating <= 5"`
	IsActive     bool      `gorm:"not null;default:true"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`

	Citizen *CitizenProfile `gorm:"foreignKey:UserID;references:UserID;belongsTo:CitizenProfile"`
}

type TransportDataSource struct {
	SourceID      string         `gorm:"type:uuid;primaryKey"`
	PartnerID     string         `gorm:"type:uuid"`
	Name          string         `gorm:"not null"`
	Type          string         `gorm:"not null;check:type IN ('API', 'File', 'Other')"`
	AuthType      string         `gorm:"not null;check:auth_type IN ('None', 'OAuth')"`
	RequiresAuth  bool           `gorm:"not null;default:false"`
	Config        datatypes.JSON `gorm:"type:jsonb"`
	Enabled       bool           `gorm:"not null;default:true"`
	LastFetchedAt *time.Time
	ErrorLog      string
	CreatedAt     time.Time `gorm:"autoCreateTime"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime"`

	Partner *Partner `gorm:"foreignKey:PartnerID;belongsTo:Partner"`
}

type CitizenShortcut struct {
	ShortcutID string    `gorm:"type:uuid;primaryKey"`
	UserID     string    `gorm:"type:uuid"`
	Name       string    `gorm:"not null"`
	IconName   string    `gorm:"not null"`
	IconColor  string    `gorm:"not null"`
	Location   Geography `gorm:"type:geography(POINT,4326);not null"`

	Citizen *CitizenProfile `gorm:"foreignKey:UserID;references:UserID;belongsTo:CitizenProfile;constraint:OnDelete:CASCADE"`
}
