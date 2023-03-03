package util

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

func ReadEnv() error {
	var (
		fileBytes []byte
		err       error
	)
	fileBytes, err = ioutil.ReadFile(".env")
	if err != nil {
		return errors.New("Failed to open .env, check the file's name or it's presence")
	}

	var newLineSplit []string = strings.Split(string(fileBytes), "\n")

	for i := range newLineSplit {
		var env = strings.SplitN(newLineSplit[i], "=", 2)
		if len(env) == 2 {
			if err = os.Setenv(env[0], env[1]); err != nil {
				return errors.New(fmt.Sprintf("Error setting env var, err: %s", err.Error()))
			}
		}
	}

	return nil
}
