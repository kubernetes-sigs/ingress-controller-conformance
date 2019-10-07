package pkg

import (
	"fmt"
)

type Config struct {
	Host string
}

type Check struct {
	Name string

	Run func(check *Check, config Config) (bool, error)

	// All checks
	checks []*Check

	// Parent check
	parent *Check
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
			fmt.Printf(err.Error())
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
