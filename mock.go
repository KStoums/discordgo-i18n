package discordgoi18n

import (
	"errors"
	"io/fs"

	"github.com/bwmarrin/discordgo"
)

func newMock() *translatorMock {
	return &translatorMock{}
}

func (mock *translatorMock) SetDefault(locale discordgo.Locale) {
	if mock.SetDefaultFunc != nil {
		mock.SetDefaultFunc(locale)
		return
	}
}

func (mock *translatorMock) LoadBundle(locale discordgo.Locale, file string) error {
	if mock.LoadBundleFunc != nil {
		return mock.LoadBundleFunc(locale, file)
	}
	return errors.New("LoadBundle not mocked")
}

func (mock *translatorMock) LoadBundleFS(locale discordgo.Locale, fs fs.FS, file string) error {
	if mock.LoadBundleFSFunc != nil {
		return mock.LoadBundleFSFunc(locale, fs, file)
	}
	return errors.New("LoadBundleFS not mocked")
}

func (mock *translatorMock) LoadBundleContent(locale discordgo.Locale, content map[string]any) error {
	if mock.LoadBundleContentFunc != nil {
		return mock.LoadBundleContentFunc(locale, content)
	}
	return errors.New("LoadBundleContent not mocked")
}

func (mock *translatorMock) Get(locale discordgo.Locale, key string, variables Vars) (string, error) {
	if mock.GetFunc != nil {
		return mock.GetFunc(locale, key, variables)
	}
	return "", errors.New("Get not mocked")
}

func (mock *translatorMock) GetArray(locale discordgo.Locale, key string, variables Vars) ([]string, error) {
	if mock.GetArrayFunc != nil {
		return mock.GetArrayFunc(locale, key, variables)
	}
	// Retourne un tableau vide + erreur
	return nil, errors.New("GetArray not mocked")
}

func (mock *translatorMock) GetDefault(key string, variables Vars) (string, error) {
	if mock.GetDefaultFunc != nil {
		return mock.GetDefaultFunc(key, variables)
	}
	return "", errors.New("GetDefault not mocked")
}

func (mock *translatorMock) GetDefaultArray(key string, variables Vars) ([]string, error) {
	if mock.GetDefaultArrayFunc != nil {
		return mock.GetDefaultArrayFunc(key, variables)
	}
	return nil, errors.New("GetDefaultArray not mocked")
}

func (mock *translatorMock) GetLocalizations(key string, variables Vars) (*map[discordgo.Locale]string, error) {
	if mock.GetLocalizationsFunc != nil {
		return mock.GetLocalizationsFunc(key, variables)
	}
	m := make(map[discordgo.Locale]string)
	return &m, nil
}
