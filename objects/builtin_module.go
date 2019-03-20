package objects

// BuiltinModule is an importable module that's written in Go.
type BuiltinModule struct {
	Name  string
	Attrs map[string]Object
}

// Import returns an immutable map for the module.
func (m *BuiltinModule) Import() (interface{}, error) {
	return m.AsImmutableMap(), nil
}

// AsImmutableMap converts builtin module into an immutable map.
func (m *BuiltinModule) AsImmutableMap() *ImmutableMap {
	attrs := make(map[string]Object, len(m.Attrs))
	for k, v := range m.Attrs {
		attrs[k] = v.Copy()
	}

	if m.Name != "" {
		attrs["__module_name__"] = &String{Value: m.Name}
	}

	return &ImmutableMap{Value: attrs}
}
