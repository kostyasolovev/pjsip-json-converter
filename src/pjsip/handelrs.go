package pjsip

import (
	"fmt"
	"log"
	"reflect"
	"regexp"
	"strings"
)

// Анмаршаллит из строк любые структуры кроме транков.
func (p *PJSIP) commonHandle(block string, m TypeRegistry, val reflect.Value) error {
	// узнаем название структуры по значению строки "type="
	fname := CamelCase(ReduceWrap(regexp.MustCompile(`type=\S+`).FindString(block), "type="))
	// создаем структуру по ее названию
	temp := m.MakeStruct(fname)
	// заполняем созданную структуру
	if err := fillStruct(block, temp); err != nil {
		return err
	}
	// аппендим в родительскую структуру
	err := Append(p, fname, temp, val)

	return err
}

// Анмаршаллит транки из строк.
func (p *PJSIP) trunkHandle(in, name string, m TypeRegistry) error {
	trunk := new(Trunk)
	trunk.Name = name
	val := reflect.ValueOf(trunk).Elem()
	// парсим на блоки
	blocks := regexp.MustCompile(`\[\S+\](.|\s)*?(\[|\s{2})`).FindAllStringIndex(in, MaxBlocksPerTrunk)
	for i := range blocks {
		// область между индексами это какая-то структура
		left, right := blocks[i][0], blocks[i][1]
		// узнаем название структуры, создаем ее, заполняем
		blockType := CamelCase(ReduceWrap(regexp.MustCompile(`type=\S+`).FindString(in[left:right]), "type="))
		if _, ok := m[blockType]; ok {
			temp := m.MakeStruct(blockType)
			if err := fillStruct(in[left:right], temp); err != nil {
				return err
			}
			// аппендим структуру в родительский транк
			if err := Append(trunk, blockType, temp, val); err != nil {
				return err
			}
		} else {
			// если тип пустой или неизвестный (отсутствует в TypeRegistry), то анмаршаллим блок в структуру Unknown
			temp := new(UnknownStruct)
			if err := fillUnknown(in, temp); err != nil {
				return err
			}

			trunk.Unknowns = append(trunk.Unknowns, (*temp))
		}
	}

	p.Trunks = append(p.Trunks, (*trunk))

	return nil
}

type Appending interface {
	Append(name string, in interface{}, val reflect.Value) error
}

// Аppends struct to parent struct.
func Append(data Appending, name string, in interface{}, val reflect.Value) error {
	if reflect.ValueOf(in).Kind() != reflect.Ptr {
		err := fmt.Errorf("failed to append: appending item must be a pointer [%s] ", name)
		return err
	}

	err := data.Append(name, in, val)

	return err
}

// аппендит структуру в Транк. Name - имя структуры, Val - reflect.Value транка.
func (t *Trunk) Append(name string, in interface{}, val reflect.Value) error {
	// Ищем в pjsip.Trunk поле с именем name, если такого поля нет - пропускаем
	f := val.FieldByName(name)
	if f.CanSet() && f.IsValid() {
		switch f.Kind() { // nolint: exhaustive
		case reflect.Slice:
			f.Set(reflect.Append(f, reflect.ValueOf(in).Elem()))
		case reflect.Struct:
			f.Set(reflect.ValueOf(in).Elem())
		default:
			err := fmt.Errorf("failed to append: unsupported field type: %s", f.Kind())
			return err
		}
	} else {
		err := fmt.Errorf("failed to append: struct's field is not settable or valid: %s", name)
		return err
	}

	return nil
}

// Аппендит интерфейс в pjsip.PJSIP. Интерфейс должен быть указателем на структуру.
func (p *PJSIP) Append(name string, in interface{}, val reflect.Value) error {
	// ищем поле структуры с похожим именем, учитываем, что имя может иметь окончание "s"
	if _, ok := val.Type().FieldByName(name); !ok {
		if _, ok := val.Type().FieldByName(name + "s"); !ok {
			err := fmt.Errorf("failed to append: pjsip doesnt have such a field [%s] ", name)
			return err
		}

		name += "s"
	}
	// проверяем на валидность и редактируемость
	f := val.FieldByName(name)
	if f.CanSet() && f.IsValid() {
		switch f.Kind() { // nolint: exhaustive
		case reflect.Slice:
			f.Set(reflect.Append(f, reflect.ValueOf(in).Elem()))
		case reflect.Struct:
			f.Set(reflect.ValueOf(in).Elem())
		default:
			err := fmt.Errorf("failed to append: unsupported pjsip field type [%s] ", f.Type().Name())
			return err
		}
	} else {
		err := fmt.Errorf("failed to append: struct field aren'n settable or valid [%s] ", name)
		return err
	}

	return nil
}

