package helpers

import (
	"fmt"
)

func WriteNestedJson(nestFields []string, value interface{}) string {
	if len(nestFields) == 0 {
		return `{}`
	}
	opener := ``
	closer := ``
	for index, field := range nestFields {
		if index == len(nestFields)-1 {
			opener = opener + fmt.Sprintf(`{"%v":`, field)
			_, ok := value.(string)
			if ok {
				closer = fmt.Sprintf(`"%v"}`, value) + closer

			} else if value == nil {
				closer = fmt.Sprintf(`%v}`, "null") + closer
			} else {
				closer = fmt.Sprintf(`%v}`, value) + closer
			}
			continue
		}
		opener = opener + fmt.Sprintf(`{"%v":`, field)
		closer = `}` + closer
	}
	return opener + closer
}
