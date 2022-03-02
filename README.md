# **pjsip.conf -> json и наоборот**
Программа преобразует содержимое файла pjsip.conf в формат JSON и наоборот;
Текст конфига парсится на блоки, повторяющиеся значения преобразуются в массивы. 
Предобработанный блок анмаршаллится в соответствующую структуру, заполненная структура аппендится в итоговую PJSIP структуру.

Программа очень гибкая и не привязана к названиям полей/структур и т.д.: для изменения, добавления, удаления структур достаточно отредактировать соответствующую инфу о структурах в файле pjsip.go, больше никаких правок не потребуется.

Предобработка блоков проводится с учетом следующих правил:

1) строки, начинающиеся с символа ";" игнорируются. Строки без символа "=" (например: *"Hello, world!\n"*) вызывают ошибку

2) все пробелы удаляются (кроме символа разделителя строк), все коменты удаляются (*"type = aor ;comment line\n" -> "type=aor\n"*)

3) Блок - кусок текста после выражения [XXXX] и до следующих квадратных скобок. Распознавание блоков использует значение строки "type=", если эта строка будет пустой, блок будет анмаршалиться в структуру Unknown. Программа может распознавать следующие типы блоков: system, global, endpoint, identify, aor, auth, registration, outbound, acl, phoneprov, domain, resource_list. В случае нестандартного значения поля type (например *"type=custom"*), блок будет обработан, но в виде структуры Unknown.

4) имеется возможность распознавания шаблонов/темплейтов
``` 
[mytrunk](!,defaulttrunk)
```

## Примеры:
* to json / struct
```go
    // переводим pjsip.conf в json
    func one(path string) ([]byte, error) {
        json, err := pjsip.ConfToJSONv2(path)
        if err != nil {
            return []byte{}, err
        }

        fmt.Println("pjsip.conf converted to json SUCCESSFULLY")

        return json, nil
    }

    // 2й вариант: анмаршаллим данные в структуру
    func two(data []byte, out *pjsip.PJSIP) error {
        if err := pjsip.Unmarshal(out, string(data)); err != nil {
            return err
        }
        return nil
    }
```
* json to pjsip.conf
```go
    // распаковываем структуру PJSIP в файл конфига по адресу outputPath
    func three(json []byte, outputPath string) error {
        if err := pjsip.JSONtoConfv2(json, outputPath); err != nil {
            return err
        }
        return nil
    }

    // 2й вариант:
    func four(p *pjsip.PJSIP, outputPath) error {
        out := ""
        pjsip.ReadStruct(p, "", &out)

        f, err := os.OpenFile(outputPath, os.O_CREATE|os.O_WRONLY, os.ModePerm)
        if err != nil {
            return err
        }

        _, err = f.Write([]byte(out))
        if err != nil {
            return err
        }
    }
```
# **freeswitch.xml -> pjsip json и наоборот**
## Примеры
* freeswitch.xml to pjsip
```go
    func five(path string) (*pjsip.PJSIP, error) {
        out := new(pjsip.PJSIP)

        data, err := ioutil.ReadFile(path)
        if err != nil {
            return out, err
        }

        if err = freeswitch.ToPJSIP(data, out); err != nil {
            return &pjsip.PJSIP{}, err
        }

        return out
    }
```

* pjsip.PJSIP to freeswitch.xml
```go
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

            file, err := os.OpenFile(pathToSave + "_" + fmt.Sprintln(i), os.O_CREATE|os.O_WRONLY, os.ModePerm)
            if err != nil {
                return err
            }

            _, err = file.Write(b)
            if err != nil {
                return err
            }
        }
        
    return nil
    }
```

# **sip_chan.conf -> pjsip.conf**
## to pjsip:
*ВНИМАНИЕ: для работы этой функции нужен python не ниже 3.9!*
```go
    // переводит sip_chan.conf в pjsip.conf
    // input (опционально) - место расположения sip_chan.conf
    // output (опционально) - место для создания pjsip.conf
    func seven(input string, output string) ([]byte, error) {
        answer, err := sip.ConvertToPjsip("./../input/sip.conf", "./../output/pjsip_from_sip.conf")
        if err != nil {
            return []byte{}, err
        }

        return answer, nil
    }
```
## to chan_sip.conf
```go
    // переводим pjsip.PJSIP в chan_sip.conf
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

        return nil
    }
```
# Синтаксис тэгов
## pjsip:
*`pjsip:"omit"`* - игнор поля при распаковке в pjsip.conf, единственно возможное значение
## fswitch:
*`fswitch:"field_name"`* - param name, который будет указан при распаковке во freeswitch.xml
Если тег отсутствует, поле при распаковке будет проигнорировано
## sip:
*`sip:"sip_name,option"`* - тэг используется только для распаковки pjsip в chan_sip.conf.
sip_name это имя поля, которое будет использовано при распаковке.
option:
| Значение option | Описание |
|:-----------:|:----------------------------------------------|
| hide | Означает, что поле будет обработано, но не будет показано в файле конфига (например type - влияет на логику обработки, но не выгружается |
| omit | Игнор |
| register | Обрабатываем register |

Если тег отсутствует, поле при распаковке в sip будет проигнорировано 

# TODO:
    1) тесты sip_to_pjsip направления
    2) проверить темплейты с запятой (!,options)
    4) может ли в chan_sip быть несколько блоков general c udpbindaddr? 
    5) темплейты при конвертации из pjsip в sip
    6) поменять строку на слайс строк в рекурсивном ридере pjsip
    8) добавить error в TypeRegistry.MakeStruct
    9) какой будет синтаксис строки host в chan_sip, если в pjsip/aor/contact несколько контактов?
    
