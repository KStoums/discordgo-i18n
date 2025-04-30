package discordgoi18n

import (
	"embed"
	"io/fs"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
)

func TestFacade(t *testing.T) {
	var expectedFile, expectedKey = "File", "Key"
	var expectedFS = embed.FS{}
	var expectedContent = make(map[string]any)
	var expectedValues Vars
	var called bool

	mock := newMock()
	mock.SetDefaultFunc = func(locale discordgo.Locale) {
		assert.Equal(t, discordgo.Italian, locale)
		called = true
	}
	mock.LoadBundleFunc = func(locale discordgo.Locale, file string) error {
		assert.Equal(t, discordgo.French, locale)
		assert.Equal(t, expectedFile, file)
		called = true
		return nil
	}
	mock.LoadBundleFSFunc = func(locale discordgo.Locale, fs fs.FS, path string) error {
		assert.Equal(t, discordgo.Czech, locale)
		assert.Equal(t, expectedFS, fs)
		assert.Equal(t, expectedFile, path)
		called = true
		return nil
	}
	mock.LoadBundleContentFunc = func(locale discordgo.Locale, content map[string]any) error {
		assert.Equal(t, discordgo.Danish, locale)
		assert.Equal(t, expectedContent, content)
		called = true
		return nil
	}
	mock.GetFunc = func(locale discordgo.Locale, key string, values Vars) string {
		assert.Equal(t, discordgo.ChineseCN, locale)
		assert.Equal(t, expectedValues, values)
		assert.Equal(t, expectedKey, key)
		called = true
		return ""
	}
	mock.GetDefaultFunc = func(key string, values Vars) string {
		assert.Equal(t, expectedValues, values)
		assert.Equal(t, expectedKey, key)
		called = true
		return ""
	}
	mock.GetLocalizationsFunc = func(key string, values Vars) *map[discordgo.Locale]string {
		assert.Equal(t, expectedValues, values)
		assert.Equal(t, expectedKey, key)
		called = true
		return nil
	}

	instance = mock

	called = false
	SetDefault(discordgo.Italian)
	assert.True(t, called)

	called = false
	assert.NoError(t, LoadBundle(discordgo.French, expectedFile))
	assert.True(t, called)

	called = false
	assert.NoError(t, LoadBundleFS(discordgo.Czech, expectedFS, expectedFile))
	assert.True(t, called)

	called = false
	assert.NoError(t, LoadBundleContent(discordgo.Danish, expectedContent))
	assert.True(t, called)

	called = false
	expectedValues = make(Vars)
	Get(discordgo.ChineseCN, expectedKey)
	assert.True(t, called)

	called = false
	expectedValues = Vars{
		"Hi": "There",
	}
	Get(discordgo.ChineseCN, expectedKey, expectedValues)
	assert.True(t, called)

	called = false
	expectedValues = Vars{
		"Hi":  "There",
		"Bye": "See u",
	}
	Get(discordgo.ChineseCN, expectedKey, Vars{"Hi": "There"}, Vars{"Bye": "See u"})
	assert.True(t, called)

	called = false
	expectedValues = make(Vars)
	GetDefault(expectedKey)
	assert.True(t, called)

	called = false
	expectedValues = Vars{
		"Hi": "There",
	}
	GetDefault(expectedKey, expectedValues)
	assert.True(t, called)

	called = false
	expectedValues = Vars{
		"Hi":  "There",
		"Bye": "See u",
	}
	GetDefault(expectedKey, Vars{"Hi": "There"}, Vars{"Bye": "See u"})
	assert.True(t, called)

	called = false
	GetLocalizations(expectedKey, Vars{"Hi": "There"}, Vars{"Bye": "See u"})
	assert.True(t, called)
}
