package main

import (
	"crypto/ecdsa"
	"database/sql"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/table"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ethereum/go-ethereum/accounts/keystore"
	"github.com/ethereum/go-ethereum/crypto"
	_ "github.com/mattn/go-sqlite3"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

const (
	WALLETS_DIR = ".wallets/keystore"
	DB_NAME     = ".wallets/wallets.db"
	MENU_WIDTH  = 60 // Fixed menu width to prevent excessive space usage
)

type Wallet struct {
	ID           int
	Address      string
	KeyStorePath string
	Mnemonic     string
}

type Database struct {
	conn *sql.DB
}

func NewDatabase(dbPath string) (*Database, error) {
	conn, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS wallets (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		address TEXT UNIQUE NOT NULL,
		keystore_path TEXT NOT NULL,
		mnemonic TEXT NOT NULL
	);
	`
	_, err = conn.Exec(createTableQuery)
	if err != nil {
		return nil, err
	}

	return &Database{conn: conn}, nil
}

func (db *Database) AddWallet(wallet Wallet) error {
	insertQuery := `
	INSERT INTO wallets (address, keystore_path, mnemonic)
	VALUES (?, ?, ?);
	`
	_, err := db.conn.Exec(insertQuery, wallet.Address, wallet.KeyStorePath, wallet.Mnemonic)
	return err
}

func (db *Database) GetAllWallets() ([]Wallet, error) {
	selectQuery := `
	SELECT id, address, keystore_path, mnemonic FROM wallets;
	`
	rows, err := db.conn.Query(selectQuery)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var wallets []Wallet
	for rows.Next() {
		var w Wallet
		err := rows.Scan(&w.ID, &w.Address, &w.KeyStorePath, &w.Mnemonic)
		if err != nil {
			return nil, err
		}
		wallets = append(wallets, w)
	}

	return wallets, nil
}

func (db *Database) Close() error {
	return db.conn.Close()
}

type menuItem struct {
	title       string
	description string
}

func (i menuItem) Title() string       { return i.title }
func (i menuItem) Description() string { return i.description }
func (i menuItem) FilterValue() string { return i.title }

type Model struct {
	db                 *Database
	keystore           *keystore.KeyStore
	currentView        string
	menuList           list.Model
	importWords        []string
	importStage        int
	textInputs         []textinput.Model
	wallets            []Wallet
	selectedWallet     *Wallet
	err                error
	passwordInput      textinput.Model
	mnemonic           string
	walletTable        table.Model
	selectedPrivateKey *ecdsa.PrivateKey
	width              int
	height             int
}

var (
	contentStyle = lipgloss.NewStyle().Padding(1).Align(lipgloss.Left)
)

func InitModel(db *Database, ks *keystore.KeyStore) Model {
	menuItems := []list.Item{
		menuItem{title: "Create New Wallet", description: "Generate a new Ethereum wallet"},
		menuItem{title: "Import Wallet via Mnemonic", description: "Import an existing wallet using a mnemonic phrase"},
		menuItem{title: "List All Wallets", description: "Display all stored wallets"},
		menuItem{title: "Exit", description: "Exit the application"},
	}

	delegate := list.NewDefaultDelegate()
	delegate.ShowDescription = true
	delegate.Styles.SelectedTitle = delegate.Styles.SelectedTitle.Foreground(lipgloss.Color("#00FF00")).Bold(true)
	delegate.Styles.SelectedDesc = delegate.Styles.SelectedDesc.Foreground(lipgloss.Color("#00FF00"))
	delegate.Styles.NormalTitle = delegate.Styles.NormalTitle.Foreground(lipgloss.Color("#FFFFFF"))
	delegate.Styles.NormalDesc = delegate.Styles.NormalDesc.Foreground(lipgloss.Color("#888888"))

	menuList := list.New(menuItems, delegate, MENU_WIDTH, 0)
	menuList.Title = "Main Menu"
	menuList.SetShowStatusBar(false)
	menuList.SetFilteringEnabled(false)
	menuList.SetShowHelp(false)
	menuList.Styles.Title = lipgloss.NewStyle().
		Background(lipgloss.Color("#25A065")).
		Foreground(lipgloss.Color("#FFFDF5")).
		Padding(0, 1).
		Bold(true)

	return Model{
		db:          db,
		keystore:    ks,
		currentView: "menu",
		menuList:    menuList,
	}
}

func (m *Model) Init() tea.Cmd {
	return nil
}

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.err != nil {
		switch msg.(type) {
		case tea.KeyMsg:
			m.err = nil
			m.currentView = "menu"
			return m, nil
		}
		return m, nil
	}

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height

		// Set menu height
		m.menuList.SetHeight(m.height - 2)

		// Calculate content width
		contentWidth := m.width - MENU_WIDTH - 2
		if contentWidth < 20 {
			contentWidth = 20 // Minimum content width
		}

		// Adjust table size if it exists
		if m.walletTable.Columns() != nil {
			m.walletTable.SetWidth(contentWidth - 4)
			m.walletTable.SetHeight(m.height - 4)
		}
		return m, nil
	}

	switch m.currentView {
	case "menu":
		var cmd tea.Cmd
		m.menuList, cmd = m.menuList.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				i, ok := m.menuList.SelectedItem().(menuItem)
				if ok {
					switch i.title {
					case "Create New Wallet":
						m.currentView = "create_wallet"
						m.initCreateWallet()
					case "Import Wallet via Mnemonic":
						m.currentView = "import_wallet"
						m.initImportWallet()
					case "List All Wallets":
						m.currentView = "list_wallets"
						m.initListWallets()
					case "Exit":
						return m, tea.Quit
					}
				}
			}
		}
		return m, cmd
	case "create_wallet_password":
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "enter" {
				password := strings.TrimSpace(m.passwordInput.Value())
				if len(password) < 8 {
					m.err = fmt.Errorf("the password must be at least 8 characters long")
					m.currentView = "menu"
					return m, nil
				}
				err := m.createWallet(password)
				if err != nil {
					m.err = err
					m.currentView = "menu"
					return m, nil
				}
				// Alteração: Exibir os detalhes da wallet após criação
				// A função createWallet já define m.selectedWallet e m.currentView
			} else {
				var cmd tea.Cmd
				m.passwordInput, cmd = m.passwordInput.Update(msg)
				return m, cmd
			}
		}
	case "import_wallet":
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "enter":
				word := strings.TrimSpace(m.textInputs[m.importStage].Value())
				if word == "" {
					m.err = fmt.Errorf("all words must be entered")
					return m, nil
				}
				m.importWords[m.importStage] = word
				m.textInputs[m.importStage].Blur()
				m.importStage++
				if m.importStage < 12 {
					m.textInputs[m.importStage].Focus()
				} else {
					m.passwordInput = textinput.New()
					m.passwordInput.Placeholder = "Enter a password to encrypt the wallet"
					m.passwordInput.CharLimit = 64
					m.passwordInput.Width = 30
					m.passwordInput.EchoMode = textinput.EchoPassword
					m.passwordInput.EchoCharacter = '•'
					m.passwordInput.Focus()
					m.currentView = "import_wallet_password"
				}
			default:
				var cmd tea.Cmd
				m.textInputs[m.importStage], cmd = m.textInputs[m.importStage].Update(msg)
				return m, cmd
			}
		}
	case "import_wallet_password":
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "enter" {
				password := strings.TrimSpace(m.passwordInput.Value())
				if len(password) < 8 {
					m.err = fmt.Errorf("the password must be at least 8 characters long")
					m.currentView = "menu"
					return m, nil
				}
				err := m.importWallet(password)
				if err != nil {
					m.err = err
					m.currentView = "menu"
					return m, nil
				}
				m.currentView = "menu"
			} else {
				var cmd tea.Cmd
				m.passwordInput, cmd = m.passwordInput.Update(msg)
				return m, cmd
			}
		}
	case "list_wallets":
		var cmd tea.Cmd
		m.walletTable, cmd = m.walletTable.Update(msg)
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "up", "k":
				m.walletTable.MoveUp(0) // Alterado de MoveUp(0) para MoveUp(1)
			case "down", "j":
				m.walletTable.MoveDown(0) // Alterado de MoveDown(0) para MoveDown(1)
			case "enter":
				selectedRow := m.walletTable.SelectedRow()
				if len(selectedRow) > 1 {
					for _, w := range m.wallets {
						if w.Address == selectedRow[1] {
							m.selectedWallet = &w
							break
						}
					}
					m.currentView = "wallet_password"
					m.initWalletPassword()
				}
			case "esc":
				m.currentView = "menu"
			}
		}
		return m, cmd
	case "wallet_password":
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "enter" {
				password := strings.TrimSpace(m.passwordInput.Value())
				if password == "" {
					m.err = fmt.Errorf("the password cannot be empty")
					m.currentView = "menu"
					return m, nil
				}
				err := m.loadWallet(password)
				if err != nil {
					m.err = err
					m.currentView = "menu"
					return m, nil
				}
				m.currentView = "wallet_details"
			} else {
				var cmd tea.Cmd
				m.passwordInput, cmd = m.passwordInput.Update(msg)
				return m, cmd
			}
		}
	case "wallet_details":
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "esc" {
				m.currentView = "list_wallets"
			}
		}
	}
	return m, nil
}

func (m *Model) initCreateWallet() {
	mnemonic, err := gerarMnemonico()
	if err != nil {
		m.err = err
		m.currentView = "menu"
		return
	}
	m.mnemonic = mnemonic
	m.passwordInput = textinput.New()
	m.passwordInput.Placeholder = "Enter a password to encrypt the wallet"
	m.passwordInput.CharLimit = 64
	m.passwordInput.Width = 30
	m.passwordInput.EchoMode = textinput.EchoPassword
	m.passwordInput.EchoCharacter = '•'
	m.passwordInput.Focus()
	m.currentView = "create_wallet_password"
}

func (m *Model) createWallet(password string) error {
	privateKeyHex, err := derivarChavePrivada(m.mnemonic)
	if err != nil {
		return err
	}
	privKey, err := hexToECDSA(privateKeyHex)
	if err != nil {
		return err
	}
	account, err := m.keystore.ImportECDSA(privKey, password)
	if err != nil {
		return err
	}
	originalPath := account.URL.Path
	newFilename := fmt.Sprintf("%s.json", account.Address.Hex())
	newPath := filepath.Join(filepath.Dir(originalPath), newFilename)
	err = os.Rename(originalPath, newPath)
	if err != nil {
		return fmt.Errorf("error renaming the wallet file: %v", err)
	}
	wallet := Wallet{
		Address:      account.Address.Hex(),
		KeyStorePath: newPath,
		Mnemonic:     m.mnemonic,
	}
	err = m.db.AddWallet(wallet)
	if err != nil {
		return err
	}

	// Alteração: Definir a wallet recém-criada como selecionada e exibir os detalhes
	m.selectedWallet = &wallet
	m.currentView = "wallet_details"

	return nil
}

func (m *Model) initImportWallet() {
	m.importWords = make([]string, 12)
	m.importStage = 0
	m.textInputs = make([]textinput.Model, 12)
	for i := 0; i < 12; i++ {
		ti := textinput.New()
		ti.Placeholder = fmt.Sprintf("Word %d", i+1)
		ti.CharLimit = 32
		ti.Width = 20
		m.textInputs[i] = ti
	}
	m.textInputs[0].Focus()
}

func (m *Model) importWallet(password string) error {
	frase := strings.Join(m.importWords, " ")
	if !bip39.IsMnemonicValid(frase) {
		return fmt.Errorf("invalid mnemonic phrase")
	}
	m.selectedWallet = &Wallet{
		Mnemonic: frase,
	}
	privateKeyHex, err := derivarChavePrivada(m.selectedWallet.Mnemonic)
	if err != nil {
		return err
	}
	privKey, err := hexToECDSA(privateKeyHex)
	if err != nil {
		return err
	}
	account, err := m.keystore.ImportECDSA(privKey, password)
	if err != nil {
		return err
	}
	originalPath := account.URL.Path
	newFilename := fmt.Sprintf("%s.json", account.Address.Hex())
	newPath := filepath.Join(filepath.Dir(originalPath), newFilename)
	err = os.Rename(originalPath, newPath)
	if err != nil {
		return fmt.Errorf("error renaming the wallet file: %v", err)
	}
	wallet := Wallet{
		Address:      account.Address.Hex(),
		KeyStorePath: newPath,
		Mnemonic:     m.selectedWallet.Mnemonic,
	}
	err = m.db.AddWallet(wallet)
	if err != nil {
		return err
	}
	return nil
}

func (m *Model) initListWallets() {
	wallets, err := m.db.GetAllWallets()
	if err != nil {
		m.err = fmt.Errorf("error loading wallets: %v", err)
		m.currentView = "menu"
		return
	}
	m.wallets = wallets
	columns := []table.Column{
		{Title: "ID", Width: 5},
		{Title: "Ethereum Address", Width: 61 - 1}, // Adjusted width based on fixed menu width
	}
	var rows []table.Row
	for _, w := range m.wallets {
		rows = append(rows, table.Row{fmt.Sprintf("%d", w.ID), w.Address})
	}
	m.walletTable = table.New(
		table.WithColumns(columns),
		table.WithRows(rows),
		table.WithFocused(true),
	)
	m.walletTable.SetWidth(MENU_WIDTH - 2) // Fixed menu width
	m.walletTable.SetHeight(m.height - 4)

	s := table.DefaultStyles()
	s.Header = s.Header.
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("240")).
		BorderBottom(true).
		Bold(false)
	s.Selected = s.Selected.
		Foreground(lipgloss.Color("229")).
		Background(lipgloss.Color("57")).
		Bold(false)
	s.Cell = s.Cell.Align(lipgloss.Left)
	m.walletTable.SetStyles(s)
}

func (m *Model) initWalletPassword() {
	m.passwordInput = textinput.New()
	m.passwordInput.Placeholder = "Enter the wallet password"
	m.passwordInput.CharLimit = 64
	m.passwordInput.Width = 30
	m.passwordInput.EchoMode = textinput.EchoPassword
	m.passwordInput.EchoCharacter = '•'
	m.passwordInput.Focus()
}

func (m *Model) loadWallet(password string) error {
	keyJSON, err := os.ReadFile(m.selectedWallet.KeyStorePath)
	if err != nil {
		return fmt.Errorf("error reading the wallet file: %v", err)
	}
	key, err := keystore.DecryptKey(keyJSON, password)
	if err != nil {
		return fmt.Errorf("incorrect password")
	}
	m.selectedPrivateKey = key.PrivateKey
	return nil
}

func (m Model) View() string {
	if m.err != nil {
		return fmt.Sprintf("Error: %v\n\nPress any key to return to the main menu.", m.err)
	}

	var menuView string
	var contentView string

	menuStyle := lipgloss.NewStyle().Padding(1).Align(lipgloss.Left).Width(MENU_WIDTH)

	switch m.currentView {
	case "menu":
		menuView = m.menuList.View()
		contentView = m.viewWelcome()
	case "create_wallet_password":
		menuView = m.menuList.View()
		contentView = m.viewCreateWalletPassword()
	case "import_wallet":
		menuView = m.menuList.View()
		contentView = m.viewImportWallet()
	case "import_wallet_password":
		menuView = m.menuList.View()
		contentView = m.viewImportWalletPassword()
	case "list_wallets":
		menuView = m.menuList.View()
		contentView = m.viewWalletList()
	case "wallet_password":
		menuView = m.menuList.View()
		contentView = m.viewWalletPassword()
	case "wallet_details":
		menuView = m.menuList.View()
		contentView = m.viewWalletDetails()
	default:
		menuView = m.menuList.View()
		contentView = "Unknown state."
	}

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		menuStyle.Render(menuView),
		contentStyle.Width(m.width-MENU_WIDTH-2).Height(m.height).Render(contentView),
	)
}

func (m Model) viewWelcome() string {
	return "Welcome to the EITA Wallet Manager!\n\nSelect an option from the menu."
}

func (m Model) viewCreateWalletPassword() string {
	var view strings.Builder
	view.WriteString(lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FF00")).Render("Mnemonic Phrase (Keep it Safe!):\n\n"))
	view.WriteString(fmt.Sprintf("%s\n\n", m.mnemonic))
	view.WriteString("Enter a password to encrypt the wallet:\n\n")
	// Alinhar o campo de entrada de senha à esquerda
	passwordStyle := lipgloss.NewStyle().Align(lipgloss.Left)
	view.WriteString(passwordStyle.Render(m.passwordInput.View()))
	view.WriteString("\n\nPress Enter to continue.")
	return view.String()
}

func (m Model) viewImportWallet() string {
	var view strings.Builder
	view.WriteString(lipgloss.NewStyle().Bold(true).Render("Import Wallet via Mnemonic Phrase\n\n"))
	for i, ti := range m.textInputs {
		if i == m.importStage {
			view.WriteString(fmt.Sprintf("Word %d: %s\n", i+1, ti.View()))
		} else {
			view.WriteString(fmt.Sprintf("Word %d: %s\n", i+1, ti.Value()))
		}
	}
	view.WriteString("\nPress Enter to continue.")
	return view.String()
}

func (m Model) viewImportWalletPassword() string {
	var view strings.Builder
	view.WriteString(lipgloss.NewStyle().Bold(true).Render("Enter a password to encrypt the wallet:\n\n"))
	// Alinhar o campo de entrada de senha à esquerda
	passwordStyle := lipgloss.NewStyle().Align(lipgloss.Left)
	view.WriteString(passwordStyle.Render(m.passwordInput.View()))
	view.WriteString("\n\nPress Enter to continue.")
	return view.String()
}

func (m Model) viewWalletList() string {
	var view strings.Builder
	view.WriteString(m.walletTable.View())
	view.WriteString("\nUse the arrow keys to navigate, Enter to view details, ESC to return to the menu.")
	return view.String()
}

func (m Model) viewWalletPassword() string {
	var view strings.Builder
	view.WriteString(lipgloss.NewStyle().Bold(true).Render("Enter the wallet password:\n\n"))
	view.WriteString(m.passwordInput.View())
	view.WriteString("\n\nPress Enter to continue.")
	return view.String()
}

func (m Model) viewWalletDetails() string {
	if m.selectedWallet == nil || m.selectedPrivateKey == nil {
		return "Select a wallet and enter the password to view the details."
	}

	publicKey := m.selectedPrivateKey.PublicKey

	detailStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#00FF00")).Align(lipgloss.Left).AlignVertical(lipgloss.Left)
	labelStyle := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#FFFFFF")).Align(lipgloss.Left).AlignVertical(lipgloss.Left)
	valueStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("#AAAAAA")).Align(lipgloss.Left).AlignVertical(lipgloss.Left)

	var view strings.Builder
	view.WriteString(detailStyle.Render("Wallet Details\n") + "\n")
	view.WriteString(labelStyle.Render("Ethereum Address: ") + valueStyle.Render(fmt.Sprintf("%s\n", m.selectedWallet.Address)) + "\n")
	view.WriteString(labelStyle.Render("Public Key: ") + valueStyle.Render(fmt.Sprintf("0x%x\n", crypto.FromECDSAPub(&publicKey))) + "\n")
	view.WriteString(labelStyle.Render("Private Key: ") + valueStyle.Render(fmt.Sprintf("0x%x\n", crypto.FromECDSA(m.selectedPrivateKey))) + "\n")
	view.WriteString(labelStyle.Render("Mnemonic Phrase: ") + valueStyle.Render(fmt.Sprintf("%s\n", m.selectedWallet.Mnemonic)))
	view.WriteString("\nPress ESC to return to the wallet list.")
	return view.String()
}

func gerarMnemonico() (string, error) {
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return "", err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", err
	}
	return mnemonic, nil
}

func derivarChavePrivada(frase string) (string, error) {
	if !bip39.IsMnemonicValid(frase) {
		return "", fmt.Errorf("invalid mnemonic phrase")
	}
	seed := bip39.NewSeed(frase, "")
	masterKey, err := bip32.NewMasterKey(seed)
	if err != nil {
		return "", err
	}
	purposeKey, err := masterKey.NewChildKey(bip32.FirstHardenedChild + 44)
	if err != nil {
		return "", err
	}
	coinTypeKey, err := purposeKey.NewChildKey(bip32.FirstHardenedChild + 60)
	if err != nil {
		return "", err
	}
	accountKey, err := coinTypeKey.NewChildKey(bip32.FirstHardenedChild + 0)
	if err != nil {
		return "", err
	}
	changeKey, err := accountKey.NewChildKey(0)
	if err != nil {
		return "", err
	}
	addressKey, err := changeKey.NewChildKey(0)
	if err != nil {
		return "", err
	}
	privateKeyBytes := addressKey.Key
	return hex.EncodeToString(privateKeyBytes), nil
}

func hexToECDSA(hexkey string) (*ecdsa.PrivateKey, error) {
	privateKeyBytes, err := hex.DecodeString(hexkey)
	if err != nil {
		return nil, err
	}
	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return nil, err
	}
	return privateKey, nil
}

func main() {
	db, err := NewDatabase(DB_NAME)
	if err != nil {
		fmt.Println("Error initializing the database:", err)
		os.Exit(1)
	}
	defer func(db *Database) {
		err := db.Close()
		if err != nil {
			fmt.Println("Error closing the database:", err)
		}
	}(db)

	if _, err := os.Stat(WALLETS_DIR); os.IsNotExist(err) {
		err := os.Mkdir(WALLETS_DIR, os.ModePerm)
		if err != nil {
			fmt.Println("Error creating wallets directory:", err)
			os.Exit(1)
		}
	}

	ks := keystore.NewKeyStore(WALLETS_DIR, keystore.StandardScryptN, keystore.StandardScryptP)

	model := InitModel(db, ks)

	p := tea.NewProgram(&model, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		fmt.Println("Error starting the program:", err)
		os.Exit(1)
	}
}
