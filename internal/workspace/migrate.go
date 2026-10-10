package workspace

import (
	"fmt"
)

// Миграции старых версий конфига в текущую. Когда появится v2,
// сюда добавится case 1 с переводом v1 -> v2.
func migrateConfig(data []byte, version int) (Config, error) {
	switch version {
	default:
		return Config{}, fmt.Errorf("%w: %d", ErrUnsupportedConfigVersion, version)
	}
}
