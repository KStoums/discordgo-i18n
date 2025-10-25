package discordgoi18n

import (
	"os"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/kstoums/discordgo-i18n/logger"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

const (
	translatorNominalCase1         = "translatorNominalCase1.json"
	translatorNominalCase2         = "translatorNominalCase2.json"
	translatorFailedUnmarshallCase = "translatorFailedUnmarshallCase.json"
	translatorFileDoesNotExistCase = "translatorFileDoesNotExistCase.json"

	content1 = `
    {
       "hi": ["this is a {{ .Test }}"],
       "with": ["all"],
       "the": ["elements", "we"],
       "can": ["find"],
       "in": ["a","json"],
       "config": ["file", "! {{ .Author }}"],
       "parse": ["{{if $foo}}{{end}}"]
    }
    `

	content2 = `
    {
       "this": ["is a {{ .Test }}"],
       "with.a.file": ["containing", "less", "variables"],
       "bye": ["see you"],
       "parse2": ["{{if $foo}}{{end}}"]
    }
    `

	badContent = `
     [
       "content",
       "not",
       "ok",
       "test"
     ]
    `
)

var (
	translatorTest *translatorImpl
	content1Map    = map[string]any{
		"hi":     []string{"this is a {{ .Test }}"},
		"with":   []string{"all"},
		"the":    []string{"elements", "we"},
		"can":    []string{"find"},
		"in":     []string{"a", "json"},
		"config": []string{"file", "! {{ .Author }}"},
		"parse":  []string{"{{if $foo}}{{end}}"},
	}
	content2Map = map[string]any{
		"this":        []string{"is a {{ .Test }}"},
		"with.a.file": []string{"containing", "less", "variables"},
		"bye":         []string{"see you"},
		"parse2":      []string{"{{if $foo}}{{end}}"},
	}
)

// Setup and teardown for tests
func setUp() {
	translatorTest = NewTranslator(&logger.DummyLogger{}).(*translatorImpl)

	// Create JSON files for testing
	for name, content := range map[string]string{
		translatorNominalCase1:         content1,
		translatorNominalCase2:         content2,
		translatorFailedUnmarshallCase: badContent,
	} {
		if err := os.WriteFile(name, []byte(content), os.ModePerm); err != nil {
			log.Fatal().Err(err).Msgf("'%s' could not be created, test stopped", name)
		}
	}
}

func tearDown() {
	translatorTest = nil
	for _, f := range []string{
		translatorNominalCase1,
		translatorNominalCase2,
		translatorFailedUnmarshallCase,
		translatorFileDoesNotExistCase,
	} {
		if err := os.Remove(f); err != nil {
			log.Warn().Err(err).Msgf("'%s' could not be deleted", f)
		}
	}
}

// Test initialization of translator
func TestNewTranslator(t *testing.T) {
	setUp()
	defer tearDown()
	assert.Empty(t, translatorTest.translations)  // No translations loaded initially
	assert.Empty(t, translatorTest.loadedBundles) // No bundles loaded initially
}

// Test setting default locale
func TestSetDefault(t *testing.T) {
	setUp()
	defer tearDown()
	assert.Equal(t, defaultLocale, translatorTest.defaultLocale)
	translatorTest.SetDefault(discordgo.Italian)
	assert.Equal(t, discordgo.Italian, translatorTest.defaultLocale)
}

// Test loading JSON bundles from files
func TestLoadBundle(t *testing.T) {
	setUp()
	defer tearDown()

	// Nonexistent file returns an error and does not modify state
	assert.Error(t, translatorTest.LoadBundle(discordgo.French, translatorFileDoesNotExistCase))
	assert.Empty(t, translatorTest.translations)
	assert.Empty(t, translatorTest.loadedBundles)

	// Malformed JSON returns an error
	assert.Error(t, translatorTest.LoadBundle(discordgo.French, translatorFailedUnmarshallCase))
	assert.Empty(t, translatorTest.translations)
	assert.Empty(t, translatorTest.loadedBundles)

	// Load valid bundles
	assert.NoError(t, translatorTest.LoadBundle(discordgo.French, translatorNominalCase1))
	assert.Equal(t, 1, len(translatorTest.loadedBundles))
	assert.Equal(t, 1, len(translatorTest.translations))
	assert.Equal(t, 7, len(translatorTest.translations[discordgo.French]))

	assert.NoError(t, translatorTest.LoadBundle(discordgo.French, translatorNominalCase2))
	assert.Equal(t, 2, len(translatorTest.loadedBundles))
	assert.Equal(t, 1, len(translatorTest.translations))
	assert.Equal(t, 4, len(translatorTest.translations[discordgo.French]))

	// Load bundles for different locale
	assert.NoError(t, translatorTest.LoadBundle(discordgo.EnglishGB, translatorNominalCase2))
	assert.Equal(t, 2, len(translatorTest.loadedBundles))
	assert.Equal(t, 2, len(translatorTest.translations))
	assert.Equal(t, 4, len(translatorTest.translations[discordgo.EnglishGB]))

	assert.NoError(t, translatorTest.LoadBundle(discordgo.EnglishGB, translatorNominalCase1))
	assert.Equal(t, 2, len(translatorTest.loadedBundles))
	assert.Equal(t, 2, len(translatorTest.translations))
	assert.Equal(t, 7, len(translatorTest.translations[discordgo.EnglishGB]))
}

// Test loading bundles from an FS
func TestLoadBundleFS(t *testing.T) {
	setUp()
	defer tearDown()
	dirFS := os.DirFS(".") // Use current directory

	assert.Error(t, translatorTest.LoadBundleFS(discordgo.French, dirFS, translatorFileDoesNotExistCase))
	assert.Error(t, translatorTest.LoadBundleFS(discordgo.French, dirFS, translatorFailedUnmarshallCase))

	// Load valid bundles
	assert.NoError(t, translatorTest.LoadBundleFS(discordgo.French, dirFS, translatorNominalCase1))
	assert.NoError(t, translatorTest.LoadBundleFS(discordgo.French, dirFS, translatorNominalCase2))
	assert.NoError(t, translatorTest.LoadBundleFS(discordgo.EnglishGB, dirFS, translatorNominalCase2))
	assert.NoError(t, translatorTest.LoadBundleFS(discordgo.EnglishGB, dirFS, translatorNominalCase1))
}

// Test loading bundles directly from content
func TestLoadBundleContent(t *testing.T) {
	setUp()
	defer tearDown()
	assert.NoError(t, translatorTest.LoadBundleContent(discordgo.French, content1Map))
	assert.NoError(t, translatorTest.LoadBundleContent(discordgo.French, content2Map))
	assert.NoError(t, translatorTest.LoadBundleContent(discordgo.EnglishGB, content2Map))
	assert.NoError(t, translatorTest.LoadBundleContent(discordgo.EnglishGB, content1Map))
}

// Test getting single translations
func TestGet(t *testing.T) {
	setUp()
	defer tearDown()

	// No bundle loaded: return key
	assert.Equal(t, "hi", translatorTest.Get(discordgo.Dutch, "hi", nil))

	// Load bundles
	assert.NoError(t, translatorTest.LoadBundle(discordgo.Dutch, translatorNominalCase1))
	assert.NoError(t, translatorTest.LoadBundle(defaultLocale, translatorNominalCase2))

	// Nonexistent key returns the key itself
	assert.Equal(t, "does_not_exist", translatorTest.Get(discordgo.Dutch, "does_not_exist", nil))

	// Template present but variables missing: return raw template
	assert.Equal(t, "this is a {{ .Test }}", translatorTest.Get(discordgo.Dutch, "hi", nil))

	// Template with variable substitution
	assert.Equal(t, "this is a test :)", translatorTest.Get(discordgo.Dutch, "hi", Vars{"Test": "test :)"}))

	// Default locale fallback
	assert.Equal(t, "bye", translatorTest.Get(discordgo.Dutch, "bye", nil))

	// Invalid template returns key
	assert.Equal(t, "parse", translatorTest.Get(discordgo.Dutch, "parse", Vars{}))

	// Missing variable: fallback to key
	assert.Equal(t, "hi", translatorTest.Get(discordgo.Dutch, "hi", Vars{}))
}

// Test getting arrays of translations
func TestGetArray(t *testing.T) {
	setUp()
	defer tearDown()

	assert.Equal(t, 1, len(translatorTest.GetArray(discordgo.Dutch, "hi", nil)))

	assert.NoError(t, translatorTest.LoadBundle(discordgo.Dutch, translatorNominalCase1))

	// Template with variable missing: returns raw template array
	assert.Equal(t, []string{"this is a {{ .Test }}"}, translatorTest.GetArray(discordgo.Dutch, "hi", nil))

	// Template with variable substituted
	assert.Equal(t, []string{"this is a valeur"}, translatorTest.GetArray(discordgo.Dutch, "hi", Vars{"Test": "valeur"}))

	// Multi-value key
	assert.Equal(t, []string{"elements", "we"}, translatorTest.GetArray(discordgo.Dutch, "the", nil))

	// Nonexistent key returns array with key
	assert.Equal(t, 1, len(translatorTest.GetArray(discordgo.Dutch, "no_exist", nil)))
}

// Test getting default locale translations
func TestGetDefault(t *testing.T) {
	setUp()
	defer tearDown()
	assert.Equal(t, "hi", translatorTest.GetDefault("hi", nil))

	assert.NoError(t, translatorTest.LoadBundle(defaultLocale, translatorNominalCase2))

	assert.Equal(t, "does_not_exist", translatorTest.GetDefault("does_not_exist", nil))
	assert.Equal(t, "is a {{ .Test }}", translatorTest.GetDefault("this", nil))
	assert.Equal(t, "is a test :)", translatorTest.GetDefault("this", Vars{"Test": "test :)"}))
	assert.Equal(t, "parse2", translatorTest.GetDefault("parse2", Vars{}))
	assert.Equal(t, "this", translatorTest.GetDefault("this", Vars{}))
}

// Test getting arrays for default locale
func TestGetDefaultArray(t *testing.T) {
	setUp()
	defer tearDown()

	// No bundle loaded: return key in array
	assert.Equal(t, 1, len(translatorTest.GetDefaultArray("hi", nil)))

	assert.NoError(t, translatorTest.LoadBundle(defaultLocale, translatorNominalCase2))
	assert.Equal(t, []string{"see you"}, translatorTest.GetDefaultArray("bye", nil))

	// Malformed template: return key in array
	assert.Equal(t, 1, len(translatorTest.GetDefaultArray("parse2", Vars{})))
}

// Test getting all localizations for a key
func TestGetLocalizations(t *testing.T) {
	setUp()
	defer tearDown()

	// No bundles loaded: return empty map
	assert.NotNil(t, translatorTest.GetLocalizations("hi", Vars{}))
	assert.Equal(t, 0, len(*translatorTest.GetLocalizations("hi", Vars{})))

	// Load bundles
	assert.NoError(t, translatorTest.LoadBundle(discordgo.Dutch, translatorNominalCase1))
	assert.NoError(t, translatorTest.LoadBundle(defaultLocale, translatorNominalCase2))

	// Key present: all locales with variable provided
	assert.NotNil(t, translatorTest.GetLocalizations("hi", Vars{"Test": "foo"}))
	assert.Equal(t, 2, len(*translatorTest.GetLocalizations("hi", Vars{"Test": "foo"})))

	// Missing variable: returns key
	assert.NotNil(t, translatorTest.GetLocalizations("hi", Vars{}))

	// Key absent in bundles: returns empty or partial map
	assert.NotNil(t, translatorTest.GetLocalizations("inconnue", Vars{}))
}
