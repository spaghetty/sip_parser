// Copyright 2011, Shelby Ramsey.   All rights reserved.
// Use of this code is governed by a BSD license that can be
// found in the LICENSE.txt file.

package sipparser

// Imports from the go standard library
import (
	"errors"
	"strings"
)

// Cseq is a struct that holds the values for a cseq header:
//  -- Val is the raw string value of the cseq hdr
//  -- Method is the SIP method
//  -- Digit is the numeric indicator for the method
type Cseq struct {
	Val    string
	Method string
	Digit  string
}

func (c *Cseq) parse() error {
	if c.Val == "" {
		return errors.New("Cseq.parse err: val can not be blank.")
	}
	s := strings.IndexRune(c.Val, ' ')
	if s == -1 {
		return errors.New("Cseq.parse err: lws err with: " + c.Val)
	}
	if s == 0 {
		return errors.New("Cseq.parse err: lws at pos 0 in val: " + c.Val)
	}
	if len(c.Val)-1 < s+1 {
		return errors.New("Cseq.parse err: first lws is end of line in val: " + c.Val)
	}
	c.Method = c.Val[s+1:]
	c.Digit = c.Val[0:s]
	return nil
}
