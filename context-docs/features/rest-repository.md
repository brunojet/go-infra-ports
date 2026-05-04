# Feature: RestRepository

## Objetivo

Implementar o adaptador concreto `restRepository` em `internal/repositories/rest/` que satisfaz a porta `RestRepository` definida em `pkg/repositories/rest/contracts/contracts.go`.

O `restRepository` usa:
- `HttpClient` (outbound port) como transporte HTTP
- `RestRegistry` (outbound port) para marshaling de request/response

---

## Decisões Arquiteturais

### Direção de dependência

```
go-infra-adapters  →  importa  →  go-infra-ports
```

`go-infra-ports` **nunca** importa `go-infra-adapters`. A direção correta é: adaptadores dependem de portas.

### Contrato HttpClient em go-infra-ports

`go-infra-ports` define seu próprio contrato outbound `HttpClient` em `pkg/http_clients/contracts/`:

```go
type HttpClient interface {
    Do(ctx context.Context, req *http.Request) (*http.Response, error)
}
```

Idêntico ao shape do `netHttpClient` de `go-infra-adapters`. O duck typing do Go garante compatibilidade em runtime na aplicação final — sem nenhum `require` entre os dois módulos.

`HttpClient` é um **outbound port** (dependência que o `RestRepository` precisa que alguém forneça), exatamente como `RestRegistry`. Não é um contrato de adapter — é um contrato de dependência do port.

### Bootstrap na aplicação final

A aplicação final é a "ponte" entre as duas libs:

```go
// na aplicação (main.go / wire / di container)
httpClient, _ := adapters.NewNetHttpClient(...)  // go-infra-adapters
registry      := ports.NewRestRegistry(...)       // go-infra-ports
repo, _       := ports.NewRestRepository(
    ports.WithHttpClient(httpClient),             // duck typing — funciona sem import entre libs
    ports.WithRegistry(registry),
)
svc := myservice.New(repo)
```

### Outros

- Sem dependência nova em `go.mod`
- HTTP method mapping hard-coded por operação (não configurável): `Create`→POST, `List`→GET, `Get`→GET, `Update`→PUT, `Save`→PATCH, `Delete`→DELETE
- URL templates usam sintaxe `{key}` resolvida via `RequestContext.Identifiers`
- `HttpClient` e `RestRegistry` são opções **obrigatórias**; `baseURL`, `paths`, `headers` são opcionais
- `RepositoryOption` segue semântica `func(*repositoryOptions) error` (permite validação inline, ex.: `WithBaseURL` rejeitar string vazia)
- `pkg/http_clients/` expõe apenas `contracts/` e `mocks/` — sem `api.go` (nenhuma factory neste repo)

---

## Plano de Implementação

### Phase 0 — Contrato HttpClient

| Arquivo | Conteúdo |
|---|---|
| `pkg/http_clients/contracts/contracts.go` | Interface `HttpClient { Do(ctx, *http.Request) (*http.Response, error) }` |
| `pkg/http_clients/contracts/contracts_test.go` | Compile-time assertion + surface check |
| `pkg/http_clients/mocks/mocks.go` | `MockHttpClient` via GoMock |

### Phase 1 — Internal: tipos, config, erros, helpers

| Arquivo | Conteúdo |
|---|---|
| `rest_repository_type.go` | Alias `HttpClient`, constantes `OperationCreate/List/Get/Update/Save/Delete`, struct `restRepository{client, registry, opts}` |
| `rest_repository_config.go` | `RepositoryOption func(*repositoryOptions) error`, struct `repositoryOptions`, `newRepositoryOptions`, `WithHttpClient`, `WithRegistry`, `WithBaseURL`, `WithPath`, `WithHeader` |
| `rest_repository_error.go` | Sentinelas: `ErrMissingHttpClient`, `ErrMissingRegistry`, `ErrBuildRequest`, `ErrExecuteRequest`, `ErrReadResponseBody`, `ErrResolveRequest`, `ErrResolveResponse` + funções de wrap |
| `rest_repository_helper.go` | `interpolatePath(template, Identifiers) string` (`{key}` → valor), `applyQueryParams(*url.URL, url.Values)` |

### Phase 2 — Internal: métodos locais e implementação

| Arquivo | Conteúdo |
|---|---|
| `rest_repository_local.go` | Métodos privados: `resolveURL`, `buildHTTPRequest`, `executeRequest`, `readBody`, `mapResponse`, `mapResponses` |
| `rest_repository.go` | `NewRestRepository(opts...) (RestRepository, error)`, 6 métodos da interface `RestRepository` |

### Phase 3 — API pública

| Arquivo | Conteúdo |
|---|---|
| `pkg/repositories/api.go` *(modify)* | Adicionar: `HttpClient`, `RepositoryOption`, `NewRestRepository`, `WithHttpClient`, `WithRegistry`, `WithBaseURL`, `WithPath`, `WithHeader` |
| `pkg/repositories/api_test.go` *(modify)* | Cobertura dos novos símbolos re-exportados |

### Phase 4 — Testes

