package pjsip

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

const (
	MAXTRUNKS         = 100 // максимальное кол-во транков(блоков)
	MaxBlocksPerTrunk = 25  // максимальное кол-во блоков в транке
)

// Простой хелпер для конвертации pjsip.conf to JSON.
func ConfToJSONv2(path string) ([]byte, error) {
	if path == "" {
		// местонахождение файла по умолчанию
		path = "./../output/pjsip.conf"
	}
	// читаем конфиг
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return []byte{}, err
	}

	P := new(PJSIP)
	// анмаршаллим в PJSIP.
	if err = Unmarshal(P, string(data)); err != nil { // nolint: gocritic
		return []byte{}, err
	}
	// маршаллим PJSIP в json
	out, err := json.Marshal(&P)

	if err != nil {
		return []byte{}, err
	}

	log.Println("OK: pjsip.conf translated to JSON successfully")

	return out, nil
}

// возможность создавать моки функций для тестов.
type Unmarshaler interface {
	trunkHandle(block, name string, m TypeRegistry) error
	commonHandle(block string, m TypeRegistry, val reflect.Value) error
}

// Анмаршаллит строку в pjsip.PJSIP.
func Unmarshal(data Unmarshaler, s string) error {
	// удаляем лишние пробелы, удаляем все комменты
	s = regexp.MustCompile(`(\p{Z})|(;.+\S+)`).ReplaceAllString(s, "")
	s = strings.TrimRight(s, "\n") + "\n\n" // независим от завершающего символа новой строки
	// создаем мапу рефлект тайпов структур, которая позволяет создать структуру по ее названию
	structs := InitRegistration()
	val := reflect.ValueOf(data).Elem()
	// пилим текст на блоки + получаем список уникальных имен
	parts, names, err := parseByNumber(s)
	if err != nil {
		return err
	}

	for i := range names {
		// конкатенируем все блоки с одинаковым именем
		block := parts.concatBlocksByName(names[i], s)
		// узнаем тип блока, от этого зависит тип обработчика
		typo := ReduceWrap(regexp.MustCompile(`type=\S+`).FindString(block), "type=")

		switch typo {
		case "system", "global", "transport", "domain", "phoneprov", "resource_list":
			if err := data.commonHandle(block, structs, val); err != nil {
				return err
			}
		default:
			if err := data.trunkHandle(block, names[i], structs); err != nil {
				return err
			}
		}
	}

	return nil
}

type Part struct {
	name string
	inds []int
}

// имплементация sort.Interface.
type Parts []Part

func (parts Parts) Len() int      { return len(parts) }
func (parts Parts) Swap(i, j int) { parts[i], parts[j] = parts[j], parts[i] }
func (parts Parts) Less(i, j int) bool {
	if parts[i].name == parts[j].name {
		return parts[i].inds[0] < parts[j].inds[0]
	}

	return parts[i].name < parts[j].name
}

// разрезает файл на блоки, возвращает слайс блоков и слайс уникальных имен.
func parseByNumber(input string) (Parts, []string, error) {
	// начало блока - любой текст в квадратных скобках
	re := regexp.MustCompile(`\[\S+\]`)
	// разделяем файл на блоки
	strnames := re.FindAllString(input, MAXTRUNKS)
	if len(strnames) == 0 {
		err := errors.New("pjsip.conf doesnt contain any trunk")
		return []Part{}, []string{}, err
	}
	// создаем массив индексов блоков в тексте
	inds := re.FindAllStringIndex(input, MAXTRUNKS)
	// создаем массив имен блоков
	names := make([]string, 0, len(strnames))
	for i := range strnames {
		names = append(names, TrimAuth(strnames[i]))
	}
	// преобразуем индексы так, чтобы они содержали область до следующего блока либо до конца текста
	for i := 0; i < len(inds)-1; i++ {
		inds[i][1] = inds[i+1][0]
	}
	// заполняем крайний индекс
	inds[len(inds)-1][1] = len(input)
	// создаем слайс Parts
	parts := make([]Part, 0, len(names))
	for i := range names {
		parts = append(parts, Part{
			name: names[i],
			inds: inds[i],
		})
	}
	// сортируем блоки, сортируем имена, удаляем повторяющиеся имена
	sort.Sort(Parts(parts))
	sort.Strings(names)
	removeDuplicates(&names)

	return parts, names, nil
}
