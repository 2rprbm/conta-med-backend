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

	// Validate phone number format
	if !domain.ValidatePhoneNumber(phoneNumber) {
		s.log.Warn("Invalid phone number format", logger.Fields{
			"phone_number": phoneNumber,
		})
		return fmt.Errorf("invalid phone number format: %s", phoneNumber)
	}

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
	
	// Get main menu options from domain
	menuOptions := domain.GetMainMenuOptions()
	mainMenuOptions := make([]string, 0, len(menuOptions))
	for key, description := range menuOptions {
		mainMenuOptions = append(mainMenuOptions, fmt.Sprintf("%s- %s", key, description))
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
	// Validate main menu option
	if !domain.ValidateMainMenuOption(message) {
		// Invalid selection, ask again with available options
		menuOptions := domain.GetMainMenuOptions()
		optionsList := make([]string, 0, len(menuOptions))
		for key, description := range menuOptions {
			optionsList = append(optionsList, fmt.Sprintf("%s- %s", key, description))
		}

		invalidMsg := "Opção inválida. Por favor, escolha uma das opções disponíveis:\n" + strings.Join(optionsList, "\n")
		if err := s.whatsappService.SendTextMessage(ctx, conversation.PhoneNumber, invalidMsg); err != nil {
			return fmt.Errorf("failed to send invalid option message: %w", err)
		}

		// Save the outbound message
		outboundMsg := domain.NewOutboundMessage(conversation.ID, conversation.PhoneNumber, invalidMsg, domain.TextMessage)
		return s.messageRepo.Save(ctx, outboundMsg)
	}

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

		// Get CRM options from domain
		crmOptions := domain.GetCRMOptions()
		crmOptionsList := make([]string, 0, len(crmOptions))
		for key, description := range crmOptions {
			crmOptionsList = append(crmOptionsList, fmt.Sprintf("%s- %s", key, description))
		}
		
		if err := s.whatsappService.SendOptionsMessage(ctx, conversation.PhoneNumber, "Digite a opção adequada:", crmOptionsList); err != nil {
			return fmt.Errorf("failed to send CRM options: %w", err)
		}

		// Save the outbound message
		outboundMsg := domain.NewOutboundMessage(conversation.ID, conversation.PhoneNumber, "Digite a opção adequada:", domain.TextMessage)
		return s.messageRepo.Save(ctx, outboundMsg)
	}
	
	return nil
}

func (s *chatbotService) processCRMSelection(ctx context.Context, conversation *domain.Conversation, message string) error {
	// Validate CRM option
	if !domain.ValidateCRMOption(message) {
		// Invalid selection, ask again with available options
		crmOptions := domain.GetCRMOptions()
		optionsList := make([]string, 0, len(crmOptions))
		for key, description := range crmOptions {
			optionsList = append(optionsList, fmt.Sprintf("%s- %s", key, description))
		}

		invalidMsg := "Opção inválida. Por favor, escolha uma das opções disponíveis:\n" + strings.Join(optionsList, "\n")
		if err := s.whatsappService.SendTextMessage(ctx, conversation.PhoneNumber, invalidMsg); err != nil {
			return fmt.Errorf("failed to send invalid CRM option message: %w", err)
		}

		// Save the outbound message
		outboundMsg := domain.NewOutboundMessage(conversation.ID, conversation.PhoneNumber, invalidMsg, domain.TextMessage)
		return s.messageRepo.Save(ctx, outboundMsg)
	}

	// Save the user's selection
	conversation.AddUserSelection("crm_selection", message)

	// Update state to ask for state
	conversation.UpdateState(domain.StateSelection)
	if err := s.conversationRepo.Update(ctx, conversation); err != nil {
		return fmt.Errorf("failed to update conversation state: %w", err)
	}

	// Ask for the user's state
	stateMsg := "Por favor, informe o estado (UF) onde você pretende atuar (ex: SP, RJ, MG):"
	if err := s.whatsappService.SendTextMessage(ctx, conversation.PhoneNumber, stateMsg); err != nil {
		return fmt.Errorf("failed to send state question: %w", err)
	}

	// Save the outbound message
	outboundMsg := domain.NewOutboundMessage(conversation.ID, conversation.PhoneNumber, stateMsg, domain.TextMessage)
	return s.messageRepo.Save(ctx, outboundMsg)
}

