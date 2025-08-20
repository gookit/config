package properties_test

import (
	"testing"

	"github.com/gookit/goutil/testutil/assert"
	"github.com/nstrlabs/config"
	"github.com/nstrlabs/config/properties"
)

func TestDriver(t *testing.T) {
	is := assert.New(t)
	is.Eq(properties.Name, properties.Driver.Name())

	c := config.NewEmpty("test")
	is.False(c.HasDecoder(properties.Name))
	c.AddDriver(properties.Driver)

	is.True(c.HasDecoder(properties.Name))
	is.True(c.HasEncoder(properties.Name))

	m := struct {
		N string
	}{}
	err := properties.Decoder([]byte(`
// comments
n=value
	`), &m)

	is.Nil(err)
	is.Eq("value", m.N)
}
