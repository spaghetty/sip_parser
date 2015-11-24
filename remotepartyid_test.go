// Copyright 2011, Shelby Ramsey.   All rights reserved.
// Use of this code is governed by a BSD license that can be
// found in the LICENSE.txt file.

package sipparser

// Imports from the go standard library
import (
	"testing"
)

func TestRpid(t *testing.T) {
	sm := &SipMsg{}
	s := "\"Unknown\" <sip:5558887777@0.0.0.0>;party=calling;screen=yes;privacy=off"
	sm.parseRemotePartyId(s)
	if sm.Error != nil {
		t.Errorf("[TestRpid] Error parsing rpid hdr: \"Unknown\" <sip:5558887777@0.0.0.0>;party=calling;screen=yes;privacy=off.  Received err: " + sm.Error.Error())
	}
	if sm.RemotePartyId.Name != "Unknown" {
		t.Errorf("[TestRpid] Error parsing rpid hdr: \"Unknown\" <sip:5558887777@0.0.0.0>;party=calling;screen=yes;privacy=off.  Name should be \"Unknown\".")
	}
	if sm.RemotePartyId.URI == nil {
		t.Errorf("[TestRpid Error parsing rpid hdr: \"Unknown\" <sip:5558887777@0.0.0.0>;party=calling;screen=yes;privacy=off.  sm.RemotePartyId.URI is nil.")
	}
}
