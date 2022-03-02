package pjsip

import (
	"errors"
	"reflect"
	"testing"
)

type mapPair struct {
	orig string
	expd map[string]string
	err  error
}

func TestToMap(t *testing.T) {
	var mcases = []mapPair{
		{"[mytrunk](template)\ntype=aor\ncontact=sip:198.51.100.1:5060\ncontact=sip:198.51.100.2:5060\n\n",
			map[string]string{"Name": "mytrunk template",
				"Type":    "aor",
				"Contact": "sip:198.51.100.1:5060 sip:198.51.100.2:5060",
			},
			nil},
		{"[mytrunk2]\ntype=identify\nsometext",
			map[string]string{},
			errors.New("invalid line")},
		{"[6001]\ntype=endpoint\ncontext=false\nfax_detect=true\n\n[",
			map[string]string{"Name": "6001",
				"Type":      "endpoint",
				"Context":   "false",
				"FaxDetect": "true"}, nil},
		{"fail string\n",
			map[string]string{},
			errors.New("invalid line")},
	}

	for _, v := range mcases {
		res, err := toMap(v.orig)
		if err == nil && v.err != nil {
			t.Errorf("incorrect error behavior: expected %s, got nil, testcase: %v", v.err, v.orig)
		}

		if err != nil && v.err == nil {
			t.Errorf("unexpected error %s, expected nil", err)
		}

		if !reflect.DeepEqual(res, v.expd) {
			t.Errorf("expected %v, got %v", v.expd, res)
		}
	}
}

func TestToMapUnknown(t *testing.T) {
	var ucases = []mapPair{
		{"[mytrunk](template)\ntype=aor\ncontact=sip:198.51.100.1:5060\ncontact=sip:198.51.100.2:5060\n\n",
			map[string]string{"Name": "mytrunk template",
				"Type":  "aor",
				"LINES": "contact=sip:198.51.100.1:5060 contact=sip:198.51.100.2:5060",
			},
			nil},
		{"[mytrunk2]\ntype=identify\nfailtext",
			map[string]string{},
			errors.New("invalid line")},
		{"[6001]\ntype=endpoint\ncontext=false\nfax_detect=true\n\n[",
			map[string]string{"Name": "6001",
				"Type":  "endpoint",
				"LINES": "context=false fax_detect=true"}, nil},
		{"failstring\n",
			map[string]string{},
			errors.New("invalid block")},
	}

	for _, c := range ucases {
		m, err := toMapUnknown(c.orig)
		if err == nil && c.err != nil {
			t.Errorf("incorrect error behavior: expected %s, got nil, testcase: %v", c.err, c.orig)
		}

		if err != nil && c.err == nil {
			t.Errorf("unexpected error %s, expected nil", err)
		}

		if !reflect.DeepEqual(m, c.expd) {
			t.Errorf("expected %v, got %v", c.expd, m)
		}
	}
}

func TestFillStruct(t *testing.T) {
	type fillpair struct {
		orig  string
		iface interface{}
		exp   interface{}
		err   error
	}

	var filltests = []fillpair{
		{"[mytrunk](template)\ntype=aor\ncontact=sip:198.51.100.1:5060\ncontact=sip:198.51.100.2:5060\n\n", &Aor{},
			&Aor{Name: []string{"mytrunk", "template"}, Type: "aor", Contact: []string{"sip:198.51.100.1:5060", "sip:198.51.100.2:5060"}}, nil},
		{"[6001]\ntype=endpoint\ncontext=false\nfax_detect=true\n\n[", &Endpoint{},
			&Endpoint{Name: []string{"6001"}, Type: "endpoint", Context: "false", FaxDetect: "true"}, nil},
		{"[mytrunk2]\ntype=identify\nfail_field=true\n", &Identify{}, nil, errors.New("struct doesnt have such a field")},
		{"[simple-trans]\ntype=transport\nprotocol=udp\nbind=0.0.0.0", Transport{}, nil, errors.New("struct must be a pointer")},
	}

	for i, v := range filltests {
		err := fillStruct(v.orig, v.iface)
		if err != nil && v.err == nil {
			t.Errorf("unexpected error: %s, testcase %d", err, i)
		}

		if err == nil && v.err != nil {
			t.Errorf("testcase %d, incorrect error behavior: expected %s, got nil", i, err)
		}

		if err == nil && !reflect.DeepEqual(v.iface, v.exp) {
			t.Errorf("testcase %d, expected %v, got %v", i, v.exp, v.iface)
		}
	}
}

func TestAppend(t *testing.T) {
	type appendCase struct {
		data   interface{}
		Parent Appending
		expd   interface{}
		err    error
	}

	type faildummy struct{ Name string }

	var appendCases = []appendCase{
		{&Trunk{Name: "test1", Endpoint: Endpoint{Name: []string{"test1"}, Allow: []string{"1", "2"}}}, &PJSIP{},
			&PJSIP{Trunks: []Trunk{{Name: "test1", Endpoint: Endpoint{Name: []string{"test1"}, Allow: []string{"1", "2"}}}}}, nil},
		{&Aor{Name: []string{"test2"}, MaxContacts: "3"}, &Trunk{},
			&Trunk{Aor: Aor{Name: []string{"test2"}, MaxContacts: "3"}}, nil},
		{Endpoint{}, &Trunk{}, &Trunk{}, errors.New("appending item must be a pointer")},
		{&faildummy{Name: "dummy"}, &PJSIP{}, &PJSIP{}, errors.New("pjsip doesnt have such a field")},
	}

	for i, v := range appendCases {
		name := ""
		if reflect.ValueOf(v.data).Kind() != reflect.Ptr {
			name = reflect.TypeOf(v.data).Name()
		} else {
			name = reflect.TypeOf(v.data).Elem().Name()
		}

		val := reflect.ValueOf(v.Parent).Elem()
		err := Append(v.Parent, name, v.data, val)

		if err == nil && v.err != nil {
			t.Fatalf("testcase %d: expected error %s, got nil", i, v.err)
		}

		if err != nil && v.err == nil {
			t.Fatalf("testcase %d: unexpected error %s", i, err)
		}

		if !reflect.DeepEqual(v.Parent, v.expd) {
			t.Fatalf("testcase %d: expected result %s, got %s", i, v.expd, v.Parent)
		}
	}
}

func TestFillUnknown(t *testing.T) {
	type fillUnkCase struct {
		data string
		expd *UnknownStruct
		err  bool
	}

	var testcases = []fillUnkCase{
		{"[asterisk-publication]\ntype=asterisk-publication\ndevicestate_publish=test\nmailboxstate_publish=pubtest\ndevice_state=no\ndevice_state_filter=S+\n",
			&UnknownStruct{Name: []string{"asterisk-publication"}, Type: "asterisk-publication", LINES: []string{
				"devicestate_publish=test", "mailboxstate_publish=pubtest", "device_state=no", "device_state_filter=S+"}}, false},
		{"Hello, world!", &UnknownStruct{}, true},
	}

	for i := range testcases {
		u := new(UnknownStruct)

		if err := fillUnknown(testcases[i].data, u) != nil; err != testcases[i].err {
			t.Fatalf("testcase %d: expected error %t, got %t", i, testcases[i].err, err)
		}

		if testcases[i].err && !reflect.DeepEqual(u, testcases[i].expd) {
			t.Fatalf("testcase %d: expected %v, got %v", i, testcases[i].expd, u)
		}
	}
}
