package domain

import "time"

// MessageType defines the type of message
type MessageType string

const (
	// TextMessage represents a text message
	TextMessage MessageType = "text"
	// ImageMessage represents an image message
	ImageMessage MessageType = "image"
	// DocumentMessage represents a document message
	DocumentMessage MessageType = "document"
	// LocationMessage represents a location message
	LocationMessage MessageType = "location"
)

// MessageDirection defines the direction of the message
type MessageDirection string

const (
	// Inbound message from user to system
	Inbound MessageDirection = "inbound"
	// Outbound message from system to user
	Outbound MessageDirection = "outbound"
)

// Message represents a message in the chatbot system
type Message struct {
	ID            string           `json:"id" bson:"_id,omitempty"`
	ConversationID string           `json:"conversation_id" bson:"conversation_id"`
	PhoneNumber   string           `json:"phone_number" bson:"phone_number"`
	Content       string           `json:"content" bson:"content"`
	Type          MessageType      `json:"type" bson:"type"`
	Direction     MessageDirection `json:"direction" bson:"direction"`
	Timestamp     time.Time        `json:"timestamp" bson:"timestamp"`
	MetaData      map[string]interface{} `json:"metadata,omitempty" bson:"metadata,omitempty"`
}

// NewMessage creates a new message
func NewMessage(conversationID, phoneNumber, content string, msgType MessageType, direction MessageDirection) *Message {
	return &Message{
		ConversationID: conversationID,
		PhoneNumber:   phoneNumber,
		Content:       content,
		Type:          msgType,
		Direction:     direction,
		Timestamp:     time.Now(),
		MetaData:      make(map[string]interface{}),
	}
}

// NewInboundMessage creates a new inbound message
func NewInboundMessage(conversationID, phoneNumber, content string, msgType MessageType) *Message {
	return NewMessage(conversationID, phoneNumber, content, msgType, Inbound)
}

// NewOutboundMessage creates a new outbound message
func NewOutboundMessage(conversationID, phoneNumber, content string, msgType MessageType) *Message {
	return NewMessage(conversationID, phoneNumber, content, msgType, Outbound)
} 