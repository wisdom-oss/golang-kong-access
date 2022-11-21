package golangkongaccess

import "errors"

var ErrConnectionNotSetUp = errors.New("the connection to the api gateway was not set up")
var ErrUnexpectedHttpCode = errors.New("an unexpected http code was received")
var ErrPortOutOfRange = errors.New("the port is outside of the supported range")
var ErrEmptyFunctionParameter = errors.New("one or more parameters supplied to this function are empty")

// Function specific errors

var ErrResourceExists = errors.New("the resource you tried to create already exists")
var ErrBadRequest = errors.New("the request made by the package was rejected by the gateway due to formal errors")
var ErrResourceNotCreated = errors.New("the resource you tried to create has not been created by the gateway")
var ErrResourceNotFound = errors.New("the resource you tried to access was not found")
var ErrResourceNotModified = errors.New("the resource you tried to update was not updated successfully")
