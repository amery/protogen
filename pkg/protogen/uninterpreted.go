package protogen

import "google.golang.org/protobuf/types/descriptorpb"

// UninterpretedOption determines if the given
// [descriptorpb.UninterpretedOption] matches the requested name.
// If it matches, it will return the value and true.
// Otherwise, it will return "" and false.
func UninterpretedOption(uo *descriptorpb.UninterpretedOption, name string) (string, bool) {
	if uo != nil {
		if sp := uo.IdentifierValue; sp != nil {
			if name == *sp {
				// match
				return string(uo.StringValue), true
			}
		}
	}
	return "", false
}