func (s *chatbotService) processStateSelection(ctx context.Context, conversation *domain.Conversation, message string) error {
	// Validate Brazilian state
	if !domain.ValidateBrazilianState(message) {
		// Invalid state, ask again with examples
		invalidMsg := "Estado inválido. Por favor, informe um estado brasileiro válido (ex: SP, RJ, MG) ou o nome completo do estado:"
		if err := s.whatsappService.SendTextMessage(ctx, conversation.PhoneNumber, invalidMsg); err != nil {
			return fmt.Errorf("failed to send invalid state message: %w", err)
		}

		// Save the outbound message
		outboundMsg := domain.NewOutboundMessage(conversation.ID, conversation.PhoneNumber, invalidMsg, domain.TextMessage)
		return s.messageRepo.Save(ctx, outboundMsg)
	}

	// Normalize and save the user's selection
	normalizedState := domain.NormalizeBrazilianState(message)
	conversation.AddUserSelection("state", normalizedState)

	// Update state to ask for city
	conversation.UpdateState(domain.CitySelection)
	if err := s.conversationRepo.Update(ctx, conversation); err != nil {
		return fmt.Errorf("failed to update conversation state: %w", err)
	}

	// Ask for the user's city
	cityMsg := "Por favor, informe o município onde você pretende atuar:"
	if err := s.whatsappService.SendTextMessage(ctx, conversation.PhoneNumber, cityMsg); err != nil {
		return fmt.Errorf("failed to send city question: %w", err)
	}

	// Save the outbound message
	outboundMsg := domain.NewOutboundMessage(conversation.ID, conversation.PhoneNumber, cityMsg, domain.TextMessage)
	return s.messageRepo.Save(ctx, outboundMsg)
}

func (s *chatbotService) processCitySelection(ctx context.Context, conversation *domain.Conversation, message string) error {
	// Validate city name
	if !domain.ValidateCityName(message) {
		// Invalid city name, ask again
		invalidMsg := "Nome do município inválido. Por favor, informe um nome de município válido:"
		if err := s.whatsappService.SendTextMessage(ctx, conversation.PhoneNumber, invalidMsg); err != nil {
			return fmt.Errorf("failed to send invalid city message: %w", err)
		}

		// Save the outbound message
		outboundMsg := domain.NewOutboundMessage(conversation.ID, conversation.PhoneNumber, invalidMsg, domain.TextMessage)
		return s.messageRepo.Save(ctx, outboundMsg)
	}

	// Save the user's selection
	conversation.AddUserSelection("city", strings.TrimSpace(message))

	return s.finishConversation(ctx, conversation)
}

func (s *chatbotService) finishConversation(ctx context.Context, conversation *domain.Conversation) error {
	// Update state to waiting for consultant
	conversation.UpdateState(domain.WaitingForConsultant)
	if err := s.conversationRepo.Update(ctx, conversation); err != nil {
		return fmt.Errorf("failed to update conversation state: %w", err)
	}

	// Create a summary of the conversation
	summaryMsg := s.createConversationSummary(conversation)

	// Thank the user and inform them they'll be contacted by a consultant
	thankYouMsg := fmt.Sprintf("Obrigado pelas informações!\n\n%s\n\nUm de nossos consultores entrará em contato com você em breve para dar continuidade ao seu atendimento.", summaryMsg)
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

// createConversationSummary creates a summary of the conversation for the user
func (s *chatbotService) createConversationSummary(conversation *domain.Conversation) string {
	summary := "Resumo das suas informações:\n"
	
	if mainMenu, exists := conversation.UserSelections["main_menu"]; exists {
		menuOptions := domain.GetMainMenuOptions()
		if description, found := menuOptions[mainMenu]; found {
			summary += fmt.Sprintf("• Opção selecionada: %s\n", description)
		}
	}
	
	if crmSelection, exists := conversation.UserSelections["crm_selection"]; exists {
		crmOptions := domain.GetCRMOptions()
		if description, found := crmOptions[crmSelection]; found {
			summary += fmt.Sprintf("• CRM: %s\n", description)
		}
	}
	
	if state, exists := conversation.UserSelections["state"]; exists {
		states := domain.GetBrazilianStates()
		if stateName, found := states[state]; found {
			summary += fmt.Sprintf("• Estado: %s (%s)\n", stateName, state)
		} else {
			summary += fmt.Sprintf("• Estado: %s\n", state)
		}
	}
	
	if city, exists := conversation.UserSelections["city"]; exists {
		summary += fmt.Sprintf("• Município: %s\n", city)
	}
	
	return summary
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