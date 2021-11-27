package forms

//holding validation error messages for forms
type errors map[string][]string

//adding error messages for a given field to the map
func (e errors) Add(field, message string) {
	e[field] = append(e[field], message)
}

//get method to retrieve the first error mes.
func (e errors) Get(field string) string {
	es := e[field]
	if len(es) == 0 {
		return ""
	}
	return es[0]
}

