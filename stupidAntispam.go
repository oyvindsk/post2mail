package post2mail

import (
	"strings"
)

// Ops, got a lot of spam right away, like 200-300 per day :( Fortunately it looks quite easy to filter (#famouslastwords)
// Lobster Thermidor aux crevettes with a Mornay sauce, garnished with truffle pâté, brandy and a fried egg on top, and Spam.

// IsSpam returnes true/false and a reason (currently just 1 is implemented)
func IsSpam(ed EmailData) (bool, string) {

	// Equal name and email/phone? That's weird
	if ed.FromEmail == ed.FromName {

		if !strings.Contains(ed.FromEmail, "@") && !strings.ContainsAny(ed.FromEmail, "0123456789") {
			return true, "name==email, no @ or numbers"
		}

	}

	if strings.HasSuffix(ed.FromEmail, "@mail.ru") {
		return true, "mail from ends with mail.ru"
	}

	return false, ""
}
