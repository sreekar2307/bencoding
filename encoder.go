package bencoding

import "fmt"

type encoder struct {
}

func (e *encoder) encodeString(str string) ([]byte, error) {
	return []byte(fmt.Sprintf("%d:%s", len(str), str)), nil
}

func (e *encoder) encodeInt(val int) ([]byte, error) {
	return []byte(fmt.Sprintf("i%de", val)), nil
}

func (e *encoder) encodeList(list []any) ([]byte, error) {
	var res []byte

	res = append(res, 'l')
	for _, val := range list {
		encodedVal, err := e.encode(val)
		if err != nil {
			return nil, err
		}
		res = append(res, encodedVal...)
	}
	res = append(res, 'e')
	return res, nil
}

func (e *encoder) encodeDict(dict map[string]any) ([]byte, error) {
	var res []byte

	res = append(res, 'd')
	for key, val := range dict {
		encodedKey, err := EncodeString(key)
		if err != nil {
			return nil, err
		}
		encodedVal, err := e.encode(val)
		if err != nil {
			return nil, err
		}
		res = append(res, encodedKey...)
		res = append(res, encodedVal...)
	}
	res = append(res, 'e')
	return res, nil
}

func (e *encoder) encode(val any) ([]byte, error) {
	valType, err := e.checkType(val)
	if err != nil {
		return nil, err
	}
	switch valType {
	case String:
		return e.encodeString(val.(string))
	case Int:
		return e.encodeInt(val.(int))
	case List:
		return e.encodeList(val.([]interface{}))
	case Dict:
		return e.encodeDict(val.(map[string]interface{}))
	}
	return nil, err
}

func (e *encoder) checkType(val any) (int, error) {
	switch val.(type) {
	case string:
		return String, nil
	case int:
		return Int, nil
	case []interface{}:
		return List, nil
	case map[string]interface{}:
		return Dict, nil
	default:
		return -1, ErrInvalidFormat
	}
}

func Encode(val any) ([]byte, error) {
	encoder := encoder{}
	return encoder.encode(val)
}

func EncodeString(str string) ([]byte, error) {
	encoder := encoder{}
	return encoder.encodeString(str)
}

func EncodeInt(val int) ([]byte, error) {
	encoder := encoder{}
	return encoder.encodeInt(val)
}

func EncodeList(val []any) ([]byte, error) {
	encoder := encoder{}
	return encoder.encodeList(val)
}
func EncodeDict(val map[string]any) ([]byte, error) {
	encoder := encoder{}
	return encoder.encodeDict(val)
}
