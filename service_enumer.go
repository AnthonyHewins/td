// Code generated by "enumer -type service -json -trimprefix service -transform snake-upper"; DO NOT EDIT.

package td

import (
	"encoding/json"
	"fmt"
	"strings"
)

const _serviceName = "UNSPECIFIEDADMINLEVELONE_EQUITIESLEVELONE_OPTIONSLEVELONE_FUTURESLEVELONE_FUTURES_OPTIONSLEVELONE_FOREXNYSE_BOOKNASDAQ_BOOKOPTIONS_BOOKCHART_EQUITYCHART_FUTURESSCREENER_EQUITYSCREENER_OPTIONACCT_ACTIVITY"

var _serviceIndex = [...]uint8{0, 11, 16, 33, 49, 65, 89, 103, 112, 123, 135, 147, 160, 175, 190, 203}

const _serviceLowerName = "unspecifiedadminlevelone_equitieslevelone_optionslevelone_futureslevelone_futures_optionslevelone_forexnyse_booknasdaq_bookoptions_bookchart_equitychart_futuresscreener_equityscreener_optionacct_activity"

func (i service) String() string {
	if i >= service(len(_serviceIndex)-1) {
		return fmt.Sprintf("service(%d)", i)
	}
	return _serviceName[_serviceIndex[i]:_serviceIndex[i+1]]
}

// An "invalid array index" compiler error signifies that the constant values have changed.
// Re-run the stringer command to generate them again.
func _serviceNoOp() {
	var x [1]struct{}
	_ = x[serviceUnspecified-(0)]
	_ = x[serviceAdmin-(1)]
	_ = x[serviceLeveloneEquities-(2)]
	_ = x[serviceLeveloneOptions-(3)]
	_ = x[serviceLeveloneFutures-(4)]
	_ = x[serviceLeveloneFuturesOptions-(5)]
	_ = x[serviceLeveloneForex-(6)]
	_ = x[serviceNyseBook-(7)]
	_ = x[serviceNasdaqBook-(8)]
	_ = x[serviceOptionsBook-(9)]
	_ = x[serviceChartEquity-(10)]
	_ = x[serviceChartFutures-(11)]
	_ = x[serviceScreenerEquity-(12)]
	_ = x[serviceScreenerOption-(13)]
	_ = x[serviceAcctActivity-(14)]
}

var _serviceValues = []service{serviceUnspecified, serviceAdmin, serviceLeveloneEquities, serviceLeveloneOptions, serviceLeveloneFutures, serviceLeveloneFuturesOptions, serviceLeveloneForex, serviceNyseBook, serviceNasdaqBook, serviceOptionsBook, serviceChartEquity, serviceChartFutures, serviceScreenerEquity, serviceScreenerOption, serviceAcctActivity}

var _serviceNameToValueMap = map[string]service{
	_serviceName[0:11]:         serviceUnspecified,
	_serviceLowerName[0:11]:    serviceUnspecified,
	_serviceName[11:16]:        serviceAdmin,
	_serviceLowerName[11:16]:   serviceAdmin,
	_serviceName[16:33]:        serviceLeveloneEquities,
	_serviceLowerName[16:33]:   serviceLeveloneEquities,
	_serviceName[33:49]:        serviceLeveloneOptions,
	_serviceLowerName[33:49]:   serviceLeveloneOptions,
	_serviceName[49:65]:        serviceLeveloneFutures,
	_serviceLowerName[49:65]:   serviceLeveloneFutures,
	_serviceName[65:89]:        serviceLeveloneFuturesOptions,
	_serviceLowerName[65:89]:   serviceLeveloneFuturesOptions,
	_serviceName[89:103]:       serviceLeveloneForex,
	_serviceLowerName[89:103]:  serviceLeveloneForex,
	_serviceName[103:112]:      serviceNyseBook,
	_serviceLowerName[103:112]: serviceNyseBook,
	_serviceName[112:123]:      serviceNasdaqBook,
	_serviceLowerName[112:123]: serviceNasdaqBook,
	_serviceName[123:135]:      serviceOptionsBook,
	_serviceLowerName[123:135]: serviceOptionsBook,
	_serviceName[135:147]:      serviceChartEquity,
	_serviceLowerName[135:147]: serviceChartEquity,
	_serviceName[147:160]:      serviceChartFutures,
	_serviceLowerName[147:160]: serviceChartFutures,
	_serviceName[160:175]:      serviceScreenerEquity,
	_serviceLowerName[160:175]: serviceScreenerEquity,
	_serviceName[175:190]:      serviceScreenerOption,
	_serviceLowerName[175:190]: serviceScreenerOption,
	_serviceName[190:203]:      serviceAcctActivity,
	_serviceLowerName[190:203]: serviceAcctActivity,
}

var _serviceNames = []string{
	_serviceName[0:11],
	_serviceName[11:16],
	_serviceName[16:33],
	_serviceName[33:49],
	_serviceName[49:65],
	_serviceName[65:89],
	_serviceName[89:103],
	_serviceName[103:112],
	_serviceName[112:123],
	_serviceName[123:135],
	_serviceName[135:147],
	_serviceName[147:160],
	_serviceName[160:175],
	_serviceName[175:190],
	_serviceName[190:203],
}

// serviceString retrieves an enum value from the enum constants string name.
// Throws an error if the param is not part of the enum.
func serviceString(s string) (service, error) {
	if val, ok := _serviceNameToValueMap[s]; ok {
		return val, nil
	}

	if val, ok := _serviceNameToValueMap[strings.ToLower(s)]; ok {
		return val, nil
	}
	return 0, fmt.Errorf("%s does not belong to service values", s)
}

// serviceValues returns all values of the enum
func serviceValues() []service {
	return _serviceValues
}

// serviceStrings returns a slice of all String values of the enum
func serviceStrings() []string {
	strs := make([]string, len(_serviceNames))
	copy(strs, _serviceNames)
	return strs
}

// IsAservice returns "true" if the value is listed in the enum definition. "false" otherwise
func (i service) IsAservice() bool {
	for _, v := range _serviceValues {
		if i == v {
			return true
		}
	}
	return false
}

// MarshalJSON implements the json.Marshaler interface for service
func (i service) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.String())
}

// UnmarshalJSON implements the json.Unmarshaler interface for service
func (i *service) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return fmt.Errorf("service should be a string, got %s", data)
	}

	var err error
	*i, err = serviceString(s)
	return err
}
