package tftp

const(
	DatagramSize = 516
	BlockSize = DatagramSize - 4 // 4 byte is reserved for the header
)

type OperationCode uint16

const (
	OpRRQ OperationCode = iota + 1
	_
	OpData
	OpAck
	OpErr
)

type ErrorCode uint16

const (
	ErrUnknown ErrorCode = iota
	ErrNotFound
	ErrAccessViolation
	ErrDiskFull
	ErrIllegalOp
	ErrUnknownId
	ErrFileExists
	ErrNoUser
)
