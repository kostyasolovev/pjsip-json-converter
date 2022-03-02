package pjsip

import (
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strings"
)

// Removes duplicates in-place.
func removeDuplicates(names *[]string) {
	for i := len(*names) - 1; i > 0; i-- {
		if (*names)[i] == (*names)[i-1] {
			*names = append((*names)[:i], (*names)[i+1:]...)
		}
	}
}

// Удаляем суффикс auth и любые неальфанюмерик-символы рядом с ним, удаляем квадратные скобки.
func TrimAuth(s string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9]*(?i)auth[^a-zA-Z0-9]*`).FindStringIndex(s)
	if len(re) == 0 {
		return strings.Trim(s, "[]")
	}

	return strings.Trim(s[:re[0]]+s[re[1]:], "[]")
}

// Конкатенирует блоки с одинаковым именем.
func (parts Parts) concatBlocksByName(name, in string) string {
	var out string

	for _, v := range parts {
		if v.name == name {
			out += in[v.inds[0]:v.inds[1]]
		}
	}

	return out
}

// Удаляет указанный префикс, например: "allow=true" -> "true".
func ReduceWrap(input, pref string) string {
	if input == "" {
		return input
	} else if len(input) <= len(pref) {
		return ""
	}

	return input[len(pref):]
}

// Переводим CamelCase в under_score (snake_case).
func Under_score(in string) string { // nolint: stylecheck, revive
	if in == "" {
		return ""
	}

	if !regexp.MustCompile(`[a-z]+`).MatchString(in) { // проверка на капс
		return in
	}

	maxexpr := 8 // максимальное количество появлений регулярки в тексте
	expr := regexp.MustCompile(`[A-Z]([a-z]+|$)`).FindAllStringIndex(in, maxexpr)

	if len(expr) == 0 {
		return in
	}

	for i := len(expr) - 1; i > 0; i-- {
		in = in[:expr[i][0]] + "_" + in[expr[i][0]:]
	}

	return strings.ToLower(in)
}

// Конвертирует under_score to CamelCase.
func CamelCase(in string) string {
	// under_score, md5_field, simple, CAPS -> UnderScore, Md5Field, Simple, CAPS
	if in == "" {
		return in
	}

	maxexpr := 8                                                         // max expression occurs: AaaBbbCccDdEe...
	re := regexp.MustCompile(`\_[a-z]+`).FindAllStringIndex(in, maxexpr) // \_[a-z]{1,}

	if len(re) == 0 {
		return strings.ToTitle(string(in[0])) + in[1:]
	}

	for i := len(re) - 1; i >= 0; i-- {
		in = in[:re[i][0]] + strings.ToTitle(string(in[re[i][0]+1])) + in[re[i][0]+2:]
	}

	return strings.ToTitle(string(in[0])) + in[1:]
}

// Обрезает текст для упрощения его вывода в консоль.
func piece(s string) string {
	if s == "" {
		return "(empty string)"
	}

	shorttext := 5
	if len(s) <= shorttext {
		return s
	}

	return s[:shorttext+len(s)/4] + "... "
}

// Сравнивает массивы строк без учета их порядка.
func CompareSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	sort.Strings(a)
	sort.Strings(b)

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}

// Переводим CamelCase в lisp-case.
func lisp_case(in string) string { // nolint: stylecheck, revive
	if in == "" {
		return ""
	}
	// проверка на капс
	if !regexp.MustCompile(`[a-z]+`).MatchString(in) {
		return in
	}

	maxexpr := 8 // максимальное количество появлений регулярки в тексте
	expr := regexp.MustCompile(`[A-Z]([a-z]+|$)`).FindAllStringIndex(in, maxexpr)

	if len(expr) == 0 {
		return in
	}

	for i := len(expr) - 1; i > 0; i-- {
		in = in[:expr[i][0]] + "-" + in[expr[i][0]:]
	}

	return strings.ToLower(in)
}

// Возвращает значения искомых полей (prefs) из указанной строки s. Порядок префиксов имеет значение
// например, имеем строку `pref1="one", pref2="two"` --> GetAttr("pref1", "pref2") --> []string{"one", "two"}.
func GetAttr(s string, prefs ...string) []string {
	res := make([]string, 0, len(prefs))
	curprefind := 0 // индекс текущего префикса

	for i := 0; i < len(s); {
		if len(res) == len(prefs) {
			break
		}
		// длина текущего префикса + 2 символа (потому что после префикса идет еще =")
		curprefle := len(prefs[curprefind]) + 2
		if i+curprefle >= len(s) {
			break
		}

		temp := make([]byte, 0, len(s[i+curprefle:]))
		// встречаем искомый префикс - начинаем запоминать символы, находящиеся внутри кавычек
		if s[i:i+curprefle] == fmt.Sprintf(`%s="`, prefs[curprefind]) {
			for j := i + curprefle; j < len(s); j++ {
				if s[j] != '"' {
					temp = append(temp, s[j])
				} else {
					// достигаем кавычек - сохраняем строку в res, переключаемся на новый тэг
					res = append(res, string(temp))
					i = j
					curprefind++
					break
				}
			}
		}
		i++
	}

	return res
}

// Возвращает значение выбранного тэга из поля с индексом index в reflect.Value структуры.
func GetTag(val reflect.Value, index int, tag string) []string {
	if tag == "" {
		return []string{}
	}
	// (начало строки либо пробел) + наш тэг + все символы до следующих кавычек
	re := regexp.MustCompile(fmt.Sprintf(`(^|\p{Z})%s:\".+?\"`, tag))
	// ищем эту регулярку в тэгах поля
	temp := re.FindString(string(val.Type().Field(index).Tag))
	if temp == "" {
		return []string{}
	}
	// если находим - удаляем префикс тэга, лишние пробелы, двоеточие и кавычки
	temp = strings.Trim(strings.TrimLeft(temp[len(tag)+1:], ": "), `"`)

	return strings.Split(temp, ",")
}

// проверяет является ли значение булевым и если необходимо меняет "yes/no/never" на "true/false".
func ConvertToFSBoolleans(in string) string {
	if in == "yes" {
		return "true"
	}

	if in == "no" || in == "never" {
		return "false"
	}

	return in
}
