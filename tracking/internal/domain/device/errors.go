package domainDevice

import "fmt"

type ErrDeviceNotFound struct {
	ID string
}

func (e *ErrDeviceNotFound) Error() string {
	return fmt.Sprintf("device with id %s not found", e.ID)
}

func IsErrDeviceNotFound(err error) bool {
	_, ok := err.(*ErrDeviceNotFound)
	return ok
}
