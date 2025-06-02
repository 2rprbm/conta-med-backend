package services

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/2rprbm/conta-med-backend/internal/core/domain"
	"github.com/2rprbm/conta-med-backend/internal/core/ports"
	"github.com/2rprbm/conta-med-backend/pkg/logger"
)

type chatbotService struct {
	messageRepo      ports.MessageRepository
	conversationRepo ports.ConversationRepository
	whatsappService  ports.WhatsAppService
	log              logger.Logger
}

// NewChatbotService creates a new instance of the chatbot service
func NewChatbotService(
	messageRepo ports.MessageRepository,
	conversationRepo ports.ConversationRepository,
	whatsappService ports.WhatsAppService,
	log logger.Logger,
) ports.ChatbotService {
	return &chatbotService{
		messageRepo:      messageRepo,
		conversationRepo: conversationRepo,
		whatsappService:  whatsappService,
		log:              log,
	}
}

// HandleIncomingMessage processes an incoming message and generates a response
func (s *chatbotService) HandleIncomingMessage(ctx context.Context, phoneNumber, message string) error {
	s.log.Info("Handling incoming message", logger.Fields{
		"phone_number": phoneNumber,
		"message":      message,
	})

	// Get or create a conversation
	conversation, err := s.GetOrCreateConversation(ctx, phoneNumber)
	if err != nil {
		return fmt.Errorf("failed to get or create conversation: %w", err)
	}

	// Create and save the incoming message
	inboundMsg := domain.NewInboundMessage(conversation.ID, phoneNumber, message, domain.TextMessage)
	if err := s.messageRepo.Save(ctx, inboundMsg); err != nil {
		return fmt.Errorf("failed to save incoming message: %w", err)
	}

	// If this is a new conversation, send welcome message
	if conversation.State == domain.Initial {
		return s.SendWelcomeMessage(ctx, phoneNumber)
	}

	// Process the user's response based on the current conversation state
	return s.ProcessUserSelection(ctx, phoneNumber, message)
}

// GetOrCreateConversation retrieves an existing conversation or creates a new one
func (s *chatbotService) GetOrCreateConversation(ctx context.Context, phoneNumber string) (*domain.Conversation, error) {
	// Try to find an active conversation for this phone number
	conversation, err := s.conversationRepo.FindActiveByPhoneNumber(ctx, phoneNumber)
	if err == nil && conversation != nil {
		return conversation, nil
	}

	// Create a new conversation if none exists
	newConversation := domain.NewConversation(phoneNumber)
	if err := s.conversationRepo.Save(ctx, newConversation); err != nil {
		return nil, fmt.Errorf("failed to save new conversation: %w", err)
	}

	return newConversation, nil
}

// SendWelcomeMessage sends the initial welcome message to a user
func (s *chatbotService) SendWelcomeMessage(ctx context.Context, phoneNumber string) error {
	conversation, err := s.conversationRepo.FindActiveByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return fmt.Errorf("failed to find active conversation: %w", err)
	}

	// Update conversation state
	conversation.UpdateState(domain.MainMenu)
	if err := s.conversationRepo.Update(ctx, conversation); err != nil {
		return fmt.Errorf("failed to update conversation state: %w", err)
	}

	// Get appropriate greeting based on time of day
	greeting := getGreetingByTime(time.Now())
	
	// Main menu options
	mainMenuOptions := []string{
		"1- Já tenho uma empresa médica constituida",
		"2- Quero abrir uma empresa",
		"3- Gostaria de tirar dúvidas",
		"4- Outros",
	}

	// Compose welcome message
	welcomeMsg := fmt.Sprintf("%s. A conta med é uma plataforma de contabilidade digital para empresas médicas.\nLogo você será redirecionado para um de nossos consultores, mas antes, digite a opção que melhor te atende:", greeting)

	// Send the welcome message with options
	if err := s.whatsappService.SendOptionsMessage(ctx, phoneNumber, welcomeMsg, mainMenuOptions); err != nil {
		return fmt.Errorf("failed to send welcome message: %w", err)
	}

	// Save the outbound message
	outboundMsg := domain.NewOutboundMessage(conversation.ID, phoneNumber, welcomeMsg, domain.TextMessage)
	return s.messageRepo.Save(ctx, outboundMsg)
}

