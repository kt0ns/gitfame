package languagebase

import (
	"strings"
)

type LanguageEntry struct {
	Name       string
	Type       string
	Extensions []string
}

type DataLoader interface {
	LoadData() ([]LanguageEntry, error)
}

type Base struct {
	extensionsMap map[string][]string
}

func New(loader DataLoader) (*Base, error) {
	entries, err := loader.LoadData()
	if err != nil {
		return nil, err
	}

	extMap := make(map[string][]string)

	for _, entry := range entries {
		key := strings.ToLower(entry.Name)
		extMap[key] = entry.Extensions
	}

	return &Base{
		extensionsMap: extMap,
	}, nil
}

func (b *Base) Resolve(language string) ([]string, bool) {
	exts, found := b.extensionsMap[strings.ToLower(language)]
	return exts, found
}
