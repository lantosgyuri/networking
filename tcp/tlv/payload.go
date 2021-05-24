package tlv

import (
	"errors"
	"fmt"
	"io"
)

const (
	BinaryType uint8 = iota + 1
	StringType

	MayPayloadSize uint32 = 10 << 20
)

var ErrMayPayloadExceeded = errors.New("maximum payload size exceeded")

// Payload is an interface for messages which uses the type-length-value scheme
type Payload interface {
	fmt.Stringer
	io.ReaderFrom
	io.WriterTo
	Bytes() []byte
}
