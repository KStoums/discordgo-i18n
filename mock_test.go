package discordgoi18n

import (
	"errors"
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

	mock.GetFunc = func(locale discordgo.Locale, key string, variables Vars) (string, error) {
		if key == "fail" {
			return "", errors.New("mock Get error")
		}
		assert.Equal(t, discordgo.SpanishES, locale)
		assert.Equal(t, "greeting", key)
		return "Hola", nil
	}

	mock.GetArrayFunc = func(locale discordgo.Locale, key string, variables Vars) ([]string, error) {
		if key == "fail_array" {
			return nil, errors.New("mock GetArray error")
		}
		assert.Equal(t, discordgo.SpanishES, locale)
		assert.Equal(t, "days", key)
		return []string{"lunes", "martes"}, nil
	}

	mock.GetDefaultFunc = func(key string, variables Vars) (string, error) {
		if key == "fail_default" {
			return "", errors.New("mock GetDefault error")
		}
		assert.Equal(t, "key", key)
		return "default", nil
	}

	mock.GetDefaultArrayFunc = func(key string, variables Vars) ([]string, error) {
		if key == "fail_default_array" {
			return nil, errors.New("mock GetDefaultArray error")
		}
		assert.Equal(t, "key_array", key)
		return []string{"uno", "dos"}, nil
	}

	mock.GetLocalizationsFunc = func(key string, variables Vars) (*map[discordgo.Locale]string, error) {
		assert.Equal(t, "welcome", key)
		m := map[discordgo.Locale]string{
			discordgo.EnglishUS: "Hello",
			discordgo.French:    "Bonjour",
		}
		return &m, nil
	}

	// TESTS ------------------

	assert.NotPanics(t, func() { mock.SetDefault(discordgo.EnglishUS) })

	assert.NoError(t, mock.LoadBundle(discordgo.French, "file.json"))

	fsys := fstest.MapFS{"bundle.json": {Data: []byte(`{"example":"value"}`)}}
	assert.NoError(t, mock.LoadBundleFS(discordgo.German, fsys, "bundle.json"))

	assert.NoError(t, mock.LoadBundleContent(discordgo.Italian, map[string]any{"hi": "ciao"}))

	// GET (success)
	res, err := mock.Get(discordgo.SpanishES, "greeting", nil)
	assert.NoError(t, err)
	assert.Equal(t, "Hola", res)

	// GET (error)
	_, err = mock.Get(discordgo.SpanishES, "fail", nil)
	assert.Error(t, err)

	// GET ARRAY (success)
	arr, err := mock.GetArray(discordgo.SpanishES, "days", nil)
	assert.NoError(t, err)
	assert.Equal(t, []string{"lunes", "martes"}, arr)

	// GET ARRAY (error)
	_, err = mock.GetArray(discordgo.SpanishES, "fail_array", nil)
	assert.Error(t, err)

	// GET DEFAULT (success)
	def, err := mock.GetDefault("key", nil)
	assert.NoError(t, err)
	assert.Equal(t, "default", def)

	// GET DEFAULT (error)
	_, err = mock.GetDefault("fail_default", nil)
	assert.Error(t, err)

	// GET DEFAULT ARRAY (success)
	defArr, err := mock.GetDefaultArray("key_array", nil)
	assert.NoError(t, err)
	assert.Equal(t, []string{"uno", "dos"}, defArr)

	// GET DEFAULT ARRAY (error)
	_, err = mock.GetDefaultArray("fail_default_array", nil)
	assert.Error(t, err)

	// GET LOCALIZATIONS
	loc, err := mock.GetLocalizations("welcome", nil)
	assert.NoError(t, err)
	assert.NotNil(t, loc)
	assert.Contains(t, *loc, discordgo.EnglishUS)
	assert.Contains(t, *loc, discordgo.French)
}
