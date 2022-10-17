package hostsfile

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
)

var (
	path    = "/etc/hosts"
	newline = "\n"
)

func init() {
	if runtime.GOOS == "windows" {
		path = os.ExpandEnv(filepath.FromSlash("${SystemRoot}/System32/drivers/etc/hosts"))
		newline = "\r\n"
	}
}

func AddAlias(address string, aliases ...string) error {
	if len(address) == 0 {
		return errors.New("address is required")
	}

	if len(aliases) == 0 {
		return errors.New("at least one alias is required")
	}

	lines, err := readLines()

	if err != nil {
		return err
	}

	lines = removeAliases(lines, aliases...)

	for _, alias := range aliases {
		line := fmt.Sprintf("%s %s", address, alias)
		lines = append(lines, line)
	}

	return writeLines(lines)
}

func RemoveByAlias(alias ...string) error {
	if len(alias) == 0 {
		return nil
	}

	lines, err := readLines()

	if err != nil {
		return err
	}

	lines = removeAliases(lines, alias...)

	return writeLines(lines)
}

func RemoveByAddress(address ...string) error {
	if len(address) == 0 {
		return nil
	}

	lines, err := readLines()

	if err != nil {
		return err
	}

	lines = removeAddresses(lines, address...)

	return writeLines(lines)
}

func removeAddresses(lines []string, address ...string) []string {
	result := make([]string, 0)

loop:
	for _, line := range lines {
		for _, val := range address {
			if matched, _ := regexp.MatchString(`^`+val+`\s+`, line); matched {
				continue loop
			}
		}

		result = append(result, line)
	}

	return result
}

func removeAliases(lines []string, alias ...string) []string {
	result := make([]string, 0)

loop:
	for _, line := range lines {
		for _, val := range alias {
			if matched, _ := regexp.MatchString(`\s+`+val+`\s+`, line); matched {
				continue loop
			}

			if matched, _ := regexp.MatchString(`\s+`+val+`$`, line); matched {
				continue loop
			}
		}

		result = append(result, line)
	}

	return result
}

func readLines() ([]string, error) {
	result := make([]string, 0)

	file, err := os.Open(path)

	if err != nil {
		return result, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		result = append(result, scanner.Text())
	}

	return result, scanner.Err()
}

func writeLines(lines []string) error {
	return os.WriteFile(path, []byte(strings.Join(lines, newline)), 0644)
}
