# Implementation Plan

- [x] 1. Remover a validação de extensão de arquivo na interface do usuário
  - Localizar e remover a verificação de extensão .json no arquivo `internal/ui/tui.go`
  - Atualizar comentários relacionados à validação de arquivo
  - _Requirements: 1.1, 1.2_

- [x] 2. Atualizar mensagens de erro relacionadas à validação de arquivo
  - Modificar as mensagens de erro no pacote de localização para focar no conteúdo do arquivo em vez da extensão
  - Atualizar a mensagem de erro para JSON inválido para "O arquivo não contém um JSON válido"
  - Atualizar a mensagem de erro para estrutura inválida para "O arquivo não contém um keystore v3 válido"
  - _Requirements: 2.1, 2.2, 2.3_

- [x] 3. Criar arquivos de teste com diferentes extensões
  - Criar cópias dos arquivos de teste existentes com extensões diferentes (.key, .keystoremaster)
  - Criar cópias dos arquivos de teste existentes sem extensão
  - Criar um arquivo de teste com nome complexo como o exemplo fornecido
  - _Requirements: 4.2, 4.3_

- [ ] 4. Implementar testes unitários para validação de arquivos com diferentes extensões
  - Adicionar testes para validar keystores com extensões não-padrão
  - Adicionar testes para validar keystores sem extensão
  - Adicionar testes para validar keystores com nomes de arquivo complexos
  - _Requirements: 4.1, 4.2, 4.3, 4.4, 4.5_

- [ ] 5. Implementar testes de integração para o fluxo completo de importação
  - Adicionar testes que verificam a importação de keystores com diferentes extensões
  - Verificar se o arquivo é copiado corretamente para o diretório de destino
  - Verificar se a carteira é criada corretamente no banco de dados
  - _Requirements: 3.1, 3.2, 3.3, 3.4_

- [ ] 6. Atualizar a documentação para refletir a nova funcionalidade
  - Atualizar README ou documentação de usuário para mencionar suporte a qualquer extensão de arquivo
  - Adicionar exemplos de uso com diferentes tipos de arquivo
  - _Requirements: 1.1, 1.2_

- [ ] 7. Verificar compatibilidade com keystores existentes
  - Testar a importação de keystores existentes com extensão .json
  - Garantir que o processo de importação continue funcionando como antes para arquivos .json
  - _Requirements: 3.1, 3.2_

- [ ] 8. Implementar validação robusta de conteúdo JSON
  - Reforçar a validação de conteúdo JSON no KeystoreValidator
  - Adicionar verificações adicionais para garantir que o arquivo contém um JSON válido
  - Melhorar as mensagens de erro para fornecer informações mais específicas sobre problemas de formato
  - _Requirements: 1.3, 1.4, 2.1, 2.2_

- [ ] 9. Testar casos de erro e recuperação
  - Implementar testes para verificar o comportamento quando arquivos inválidos são fornecidos
  - Verificar se as mensagens de erro são claras e úteis
  - Testar a recuperação após tentativas de importação com arquivos inválidos
  - _Requirements: 1.4, 2.1, 2.2, 2.3, 4.4, 4.5_

- [ ] 10. Realizar testes de regressão
  - Executar todos os testes existentes para garantir que não houve regressões
  - Verificar se todas as funcionalidades existentes continuam funcionando corretamente
  - _Requirements: 3.1, 3.2, 3.3, 3.4_