package JSON

import (
	"encoding/json"

	"gitlab.com/slon/shad-go/gitfame/internal/config/languagebase"
)

type JSONLoader struct {
	Data []byte
}

func (j *JSONLoader) LoadData() ([]languagebase.LanguageEntry, error) {
	var entries []languagebase.LanguageEntry
	if err := json.Unmarshal(j.Data, &entries); err != nil {
		return nil, err
	}

	return entries, nil
}
