package freeswitch

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"pjsip-handler/src/pjsip"
	"reflect"
	"strings"
)

// Вспомогательные структуры, позволяющие использовать пакет encoding/xml для маршаллинга в xml.
type Include struct {
	Gateway Gateway `xml:"gateway"`
	Param   []Param `xml:"param"`
}
type Param struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type Gateway struct {
	Name string `xml:"name,attr"`
}

// unmarshals FS config to Go struct.
func SimpleXMLUnmarshal(path string, s *Include) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = xml.Unmarshal(data, s)
	if err != nil {
		return err
	}

	return nil
}

// Конвертирует PJSIP в Include. Include предназначен для маршаллинга в xml.
func PjsipToFS(in *pjsip.PJSIP) ([]Include, error) {
	trunksCount := len(in.Trunks)
	// идея в том, что PJSIP может содержать несколько транков,
	// соответственно при распаковке во freeswitch мы можем получить несколько конфигов
	result := make([]Include, 0, trunksCount)

	for i := range in.Trunks {
		mp := make([]string, 0, GetTotalFSKeysNumber())
		// парсим PJSIP на строки
		ReadPjsipAsFS(in.Trunks[i], "", &mp)

		inc := Include{}
		// а строки добавляем в Include
		if err := unmarshalToInclude(mp, &inc); err != nil {
			return []Include{}, err
		}

		result = append(result, inc)
	}

	return result, nil
}

// Преобразует транк pjsip в мапу, для последующего анмаршаллинга в xml структуру.
// ent - преобразуемая структура;
// parent - название родительского поля;
// out - слайс, в который отправляем результаты.
func ReadPjsipAsFS(ent interface{}, parent string, out *[]string) {
	val := reflect.ValueOf(ent)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() == reflect.String {
		line := parent + ":" + val.String()
		*out = append(*out, line)

		return
	}
	// итерируемся по полям структуры
	for i := 0; i < val.NumField(); i++ {
		field, tag := val.Field(i), pjsip.GetTag(val, i, "fswitch")
		// если тэга нет или поле пустое - пропускаем его
		if len(tag) == 0 || field.IsZero() {
			continue
		}

		switch field.Kind() { // nolint: exhaustive
		case reflect.String:
			// берем значение поля, если оно булевое - переводим в true/false
			line := pjsip.ConvertToFSBoolleans(field.String())

			*out = append(*out, tag[0]+":"+line)
		case reflect.Slice:
			for j := 0; j < field.Len(); j++ {
				ReadPjsipAsFS(field.Index(j).Interface(), tag[0], out)
			}
		case reflect.Struct:
			ReadPjsipAsFS(field.Interface(), tag[0], out)
		default:
			log.Printf("unsupported field type %s (%s)", field.Kind(), field.Type().Name())
		}
	}
}

// анмаршаллит слайс строк формата "paramName:paramValue" в структуру Include.
func unmarshalToInclude(data []string, inc *Include) error {
	for i := range data {
		params := strings.Split(data[i], ":")

		paramName, paramVal := params[0], params[1]
		// имя параметра получаем из тэга fswitch, если тэг заполнен неправильно, он может оказаться пустым
		if paramName == "" {
			err := fmt.Errorf("param name must be non-empty: %v", params)
			return err
		}

		if paramName == "gateway" {
			gate := Gateway{Name: paramVal}
			inc.Gateway = gate

			continue
		}

		p := Param{paramName, paramVal}
		inc.Param = append(inc.Param, p)
	}

	return nil
}
