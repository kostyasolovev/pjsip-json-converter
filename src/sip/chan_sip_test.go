package sip

import (
	"pjsip-handler/src/pjsip"
	"reflect"
	"testing"
)

func TestSipSetNat(t *testing.T) {
	type natCase struct {
		input pjsip.Trunk
		expd  string
	}

	var testcases = []natCase{
		{pjsip.Trunk{Endpoint: pjsip.Endpoint{RtpSymmetric: "yes", ForceRport: "yes", RewriteContact: "yes"}}, "nat=yes"},
		{pjsip.Trunk{}, ""},
		{pjsip.Trunk{Endpoint: pjsip.Endpoint{RtpSymmetric: "yes", RewriteContact: "yes"}}, "nat=yes"},
		{pjsip.Trunk{Endpoint: pjsip.Endpoint{RtpSymmetric: "no", ForceRport: "no", RewriteContact: "no"}}, "nat=no"},
		{pjsip.Trunk{Endpoint: pjsip.Endpoint{RewriteContact: "yes"}}, "nat=route"},
	}

	for i := range testcases {
		if nat := sipSetNat(reflect.ValueOf(testcases[i].input)); nat != testcases[i].expd {
			t.Fatalf("testcase %d: expected %s, got %s", i, testcases[i].expd, nat)
		}
	}
}

func TestSipSetHost(t *testing.T) {
	type hostCase struct {
		input pjsip.Trunk
		expd  string
	}

	var testcases = []hostCase{
		{pjsip.Trunk{Aor: pjsip.Aor{Contact: []string{"sip:6002@192.0.2.1:5060"}}}, "host=192.0.2.1\n"},
		{pjsip.Trunk{}, "host=dynamic\n"},
	}

	for i := range testcases {
		if res := sipSetHost(reflect.ValueOf(testcases[i].input)); res != testcases[i].expd {
			t.Fatalf("testcase %d: expected %s, got %s", i, testcases[i].expd, res)
		}
	}
}
func TestReadPjsipAsSIP(t *testing.T) {
	type readCase struct {
		pjsip interface{}
		expd  []string
	}

	var testcases = []readCase{
		{&pjsip.PJSIP{Transports: []pjsip.Transport{{Protocol: "udp", Bind: "1.1.1.1", Password: "1234"}}},
			[]string{"[general]\n", "udpbindaddr=1.1.1.1\n"}},
		{pjsip.Trunk{Name: "6000", Endpoint: pjsip.Endpoint{Context: "from-external", Disallow: []string{"all"}, Allow: []string{"ulaw"}, UsePtime: "yes"}},
			[]string{"[6000]\n", "type=friend\n", "host=dynamic\n", "context=from-external\n", "disallow=all\n", "allow=ulaw\n"}},
		{pjsip.Trunk{Name: "6004", Endpoint: pjsip.Endpoint{RtpSymmetric: "yes", RewriteContact: "yes"}},
			[]string{"[6004]\n", "type=friend\n", "host=dynamic\n", "nat=yes\n"}},
		{&pjsip.Trunk{Name: "6015", Endpoint: pjsip.Endpoint{ForceRport: "no"}},
			[]string{"[6015]\n", "type=friend\n", "host=dynamic\n", "nat=no\n"}},
		{pjsip.Trunk{Name: "6005", Endpoint: pjsip.Endpoint{RewriteContact: "yes"}},
			[]string{"[6005]\n", "type=friend\n", "host=dynamic\n", "nat=route\n"}},
		{pjsip.Trunk{Name: "6006", Endpoint: pjsip.Endpoint{DirectMedia: "cluecon", Context: "from-external"},
			Registration: pjsip.Registration{Type: "registration", OutboundAuth: "6006", ServerUri: "sip:myaccountname@100.0.0.1:5060", ClientUri: "sip:myaccountname@100.0.0.1:5060"},
			Auth:         pjsip.Auth{Type: "auth", AuthType: "userpass", Password: "1234567890", Username: "myaccountname"},
			Aor:          pjsip.Aor{Type: "aor", Contact: []string{"sip:6002@100.0.0.1:5060"}},
			Identify:     pjsip.Identify{Type: "identify", Endpoint: "6006", Match: []string{"100.0.0.1"}}},
			[]string{"register => myaccountname:1234567890@100.0.0.1:5060\n", "[6006]\n", "type=friend\n", "host=100.0.0.1\n", "directmedia=cluecon\n", "context=from-external\n",
				"secret=1234567890\n", "username=myaccountname\n"}},
	}

	for i := range testcases {
		res := make([]string, 0, 20)

		ReadPjsipAsSIP(testcases[i].pjsip, "", &res)

		if !pjsip.CompareSlices(res, testcases[i].expd) {
			t.Fatalf("testcase %d: expected %v, got %v", i, testcases[i].expd, res)
		}
	}
}
