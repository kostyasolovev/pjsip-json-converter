package pjsip

import (
	"reflect"
	"testing"
)

type translatePair struct {
	orig   string
	expect string
}

func TestCamelCase(t *testing.T) {
	var CamelTests = []translatePair{
		{"under_score", "UnderScore"},
		{"md5_flag", "Md5Flag"},
		{"simple", "Simple"},
		{"very_long_line100", "VeryLongLine100"},
		{"CAPS", "CAPS"},
		{"last_b", "LastB"},
	}

	for i := range CamelTests {
		res := CamelCase(CamelTests[i].orig)
		if res != CamelTests[i].expect {
			t.Errorf("Expected %s, got %s\n", CamelTests[i].expect, res)
		}
	}
}

func TestUnderScore(t *testing.T) {
	var UTests = []translatePair{
		{"CamelCase", "camel_case"},
		{"Md5Flag", "md5_flag"},
		{"Simple", "simple"},
		{"VeryLongLine100", "very_long_line100"},
		{"CAPS", "CAPS"},
		{"LastV", "last_v"},
	}

	for i := range UTests {
		res := Under_score(UTests[i].orig)
		if res != UTests[i].expect {
			t.Errorf("Expected %s, got %s\n", UTests[i].expect, res)
		}
	}
}

func TestTrimAuth(t *testing.T) {
	var trims = []translatePair{
		{"[auth6000]", "6000"},
		{"[auth_6001]", "6001"},
		{"[mytrunk-auth]", "mytrunk"},
		{"[mytrunk2auth]", "mytrunk2"},
	}

	for i := range trims {
		res := TrimAuth(trims[i].orig)
		if res != trims[i].expect {
			t.Errorf("Expected %s, got %s\n", trims[i].expect, res)
		}
	}
}

func TestLispCase(t *testing.T) {
	var LispTests = []translatePair{
		{"LispCase", "lisp-case"},
		{"Md5Flag", "md5-flag"},
		{"Simple", "simple"},
		{"VeryLongLine100", "very-long-line100"},
		{"CAPS", "CAPS"},
		{"LastV", "last-v"},
	}

	for i := range LispTests {
		res := lisp_case(LispTests[i].orig)
		if res != LispTests[i].expect {
			t.Errorf("Expected %s, got %s\n", LispTests[i].expect, res)
		}
	}
}

func TestGetTag(t *testing.T) {
	type tagDummy struct {
		Name  string `json:",omitempty" sip:"type,testsip" pjsip:",omit"`
		Type  string `pjsip:"some,testpjsip"`
		Empty string
	}

	type expdTags struct {
		tags []string
	}

	var expdtags = []expdTags{
		{[]string{"", "omitempty"}},
		{[]string{`type`, "testsip"}},
		{[]string{"", "omit"}},
		{[]string{}},
		{[]string{}},
		{[]string{}},
		{[]string{"some", "testpjsip"}},
		{[]string{}},
		{[]string{}},
		{[]string{}},
		{[]string{}},
		{[]string{}},
	}

	d := &tagDummy{"Victoria", "dummy", "someval"}
	tagToLook := []string{"json", "sip", "pjsip", ""}

	val := reflect.ValueOf(d).Elem()
	for i := 0; i < val.NumField(); i++ {
		for j := 0; j < len(tagToLook); j++ {
			tag := GetTag(val, i, tagToLook[j])

			if !CompareSlices(tag, expdtags[i*len(tagToLook)+j].tags) {
				t.Fatalf("testcase %d: expected %v, got %v", i*len(tagToLook)+j, expdtags[i*len(tagToLook)+j].tags, tag)
			}
		}
	}
}

func TestGetAttr(t *testing.T) {
	type getAttrPair struct {
		orig string
		expd []string
	}

	var attrpairs = []getAttrPair{
		{`param name="username" value="cluecon"/>\n<!--/// auth realm: *optional* same as gateway name, if blank ///-->\n<param name="realm"`,
			[]string{"username", "cluecon"}},
		{`param name="" value="cluecon2" some text`, []string{"", "cluecon2"}},
		{`absolutely no prefs here`, []string{}},
	}

	for i := range attrpairs {
		res := GetAttr(attrpairs[i].orig, "name", "value")

		if !reflect.DeepEqual(res, attrpairs[i].expd) {
			t.Fatalf("testcase %d: expected %v, got %v", i, attrpairs[i].expd, res)
		}
	}
}

func TestCompareSlices(t *testing.T) {
	type compSliceCase struct {
		a    []string
		b    []string
		expd bool
	}

	var testcases = []compSliceCase{
		{[]string{"from-user:someuser", "gateway:test1", "username:somename"}, []string{"from-user:someuser", "gateway:test1", "username:somename"}, true},
		{[]string{"gateway:test3", "register-proxy:someproxy", "register:true", "caller-id-in-from:false"}, []string{"gateway:test3", "caller-id-in-from:false", "register-proxy:someproxy", "register-proxy:true"}, false},
		{[]string{"a", "a", "b"}, []string{"a", "b", "a"}, true},
		{[]string{"a", "b", "c"}, []string{"a", "c", "c"}, false},
		{[]string{"a", "b"}, []string{"a"}, false},
	}

	for i := range testcases {
		if res := CompareSlices(testcases[i].a, testcases[i].b); res != testcases[i].expd {
			t.Fatalf("testcase %d: expected %t, got %t", i, testcases[i].expd, res)
		}
	}
}

func TestConvertToFSBoolleans(t *testing.T) {
	var tcases = map[string]string{"yes": "true", "no": "false", "never": "false", "noboo": "noboo"}
	for k, v := range tcases {
		if res := ConvertToFSBoolleans(k); res != v {
			t.Fatalf("expected %s, got %s", v, res)
		}
	}
}