| Arquivo | Cobertura |
|---|---|
| `rest_repository_type_test.go` | Constantes e struct |
| `rest_repository_config_test.go` | Aplicação de options, erros de validação |
| `rest_repository_error_test.go` | Sentinelas e wrapping |
| `rest_repository_helper_test.go` | `interpolatePath`, `applyQueryParams` |
| `rest_repository_local_test.go` | `resolveURL`, `buildHTTPRequest`, `readBody`, `mapResponse` |
| `rest_repository_test.go` | Table-driven, 6 operações, usando `MockHttpClient` + `MockRestRegistry` |

---

## Arquivos de Referência

| Arquivo | Padrão a seguir |
|---|---|
| `internal/repositories/rest/rest_registry_config.go` | `RegistryOption` (functional options) |
| `internal/repositories/rest/rest_registry_error.go` | Sentinelas + wrap functions |
| `internal/repositories/rest/rest_registry.go` | Construtor com options |
| `pkg/repositories/api.go` | Facade de re-exportação |
| `pkg/types/request_context.go` | `Identifiers`, `Query`, `Headers` |
| `pkg/repositories/rest/contracts/contracts.go` | `RestRepository`, `RestRequest`, `RestResponse` |
| `pkg/repositories/rest/mocks/mocks.go` | Padrão GoMock |

---

## Verificação

```sh
go build ./...
go test ./internal/repositories/rest/... ./pkg/repositories/... ./pkg/http_clients/...
golangci-lint run ./...
```

Cobertura mínima esperada: ≥ 80% em `internal/repositories/rest/`.

