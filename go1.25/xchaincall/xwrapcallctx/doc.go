package xwrapcallctx

// Rule of thumb for component interfaces:
// - interface including new parts: {base}{Component} (lowercased)
//   - preferred way: include {base} in package name, i.e. `...{base}.{component}`
// - interface including base and new parts: {Base}With{Component}
// - implementation including new parts: {Base}{Component}
//   - {Base} should be repeated even when it is included in package name, i.e. `...{base}.{Base}{Component}`
//     - rationale: differ from others structs of this package belonging to different "category" of structs
//   - reads as base's component (e.g. ContextAbort)
