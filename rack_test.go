// Copyright 2011, Shelby Ramsey.   All rights reserved.
// Use of this code is governed by a BSD license that can be
// found in the LICENSE.txt file.

package sipparser

// Imports from the go standard library
import (
	"testing"
)

func TestRack(t *testing.T) {
	sm := &SipMsg{}
	s := "776656 1 INVITE"
	sm.parseRack(s)
	if sm.Error != nil {
		t.Errorf("[TestRack] Error parsing rack hdr: 776656 1 INVITE.  Received err: " + sm.Error.Error())
	}
	if sm.Rack.RseqVal != "776656" {
		t.Errorf("[TestRack] Error parsing rack hdr: 776656 1 INVITE.  RseqVal should be 776656 but received: ", sm.Rack.RseqVal)
	}
	if sm.Rack.CseqVal != "1" {
		t.Errorf("[TestRack] Error parsing rack hdr: 776656 1 INVITE.  CseqVal should be 1 but received: ", sm.Rack.CseqVal)
	}
	if sm.Rack.CseqMethod != "INVITE" {
		t.Errorf("[TestRack] Error parsing rack hdr: 776656 1 INVITE.  CseqMethod should be \"INVITE\" but received:", sm.Rack.CseqMethod)
	}
}
