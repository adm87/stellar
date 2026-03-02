package errs

import "strconv"

// --------------------------------------------------------------------------------
// BootFailure Error
// --------------------------------------------------------------------------------

type BootFailure struct {
	Message string
}

func (e BootFailure) Error() string {
	return "boot failure: " + e.Message
}

// --------------------------------------------------------------------------------
// DuplicateEntry Error
// --------------------------------------------------------------------------------

type DuplicateEntry struct {
	Message string
}

func (e DuplicateEntry) Error() string {
	return "duplicate entry: " + e.Message
}

// --------------------------------------------------------------------------------
// Fatal Error
// --------------------------------------------------------------------------------

type Fatal struct {
	Message string
}

func (e Fatal) Error() string {
	return "fatal error: " + e.Message
}

// --------------------------------------------------------------------------------
// IndexOutOfBounds Error
// --------------------------------------------------------------------------------

type IndexOutOfBounds struct {
	Message string
}

func (e IndexOutOfBounds) Error() string {
	return "index out of bounds: " + e.Message
}

// --------------------------------------------------------------------------------
// InvalidArg Error
// --------------------------------------------------------------------------------

type InvalidArg struct {
	Message string
}

func (e InvalidArg) Error() string {
	return "invalid argument: " + e.Message
}

// --------------------------------------------------------------------------------
// InvalidOperation Error
// --------------------------------------------------------------------------------

type InvalidOperation struct {
	Message string
}

func (e InvalidOperation) Error() string {
	return "invalid operation: " + e.Message
}

// --------------------------------------------------------------------------------
// MaxCapacity Error
// --------------------------------------------------------------------------------

type MaxCapacity struct {
	Message  string
	Capacity int
}

func (e MaxCapacity) Error() string {
	return "max capacity exceeded: " + strconv.Itoa(e.Capacity) + " - " + e.Message
}
