package domain

import "time"

// ConversationStatus defines the status of a conversation
type ConversationStatus string

const (
	// Active conversation is ongoing
	Active ConversationStatus = "active"
	// Completed conversation is finished
	Completed ConversationStatus = "completed"
	// Pending conversation is waiting for further action
	Pending ConversationStatus = "pending"
)

// ConversationState defines the state of the conversation in the chatbot flow
type ConversationState string

const (
	// Initial state when conversation starts
	Initial ConversationState = "initial"
	// MainMenu state when displaying main options
	MainMenu ConversationState = "main_menu"
	// CompanyTypeSelection state for company type selection
	CompanyTypeSelection ConversationState = "company_type_selection"
	// CRMSelection state for CRM selection
	CRMSelection ConversationState = "crm_selection"
	// StateSelection state for selecting Brazilian state
	StateSelection ConversationState = "state_selection"
	// CitySelection state for selecting city
	CitySelection ConversationState = "city_selection"
	// WaitingForConsultant state when waiting for a human consultant
	WaitingForConsultant ConversationState = "waiting_for_consultant"
	// ConsultantAssigned state when a consultant is assigned
	ConsultantAssigned ConversationState = "consultant_assigned"
)

// Conversation represents a conversation between a user and the chatbot
type Conversation struct {
	ID              string             `json:"id" bson:"_id,omitempty"`
	PhoneNumber     string             `json:"phone_number" bson:"phone_number"`
	Status          ConversationStatus `json:"status" bson:"status"`
	State           ConversationState  `json:"state" bson:"state"`
	StartedAt       time.Time          `json:"started_at" bson:"started_at"`
	LastUpdatedAt   time.Time          `json:"last_updated_at" bson:"last_updated_at"`
	EndedAt         *time.Time         `json:"ended_at,omitempty" bson:"ended_at,omitempty"`
	UserSelections  map[string]string  `json:"user_selections" bson:"user_selections"`
	ConsultantID    string             `json:"consultant_id,omitempty" bson:"consultant_id,omitempty"`
}

// NewConversation creates a new conversation
func NewConversation(phoneNumber string) *Conversation {
	now := time.Now()
	return &Conversation{
		PhoneNumber:    phoneNumber,
		Status:         Active,
		State:          Initial,
		StartedAt:      now,
		LastUpdatedAt:  now,
		UserSelections: make(map[string]string),
	}
}

// UpdateState updates the conversation state and last updated time
func (c *Conversation) UpdateState(state ConversationState) {
	c.State = state
	c.LastUpdatedAt = time.Now()
}

// AddUserSelection adds a user selection to the conversation
func (c *Conversation) AddUserSelection(key, value string) {
	c.UserSelections[key] = value
	c.LastUpdatedAt = time.Now()
}

// AssignConsultant assigns a consultant to the conversation
func (c *Conversation) AssignConsultant(consultantID string) {
	c.ConsultantID = consultantID
	c.UpdateState(ConsultantAssigned)
}

// Complete marks the conversation as completed
func (c *Conversation) Complete() {
	c.Status = Completed
	now := time.Now()
	c.EndedAt = &now
	c.LastUpdatedAt = now
} 