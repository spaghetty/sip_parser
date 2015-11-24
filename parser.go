package sipparser

import (
	"errors"
	"fmt"
	"strings"
)

const (
	sipParseStateStartLine    = "SipParseStateStartLine"
	sipParseStateCrlf         = "SipMsgStateCrlf"
	sipParseStateBody         = "SipMsgStateBody"
	sipParseStateHeaders      = "SipMsgStateHeaders"
	sipParseStateParseHeaders = "SipMsgStateParseHeaders"
	uintCR                    = '\r'
	uintLF                    = '\n'
	CR                        = "\r"
	LF                        = "\n"
	CALLING_PARTY_DEFAULT     = "default"
	CALLING_PARTY_RPID        = "rpid"
	CALLING_PARTY_PAID        = "paid"
)

type CallingPartyInfo struct {
	Name      string
	Number    string
	Anonymous bool
}

type Header struct {
	Header string
	Val    string
}

func (h *Header) String() string {
	return fmt.Sprintf("%s: %s", h.Header, h.Val)
}

type sipParserStateFn func(s *SipMsg) sipParserStateFn

type SipMsg struct {
	State              string
	Error              error
	Msg                string
	CallingParty       *CallingPartyInfo
	Body               string
	StartLine          *StartLine
	Headers            []*Header
	Accept             *Accept
	AlertInfo          string
	Allow              []string
	AllowEvents        []string
	Authorization      *Authorization
	ContentDisposition *ContentDisposition
	ContentLength      string
	ContentLengthInt   int
	ContentType        string
	From               *From
	MaxForwards        string
	MaxForwardsInt     int
	Organization       string
	To                 *From
	Contact            *From
	ContactVal         string
	CallId             string
	Cseq               *Cseq
	Rack               *Rack
	Reason             *Reason
	Rseq               string
	RseqInt            int
	RecordRoute        []*URI
	Route              []*URI
	Via                []*Via
	Require            []string
	Supported          []string
	Privacy            string
	ProxyAuthenticate  *Authorization
	ProxyRequire       []string
	RemotePartyIdVal   string
	RemotePartyId      *RemotePartyId
	PAssertedIdVal     string
	PAssertedId        *PAssertedId
	Unsupported        []string
	UserAgent          string
	Server             string
	Subject            string
	Warning            *Warning
	WWWAuthenticate    *Authorization
	eof                int
	hdr                string
	hdrv               string
}

func (s *SipMsg) run() {
	s.Headers = make([]*Header, 0)
	for state := parseSip; state != nil; {
		state = state(s)
	}
}

func (s *SipMsg) addError(err string) sipParserStateFn {
	s.Error = errors.New(err)
	return nil
}

func (s *SipMsg) addErrorNoReturn(err string) {
	s.Error = errors.New(err)
}

