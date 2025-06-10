# Diretrizes de Debugging - Abordagens que NÃO Funcionaram

## ❌ Problemas com Integração huh + bubbletea

### 1. **Early Return no TUI (FALHOU)**
**Tentativa**: Fazer early return quando componentes retornam comandos
```go
// ❌ FALHOU - Impede processamento completo de mensagens
if cmd != nil {
    return m, cmd
}
```
**Problema**: Impede que o form processe todas as mensagens necessárias para navegação e validação.

### 2. **Interceptação de Escape Antes do Form (FALHOU)**
**Tentativa**: Processar escape key antes do form.Update()
```go
// ❌ FALHOU - Quebra navegação interna do huh
if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" {
    return c, func() tea.Msg { return BackToMenuMsg{} }
}
form, cmd := c.form.Update(msg)
```
**Problema**: O huh precisa processar escape para navegação interna entre campos.

### 3. **Múltiplos Grupos Separados (FALHOU PARCIALMENTE)**
**Tentativa**: Criar um grupo separado para cada campo
```go
// ❌ FALHOU - Causa problemas de navegação
huh.NewForm(
    huh.NewGroup(huh.NewInput().Title("Nome")),
    huh.NewGroup(huh.NewInput().Title("Senha")),
    huh.NewGroup(huh.NewInput().Title("Chave")),
)
```
**Problema**: Campos relacionados devem estar no mesmo grupo para navegação fluida.

### 4. **SuggestionsFunc com API Síncrona (FALHOU)**
**Tentativa**: Chamar API blockchain diretamente na SuggestionsFunc
```go
// ❌ FALHOU - Causa travamentos
SuggestionsFunc(func() []string {
    suggestions, _ := c.chainListService.SearchNetworksByName(c.networkName)
    // Processamento síncrono que trava a UI
    return names
}, &c.networkName)
```
**Problema**: SuggestionsFunc é chamada frequentemente e não deve fazer operações lentas.

## ✅ Soluções que Funcionaram

### 1. **Form Processa Primeiro, Escape Depois**
```go
// ✅ FUNCIONA - Permite navegação interna do huh
form, cmd := c.form.Update(msg)
if f, ok := form.(*huh.Form); ok {
    c.form = f
    cmds = append(cmds, cmd)
}

// Só processa escape se form não está em uso ativo
if keyMsg, ok := msg.(tea.KeyMsg); ok && keyMsg.String() == "esc" && c.form.State == huh.StateNormal {
    return c, func() tea.Msg { return BackToMenuMsg{} }
}
```

### 2. **Campos Relacionados no Mesmo Grupo**
```go
// ✅ FUNCIONA - Navegação fluida
huh.NewForm(
    huh.NewGroup(
        huh.NewInput().Title("Nome").Value(&c.name),
        huh.NewInput().Title("Senha").Value(&c.password),
        huh.NewInput().Title("Chave").Value(&c.key),
    ),
)
```

### 3. **Sugestões Estáticas ou Pré-carregadas**
```go
// ✅ FUNCIONA - Performance consistente
Suggestions([]string{
    "Polygon", "Binance Smart Chain", "Ethereum",
    // Lista pré-definida de redes populares
})
```

## 🔄 Integração huh + bubbletea - Padrão Correto

### Model Structure
```go
type Component struct {
    form *huh.Form
    // outros campos
}

func (c *Component) Init() tea.Cmd {
    return c.form.Init() // ESSENCIAL para foco inicial
}

func (c *Component) Update(msg tea.Msg) (*Component, tea.Cmd) {
    // 1. Processar mensagens específicas primeiro (WindowSize, custom msgs)
    // 2. Sempre deixar form.Update() processar TODAS as mensagens
    // 3. Só interceptar keys depois se necessário e com condições específicas
    // 4. Verificar form.State para ações baseadas em estado
}
```

### TUI Integration
```go
// ✅ PADRÃO CORRETO - Sempre chamar Init() ao entrar na view
case 1: // Create New Wallet
    m.currentView = CreateWalletView
    m.createWalletComponent.Reset()
    return m, m.createWalletComponent.Init() // ESSENCIAL

// ✅ PADRÃO CORRETO - Não fazer early return em form views
case CreateWalletView:
    updatedComponent, componentCmd := m.createWalletComponent.Update(msg)
    m.createWalletComponent = *updatedComponent
    cmd = componentCmd
    // Continue processando outras mensagens - NÃO fazer early return
```

## 🚫 Anti-Patterns a Evitar

1. **Nunca** interceptar keys antes do form processar
2. **Nunca** fazer early return baseado em comandos em form views
3. **Nunca** chamar APIs síncronas em SuggestionsFunc
4. **Nunca** esquecer de chamar form.Init() ao entrar na view
5. **Nunca** criar grupos desnecessariamente separados para campos relacionados

## ✅ Best Practices Confirmadas

1. **Sempre** chamar form.Init() ao entrar em form views
2. **Sempre** deixar form.Update() processar primeiro
3. **Sempre** verificar form.State antes de interceptar keys
4. **Sempre** agrupar campos relacionados no mesmo huh.Group
5. **Sempre** usar sugestões estáticas ou pré-carregadas para performance