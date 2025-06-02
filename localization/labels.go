package localization

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

var Labels map[string]string

func SetLanguage(lang string, appDir string) error {
	labelsPath := filepath.Join(appDir, "locales", fmt.Sprintf("%s.yaml", lang))
	fmt.Printf("\n\n\n%s", labelsPath)
	if _, err := os.Stat(labelsPath); os.IsNotExist(err) {
		err := createDefaultLabels(lang, labelsPath)
		if err != nil {
			return err
		}
	}

	labelsFile, err := os.ReadFile(labelsPath)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(labelsFile, &Labels)
	if err != nil {
		return err
	}

	return nil
}

func createDefaultLabels(lang, labelsPath string) error {
	defaultLabels := map[string]string{}
	switch lang {
	case "en":
		defaultLabels = map[string]string{

			"welcome_message":          "Welcome to the BLOCO wallet Manager!\n\nSelect an option from the menu.",
			"mnemonic_phrase":          "Mnemonic Phrase (Keep it Safe!):",
			"enter_password":           "Enter a password to encrypt the wallet:",
			"press_enter":              "Press Enter to continue.",
			"import_wallet_title":      "Import an existing Wallet",
			"wallet_list_instructions": "Use the arrow keys to navigate, Enter to view details, 'd' to delete a wallet, 'esc' to return to the menu.",
			"status_bar_instructions":  "View: %s | Press 'esc' to return | Press 'q' to quit",
			"wallet_list_status_bar":   "View: %s | Press 'd' to delete | Press 'esc' to return | Press 'q' to quit",
			"enter_wallet_password":    "Enter the wallet password:",
			"select_wallet_prompt":     "Select a wallet and enter the password to view the details.",
			"wallet_details_title":     "Wallet Details",
			"ethereum_address":         "Ethereum Address:",
			"public_key":               "Public Key:",
			"private_key":              "Private Key:",
			"mnemonic_phrase_label":    "Mnemonic Phrase:",
			"press_esc":                "Press ESC to return to the wallet list.",
			"main_menu_title":          "Main Menu",
			"create_new_wallet":        "Create New",
			"create_new_wallet_desc":   "Generate a new Ethereum wallet",
			"import_wallet":            "Import Wallet",
			"import_wallet_desc":       "Import an existing wallet",
			"import_method_title":      "Select Import Method",
			"import_mnemonic":          "Mnemonic Phrase",
			"import_mnemonic_desc":     "Import using 12-word mnemonic phrase",
			"import_private_key":       "Private Key",
			"import_private_key_desc":  "Import using a private key",
			"back_to_menu":             "Back to Main Menu",
			"back_to_menu_desc":        "Return to the main menu",
			"private_key_title":        "Import Wallet via Private Key",
			"enter_private_key":        "Enter the private key (with or without 0x prefix):",
			"invalid_private_key":      "Invalid private key format",
			"list_wallets":             "List Wallets",
			"list_wallets_desc":        "Display all stored wallets",
			"exit":                     "Exit",
			"exit_desc":                "Exit the application",
			"error_message":            "Error: %v\n\nPress any key to return to the main menu.",
			"unknown_state":            "Unknown state.",
			"word":                     "Word",
			"password_too_short":       "The password must be at least 8 characters long.",
			"all_words_required":       "All words must be entered.",
			"error_loading_wallets":    "Error loading wallets: %v",
			"password_cannot_be_empty": "The password cannot be empty.",
			"version":                  "0.2.0",
			"menu":                     "Menu",
			"create_wallet_password":   "Create Wallet Password",
			"import_wallet_password":   "Import Wallet Password",
			"import_method_selection":  "Import Method Selection",
			"import_private_key_view":  "Import Private Key",
			"wallet_password":          "Wallet Password",
			"wallet_details":           "Wallet Details",
			"id":                       "ID",
			"confirm_delete_wallet":    "Are you sure you want to delete this wallet?",
			"confirm":                  "Confirm",
			"cancel":                   "Cancel",
		}
	case "pt":
		defaultLabels = map[string]string{
			"welcome_message":           "Bem-vindo ao Administrador de Carteiras BLOCO!\n\nSelecione uma opção do menu.",
			"mnemonic_phrase":           "Frase Mnemotécnica (Mantenha-a Segura!):",
			"enter_password":            "Digite uma senha para encriptar a carteira:",
			"press_enter":               "Pressione Enter para continuar.",
			"import_wallet_title":       "Importar carteira pré existente",
			"wallet_list_instructions":  "Use as teclas de seta para navegar, Enter para ver detalhes, ESC para voltar ao menu.",
			"status_bar_instructions":   "Visualização: %s | Pressione 'esc' ou 'backspace' para retornar | Pressione 'q' para sair",
			"wallet_list_status_bar":    "Visualização: %s | Pressione 'd' para excluir | Pressione 'esc' para retornar | Pressione 'q' para sair",
			"enter_wallet_password":     "Digite a senha da carteira:",
			"select_wallet_prompt":      "Selecione uma carteira e digite a senha para ver os detalhes.",
			"wallet_details_title":      "Detalhes da Carteira",
			"ethereum_address":          "Endereço Ethereum:",
			"public_key":                "Chave Pública:",
			"private_key":               "Chave Privada:",
			"mnemonic_phrase_label":     "Frase Mnemotécnica:",
			"press_esc":                 "Pressione ESC para voltar à lista de carteiras.",
			"main_menu_title":           "Menu Principal",
			"create_new_wallet":         "Criar Carteira",
			"create_new_wallet_desc":    "Criar uma nova carteira Ethereum",
			"import_wallet":             "Importar Carteira",
			"import_wallet_desc":        "Importar uma carteira existente",
			"import_method_title":       "Selecione o Método de Importação",
			"import_mnemonic":           "Frase Mnemônica",
			"import_mnemonic_desc":      "Importar usando frase mnemônica de 12 palavras",
			"import_private_key":        "Chave Privada",
			"import_private_key_desc":   "Importar usando uma chave privada",
			"back_to_menu":              "Voltar ao Menu Principal",
			"back_to_menu_desc":         "Retornar ao menu principal",
			"private_key_title":         "Importar Carteira via Chave Privada",
			"enter_private_key":         "Digite a chave privada (com ou sem prefixo 0x):",
			"invalid_private_key":       "Formato de chave privada inválido",
			"list_wallets":              "Listar Carteiras",
			"list_wallets_desc":         "Exibir todas as carteiras armazenadas",
			"exit":                      "Sair",
			"exit_desc":                 "Sair da aplicação",
			"error_message":             "Erro: %v\n\nPressione qualquer tecla para voltar ao menu principal.",
			"unknown_state":             "Estado desconhecido.",
			"word":                      "Palavra",
			"password_too_short":        "A senha deve ter pelo menos 8 caracteres.",
			"all_words_required":        "Todas as palavras devem ser inseridas.",
			"error_loading_wallets":     "Erro ao carregar as carteiras: %v",
			"password_cannot_be_empty":  "A senha não pode estar vazia.",
			"version":                   "0.1.0",
			"id":                        "ID",
			"confirm_delete_wallet":     "Tem certeza de que deseja excluir esta carteira?",
			"confirm":                   "Confirmar",
			"cancel":                    "Cancelar",
			"list_wallets_title":        "Lista de Carteiras",
			"list_wallets_instructions": "Use as setas ↑↓ para navegar, Enter para selecionar, 'd' ou 'delete' para excluir uma carteira, ESC para voltar ao menu.",
		}
	case "es":
		defaultLabels = map[string]string{
			"welcome_message":          "¡Bienvenido al Administrador de Carteras BLOCO!\n\nSeleccione una opción del menú.",
			"mnemonic_phrase":          "Frase Mnemotécnica (¡Guárdela de Forma Segura!):",
			"enter_password":           "Ingrese una contraseña para encriptar la cartera:",
			"press_enter":              "Presione Enter para continuar.",
			"import_wallet_title":      "Importar Cartera mediante Frase Mnemotécnica",
			"wallet_list_instructions": "Use las teclas de flecha para navegar, Enter para ver detalles, 'd' o 'delete' para eliminar una cartera, ESC para volver al menú.",
			"status_bar_instructions":  "Vista: %s | Presione 'esc' o 'backspace' para regresar | Presione 'q' para salir",
			"wallet_list_status_bar":   "Vista: %s | Presione 'd' para eliminar | Presione 'esc' para regresar | Presione 'q' para salir",
			"enter_wallet_password":    "Ingrese la contraseña de la cartera:",
			"select_wallet_prompt":     "Seleccione una cartera e ingrese la contraseña para ver los detalles.",
			"wallet_details_title":     "Detalles de la Cartera",
			"ethereum_address":         "Dirección Ethereum:",
			"public_key":               "Clave Pública:",
			"private_key":              "Clave Privada:",
			"mnemonic_phrase_label":    "Frase Mnemotécnica:",
			"press_esc":                "Presione ESC para volver a la lista de carteras.",
			"main_menu_title":          "Menú Principal",
			"create_new_wallet":        "Crear Nueva Cartera",
			"create_new_wallet_desc":   "Generar una nueva cartera de Ethereum",
			"import_wallet":            "Importar Cartera",
			"import_wallet_desc":       "Importar una cartera existente",
			"import_method_title":      "Seleccione el Método de Importación",
			"import_mnemonic":          "Frase Mnemotécnica",
			"import_mnemonic_desc":     "Importar usando frase mnemotécnica de 12 palabras",
			"import_private_key":       "Clave Privada",
			"import_private_key_desc":  "Importar usando una clave privada",
			"back_to_menu":             "Volver al Menú Principal",
			"back_to_menu_desc":        "Regresar al menú principal",
			"private_key_title":        "Importar Cartera mediante Clave Privada",
			"enter_private_key":        "Ingrese la clave privada (con o sin prefijo 0x):",
			"invalid_private_key":      "Formato de clave privada inválido",
			"list_wallets":             "Listar Todas las Carteras",
			"list_wallets_desc":        "Mostrar todas las carteras almacenadas",
			"exit":                     "Salir",
			"exit_desc":                "Salir de la aplicación",
			"error_message":            "Error: %v\n\nPresione cualquier tecla para volver al menú principal.",
			"unknown_state":            "Estado desconocido.",
			"word":                     "Palabra",
			"password_too_short":       "La contraseña debe tener al menos 8 caracteres.",
			"all_words_required":       "Todas las palabras deben ser ingresadas.",
			"error_loading_wallets":    "Error al cargar las carteras: %v",
			"password_cannot_be_empty": "La contraseña no puede estar vacía.",
			"version":                  "0.1.0",
			"id":                       "ID",
			"confirm_delete_wallet":    "¿Está seguro de que desea eliminar esta cartera?",
			"confirm":                  "Confirmar",
			"cancel":                   "Cancelar",
		}
	default:
		return fmt.Errorf("unsupported language: %s", lang)
	}

	data, err := yaml.Marshal(defaultLabels)
	if err != nil {
		return err
	}

	appDir := filepath.Dir(labelsPath)
	localesDir := filepath.Join(appDir, "locales")
	if _, err := os.Stat(localesDir); os.IsNotExist(err) {
		err := os.MkdirAll(localesDir, os.ModePerm)
		if err != nil {
			return err
		}
	}

	err = os.WriteFile(labelsPath, data, 0644)
	if err != nil {
		return err
	}

	return nil
}
