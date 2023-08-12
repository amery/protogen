# protogen

`protogen` is a Go library to write code generators on any language.

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

gen, err := protogen.NewPlugin(opts, req)
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
