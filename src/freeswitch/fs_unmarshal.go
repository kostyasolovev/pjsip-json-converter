package freeswitch

import (
	"fmt"
	"pjsip-handler/src/pjsip"
	"reflect"
	"regexp"
	"strings"
)

// Анмаршаллим xml конфиг в структуру PJSIP.
func ToPJSIP(data []byte, out *pjsip.PJSIP) error {
	m, err := freeSwitchToMap(data)
	if err != nil {
		return err
	}

	trunk := new(pjsip.Trunk)
	trunk.Name = m["TrunkName"]
	// структуры, создаваемые далее, будут иметь статик тайп интерфейс
	// соответственно для них недоступны будут стандартные операции со структурами
	// чтобы зааппендить такую структуру в родительский Trunk нужен его reflect.Value
	trunkVal := reflect.ValueOf(trunk).Elem()
	// typos - мапа рефлект типов для создания структур
	typos := pjsip.InitRegistration()
	seek := make(map[string]bool)
	// ключ мапы выглядит как "parent/child", где parent - имя структуры, child - имя поля структуры
	for k := range m {
		if k == "TrunkName" {
			continue
		}

		fieldNames := strings.Split(k, "/")
		if _, ok := seek[fieldNames[0]]; ok {
			continue
		}
		// ключи в мапе смешаны, чтобы сфокусироваться на конкретном ключе запоминаем его в переменную cur
		cur := fieldNames[0]
		// запоминаем, какие ключи мы уже смотрели
		seek[fieldNames[0]] = true
		// а теперь ищем ключи, удовлетворяющие cur
		source := make(map[string]string)

		for key, value := range m {
			if parent := strings.Split(key, "/"); parent[0] == cur {
				source[pjsip.CamelCase(parent[1])] = value
			}
		}
		// cur - название структуры, ищем ее в мапе и создаем
		temp := typos.MakeStruct(pjsip.CamelCase(cur))
		if err := fillStructFromXMLMap(temp, source); err != nil {
			return err
		}
		// аппендим созданную структуру в родительский транк
		if err := trunk.Append(pjsip.CamelCase(cur), temp, trunkVal); err != nil {
			return err
		}
	}

	out.Trunks = append(out.Trunks, (*trunk))

	return nil
}

// Заполняет выбранную структуру значениями из мапы. Ключ мапы = название поля структуры.
func fillStructFromXMLMap(temp interface{}, data map[string]string) error {
	val := reflect.ValueOf(temp)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	// итерируемся по мапе
	for k, v := range data {
		// узнаем, есть ли в выбранной структуре поле с таким названием
		field := val.FieldByName(k)
		if !field.CanSet() || !field.IsValid() {
			err := fmt.Errorf("invalid field name: %s", k)
			return err
		}
		// если есть - устанавливаем значение в поле структуры
		switch field.Kind() { // nolint: exhaustive
		case reflect.String:
			field.SetString(v)
		case reflect.Slice:
			field.Set(reflect.Append(field, reflect.ValueOf(v)))
		default:
			err := fmt.Errorf("usupported field type: %s", field.Kind())
			return err
		}
	}

	return nil
}

// Парсит xml файл в мапу с ключами pjsip.
func freeSwitchToMap(data []byte) (map[string]string, error) {
	m := make(map[string]string)
	// заполняем мапу. Во-первых, Gateway (Trunk name)
	gateInd := regexp.MustCompile(`gateway.name=\".+?\"`).FindIndex(data)
	if len(gateInd) != 0 {
		gateVal := strings.Trim(pjsip.ReduceWrap(string(data[gateInd[0]:gateInd[1]]), "gateway.name="), `"`)
		m["TrunkName"] = gateVal
		// теперь этот кусок можно не учитывать
		data = data[gateInd[1]:]
	}
	// далее парсим параметры param
	lines := strings.Split(string(data), "param")
	for i := range lines {
		// каждый лайн param содержит атрибуты name value, парсим их
		namevalue := pjsip.GetAttr(lines[i], "name", "value")
		if len(namevalue) == 0 {
			continue
		}
		// ругаемся на лайн с пустым именем
		if namevalue[0] == "" {
			err := fmt.Errorf("param name must be non-empty ...%s... ", lines[i])
			return map[string]string{}, err
		}
		// проверяем правильное ли имя параметра
		v, ok := FreeSwitchPJSIPEquivalents[namevalue[0]]
		if !ok {
			err := fmt.Errorf("invalid param name %s", namevalue[0])
			return map[string]string{}, err
		}
		// в pjsip булевые значения имеют формат yes/no, а у freeswitch - true/false, меняем если нужно
		linehasvalue := 2
		if len(namevalue) == linehasvalue {
			if namevalue[1] == "true" {
				namevalue[1] = "yes"
			} else if namevalue[1] == "false" {
				namevalue[1] = "no"
			}
		}
		// ключ это маршрут "Х/y", где Х - название структуры, у - название поля структуры
		pjsipkey := strings.Join(v, "/")
		m[pjsipkey] = namevalue[1]
	}

	return m, nil
}
