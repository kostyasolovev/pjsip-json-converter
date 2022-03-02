package pjsip

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"reflect"
	"time"
)

const RecursionSpeedLimit int = 10 // сдерживает рекурсию, благодаря чему блоки читаются по порядку

// Простой хелпер, конвертирует json(pjsip.PJSIP) в файл конфига pjsip.conf по указанному пути.
func JSONtoConfv2(in []byte, path string) error {
	fmt.Println("translating in process...")

	p := new(PJSIP)
	if err := json.Unmarshal(in, p); err != nil {
		return err
	}

	var out = "; *** Created by Business-Line ***\n"
	// ключевая функция - маршаллим структуру в строку по формату pjsip конфига.
	ReadStruct(p, "", &out)

	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}
	// сохраняем результат в файл
	_, err = f.Write([]byte(out))
	if err != nil {
		return err
	}

	fmt.Println("OK: pjsip.conf was writed successfully")

	return nil
}

// Рекурсивно читает структуры. Приводит строки к формату конфига pjsip.conf
// Reads struct recursively. Распаковывает структуру рекурсивно и приводит к виду pjsip.conf.
func ReadStruct(st interface{}, parent string, out *string) { // nolint: gocognit, gocyclo
	val := reflect.ValueOf(st)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}
	// мы выходим на этот кейс при распаковке слайсов строк
	if val.Kind() == reflect.String {
		*out += fmt.Sprintf("%v=%v\n", parent, val.String())
		return
	}
	// итерируемся по полям структуры
	for i := 0; i < val.NumField(); i++ {
		f, fname, tag := val.Field(i), val.Type().Field(i).Name, GetTag(val, i, "pjsip")
		// если значение тэга pjsip == "omit" -> пропускаем
		if len(tag) != 0 && tag[0] == "omit" || f.IsZero() {
			continue
		}
		// при рекурсивном чтении сложной структуры мы можем встретить структуру, слайс (структур или строк) или строку
		switch f.Kind() { // nolint: exhaustive
		case reflect.Struct:
			time.Sleep(time.Duration(RecursionSpeedLimit) * time.Millisecond) // чтобы структуры читались по порядку, а не одновременно
			ReadStruct(f.Interface(), fname, out)
		case reflect.Slice:
			// кейс для темплейтов [name](template)
			if fname == "Name" { // nolint: goconst
				hastemplate, notemplate := 2, 1
				if f.Len() == hastemplate {
					*out += fmt.Sprintf("\n[%s](%s)\n", f.Index(0).Interface(), f.Index(1).Interface())
				} else if f.Len() == notemplate {
					*out += fmt.Sprintf("\n[%s]\n", f.Index(0).Interface())
				}

				continue
			}
			// кейс для Unknown блоков
			if fname == "LINES" {
				readUnknownLines(f, out)
				continue
			}
			// для всего остального
			for j := 0; j < f.Len(); j++ {
				ReadStruct(f.Index(j).Interface(), Under_score(fname), out)
			}
		case reflect.String:
			if fname == "Name" {
				if parent == "trunks" {
					continue
				}

				*out += fmt.Sprintf("\n[%v]\n", f)
			} else {
				*out += fmt.Sprintf("%v=%v\n", Under_score(fname), f)
			}
		}
	}
}

// распаковывает Unknown структуры в строки по правилам pjsip.conf.
func readUnknownLines(val reflect.Value, out *string) {
	if val.Kind() != reflect.Slice {
		log.Fatalf("failed to read block: %v\n", val.Interface())
		return
	}

	for i := 0; i < val.Len(); i++ {
		*out += fmt.Sprintf("%v\n", val.Index(i).Interface())
	}
}
