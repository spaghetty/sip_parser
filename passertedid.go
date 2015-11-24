// Copyright 2011, Shelby Ramsey.   All rights reserved.
// Use of this code is governed by a BSD license that can be
// found in the LICENSE.txt file.

package sipparser

// Imports from the go standard library
import (
	"errors"
)

// pAssertedIdStateFn is just a fn type 
type pAssertedIdStateFn func(p *PAssertedId) pAssertedIdStateFn

// PAssertedId is a struct that holds:
// -- Error is just an error
// -- Val is the raw value
// -- Name is the name value from the p-asserted-id hdr
// -- URI is the parsed uri from the p-asserted-id hdr
// -- Params is a slice of the *Params from the p-asserted-id hdr
type PAssertedId struct {
	Error   error
	Val     string
	Name    string
	URI     *URI
	Params  []*Param
	nameInt int
}

// addParam adds the *Param to the Params field
func (p *PAssertedId) addParam(s string) {
	np := getParam(s)
	if p.Params == nil {
		p.Params = []*Param{np}
		return
	}
	p.Params = append(p.Params, np)
}

// parse actually parsed the .Val field of the PAssertedId struct
func (p *PAssertedId) parse() {
	for state := parsePAssertedId; state != nil; {
		state = state(p)
	}
}

func parsePAssertedId(p *PAssertedId) pAssertedIdStateFn {
	if p.Error != nil {
		return nil
	}
	p.Name, p.nameInt = getName(p.Val)
	return parsePAssertedIdGetUri
}

func parsePAssertedIdGetUri(p *PAssertedId) pAssertedIdStateFn {
	left := 0
	right := 0
	for i := range p.Val {
		if p.Val[i] == '<' && left == 0 {
			left = i
		}
		if p.Val[i] == '>' && right == 0 {
			right = i
		}
	}
	if left < right {
		p.URI = ParseURI(p.Val[left+1 : right])
		if p.URI.Error != nil {
			p.Error = errors.New("parseRpidGetUri err: received err getting uri: " + p.URI.Error.Error())
			return nil
		}
		return parsePAssertedIdGetParams
	}
	p.Error = errors.New("parseRpidGetUri err: could not locate bracks.  no uri found.")
	return nil
}

func parsePAssertedIdGetParams(p *PAssertedId) pAssertedIdStateFn {
	var pos []int
	right := 0
	for i := range p.Val {
		if p.Val[i] == '>' && right == 0 {
			right = i
		}
	}
	if len(p.Val) > right+1 {
		pos = make([]int, 0)
		for i := range p.Val[right+1:] {
			if p.Val[right+1:][i] == ';' {
				pos = append(pos, i)
			}
		}
	}
	if pos == nil {
		return nil
	}
	for i := range pos {
		if len(pos)-1 == i {
			if len(p.Val[right+1:])-1 > pos[i]+1 {
				p.addParam(p.Val[right+1:][pos[i]+1:])
			}
		}
		if len(pos)-1 > i {
			p.addParam(p.Val[right+1:][pos[i]+1 : pos[i+1]])
		}
	}
	return nil
}
