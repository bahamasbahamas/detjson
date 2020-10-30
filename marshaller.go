package detjson

import (
	"bytes"
	"encoding/json"
	"errors"
	"sort"
	"strconv"
	"strings"
)

// Marshaller type.
type Marshaller struct {
	Map    map[string]interface{}
	jsonString string
}

// NewMarshaller Constructor.
func NewMarshaller(jsonString string) *Marshaller {
	return &Marshaller{
		jsonString: jsonString,
	}
}

// UnMarshal parses the JSON-encoded data and stores the result in the Map
func (m *Marshaller) UnMarshal() error {
	decoder := json.NewDecoder(strings.NewReader(m.GetJSONString()))
	decoder.UseNumber()
	if err := decoder.Decode(&m.Map); err != nil {
		return err
	}
	return nil
}

// MarshalOrdered Marshal returns the ordered JSON encoding of Map .
func (m *Marshaller) MarshalOrdered() error {
	data := m.Map
	str, err := m.marshalRecursiveFunction(&bytes.Buffer{}, data, 0)
	if err != nil {
		return err
	}
	m.SetJSONString(str)
	return nil
}

// sortKeys sorts keys alphabetically and sorts types regarding the deepness of the level.
// even level causes primitive types like strings and integers were pushed to the top
// uneven level causes objects and maps were pushed to the top
func (m *Marshaller) sortKeys(mapToSort map[string]interface{}, level int) []string {
	simpleValuesKeys := make([]string, 0)
	objectsValuesKeys := make([]string, 0)
	for key, value := range mapToSort {
		switch value.(type) {
		case bool:
			simpleValuesKeys = append(simpleValuesKeys, key)
		case int:
			simpleValuesKeys = append(simpleValuesKeys, key)
		case json.Number:
			simpleValuesKeys = append(simpleValuesKeys, key)
		case nil:
			simpleValuesKeys = append(simpleValuesKeys, key)
		case string:
			simpleValuesKeys = append(simpleValuesKeys, key)
		case []interface{}:
			objectsValuesKeys = append(objectsValuesKeys, key)
		case map[string]interface{}:
			objectsValuesKeys = append(objectsValuesKeys, key)
		}
	}
	sort.Strings(simpleValuesKeys)
	sort.Strings(objectsValuesKeys)
	if level%2 == 0 {
		simpleValuesKeys = append(simpleValuesKeys, objectsValuesKeys...)
		return simpleValuesKeys
	}
	objectsValuesKeys = append(objectsValuesKeys, simpleValuesKeys...)
	return objectsValuesKeys
}

// marshalRecursiveFunction recursive help function for MarshalOrdered.
func (m *Marshaller) marshalRecursiveFunction(buf *bytes.Buffer, data map[string]interface{}, level int) (string, error) {
	sortMap := make([]string, 0)
	counter := 0
	sortMap = append(sortMap, m.sortKeys(data, level)...)
	buf.WriteString("{")
	for i := 0; i < len(sortMap); i++ {
		counter++
		buf.WriteString("\"" + sortMap[i] + "\"" + ":")
		value := data[sortMap[i]]
		lenSortMap := len(sortMap)
		switch val := value.(type) {
		case bool:
			boolAsString := strconv.FormatBool(val)
			buf.WriteString(boolAsString)
		case int:
			buf.Write([]byte((strconv.Itoa(val))))
		case json.Number:
			buf.Write([]byte(val))
		case nil:
			buf.WriteString("null")
		case string:
			buf.WriteString("\"" + val + "\"")
		case []interface{}:
			err := m.handleInterfaceArray(val, buf, level)
			if err != nil {
				return "", err
			}
		case map[string]interface{}:
			str, err := m.marshalRecursiveFunction(&bytes.Buffer{}, val, level+1)
			if err != nil {
				return "", err
			}
			buf.WriteString(str)
		default:
			return "", errors.New("wrong type. cannot marshall json ordered")
		}
		m.writeCommaIfNecessary(i, lenSortMap, buf)
	}
	buf.WriteString("}")
	return buf.String(), nil
}

// writeCommaIfNecessary compares the position of the element with the length of the list
// and writes comma to the buf is the element is not the last element in the list
func (m *Marshaller) writeCommaIfNecessary(pos int, length int, buf *bytes.Buffer) {
	if pos < length-1 {
		buf.WriteString(",")
	}
}

// handleInterfaceArray marshal the array Interface ([]interface{}).
func (m *Marshaller) handleInterfaceArray(val []interface{}, buf *bytes.Buffer, level int) error {
	buf.WriteString("[")
	// für jedes Feld im Array, zB. für alle Traveler
	for k := 0; k < len(val); k++ {
		lenVal := len(val)
		switch val := val[k].(type) {
		case bool:
			boolAsString := strconv.FormatBool(val)
			buf.WriteString(boolAsString)
		case int:
			buf.Write([]byte((strconv.Itoa(val))))
		case json.Number:
			buf.Write([]byte(val))
		case nil:
			buf.WriteString("null")
		case string:
			buf.WriteString("\"" + val + "\"")
		case map[string]interface{}:
			str, err := m.marshalRecursiveFunction(&bytes.Buffer{}, val, level+1)
			if err != nil {
				return err
			}
			buf.WriteString(str)
		default:
			return errors.New("cannot marshall interface array")
		}
		m.writeCommaIfNecessary(k, lenVal, buf)
	}
	buf.WriteString("]")
	return nil
}

// GetJSONString.
func (m Marshaller) GetJSONString() string {
	return m.jsonString
}

// SetJSONString.
func (m *Marshaller) SetJSONString(jsonString string) {
	m.jsonString = jsonString
}

// Marshal Marshal.
func (m *Marshaller) Marshal() error {
	jsonBytes, err := json.Marshal(m.Map)
	if err != nil {
		return err
	}
	m.SetJSONString(string(jsonBytes))
	return nil
}