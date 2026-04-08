package verb

// the following will enable quick additions for functions
// to the global template that every
// page and component will copy from.

func (v *Verb) Func(name string, f any) {
	v.functions[name] = f
}
