package main

import (
	"github.com/amery/protogen/pkg/protogen"
)

// INDENT is the intentation of each level
const INDENT = `    `

func logRequest(gen *protogen.Plugin) {
	logParameters(gen)

	gen.ForEachFile(func(f *protogen.File) {
		logFile(gen, ``, f)
	})
}

func logParameters(gen *protogen.Plugin) {
	gen.Println("parameters:")
	for k, v := range gen.Params() {
		gen.Printf("  %q: %q", k, v)
	}
}

func logFile(gen *protogen.Plugin, indent string, f *protogen.File) {
	indent2 := indent + INDENT

	if f.Generate() {
		gen.Printf("%s%s: %q", indent, "source", f.Name())
	} else {
		gen.Printf("%s%s: %q", indent, "import", f.Name())
	}

	for _, dep := range f.Dependencies() {
		gen.Printf("%s%s: %q", indent2, "import", dep.Name())
	}

	logEnums(gen, indent2, f.Enums())
	logMessages(gen, indent2, f.Messages())
}

func logEnums(gen *protogen.Plugin, indent string, enums []*protogen.Enum) {
	indent2 := indent + INDENT

	for _, p := range enums {
		gen.Printf("%s%s: %s (%s)", indent, "enum", p.Name(), p.FullName())
		logEnumValues(gen, indent2, p.Values())
	}
}

func logEnumValues(gen *protogen.Plugin, indent string, values []*protogen.EnumValue) {
	for _, p := range values {
		gen.Printf("%s[%v]: %s", indent, p.Number(), p.Name())
	}
}

func logMessages(gen *protogen.Plugin, indent string, msgs []*protogen.Message) {
	indent2 := indent + INDENT

	for _, p := range msgs {
		gen.Printf("%s%s: %s (%s)", indent, "message", p.Name(), p.FullName())
		logEnums(gen, indent2, p.Enums())
		logMessages(gen, indent2, p.Messages())
		logFields(gen, indent2, p.Fields())
	}
}

func logFields(gen *protogen.Plugin, indent string, fields []*protogen.Field) {
	for _, p := range fields {
		gen.Printf("%s%s: %s", indent, "field", p.Name())
	}
}
