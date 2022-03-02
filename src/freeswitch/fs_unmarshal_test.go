package freeswitch

import (
	"errors"
	"pjsip-handler/src/pjsip"
	"reflect"
	"testing"
)

func TestFreeSwitchToMap(t *testing.T) {
	type ToMapCase struct {
		data string
		expd map[string]string
		err  error
	}

	var testcases = []ToMapCase{
		{`<include>\n<gateway name="asterlink.com"/>\n\t<param name="username" value="cluecon"/>\n\t<param name="realm" value="aster-realm"/>\n\t<param name="from-domain" value="test-domain"/>`,
			map[string]string{"TrunkName": "asterlink.com", "auth/username": "cluecon", "auth/realm": "aster-realm", "endpoint/from_domain": "test-domain"}, nil},
		{`abcdsda -a=asda asdjald gateway name="test2"\n\t<param name="" value="fail"/>\n\t<param name="password" value="1234"/>\n\t`,
			map[string]string{}, errors.New("param name must be non-empty")},
		{`<include>\n<gateway name="asterlink.com"/>\n\t<param name="invalid" value="someval"/>\n\t`,
			map[string]string{}, errors.New("invalid param name")},
		{`<include>\n\t<param name="register" value="false"/>\n\t<param name="caller-id-in-from" value="true"/>`,
			map[string]string{"registration/line": "no", "endpoint/caller_id_in_from": "yes"}, nil},
	}

	for i := range testcases {
		m, err := freeSwitchToMap([]byte(testcases[i].data))
		if err != nil && testcases[i].err == nil {
			t.Fatalf("testcase %d: unexpected error %s", i, err)
		}

		if err == nil && testcases[i].err != nil {
			t.Fatalf("testcase %d: expected error %s, got nil", i, testcases[i].err)
		}

		if !reflect.DeepEqual(m, testcases[i].expd) {
			t.Fatalf("testcase %d: expected result %v, got %v", i, testcases[i].expd, m)
		}
	}
}

func TestFillStructFromXML(t *testing.T) {
	type fillStructCase struct {
		tofill interface{}
		data   map[string]string
		expd   interface{}
		err    bool
	}

	var testcases = []fillStructCase{
		{&pjsip.Endpoint{}, map[string]string{"Context": "cluecon", "OutboundProxy": "someproxy"},
			&pjsip.Endpoint{Context: "cluecon", OutboundProxy: "someproxy"}, false},
		{&pjsip.Endpoint{}, map[string]string{"FromUser": "someuser", "fail_field": "fail"},
			&pjsip.Endpoint{}, true},
		{&pjsip.Registration{}, map[string]string{"Transport": "udp", "RetryInterval": "60"},
			&pjsip.Registration{Transport: []string{"udp"}, RetryInterval: "60"}, false},
	}

	for i, c := range testcases {
		err := fillStructFromXMLMap(c.tofill, c.data)
		if err != nil && c.err == false {
			t.Fatalf("testcase %d unexpected error %s", i, err)
		}

		if err == nil && c.err == true {
			t.Fatalf("testcase %d: expected error, got nil", i)
		}

		if err == nil && !reflect.DeepEqual(c.tofill, c.expd) {
			t.Fatalf("testcase %d: expected %v, got %v", i, c.expd, c.tofill)
		}
	}
}

func TestToPjsip(t *testing.T) {
	type toPjsipCase struct {
		data string
		expd *pjsip.PJSIP
		err  bool
	}

	var testcases = []toPjsipCase{
		{`<include>\n\t<gateway name="test1"/>\n\t<param name="username" value="user1"/>\n\t<param name="from-user" value="cluecon"/>\n\t<param name="password" value="2007"/>\n<include>`,
			&pjsip.PJSIP{Trunks: []pjsip.Trunk{{Name: "test1", Auth: pjsip.Auth{Username: "user1", Password: "2007"}, Endpoint: pjsip.Endpoint{FromUser: "cluecon"}}}}, false},
		{`<include>\n\t<gateway name="test2"/>\n\t<param name="from-domain" value="asterlink.com"/>\n\t<param name="extension" value="cluecon"/>\n\t<param name="register-proxy" value="mysbc.com"/>\n<include>`,
			&pjsip.PJSIP{Trunks: []pjsip.Trunk{{Name: "test2", Endpoint: pjsip.Endpoint{FromDomain: "asterlink.com", Context: "cluecon"}, Registration: pjsip.Registration{OutboundProxy: "mysbc.com"}}}}, false},
		{`<include>\n\t<gateway name="test3"/>\n\t<param name="fail" value="fail"/>\n<include>`,
			&pjsip.PJSIP{}, true},
	}

	for i := range testcases {
		temp := new(pjsip.PJSIP)

		err := ToPJSIP([]byte(testcases[i].data), temp)
		if err != nil && !testcases[i].err {
			t.Fatalf("testcase %d: unexpected error %s", i, err)
		}

		if err == nil && testcases[i].err {
			t.Fatalf("testcase %d: expected error, got nil", i)
		}

		if err == nil && !reflect.DeepEqual(temp, testcases[i].expd) {
			t.Fatalf("testcase %d: expected %v, got %v", i, testcases[i].expd, temp)
		}
	}
}
