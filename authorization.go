// Copyright 2011, Shelby Ramsey.   All rights reserved.
// Use of this code is governed by a BSD license that can be
// found in the LICENSE.txt file.

package sipparser

// Imports from go standard library
import (
	"errors"
	"strings"
)

type Authorization struct {
	Val         string   "val"
	Credentials string   "credentials"
	Params      []*Param "params"
}

func (a *Authorization) GetParam(param string) *Param {
	if a.Params == nil {
		return nil
	}
	for i := range a.Params {
		if a.Params[i].Param == param {
			return a.Params[i]
		}
	}
	return nil
}

func (a *Authorization) parse() error {
	pos := strings.IndexRune(a.Val, ' ')
	if pos == -1 {
		return errors.New("Authorization.parse err: no LWS found.")
	}
	a.Credentials = a.Val[0:pos]
	if len(a.Val)-1 <= pos {
		return errors.New("Authorization.parse err: no digest-resp found.")
	}
	a.Params = make([]*Param, 0)
	parts := strings.Split(a.Val[pos+1:], ",")
	for i := range parts {
		a.Params = append(a.Params, getParam(strings.Replace(parts[i], "\"", "", -1)))
	}
	return nil
}
