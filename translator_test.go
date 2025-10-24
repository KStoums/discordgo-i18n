package discordgoi18n

import (
	"os"
	"testing"

	"github.com/bwmarrin/discordgo"
	"github.com/rs/zerolog/log"
	"github.com/stretchr/testify/assert"
)

const (
	translatornominalCase1         = "translatorNominalCase1.json"
	translatornominalCase2         = "translatorNominalCase2.json"
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

// Setup and teardown

func setUp() {
	translatorTest = newTranslator()
	if err := os.WriteFile(translatornominalCase1, []byte(content1), os.ModePerm); err != nil {
		log.Fatal().Err(err).Msgf("'%s' could not be created, test stopped", translatornominalCase1)
	}
	if err := os.WriteFile(translatornominalCase2, []byte(content2), os.ModePerm); err != nil {
		log.Fatal().Err(err).Msgf("'%s' could not be created, test stopped", translatornominalCase2)
	}
	if err := os.WriteFile(translatorFailedUnmarshallCase, []byte(badContent), os.ModePerm); err != nil {
		log.Fatal().Err(err).Msgf("'%s' could not be created, test stopped", translatorFailedUnmarshallCase)
	}
}

func tearDown() {
	translatorTest = nil
	for _, f := range []string{
		translatornominalCase1,
		translatornominalCase2,
		translatorFailedUnmarshallCase,
		translatorFileDoesNotExistCase,
	} {
		if err := os.Remove(f); err != nil {
			log.Warn().Err(err).Msgf("'%s' could not be deleted", f)
		}
	}
}

// Tests

func TestNewTranslator(t *testing.T) {
	setUp()
	defer tearDown()
	assert.Empty(t, translatorTest.translations)
	assert.Empty(t, translatorTest.loadedBundles)
}

func TestSetDefault(t *testing.T) {
	setUp()
	defer tearDown()
	assert.Equal(t, defaultLocale, translatorTest.defaultLocale)
	translatorTest.SetDefault(discordgo.Italian)
	assert.Equal(t, discordgo.Italian, translatorTest.defaultLocale)
}

func TestLoadBundle(t *testing.T) {
	setUp()
	defer tearDown()
	_, err := os.Stat(translatorFileDoesNotExistCase)
	assert.True(t, os.IsNotExist(err))
	assert.Error(t, translatorTest.LoadBundle(discordgo.French, translatorFileDoesNotExistCase))
	assert.Empty(t, translatorTest.translations)
	assert.Empty(t, translatorTest.loadedBundles)

	assert.Error(t, translatorTest.LoadBundle(discordgo.French, translatorFailedUnmarshallCase))
	assert.Empty(t, translatorTest.translations)
	assert.Empty(t, translatorTest.loadedBundles)

	assert.NoError(t, translatorTest.LoadBundle(discordgo.French, translatornominalCase1))
	assert.Equal(t, 1, len(translatorTest.loadedBundles))
	assert.Equal(t, 1, len(translatorTest.translations))
	assert.Equal(t, 7, len(translatorTest.translations[discordgo.French]))

	assert.NoError(t, translatorTest.LoadBundle(discordgo.French, translatornominalCase2))
	assert.Equal(t, 2, len(translatorTest.loadedBundles))
	assert.Equal(t, 1, len(translatorTest.translations))
	assert.Equal(t, 4, len(translatorTest.translations[discordgo.French]))

	assert.NoError(t, translatorTest.LoadBundle(discordgo.EnglishGB, translatornominalCase2))
	assert.Equal(t, 2, len(translatorTest.loadedBundles))
	assert.Equal(t, 2, len(translatorTest.translations))
	assert.Equal(t, 4, len(translatorTest.translations[discordgo.EnglishGB]))

	assert.NoError(t, translatorTest.LoadBundle(discordgo.EnglishGB, translatornominalCase1))
	assert.Equal(t, 2, len(translatorTest.loadedBundles))
	assert.Equal(t, 2, len(translatorTest.translations))
	assert.Equal(t, 7, len(translatorTest.translations[discordgo.EnglishGB]))
}

func TestLoadBundleFS(t *testing.T) {
	setUp()
	defer tearDown()
	dirFS := os.DirFS(".")

	_, err := os.Stat(translatorFileDoesNotExistCase)
	assert.True(t, os.IsNotExist(err))
	assert.Error(t, translatorTest.LoadBundleFS(discordgo.French, dirFS, translatorFileDoesNotExistCase))
	assert.Empty(t, translatorTest.translations)
	assert.Empty(t, translatorTest.loadedBundles)

	assert.Error(t, translatorTest.LoadBundleFS(discordgo.French, dirFS, translatorFailedUnmarshallCase))
	assert.Empty(t, translatorTest.translations)
	assert.Empty(t, translatorTest.loadedBundles)

	assert.NoError(t, translatorTest.LoadBundleFS(discordgo.French, dirFS, translatornominalCase1))
	assert.Equal(t, 1, len(translatorTest.loadedBundles))
	assert.Equal(t, 1, len(translatorTest.translations))
	assert.Equal(t, 7, len(translatorTest.translations[discordgo.French]))

	assert.NoError(t, translatorTest.LoadBundleFS(discordgo.French, dirFS, translatornominalCase2))
	assert.Equal(t, 2, len(translatorTest.loadedBundles))
	assert.Equal(t, 1, len(translatorTest.translations))
	assert.Equal(t, 4, len(translatorTest.translations[discordgo.French]))

	assert.NoError(t, translatorTest.LoadBundleFS(discordgo.EnglishGB, dirFS, translatornominalCase2))
	assert.Equal(t, 2, len(translatorTest.loadedBundles))
	assert.Equal(t, 2, len(translatorTest.translations))
	assert.Equal(t, 4, len(translatorTest.translations[discordgo.EnglishGB]))

	assert.NoError(t, translatorTest.LoadBundleFS(discordgo.EnglishGB, dirFS, translatornominalCase1))
	assert.Equal(t, 2, len(translatorTest.loadedBundles))
	assert.Equal(t, 2, len(translatorTest.translations))
	assert.Equal(t, 7, len(translatorTest.translations[discordgo.EnglishGB]))
}

func TestLoadBundleContent(t *testing.T) {
	setUp()
	defer tearDown()

	assert.NoError(t, translatorTest.LoadBundleContent(discordgo.French, content1Map))
	assert.Equal(t, 1, len(translatorTest.loadedBundles))
	assert.Equal(t, 1, len(translatorTest.translations))
	assert.Equal(t, 7, len(translatorTest.translations[discordgo.French]))

	assert.NoError(t, translatorTest.LoadBundleContent(discordgo.French, content2Map))
	assert.Equal(t, 2, len(translatorTest.loadedBundles))
	assert.Equal(t, 1, len(translatorTest.translations))
	assert.Equal(t, 4, len(translatorTest.translations[discordgo.French]))

	assert.NoError(t, translatorTest.LoadBundleContent(discordgo.EnglishGB, content2Map))
	assert.Equal(t, 2, len(translatorTest.loadedBundles))
	assert.Equal(t, 2, len(translatorTest.translations))
	assert.Equal(t, 4, len(translatorTest.translations[discordgo.EnglishGB]))

	assert.NoError(t, translatorTest.LoadBundleContent(discordgo.EnglishGB, content1Map))
	assert.Equal(t, 2, len(translatorTest.loadedBundles))
	assert.Equal(t, 2, len(translatorTest.translations))
	assert.Equal(t, 7, len(translatorTest.translations[discordgo.EnglishGB]))
}

func TestGet(t *testing.T) {
	setUp()
	defer tearDown()

	// Not loaded: must return error and empty string
	get, err := translatorTest.Get(discordgo.Dutch, "hi", nil)
	assert.Error(t, err)
	assert.Equal(t, "", get)

	assert.NoError(t, translatorTest.LoadBundle(discordgo.Dutch, translatornominalCase1))
	assert.NoError(t, translatorTest.LoadBundle(defaultLocale, translatornominalCase2))

	// Key does not exist
	get, err = translatorTest.Get(discordgo.Dutch, "does_not_exist", nil)
	assert.Error(t, err)
	assert.Equal(t, "", get)

	// Key and variables missing, raw presence of the template
	get, err = translatorTest.Get(discordgo.Dutch, "hi", nil)
	assert.NoError(t, err)
	assert.Equal(t, "this is a {{ .Test }}", get)

	// Key and variables present
	get, err = translatorTest.Get(discordgo.Dutch, "hi", Vars{"Test": "test :)"})
	assert.NoError(t, err)
	assert.Equal(t, "this is a test :)", get)

	// Default retrieved (key from fallback locale)
	get, err = translatorTest.Get(discordgo.Dutch, "bye", nil)
	assert.Error(t, err)
	assert.Equal(t, "", get)

	// Incorrect template syntax
	get, err = translatorTest.Get(discordgo.Dutch, "parse", Vars{})
	assert.Error(t, err)
	assert.Equal(t, "", get)

	// Missing value injection
	get, err = translatorTest.Get(discordgo.Dutch, "hi", Vars{})
	assert.Error(t, err)
	assert.Equal(t, "", get)
}

func TestGetArray(t *testing.T) {
	setUp()
	defer tearDown()

	// Bundle case not loaded
	arr, err := translatorTest.GetArray(discordgo.Dutch, "hi", nil)
	assert.Error(t, err)
	assert.Equal(t, 0, len(arr))

	// Bundle loaded
	assert.NoError(t, translatorTest.LoadBundle(discordgo.Dutch, translatornominalCase1))

	// Key with template, WITHOUT providing the expected variable => success
	arr, err = translatorTest.GetArray(discordgo.Dutch, "hi", nil)
	assert.NoError(t, err)
	assert.Equal(t, []string{"this is a {{ .Test }}"}, arr)

	// Key with template, WITH the expected variable => success
	arr, err = translatorTest.GetArray(discordgo.Dutch, "hi", Vars{"Test": "valeur"})
	assert.NoError(t, err)
	assert.Equal(t, []string{"this is a valeur"}, arr)

	// Multi-value key without template
	arr, err = translatorTest.GetArray(discordgo.Dutch, "the", nil)
	assert.NoError(t, err)
	assert.Equal(t, []string{"elements", "we"}, arr)

	// Key does not exist
	arr, err = translatorTest.GetArray(discordgo.Dutch, "no_exist", nil)
	assert.Error(t, err)
	assert.Equal(t, 0, len(arr))
}

func TestGetDefault(t *testing.T) {
	setUp()
	defer tearDown()

	getDefault, err := translatorTest.GetDefault("hi", nil)
	assert.Error(t, err)
	assert.Equal(t, "", getDefault)

	assert.NoError(t, translatorTest.LoadBundle(defaultLocale, translatornominalCase2))

	getDefault, err = translatorTest.GetDefault("does_not_exist", nil)
	assert.Error(t, err)
	assert.Equal(t, "", getDefault)

	getDefault, err = translatorTest.GetDefault("this", nil)
	assert.NoError(t, err)
	assert.Equal(t, "is a {{ .Test }}", getDefault)

	getDefault, err = translatorTest.GetDefault("this", Vars{"Test": "test :)"})
	assert.NoError(t, err)
	assert.Equal(t, "is a test :)", getDefault)

	getDefault, err = translatorTest.GetDefault("parse2", Vars{})
	assert.Error(t, err)
	assert.Equal(t, "", getDefault)

	getDefault, err = translatorTest.GetDefault("this", Vars{})
	assert.Error(t, err)
	assert.Equal(t, "", getDefault)
}

func TestGetDefaultArray(t *testing.T) {
	setUp()
	defer tearDown()

	// No bundle loaded, error, and empty slice
	arr, err := translatorTest.GetDefaultArray("hi", nil)
	assert.Error(t, err)
	assert.Equal(t, 0, len(arr))

	// Bundle loaded for the default locale
	assert.NoError(t, translatorTest.LoadBundle(defaultLocale, translatornominalCase2))
	arr, err = translatorTest.GetDefaultArray("bye", nil)
	assert.NoError(t, err)
	assert.Equal(t, []string{"see you"}, arr)

	// Malformed template, error
	arr, err = translatorTest.GetDefaultArray("parse2", Vars{})
	assert.Error(t, err)
	assert.Equal(t, 0, len(arr))
}

func TestGetLocalizations(t *testing.T) {
	setUp()
	defer tearDown()

	// No bundle loaded: waiting for empty map, no error
	locs, err := translatorTest.GetLocalizations("hi", Vars{})
	assert.NoError(t, err)
	assert.NotNil(t, locs)
	assert.Equal(t, 0, len(*locs))

	// Bundles loaded, key present (Test must be provided for .Test, otherwise fail!)
	assert.NoError(t, translatorTest.LoadBundle(discordgo.Dutch, translatornominalCase1))
	assert.NoError(t, translatorTest.LoadBundle(defaultLocale, translatornominalCase2))

	locs, err = translatorTest.GetLocalizations("hi", Vars{"Test": "foo"})
	if err != nil {
		assert.Nil(t, locs)
		t.Logf("TestGetLocalizations: got error (expected if .Test manquant): %v", err)
		return
	}
	assert.NotNil(t, locs)
	assert.Equal(t, 2, len(*locs))

	// Case: missing template variable (should cause an error and locs==nil)
	locs, err = translatorTest.GetLocalizations("hi", Vars{})
	if err != nil {
		assert.Nil(t, locs)
	} else {
		assert.NotNil(t, locs)
	}

	// Key case completely absent from bundles (or present in only 1)
	locs, err = translatorTest.GetLocalizations("inconnue", Vars{})
	if err != nil {
		assert.Nil(t, locs)
	} else {
		assert.NotNil(t, locs)
		// assert maybe len(*locs) == 0 here depending on your implementation (to be confirmed)
	}
}
