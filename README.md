# 📱 ContaMed - Chatbot WhatsApp

<div align="center">
  
  ![Version](https://img.shields.io/badge/version-0.1.0-blue.svg?cacheSeconds=2592000)
  ![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go&logoColor=white)
  ![License](https://img.shields.io/badge/license-Proprietary-red)

  Sistema de chatbot para WhatsApp da ContaMed, uma plataforma de contabilidade digital para empresas médicas.
</div>

## 🚀 Visão Geral

Este projeto implementa um backend para chatbot do WhatsApp que permite interações automatizadas com clientes da ContaMed. O sistema utiliza a API do WhatsApp Cloud e armazena as conversas em um banco de dados MongoDB.

## 📝 Status do Projeto

### Sprints Concluídas

#### ✅ Sprint 1: Configuração Inicial
- Estrutura do projeto em arquitetura hexagonal
- Configuração inicial (logger, config, servidor HTTP)
- Integração com ambiente de desenvolvimento

#### ✅ Sprint 2: Core do Chatbot
- Domínio do chatbot (mensagens, conversações, estados)
- Portas e adaptadores para serviços e repositórios
- Implementação da lógica de fluxo de conversação
- Handler de webhook para integração com WhatsApp

#### 🔄 Sprint 3: Repositórios e Persistência (Em andamento)
- Implementação dos repositórios MongoDB
- Armazenamento e recuperação de conversas
- Testes de integração

## ✨ Funcionalidades

- 💬 Resposta automática a mensagens do WhatsApp
- 🔄 Fluxo de conversação com menu de opções
- 💾 Armazenamento de conversas e mensagens
- 🔌 Integração com a API oficial do WhatsApp

## 🏗️ Arquitetura

O projeto é estruturado seguindo os princípios da arquitetura hexagonal (ports and adapters) e arquitetura limpa:

- **Domain** 📊 - Contém as entidades de negócio e regras de domínio
- **Application** ⚙️ - Contém os casos de uso e regras de aplicação
- **Adapters** 🔄 - Implementa as interfaces de entrada e saída
  - **Primary Adapters** 📥 - HTTP, CLI (interfaces de entrada)
  - **Secondary Adapters** 📤 - WhatsApp API, MongoDB (interfaces de saída)

## 🛠️ Tecnologias

- 💻 **Linguagem**: [Go](https://golang.org/)
- 🌐 **Framework Web**: [Chi Router](https://github.com/go-chi/chi)
- 🗄️ **Banco de Dados**: [MongoDB](https://www.mongodb.com/)
- 🧪 **Teste**: [Testify](https://github.com/stretchr/testify), [go-uber/mock](https://github.com/uber-go/mock)
- 📚 **Documentação**: [Swagger](https://swagger.io/)

## 📋 Requisitos

- Go 1.21 ou superior
- MongoDB
- Conta no WhatsApp Business API

## 🔧 Configuração

1. Clone o repositório
   ```bash
   git clone https://github.com/2rprbm/conta-med-backend.git
   cd conta-med-backend
   ```

2. Copie o arquivo `.env.example` para `.env` e configure as variáveis de ambiente
   ```bash
   cp .env.example .env
   # Edite o arquivo .env com suas configurações
   ```

3. Execute `go mod download` para instalar as dependências
   ```bash
   go mod download
   ```

4. Execute `go run cmd/server/main.go` para iniciar o servidor
   ```bash
   go run cmd/server/main.go
   ```

## 🔐 Variáveis de Ambiente

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
WHATSAPP_API_VERSION=v17.0
WHATSAPP_BASE_URL=https://graph.facebook.com

# Logging Configuration
LOG_LEVEL=debug
```

## 👨‍💻 Desenvolvimento

### 🧪 Executando testes

```bash
go test ./...
```

### 📚 Gerando documentação Swagger

```bash
swag init -g cmd/server/main.go -o docs
```

## 🤖 Fluxo do Chatbot

1. Ao receber uma mensagem, o chatbot responde com uma saudação personalizada (bom dia/tarde/noite) e apresenta as opções:
   - 1️⃣ Já tenho uma empresa médica constituída
   - 2️⃣ Quero abrir uma empresa
   - 3️⃣ Gostaria de tirar dúvidas
   - 4️⃣ Outros

2. Caso o usuário escolha a opção 2, o chatbot pergunta sobre o CRM:
   - 1️⃣ Já tenho CRM
   - 2️⃣ Ainda não possuo CRM

3. Em seguida, pergunta o Estado e Município de atuação.

## 📜 Licença

Este projeto é proprietário e confidencial.

## 📞 Contato

Para mais informações, entre em contato com a equipe de desenvolvimento da ContaMed. 