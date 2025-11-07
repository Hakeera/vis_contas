# vis_contas
Visualização de contas, boletos a pagar, agendados e pagos

---

## Modelos (Models)

As entidades principais são definidas em `internal/model` e mapeadas via **GORM**.

### Fatura
Representa uma transação financeira, armazenada em `financeiro.faturas`.

Campos principais:
- `Vencimento` — data de vencimento.
- `Valor` — valor monetário da fatura.
- `Situacao` — estado atual (padrão: `Pendente`).
- `Destinatario`, `Categoria`, `Empresa` — metadados da transação.

### User
Representa os usuários do sistema (`financeiro.users`).

Campos principais:
- `Username` — identificador único do usuário.
- `PasswordHash` — senha criptografada com bcrypt.
- `Role` — papel do usuário (ex: `user`, `admin`).
- `CreatedAt`, `UpdatedAt` — controle de timestamps.

Esses modelos são automaticamente migrados no banco quando `config.InitDB()` é executado.

---

## Rotas e Endpoints

As rotas são definidas em `internal/routes/routes.go` e organizadas em dois grupos:

### Rotas Públicas
| Método | Endpoint | Descrição |
|---------|-----------|-----------|
| **GET** | `/user` | Exibe a página de login |
| **POST** | `/login` | Autentica usuário e gera token JWT |

### Rotas Protegidas (JWT)
| Método | Endpoint | Descrição |
|---------|-----------|-----------|
| **GET** | `/` | Página inicial (dashboard) |
| **GET** | `/load_table` | Atualiza tabela de faturas (via HTMX) |
| **PUT** | `/alternar_sit/:id` | Alterna situação da fatura (Pendente ↔️ Paga) |
| **POST** | `/logout` | Encerra sessão e limpa cookie JWT |

### Middlewares Globais
- `Logger()` → Loga requisições HTTP.  
- `Recover()` → Trata panics e evita crash do servidor.  
- `CORS()` → Permite requisições entre origens distintas.

