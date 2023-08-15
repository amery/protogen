# protogen

[![Go Reference][godoc-badge]][godoc]
[![Go Report Card][goreport-badge]][goreport]

`protogen` is a Go library to write proto code generators for any language.

[godoc]: https://pkg.go.dev/github.com/amery/protogen
[godoc-badge]: https://pkg.go.dev/badge/github.com/amery/protogen.svg
[goreport]: https://goreportcard.com/report/github.com/amery/protogen
[goreport-badge]: https://goreportcard.com/badge/github.com/amery/protogen

## Setting `protogen` up

The recommended way is to begin with a `protogen.Options` instance and run it against
a handler.

```go
opts := &protogen.Options{
    // ...
}

err := opts.Run(func (gen *protogen.Plugin) error {
    // ...
})
if err != nil {
    log.Fatal(err)
}
```

In this scenario the `CodeGeneratorRequest` will be parsed from `opts.Stdin` and
the `CodeGeneratorResponse` will be written to `opts.Stdout` as expected by `protoc`.

For testing, you can provide a preparsed `CodeGeneratorRequest` manually.

```go
opts := &protogen.Options{
    // ...
}

req := &pluginpb.CodeGeneratorRequest{
    // ...
}

gen, err := opts.New(req)
if err != nil {
    log.Fatal(err)
}

err = h(gen)
if err != nil {
    log.Fatal(err)
}

// extract the `pluginpb.CodeGEneratorResponse` object
resp := gen.Response()
// or write it to opts.Stdout
_, err = gen.Write()
if err != nil {
    log.Fatal(err)
}
```
