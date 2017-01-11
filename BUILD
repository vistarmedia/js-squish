load('@io_bazel_rules_go//go:def.bzl',
  'go_binary',
  'go_library',
  'go_test')

load('//tool/build_rule:doc.bzl', 'doc', 'go_doc')
load('//tool/build_rule:markdown.bzl', 'markdown')
load('//tool/build_rule:go.bzl', 'go_bindata')


go_bindata(
  name    = 'file',
  package = 'file',
  files   = ['preamble.js'],
)

go_library(
  name = 'go_default_library',
  srcs = [
    'ast.go',
    'file_set.go',
    'jssquish.go',
    'parser.go',
    'repository.go',
    'resolver.go',
    'writer.go',
  ],
  deps = [
    ':file',

    '@otto//ast:go_default_library',
    '@otto//parser:go_default_library',
  ],
)

go_test(
  name = 'test',
  size = 'small',
  srcs = [
    'parser_test.go',
    'repository_test.go',
    'resolver_test.go',
    'test.go',
  ],
  deps = [
    '//vendor/github.com/onsi:ginkgo',
    '//vendor/github.com/onsi:gomega',
  ],
  library = 'go_default_library',
)

go_binary(
  name       = 'js-squish',
  srcs       = ['main.go'],
  deps       = [':go_default_library'],
  visibility = ['//visibility:public'],
)

markdown(
  name = 'readme',
  srcs = ['README.md'],
)

go_doc(
  name    = 'godoc',
  library = ':go_default_library',
)

doc(
  srcs=[
    ':godoc',
    ':readme',
  ]
)

test_suite(
  name  = 'tests',
  tests = [
    ':test',
    '//tool/js-squish/example:tests',
  ],
)
