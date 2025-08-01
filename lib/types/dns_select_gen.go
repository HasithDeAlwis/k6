// Code generated by "enumer -type=DNSSelect -trimprefix DNS -output dns_select_gen.go"; DO NOT EDIT.

package types

import (
	"fmt"
)

const _DNSSelectName = "firstroundRobinrandom"

var _DNSSelectIndex = [...]uint8{0, 5, 15, 21}

func (i DNSSelect) String() string {
	i -= 1
	if i >= DNSSelect(len(_DNSSelectIndex)-1) {
		return fmt.Sprintf("DNSSelect(%d)", i+1)
	}
	return _DNSSelectName[_DNSSelectIndex[i]:_DNSSelectIndex[i+1]]
}

var _DNSSelectValues = []DNSSelect{1, 2, 3}

var _DNSSelectNameToValueMap = map[string]DNSSelect{
	_DNSSelectName[0:5]:   1,
	_DNSSelectName[5:15]:  2,
	_DNSSelectName[15:21]: 3,
}

// DNSSelectString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func DNSSelectString(s string) (DNSSelect, error) {
	if val, ok := _DNSSelectNameToValueMap[s]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to DNSSelect values", s)
}

// DNSSelectValues returns all values of the enum
func DNSSelectValues() []DNSSelect {
	return _DNSSelectValues
}

// IsADNSSelect returns "true" if the value is listed in the enum definition. "false" otherwise
func (i DNSSelect) IsADNSSelect() bool {
	for _, v := range _DNSSelectValues {
		if i == v {
			return true
		}
	}
	return false
}
