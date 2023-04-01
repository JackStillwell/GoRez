package internal

import (
	"encoding/json"
	"fmt"
)

func UnmarshalObjs[T any](rawObjs [][]byte, errs []error) ([]*T, []error) {
	objs := make([]*T, len(rawObjs))
	for i, rawObj := range rawObjs {
		if rawObj != nil {
			var unmarshaledObj T
			err := json.Unmarshal(rawObj, &unmarshaledObj)
			if err != nil {
				errs[i] = fmt.Errorf("unmarshaling response: %w", err)
			} else {
				objs[i] = &unmarshaledObj
			}
		}
	}
	return objs, errs
}
