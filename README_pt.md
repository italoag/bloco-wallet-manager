<!--suppress HtmlDeprecatedAttribute -->

<h1 align="center">🔐 BLOCO | Wallet Manager</h1>

<p align="center">
<a href="https://github.com/italoag/bloco-wallet-manager/releases" style="text-decoration: none">
<img src="https://img.shields.io/github/v/release/italoag/bloco-wallet-manager?style=flat-square" alt="Latest Release">
</a>
<a href="https://github.com/italoag/bloco-wallet-manager/stargazers" style="text-decoration: none">
<img src="https://img.shields.io/github/stars/italoag/bloco-wallet-manager.svg?style=flat-square" alt="Stars">
</a>
<a href="https://github.com/italoag/bloco-wallet-manager/fork" style="text-decoration: none">
<img src="https://img.shields.io/github/forks/italoag/bloco-wallet-manager.svg?style=flat-square" alt="Forks">
</a>
<a href="https://opensource.org/licenses/MIT" style="text-decoration: none">
<img src="https://img.shields.io/badge/License-MIT-blue.svg?style=flat-square" alt="License: MIT">
</a>
<br/>
<a href="https://github.com/italoag/bloco-wallet-manager/releases" style="text-decoration: none">
<img src="https://img.shields.io/badge/platform-windows%20%7C%20macos%20%7C%20linux-informational?style=for-the-badge" alt="Downloads">
</a>
 <a href="https://twitter.com/0xitalo">
        <img src="https://img.shields.io/badge/Twitter-%400xitalo-1DA1F2?logo=twitter&style=for-the-badge" alt=""/>
    </a>
<br/>

### Índice
- [Introdução](#introdução)
- [Recursos](#recursos)
- [Instalação](#instalação)
- [Uso](#uso)
- [Roteiro](#roteiro)
- [Contribuindo](#contribuindo)
- [Licença](#licença)

### Introdução

**BLOCO Wallet Manager** é um aplicativo de Interface de Linha de Comando (CLI) projetado para gerenciar carteiras de criptomoedas compatíveis com a rede Ethereum, aderindo ao padrão KeyStoreV3. Desenvolvido em GoLang, o BLOCO oferece uma Interface de Usuário de Terminal (TUI) para gerenciamento de carteiras de forma eficiente. Futuras integrações incluirão cofres externos como Hashicorp Vault, Amazon KMS, Cloud HSM e Azure Key Vault.

### Recursos

- **Gerenciamento de Carteiras**
    - Criar novas carteiras compatíveis com Ethereum.
    - Importar carteiras usando Mnemônicos ou Chaves Privadas.
    - Exportar carteiras no formato KeyStoreV3.
    - Excluir, bloquear e desbloquear endereços de carteiras.
    - Listar todas as carteiras gerenciadas.

- **Segurança**
    - Compatibilidade com KeyStoreV3 para armazenamento seguro de chaves.
    - Integrações planejadas com cofres externos:
        - Hashicorp Vault
        - Amazon KMS
        - Cloud HSM
        - Azure Key Vault

- **Consulta de Saldo**
    - Consultar o saldo de carteiras compatíveis com Ethereum.

- **Extensibilidade**
    - Suporte para redes blockchain adicionais (futuro).
    - Suporte para múltiplas curvas criptográficas e algoritmos de assinatura:
        - Curvas: secp256k1, secp256r1, ed25519
        - Algoritmos: ECDSA, EdDSA

### Instalação

Certifique-se de ter o [Go](https://golang.org/doc/install) instalado no seu sistema.

```bash
git clone https://github.com/italoag/bloco-wallet-manager.git
cd bloco-wallet-manager
go build -o bwm
```

Mova o executável para um diretório no seu PATH para fácil acesso:

```bash
mv bwm /usr/local/bin/
```

### Uso

Execute o BLOCO Wallet Manager usando o terminal:

```bash
bwm
```

Navegue pela TUI para gerenciar suas carteiras. Os comandos disponíveis incluem:

- **Criar Carteira:** Inicializar uma nova carteira compatível com Ethereum.
- **Importar Carteira:** Importar carteiras existentes usando Mnemônicos ou Chaves Privadas.
- **Listar Carteiras:** Exibir todas as carteiras gerenciadas.
- **Excluir Carteira:** Remover uma carteira do gerenciador.
- **Bloquear/Desbloquear Carteira:** Desativar ou ativar temporariamente um endereço de carteira.
- **Verificar Saldo:** Visualizar o saldo de uma carteira selecionada.

### Roteiro

**Funcionalidades Futuras:**

- **Integrações com Cofres:**
    - Hashicorp Vault
    - Amazon KMS
    - Cloud HSM
    - Azure Key Vault

- **Suporte Multi-Rede:**
    - Integração com redes blockchain adicionais.

- **Criptografia Avançada:**
    - Suporte para as curvas secp256r1 e ed25519.
    - Implementação dos algoritmos de assinatura ECDSA e EdDSA.

- **Recursos de Segurança Aprimorados:**
    - Autenticação de dois fatores para acesso à carteira.
    - Suporte para carteiras multiassinatura.

- **Melhorias na Experiência do Usuário:**
    - TUI aprimorada com navegação mais intuitiva.
    - Históricos detalhados de transações e análises.

### Contribuindo

Contribuições são bem-vindas! Siga estes passos:

1. Faça um fork do repositório.
2. Crie uma nova branch para sua funcionalidade ou correção de bug.
3. Faça commits com mensagens claras.
4. Envie um pull request detalhando suas alterações.

Para mais informações, consulte o arquivo [CONTRIBUTING.md](CONTRIBUTING.md).

### Licença

Este projeto está licenciado sob a [Licença MIT](LICENSE).