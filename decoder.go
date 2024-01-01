package bencoding

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
)

const (
	Int = iota
	String
	List
	Dict
)

type decoder struct {
	reader *bufio.Reader
}

func newDecoder(r io.Reader) *decoder {
	if reader, ok := r.(*bufio.Reader); ok {
		return &decoder{reader}
	}
	return &decoder{bufio.NewReader(r)}
}

func (d *decoder) decodeString(firstByte byte) (string, error) {
	var (
		strLen   = []byte{firstByte}
		strBytes []byte
	)

	for {
		b, err := d.readByte()
		if err != nil {
			return "", err
		}
		if b == ':' {
			break
		}
		strLen = append(strLen, b)
	}

	strLenAsInt, err := strconv.Atoi(string(strLen))
	if err != nil {
		return "", err
	}

	for i := 0; i < strLenAsInt; i++ {
		b, err := d.readByte()
		if err != nil {
			return "", err
		}
		strBytes = append(strBytes, b)
	}
	return string(strBytes), nil
}

func (d *decoder) decodeInt() (int, error) {
	var (
		decodedIntBytes []byte
		err             error
		prevByte        byte
		nextByte        byte
	)

	// "ie" is not valid

	{
		nextByte, err = d.readByte()
		if err != nil {
			return 0, err
		}
		if nextByte == 'e' {
			return 0, errors.Join(fmt.Errorf("ie"), ErrInvalidFormat)
		}
		decodedIntBytes = append(decodedIntBytes, nextByte)
	}

	// "i-0..|i-e|i0{0-9}.." is not valid

	{

		prevByte = nextByte
		nextByte, err = d.readByte()
		if err != nil {
			return 0, err
		}

		if nextByte == '0' && prevByte == '-' {
			return 0, errors.Join(fmt.Errorf("i-0"), ErrInvalidFormat)
		}

		if nextByte == 'e' && prevByte == '-' {
			return 0, errors.Join(fmt.Errorf("i-e"), ErrInvalidFormat)
		}

		if nextByte >= '0' && nextByte <= '9' && prevByte == '0' {
			return 0, errors.Join(fmt.Errorf("i0\\d+"), ErrInvalidFormat)
		}
	}

	if nextByte == 'e' {
		return strconv.Atoi(string(decodedIntBytes))
	}

	decodedIntBytes = append(decodedIntBytes, nextByte)

	for {
		nextByte, err = d.readByte()
		if err != nil {
			return 0, err
		}
		if nextByte == 'e' {
			break
		}
		decodedIntBytes = append(decodedIntBytes, nextByte)
	}

	return strconv.Atoi(string(decodedIntBytes))
}

func (d *decoder) decodeList() ([]any, error) {
	var list []any

	for {
		nextByte, err := d.readByte()
		if err != nil {
			return nil, err
		}
		if nextByte == 'e' {
			break
		}
		el, err := d.decode(nextByte)
		if err != nil {
			return nil, err
		}
		list = append(list, el)
	}
	return list, nil
}

func (d *decoder) decodeDict() (map[string]any, error) {
	var (
		nextByte byte
		err      error
		dict     = make(map[string]any)
	)
	for {
		// form key
		nextByte, err = d.readByte()
		if err != nil {
			return nil, err
		}
		if nextByte == 'e' {
			break
		}
		keyType, err := d.checkType(nextByte)
		if err != nil {
			return nil, err
		}
		if keyType != String {
			return nil, ErrInvalidFormat
		}
		key, err := d.decode(nextByte)
		if err != nil {
			return nil, err
		}

		// form value
		nextByte, err = d.readByte()
		if err != nil {
			return nil, err
		}
		if nextByte == 'e' {
			break
		}
		value, err := d.decode(nextByte)
		if err != nil {
			return nil, err
		}

		dict[key.(string)] = value
	}
	return dict, nil
}

func (d *decoder) readByte() (byte, error) {
	return d.reader.ReadByte()
}

func (d *decoder) checkType(b byte) (int, error) {
	switch {
	case b >= '0' && b <= '9':
		return String, nil
	case b == 'l':
		return List, nil
	case b == 'i':
		return Int, nil
	case b == 'd':
		return Dict, nil
	}

	return 0, ErrInvalidFormat

}

func (d *decoder) decode(firstByte byte) (any, error) {
	valType, err := d.checkType(firstByte)
	if err != nil {
		return nil, err
	}
	switch valType {
	case Int:
		return d.decodeInt()
	case String:
		return d.decodeString(firstByte)
	case List:
		return d.decodeList()
	case Dict:
		return d.decodeDict()
	}
	return nil, ErrInvalidFormat
}

// Decode reads either list, dict, string, int from the reader
func Decode(r io.Reader) (any, error) {
	decoder := newDecoder(r)
	readByte, err := decoder.readByte()
	if err != nil {
		return nil, err
	}
	return decoder.decode(readByte)
}

// DecodeString reads a string from the reader
func DecodeString(r io.Reader) (string, error) {
	decoder := newDecoder(r)
	readByte, err := decoder.readByte()
	if err != nil {
		return "", err
	}
	resultType, err := decoder.checkType(readByte)
	if err != nil {
		return "", err
	}
	if resultType != String {
		return "", ErrInvalidFormat
	}

	return decoder.decodeString(readByte)
}

// DecodeInt reads a int from the reader
func DecodeInt(r io.Reader) (int, error) {
	decoder := newDecoder(r)
	readByte, err := decoder.readByte()
	if err != nil {
		return 0, err
	}
	resultType, err := decoder.checkType(readByte)
	if err != nil {
		return 0, err
	}
	if resultType != Int {
		return 0, ErrInvalidFormat
	}

	return decoder.decodeInt()
}

// DecodeList reads a list from the reader
func DecodeList(r io.Reader) ([]any, error) {
	decoder := newDecoder(r)
	readByte, err := decoder.readByte()
	if err != nil {
		return nil, err
	}
	resultType, err := decoder.checkType(readByte)
	if err != nil {
		return nil, err
	}
	if resultType != List {
		return nil, ErrInvalidFormat
	}

	return decoder.decodeList()
}

// DecodeDict reads a dict from the reader
func DecodeDict(r io.Reader) (map[string]any, error) {
	decoder := newDecoder(r)
	readByte, err := decoder.readByte()
	if err != nil {
		return nil, err
	}
	resultType, err := decoder.checkType(readByte)
	if err != nil {
		return nil, err
	}
	if resultType != Dict {
		return nil, ErrInvalidFormat
	}

	return decoder.decodeDict()
}
