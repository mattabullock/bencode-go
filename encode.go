package bencode

import (
	"bytes"
	"strconv"
)

func Encode(enc interface{}) []byte {
	bytes := []byte{}
	switch enc.(type) {
	case string:
		bytes = EncodeString(enc.(string))
	case int:
		bytes = EncodeInt(enc.(int))
	case map[string]interface{}:
		bytes = EncodeMap(enc.(map[string]interface{}))
	}

	return bytes
}

func EncodeString(s string) []byte {
	encString := strconv.Itoa(len(s)) + ":" + s

	return []byte(encString)
}

func EncodeInt(i int) []byte {
	encInt := "i" + strconv.Itoa(i) + "e"

	return []byte(encInt)
}

func EncodeMap(m map[string]interface{}) []byte {
	var buffer bytes.Buffer

	buffer.WriteString("d")

	for k, v := range m {
		buffer.Write(EncodeString(k))
		switch v.(type) {
		case string:
			buffer.Write(EncodeString(v.(string)))
		case int:
			buffer.Write(EncodeInt(v.(int)))
		case map[string]interface{}:
			buffer.Write(EncodeMap(v.(map[string]interface{})))
		case []interface{}:
			buffer.Write(EncodeList(v.([]interface{})))
		}
	}

	buffer.WriteString("e")

	return buffer.Bytes()
}

func EncodeList(l []interface{}) []byte {
	var buffer bytes.Buffer

	buffer.WriteString("l")

	for _, v := range l {
		switch v.(type) {
		case string:
			buffer.Write(EncodeString(v.(string)))
		case int:
			buffer.Write(EncodeInt(v.(int)))
		case map[string]interface{}:
			buffer.Write(EncodeMap(v.(map[string]interface{})))
		case []interface{}:
			buffer.Write(EncodeList(v.([]interface{})))
		}
	}

	buffer.WriteString("e")

	return buffer.Bytes()
}
