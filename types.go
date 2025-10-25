package discordgoi18n

import (
	"io/fs"

	"github.com/bwmarrin/discordgo"
	"github.com/kstoums/discordgo-i18n/logger"
)

// Vars is the collection used to inject variables during translation.
// This type only exists to be less verbose.
type Vars map[string]any

type Translator interface {
	SetDefault(locale discordgo.Locale)
	LoadBundle(locale discordgo.Locale, path string) error
	LoadBundleFS(locale discordgo.Locale, fs fs.FS, path string) error
	LoadBundleContent(locale discordgo.Locale, content map[string]any) error
	Get(locale discordgo.Locale, key string, values Vars) string
	GetArray(locale discordgo.Locale, key string, values Vars) []string
	GetDefault(key string, values Vars) string
	GetDefaultArray(key string, values Vars) []string
	GetLocalizations(key string, variables Vars) *map[discordgo.Locale]string
}

type translatorImpl struct {
	defaultLocale discordgo.Locale
	translations  map[discordgo.Locale]bundle
	loadedBundles map[string]bundle
	logger        logger.Logger
}

type translatorMock struct {
	SetDefaultFunc        func(locale discordgo.Locale)
	LoadBundleFunc        func(locale discordgo.Locale, path string) error
	LoadBundleFSFunc      func(locale discordgo.Locale, fs fs.FS, path string) error
	LoadBundleContentFunc func(locale discordgo.Locale, content map[string]any) error
	GetFunc               func(locale discordgo.Locale, key string, values Vars) string
	GetArrayFunc          func(locale discordgo.Locale, key string, values Vars) []string
	GetDefaultFunc        func(key string, values Vars) string
	GetDefaultArrayFunc   func(key string, values Vars) []string
	GetLocalizationsFunc  func(key string, variables Vars) *map[discordgo.Locale]string
}

type bundle map[string][]string

type source string

const (
	osSource      source = "os"
	fsSource      source = "fs"
	contentSource source = "content"
)
