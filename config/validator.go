package config

import (
	"fmt"
	"sort"
	"strings"
)

// Manage validation results and help for built an error object to notify validation errors.
//
// Usage:
//    c := &DummyConfig{}
//    v := &Validator{}
//    v.Validate("name", c.Name != "").SetReason("name is empty")
//    _, err := net.ParseIP(c.Address)
//    v.ValidateError("address", err).SetReason("invalid IP address")
//    v.Merge("vms", c.VMs.Validate())
//    return v.Error()
type Validator struct {
	e errors_
}
type FieldName string
type errors_ struct {
	m map[FieldName][]*Error
}
type Error struct {
	Field FieldName
	// [Optional] Reason why validation error occurred.
	Reason string
	// [Optional] Original error object.
	Err error
}
type DetailError interface {
	Error() string
	DetailError() string
}

func (v *Validator) Validate(field string, ok bool) *Error {
	if ok {
		return nil
	}
	ee := &Error{
		Field: FieldName(field),
	}
	v.e.append(ee)
	return ee
}
func (v *Validator) ValidateError(field string, err error) *Error {
	if err == nil {
		return nil
	}
	ee := &Error{
		Field: FieldName(field),
		Err:   err,
	}
	v.e.append(ee)
	return ee
}
func (v *Validator) Merge(field string, err error) {
	if err == nil {
		return
	}
	switch e := err.(type) {
	case *errors_:
		v.e.merge(field, e)
	default:
		v.e.append(&Error{
			Field: FieldName(field),
			Err:   err,
		})
	}
}
func (v *Validator) Error() error {
	if len(v.e.m) == 0 {
		return nil
	}
	return &v.e
}
func (fn *FieldName) AppendPrefix(prefix string) FieldName {
	if prefix == "" {
		return *fn
	}
	return FieldName(prefix + "." + string(*fn))
}
func (e *errors_) append(ee *Error) {
	if e.m == nil {
		e.m = map[FieldName][]*Error{}
	}
	if e.m[ee.Field] == nil {
		e.m[ee.Field] = []*Error{}
	}
	e.m[ee.Field] = append(e.m[ee.Field], ee)
}
func (e *errors_) merge(prefix string, src *errors_) {
	if e.m == nil {
		e.m = map[FieldName][]*Error{}
	}
	for key, errs := range src.m {
		newKey := key.AppendPrefix(prefix)
		if e.m[newKey] != nil {
			// Key already exists.
			panic(fmt.Errorf("config errors: found an duplicated key(%s)", newKey))
		}
		for i := range errs {
			errs[i].Field = errs[i].Field.AppendPrefix(prefix)
		}
		e.m[newKey] = errs
	}
}
func (e *errors_) Error() string {
	if e == nil {
		return ""
	}
	return fmt.Sprintf("found invalid values in %s", strings.Join(e.keys(), ", "))
}
func (e *errors_) DetailError() string {
	b := strings.Builder{}
	fmt.Fprintf(&b, "%s:", e.Error())
	for _, key := range e.keys() {
		errs := e.m[FieldName(key)]
		for i := range errs {
			fmt.Fprintf(&b, "\n  - %s", errs[i])
		}
	}
	return b.String()
}
func (e *errors_) keys() []string {
	keys := make([]string, 0, len(e.m))
	for key := range e.m {
		keys = append(keys, string(key))
	}
	sort.Strings(keys)
	return keys
}
func (e *Error) Error() string {
	return fmt.Sprintf("%s: %s", e.Field, e.detail())
}
func (e *Error) SetReason(reason string) {
	if e == nil {
		return
	}
	e.Reason = reason
}
func (e *Error) detail() string {
	if e.Reason != "" {
		return e.Reason
	} else if e.Err != nil {
		return e.Err.Error()
	} else {
		return "unknown error"
	}
}
