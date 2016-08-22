package main

type Value struct {
	valueType   string
	stringValue string
	intValue    uint64
	listValue   []Value
	mapValue    map[string]Value
}

func (v Value) StringValue() string {
	return v.stringValue
}

func (v Value) IntValue() uint64 {
	return v.intValue
}

func (v Value) ListValue() []Value {
	return v.listValue
}

func (v Value) MapValue() map[string]Value {
	return v.mapValue
}

func (v Value) ValueType() string {
	return v.valueType
}
