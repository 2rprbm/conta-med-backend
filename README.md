# 🤖 ContaMed WhatsApp Chatbot

> **Versão 0.3.0** - Sistema de chatbot inteligente para WhatsApp Business com persistência MongoDB completa

Um chatbot avançado para WhatsApp Business desenvolvido em Go, projetado especificamente para empresas de contabilidade médica. O sistema oferece fluxo de conversação inteligente, validações robustas e persistência completa de dados.

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

### 🎯 **Core do Chatbot**
- ✅ **Saudação inteligente** baseada no horário (Bom dia/Boa tarde/Boa noite)
- ✅ **Fluxo de conversação estruturado** com validações em tempo real
- ✅ **Menu interativo** com 4 opções principais
- ✅ **Coleta de informações** para abertura de empresas médicas
- ✅ **Validação de dados** brasileiros (telefone, estados, municípios)
- ✅ **Resumo automático** das informações coletadas

### 🗄️ **Persistência de Dados**
- ✅ **MongoDB integrado** com driver oficial
- ✅ **Repositórios completos** para mensagens e conversações
- ✅ **Índices otimizados** para performance
- ✅ **Health check** com monitoramento de conectividade
- ✅ **Graceful shutdown** com fechamento adequado de conexões

### 🔍 **Validações Implementadas**
- ✅ **Telefones brasileiros**: Móvel/fixo com/sem código do país
- ✅ **Estados**: Siglas (SP, RJ) ou nomes completos com normalização
- ✅ **Municípios**: Validação de formato e caracteres especiais
- ✅ **Opções de menu**: Validação rigorosa de entrada do usuário

### 🛡️ **Segurança e Qualidade**
- ✅ **Verificação de webhook** com assinatura HMAC SHA-256
- ✅ **Testes unitários** abrangentes (76% de cobertura)
- ✅ **Arquitetura hexagonal** para baixo acoplamento
- ✅ **Logs estruturados** para debugging e monitoramento

## 🏗️ Arquitetura

### **Arquitetura Hexagonal (Clean Architecture)**
```
├── cmd/server/             # Ponto de entrada da aplicação
├── internal/
│   ├── core/
│   │   ├── domain/         # Entidades e regras de negócio + validações
│   │   ├── services/       # Lógica de aplicação
│   │   └── ports/          # Interfaces (contratos)
│   ├── adapters/
│   │   ├── primary/        # Adaptadores de entrada (HTTP)
│   │   └── secondary/      # Adaptadores de saída (MongoDB, WhatsApp)
├── pkg/                    # Pacotes utilitários
│   ├── mongodb/            # Cliente MongoDB centralizado
│   └── logger/             # Sistema de logging
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
- **MongoDB 4.4** ou superior (ou MongoDB Atlas)
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
   # Crie um arquivo .env na raiz do projeto com:
   
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
   WHATSAPP_ACCESS_TOKEN=your_access_token
   WHATSAPP_PHONE_NUMBER_ID=your_phone_number_id
   WHATSAPP_WEBHOOK_VERIFY_TOKEN=your_webhook_token
   WHATSAPP_API_VERSION=v17.0
   WHATSAPP_BASE_URL=https://graph.facebook.com
   
   # Logging Configuration
   LOG_LEVEL=debug
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

| Variável | Descrição | Exemplo |
|----------|-----------|---------|
| `SERVER_PORT` | Porta do servidor | `8080` |
| `MONGODB_URI` | URI de conexão MongoDB | `mongodb://localhost:27017` |
| `MONGODB_DATABASE` | Nome do banco de dados | `medical_scheduler` |
| `WHATSAPP_ACCESS_TOKEN` | Token de acesso WhatsApp | `EAAJqTNxxCpY...` |
| `WHATSAPP_PHONE_NUMBER_ID` | ID do número WhatsApp | `123456789` |
| `WHATSAPP_WEBHOOK_VERIFY_TOKEN` | Token de verificação | `meu_token_secreto` |
| `LOG_LEVEL` | Nível de log | `debug`, `info`, `warn`, `error` |

## 👨‍💻 Desenvolvimento

### 🧪 Executando testes

```bash
# Todos os testes
go test ./...

# Com cobertura detalhada
go test -coverprofile="coverage.out" ./...
go tool cover -func="coverage.out"

# Coverage HTML report
go tool cover -html="coverage.out"
```

