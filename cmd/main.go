package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"pjsip-handler/src/freeswitch"
	"pjsip-handler/src/pjsip"
	"pjsip-handler/src/sip"
	"strings"
)

// тестим примеры из README.
func main() {
	// pjsip.conf to json
	jsone, err := one("./../input/pjsip.conf")
	if err != nil {
		log.Println(err)
	}

	// pjsip.conf to json, 2й вариант
	data, err := ioutil.ReadFile("./../input/pjsip2.conf")
	if err != nil {
		log.Println(err)
	}

	p := new(pjsip.PJSIP)
	if err := two(data, p); err != nil {
		log.Println(err)
	}

	jsone, err = json.Marshal(p)
	if err != nil {
		log.Println(err)
	}

	_ = jsone

	// наоборот: из json в pjsip.conf
	if err = three(jsone, "./../output/pjsip_1.conf"); err != nil {
		log.Println(err)
	}
	// FREESWITCH -> pjsip.PJSIP
	p, err = five("./../input/freeswitch_example.xml")
	if err != nil {
		log.Println(err)
	}

	// pjsip.PJSIP -> FREESWITCH.xml
	if err = six(p, "./../output/fs_out.xml"); err != nil {
		log.Println(err)
	}

	// chan_sip.conf -> pjsip.conf
	if answer, err := seven("./../input/sip.conf", "./../output/pjsip_from_sip.conf"); err != nil {
		log.Println(err)
	} else {
		fmt.Println(string(answer))
	}

	// pjsip.Pjsip -> chan_sip.conf
	if err = eight(p, "./../output/sip_from_pjsip.conf"); err != nil {
		log.Println(err)
	}
}

// переводим pjsip.conf в json
func one(path string) ([]byte, error) {
	json, err := pjsip.ConfToJSONv2(path)
	if err != nil {
		return []byte{}, err
	}

	fmt.Println("func ONE SUCCSESS: pjsip.conf converted to json")

	return json, nil
}

// 2й вариант: анмаршаллим данные в структуру
func two(data []byte, out *pjsip.PJSIP) error {
	if err := pjsip.Unmarshal(out, string(data)); err != nil {
		return err
	}

	fmt.Println("func TWO SUCCSESS: pjsip.conf converted to json")

	return nil
}

// распаковываем структуру PJSIP в файл конфига по адресу outputPath
func three(jsone []byte, outputPath string) error {
	if err := pjsip.JSONtoConfv2(jsone, outputPath); err != nil {
		return err
	}

	fmt.Println("func THREE SUCCSESS: json converted to pjsip.conf")

	return nil
}

// FREESWITCH.xml to pjsip.PJSIP
func five(path string) (*pjsip.PJSIP, error) {
	out := new(pjsip.PJSIP)

	data, err := ioutil.ReadFile(path)
	if err != nil {
		return out, err
	}

	if err = freeswitch.ToPJSIP(data, out); err != nil {
		return &pjsip.PJSIP{}, err
	}

	log.Println("func FIVE SUCCSESS: freeswitch.xml converted to pjsip.PJSIP")

	return out, nil
}

func six(pj *pjsip.PJSIP, pathToSave string) error {
	// идея в том, что структура pjsip может содержать 1+ транков, поэтому на выходе
	// может быть несколько xml файлов
	confs, err := freeswitch.PjsipToFS(pj)
	if err != nil {
		return err
	}

	for i := range confs {
		b, err := xml.Marshal(confs[i])
		if err != nil {
			return err
		}

		file, err := os.OpenFile(pathToSave+"_"+fmt.Sprintln(i), os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			return err
		}

		_, err = file.Write(b)
		if err != nil {
			return err
		}
	}

	log.Println("func SIX SUCCSESS: pjsip.PJSIP converted to freeswitch.xml")

	return nil
}

// переводит sip_chan.conf в pjsip.conf
// input (опционально) - место расположения sip_chan.conf
// output (опционально) - место для создания pjsip.conf
func seven(input string, output string) ([]byte, error) {
	answer, err := sip.ConvertToPjsip("./../input/sip.conf", "./../output/pjsip_from_sip.conf")
	if err != nil {
		return []byte{}, err
	}

	log.Println("func SEVEN SUCCSESS: chan_sip.conf converted to pjsip.conf")

	return answer, nil
}

// pjsip.PJSIP -> chan_sip.conf
func eight(p *pjsip.PJSIP, path string) error {
	// слайс строк, в который будем писать
	out := make([]string, 0, len(p.Trunks)*15) // капасити для оптимизации памяти, можно просто out := []string{}

	sip.ReadPjsipAsSIP(p, "", &out) // ключевая функция - парсит структуру в текст sip.conf

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return err
	}

	_, err = file.Write([]byte(strings.Join(out, "")))
	if err != nil {
		return err
	}

	log.Println("func EIGHT SUCCSESS: sip.conf converted to pjsip.PJSIP")

	return nil
}
