package discordgoi18n

import (
	"io/fs"
	"testing"
	"testing/fstest"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestMock(t *testing.T) {
	mock := newMock()

	// Implémentations des mocks
	mock.SetDefaultFunc = func(locale discordgo.Locale) {
		assert.Equal(t, discordgo.EnglishUS, locale)
	}

	mock.LoadBundleFunc = func(locale discordgo.Locale, file string) error {
		assert.Equal(t, discordgo.French, locale)
		assert.Equal(t, "file.json", file)
		return nil
	}

	mock.LoadBundleFSFunc = func(locale discordgo.Locale, f fs.FS, file string) error {
		assert.Equal(t, discordgo.German, locale)
		assert.NotNil(t, f)
		assert.Equal(t, "bundle.json", file)
		return nil
	}

	mock.LoadBundleContentFunc = func(locale discordgo.Locale, content map[string]any) error {
		assert.Equal(t, discordgo.Italian, locale)
		assert.NotNil(t, content)
		return nil
	}

	mock.GetFunc = func(locale discordgo.Locale, key string, variables Vars) string {
		if key == "fail" {
			return key
		}
		assert.Equal(t, discordgo.SpanishES, locale)
		assert.Equal(t, "greeting", key)
		return "Hola"
	}

	mock.GetArrayFunc = func(locale discordgo.Locale, key string, variables Vars) []string {
		if key == "fail_array" {
			return []string{key}
		}
		assert.Equal(t, discordgo.SpanishES, locale)
		assert.Equal(t, "days", key)
		return []string{"lunes", "martes"}
	}

	mock.GetDefaultFunc = func(key string, variables Vars) string {
		if key == "fail_default" {
			return key
		}
		assert.Equal(t, "key", key)
		return "default"
	}

	mock.GetDefaultArrayFunc = func(key string, variables Vars) []string {
		if key == "fail_default_array" {
			return []string{key}
		}
		assert.Equal(t, "key_array", key)
		return []string{"uno", "dos"}
	}

	mock.GetLocalizationsFunc = func(key string, variables Vars) *map[discordgo.Locale]string {
		assert.Equal(t, "welcome", key)
		m := map[discordgo.Locale]string{
			discordgo.EnglishUS: "Hello",
			discordgo.French:    "Bonjour",
		}
		return &m
	}

	// TESTS ------------------

	assert.NotPanics(t, func() { mock.SetDefault(discordgo.EnglishUS) })

	assert.NoError(t, mock.LoadBundle(discordgo.French, "file.json"))

	fsys := fstest.MapFS{"bundle.json": {Data: []byte(`{"example":"value"}`)}}
	assert.NoError(t, mock.LoadBundleFS(discordgo.German, fsys, "bundle.json"))

	assert.NoError(t, mock.LoadBundleContent(discordgo.Italian, map[string]any{"hi": "ciao"}))

	// GET (success)
	assert.Equal(t, "Hola", mock.Get(discordgo.SpanishES, "greeting", nil))

	// GET (error)
	translation := mock.Get(discordgo.SpanishES, "fail", nil)
	assert.Equal(t, "fail", translation)

	// GET ARRAY (success)
	assert.Equal(t, []string{"lunes", "martes"}, mock.GetArray(discordgo.SpanishES, "days", nil))

	// GET ARRAY (error)
	translations := mock.GetArray(discordgo.SpanishES, "fail_array", nil)
	assert.Equal(t, []string{"fail_array"}, translations)

	// GET DEFAULT (success)
	assert.Equal(t, "default", mock.GetDefault("key", nil))

	// GET DEFAULT (error)
	translation = mock.GetDefault("fail_default", nil)
	assert.Equal(t, "fail_default", translation)

	// GET DEFAULT ARRAY (success)
	assert.Equal(t, []string{"uno", "dos"}, mock.GetDefaultArray("key_array", nil))

	// GET DEFAULT ARRAY (error)
	translations = mock.GetDefaultArray("fail_default_array", nil)
	assert.Equal(t, []string{"fail_default_array"}, translations)

	// GET LOCALIZATIONS
	assert.NotNil(t, mock.GetLocalizations("welcome", nil))
	assert.Contains(t, *mock.GetLocalizations("welcome", nil), discordgo.EnglishUS)
	assert.Contains(t, *mock.GetLocalizations("welcome", nil), discordgo.French)
}
