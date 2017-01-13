def _js_squish_impl(ctx):
  bin_target = ctx.attr.src

  arguments = [
    '-jstar',       bin_target.js_tar.path,
    '-output',      ctx.outputs.out.path,
    '-entrypoint',  bin_target.main.path,
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
    'src':   attr.label(providers=['js_tar', 'main']),

    '_js_squish': attr.label(
      default     = Label('//tool/js-squish'),
      cfg         = 'host',
      executable  = True),
  },
  outputs = {
    'out':  '%{name}.js'
  },
)
