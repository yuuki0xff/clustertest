package config

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestError_Error(t *testing.T) {
	type fields struct {
		Field  FieldName
		Reason string
		Err    error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"not detail info",
			fields{Field: FieldName("field")},
			"field: unknown error",
		}, {
			"an reason",
			fields{Field: FieldName("field"), Reason: "invalid url"},
			"field: invalid url",
		}, {
			"an error",
			fields{Field: FieldName("field"), Err: fmt.Errorf("not found resource")},
			"field: not found resource",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Error{
				Field:  tt.fields.Field,
				Reason: tt.fields.Reason,
				Err:    tt.fields.Err,
			}
			assert.Equal(t, tt.want, e.Error())
		})
	}
}

func Test_errors__Error(t *testing.T) {
	type fields struct {
		m map[FieldName][]*Error
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			"an error",
			fields{
				m: map[FieldName][]*Error{
					"field": []*Error{
						{Field: "field", Reason: "invalid url"},
					},
				},
			},
			"found invalid values in field",
		}, {
			"multiple errors in a field",
			fields{
				m: map[FieldName][]*Error{
					"field": {
						{Field: "field", Reason: "invalid schema"},
						{Field: "field", Reason: "invalid host name"},
					},
				},
			},
			"found invalid values in field",
		}, {
			"errors in multiple fields",
			fields{
				m: map[FieldName][]*Error{
					"url": {
						{Field: "foo", Reason: "invalid ulr"},
					},
					"retry": {
						{Field: "bar", Reason: "retry must be greater than 0"},
					},
					"follow_redirect": {
						{Field: "follow_redirect", Reason: "invalid boolean value"},
					},
				},
			},
			"found invalid values in follow_redirect, retry, url",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &errors_{
				m: tt.fields.m,
			}
			assert.EqualError(t, e, tt.want)
		})
	}
}

func TestValidator_Merge(t *testing.T) {
	tests := []struct {
		name       string
		left       map[FieldName][]*Error
		rightField string
		right      error
		want       string
		wantDetail string
	}{
		{
			"merge an errors_ object",
			map[FieldName][]*Error{
				"users.count": {{Field: "users.count", Reason: "count must be greater than 0"}},
			},
			"admin",
			&errors_{
				m: map[FieldName][]*Error{
					"name": {{Field: "name", Reason: "too long"}},
				},
			},
			"found invalid values in admin.name, users.count",
			"found invalid values in admin.name, users.count:\n" +
				"  - admin.name: too long\n" +
				"  - users.count: count must be greater than 0",
		}, {
			"merge general error",
			map[FieldName][]*Error{
				"users.count": {{Field: "users.count", Reason: "count must be greater than 0"}},
			},
			"admin",
			fmt.Errorf("bla bla bla"),
			"found invalid values in admin, users.count",
			"found invalid values in admin, users.count:\n" +
				"  - admin: bla bla bla\n" +
				"  - users.count: count must be greater than 0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &Validator{
				e: errors_{
					m: tt.left,
				},
			}
			v.Merge(tt.rightField, tt.right)
			assert.EqualError(t, v.Error(), tt.want)
			assert.Equal(t, tt.wantDetail, v.Error().(DetailError).DetailError())
		})
	}
}
