package forms

type errors map[string][]string

// Add adds an error message to a given html form field
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

// Get returns the first error message of a given html form field
func (e errors) Get(field string) string {
	ss := e[field]
	if len(ss) == 0 {
		return ""
	}

	return ss[0]
}
