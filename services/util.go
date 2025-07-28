package services

import "fmt"

func convertInterfaceSliceToStringSlice(interfaceSlice []interface{}) ([]string, error) {
	stringSlice := make([]string, len(interfaceSlice))

	for i, v := range interfaceSlice {
		str, ok := v.(string)
		if !ok {
			return nil, fmt.Errorf("element at index %d is not a string", i)
		}
		stringSlice[i] = str
	}

	return stringSlice, nil
}
