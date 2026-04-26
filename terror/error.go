// Copyright 2020 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package terror

import (
	"errors"
	"fmt"
	"strconv"
)

// ErrCodeText is a textual error code that represents a specific error type
// in an error class.
type ErrCodeText string

// ErrorID is a textual identifier of an Error prototype.
type ErrorID string

// RFCErrorCode is the textual error code in the form {Class}:{Code}.
type RFCErrorCode string

// Error is the prototype of a class of errors. Use Normalize or
// ErrClass.New / ErrClass.NewStd to create one.
type Error struct {
	code     ErrCode
	codeText ErrCodeText
	message  string
	args     []any
}

func (e *Error) Code() ErrCode { return e.code }

func (e *Error) RFCCode() RFCErrorCode { return RFCErrorCode(e.ID()) }

// ID returns the unique identifier of this error.
func (e *Error) ID() ErrorID {
	if e.codeText != "" {
		return ErrorID(e.codeText)
	}
	return ErrorID(strconv.Itoa(int(e.code)))
}

func (e *Error) Error() string {
	if e == nil {
		return "<nil>"
	}
	return fmt.Sprintf("[%s]%s", e.RFCCode(), e.GetMsg())
}

func (e *Error) GetMsg() string {
	if len(e.args) > 0 {
		return fmt.Sprintf(e.message, e.args...)
	}
	return e.message
}

// GenWithStack returns a copy of e with a new message format and arguments.
func (e *Error) GenWithStack(format string, args ...any) error {
	err := *e
	err.message = format
	err.args = args
	return &err
}

// GenWithStackByArgs returns a copy of e with new arguments.
func (e *Error) GenWithStackByArgs(args ...any) error {
	err := *e
	err.args = args
	return &err
}

// FastGen returns a copy of e with a new message format and arguments.
func (e *Error) FastGen(format string, args ...any) error {
	return e.GenWithStack(format, args...)
}

// FastGenByArgs returns a copy of e with new arguments.
func (e *Error) FastGenByArgs(args ...any) error {
	return e.GenWithStackByArgs(args...)
}

// Is reports whether other is an Error with the same ID. Allows usage with
// errors.Is from the standard library.
func (e *Error) Is(other error) bool {
	err, ok := other.(*Error)
	if !ok {
		return false
	}
	return (e == nil && err == nil) || (e != nil && err != nil && e.ID() == err.ID())
}

// Equal reports whether other is or wraps an Error with the same ID as e.
func (e *Error) Equal(other error) bool {
	if other == nil {
		return false
	}
	inner := cause(other)
	if inner == nil {
		return false
	}
	inErr, ok := inner.(*Error)
	if !ok {
		return false
	}
	return e.ID() == inErr.ID()
}

// NotEqual is the inverse of Equal.
func (e *Error) NotEqual(other error) bool { return !e.Equal(other) }

// NormalizeOption configures an Error created via Normalize.
type NormalizeOption func(*Error)

// RFCCodeText sets the textual RFC code of the Error.
func RFCCodeText(s string) NormalizeOption {
	return func(e *Error) { e.codeText = ErrCodeText(s) }
}

// MySQLErrorCode sets the numeric MySQL error code of the Error.
func MySQLErrorCode(code int) NormalizeOption {
	return func(e *Error) { e.code = ErrCode(code) }
}

// Normalize creates a new Error prototype.
func Normalize(message string, opts ...NormalizeOption) *Error {
	e := &Error{message: message}
	for _, opt := range opts {
		opt(e)
	}
	return e
}

// cause walks the error chain via errors.Unwrap and returns the deepest
// non-nil error.
func cause(err error) error {
	for err != nil {
		next := errors.Unwrap(err)
		if next == nil {
			return err
		}
		err = next
	}
	return err
}
