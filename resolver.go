package jssquish

import (
	"encoding/json"
	"fmt"
	"path"
	"path/filepath"
	"strings"
)

type resolveError struct {
	require string
	from    string
}

func (re *resolveError) Error() string {
	return fmt.Sprintf("Could not resolve '%s' from '%s'", re.require, re.from)
}

// Module resolution algorithm used in node.js. The basic algorithm is defined
// at https://nodejs.org/api/modules.html
// The same `require` statement can have different resolutions depending on
// which file is doing the requiring. The simplest example is relative imports.
// Once an import has been made absolute, the queried `require` value will be
// cached. The same resulting qualified value will be returned on the next
// invocation with recomputing.
type Resolver struct {
	repo  Repository
	cache map[string]string
}

// Creates a new `Resolver`. This instance will not share a cache with any
// previous instances.
func NewResolver(repo Repository) *Resolver {
	return &Resolver{
		repo:  repo,
		cache: make(map[string]string),
	}
}

// Resolve the `require` request from the given file or directory. This largely
// implements the node resolution algorithm, but excludes certain lookups such
// as built-ins and `.node` files. From the site:
//
//		require(X) from module at path Y
//		1. If X is a core module,
//			a. return the core module
//			b. STOP
//		2. If X begins with './' or '/' or '../'
//			a. LOAD_AS_FILE(Y + X)
//			b. LOAD_AS_DIRECTORY(Y + X)
//		3. LOAD_NODE_MODULES(X, dirname(Y))
//		4. THROW "not found"
func (r *Resolver) Resolve(require, from string) (string, error) {
	if fq, ok := r.cache[require]; ok {
		return fq, nil
	}

	if strings.HasPrefix(require, "./") || strings.HasPrefix(require, "/") ||
		strings.HasPrefix(require, "../") {

		absolute := filepath.Clean(path.Join(from, require))
		if fq, ok := r.cache[absolute]; ok {
			return fq, nil
		}

		if fq, ok := r.resolveAsFile(absolute); ok {
			r.cache[absolute] = fq
			return fq, nil
		}

		if fq, ok := r.resolveAsDirectory(absolute); ok {
			r.cache[absolute] = fq
			return fq, nil
		}
		return "", &resolveError{require, from}
	}

	if fq, ok := r.resolveAsModule(require, path.Dir(from)); ok {
		r.cache[require] = fq
		return fq, nil
	}

	return "", &resolveError{require, from}
}

// Loads the qualified require value assuming its a file. JSON files will be
// loaded, but they will not be converted into JavaScript objects (ie: no export
// statement). Similarly, no check is made for `.node` files. The basic
// algorithm from the node site:
//
//		LOAD_AS_FILE(X)
//		1. If X is a file, load X as JavaScrip text. STOP
//		2. If X.js is a file, load X.js as JavaScript text. STOP
//		2. If X.json is a file, load X.json to a JavaScript Object. STOP
//		4. If X.node is a file, load X.node as a binary addon. STOP
func (r *Resolver) resolveAsFile(require string) (string, bool) {
	if r.repo.IsFile(require) {
		return require, true
	}

	check := require + ".js"
	if r.repo.IsFile(check) {
		return check, true
	}

	check = require + ".json"
	if r.repo.IsFile(check) {
		return check, true
	}

	return "", false
}

// Loads the qualified require value assuming its a directory. As with the rules
// above, it will not define and export for json files, and will ignore `.node`
// binary files. The stated algoithm:
//
//		LOAD_AS_DIRECTORY(X)
//		1. If X/package.json is a file,
//			a. Parse X/package.json, and look for a "main" field
//			b. let M = X + (json main field)
//			c. LOAD_AS_FILE(M)
//		2. If X/index.js is a file, load X/index.json as JavaScript text. STOP
//		3. If X/index.json is a file, load X/index.json to a JavaScript object.
//			 STOP
//		4. If X/index.node is a file, load X/index.node as a binary addon. STOP
func (r *Resolver) resolveAsDirectory(require string) (string, bool) {
	pkgPath := path.Join(require, "package.json")
	if r.repo.IsFile(pkgPath) {
		pkgBody, err := r.repo.Open(pkgPath)
		if err != nil {
			return "", false
		}
		defer pkgBody.Close()

		main := struct {
			Main string `json:"main"`
		}{}
		if err := json.NewDecoder(pkgBody).Decode(&main); err != nil {
			return "", false
		}

		mainPath := path.Join(require, main.Main)
		if fq, ok := r.resolveAsFile(mainPath); ok {
			return fq, ok
		}
	}

	check := path.Join(require, "index.js")
	if r.repo.IsFile(check) {
		return check, true
	}

	check = path.Join(require, "index.json")
	if r.repo.IsFile(check) {
		return check, true
	}

	return "", false
}

// Uses a simplified `node_modules` algorithm to resolve external dependencies.
// In the `js_tar` world, there is not concept of `node_modules`. Because there
// is only one search space, and only one version of any library loaded, this
// only searches `.` as a path. There is a very good chance that cleaning up
// paths before querying can remove the need for this rule altogether. Once
// again, from the site:
//
//		LOAD_NODE_MODULES(X, START)
//		1. let DIRS=NODE_MODULES_PATHS(START)
//		2 for each DIR in DIRS:
//			a. LOAD_AS_FILE(DIR/X)
//			b. LOAD_AS_DIRECTORY(DIR/X)
//
// The NODE_MODULES_PATHS implementation has been omitted.
func (r *Resolver) resolveAsModule(require, start string) (string, bool) {
	dirs := []string{"."}

	for _, dir := range dirs {
		absolute := path.Join(dir, require)

		if fq, ok := r.resolveAsFile(absolute); ok {
			return fq, ok
		}

		if fq, ok := r.resolveAsDirectory(absolute); ok {
			return fq, ok
		}
	}

	return "", false
}
