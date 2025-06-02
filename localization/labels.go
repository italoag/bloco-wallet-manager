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
			"welcome_message":              "Welcome to the BLOCO wallet Manager!\n\nSelect an option from the menu.",
			"mnemonic_phrase":              "Mnemonic Phrase (Keep it Safe!):",
			"enter_password":               "Enter a password to encrypt the wallet:",
			"press_enter":                  "Press Enter to continue.",
			"import_wallet_title":          "Import Wallet via Mnemonic Phrase",
			"wallet_list_instructions":     "Use the arrow keys to navigate, Enter to view details, ESC to return to the menu.",
			"enter_wallet_password":        "Enter the wallet password:",
			"select_wallet_prompt":         "Select a wallet and enter the password to view the details.",
			"wallet_details_title":         "Wallet Details",
			"ethereum_address":             "Ethereum Address:",
			"public_key":                   "Public Key:",
			"private_key":                  "Private Key:",
			"mnemonic_phrase_label":        "Mnemonic Phrase:",
			"press_esc":                    "Press ESC to return to the wallet list.",
			"main_menu_title":              "Main Menu",
			"create_new_wallet":            "Create New",
			"create_new_wallet_desc":       "Generate a new Ethereum wallet",
			"import_wallet":                "Import Mnemonic",
			"import_wallet_desc":           "Import a wallet using a mnemonic phrase",
			"list_wallets":                 "List Wallets",
			"list_wallets_desc":            "Display all stored wallets",
			"delete_wallet":                "Delete wallet",
			"delete_wallet_desc":           "Delete an existent wallet from database",
			"exit":                         "Exit",
			"exit_desc":                    "Exit the application",
			"error_message":                "Error: %v\n\nPress any key to return to the main menu.",
			"unknown_state":                "Unknown state.",
			"word":                         "Word",
			"password_too_short":           "The password must be at least 8 characters long.",
			"all_words_required":           "All words must be entered.",
			"error_loading_wallets":        "Error loading wallets: %v",
			"password_cannot_be_empty":     "The password cannot be empty.",
			"version":                      "0.1.0",
			"menu":                         "Menu",
			"create_wallet_password":       "Create Wallet Password",
			"import_wallet_password":       "Import Wallet Password",
			"wallet_password":              "Wallet Password",
			"wallet_details":               "Wallet Details",
			"id":                           "ID",
			"operation":                    "Operation",
			"wallet_details_option":        "Wallet details",
			"wallet_delete_option":         "Delete wallet",
			"delete_confirmation_message":  "Are you sure you want to delete your wallet? This will not affect your keystore.",
			"type_delete":                  "Type 'DELETE' to confirm",
			"delete_wallet_confirmation":   "Delete wallet confirmation",
			"wallet_select_operation_view": "Select operation",
		}
	case "pt":
		defaultLabels = map[string]string{
			"welcome_message":              "Bem-vindo ao Administrador de Carteiras BLOCO!\n\nSelecione uma opção do menu.",
			"mnemonic_phrase":              "Frase Mnemotécnica (Mantenha-a Segura!):",
			"enter_password":               "Digite uma senha para encriptar a carteira:",
			"press_enter":                  "Pressione Enter para continuar.",
			"import_wallet_title":          "Importar Carteira via Frase Mnemotécnica",
			"wallet_list_instructions":     "Use as teclas de seta para navegar, Enter para ver detalhes, ESC para voltar ao menu.",
			"enter_wallet_password":        "Digite a senha da carteira:",
			"select_wallet_prompt":         "Selecione uma carteira e digite a senha para ver os detalhes.",
			"wallet_details_title":         "Detalhes da Carteira",
			"ethereum_address":             "Endereço Ethereum:",
			"public_key":                   "Chave Pública:",
			"private_key":                  "Chave Privada:",
			"mnemonic_phrase_label":        "Frase Mnemotécnica:",
			"press_esc":                    "Pressione ESC para voltar à lista de carteiras.",
			"main_menu_title":              "Menu Principal",
			"create_new_wallet":            "Criar Carteira",
			"create_new_wallet_desc":       "Criar uma nova carteira Ethereum",
			"import_wallet":                "Importar via Mnemônico",
			"import_wallet_desc":           "Importar uma carteira usando uma frase mnemônica",
			"list_wallets":                 "Listar Carteiras",
			"list_wallets_desc":            "Exibir todas as carteiras armazenadas",
			"delete_wallet":                "Deletar Carteira",
			"delete_wallet_desc":           "Deleta uma carteira existente da base de dados",
			"exit":                         "Sair",
			"exit_desc":                    "Sair da aplicação",
			"error_message":                "Erro: %v\n\nPressione qualquer tecla para voltar ao menu principal.",
			"unknown_state":                "Estado desconhecido.",
			"word":                         "Palavra",
			"password_too_short":           "A senha deve ter pelo menos 8 caracteres.",
			"all_words_required":           "Todas as palavras devem ser inseridas.",
			"error_loading_wallets":        "Erro ao carregar as carteiras: %v",
			"password_cannot_be_empty":     "A senha não pode estar vazia.",
			"version":                      "0.1.0",
			"id":                           "ID",
			"operation":                    "Operação",
			"wallet_details_option":        "Detalhes da carteira",
			"wallet_delete_option":         "Excluir carteira",
			"delete_confirmation_message":  "Você tem certeza de que deseja excluir sua carteira? Isso não afetará seu keystore.",
			"type_delete":                  "Digite 'DELETE' para confirmar",
			"delete_wallet_confirmation":   "Confirmação de exclusão da carteira",
			"wallet_select_operation_view": "Selecionar operação",
		}
	case "es":
		defaultLabels = map[string]string{
			"welcome_message":              "¡Bienvenido al Administrador de Carteras BLOCO!\n\nSeleccione una opción del menú.",
			"mnemonic_phrase":              "Frase Mnemotécnica (¡Guárdela de Forma Segura!):",
			"enter_password":               "Ingrese una contraseña para encriptar la cartera:",
			"press_enter":                  "Presione Enter para continuar.",
			"import_wallet_title":          "Importar Cartera mediante Frase Mnemotécnica",
			"wallet_list_instructions":     "Use las teclas de flecha para navegar, Enter para ver detalles, ESC para volver al menú.",
			"enter_wallet_password":        "Ingrese la contraseña de la cartera:",
			"select_wallet_prompt":         "Seleccione una cartera e ingrese la contraseña para ver los detalles.",
			"wallet_details_title":         "Detalles de la Cartera",
			"ethereum_address":             "Dirección Ethereum:",
			"public_key":                   "Clave Pública:",
			"private_key":                  "Clave Privada:",
			"mnemonic_phrase_label":        "Frase Mnemotécnica:",
			"press_esc":                    "Presione ESC para volver a la lista de carteras.",
			"main_menu_title":              "Menú Principal",
			"create_new_wallet":            "Crear Nueva Cartera",
			"create_new_wallet_desc":       "Generar una nueva cartera de Ethereum",
			"import_wallet":                "Importar Cartera mediante Mnemónico",
			"import_wallet_desc":           "Importar una cartera existente usando una frase mnemotécnica",
			"list_wallets":                 "Listar Todas las Carteras",
			"list_wallets_desc":            "Mostrar todas las carteras almacenadas",
			"delete_wallet":                "Borrar Cartera",
			"delete_wallet_desc":           "Eliminar una cartera existente de la base de datos",
			"exit":                         "Salir",
			"exit_desc":                    "Salir de la aplicación",
			"error_message":                "Error: %v\n\nPresione cualquier tecla para volver al menú principal.",
			"unknown_state":                "Estado desconocido.",
			"word":                         "Palabra",
			"password_too_short":           "La contraseña debe tener al menos 8 caracteres.",
			"all_words_required":           "Todas las palabras deben ser ingresadas.",
			"error_loading_wallets":        "Error al cargar las carteras: %v",
			"password_cannot_be_empty":     "La contraseña no puede estar vacía.",
			"version":                      "0.1.0",
			"id":                           "ID",
			"operation":                    "Operaciones",
			"wallet_details_option":        "Detalles de la cartera",
			"wallet_delete_option":         "Eliminar carteira",
			"delete_confirmation_message":  "Está seguro de que desea eliminar su cartera? Esto no afectará su keystore.",
			"type_delete":                  "Escriba 'DELETE' para confirmar",
			"delete_wallet_confirmation":   "Confirmación de eliminación de la cartera",
			"wallet_select_operation_view": "Seleccionar operación",
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
