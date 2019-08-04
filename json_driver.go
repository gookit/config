package config

// default json driver(encoder/decoder)
import (
	"encoding/json"

	"github.com/gookit/goutil/jsonutil"
)

// JSONAllowComments support write comments on json file.
var JSONAllowComments = true

// JSONDecoder for json decode
var JSONDecoder Decoder = func(data []byte, v interface{}) (err error) {
	if JSONAllowComments {
		str := jsonutil.StripComments(string(data))
		return json.Unmarshal([]byte(str), v)
	}

	return json.Unmarshal(data, v)
}

// JSONEncoder for json encode
var JSONEncoder Encoder = json.Marshal

// JSONDriver instance fot json
var JSONDriver = &jsonDriver{name: JSON}

// jsonDriver for json format content
type jsonDriver struct {
	name          string
	ClearComments bool
}

// Name of the driver
func (d *jsonDriver) Name() string {
	return d.name
}

// GetDecoder for the driver
func (d *jsonDriver) GetDecoder() Decoder {
	return JSONDecoder
}

// GetEncoder for the driver
func (d *jsonDriver) GetEncoder() Encoder {
	return JSONEncoder
}