// ProcessUserSelection processes a user's menu selection
func (s *chatbotService) ProcessUserSelection(ctx context.Context, phoneNumber, message string) error {
	conversation, err := s.conversationRepo.FindActiveByPhoneNumber(ctx, phoneNumber)
	if err != nil {
		return fmt.Errorf("failed to find active conversation: %w", err)
	}

	// Trim and normalize the message
	message = strings.TrimSpace(message)

	switch conversation.State {
	case domain.MainMenu:
		return s.processMainMenuSelection(ctx, conversation, message)
	case domain.CompanyTypeSelection:
		return s.processCRMSelection(ctx, conversation, message)
	case domain.CRMSelection:
		return s.processStateSelection(ctx, conversation, message)
	case domain.StateSelection:
		return s.processCitySelection(ctx, conversation, message)
	case domain.CitySelection:
		return s.finishConversation(ctx, conversation)
	default:
		// For any other state, pass the conversation to a human consultant
		return s.transferToConsultant(ctx, conversation)
	}
}

func (s *chatbotService) processMainMenuSelection(ctx context.Context, conversation *domain.Conversation, message string) error {
	// Save the user's selection
	conversation.AddUserSelection("main_menu", message)

	// Process based on the selection
	switch message {
	case "1", "3", "4":
		// Options 1, 3, and 4 go directly to a consultant
		return s.transferToConsultant(ctx, conversation)
	case "2":
		// Option 2 - "Quero abrir uma empresa" - continues to the CRM selection
		conversation.UpdateState(domain.CompanyTypeSelection)
		if err := s.conversationRepo.Update(ctx, conversation); err != nil {
			return fmt.Errorf("failed to update conversation state: %w", err)
		}

		// Send CRM selection options
		crmOptions := []string{
			"1- Já tenho CRM",
			"2- Ainda nao possuo CRM",
		}
		
		if err := s.whatsappService.SendOptionsMessage(ctx, conversation.PhoneNumber, "Digite a opção adequada:", crmOptions); err != nil {
			return fmt.Errorf("failed to send CRM options: %w", err)
		}

		// Save the outbound message
		outboundMsg := domain.NewOutboundMessage(conversation.ID, conversation.PhoneNumber, "Digite a opção adequada:", domain.TextMessage)
		return s.messageRepo.Save(ctx, outboundMsg)
	default:
		// Invalid selection, ask again
		if err := s.whatsappService.SendTextMessage(ctx, conversation.PhoneNumber, "Opção inválida. Por favor, escolha uma das opções disponíveis (1-4):"); err != nil {
			return fmt.Errorf("failed to send invalid option message: %w", err)
		}

		// Save the outbound message
		outboundMsg := domain.NewOutboundMessage(conversation.ID, conversation.PhoneNumber, "Opção inválida. Por favor, escolha uma das opções disponíveis (1-4):", domain.TextMessage)
		return s.messageRepo.Save(ctx, outboundMsg)
	}
}

func (s *chatbotService) processCRMSelection(ctx context.Context, conversation *domain.Conversation, message string) error {
	// Save the user's selection
	conversation.AddUserSelection("crm_selection", message)

	// Update state to ask for state
	conversation.UpdateState(domain.StateSelection)
	if err := s.conversationRepo.Update(ctx, conversation); err != nil {
		return fmt.Errorf("failed to update conversation state: %w", err)
	}

	// Ask for the user's state
	if err := s.whatsappService.SendTextMessage(ctx, conversation.PhoneNumber, "Por favor, informe o estado (UF) onde você pretende atuar:"); err != nil {
		return fmt.Errorf("failed to send state question: %w", err)
	}

	// Save the outbound message
	outboundMsg := domain.NewOutboundMessage(conversation.ID, conversation.PhoneNumber, "Por favor, informe o estado (UF) onde você pretende atuar:", domain.TextMessage)
	return s.messageRepo.Save(ctx, outboundMsg)
}

