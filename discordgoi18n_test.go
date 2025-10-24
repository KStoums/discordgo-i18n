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
	// ---------------------- Mock impls ----------------------
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

	mock.LoadBundleFSFunc = func(locale discordgo.Locale, f fs.FS, path string) error {
		assert.Equal(t, discordgo.Czech, locale)
		assert.Equal(t, expectedFS, f)
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

	mock.GetFunc = func(locale discordgo.Locale, key string, values Vars) (string, error) {
		assert.Equal(t, discordgo.ChineseCN, locale)
		assert.Equal(t, expectedValues, values)
		assert.Equal(t, expectedKey, key)
		called = true
		return "val", nil
	}

	mock.GetArrayFunc = func(locale discordgo.Locale, key string, values Vars) ([]string, error) {
		assert.Equal(t, discordgo.ChineseCN, locale)
		assert.Equal(t, expectedValues, values)

		// Allows both key values depending on the scenario
		if key == "fail_array" {
			called = true
			return nil, assert.AnError
		}

		assert.Equal(t, expectedKey, key) // only checks in case of success
		called = true
		return []string{"a", "b"}, nil
	}

	mock.GetDefaultFunc = func(key string, values Vars) (string, error) {
		called = true
		assert.Equal(t, expectedValues, values)
		assert.Equal(t, expectedKey, key)
		return "default", nil
	}

	mock.GetDefaultArrayFunc = func(key string, values Vars) ([]string, error) {
		assert.Equal(t, expectedValues, values)
		if key == "fail_default_array" {
			called = true
			return nil, assert.AnError
		}
		assert.Equal(t, expectedKey, key)
		called = true
		return []string{"x", "y"}, nil
	}

	mock.GetLocalizationsFunc = func(key string, values Vars) (*map[discordgo.Locale]string, error) {
		called = true
		assert.Equal(t, expectedValues, values)
		assert.Equal(t, expectedKey, key)
		m := map[discordgo.Locale]string{
			discordgo.EnglishUS: "ok",
		}
		return &m, nil
	}

	// Replace the actual instance
	instance = mock

	// ---------------------- Facade tests ----------------------
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

	// Test Get()
	called = false
	expectedValues = Vars{"Hi": "There"}
	val, err := Get(discordgo.ChineseCN, expectedKey, expectedValues)
	assert.NoError(t, err)
	assert.Equal(t, "val", val)
	assert.True(t, called)

	// ✓ Test GetArray success
	called = false
	expectedValues = Vars{"Array": "OK"}
	arr, err := GetArray(discordgo.ChineseCN, expectedKey, expectedValues)
	assert.NoError(t, err)
	assert.Equal(t, []string{"a", "b"}, arr)
	assert.True(t, called)

	// ✗ Test GetArray error
	called = false
	expectedValues = Vars{"Array": "Fail"}
	_, err = GetArray(discordgo.ChineseCN, "fail_array", expectedValues)
	assert.Error(t, err)
	assert.True(t, called)

	// ✓ Test GetDefault
	called = false
	expectedValues = Vars{"Default": "One"}
	def, err := GetDefault(expectedKey, expectedValues)
	assert.NoError(t, err)
	assert.Equal(t, "default", def)
	assert.True(t, called)

	// ✓ Test GetDefaultArray success
	called = false
	expectedValues = Vars{"Array": "Default"}
	arr, err = GetDefaultArray(expectedKey, expectedValues)
	assert.NoError(t, err)
	assert.Equal(t, []string{"x", "y"}, arr)
	assert.True(t, called)

	// ✗ Test GetDefaultArray error
	called = false
	expectedValues = Vars{"Array": "Fail"}
	_, err = GetDefaultArray("fail_default_array", expectedValues)
	assert.Error(t, err)
	assert.True(t, called)

	// ✓ Test GetLocalizations
	called = false
	expectedValues = Vars{"A": "B"}
	locs, err := GetLocalizations(expectedKey, expectedValues)
	assert.NoError(t, err)
	assert.True(t, called)
	assert.NotNil(t, locs)
	assert.Contains(t, *locs, discordgo.EnglishUS)
}
