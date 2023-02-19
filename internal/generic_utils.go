package internal

import (
	"encoding/json"
	"fmt"
)

func UnmarshalObjs[T any](rawObjs [][]byte, errs []error) ([]*T, []error) {
	objs := make([]*T, len(rawObjs))
	for i, rawObj := range rawObjs {
		if rawObj != nil {
			obj := objs[i]
			err := json.Unmarshal(rawObj, obj)
			if err != nil {
				errs[i] = fmt.Errorf("unmarshaling response: %w", err)
			}
		}
	}
	return objs, errs
}
