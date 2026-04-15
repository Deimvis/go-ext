# (Ext)ended Golang Library

Extended Golang Library is an extension for standard golang library.

## Best stuff

#### `xmust` — panic-on-failure checks

```go
conn := xmust.Do(sql.Open("postgres", dsn)) // panics if err != nil
val  := xmust.Ok(cache.Get("key"))          // panics if ok == false

xmust.Eq(got, want, "unexpected result")
// can print values that failed a check
xmust.Eq(got, want, "unexpected result", xcheck.PrintValues)
// can get check error instead of panic with xshould
err := xshould.Eq(got, want, "unexpected result")
```

#### `xoptional` — optional values with JSON/YAML encoding support

```go
v := xoptional.New(42)
v.HasValue() // true
v.Value()    // 42

empty := xoptional.New[string]()
empty.HasValue() // false

val := xoptional.ValueOr(empty, "default") // "default"
```

#### `xwrapcall` — wrapped call framework

```go
type Context = context.Context
// may use custom context, e.g. implement xwrapcall.AbortableContext
// for explicit call abortion semantics, early return hooks, etc
fn := xwrapcall.New[Context]().
    With(
        func(c Context, next xwrapcall.Next[Context]) (Context, error) {
            c = context.WithValue(c, authKey{}, "user-123")
            return next(c)  // call next middleware / action
        },
        func(c Context, next xwrapcall.Next[Context]) (Context, error) {
            start := time.Now()
            c, err := next(c)
            log.Printf("took %s", time.Since(start))
            return c, err
        },
    ).
    Do(myAction)
err := fn(ctx)
// and wrap again :)
fn = xwrapcall.New[Context]().With(smthNew).Do(fn)
```

#### `xfb` — fallback values

Concise conditional fallback — a readable alternative to ternary expressions.

```go
xfb.On(isZero, val, fallback)          // return val if !isZero(val), else fallback
ptr = xfb.OnNil(ptr, defaultPtr)             // ptr or use default ptr fallback
v = xfb.OnNilv(ptr, defaultValue)          // dereference ptr or use default
```

#### `xio` — composable io.Reader wraps

```xslices — slice helpers
var r io.Reader
r, rCount := xio.WrapReader(r, xio.CountingReaderWrap)
r, rChecksum := xio.WrapReader(r, md5ReaderWrap)
r, _ = xio.WrapReader(r, skipFirstByteWrap)
io.ReadAll(r)
fmt.Println(rCount.BytesRead(), rChecksum.Sum())
```

#### `xiter` — iterator helpers

```go
names := xiter.Map(slices.Values(users), func(u User) string {
    return u.Name
})
allPassed := xiter.Reduce(slices.Values(results), xbooliter.All, true)
```

#### `xslices` — slice helpers

```go
age2users := xslices.GroupBy(users, func(u User) string {
    return u.Age
})
```

#### `xinvar` — zero-cost debug invariants

Invariant checks that **compile away** in release builds via build tags. Full validation in development, zero overhead in production.

```go
// These are real checks with `-tags debug`, no-ops otherwise:
xinvar.NotNilPtr(cfg)
xinvar.Eq(len(items), expected)
xinvar.Lt(idx, len(buf))
xinvar.Implements[io.Reader](val)
```

#### `xreflect.FillNilStructPointers` — safe nested struct access

```go
type SearchOptions struct {
  Filters *struct{
    Category *struct{
      Eq *string
    }
  }
  ...
}

var opts SearchOptions
xreflect.FillNilStructPointers(&opts)
cf = opts.Filters.Category // no panic, because nils are replaced with default values
if cf.Eq != nil { ... } // Eq is not struct pointer, its nil is not filled
```


## Conventions

<img src="docs/assets/diablo2_hell.jpg" style="zoom:50%">

<div style="text-align: center"><b>HELL STARTS HERE</b></div>

* Separate module is used for each Go version that is significant for the library implementation (e.g. `go1.25/`).
* Backward compatibility is not a primary concern - anything can be refactored once a better implementation is recognized, even if it breaks compatibility. But the major version will be incremented accordingly.
* Library contains some highly experimental features which can be identified by source file path (e.g. `_exp` file name suffix) or by comments at the top of the source file.

