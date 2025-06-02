# ContaMed - Chatbot WhatsApp

Sistema de chatbot para WhatsApp da ContaMed, uma plataforma de contabilidade digital para empresas médicas.

## Visão Geral

Este projeto implementa um backend para chatbot do WhatsApp que permite interações automatizadas com clientes da ContaMed. O sistema utiliza a API do WhatsApp Cloud e armazena as conversas em um banco de dados MongoDB.

## Funcionalidades

- Resposta automática a mensagens do WhatsApp
- Fluxo de conversação com menu de opções
- Armazenamento de conversas e mensagens
- Integração com a API oficial do WhatsApp

## Arquitetura

O projeto é estruturado seguindo os princípios da arquitetura hexagonal (ports and adapters) e arquitetura limpa:

- **Domain**: Contém as entidades de negócio e regras de domínio
- **Application**: Contém os casos de uso e regras de aplicação
- **Adapters**: Implementa as interfaces de entrada e saída
  - **Primary Adapters**: HTTP, CLI (interfaces de entrada)
  - **Secondary Adapters**: WhatsApp API, MongoDB (interfaces de saída)

## Tecnologias

- Linguagem: Go
- Framework Web: Chi Router
- Banco de Dados: MongoDB
- Teste: Testify, go-uber/mock
- Documentação: Swagger

## Requisitos

- Go 1.21 ou superior
- MongoDB
- Conta no WhatsApp Business API

## Configuração

1. Clone o repositório
2. Copie o arquivo `.env.example` para `.env` e configure as variáveis de ambiente
3. Execute `go mod download` para instalar as dependências
4. Execute `go run cmd/server/main.go` para iniciar o servidor

## Variáveis de Ambiente

```
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=localhost
ENV=development

# MongoDB Configuration
MONGODB_URI=mongodb+srv://user:password@cluster.mongodb.net/
MONGODB_DATABASE=database_name
MONGODB_TIMEOUT=10

# WhatsApp API Configuration
WHATSAPP_APP_ID=your_app_id
WHATSAPP_APP_SECRET=your_app_secret
WHATSAPP_ACCESS_TOKEN=your_access_token
WHATSAPP_PHONE_NUMBER_ID=your_phone_number_id
WHATSAPP_WEBHOOK_VERIFY_TOKEN=your_custom_webhook_verify_token

# Logging Configuration
LOG_LEVEL=debug
```

## Desenvolvimento

### Executando testes

```
go test ./...
```

### Gerando documentação Swagger

```
swag init -g cmd/server/main.go -o docs
```

## Fluxo do Chatbot

1. Ao receber uma mensagem, o chatbot responde com uma saudação personalizada (bom dia/tarde/noite) e apresenta as opções:
   - 1: Já tenho uma empresa médica constituída
   - 2: Quero abrir uma empresa
   - 3: Gostaria de tirar dúvidas
   - 4: Outros

2. Caso o usuário escolha a opção 2, o chatbot pergunta sobre o CRM:
   - 1: Já tenho CRM
   - 2: Ainda não possuo CRM

3. Em seguida, pergunta o Estado e Município de atuação.

## Licença

Este projeto é proprietário e confidencial.

## Contato

Para mais informações, entre em contato com a equipe de desenvolvimento da ContaMed. 