func (s *chatbotService) processStateSelection(ctx context.Context, conversation *domain.Conversation, message string) error {
	// Save the user's selection
	conversation.AddUserSelection("state", message)

	// Update state to ask for city
	conversation.UpdateState(domain.CitySelection)
	if err := s.conversationRepo.Update(ctx, conversation); err != nil {
		return fmt.Errorf("failed to update conversation state: %w", err)
	}

	// Ask for the user's city
	if err := s.whatsappService.SendTextMessage(ctx, conversation.PhoneNumber, "Por favor, informe o município onde você pretende atuar:"); err != nil {
		return fmt.Errorf("failed to send city question: %w", err)
	}

	// Save the outbound message
	outboundMsg := domain.NewOutboundMessage(conversation.ID, conversation.PhoneNumber, "Por favor, informe o município onde você pretende atuar:", domain.TextMessage)
	return s.messageRepo.Save(ctx, outboundMsg)
}

func (s *chatbotService) processCitySelection(ctx context.Context, conversation *domain.Conversation, message string) error {
	// Save the user's selection
	conversation.AddUserSelection("city", message)

	return s.finishConversation(ctx, conversation)
}

func (s *chatbotService) finishConversation(ctx context.Context, conversation *domain.Conversation) error {
	// Update state to waiting for consultant
	conversation.UpdateState(domain.WaitingForConsultant)
	if err := s.conversationRepo.Update(ctx, conversation); err != nil {
		return fmt.Errorf("failed to update conversation state: %w", err)
	}

	// Thank the user and inform them they'll be contacted by a consultant
	thankYouMsg := "Obrigado pelas informações! Um de nossos consultores entrará em contato com você em breve para dar continuidade ao seu atendimento."
	if err := s.whatsappService.SendTextMessage(ctx, conversation.PhoneNumber, thankYouMsg); err != nil {
		return fmt.Errorf("failed to send thank you message: %w", err)
	}

	// Save the outbound message
	outboundMsg := domain.NewOutboundMessage(conversation.ID, conversation.PhoneNumber, thankYouMsg, domain.TextMessage)
	return s.messageRepo.Save(ctx, outboundMsg)
}

func (s *chatbotService) transferToConsultant(ctx context.Context, conversation *domain.Conversation) error {
	// Update state to waiting for consultant
	conversation.UpdateState(domain.WaitingForConsultant)
	if err := s.conversationRepo.Update(ctx, conversation); err != nil {
		return fmt.Errorf("failed to update conversation state: %w", err)
	}

	// Inform the user they'll be contacted by a consultant
	transferMsg := "Você será atendido por um de nossos consultores em breve. Agradecemos pela paciência."
	if err := s.whatsappService.SendTextMessage(ctx, conversation.PhoneNumber, transferMsg); err != nil {
		return fmt.Errorf("failed to send transfer message: %w", err)
	}

	// Save the outbound message
	outboundMsg := domain.NewOutboundMessage(conversation.ID, conversation.PhoneNumber, transferMsg, domain.TextMessage)
	return s.messageRepo.Save(ctx, outboundMsg)
}

// getGreetingByTime returns an appropriate greeting based on the time of day
func getGreetingByTime(t time.Time) string {
	hour := t.Hour()

	if hour >= 5 && hour < 12 {
		return "Bom dia"
	} else if hour >= 12 && hour < 18 {
		return "Boa tarde"
	} else {
		return "Boa noite"
	}
} 