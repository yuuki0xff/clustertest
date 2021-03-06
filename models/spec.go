package models

import (
	"fmt"
)

type SpecType string

// Spec represents an infrastructure specification of clustered system.
// The implementations of Spec interface includes Provisioner specific data.
type Spec interface {
	fmt.Stringer
	Type() SpecType
	LoadDefault(default_ Spec) error
	Validate() error
	ValidateDefault() error
}

// InfraConfig represents current infrastructure configuration.
type InfraConfig interface {
	fmt.Stringer
	Spec() Spec
}
