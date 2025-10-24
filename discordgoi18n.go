package discordgoi18n

import (
	"io/fs"

	"maps"

	"github.com/bwmarrin/discordgo"
)

//nolint:gochecknoglobals // False positive, cannot be overridden.
var instance translator

func init() {
	instance = newTranslator()
}

// SetDefaults sets the locale used as a fallback.
// Not thread-safe; designed to be called during initialization.
func SetDefault(locale discordgo.Locale) {
	instance.SetDefault(locale)
}

// LoadBundle loads a translation file corresponding to a specified locale.
// Not thread-safe; designed to be called during initialization.
func LoadBundle(locale discordgo.Locale, path string) error {
	return instance.LoadBundle(locale, path)
}

// LoadBundleFS loads a translation file corresponding to a specified locale through fs.FS.
// Not thread-safe; designed to be called during initialization.
func LoadBundleFS(locale discordgo.Locale, fs fs.FS, path string) error {
	return instance.LoadBundleFS(locale, fs, path)
}

// LoadBundleContent loads content corresponding to a specified locale.
// Not thread-safe; designed to be called during initialization.
// Does not work with string array
func LoadBundleContent(locale discordgo.Locale, content map[string]any) error {
	return instance.LoadBundleContent(locale, content)
}

// Get gets a translation corresponding to a locale and a key.
// Optional Vars parameter is used to inject variables in the translation.
// When a key does not match any translations in the desired locale,
// the default locale is used instead. If the situation persists with the fallback,
// key is returned. If more than one translation is available for dedicated key,
// it is picked randomly. Thread-safe.
func Get(locale discordgo.Locale, key string, values ...Vars) (string, error) {
	args := make(Vars)

	for _, variables := range values {
		maps.Copy(args, variables)
	}

	return instance.Get(locale, key, args)
}

// GetArray gets array of translations corresponding to a locale and a key.
// Optional Vars parameter is used to inject variables in the translation.
// When a key does not match any translations in the desired locale,
// the default locale is used instead. If the situation persists with the fallback,
// key is returned. If more than one translation is available for dedicated key,
func GetArray(locale discordgo.Locale, key string, values ...Vars) ([]string, error) {
	args := make(Vars)

	for _, variables := range values {
		maps.Copy(args, variables)
	}

	return instance.GetArray(locale, key, args)
}

// GetDefault gets a translation corresponding to default locale and a key.
// Optional Vars parameter is used to inject variables in the translation.
// When a key does not match any translations in the default locale,
// key is returned. If more than one translation is available for dedicated key,
// it is picked randomly. Thread-safe.
func GetDefault(key string, values ...Vars) (string, error) {
	args := make(Vars)

	for _, variables := range values {
		maps.Copy(args, variables)
	}

	return instance.GetDefault(key, args)
}

func GetDefaultArray(key string, values ...Vars) ([]string, error) {
	args := make(Vars)

	for _, variables := range values {
		maps.Copy(args, variables)
	}

	return instance.GetDefaultArray(key, args)
}

// GetLocalizations retrieves translations from every loaded bundles.
// Aims to simplify discordgo.ApplicationCommand instanciations by providing
// localizations structures that can be used for any localizable field (example:
// command name, description, etc). Thread-safe.
func GetLocalizations(key string, values ...Vars) (*map[discordgo.Locale]string, error) {
	args := make(Vars)

	for _, variables := range values {
		maps.Copy(args, variables)
	}

	return instance.GetLocalizations(key, args)
}
