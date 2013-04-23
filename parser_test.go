package sipparser

// Imports from the go standard library
import (
	//"fmt"
	"testing"
)

func TestHeader(t *testing.T) {
    h := Header{"t","v",}
    if h.String() != "t: v" {
        t.Errorf("Error with header.String() method.  Unexpected res.")
    }
}

// makes sure that the body is what is expected
func TestBody(t *testing.T) {
    s := ParseMsg("fake\r\nheader\r\n\r\nbody ...\r\n\r\n")
    if s.Body != "body ...\r\n\r\n" {
        t.Errorf("[TestBody] Error getting the right body from the string.")    
    }
}

// actual msg testing 
func TestParseMsg(t *testing.T) {
    m := "SIP/2.0 200 OK\r\nVia: SIP/2.0/UDP 0.0.0.0:5060;branch=z9hG4bK24477ab511325213INV52e94be64e6687e3;received=0.0.0.0\r\nContact: <sip:10003053258853@0.0.0.0:6060>\r\nTo: <sip:10003053258853@0.0.0.0;user=phone;noa=national>;tag=a94c095b773be1dd6e8d668a785a9c843f6f2cc0\r\nFrom: <sip:8173383772@0.0.0.0;user=phone;noa=national>;tag=52e94be6-co2998-INS002\r\nCall-ID: 111118149-3524331107-398662@barinfo.fooinfous.com\r\nCSeq: 299801 INVITE\r\nAccept: application/sdp, application/dtmf-relay, text/plain\r\nX-Nonsense-Hdr: nonsense\r\nAllow: PRACK, INVITE, BYE, REGISTER, ACK, OPTIONS, CANCEL, SUBSCRIBE, NOTIFY, INFO, REFER, UPDATE\r\nContent-Type: application/sdp\r\nServer: Dialogic-SIP/10.5.3.231 IMGDAL0001 0\r\nSupported: 100rel, path, replaces, timer, tdialog\r\nContent-Length: 239\r\n\r\nv=0\r\no=Dialogic_SDP 1452654 0 IN IP4 0.0.0.0\r\ns=Dialogic-SIP\r\nc=IN IP4 4.71.122.135\r\nt=0 0\r\nm=audio 11676 RTP/AVP 0 101\r\na=rtpmap:0 PCMU/8000\r\na=rtpmap:101 telephone-event/8000\r\na=fmtp:101 0-15\r\na=silenceSupp:off - - - -\r\na=ptime:20\r\n"
    s := ParseMsg(m)
    if s.Error != nil {
	t.Errorf("[TestParseMsg] Error parsing msg. Recevied: " + s.Error.Error())
    }
    if len(s.Body) == 0 {
	t.Errorf("[TestParseMsg] Error parsing msg.  Body should have a length.")
    }
    if len(s.Headers) == 0 {
	t.Errorf("[TestParseMsg] Error parsing msg.  Does not appear to be any headers.")
    }
    if s.Via == nil || len(s.Via) == 0 {
	t.Errorf("[TestParseMsg] Error parsing msg.  Does not appear to be any vias parsed.")
	//fmt.Println("msg:", s.Msg)
	//fmt.Println("body:", s.Body)
	//fmt.Println("via:", s.Via)
	//fmt.Println("crlf:", s.crlf)
    }
    if s.ContentLength != "239" {
	t.Errorf("[TestParseMsg] Error parsing msg.  Content length should be 239.  Received: " + s.ContentLength)
    } 
    if len(s.Supported) != 5 {
	t.Errorf("[TestParseMsg] Error parsing msg.  s.Support should have length of 5.")
    }
}

// testing the GetCallingParty functionality
func TestGetCallingParty(t *testing.T) {
    rpid := "\"UNKNOWN\" <sip:8885551000@0.0.0.0>;party=calling;screen=yes;privacy=off"
    s := &SipMsg{RemotePartyIdVal: rpid,}
    err := s.GetCallingParty(CALLING_PARTY_RPID)
    if err != nil {
	t.Errorf("[TestGetCallingParty] Err with GetCallingParty.  rcvd: " + err.String())
    }
    if s.CallingParty == nil {
	t.Errorf("[TestGetCallingParty] Err calling GetCallingParty.  CallingParty field should not be nil.")
    }
    if s.CallingParty.Name != "UNKNOWN" {
	t.Errorf("[TestGetCallingParty] Err calling GetCallingParty.  Name should be \"UNKNOWN\".")
    }
    if s.CallingParty.Number != "8885551000" {
	t.Errorf("[TestGetCallingParty] Err with GetCallingParty. Number should be \"8885551000\".")
    }
    paid := "<sip:8884441000@0.0.0.0:5060;user=phone>"
    s = &SipMsg{PAssertedIdVal: paid,}
    err = s.GetCallingParty(CALLING_PARTY_PAID)
    if err != nil {
	t.Errorf("[TestGetCallingParty] Err with GetCallingParty on paid.  rcvd: " + err.String())
    }
    if s.CallingParty == nil {
	t.Errorf("[TestGetCallingParty] Err with GetCallingParty on paid.  No CallingPartyInfo.")
    }
    if s.CallingParty.Name != "" {
	t.Errorf("[TestGetCallingParty] Err with GetCallingParty on paid.  Name should be \"\".")
    }
    if s.CallingParty.Number != "8884441000" {
	t.Errorf("[TestGetCallingParty] Err with GetCallingParty on paid.  Number should be \"8884441000\".")
    }
    s = &SipMsg{}
    s.parseFrom("\"5556661000\" <sip:5556661000@0.0.0.0>;tag=ZN21rHN5B7U0K") 
    err = s.GetCallingParty("")
    if err != nil {
	t.Errorf("[TestGetCallingParty] Err calling GetCallingParty on default.  rcvd: " + err.String())
    }
    if s.CallingParty == nil {
	t.Errorf("[TestGetCallingParty] Err calling GetCallingParty on default.  No CallingPartyInfo.")
    }
    if s.CallingParty.Name != "5556661000" {
	t.Errorf("[TestGetCallingParty] Err calling GetCallingParty on default.  Name should be \"5556661000\".")
    }
    if s.CallingParty.Number != "5556661000" {
	t.Errorf("[TestGetCallingParty] Err calling GetCallingParty on default. Number should be \"5556661000\".")
    }
}
