package map_string_interface

type MapStringInterface[T any] interface {
	MarshalMap() (map[string]interface{}, error)
	UnmarshalMap(map[string]interface{}) error
}
