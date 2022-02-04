package i18n

import (
	"github.com/nicksnyder/go-i18n/i18n"
)

func MustLoadTranslationFile(filepath string) {
	i18n.MustLoadTranslationFile(filepath)
}

func LoadTranslationFile(filepath string) {
	i18n.LoadTranslationFile(filepath)
}

func Tfunc(languageSource string, languageSources ...string) (i18n.TranslateFunc, error) {
	return i18n.Tfunc(languageSource, languageSources...)
}

func IdentityTfunc() i18n.TranslateFunc {
	return i18n.IdentityTfunc()
}

///V2 ReplacementFunctions:

func MustLoadTranslationFileV2(filepath string) {
}

func LoadTranslationFileV2(filepath string) {

}

func TfuncV2(languageSource string, languageSources ...string) {

}

func IdentityTfuncV2() {

}
