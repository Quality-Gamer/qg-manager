package model

type ManagerMatch struct {
	Id             string              `json:"id"`
	ChallengeId    int                 `json:"challenge_id"`
	UserId         int                 `json:"user_id"`
	ProjectId      int                 `json:"project_id"`
	Week           int                 `json:"weak"`
	Progress       float64             `json:"progress"`
	Money          int                 `json:"money"`
	Time           int                 `json:"time"`
	Team           Team                `json:"team"`
	License        License             `json:"license"`
	Action         Action              `json:"action"`
	Event          Event               `json:"event"`
	Occurrence     []ManagerOccurrence `json:"manager_occurrence"`
	UserOccurrence []UserOccurrence    `json:"user_occurrence"`
}

type Team struct {
	Backend      int `json:"backend"`
	Frontend     int `json:"frontend"`
	Designer     int `json:"designer"`
	Tester       int `json:"tester"`
	ProductOwner int `json:"product_owner"`
	RiskAnalyst  int `json:"risk_analyst"`
}

type License struct {
	Ide            int `json:"ide"`
	DesignSoftware int `json:"design_software"`
}

type Action struct {
	Scrum           int `json:"scrum"`
	Delivery        int `json:"delivery"`
	CustomerContact int `json:"customer_contact"`
	RiskAnalysis    int `json:"risk_analysis"`
}

type Event struct {
	Name   string `json:"name"`
	Status string `json:"status"`
}
