# JS Squish
Concatenates a bunch of Javascript files together into a single deployable
source file. It's like [Browserify](http://browserify.org/), but faster.

Right now, this acts as a stand-alone binary with a rule implementation, but
once we have something substantial building with it, build should live in
`@io_bazel_rules_js//js/tool`.

## CLI Usage
From the command line, `jssquish` takes a `js_tar` and entrypoint, and creates
an output file. For more information on what a `js_tar` is, see the README for
`//tool/build_rule/rules_js`.

    ```sh
    bazel run //tool/js-squish -- -h
    Usage of js-squish:
      -entrypoint string
          Entrypoint (default "index.js")
      -jstar string
          Path to JSTar
      -output string
          Squished JS Output
    ```

## Build artifact usage
To generate js-squish'd files, include the rule file included in this module and
use the `js_squish` rule.


    ```python
    load('@io_bazel_rules_js//js:def.bzl', 'js_binary')
    load('//tool/js-squish:rules.bzl', 'js_squish')

    js_binary(
      name = 'my-prog',
      src  = 'cool_proj.js',
      deps = ['@react//:lib'],
      main = 'my/cool/prod',
    )

    js_squish(
      name = 'my-prog.dist',
      src  = ':my-prog',
    )
    ```

This will create a build artifact named `my-prog.dist.js`
