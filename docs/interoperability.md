# Interoperability  

## Table of Contents

- [Using Scripts](#using-scripts)
  - [Type Conversion Table](#type-conversion-table)
  - [User Types](#user-types)
- [Sandbox Environments](#sandbox-environments)
- [Concurrency](#concurrency)
- [Compiler and VM](#compiler-and-vm)

## Using Scripts

Embedding and executing the Tengo code in Go is very easy. At a high level,
this process is like:

- create a [Script](https://godoc.org/github.com/d5/tengo#Script) instance with
your code,
- _optionally_ add some
[Script Variables](https://godoc.org/github.com/d5/tengo#Variable) to Script,
- compile or directly run the script,
- retrieve _output_ values from the
[Compiled](https://godoc.org/github.com/d5/tengo#Compiled) instance.

The following is an example where a Tengo script is compiled and run with no
input/output variables.

```golang
import "github.com/d5/tengo/v2"

var code = `
reduce := func(seq, fn) {
    s := 0
    for x in seq { fn(x, s) }
    return s
}

print(reduce([1, 2, 3], func(x, s) { s += x }))
`

func main() {
    s := tengo.NewScript([]byte(code))
    if _, err := s.Run(); err != nil {
        panic(err)
    }
}
```

Here's another example where an input variable is added to the script, and, an
output variable is accessed through
[Variable.Int](https://godoc.org/github.com/d5/tengo#Variable.Int) function:

```golang
import (
    "fmt"

    "github.com/d5/tengo/v2"
)

func main() {
    s := tengo.NewScript([]byte(`a := b + 20`))

    // define variable 'b'
    _ = s.Add("b", 10)

    // compile the source
    c, err := s.Compile()
    if err != nil {
        panic(err)
    }

    // run the compiled bytecode
    // a compiled bytecode 'c' can be executed multiple times without re-compiling it
    if err := c.Run(); err != nil {
        panic(err)
    }

    // retrieve value of 'a'
    a := c.Get("a")
    fmt.Println(a.Int())           // prints "30"

    // re-run after replacing value of 'b'
    if err := c.Set("b", 20); err != nil {
        panic(err)
    }
    if err := c.Run(); err != nil {
        panic(err)
    }
    fmt.Println(c.Get("a").Int())  // prints "40"
}
```

A variable `b` is defined by the user before compilation using
[Script.Add](https://godoc.org/github.com/d5/tengo#Script.Add) function. Then a
compiled bytecode `c` is used to execute the bytecode and get the value of
global variables. In this example, the value of global variable `a` is read
using [Compiled.Get](https://godoc.org/github.com/d5/tengo#Compiled.Get)
function. See
[documentation](https://godoc.org/github.com/d5/tengo#Variable) for the
full list of variable value functions.

Value of the global variables can be replaced using
[Compiled.Set](https://godoc.org/github.com/d5/tengo#Compiled.Set) function.
But it will return an error if you try to set the value of un-defined global
variables _(e.g. trying to set the value of `x` in the example)_.  

### Type Conversion Table

When adding a Variable
_([Script.Add](https://godoc.org/github.com/d5/tengo#Script.Add))_, Script
converts Go values into Tengo values based on the following conversion table.

| Go Type | Tengo Type | Note |
| :--- | :--- | :--- |
|`nil`|`Undefined`||
|`string`|`String`||
|`int64`|`Int`||
|`int`|`Int`||
|`bool`|`Bool`||
|`rune`|`Char`||
|`byte`|`Char`||
|`float64`|`Float`||
|`[]byte`|`Bytes`||
|`time.Time`|`Time`||
|`error`|`Error{String}`|use `error.Error()` as String value|
|`map[string]Object`|`Map`||
|`map[string]interface{}`|`Map`|individual elements converted to Tengo objects|
|`[]Object`|`Array`||
|`[]interface{}`|`Array`|individual elements converted to Tengo objects|
|`Object`|`Object`|_(no type conversion performed)_|

### User Types

Users can add and use a custom user type in Tengo code by implementing
[Object](https://godoc.org/github.com/d5/tengo#Object) interface. Tengo runtime
will treat the user types in the same way it does to the runtime types with no
performance overhead. See
[Object Types](https://github.com/d5/tengo/blob/master/docs/objects.md) for
more details.

## Sandbox Environments

To securely compile and execute _potentially_ unsafe script code, you can use
the following Script functions.

### Script.SetImports(modules *objects.ModuleMap)

SetImports sets the import modules with corresponding names. Script **does not**
include any modules by default. You can use this function to include the
[Standard Library](https://github.com/d5/tengo/blob/master/docs/stdlib.md).

```golang
s := tengo.NewScript([]byte(`math := import("math"); a := math.abs(-19.84)`))

s.SetImports(stdlib.GetModuleMap("math"))
// or, to include all stdlib at once
s.SetImports(stdlib.GetModuleMap(stdlib.AllModuleNames()...))
```

You can also include Tengo's written module using `objects.SourceModule`
(which implements `objects.Importable`).

```golang
s := tengo.NewScript([]byte(`double := import("double"); a := double(20)`))

mods := tengo.NewModuleMap()
mods.AddSourceModule("double", []byte(`export func(x) { return x * 2 }`))
s.SetImports(mods)
```

To dynamically load or generate code for imported modules, implement and
provide a `tengo.ModuleGetter`.

```golang
type DynamicModules struct {
  mods tengo.ModuleGetter
  fallback func (name string) tengo.Importable
}
func (dm *DynamicModules) Get(name string) tengo.Importable {
  if mod := dm.mods.Get(name); mod != nil {
    return mod
  }
  return dm.fallback()
}
// ...
mods := &DynamicModules{
  mods: stdlib.GetModuleMap("math"),
  fallback: func(name string) tengo.Importable {
    src := ... // load or generate src for `name`
    return &tengo.SourceModule{Src: src}
  },
}
s := tengo.NewScript(`foo := import("foo")`)
s.SetImports(mods)
```

### Script.SetMaxAllocs(n int64)

SetMaxAllocs sets the maximum number of object allocations. Note this is a
cumulative metric that tracks only the object creations. Set this to a negative
number (e.g. `-1`) if you don't need to limit the number of allocations.

### Script.EnableFileImport(enable bool)

EnableFileImport enables or disables module loading from the local files. It's
disabled by default.

### tengo.MaxStringLen

Sets the maximum byte-length of string values. This limit applies to all
running VM instances in the process. Also it's not recommended to set or update
this value while any VM is executing.

### tengo.MaxBytesLen

Sets the maximum length of bytes values. This limit applies to all running VM
instances in the process. Also it's not recommended to set or update this value
while any VM is executing.

## Concurrency

A compiled script (`Compiled`) can be used to run the code multiple
times by a goroutine. If you want to run the compiled script by multiple
goroutine, you should use `Compiled.Clone` function to make a copy of Compiled
instances.

### Compiled.Clone()

Clone creates a new copy of Compiled instance. Cloned copies are safe for
concurrent use by multiple goroutines. 

```golang
for i := 0; i < concurrency; i++ {
    go func(compiled *tengo.Compiled) {
        // inputs
        _ = compiled.Set("a", rand.Intn(10))
        _ = compiled.Set("b", rand.Intn(10))
        _ = compiled.Set("c", rand.Intn(10))

        if err := compiled.Run(); err != nil {
            panic(err)
        }

        // outputs
        d = compiled.Get("d").Int()
        e = compiled.Get("e").Int()
    }(compiled.Clone()) // Pass the cloned copy of Compiled
}
```

## Compiler and VM

Although it's not recommended, you can directly create and run the Tengo
[Compiler](https://godoc.org/github.com/d5/tengo#Compiler), and
[VM](https://godoc.org/github.com/d5/tengo#VM) for yourself instead of using
Scripts and Script Variables. It's a bit more involved as you have to manage
the symbol tables and global variables between them, but, basically that's what
Script and Script Variable is doing internally.

_TODO: add more information here_
