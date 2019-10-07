package pkg

import (
	"fmt"
	"net/http"
)

func init() {
	Checks.AddCheck(basicCheck)
}

var basicCheck = &Check{
	Name: "basic-http",
	Run: func(check *Check, config Config) (success bool, err error) {
		resp, err := http.Get(fmt.Sprintf("http://%s", config.Host))
		if err != nil {
			return
		}
		if resp.StatusCode == 200 {
			success = true
		}
		return
	},
}
