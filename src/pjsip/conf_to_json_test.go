package pjsip

import (
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"testing"
)

func TestParseByNumber(t *testing.T) {
	type blockie struct {
		orig       string
		expdnames  []string
		expdblocks []string // nolint: structcheck // позднее будут дописаны тестфункция для блоков
		err        error    // nolint: structcheck // позднее будут тесткейсы с ошибками
	}

	var blockies = []blockie{
		{"[6001]\ntest\n[auth6003]\ntest\n[6003]\ntest\n[6002]\ntest\n[transport-udp]",
			[]string{"6001", "6002", "6003", "transport-udp"},
			[]string{"[6001]\ntest\n", "[auth6003]\ntest\n", "[6003]\ntest\n", "[6002]\ntest\n", "[transport-udp]"},
			nil},
		{"failstring\nfailstring\n", []string{}, []string{}, errors.New("pjsip.conf doesnt contain any trunk")},
	}

	for i, v := range blockies {
		_, names, err := parseByNumber(v.orig)
		if err == nil && v.err != nil {
			t.Fatalf("incorrect error behavior: expected %s, got nil, testcase %d", v.err, i)
		}

		if err != nil && v.err == nil {
			t.Fatalf("unexpected error: %s, testcase %d", err, i)
		}

		if !reflect.DeepEqual(names, v.expdnames) {
			t.Fatalf("wrong names: expected %v, got %v, testcase %d", v.expdnames, names, i)
		}
	}
}

type PjsipMock struct {
	Name string
	Type []string
}

// Имплементация интерфейса Unmarshaller.
func (p *PjsipMock) trunkHandle(block, name string, t TypeRegistry) error {
	m := map[string]bool{"aor": false, "auth": false, "transport": true, "acl": false,
		"phoneprov": true, "global": true, "system": true, "domain": true, "registration": false,
		"identify": false, "endpoint": false, "resource_list": true}

	inds := regexp.MustCompile(`\[\S+\]`).FindAllStringIndex(block, 20)
	for i := 0; i < len(inds); i++ {
		left, right := inds[i][1], 0
		if i == len(inds)-1 {
			right = len(block) - 1
		} else {
			right = inds[i+1][0]
		}

		typo := ReduceWrap(regexp.MustCompile(`type=\S+`).FindString(block[left:right]), "type=")
		b, ok := m[typo]

		if !ok {
			p.Type = append(p.Type, "unknown")
			continue
		}

		if b == true {
			err := fmt.Errorf("wrong type passed to trunkhandler: %v", typo)
			return err
		}

		p.Type = append(p.Type, typo)
	}

	return nil
}

// Имплементация интерфейса Unmarshaller.
func (p *PjsipMock) commonHandle(block string, t TypeRegistry, val reflect.Value) error {
	m := map[string]bool{"aor": false, "auth": false, "transport": true, "acl": true,
		"phoneprov": true, "global": true, "system": true, "domain": true, "registration": false,
		"identify": false, "endpoint": false, "resource_list": true}

	typo := ReduceWrap(regexp.MustCompile(`type=\S+`).FindString(block), "type=")
	if typo == "" {
		err := fmt.Errorf("passed block doesnt has a type: %s", piece(block))
		return err
	}

	if correctType, ok := m[typo]; !ok || correctType == false {
		err := fmt.Errorf("wrong type passed to commonhandler %s", typo)
		return err
	}

	p.Type = append(p.Type, typo)

	return nil
}

func TestUnmarshal(t *testing.T) {
	type marshalpair struct {
		orig string
		expd []string
		err  error // nolint: structcheck // будут добавлены тесткейсы с ошибками
	}

	var marshalpairs = []marshalpair{
		{"[6001]\ntype=endpoint\n[auth6003]\ntype=auth\n[6003]\ncontext=true\n[6002]\ntype=identify\n[transport-udp]\ntype=transport",
			[]string{"endpoint", "auth", "unknown", "identify", "transport"}, nil},
		{"[mytrunk]\ntype=system\n[auth6003]\ntype=auth\n[test]\ncontext=true\n[bob]\ntype=custom\n[glob]\ntype=global\n",
			[]string{"system", "auth", "unknown", "unknown", "global"}, nil},
		{"[8920203user]\ntype=registration\n[user007]\ntype=endpoint\n[system]\ncontext=true\n[lucia-user]\ntype=acl\n[phon1]\ntype=phoneprov\n",
			[]string{"registration", "endpoint", "unknown", "acl", "phoneprov"}, nil},
	}

	for i, v := range marshalpairs {
		p := new(PjsipMock)
		if err := Unmarshal(p, v.orig); err != nil {
			t.Fatalf("testcase %d, unexpected error: %s", i, err)
		}

		if !CompareSlices(p.Type, v.expd) {
			t.Errorf("testcase %d: expected %v, got %v", i, v.expd, p.Type)
		}
	}
}
