// Copyright 2011, Shelby Ramsey.   All rights reserved.
// Use of this code is governed by a BSD license that can be
// found in the LICENSE.txt file.

package sipparser

import (
	"testing"
)

func TestWarning(t *testing.T) {
	s := "301 isi.edu \"Incompatible network address type 'E.164'\""
	w := &Warning{Val: s}
	err := w.parse()
	if err != nil {
		t.Errorf("[TestWarning] Error parsing warning.  Got err: " + err.Error())
	}
	if w.Code != "301" {
		t.Errorf("[TestWarning] Error parsing warning.  Code is not \"301\".  Rcvd: " + w.Code)
	}
	if w.Agent != "isi.edu" {
		t.Errorf("[TestWarning] Error parsing warning.  Agent is not \"isi.edu\". Rcvd: " + w.Agent)
	}
	if w.Text != "Incompatible network address type 'E.164'" {
		t.Errorf("[TestWarning] Error parsing warning.  Text should be \"Incompatible network address type 'E.164'\".  Rcvd: " + w.Text)
	}
}
