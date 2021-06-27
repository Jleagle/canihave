package helpers

import (
	"strings"

	"github.com/mssola/user_agent"
)

func IsBot(userAgent string) bool {

	userAgent = strings.ToLower(userAgent)

	if userAgent == "" ||
		strings.Contains(userAgent, "bot") ||
		strings.Contains(userAgent, "crawl") {
		return true
	}

	return user_agent.New(userAgent).Bot()
}
