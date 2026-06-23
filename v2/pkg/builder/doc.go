/*
Package builder materializes config entries from res into runtime resources.

Configs are registered in res by library use init (AddWithTags) and by
ecfg.Register from AppConfig. Build walks the registry, calls Build() on
every entry that implements Builder, registers the result (inheriting entry
tags), and removes the config entry.

For generic Builder configs, Build first inspects instantiated type parameters.
When a type parameter implements [Resourcer], BuildResource is called and the
prototype is registered without tags. Non-Resourcer type parameters are skipped.

Builder does not perform DI — wiring happens later via sdi.

Org rule: ordinary configs implement [Builder]; generic configs may declare a
type parameter T that implements [Resourcer] for automatic prototype registration.
Resourcer is never a separate AppConfig field.

Typical pipeline:

	cfg, _ := ecfg.Parse(...)
	ecfg.Register(cfg, res.Default)
	builder.Build(res.Default)
	res.Default.Transform(...)
	sdi.Resolve(res.Default)
*/
package builder
