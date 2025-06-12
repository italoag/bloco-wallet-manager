package localization

import (
	"fmt"

	"github.com/nicksnyder/go-i18n/v2/i18n"
)

// T retorna a mensagem traduzida para o ID especificado
// data é um mapa de variáveis de template que serão substituídas na mensagem
func T(messageID string, data map[string]interface{}) string {
	if localizer == nil {
		return messageID
	}

	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: data,
	})

	if err != nil {
		return messageID
	}

	return msg
}

// TP retorna a mensagem traduzida plural para o ID especificado
// count é o número que determina a forma plural a ser usada
// data é um mapa de variáveis de template que serão substituídas na mensagem
func TP(messageID string, count interface{}, data map[string]interface{}) string {
	if localizer == nil {
		return messageID
	}

	// Se data for nil, inicialize-o
	if data == nil {
		data = make(map[string]interface{})
	}

	// Adicionar o count aos dados do template
	data["Count"] = count

	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		PluralCount:  count,
		TemplateData: data,
	})

	if err != nil {
		return fmt.Sprintf("%s (error: %v)", messageID, err)
	}

	return msg
}

// ChangeLanguage altera o idioma do localizer
func ChangeLanguage(lang string) {
	if bundle == nil {
		return
	}

	localizer = i18n.NewLocalizer(bundle, lang)

	// Atualizar o mapa global Labels para refletir o novo idioma
	err := populateLabelsMap()
	if err != nil {
		return
	}
}
