// Code generated by "enumer -type UserPrincipalField -trimprefix UserPrincipalField -transform title-lower"; DO NOT EDIT.

package td

import (
	"fmt"
	"strings"
)

const _UserPrincipalFieldName = "unspecifiedstreamerSubscriptionKeysstreamerConnectionInfopreferencessurrogateIds"

var _UserPrincipalFieldIndex = [...]uint8{0, 11, 35, 57, 68, 80}

const _UserPrincipalFieldLowerName = "unspecifiedstreamersubscriptionkeysstreamerconnectioninfopreferencessurrogateids"

func (i UserPrincipalField) String() string {
	if i >= UserPrincipalField(len(_UserPrincipalFieldIndex)-1) {
		return fmt.Sprintf("UserPrincipalField(%d)", i)
	}
	return _UserPrincipalFieldName[_UserPrincipalFieldIndex[i]:_UserPrincipalFieldIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _UserPrincipalFieldNoOp() {
	var x [1]struct{}
	_ = x[UserPrincipalFieldUnspecified-(0)]
	_ = x[UserPrincipalFieldStreamerSubscriptionKeys-(1)]
	_ = x[UserPrincipalFieldStreamerConnectionInfo-(2)]
	_ = x[UserPrincipalFieldPreferences-(3)]
	_ = x[UserPrincipalFieldSurrogateIds-(4)]
}

var _UserPrincipalFieldValues = []UserPrincipalField{UserPrincipalFieldUnspecified, UserPrincipalFieldStreamerSubscriptionKeys, UserPrincipalFieldStreamerConnectionInfo, UserPrincipalFieldPreferences, UserPrincipalFieldSurrogateIds}

var _UserPrincipalFieldNameToValueMap = map[string]UserPrincipalField{
	_UserPrincipalFieldName[0:11]:       UserPrincipalFieldUnspecified,
	_UserPrincipalFieldLowerName[0:11]:  UserPrincipalFieldUnspecified,
	_UserPrincipalFieldName[11:35]:      UserPrincipalFieldStreamerSubscriptionKeys,
	_UserPrincipalFieldLowerName[11:35]: UserPrincipalFieldStreamerSubscriptionKeys,
	_UserPrincipalFieldName[35:57]:      UserPrincipalFieldStreamerConnectionInfo,
	_UserPrincipalFieldLowerName[35:57]: UserPrincipalFieldStreamerConnectionInfo,
	_UserPrincipalFieldName[57:68]:      UserPrincipalFieldPreferences,
	_UserPrincipalFieldLowerName[57:68]: UserPrincipalFieldPreferences,
	_UserPrincipalFieldName[68:80]:      UserPrincipalFieldSurrogateIds,
	_UserPrincipalFieldLowerName[68:80]: UserPrincipalFieldSurrogateIds,
}

var _UserPrincipalFieldNames = []string{
	_UserPrincipalFieldName[0:11],
	_UserPrincipalFieldName[11:35],
	_UserPrincipalFieldName[35:57],
	_UserPrincipalFieldName[57:68],
	_UserPrincipalFieldName[68:80],
}

// UserPrincipalFieldString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func UserPrincipalFieldString(s string) (UserPrincipalField, error) {
	if val, ok := _UserPrincipalFieldNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _UserPrincipalFieldNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to UserPrincipalField values", s)
}

// UserPrincipalFieldValues returns all values of the enum
func UserPrincipalFieldValues() []UserPrincipalField {
	return _UserPrincipalFieldValues
}

// UserPrincipalFieldStrings returns a slice of all String values of the enum
func UserPrincipalFieldStrings() []string {
	strs := make([]string, len(_UserPrincipalFieldNames))
	copy(strs, _UserPrincipalFieldNames)
	return strs
}

// IsAUserPrincipalField returns "true" if the value is listed in the enum definition. "false" otherwise
func (i UserPrincipalField) IsAUserPrincipalField() bool {
	for _, v := range _UserPrincipalFieldValues {
		if i == v {
			return true
		}
	}
	return false
}