### 📊 Cobertura Atual de Testes

| Módulo | Cobertura | Status | Observações |
|--------|-----------|--------|-------------|
| **config** | 100.0% | ✅ | Completo |
| **handlers** | 76.0% | ✅ | Acima do objetivo |
| **middleware** | 100.0% | ✅ | Completo |
| **mongodb repos** | 79.6% | 🟡 | Próximo do objetivo (80%) |
| **mongodb client** | 39.7% | 🔴 | **Necessita melhoria** |
| **domain** | 74.6% | 🟡 | Necessita melhoria |
| **services** | 70.5% | 🟡 | Necessita melhoria |
| **logger** | 95.7% | ✅ | Excelente |
| **whatsapp** | 37.9% | 🔴 | **Prioridade para melhoria** |
| **GERAL** | **76%** | 🟡 | **Objetivo: 80%+** |

### 📚 Gerando documentação Swagger

```bash
swag init -g cmd/server/main.go -o docs
```

### 🔍 Health Check

```bash
# Health check detalhado
curl http://localhost:8080/health

# Ping simples
curl http://localhost:8080/ping
```

**Resposta do Health Check:**
```json
{
  "status": "healthy",
  "timestamp": "2024-01-01T00:00:00Z",
  "services": {
    "mongodb": "healthy"
  },
  "version": "v0.3.0"
}
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

## 🗄️ Estrutura do Banco de Dados

### **Collection: conversations**
```javascript
{
  "_id": ObjectId("..."),
  "phone_number": "+5511999999999",
  "status": "active|completed|pending",
  "state": "initial|main_menu|company_type_selection|...",
  "started_at": ISODate("..."),
  "last_updated_at": ISODate("..."),
  "ended_at": ISODate("..."), // nullable
  "user_selections": {
    "main_menu": "2",
    "crm_selection": "1",
    "state": "SP",
    "city": "São Paulo"
  },
  "consultant_id": "consultant_123" // opcional
}
```

### **Collection: messages**
```javascript
{
  "_id": ObjectId("..."),
  "conversation_id": "conv_123",
  "phone_number": "+5511999999999",
  "content": "Olá, gostaria de abrir uma empresa",
  "type": "text|image|document|location",
  "direction": "inbound|outbound",
  "timestamp": ISODate("..."),
  "metadata": {} // opcional
}
```

### **Índices Criados**
- **conversations**: `phone_number`, `phone_number + status`, `last_updated_at`
- **messages**: `conversation_id`, `phone_number`, `timestamp`, `conversation_id + timestamp`

## 🚀 Deploy

### **Desenvolvimento**
```bash
go run cmd/server/main.go
```

### **Produção**
```bash
# Build
go build -o conta-med-chatbot cmd/server/main.go

# Execute
./conta-med-chatbot
```

### **Docker** (Futuro)
```dockerfile
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o conta-med-chatbot cmd/server/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/conta-med-chatbot .
CMD ["./conta-med-chatbot"]
```

## 📈 Roadmap

### ✅ **v0.3.0 - Atual**
- Integração MongoDB completa
- Health check avançado
- Validações robustas
- Testes unitários abrangentes

### 🔄 **v0.4.0 - Em Desenvolvimento**
- Sistema de gerenciamento de atendimento
- Interface para funcionários
- Transferência de conversas entre atendentes
- Notificações em tempo real

### 🔮 **v0.5.0 - Planejado**
- Cache Redis para performance
- Métricas e monitoramento
- Backup automático
- API REST para integração

## 🤝 Contribuição

1. Fork o projeto
2. Crie uma branch para sua feature (`git checkout -b feature/AmazingFeature`)
3. Commit suas mudanças (`git commit -m 'Add some AmazingFeature'`)
4. Push para a branch (`git push origin feature/AmazingFeature`)
5. Abra um Pull Request

## 📄 Licença

Este projeto está sob a licença MIT. Veja o arquivo [LICENSE](LICENSE) para mais detalhes.

## 📞 Suporte

- **Email**: suporte@contamed.com.br
- **WhatsApp**: +55 11 99999-9999
- **Issues**: [GitHub Issues](https://github.com/2rprbm/conta-med-backend/issues)

---

**Desenvolvido com ❤️ pela equipe ContaMed** 