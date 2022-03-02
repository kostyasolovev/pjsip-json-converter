package freeswitch

import (
	"pjsip-handler/src/pjsip"
	"reflect"
	"sort"
	"testing"
)

type pjsipAsFs struct {
	trunk interface{}
	expd  []string
}

func TestReadPjsipAsFS(t *testing.T) {
	var testcases = []pjsipAsFs{
		{&pjsip.Trunk{Name: "test1", Endpoint: pjsip.Endpoint{FromUser: "someuser", Allow: []string{"all"}}, Auth: pjsip.Auth{Username: "somename", Md5Cred: "md5test"}},
			[]string{"gateway:test1", "from-user:someuser", "username:somename"}},
		{&pjsip.Endpoint{FromDomain: "somedomain", Allow: []string{"all"}},
			[]string{"from-domain:somedomain"}},
		{&pjsip.Trunk{Name: "test3", Registration: pjsip.Registration{OutboundProxy: "someproxy", Line: "yes"}, Endpoint: pjsip.Endpoint{CallerIdInFrom: "no"}},
			[]string{"gateway:test3", "register-proxy:someproxy", "register:true", "caller-id-in-from:false"}},
		{&pjsip.Trunk{Name: "test4", Registration: pjsip.Registration{Transport: []string{"udp"}}, Endpoint: pjsip.Endpoint{FromUser: "testuser4"}, Aor: pjsip.Aor{ExtraContactParams: "info", QualifyFrequency: "100"}},
			[]string{"gateway:test4", "register-transport:udp", "from-user:testuser4", "contact-params:info", "ping:100"}},
	}

	for i := range testcases {
		out := make([]string, 0, 20)
		ReadPjsipAsFS(testcases[i].trunk, "", &out)

		if !pjsip.CompareSlices(out, testcases[i].expd) {
			t.Fatalf("testcase %d: expected %v, got %v", i, testcases[i].expd, out)
		}
	}
}

type unmToIncCase struct {
	data []string
	expd *Include
	err  bool
}

func TestUnmarshalToInclude(t *testing.T) {
	var testcases = []unmToIncCase{
		{[]string{"gateway:trunk1", "username:arnold", "from-domain:someval", "password:1234qwer!@#$"},
			&Include{Gateway: Gateway{Name: "trunk1"}, Param: []Param{{Name: "username", Value: "arnold"}, {Name: "from-domain", Value: "someval"}, {Name: "password", Value: "1234qwer!@#$"}}}, false},
		{[]string{"gateway:trunk2", "extension:cluecon", ":fail"},
			&Include{}, true},
	}

	for i := range testcases {
		inc := new(Include)

		err := unmarshalToInclude(testcases[i].data, inc)
		if iserr := err != nil; iserr != testcases[i].err {
			t.Fatalf("testcase %d wrong error behavior: expected %t, got %t", i, testcases[i].err, iserr)
		}

		if !testcases[i].err {
			if !reflect.DeepEqual(inc, testcases[i].expd) {
				t.Fatalf("testcase %d: expected %v, got %v", i, testcases[i].expd, inc)
			}
		}
	}
}

type pjsipToFSCase struct {
	data *pjsip.PJSIP
	expd Include
	err  bool
}

func TestPjsipToFS(t *testing.T) {
	var testcases = []pjsipToFSCase{
		{&pjsip.PJSIP{Trunks: []pjsip.Trunk{{Name: "mytrunk1", Endpoint: pjsip.Endpoint{FromUser: "spiderman", DirectMedia: "g850"}, Aor: pjsip.Aor{ExtraContactParams: "cluecon", PingMin: "5"}}}},
			Include{Gateway: Gateway{"mytrunk1"}, Param: []Param{{Name: "from-user", Value: "spiderman"}, {"contact-params", "cluecon"}, {"ping-min", "5"}}}, false},
		{&pjsip.PJSIP{Trunks: []pjsip.Trunk{{Name: "mytrunk2", Aor: pjsip.Aor{TestFailField: "100500"}}}},
			Include{}, true},
		{&pjsip.PJSIP{Trunks: []pjsip.Trunk{{Name: "mytrunk3", Auth: pjsip.Auth{Username: "alien", Realm: "www.example.org"}, Registration: pjsip.Registration{OutboundProxy: "someproxy", Line: "yes", Transport: []string{"udp"}}}}},
			Include{Gateway: Gateway{"mytrunk3"}, Param: []Param{{"username", "alien"}, {"realm", "www.example.org"}, {"register-proxy", "someproxy"}, {"register", "true"}, {"register-transport", "udp"}}}, false},
	}

	for i := range testcases {
		out, err := PjsipToFS(testcases[i].data)
		if iserr := err != nil; iserr != testcases[i].err {
			t.Fatalf("testcase %d wrong error behavior: expected %t, got %t", i, testcases[i].err, iserr)
		}
		// порядок []Param в тесткейсах и в создаваемых функцией слайсах разный, поэтому нужно их сортировать прежде чем сравнивать
		if err == nil && !reflect.DeepEqual(out[0].Gateway, testcases[i].expd.Gateway) {
			got, expd := testParams{}, testParams{}

			copy(got, out[0].Param)
			copy(expd, testcases[i].expd.Param)
			sort.Sort(got)
			sort.Sort(expd)

			if !reflect.DeepEqual(got, expd) {
				t.Fatalf("testcase %d: expected %v, got %v", i, testcases[i].expd, out[0])
			}
		}
	}
}

type testParams []Param

// implementing Sort interface.
func (prms testParams) Len() int           { return len(prms) }
func (prms testParams) Swap(i, j int)      { prms[i], prms[j] = prms[j], prms[i] }
func (prms testParams) Less(i, j int) bool { return prms[i].Name < prms[j].Name }
