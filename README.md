# 📱 ContaMed - Chatbot WhatsApp

<div align="center">
  
  ![Version](https://img.shields.io/badge/version-0.2.0-blue.svg?cacheSeconds=2592000)
  ![Go Version](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat&logo=go&logoColor=white)
  ![Coverage](https://img.shields.io/badge/coverage-75%25-yellowgreen)
  ![License](https://img.shields.io/badge/license-Proprietary-red)

  Sistema de chatbot para WhatsApp da ContaMed, uma plataforma de contabilidade digital para empresas médicas.
</div>

## 🚀 Visão Geral

Este projeto implementa um backend para chatbot do WhatsApp que permite interações automatizadas com clientes da ContaMed. O sistema utiliza a API do WhatsApp Cloud e armazena as conversas em um banco de dados MongoDB com validações robustas de entrada e fluxo de conversação inteligente.

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

#### ✅ Sprint 3: Repositórios e Validações
- **Repositórios MongoDB completos** com operações CRUD
- **Sistema de validações robusto** para entrada do usuário
- **Melhorias no fluxo de conversação** com tratamento de erros
- **Cobertura de testes expandida** (75% no geral)
- **Testes de integração** com MongoDB

### 🔜 Próximas Sprints
- Melhoria da cobertura de testes (objetivo: 80%+)
- Implementação de índices MongoDB para performance
- Testes end-to-end completos
- Sistema de métricas e monitoramento

## ✨ Funcionalidades

### 💬 Chatbot Inteligente
- Resposta automática com saudação personalizada por horário
- Fluxo de conversação com menu de opções dinâmico
- Validação robusta de entrada do usuário
- Tratamento de erros com mensagens personalizadas
- Resumo automático das informações coletadas

### 🔒 Validações Implementadas
- **Telefones brasileiros**: Móvel e fixo, com/sem código do país
- **Estados brasileiros**: Siglas (SP, RJ) e nomes completos
- **Municípios**: Validação de formato e caracteres especiais
- **Opções de menu**: Validação de escolhas válidas
- **Normalização automática**: Estados para formato padrão

### 💾 Persistência de Dados
- Armazenamento completo de conversas e mensagens
- Busca eficiente por conversações ativas
- Histórico completo de mensagens ordenado por timestamp
- Geração automática de IDs únicos

### 🧪 Qualidade e Testes
- **75% de cobertura geral** de testes
- Testes unitários abrangentes
- Testes de integração com MongoDB
- Mocks centralizados para facilitar manutenção

## 🏗️ Arquitetura

O projeto é estruturado seguindo os princípios da arquitetura hexagonal (ports and adapters) e arquitetura limpa:

- **Domain** 📊 - Entidades de negócio, validações e regras de domínio
- **Application** ⚙️ - Casos de uso e regras de aplicação
- **Adapters** 🔄 - Implementa as interfaces de entrada e saída
  - **Primary Adapters** 📥 - HTTP, CLI (interfaces de entrada)
  - **Secondary Adapters** 📤 - WhatsApp API, MongoDB (interfaces de saída)

### 📦 Estrutura de Pastas

```
conta-med-backend/
├── cmd/server/             # Ponto de entrada da aplicação
├── config/                 # Configurações da aplicação
├── internal/
│   ├── core/
│   │   ├── domain/         # Entidades e validações de negócio
│   │   ├── ports/          # Interfaces (contratos)
│   │   └── services/       # Lógica de aplicação
│   ├── adapters/
│   │   ├── primary/        # Adaptadores de entrada (HTTP)
│   │   └── secondary/      # Adaptadores de saída (MongoDB, WhatsApp)
├── pkg/                    # Pacotes utilitários
└── docs/                   # Documentação Swagger
```

## 🛠️ Tecnologias

- 💻 **Linguagem**: [Go 1.21+](https://golang.org/)
- 🌐 **Framework Web**: [Chi Router](https://github.com/go-chi/chi)
- 🗄️ **Banco de Dados**: [MongoDB](https://www.mongodb.com/) com driver oficial
- 🧪 **Testes**: [Testify](https://github.com/stretchr/testify) para assertions
- 📚 **Documentação**: [Swagger](https://swagger.io/) para API docs
- 🔍 **Validações**: Sistema personalizado de validações de domínio

## 📋 Requisitos

- **Go 1.21** ou superior
- **MongoDB 4.4** ou superior
- **Conta no WhatsApp Business API**
- **Conexão com internet** para integração WhatsApp

## 🔧 Configuração

1. **Clone o repositório**
   ```bash
   git clone https://github.com/2rprbm/conta-med-backend.git
   cd conta-med-backend
   ```

2. **Configure as variáveis de ambiente**
   ```bash
   cp .env.example .env
   # Edite o arquivo .env com suas configurações
   ```

3. **Instale as dependências**
   ```bash
   go mod download
   ```

4. **Execute os testes**
   ```bash
   go test ./...
   ```

5. **Inicie o servidor**
   ```bash
   go run cmd/server/main.go
   ```

## 🔐 Variáveis de Ambiente

```bash
# Server Configuration
SERVER_PORT=8080
SERVER_HOST=localhost
ENV=development

# MongoDB Configuration
MONGODB_URI=mongodb+srv://2rprbm:2rprbm@cluster0.f4v7rrn.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0
MONGODB_DATABASE=medical_scheduler
MONGODB_TIMEOUT=10

# WhatsApp API Configuration
WHATSAPP_APP_ID=679828301875862
WHATSAPP_APP_SECRET=dac37016c2ffaf8655344289faf3b39e
WHATSAPP_ACCESS_TOKEN=your_temporary_token
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
# Todos os testes
go test ./...

# Com cobertura detalhada
go test -cover ./...

# Coverage HTML report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### 📊 Cobertura Atual de Testes

| Módulo | Cobertura | Status |
|--------|-----------|--------|
| config | 100.0% | ✅ |
| handlers | 82.4% | ✅ |
| middleware | 100.0% | ✅ |
| mongodb repos | 79.6% | 🟡 |
| domain | 74.6% | 🟡 |
| services | 70.5% | 🟡 |
| logger | 95.7% | ✅ |
| whatsapp | 37.9% | 🔴 |

### 📚 Gerando documentação Swagger

```bash
swag init -g cmd/server/main.go -o docs
```

## 🤖 Fluxo do Chatbot

### 1. 👋 Saudação Inicial
- Saudação personalizada por horário (Bom dia/Boa tarde/Boa noite)
- Apresentação da ContaMed
- Menu de opções principais

### 2. 📋 Menu Principal
O usuário pode escolher entre:
- **1️⃣ Já tenho uma empresa médica constituída** → Direcionado para consultor
- **2️⃣ Quero abrir uma empresa** → Continua para coleta de informações
- **3️⃣ Gostaria de tirar dúvidas** → Direcionado para consultor
- **4️⃣ Outros** → Direcionado para consultor

### 3. 🏥 Fluxo "Abrir Empresa" (Opção 2)

#### 📋 Informações sobre CRM
- **1️⃣ Já tenho CRM**
- **2️⃣ Ainda não possuo CRM**

#### 🗺️ Localização
1. **Estado**: Aceita siglas (SP, RJ) ou nomes completos (São Paulo, Rio de Janeiro)
2. **Município**: Validação de formato e caracteres especiais

### 4. ✅ Finalização
- Resumo das informações coletadas
- Confirmação de direcionamento para consultor
- Armazenamento completo da conversa

### 🔍 Validações Implementadas

- **Números de telefone**: +5511999999999, 5511999999999, 11999999999
- **Estados brasileiros**: SP, São Paulo, rj, Rio de Janeiro
- **Nomes de cidades**: São Paulo, Rio de Janeiro, Belo Horizonte
- **Opções de menu**: 1, 2, 3, 4 (com trim automático)

## 🚀 Performance e Escalabilidade

- **Conexões MongoDB** com pool de conexões configurável
- **Timeouts apropriados** para todas as operações
- **Logs estruturados** para monitoramento
- **Validações client-side** para reduzir carga do servidor
- **Geração eficiente de IDs** usando ObjectID do MongoDB

## 📜 Licença

Este projeto é proprietário e confidencial da ContaMed.

## 🤝 Contribuição

Para contribuir com o projeto:

1. Mantenha a cobertura de testes acima de 70%
2. Siga os padrões de arquitetura hexagonal
3. Implemente validações adequadas
4. Adicione testes para novas funcionalidades
5. Mantenha a documentação atualizada

## 📞 Contato

Para mais informações, entre em contato com a equipe de desenvolvimento da ContaMed. 