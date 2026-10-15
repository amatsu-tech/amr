package cli

import "runtime/debug"

// Заполняются через -ldflags при сборке релиза.
var (
	version = "v0.0.0-dev"
	commit  = "unknown"
)

// buildVersion — версия из -ldflags, а если её нет (go install) — версия модуля.
func buildVersion() string {
	if version != "v0.0.0-dev" {
		return version
	}

	if info, ok := debug.ReadBuildInfo(); ok && info.Main.Version != "" && info.Main.Version != "(devel)" {
		return info.Main.Version
	}

	return version
}

// buildCommit — коммит из -ldflags, а если его нет — из данных сборки Go.
func buildCommit() string {
	if commit != "unknown" {
		return commit
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" && len(setting.Value) >= 7 {
				return setting.Value[:7]
			}
		}
	}

	return commit
}
