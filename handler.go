package main

type Handler interface {
	Handle(data interface{}) (interface{}, error)
	SetNext(next Handler)
}