func (s *SipMsg) addHdr(str string) {
	if str == "" {
		return
	}
	sp := strings.IndexRune(str, ':')
	if sp == -1 {
		s.Error = errors.New("addHdr err: no semi found in: " + str)
		return
	}
	s.hdr = strings.ToLower(strings.TrimSpace(str[0:sp]))
	switch {
	case len(str)-1 > sp+1:
		s.hdrv = cleanWs(str[sp+1:])
	default:
		s.hdrv = ""
	}
	switch {
	case s.hdr == SIP_HDR_ACCEPT:
		s.parseAccept(s.hdrv)
	case s.hdr == SIP_HDR_ALLOW:
		s.parseAllow(s.hdrv)
	case s.hdr == SIP_HDR_ALLOW_EVENTS || s.hdr == SIP_HDR_ALLOW_EVENTS_CMP:
		s.parseAllowEvents(s.hdrv)
	case s.hdr == SIP_HDR_AUTHORIZATION:
		s.parseAuthorization(s.hdrv)
	case s.hdr == SIP_HDR_CALL_ID || s.hdr == SIP_HDR_CALL_ID_CMP:
		s.CallId = s.hdrv
	case s.hdr == SIP_HDR_CONTACT || s.hdr == SIP_HDR_CONTACT_CMP:
		s.ContactVal = s.hdrv
	case s.hdr == SIP_HDR_CONTENT_DISPOSITION:
		s.parseContentDisposition(s.hdrv)
	case s.hdr == SIP_HDR_CONTENT_LENGTH || s.hdr == SIP_HDR_CONTENT_LENGTH_CMP:
		s.ContentLength = s.hdrv
	case s.hdr == SIP_HDR_CSEQ:
		s.parseCseq(s.hdrv)
	case s.hdr == SIP_HDR_FROM || s.hdr == SIP_HDR_FROM_CMP:
		s.parseFrom(s.hdrv)
	case s.hdr == SIP_HDR_MAX_FORWARDS:
		s.MaxForwards = s.hdrv
	case s.hdr == SIP_HDR_CONTENT_TYPE:
		s.ContentType = s.hdrv		
	case s.hdr == SIP_HDR_ORGANIZATION:
		s.Organization = s.hdrv
	case s.hdr == SIP_HDR_P_ASSERTED_IDENTITY:
		s.PAssertedIdVal = s.hdrv
	case s.hdr == SIP_HDR_PRIVACY:
		s.Privacy = s.hdrv
	case s.hdr == SIP_HDR_PROXY_AUTHENTICATE:
		s.parseProxyAuthenticate(s.hdrv)
	case s.hdr == SIP_HDR_RACK:
		s.parseRack(s.hdrv)
	case s.hdr == SIP_HDR_REASON:
		s.parseReason(s.hdrv)
	case s.hdr == SIP_HDR_RECORD_ROUTE:
		s.parseRecordRoute(s.hdrv)
	case s.hdr == SIP_HDR_REMOTE_PARTY_ID:
		s.RemotePartyIdVal = s.hdrv
	case s.hdr == SIP_HDR_ROUTE:
		s.parseRoute(s.hdrv)
	case s.hdr == SIP_HDR_SERVER:
		s.Server = s.hdrv
	case s.hdr == SIP_HDR_SUPPORTED:
		s.parseSupported(s.hdrv)
	case s.hdr == SIP_HDR_TO || s.hdr == SIP_HDR_TO_CMP:
		s.parseTo(s.hdrv)
	case s.hdr == SIP_HDR_UNSUPPORTED:
		s.parseUnsupported(s.hdrv)
	case s.hdr == SIP_HDR_USER_AGENT:
		s.UserAgent = s.hdrv
	case s.hdr == SIP_HDR_VIA || s.hdr == SIP_HDR_VIA_CMP:
		s.parseVia(s.hdrv)
	case s.hdr == SIP_HDR_WARNING:
		s.parseWarning(s.hdrv)
	case s.hdr == SIP_HDR_WWW_AUTHENTICATE:
		s.parseWWWAuthenticate(s.hdrv)
	default:
		s.Headers = append(s.Headers, &Header{s.hdr, s.hdrv})
	}
}

func (s *SipMsg) GetRURIParamBool(str string) bool {
	if s.StartLine == nil || s.StartLine.URI == nil {
		return false
	}
	for i := range s.StartLine.URI.UriParams {
		if s.StartLine.URI.UriParams[i].Param == str {
			return true
		}
	}
	return false
}

func (s *SipMsg) GetRURIParamVal(str string) string {
	if s.StartLine == nil || s.StartLine.URI == nil {
		return ""
	}
	for i := range s.StartLine.URI.UriParams {
		if s.StartLine.URI.UriParams[i].Param == str {
			return s.StartLine.URI.UriParams[i].Val
		}
	}
	return ""
}

func (s *SipMsg) GetCallingParty(str string) error {
	switch {
	case str == CALLING_PARTY_RPID:
		return s.getCallingPartyRpid()
	case str == CALLING_PARTY_PAID:
		return s.getCallingPartyPaid()
	default:
		return s.getCallingPartyDefault()
	}
	return errors.New("GetCallingParty err: unknown err.  Fell through switch stmnt.")
}

func (s *SipMsg) getCallingPartyDefault() error {
	if s.From == nil {
		return errors.New("getCallingPartyDefault err: no from header found.")
	}
	if s.From.URI == nil {
		return errors.New("getCallingPartyDefault err: no uri found in from header.")
	}
	s.CallingParty = &CallingPartyInfo{Name: s.From.Name, Number: s.From.URI.User}
	return nil
}

