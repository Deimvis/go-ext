# xinvar (xinvarariants)

xinvar (xinvarariants) — Go package with invariants.
Extends Go xcheck module.

## Invariants

Invariants are checks that work with debug build tag only and stop runtime immediately on fail.
Unlike errors and even panics invariants are meant to be non-interceptable failures that must fail as soon as possible
(following [fail fast](https://en.wikipedia.org/wiki/Fail-fast_system) methodology).

### Usage

Avoid misusing invariants. They are meant to be used for code readability, not for any type of input/output validation.

Valid use case:

```go
func Foo(v *int) int {
    if v == nil {
        tmp := 42
        v = &tmp
    }
    
    // <a bunch of code here>
    
    xinvar.NotNilPtr(v)
    // we just remind readers that v is not nil and it's safe to dereference the pointer
    return *v
}
```

Invalid use case:

```go
func Foo(v *int) int {
    xinvar.NotNilPtr(v)
    // we know nothing about v, user could pass nil variable.
    // without debug build tag this check will be skipped and
    // we end up with invalid state.
}
```

But valid use case:

```go
package internal_impl

func Fn1() int {
    // ...
    x := 42
    res := foo(&x)
    // ...
}

func Fn2() int {
    // ...
    x := 24
    res := foo(&x)
    // ...
}

func foo(v *int) int {
    xinvar.NotNilPtr(v)
    // it's an internal implementation, we know all foo use cases
    // and it's our algorithm invariant that value passed to foo is not nil,
    // so we're legit to add this invariant for readability.
}
```
