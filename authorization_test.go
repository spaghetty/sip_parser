package sipparser

// Imports from go standard library
import (
	"testing"
)

func TestAuthorization(t *testing.T) {
	val := "Digest username=\"foobaruser124\", realm=\"FOOBAR\", algorithm=MD5, uri=\"sip:foo.bar.com\", nonce=\"4f6d7a1d\", response=\"6a79a5c75572b0f6a18963ae04e971cf\", opaque=\"\""
	a := &Authorization{Val: val}
	err := a.parse()
	if err != nil {
		t.Errorf("[TestAuthorization] Err parsing authorization hdr.  Received: " + err.Error())
	}
	if a.Credentials != "Digest" {
		t.Errorf("[TestAuthorization] Err parsing authorization hdr.  Credentials should be \"Digest\" but rcvd: " + a.Credentials)
	}
	if a.GetParam("realm").Val != "FOOBAR" {
		t.Errorf("[TestAuthorization] Err parsing authorization hdr.  Called a.GetParam(\"realm\") and did not get \"FOOBAR\".  rcvd: " + a.GetParam("realm").Val)
	}
}
