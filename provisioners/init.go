package provisioners

import (
	"github.com/pkg/errors"
	"github.com/yuuki0xff/clustertest/models"
)

var Provisioners = map[models.SpecType]Initializer{}

type Initializer func(prefix string, spec models.Spec) models.Provisioner

func New(prefix string, spec models.Spec) (models.Provisioner, error) {
	fn := Provisioners[spec.Type()]
	if fn == nil {
		return nil, errors.Errorf("provisioner does not support SpecType(%s)", spec.Type())
	}
	return fn(prefix, spec), nil
}