func (s *SipMsg) getCallingPartyPaid() error {
	if s.PAssertedId == nil {
		if s.PAssertedIdVal == "" {
			return s.getCallingPartyDefault()
		}
		s.parsePAssertedId(s.PAssertedIdVal)
		if s.Error != nil {
			return s.Error
		}
		if s.PAssertedId.URI == nil {
			return errors.New("getCallingPartyPaid err: p-asserted-id uri is nil.")
		}
		s.CallingParty = &CallingPartyInfo{Name: s.PAssertedId.Name, Number: s.PAssertedId.URI.User}
		return nil
	}
	if s.PAssertedId.URI == nil {
		return errors.New("getCallingPartyPaid err: p-asserted-id uri is nil.")
	}
	s.CallingParty = &CallingPartyInfo{Name: s.PAssertedId.Name, Number: s.PAssertedId.URI.User}
	return nil
}

func (s *SipMsg) getCallingPartyRpid() error {
	if s.RemotePartyId == nil {
		if s.RemotePartyIdVal == "" {
			return s.getCallingPartyDefault()
		}
		s.parseRemotePartyId(s.RemotePartyIdVal)
		if s.Error != nil {
			return s.Error
		}
		if s.RemotePartyId.URI == nil {
			return errors.New("getCallingPartyRpid err: remote party id uri is nil.")
		}
		s.CallingParty = &CallingPartyInfo{Name: s.RemotePartyId.Name, Number: s.RemotePartyId.URI.User}
		return nil
	}
	if s.RemotePartyId.URI == nil {
		return errors.New("getCallingPartyRpid err: remote party id uri is nil.")
	}
	s.CallingParty = &CallingPartyInfo{Name: s.RemotePartyId.Name, Number: s.RemotePartyId.URI.User}
	return nil
}

func (s *SipMsg) parseAccept(str string) {
	s.Accept = &Accept{Val: str}
	s.Accept.parse()
}

func (s *SipMsg) parseAllow(str string) {
	s.Allow = getCommaSeperated(str)
	if s.Allow == nil {
		s.Allow = []string{str}
	}
}

func (s *SipMsg) parseAllowEvents(str string) {
	s.AllowEvents = getCommaSeperated(str)
	if s.AllowEvents == nil {
		s.Allow = []string{str}
	}
}

func (s *SipMsg) parseAuthorization(str string) {
	s.Authorization = &Authorization{Val: str}
	s.Error = s.Authorization.parse()
}

func (s *SipMsg) parseContact(str string) {
	s.Contact = getFrom(str)
	if s.Contact.Error != nil {
		s.Error = s.Contact.Error
	}
}

func (s *SipMsg) ParseContact(str string) {
	s.parseContact(str)
}

func (s *SipMsg) parseContentDisposition(str string) {
	s.ContentDisposition = &ContentDisposition{Val: str}
	s.ContentDisposition.parse()
}

func (s *SipMsg) parseCseq(str string) {
	s.Cseq = &Cseq{Val: str}
	s.Error = s.Cseq.parse()
}

func (s *SipMsg) parseFrom(str string) {
	s.From = getFrom(str)
	if s.From.Error != nil {
		s.Error = s.From.Error
	}
}

func (s *SipMsg) parsePAssertedId(str string) {
	s.PAssertedId = &PAssertedId{Val: str}
	s.PAssertedId.parse()
	if s.PAssertedId.Error != nil {
		s.Error = s.PAssertedId.Error
	}
}

func (s *SipMsg) ParsePAssertedId(str string) {
	s.parsePAssertedId(str)
}

func (s *SipMsg) parseProxyAuthenticate(str string) {
	s.ProxyAuthenticate = &Authorization{Val: str}
	s.Error = s.ProxyAuthenticate.parse()
}

func (s *SipMsg) parseRack(str string) {
	s.Rack = &Rack{Val: str}
	s.Error = s.Rack.parse()
}

func (s *SipMsg) parseReason(str string) {
	s.Reason = &Reason{Val: str}
	s.Reason.parse()
}

func (s *SipMsg) parseRecordRoute(str string) {
	cs := getCommaSeperated(str)
	for rt := range cs {
		left := 0
		right := 0
		for i := range cs[rt] {
			if cs[rt][i] == '<' && left == 0 {
				left = i
			}
			if cs[rt][i] == '>' && right == 0 {
				right = i
			}
		}
		if left < right {
			u := ParseURI(cs[rt][left+1 : right])
			if u.Error != nil {
				s.Error = errors.New("parseRecordRoute err: received err parsing uri: " + u.Error.Error())
				return
			}
			if s.RecordRoute == nil {
				s.RecordRoute = []*URI{u}
			}
			s.RecordRoute = append(s.RecordRoute, u)
		}
	}
	return
}

