package pjsip

import (
	"regexp"
	"testing"
)

func TestReadStruct(t *testing.T) {
	type unmarshalPair struct {
		orig *PJSIP
		expd string
	}

	var unmarshalcases = []unmarshalPair{
		{&PJSIP{Transports: []Transport{{Name: []string{"simple-trans"}, Type: "transport", Bind: "0:0:0:0"},
			{Name: []string{"transport_udp"}, Type: "transport", Bind: "1:1:1:1"}},
			Trunks: []Trunk{{Name: "6001", Auth: Auth{Name: []string{"6001", "template"}, Type: "auth"}},
				{Name: "mytrunk1", Aor: Aor{Name: []string{"mytrunk1"}, Type: "aor"}},
			}},
			"[simple-trans]\ntype=transport\nbind=0:0:0:0\n\n[transport_udp]\ntype=transport\nbind=1:1:1:1\n\n[6001](template)\ntype=auth\n\n[mytrunk1]\ntype=aor\n"},
		{&PJSIP{Trunks: []Trunk{{Name: "6002", Endpoint: Endpoint{Name: []string{"6002", "!"}, Type: "endpoint", Disallow: []string{"ulaw", "alaw"}}},
			{Name: "mytrunk2", Unknowns: []UnknownStruct{{Name: []string{"mytrunk2"}, Type: "custom", LINES: []string{"field_1=true", "field_2=false"}}}}},
			System: System{Name: "system2", Type: "system"},
		},
			"[system2]\ntype=system\n\n[6002](!)\ntype=endpoint\ndisallow=ulaw\ndisallow=alaw\n\n[mytrunk2]\ntype=custom\nfield_1=true\nfield_2=false\n"},
	}

	for i, v := range unmarshalcases {
		out := ""

		ReadStruct(v.orig, "", &out)
		// пилим результат на блоки и сравниваем без учета порядка блоков, т.к. порядок не имеет значения
		outblocks := regexp.MustCompile(`\[\S+\](.|\s)*?(\[|\s{2}|$)`).FindAllString(out, 25)
		expdblocks := regexp.MustCompile(`\[\S+\](.|\s)*?(\[|\s{2}|$)`).FindAllString(v.expd, 25)

		if !CompareSlices(outblocks, expdblocks) {
			re := regexp.MustCompile(`\s`)
			t.Fatalf("testcase %d: expected: %s\n got: %s", i, re.ReplaceAllString(v.expd, "__"), re.ReplaceAllString(out, "__"))
		}
	}
}