// Aнмаршаллит строку в структуру.
func fillStruct(block string, in interface{}) error {
	if reflect.ValueOf(in).Kind() != reflect.Ptr {
		err := fmt.Errorf("input interface must be a pointer %T", in)
		return err
	}

	blockMap, err := toMap(block)
	if err != nil {
		return err
	}

	if len(blockMap) == 0 {
		err := fmt.Errorf("invalid block: %s", piece(block))
		return err
	}

	val := reflect.ValueOf(in).Elem()
	for k, v := range blockMap {
		field := val.FieldByName(k)
		if field.CanSet() && field.IsValid() {
			switch field.Kind() { // nolint: exhaustive
			case reflect.Slice:
				field.Set(reflect.ValueOf(strings.Split(v, " ")))
			case reflect.String:
				field.SetString(v)
			default:
				log.Println("unsupported field type", field.Kind())
			}
		} else {
			err := fmt.Errorf("struct's field is not settable or valid: %s", k)
			return err
		}
	}

	return nil
}

// Преобразует строку в мапу.
func toMap(input string) (map[string]string, error) {
	input = strings.TrimRight(input, "\n[")
	rows := strings.Split(input, "\n")
	blockMap := make(map[string]string)
	i := 0
	// заполняем ключ Name
	name := strings.Split(regexp.MustCompile(`\[\S+\](\(.+\))*`).FindString(input), "(")
	if len(name) != 0 && name[0] != "" {
		for j := range name {
			name[j] = strings.Trim(name[j], "[]()")
		}

		blockMap["Name"] = strings.Join(name, " ")
		i++
	}
	// итерируемся по строкам, заполняем остальные ключи
	for i < len(rows) {
		validlen := 2
		// разрезаем строку по символу "=", если полученный слайс имеет длину меньше 2х, значит строка невалидная
		kv := strings.Split(rows[i], "=") // пара ключ-значение
		if len(kv) != validlen {
			err := fmt.Errorf("invalid line: [%s] ", rows[i])
			return map[string]string{}, err
		}
		// если ключ повторяется, значит имеем повторяющиеся параметры ("allow=ulaw\nallow=mulaw")
		// добавляем их в к предыдущему значению через пробел
		// в будущем такие значения будут преобразованы в слайсы
		kv[0] = CamelCase(kv[0])
		if _, ok := blockMap[kv[0]]; ok {
			blockMap[kv[0]] = fmt.Sprintf("%s %s", blockMap[kv[0]], kv[1])
		} else {
			blockMap[kv[0]] = kv[1]
		}
		i++
	}

	return blockMap, nil
}

// Анмаршаллит строку в Unknown struct.
func fillUnknown(input string, u *UnknownStruct) error {
	m, err := toMapUnknown(input)
	if err != nil {
		return err
	}

	if len(m) == 0 {
		err := fmt.Errorf("block is invalid: %s", piece(input))
		return err
	}

	val := reflect.ValueOf(u).Elem()

	for k, v := range m {
		f := val.FieldByName(k)
		if f.CanSet() && f.IsValid() {
			switch f.Kind() { // nolint: exhaustive
			case reflect.Slice:
				f.Set(reflect.ValueOf(strings.Split(v, " ")))
			case reflect.String:
				f.SetString(v)
			default:
				err := fmt.Errorf("struct's field is not settable or valid: %s", k)
				return err
			}
		}
	}

	return nil
}

// Создает мапу из строки для транков неизвестного типа.
func toMapUnknown(input string) (map[string]string, error) {
	input = strings.TrimRight(input, "\n[")
	m := make(map[string]string)
	name := strings.Split(regexp.MustCompile(`\[\S+\](\(.+\))*`).FindString(input), "(")
	// принцип как у функции toMAp
	if len(name) != 0 && name[0] != "" {
		for i := range name {
			name[i] = strings.Trim(name[i], "[]()")
		}

		m["Name"] = strings.Join(name, " ")
		input = input[regexp.MustCompile(`\s`).FindStringIndex(input)[1]:]
	}
	// заполняем ключ "Type"
	typo := regexp.MustCompile(`type=\S+`).FindStringIndex(input)
	if len(typo) != 0 {
		m["Type"] = ReduceWrap(input[typo[0]:typo[1]], "type=")
		input = input[:typo[0]] + input[typo[1]+1:]
	}
	// заполняем ключ "LINES"
	lines := strings.Split(input, "\n")
	for i := range lines {
		if regexp.MustCompile(`=`).MatchString(lines[i]) {
			if _, ok := m["LINES"]; !ok {
				m["LINES"] = lines[i]
				continue
			}

			m["LINES"] = fmt.Sprintf("%s %s", m["LINES"], lines[i])
		} else {
			err := fmt.Errorf("invalid line: %s", lines[i])
			return map[string]string{}, err
		}
	}

	return m, nil
}
