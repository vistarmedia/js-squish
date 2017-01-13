load('@io_bazel_rules_js//js:def.bzl', 'js_library', 'js_binary')

load('//tool/js-squish:rules.bzl', 'js_squish')

js_library(
  name = 'simple_export',
  srcs = ['simple_export.js'],
)

js_library(
  name = 'simple_import',
  srcs = ['simple_import.js'],
  deps = [':simple_export']
)

js_binary(
  name = 'require_once',
  src  = 'require_once.js',
  deps = [':simple_import'],
)

js_squish(
  name = 'once.squished',
  src  = ':require_once',
)

js_binary(
  name = 'require_twice',
  src  = 'require_twice.js',
  deps = [':simple_import'],
)

js_squish(
  name = 'twice.squished',
  src  = ':require_twice',
)

sh_test(
  name = 'test',
  size = 'small',
  srcs = ['test.sh'],
  data = [
    ':once.squished',
    ':twice.squished',

    '@io_bazel_rules_js//js/toolchain:node',
  ],
)

test_suite(
  name  = 'tests',
  tests = [':test'],
)
