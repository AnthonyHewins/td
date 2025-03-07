// Code generated by "enumer -type command -json -trimprefix command -transform upper"; DO NOT EDIT.

package td

import (
	"encoding/json"
	"fmt"
	"strings"
)

const _commandName = "UNSPECIFIEDLOGINSUBSADDUNSUBSVIEWLOGOUT"

var _commandIndex = [...]uint8{0, 11, 16, 20, 23, 29, 33, 39}

const _commandLowerName = "unspecifiedloginsubsaddunsubsviewlogout"

func (i command) String() string {
	if i >= command(len(_commandIndex)-1) {
		return fmt.Sprintf("command(%d)", i)
	}
	return _commandName[_commandIndex[i]:_commandIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _commandNoOp() {
	var x [1]struct{}
	_ = x[commandUnspecified-(0)]
	_ = x[commandLogin-(1)]
	_ = x[commandSubs-(2)]
	_ = x[commandAdd-(3)]
	_ = x[commandUnsubs-(4)]
	_ = x[commandView-(5)]
	_ = x[commandLogout-(6)]
}

var _commandValues = []command{commandUnspecified, commandLogin, commandSubs, commandAdd, commandUnsubs, commandView, commandLogout}

var _commandNameToValueMap = map[string]command{
	_commandName[0:11]:       commandUnspecified,
	_commandLowerName[0:11]:  commandUnspecified,
	_commandName[11:16]:      commandLogin,
	_commandLowerName[11:16]: commandLogin,
	_commandName[16:20]:      commandSubs,
	_commandLowerName[16:20]: commandSubs,
	_commandName[20:23]:      commandAdd,
	_commandLowerName[20:23]: commandAdd,
	_commandName[23:29]:      commandUnsubs,
	_commandLowerName[23:29]: commandUnsubs,
	_commandName[29:33]:      commandView,
	_commandLowerName[29:33]: commandView,
	_commandName[33:39]:      commandLogout,
	_commandLowerName[33:39]: commandLogout,
}

var _commandNames = []string{
	_commandName[0:11],
	_commandName[11:16],
	_commandName[16:20],
	_commandName[20:23],
	_commandName[23:29],
	_commandName[29:33],
	_commandName[33:39],
}

// commandString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func commandString(s string) (command, error) {
	if val, ok := _commandNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _commandNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to command values", s)
}

// commandValues returns all values of the enum
func commandValues() []command {
	return _commandValues
}

// commandStrings returns a slice of all String values of the enum
func commandStrings() []string {
	strs := make([]string, len(_commandNames))
	copy(strs, _commandNames)
	return strs
}

// IsAcommand returns "true" if the value is listed in the enum definition. "false" otherwise
func (i command) IsAcommand() bool {
	for _, v := range _commandValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for command
func (i command) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for command
func (i *command) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("command should be a string, got %s", data)
	}

	var err error
	*i, err = commandString(s)
	return err
}
