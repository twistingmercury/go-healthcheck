// Code generated by go-enum
// DO NOT EDIT!

package healthcheck

import (
	"fmt"
)

const (
	// HealthStatusNotSet is a HealthStatus of type NotSet
	HealthStatusNotSet HealthStatus = iota
	// HealthStatusOK is a HealthStatus of type OK
	HealthStatusOK
	// HealthStatusWarning is a HealthStatus of type Warning
	HealthStatusWarning
	// HealthStatusCritical is a HealthStatus of type Critical
	HealthStatusCritical
)

const _HealthStatusName = "NotSetOKWarningCritical"

var _HealthStatusMap = map[HealthStatus]string{
	0: _HealthStatusName[0:6],
	1: _HealthStatusName[6:8],
	2: _HealthStatusName[8:15],
	3: _HealthStatusName[15:23],
}

// String implements the Stringer interface.
func (x HealthStatus) String() string {
	if str, ok := _HealthStatusMap[x]; ok {
		return str
	}
	return fmt.Sprintf("HealthStatus(%d)", x)
}

var _HealthStatusValue = map[string]HealthStatus{
	_HealthStatusName[0:6]:   0,
	_HealthStatusName[6:8]:   1,
	_HealthStatusName[8:15]:  2,
	_HealthStatusName[15:23]: 3,
}

// ParseHealthStatus attempts to convert a string to a HealthStatus
func ParseHealthStatus(name string) (HealthStatus, error) {
	if x, ok := _HealthStatusValue[name]; ok {
		return x, nil
	}
	return HealthStatus(0), fmt.Errorf("%s is not a valid HealthStatus", name)
}

// MarshalText implements the text marshaller method
func (x HealthStatus) MarshalText() ([]byte, error) {
	return []byte(x.String()), nil
}

// UnmarshalText implements the text unmarshaller method
func (x *HealthStatus) UnmarshalText(text []byte) error {
	name := string(text)
	tmp, err := ParseHealthStatus(name)
	if err != nil {
		return err
	}
	*x = tmp
	return nil
}