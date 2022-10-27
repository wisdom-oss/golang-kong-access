package golangkongaccess

import "errors"

var ErrConnectionNotSetUp = errors.New("the connection to the api gateway was not set up")
var ErrPortOutOfRange = errors.New("the port is outside of the supported range")
