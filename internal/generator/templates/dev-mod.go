package templates

import (
	"html/template"
	"os"
	"path/filepath"
	"runtime"
)

const devModFile = `module {{.PackageName}}

require github.com/nathan-fiscaletti/coattail-go v0.0.0
replace github.com/nathan-fiscaletti/coattail-go => {{.RelativePath}}

go {{.GoVersion}}
`

func getModulePath() string {
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		return "Unknown path"
	}
	return filepath.Dir(file)
}

type relativePathModTemplateData struct {
	ModTemplateData
	RelativePath string
}

// Writes a development mod file when running from a (devel) build
func WriteDevModFile(path string, data ModTemplateData) error {
	relativePath, err := filepath.Rel(filepath.Dir(path), filepath.Join(getModulePath(), "..", "..", ".."))
	if err != nil {
		return err
	}

	finalData := relativePathModTemplateData{
		ModTemplateData: data,
		RelativePath:    relativePath,
	}

	// use go text-template to render the template
	tmpl, err := template.New("devmod").Parse(devModFile)
	if err != nil {
		return err
	}

	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	return tmpl.Execute(f, finalData)
}
