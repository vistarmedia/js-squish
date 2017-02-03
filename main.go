package main

import (
	"flag"
	"log"
	"os"

	"vistarmedia.com/tool/js-squish"
)

var (
	jsTarName   string
	entrypoint  string
	outputName  string
	environment string
)

func init() {
	flag.StringVar(&jsTarName, "jstar", "", "Path to JSTar")
	flag.StringVar(&entrypoint, "entrypoint", "index.js", "Entrypoint")
	flag.StringVar(&outputName, "output", "", "Squished JS Output")
	flag.StringVar(&environment, "environment", "", "NODE_ENV")
}

func main() {
	flag.Parse()

	if jsTarName == "" || entrypoint == "" || outputName == "" {
		flag.PrintDefaults()
		os.Exit(2)
	}
	var env *string
	if environment != "" {
		env = &environment
	}

	repoFile, err := os.Open(jsTarName)
	if err != nil {
		log.Fatal(err)
	}

	// Note: This will keep the entire source in memory. the `DiskJsTarRepository`
	// implementation is very close to the same speed, but runs afowl of Bazel's
	// sandboxing rules. If you're having memory issues running this probably,
	// this is the likely culprit. See `repository.go` for more details.
	repo, err := jssquish.NewMemoryJsTarRepository(repoFile)
	if err != nil {
		log.Fatal(err)
	}
	defer repo.Close()

	out, err := os.Create(outputName)
	if err != nil {
		log.Fatal(err)
	}

	if err := jssquish.Main(repo, entrypoint, env, out); err != nil {
		log.Fatal(err)
	}
}
