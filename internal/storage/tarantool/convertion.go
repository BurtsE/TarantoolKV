package tarantool

import (
	"TarantoolKV/internal/application/core/domain"
	"fmt"
)

func convertToDomain(value [][]interface{}) domain.Entity {
	entity := domain.Entity{}
	key := value[0][0].(string)
	m := value[0][1].(map[interface{}]interface{})
	entity.Key = key
	entity.Value = convertMap(m)
	return entity
}

func convertMap(m map[interface{}]interface{}) map[string]interface{} {
	res := make(map[string]interface{})
	for k, v := range m {
		key := fmt.Sprintf("%v", k)
		switch val := v.(type) {
		case map[interface{}]interface{}:
			res[key] = convertMap(val)
		case []interface{}:
			res[key] = convertSlice(val)
		default:
			res[key] = val
		}
	}
	return res
}

func convertSlice(s []interface{}) []interface{} {
	res := make([]interface{}, len(s))
	for i, v := range s {
		switch val := v.(type) {
		case map[interface{}]interface{}:
			res[i] = convertMap(val)
		case []interface{}:
			res[i] = convertSlice(val)
		default:
			res[i] = val
		}
	}
	return res
}
