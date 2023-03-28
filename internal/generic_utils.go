package internal

import (
	"encoding/json"
	"fmt"
	"log"
)

func UnmarshalObjs[T any](rawObjs [][]byte, errs []error) ([]*T, []error) {
	objs := make([]*T, len(rawObjs))
	for i, rawObj := range rawObjs {
		if rawObj != nil {
			var unmarshaledObj T
			// log.Println("rawobj equals:", string(rawObj))
			err := json.Unmarshal(rawObj, &unmarshaledObj)
			if err != nil {
				log.Println("error unmarshaling:", err)
				errs[i] = fmt.Errorf("unmarshaling response: %w", err)
			} else {
				objs[i] = &unmarshaledObj
			}
		}
	}
	return objs, errs
}
