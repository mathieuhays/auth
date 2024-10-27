package forms

import (
	"net/http"
)

type Form struct {
	Fields           map[string]*Field
	Error            error
	validateCallback func(form *Form)
}

type Field struct {
	Name     string
	Value    string
	Error    error
	Validate func(field *Field, form *Form)
}

func (f *Form) LoadValuesFromRequest(request *http.Request) {
	for _, field := range f.Fields {
		field.Value = request.PostFormValue(field.Name)
	}
}

func (f *Form) SetValidation(validate func(form *Form)) {
	f.validateCallback = validate
}

func (f *Form) Validate() {
	f.ValidateFields()
	if f.validateCallback != nil {
		f.validateCallback(f)
	}
}

func (f *Form) ValidateFields() {
	for _, field := range f.Fields {
		if field.Validate != nil {
			field.Validate(field, f)
		}
	}
}

func (f *Form) HasErrors() bool {
	if f.Error != nil {
		return true
	}

	for _, field := range f.Fields {
		if field.Error != nil {
			return true
		}
	}

	return false
}

func NewForm(fields ...Field) *Form {
	form := Form{
		Fields: map[string]*Field{},
	}

	for _, field := range fields {
		form.Fields[field.Name] = &field
	}

	return &form
}
