// Copyright 2011, Shelby Ramsey.   All rights reserved.
// Use of this code is governed by a BSD license that can be
// found in the LICENSE.txt file.

package sipparser

// Imports from the go standard library
import (
	"testing"
)

func TestVia(t *testing.T) {
	sm := &SipMsg{}
	s := "SIP/2.0/UDP 0.0.0.0:5060;branch=z9hG4bK05B1a4c756d527cb513"
	sm.parseVia(s)
	if sm.Error != nil {
		t.Errorf("[TestVia] Error parsing via.  Received: " + sm.Error.Error())
	}
	if sm.Via[0].Proto != "SIP" {
		t.Errorf("[TestVia] Error parsing via \"SIP/2.0/UDP 0.0.0.0:5060;branch=z9hG4bK05B1a4c756d527cb513\".  sm.Via[0].Proto should be \"SIP\" but received: " + sm.Via[0].Proto)
	}
	if sm.Via[0].Version != "2.0" {
		t.Errorf("[TestVia] Error parsing via \"SIP/2.0/UDP 0.0.0.0:5060;branch=z9hG4bK05B1a4c756d527cb513\".  sm.Via[0].Version should be \"2.0\" but received: " + sm.Via[0].Version)
	}
	if sm.Via[0].Transport != "UDP" {
		t.Errorf("[TestVia] Error parsing via \"SIP/2.0/UDP 0.0.0.0:5060;branch=z9hG4bK05B1a4c756d527cb513\".  sm.Via[0].Transport should be \"UDP\" but received: " + sm.Via[0].Transport)
	}
	if sm.Via[0].SentBy != "0.0.0.0:5060" {
		t.Errorf("[TestVia] Error parsing via \"SIP/2.0/UDP 0.0.0.0:5060;branch=z9hG4bK05B1a4c756d527cb513\".  Sent by should be \"0.0.0.0:5060\" but received: " + sm.Via[0].SentBy + ".")
	}
}
