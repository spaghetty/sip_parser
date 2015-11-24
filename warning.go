// Copyright 2011, Shelby Ramsey.   All rights reserved.
// Use of this code is governed by a BSD license that can be
// found in the LICENSE.txt file.

package sipparser

// Imports from the go standard library
import (
	"errors"
	"strings"
)

type Warning struct {
	Val   string
	Code  string
	Agent string
	Text  string
}

func (w *Warning) parse() error {
	parts := strings.SplitN(w.Val, " ", 3)
	if len(parts) != 3 {
		return errors.New("Warning.parse err: split on LWS was not correct.")
	}
	w.Code = parts[0]
	w.Agent = parts[1]
	w.Text = strings.Replace(parts[2], "\"", "", -1)
	return nil
}