Read [](file:///c%3A/Projects/go/go-infra-ports/internal/services/rest/rest_service_local.go)

Read [](file:///c%3A/Projects/go/go-infra-ports/go.mod)

Read [](file:///c%3A/Projects/go/go-infra-ports/pkg/repositories/api.go)

Read [](file:///c%3A/Projects/go/go-infra-ports/internal/repositories/rest/rest_registry.go)

Read [](file:///c%3A/Projects/go/go-infra-ports/internal/repositories/rest/rest_registry_type.go)

Read [](file:///c%3A/Projects/go/go-infra-ports/internal/repositories/rest/rest_registry_config.go)

Read [](file:///c%3A/Projects/go/go-infra-ports/internal/repositories/rest/rest_registry_local.go)

Read [](file:///c%3A/Projects/go/go-infra-ports/internal/repositories/rest/rest_registry_helper.go)

Read [](file:///c%3A/Projects/go/go-infra-ports/internal/repositories/rest/rest_registry_error.go)

Read [](file:///c%3A/Projects/go/go-infra-ports/internal/services/rest/rest_service.go)

Read [](file:///c%3A/Projects/go/go-infra-ports/internal/services/rest/rest_service_type.go)

Read [](file:///c%3A/Projects/go/go-infra-ports/internal/services/rest/rest_service_config.go)

Read [](file:///c%3A/Projects/go/go-infra-ports/pkg/services/rest/api.go)

Read [](file:///c%3A/Projects/go/go-infra-ports/pkg/services/contracts/contracts.go)

Searched for files matching `**/pkg/repositories/**/mocks.go`, 1 match

Read [](file:///c%3A/Projects/go/go-infra-ports/pkg/repositories/rest/mocks/mocks.go#1-1), lines 1 to 200

Read [](file:///c%3A/Projects/go/go-infra-ports/pkg/repositories/rest/mocks/mocks.go#200-200), lines 200 to 422



Fetched 2 resources

Read [](file:///c%3A/Users/bruno/AppData/Roaming/Code/User/workspaceStorage/bd80977c4d6ebc8a85b184bacc12228c/GitHub.copilot-chat/chat-session-resources/3be05195-2550-41d6-b734-bd1a420fd400/toolu_011QpvGQvJSYfBPEm1mnNMgu__vscode-1777807519322/content.txt)

Let me fetch the remaining context I need:   



Read [](file:///c%3A/Projects/go/go-infra-ports/pkg/repositories/api.go#1-1), lines 1 to 60

Fetched https://github.com/brunojet/go-infra-adapters/blob/331dd5c1d355ee872f30e54ede491b16e39eebb5/internal/http_clients/net_http/net_http_client_config.go

Read [](file:///c%3A/Projects/go/go-infra-ports/pkg/types/request_context.go)

Now I have a complete picture. Let me build the plan. 

Read memory

Now I have everything I need to write the plan. 

Created memory file [](file:///c%3A/Users/bruno/AppData/Roaming/Code/User/workspaceStorage/bd80977c4d6ebc8a85b184bacc12228c/GitHub.copilot-chat/memory-tool/memories/M2JlMDUxOTUtMjU1MC00MWQ2LWI3MzQtYmQxYTQyMGZkNDAw/plan.md)

## Avaliação do `http_client`

A implementação em `go-infra-adapters` é sólida:
- Interface mínima `Do(ctx, req) (*http.Response, error)` — nível certo de abstração
- Functional options (`HttpClientOption`) — mesma linguagem do projeto
- Injeção de `RoundTripper` — permite circuit breaker/retries como middleware
- Merge de headers config+request com precedência correta
- Resolução de URL relativa contra `baseURL`
- Sem dependências externas de HTTP

**Decisão arquitetural importante**: `go-infra-adapters` **não** importa `go-infra-ports`. Se `go-infra-ports` importasse `go-infra-adapters`, a direção seria errada (porta dependendo de adaptador). A solução: definir `HttpClient` neste repo — duck typing do Go garante compatibilidade implícita.

---

## Plan: Implement RestRepository (net/http adapter)

Implementar o adaptador concreto `restRepository` em rest que satisfaz a porta `RestRepository` usando `HttpClient` como transporte e `RestRegistry` para marshaling.

**Steps**

### Phase 0 — Contrato HttpClient (novo pacote)
1. Criar `pkg/http_clients/contracts/contracts.go` — interface `HttpClient { Do(ctx, *http.Request) (*http.Response, error) }`
2. Criar `pkg/http_clients/contracts/contracts_test.go` — check de compilação/superfície
3. Criar `pkg/http_clients/mocks/mocks.go` — `MockHttpClient` via GoMock

### Phase 1 — internal: tipos, config, erros, helpers
4. Criar `internal/repositories/rest/rest_repository_type.go` — alias `HttpClient`, constantes `OperationCreate/List/Get/Update/Save/Delete`, struct `restRepository{client, registry, opts}`
5. Criar `internal/repositories/rest/rest_repository_config.go` — `RepositoryOption func(*repositoryOptions) error`, struct `repositoryOptions`, `newRepositoryOptions`, `WithHttpClient`, `WithRegistry`, `WithBaseURL`, `WithPath`, `WithHeader`
6. Criar `internal/repositories/rest/rest_repository_error.go` — erros sentinela (`ErrMissingHttpClient`, `ErrMissingRegistry`, `ErrBuildRequest`, `ErrExecuteRequest`, `ErrReadResponseBody`, `ErrResolveRequest`, `ErrResolveResponse`) + funções de wrap
7. Criar `internal/repositories/rest/rest_repository_helper.go` — `interpolatePath(template, Identifiers) string` (`{key}` → valor), `applyQueryParams(*url.URL, url.Values)`

### Phase 2 — internal: métodos locais e implementação principal *(depends on 4–7)*
8. Criar `internal/repositories/rest/rest_repository_local.go` — métodos privados: `resolveURL`, `buildHTTPRequest`, `executeRequest`, `readBody`, `mapResponse`, `mapResponses`
9. Criar `internal/repositories/rest/rest_repository.go` — `NewRestRepository(opts...) (RestRepository, error)`, 6 métodos (`Create`→POST, `List`→GET collection, `Get`→GET single, `Update`→PUT, `Save`→PATCH, `Delete`→DELETE)

### Phase 3 — API pública *(depends on 9)*
10. Modificar api.go — adicionar `HttpClient`, `RepositoryOption`, `NewRestRepository`, `WithHttpClient`, `WithRegistry`, `WithBaseURL`, `WithPath`, `WithHeader`
11. Atualizar api_test.go

### Phase 4 — arquivos de teste *(parallel, depends on each source file)*
12. `rest_repository_type_test.go`, `rest_repository_config_test.go`, `rest_repository_error_test.go`, `rest_repository_helper_test.go`, `rest_repository_local_test.go`, `rest_repository_test.go` (table-driven, usando `MockHttpClient` + `MockRestRegistry`)

---

**Relevant files**
- rest_registry_config.go — padrão `RegistryOption` a espelhar
- rest_registry_error.go — padrão de erros a espelhar
- rest_registry.go — padrão de construtor a espelhar
- api.go — facade a estender
- request_context.go — `Identifiers`, `Query`, `Headers`
- contracts.go — `RestRepository`, `RestRequest`, `RestResponse`
- mocks.go — padrão de mock a estender

**Verification**
1. `go build go.` — sem erros
2. `go test repositories.` — cobertura ≥ 80%
3. `go test pkg.` e `./pkg/http_clients/...`
4. `golangci-lint run go.` — sem issues

**Decisions**
- Sem dependência nova em go.mod — `HttpClient` definido neste repo
- HTTP method mapping hard-coded por operação (não configurável)
- URL templates usam sintaxe `{key}` resolvida via `RequestContext.Identifiers`
- `HttpClient` e `RestRegistry` são opções **obrigatórias**; demais são opcionais

**Further Considerations**
1. **Semântica de erro em `RepositoryOption`**: O padrão atual do projeto usa `func(*registryOptions)` (sem erro). `HttpClientOption` em go-infra-adapters usa `func(*config) error`. Recomendo: `func(*repositoryOptions) error` para permitir validação nas options (ex.: `WithBaseURL` rejeitar string vazia). Confirma?
2. **`pkg/http_clients/` público**: Apenas `contracts/` e `mocks/` por enquanto (sem api.go), já que não há factory neste repo. Confirma?You've used 95% of your weekly rate limit. Your weekly rate limit will reset on May 3 at 9:00 PM. [Learn More](https://aka.ms/github-copilot-rate-limit-error)