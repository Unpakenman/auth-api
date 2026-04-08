package pg

import (
	"bytes"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

// Repository stores SQL templates.
type Repository struct {
	templates map[string]*template.Template // namespace: template
}

type Parser interface {
	AddRoot(root string, pattern string) (err error)
	AddFSRoot(paths []string, files fs.FS, pattern string) error
	AddFiles(mappingNamespaceFiles []MappingNamespaceFiles, pattern string) error
	Get(name string) (string, error)
	Exec(name string, data interface{}) (string, error)
	Parse(name string, data interface{}) (string, error)
}

// NewParser creates a new Parser.
func NewParser() Parser {
	return &Repository{
		templates: make(map[string]*template.Template),
	}
}

// Add adds a root directory to the repository, recursively. Match only the
// given file extension. Blocks on the same namespace will be overridden. Does
// not follow symbolic links.
func (r *Repository) AddRoot(root string, pattern string) (err error) {
	// List the directories
	dirs := []string{}
	err = filepath.Walk(
		root,
		func(path string, info os.FileInfo, e error) error {
			if e != nil {
				return nil
			}
			if info.IsDir() {
				dirs = append(dirs, path)
			}
			return nil
		},
	)
	if err != nil {
		return err
	}

	// Add each sub-directory as namespace
	for _, dir := range dirs {
		d := strings.Split(dir, string(os.PathSeparator))
		ro := strings.Split(root, string(os.PathSeparator))
		namespace := strings.Join(d[len(ro):], "/")
		err = r.addDir(dir, namespace, pattern)
		if err != nil {
			return err
		}
	}
	return nil
}

// addDir parses a directory.
func (r *Repository) addDir(path, namespace, pattern string) error {
	// Parse the template
	t, err := template.ParseGlob(filepath.Join(path, pattern))
	if err != nil {
		r.templates[namespace] = template.New("")
		return err
	}
	r.templates[namespace] = t
	return nil
}

// AddFSRoot parses a directory.
func (r *Repository) AddFSRoot(paths []string, files fs.FS, pattern string) error {
	namespace := ""
	for _, path := range paths {
		t, err := template.ParseFS(files, path+pattern)
		if err != nil {
			r.templates[namespace] = template.New("")
			return err
		}
		r.templates[namespace] = t
	}
	return nil
}

// AddFiles mappping namespace and files
func (r *Repository) AddFiles(mappingNamespaceFiles []MappingNamespaceFiles, pattern string) error {
	for _, mapping := range mappingNamespaceFiles {
		t, err := template.ParseFS(mapping.QueryFiles, mapping.PathToDbQueries+pattern)
		if err != nil {
			r.templates[mapping.Namespace] = template.New("")
			return err
		}
		r.templates[mapping.Namespace] = t
	}
	return nil
}

// Get is a shortcut for r.Exec(), passing nil as data.
func (r *Repository) Get(name string) (string, error) {
	return r.Exec(name, nil)
}

// Exec is a shortcut for r.Parse(), but panics if an error occur.
func (r *Repository) Exec(name string, data interface{}) (string, error) {
	s, err := r.Parse(name, data)
	if err != nil {
		return "", err
	}
	return s, nil
}

// Parse executes the template and returns the resulting SQL or an error.
func (r *Repository) Parse(name string, data interface{}) (string, error) {
	// Prepare namespace and block name
	if name == "" {
		return "", fmt.Errorf("unnamed block")
	}
	path := strings.Split(name, "/")
	namespace := strings.Join(path[0:len(path)-1], "/")
	if namespace == "." {
		namespace = ""
	}
	block := path[len(path)-1]

	// Execute the template
	var b bytes.Buffer
	t, ok := r.templates[namespace]
	if !ok {
		return "", fmt.Errorf("unknown namespace \"%s\"", namespace)
	}
	err := t.ExecuteTemplate(&b, block, data)
	if err != nil {
		return "", err
	}
	return b.String(), nil
}
