package models

import "database/sql"

// Region represents a geographic region
type Region struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// District represents a district within a region
type District struct {
	ID       int    `json:"id"`
	RegionID int    `json:"regionId"`
	Name     string `json:"name"`
}

// Subcounty represents a subcounty within a district
type Subcounty struct {
	ID         int    `json:"id"`
	DistrictID int    `json:"districtId"`
	Name       string `json:"name"`
}

// Facility represents a health facility
type Facility struct {
	ID          int    `json:"id"`
	SubcountyID int    `json:"subcountyId"`
	Name        string `json:"name"`
}

// HealthWorker represents a health worker
type HealthWorker struct {
	ID           int            `json:"id"`
	FullName     string         `json:"fullName"`
	Email        sql.NullString `json:"email"`
	PhoneNumber  sql.NullString `json:"phoneNumber"`
	FacilityID   int            `json:"facilityId"`
	FacilityName string         `json:"facilityName,omitempty"`
	CreatedAt    string         `json:"createdAt,omitempty"`
	UpdatedAt    string         `json:"updatedAt,omitempty"`
}

// AssessmentType represents a type of assessment
type AssessmentType struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
	Code string `json:"code"`
}

// ThematicArea represents a thematic area within an assessment
type ThematicArea struct {
	ID               int    `json:"id"`
	AssessmentTypeID int    `json:"assessmentTypeId"`
	Name             string `json:"name"`
	DisplayOrder     int    `json:"displayOrder"`
}

// Question represents a question in an assessment
type Question struct {
	ID             int    `json:"id"`
	ThematicAreaID int    `json:"thematicAreaId"`
	QuestionText   string `json:"text"`
	ScoreWeight    int    `json:"scoreWeight"`
	IsCritical     bool   `json:"isCritical"`
	IsImportant    bool   `json:"isImportant"`
	DisplayOrder   int    `json:"displayOrder"`
}

// Assessment represents a completed assessment
type Assessment struct {
	ID                 int            `json:"id"`
	HealthWorkerID     int            `json:"healthWorkerId"`
	HealthWorkerName   string         `json:"healthWorkerName,omitempty"`
	FacilityID         int            `json:"facilityId"`
	AssessmentTypeID   int            `json:"assessmentTypeId"`
	AssessorName       sql.NullString `json:"assessorName"`
	ClientName         sql.NullString `json:"clientName"`
	Notes              sql.NullString `json:"notes"`
	TotalPossibleScore int            `json:"totalPossible"`
	AchievedScore      int            `json:"achieved"`
	PercentageScore    float64        `json:"percentage"`
	PerformanceLevel   string         `json:"performanceLevel"`
	CreatedAt          string         `json:"createdAt"`
	FacilityName       string         `json:"facilityName"`
	AssessmentType     string         `json:"assessmentType"`
}

// AssessmentResponse represents a response to a question
type AssessmentResponse struct {
	ID               int    `json:"id"`
	AssessmentID     int    `json:"assessmentId"`
	QuestionID       int    `json:"questionId"`
	Response         string `json:"response"` // "Yes", "No", or "NA"
	PointsEarned     int    `json:"pointsEarned"`
	QuestionText     string `json:"questionText"`
	ScoreWeight      int    `json:"scoreWeight"`
	ThematicAreaID   int    `json:"thematicAreaId"`
	ThematicAreaName string `json:"thematicAreaName"`
}

// ThematicAreaScore represents the score for a thematic area
type ThematicAreaScore struct {
	ID               int     `json:"id"`
	AssessmentID     int     `json:"assessmentId"`
	ThematicAreaID   int     `json:"thematicAreaId"`
	ThematicAreaName string  `json:"name"`
	PossibleScore    int     `json:"possible"`
	AchievedScore    int     `json:"achieved"`
	PercentageScore  float64 `json:"percentage"`
}
