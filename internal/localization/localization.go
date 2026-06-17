package localization

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"golang.org/x/text/message/catalog"
	"gopkg.in/yaml.v2"
)

type dictionary struct {
	Data map[string]string
}

func (d *dictionary) Lookup(key string) (string, bool) {
	v, ok := d.Data[key]
	return v, ok
}

var globalCatalog *catalog.Builder

// Init loads every pkg/translations/<lang>.yaml file into a global message
// catalog. Called once via fx.Invoke at startup.
func Init() {
	globalCatalog = catalog.NewBuilder()

	translations, err := parseYAMLDict()
	if err != nil {
		panic(err)
	}
	for langCode, dict := range translations {
		langTag := language.MustParse(langCode)
		for key, value := range dict.(*dictionary).Data {
			globalCatalog.SetString(langTag, key, value)
		}
	}
}

func parseYAMLDict() (map[string]catalog.Dictionary, error) {
	_, filename, _, _ := runtime.Caller(0)
	projectRoot := filepath.Dir(filepath.Dir(filepath.Dir(filename)))
	dir := filepath.Join(projectRoot, "pkg", "translations")

	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	translations := map[string]catalog.Dictionary{}
	for _, f := range files {
		yamlFile, err := os.ReadFile(filepath.Join(dir, f.Name()))
		if err != nil {
			return nil, err
		}
		data := map[string]string{}
		if err := yaml.Unmarshal(yamlFile, &data); err != nil {
			return nil, err
		}
		lang := strings.Split(f.Name(), ".")[0]
		translations[lang] = &dictionary{Data: data}
	}
	return translations, nil
}

// GetMessage returns the localized, formatted string for a tag and language.
func GetMessage(tag, languageCode string, args ...interface{}) string {
	langTag, err := language.Parse(languageCode)
	if err != nil {
		fmt.Printf("invalid language code: %s\n", languageCode)
		return ""
	}
	p := message.NewPrinter(langTag, message.Catalog(globalCatalog))
	return p.Sprintf(tag, args...)
}
