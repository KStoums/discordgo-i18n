package discordgoi18n

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/fs"
	"math/rand"
	"os"
	"strings"
	"text/template"

	"github.com/bwmarrin/discordgo"
	"github.com/kstoums/discordgo-i18n/logger"
)

const (
	defaultLocale   = discordgo.EnglishUS
	leftDelim       = "{{"
	rightDelim      = "}}"
	keyDelim        = "."
	executionPolicy = "missingkey=error"
)

func NewTranslator(logger logger.Logger, defaultLocale discordgo.Locale) Translator {
	return &translatorImpl{
		defaultLocale: defaultLocale,
		translations:  make(map[discordgo.Locale]bundle),
		loadedBundles: make(map[string]bundle),
		logger:        logger,
	}
}

func (translator *translatorImpl) LoadBundle(locale discordgo.Locale, path string) error {
	cachePath := translator.buildCachePath(path, osSource)
	loadedBundle, found := translator.loadedBundles[cachePath]
	if !found {
		buf, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return translator.loadBundleBuf(locale, buf, cachePath)
	}

	translator.translations[locale] = loadedBundle
	return nil
}

func (translator *translatorImpl) LoadBundleFS(locale discordgo.Locale, fsys fs.FS, path string) error {
	cachePath := translator.buildCachePath(path, fsSource)
	loadedBundle, found := translator.loadedBundles[cachePath]
	if !found {
		buf, err := fs.ReadFile(fsys, path)
		if err != nil {
			return err
		}

		return translator.loadBundleBuf(locale, buf, cachePath)
	}

	translator.translations[locale] = loadedBundle
	return nil
}

func (translator *translatorImpl) LoadBundleContent(locale discordgo.Locale, content map[string]any) error {
	cachePath := translator.buildCachePath(fmt.Sprintf("%p", content), contentSource)
	loadedBundle, found := translator.loadedBundles[cachePath]
	if !found {
		newBundle := translator.mapBundleStructure(content)
		translator.loadedBundles[cachePath] = newBundle
		translator.translations[locale] = newBundle
		return nil
	}

	translator.translations[locale] = loadedBundle
	return nil
}

func (translator *translatorImpl) Get(locale discordgo.Locale, key string, variables Vars) string {
	bundles, found := translator.translations[locale]
	if !found {
		if locale != translator.defaultLocale {
			translator.logger.Error().Err(fmt.Errorf("bundle '%s' is not loaded, trying to translate key '%s' in '%s'",
				locale, key, translator.defaultLocale))
			return key
		}

		translator.logger.Error().Err(fmt.Errorf("bundle '%s' is not loaded, cannot translate '%s', key returned", locale, key))
		return key
	}

	raws, found := bundles[key]
	if !found || len(raws) == 0 {
		if locale != translator.defaultLocale {
			translator.logger.Error().Err(fmt.Errorf("no label found for key '%s' in '%s', trying to translate it in %s",
				key, locale, translator.defaultLocale))
			return key
		}

		translator.logger.Error().Err(fmt.Errorf("no label found for key '%s' in '%s'", key, locale))
		return key
	}

	//nolint:gosec // No need to have a strong random number generator here.
	raw := raws[rand.Intn(len(raws))]

	if variables != nil && strings.Contains(raw, leftDelim) {
		t, err := template.New("").Delims(leftDelim, rightDelim).Option(executionPolicy).Parse(raw)
		if err != nil {
			translator.logger.Error().Err(err).Msgf("Cannot parse raw corresponding to key '%s' in '%s'", locale, key)
			return key
		}

		var buf bytes.Buffer
		err = t.Execute(&buf, variables)
		if err != nil {
			translator.logger.Error().Err(err).Msgf("Cannot inject variables in raw corresponding to key '%s' in '%s', raw returned", locale, key)
			return key
		}
		return buf.String()
	}

	return raw
}

func (translator *translatorImpl) GetArray(locale discordgo.Locale, key string, variables Vars) []string {
	bundles, found := translator.translations[locale]
	if !found {
		if locale != translator.defaultLocale {
			translator.logger.Error().Err(fmt.Errorf("bundle '%s' is not loaded, trying to translate key '%s' in '%s'",
				locale, key, translator.defaultLocale))
			return []string{key}
		}

		translator.logger.Error().Err(fmt.Errorf("bundle '%s' is not loaded, cannot translate '%s', key returned", locale, key))
		return []string{key}
	}

	raws, found := bundles[key]
	if !found || len(raws) == 0 {
		if locale != translator.defaultLocale {
			translator.logger.Error().Err(fmt.Errorf("no label found for key '%s' in '%s'", key, locale))
			return []string{key}
		}

		translator.logger.Error().Err(fmt.Errorf("no label found for key '%s' in '%s'", key, locale))
		return []string{key}
	}

	//nolint:gosec // No need to have a strong random number generator here.
	for i, raw := range raws {
		if variables != nil && strings.Contains(raw, leftDelim) {
			t, err := template.New("").Delims(leftDelim, rightDelim).Option(executionPolicy).Parse(raw)
			if err != nil {
				translator.logger.Error().Err(err).Msgf("Cannot parse raw corresponding to key '%s' in '%s'", key, locale)
				return []string{key}
			}

			var buf bytes.Buffer
			err = t.Execute(&buf, variables)
			if err != nil {
				translator.logger.Error().Err(err).Msgf("Cannot inject variables in raw corresponding to key '%s' in '%s'", key, locale)
				return []string{key}
			}

			raws[i] = buf.String()
		}
	}

	return raws
}

func (translator *translatorImpl) GetDefault(key string, variables Vars) string {
	return translator.Get(translator.defaultLocale, key, variables)
}

func (translator *translatorImpl) GetDefaultArray(key string, variables Vars) []string {
	return translator.GetArray(translator.defaultLocale, key, variables)
}

func (translator *translatorImpl) GetLocalizations(key string, variables Vars) *map[discordgo.Locale]string {
	localizations := make(map[discordgo.Locale]string)

	for locale := range translator.translations {
		localizations[locale] = translator.Get(locale, key, variables)
	}

	return &localizations
}

func (translator *translatorImpl) loadBundleBuf(locale discordgo.Locale, buf []byte, cachePath string) error {
	var jsonContent map[string]any
	err := json.Unmarshal(buf, &jsonContent)
	if err != nil {
		return err
	}

	newBundle := translator.mapBundleStructure(jsonContent)

	translator.logger.Debug().Msgf("Bundle '%s' loaded with '%s' content", locale, cachePath)
	translator.loadedBundles[cachePath] = newBundle
	translator.translations[locale] = newBundle
	return nil
}

func (translator *translatorImpl) mapBundleStructure(jsonContent map[string]any) bundle {
	bundle := make(map[string][]string)
	for key, content := range jsonContent {
		switch v := content.(type) {
		case string:
			bundle[key] = []string{v}
		case []any:
			values := make([]string, 0)
			for _, value := range v {
				values = append(values, fmt.Sprintf("%v", value))
			}
			bundle[key] = values
		case map[string]any:
			subValues := translator.mapBundleStructure(v)
			for subKey, subValue := range subValues {
				bundle[fmt.Sprintf("%s%s%s", key, keyDelim, subKey)] = subValue
			}
		default:
			bundle[key] = []string{fmt.Sprintf("%v", v)}
		}
	}

	return bundle
}

func (translator *translatorImpl) buildCachePath(path string, source source) string {
	return fmt.Sprintf("%v:%v", source, path)
}
