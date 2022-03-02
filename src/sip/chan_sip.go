package sip

import (
	"fmt"
	"log"
	"pjsip-handler/src/pjsip"
	"reflect"
	"regexp"
	"strings"
	"time"
)

// рекурсивно читает структуру PJSIP, отправляет результат в слайс строк.
func ReadPjsipAsSIP(ent interface{}, parent string, out *[]string) { // nolint: gocyclo
	val := reflect.ValueOf(ent)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() == reflect.String {
		if val.String() == "" {
			return
		}

		*out = append(*out, parent+"="+val.String()+"\n")

		return
	}

	if val.Type() == reflect.TypeOf(pjsip.Trunk{}) {
		if register := sipSetRegister(val); register != "" {
			*out = append(*out, fmt.Sprintf("register => %s\n", register))
		}

		name := val.FieldByName("Name").String()

		*out = append(*out, fmt.Sprintf("[%s]\n", name), "type=friend\n")

		if nat := sipSetNat(val); nat != "" {
			*out = append(*out, nat+"\n")
		}

		host := sipSetHost(val)
		*out = append(*out, host)
	}

	if val.Type() == reflect.TypeOf(pjsip.Transport{}) {
		*out = append(*out, fmt.Sprintf("[%s]\n", "general")) // TODO: remove hardcoding
	}

	for i := 0; i < val.NumField(); i++ {
		field, tags := val.Field(i), pjsip.GetTag(val, i, "sip")
		if field.IsZero() || len(tags) == 0 || tags[1] == "omit" {
			continue
		}

		switch field.Kind() { // nolint: exhaustive // баг линтера: требует проверить все 33 кейса, несмотря на default
		case reflect.Struct:
			// ограничиваем скорость рекурсии, чтобы читать по порядку
			time.Sleep(time.Duration(pjsip.RecursionSpeedLimit) * time.Millisecond)
			ReadPjsipAsSIP(field.Interface(), tags[0], out)
		case reflect.Slice:
			for j := 0; j < field.Len(); j++ {
				temp := field.Index(j).Interface()
				ReadPjsipAsSIP(temp, tags[0], out)
			}
		case reflect.String:
			if tags[1] != "" { // вторая часть тэга означает, что там какой-то спецкейс
				sipTagOptionHandle(tags[1], tags[0], out)

				continue
			}

			*out = append(*out, tags[0]+"="+field.String()+"\n")
		default:
			log.Printf("usupported field type %s", field.Kind())
		}
	}
}

// обрабатываем вторую часть тэга (тип поля строка).
func sipTagOptionHandle(tag, parent string, out *[]string) {
	switch tag {
	case "hide":
		return
	case "name":
		*out = append(*out, fmt.Sprintf("[%s]", parent))
	}
}

// устанавливаем строку host.
func sipSetHost(val reflect.Value) string {
	host := ""
	// проверяем поле Contact в структуре Aor: в случае, если тип слайс и его длина больше 0, используем его
	if contactIsNil := val.FieldByName("Aor").FieldByName("Contact"); contactIsNil.Kind() == reflect.Slice && contactIsNil.Len() != 0 {
		if contactFieldVal := val.FieldByName("Aor").FieldByName("Contact").Index(0).String(); contactFieldVal != "" {
			host = regexp.MustCompile(`(\d{1,3}\.){3}\d{1,3}`).FindString(contactFieldVal)
			if host != "" {
				host = "host=" + host + "\n"
				return host
			}
		}
	}

	return "host=dynamic\n"
}

// определяем строку nat.
func sipSetNat(val reflect.Value) string {
	var defaults = []string{"0", "1", "0"}

	if val.Type() != reflect.TypeOf(pjsip.Trunk{}) {
		return ""
	}

	if temp := val.FieldByName("Endpoint").FieldByName("RtpSymmetric"); temp.Interface() == "yes" { // nolint: goconst
		defaults[0] = "1"
	}

	if temp2 := val.FieldByName("Endpoint").FieldByName("ForceRport"); temp2.Interface() == "no" {
		defaults[1] = "0"
	}

	if temp3 := val.FieldByName("Endpoint").FieldByName("RewriteContact"); temp3.Interface() == "yes" {
		defaults[2] = "1"
	}

	switch strings.Join(defaults, "") {
	case "111":
		return "nat=yes"
	case "000":
		return "nat=no"
	case "011":
		return "nat=route"
	}

	return ""
}

// обрабатываем строку register=>sip:...
func sipSetRegister(val reflect.Value) string {
	if val.Type() != reflect.TypeOf(pjsip.Trunk{}) {
		return ""
	}
	// проверяем есть ли что-нибудь в полях ServerUri ClientUri блока Registration
	reg := ""
	if reg = val.FieldByName("Registration").FieldByName("ServerUri").String(); reg == "" {
		if reg = val.FieldByName("Registration").FieldByName("ClientUri").String(); reg == "" {
			return reg
		}
	}

	auth := ""
	// смотрим поля Username и Password. Если непустые - добавляем в результат
	auth += val.FieldByName("Auth").FieldByName("Username").String()
	if auth != "" {
		auth += ":" + val.FieldByName("Auth").FieldByName("Password").String()
	}

	if auth != "" {
		reg = auth + "@" + regexp.MustCompile(`.+@`).ReplaceAllString(reg, "")
	}

	return reg
}