func (s *SipMsg) parseRemotePartyId(str string) {
	s.RemotePartyId = &RemotePartyId{Val: str}
	s.RemotePartyId.parse()
	if s.RemotePartyId.Error != nil {
		s.Error = s.RemotePartyId.Error
	}
}

func (s *SipMsg) ParseRemotePartyId(str string) {
	s.parseRemotePartyId(str)
}

func (s *SipMsg) parseRequire(str string) {
	s.Require = getCommaSeperated(str)
	if s.Require == nil {
		s.Require = []string{str}
	}
}

func (s *SipMsg) parseRoute(str string) {
	cs := getCommaSeperated(str)
	for rt := range cs {
		left := 0
		right := 0
		for i := range cs[rt] {
			if cs[rt][i] == '<' && left == 0 {
				left = i
			}
			if cs[rt][i] == '>' && right == 0 {
				right = i
			}
		}
		if left < right {
			u := ParseURI(cs[rt][left+1 : right])
			if u.Error != nil {
				s.Error = errors.New("parseRoute err: received err parsing uri: " + u.Error.Error())
				return
			}
			if s.Route == nil {
				s.Route = []*URI{u}
			}
			s.Route = append(s.Route, u)
		}
	}
}

func (s *SipMsg) parseStartLine(str string) {
	s.State = sipParseStateStartLine
	s.StartLine = ParseStartLine(str)
	if s.StartLine.Error != nil {
		s.Error = errors.New("parseStartLine err: received err while parsing start line: " + s.StartLine.Error.Error())
	}
}

func (s *SipMsg) parseSupported(str string) {
	s.Supported = getCommaSeperated(str)
	if s.Supported == nil {
		s.Supported = []string{str}
	}
}

func (s *SipMsg) parseTo(str string) {
	s.To = getFrom(str)
	if s.To.Error != nil {
		s.Error = s.To.Error
	}
}

func (s *SipMsg) parseUnsupported(str string) {
	s.Unsupported = getCommaSeperated(str)
	if s.Unsupported == nil {
		s.Unsupported = []string{str}
	}
}

func (s *SipMsg) parseVia(str string) {
	v := &Via{Via: str}
	v.parse()
	if v.Error != nil {
		s.Error = v.Error
		return
	}
	if s.Via == nil {
		s.Via = []*Via{v}
		return
	}
	s.Via = append(s.Via, v)
}

func (s *SipMsg) parseWarning(str string) {
	s.Warning = &Warning{Val: str}
	s.Error = s.Warning.parse()
}

func (s *SipMsg) parseWWWAuthenticate(str string) {
	s.WWWAuthenticate = &Authorization{Val: str}
	s.Error = s.WWWAuthenticate.parse()
}

func getBody(s *SipMsg) sipParserStateFn {
	s.State = sipParseStateBody
	if len(s.Msg)-1 > s.eof+4 {
		s.Body = s.Msg[s.eof+4:]
	}
	return getHeaders
}

func getHeaders(s *SipMsg) sipParserStateFn {
	s.State = sipParseStateHeaders
	var lasth string
	hdrs := strings.Split(s.Msg[0:s.eof], "\r\n")
	for i := range hdrs {
		switch {
		case i == 0:
			s.parseStartLine(hdrs[0])
			if s.Error != nil {
				return nil
			}
		case i == 1:
			lasth = hdrs[i]
		case i > 1:
			if len(hdrs[i]) > 4 {
				switch {
				case hdrs[i][0] == '\t':
					lasth = lasth + hdrs[i][1:]
				case hdrs[i][0:4] == "    ":
					lasth = lasth + hdrs[i][4:]
				default:
					s.addHdr(lasth)
					if s.Error != nil {
						return nil
					}
					lasth = hdrs[i]
				}
			}
			if len(hdrs[i]) < 4 {
				s.addHdr(lasth)
				lasth = hdrs[i]
			}
		}
	}
	s.addHdr(lasth)
	return nil
}

func ParseMsg(str string) (s *SipMsg) {
	s = &SipMsg{Msg: str, eof: strings.Index(str, "\r\n\r\n")}
	if s.eof == -1 {
		s.Error = errors.New("ParseMsg: err parsing msg.  No SIP eof found.")
		return s
	}
	s.run()
	return s
}

func parseSip(s *SipMsg) sipParserStateFn {
	if s.Error != nil {
		return nil
	}
	return getBody
}
