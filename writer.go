package jssquish

import (
	"encoding/json"
	"fmt"
	"io"
	"text/template"

	"vistarmedia.com/tool/js-squish/file"
)

var (
	entryPre = template.Must(
		template.New("entry").Parse("{{.Id}}: [function(require,module,exports) {\n"),
	)
	entryPost = template.Must(
		template.New("entry").Parse("\n}, {{.Imports}}]"),
	)

	preamble = file.MustAsset("tool/js-squish/preamble.js")

	preambleTemplate = template.Must(
		template.New("entry").Parse(string(preamble)),
	)

	postamble = `},{}, [0]);`
)

type Writer struct {
	w           io.Writer
	firstModule bool
}

func NewWriter(w io.Writer) *Writer {
	return &Writer{w, true}
}

func (w *Writer) Open() error {
	return w.OpenWithEnvironment(nil)
}

func (w *Writer) OpenWithEnvironment(environment *string) error {
	// Write the preamble function, and start to invoke the function with the
	// first argument as an object of modules
	var env string
	if environment != nil {
		env = "'" + *environment + "'"
	} else {
		env = "undefined"
	}

	entry := struct {
		Environment string
	}{env}

	if err := preambleTemplate.Execute(w.w, entry); err != nil {
		return err
	}
	_, err := fmt.Fprint(w.w, "({")
	return err
}

func (w *Writer) Close() error {
	// Close the object of modules, and pass the other two arguments to the anon
	// function defined in the preamble (module cache, and starting module index
	// -- always 0).
	_, err := fmt.Fprint(w.w, "},{},[0]);")
	return err
}

func (w *Writer) Write(src io.Reader, id int, deps map[string]*srcEntry) error {
	// Serialize imports as a json object
	importsMap, err := w.importsMap(deps)
	if err != nil {
		return err
	}

	entry := struct {
		Id      int
		Imports string
	}{id, importsMap}

	// If we are the first module written, prefix with a newline. If not, prefix
	// with a comma and newline
	if w.firstModule {
		if _, err = fmt.Fprint(w.w, "\n"); err != nil {
			return err
		}
		w.firstModule = false
	} else {
		if _, err = fmt.Fprint(w.w, ",\n"); err != nil {
			return err
		}
	}

	// Write entry preamble
	if err = entryPre.Execute(w.w, entry); err != nil {
		return err
	}

	// Write entry body
	if _, err = io.Copy(w.w, src); err != nil {
		return err
	}

	// Write entry postamble
	if err = entryPost.Execute(w.w, entry); err != nil {
		return err
	}

	return nil
}

func (w *Writer) importsMap(deps map[string]*srcEntry) (string, error) {
	imports := make(map[string]int, len(deps))
	for impt, entry := range deps {
		imports[impt] = entry.id
	}
	if bs, err := json.Marshal(imports); err != nil {
		return "", err
	} else {
		return string(bs), nil
	}
}
