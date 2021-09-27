package golib

func Dump(v interface{}) {
	println(JsonString(v))
}
