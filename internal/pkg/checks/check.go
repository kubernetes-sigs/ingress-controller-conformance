package checks

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Config struct {
}

type Check struct {
	Name string

	Run func(check *Check, config Config) (bool, error)

	// All checks
	checks []*Check

	// Parent check
	parent *Check
}

type CapturedRequest struct {
	StatusCode int
	TestId     string
	Path       string
	Host       string
}

func captureRequest(location string, hostOverride string) (data CapturedRequest, err error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", location, nil)
	if err != nil {
		return
	}
	if hostOverride != "" {
		req.Host = hostOverride
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()
	if err != nil {
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return
	}

	data.StatusCode = resp.StatusCode
	return
}

type assertionSet []error
type assert struct {
	expect        interface{}
	actual        interface{}
	errorTemplate string
}

func (a *assertionSet) equals(assert assert) {
	if assert.expect != assert.actual {
		err := fmt.Errorf(assert.errorTemplate, assert.actual, assert.expect)
		*a = append(*a, err)
	}
}

func (a *assertionSet) Error() (err string) {
	for i, e := range *a {
		err += fmt.Sprintf("\t%d) Assertion failed: %s\n", i+1, e.Error())
	}
	return
}

func (c *Check) AddCheck(checks ...*Check) {
	for i, x := range checks {
		if checks[i] == c {
			panic("Checks can't be a child of itself")
		}
		checks[i].parent = c
		c.checks = append(c.checks, x)
	}
}

var Checks = &Check{
	Name: "all",
}

func (c Check) Verify(filterOnCheckName string, config Config) (successCount int, failureCount int, err error) {
	if filterOnCheckName != c.Name && filterOnCheckName != "" {
		for _, check := range c.checks {
			s, f, err := check.Verify(filterOnCheckName, config)
			successCount += s
			failureCount += f
			if err != nil {
				fmt.Printf(err.Error())
			}
		}

		return
	}

	fmt.Printf("Running %s verifications...\n", c.Name)
	if c.Run != nil {
		success, err := c.Run(&c, config)
		if err != nil {
			fmt.Printf("  %s\n", err.Error())
		}

		if success {
			successCount++
			fmt.Printf("  Check passed: %s\n", c.Name)
		} else {
			failureCount++
			fmt.Printf("  Check failed: %s\n", c.Name)
		}
	}

	for _, check := range c.checks {
		s, f, err := check.Verify("", config)
		if err != nil {
			fmt.Printf(err.Error())
		}
		successCount += s
		failureCount += f
	}
	return
}
