package sip

import (
	"fmt"
	"os"
	"os/exec"
)

// Converts sip.conf to pjsip.conf. Optional args: [input path] [output path].
func ConvertToPjsip(args ...string) ([]byte, error) {
	optsNumber := 3
	opts := make([]string, 0, optsNumber)
	opts = append(opts, "./../src/sip/sip_to_pjsip.py")

	arg1 := "./../input/sip.conf"             // дефолтный путь исходного файла с конфигом
	arg2 := "./../output/pjsip_from_sip.conf" // дефолтный путь для создания pjsip.conf

	if len(args) != 0 {
		arg1 = args[0]
		if _, err := os.Stat(arg1); os.IsNotExist(err) {
			err := fmt.Errorf("the file [%s] isn't exist", arg1)
			return []byte{}, err
		}

		if len(args) > 1 {
			arg2 = args[1]
		}
	}

	opts = append(opts, arg1, arg2)

	cmd := exec.Command("python3.9", opts...)

	answer, err := cmd.Output()
	if err != nil {
		return []byte{}, err
	}

	return answer, nil
}
