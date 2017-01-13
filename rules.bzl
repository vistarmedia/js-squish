load('@io_bazel_rules_js//js/private:rules.bzl', 'js_bin_providers')


def _js_squish_impl(ctx):
  bin_target = ctx.attr.js_binary

  arguments = [
    '-jstar',       bin_target.js_tar.path,
    '-output',      ctx.outputs.out.path,
    '-entrypoint',  bin_target.main,
  ]

  ctx.action(
    inputs     = [ctx.executable._js_squish, bin_target.js_tar],
    outputs    = [ctx.outputs.out],
    executable = ctx.executable._js_squish,
    arguments  = arguments,
    mnemonic   = 'JsSquish',
  )

  return struct(
    files    = set([ctx.outputs.out]),
    runfiles = ctx.runfiles(files = [ctx.outputs.out]),
  )


js_squish = rule(
  _js_squish_impl,
  attrs = {
    'js_binary':  attr.label(providers=js_bin_providers),

    '_js_squish': attr.label(
      default     = Label('//tool/js-squish'),
      cfg         = 'host',
      executable  = True),
  },
  outputs = {
    'out':  '%{name}.js'
  },
)
