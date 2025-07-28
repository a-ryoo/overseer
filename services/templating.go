package services

import (
	log "github.com/sirupsen/logrus"
	"os"
	"regexp"
	"strings"
)

type TemplatingService struct{}

func NewTemplatingService() *TemplatingService {
	return &TemplatingService{}
}

func (s *TemplatingService) RenderFile(src, dest string, vaultFunc func(store, path, key string) string) string {
	content, err := os.ReadFile(src)
	if err != nil {
		log.Fatalf("failed to read template file: %v", err)
	}

	text := string(content)

	re := regexp.MustCompile(`\{\{\s*([a-zA-Z0-9/_\-]+)\s*#\s*([a-zA-Z0-9_\-]+)\s*\}\}`)

	result := re.ReplaceAllStringFunc(text, func(match string) string {
		groups := re.FindStringSubmatch(match)
		if len(groups) != 3 {
			return match // skip malformed ones
		}

		storePath := groups[1]
		key := groups[2]

		parts := strings.SplitN(storePath, "/", 2)
		if len(parts) != 2 {
			return match
		}
		store := parts[0]
		path := parts[1]

		value := vaultFunc(store, path, key)
		return value
	})

	return result
}
