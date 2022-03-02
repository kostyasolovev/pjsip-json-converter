package pjsip

import "reflect"

// Contains structs registry.
type TypeRegistry map[string]reflect.Type

// Creates empty Registry instance.
func newRegMap() TypeRegistry { return make(TypeRegistry) }

// Puts struct type into Registry.
func (t *TypeRegistry) register(in interface{}) {
	typo := reflect.TypeOf(in).Elem()
	(*t)[typo.Name()] = typo
}

// Inits and fills Registry of structs
// Создает мапу[имя_структуры]reflect.Type(структуры).
func InitRegistration() TypeRegistry {
	t := newRegMap()
	t.register((*Endpoint)(nil))
	t.register((*Trunk)(nil))
	t.register((*Auth)(nil))
	t.register((*Aor)(nil))
	t.register((*Acl)(nil))
	t.register((*Transport)(nil))
	t.register((*Identify)(nil))
	t.register((*Registration)(nil))
	t.register((*ResourceList)(nil))
	t.register((*Phoneprov)(nil))
	t.register((*Domain)(nil))
	t.register((*System)(nil))
	t.register((*Global)(nil))
	t.register((*UnknownStruct)(nil))

	return t
}

// Creates struct according to chosen name, returns pointer to it.
// Создает структуры по их названию.
func (t *TypeRegistry) MakeStruct(name string) interface{} {
	return reflect.New((*t)[name]).Interface()
}
