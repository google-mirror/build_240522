package canoninja

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"os"
)

var (
	rulePrefix  = []byte("rule ")
	buildPrefix = []byte("build ")
	phonyRule   = []byte("phony")
)

func Generate(path string) error {
	var buffer []byte
	var err error
	if buffer, err = os.ReadFile(path); err != nil {
		return err
	}

	// Break file into lines
	from := 0
	var lines [][]byte
	for from < len(buffer) {
		line := getLine(buffer[from:])
		lines = append(lines, line)
		from += len(line)
	}

	// FOr each rule, calculate and remember its digest
	ruleDigest := make(map[string]string)
	for i := 0; i < len(lines); {
		if bytes.HasPrefix(lines[i], rulePrefix) {
			// Find ruleName
			rn := ruleName(lines[i])
			if len(rn) == 0 {
				return fmt.Errorf("%s:%d: rule name is missing or on the next line", path, i+1)
			}
			sRuleName := string(rn)
			if _, ok := ruleDigest[sRuleName]; ok {
				return fmt.Errorf("%s:%d: the rule %s has been already defined", path, i+1, sRuleName)
			}
			// Calculate rule text digest as a digests of line digests.
			var digests []byte
			doDigest := func(b []byte) {
				h := sha1.New()
				h.Write(lines[i])
				digests = h.Sum(digests)

			}
			doDigest(lines[i])
			for i++; i < len(lines) && lines[i][0] == ' '; i++ {
				doDigest(lines[i])
			}
			h := sha1.New()
			h.Write(digests)
			ruleDigest[sRuleName] = "R" + hex.EncodeToString(h.Sum(nil))

		} else {
			i++
		}
	}

	// Rewrite rule names.
	for i, line := range lines {
		if bytes.HasPrefix(line, buildPrefix) {
			brn := getBuildRuleName(line)
			if bytes.Equal(brn, phonyRule) {
				os.Stdout.Write(line)
				continue
			}
			if len(brn) == 0 {
				return fmt.Errorf("%s:%d: build statement lacks rule name", path, i+1)
			}
			os.Stdout.Write(line[0 : cap(line)-cap(brn)])
			if digest, ok := ruleDigest[string(brn)]; ok {
				os.Stdout.WriteString(digest)
			} else {
				return fmt.Errorf("%s:%d: no rule for this build target", path, i+1)
			}
			os.Stdout.Write(line[cap(line)+len(brn)-cap(brn):])
		} else if bytes.HasPrefix(line, rulePrefix) {
			rn := ruleName(line)
			// Write everything before it
			os.Stdout.Write(line[0 : cap(line)-cap(rn)])
			os.Stdout.WriteString(ruleDigest[string(rn)])
			//goland:noinspection GoUnhandledErrorResult
			os.Stdout.Write(line[cap(line)+len(rn)-cap(rn):])
		} else {
			//goland:noinspection GoUnhandledErrorResult
			os.Stdout.Write(line)
		}
	}
	return nil
}

func getLine(b []byte) []byte {
	if n := bytes.IndexByte(b, '\n'); n >= 0 {
		return b[:n+1]
	}
	return b
}

// Returns build statement's rule name
func getBuildRuleName(line []byte) []byte {
	n := bytes.IndexByte(line, ':')
	if n <= 0 {
		return nil
	}
	ruleName := line[n+1:]
	if ruleName[0] == ' ' {
		ruleName = bytes.TrimLeft(ruleName, " ")
	}
	if n := bytes.IndexAny(ruleName, " \t\r\n"); n >= 0 {
		ruleName = ruleName[0:n]
	}
	return ruleName
}

// Returns rule statement's rule name
func ruleName(lineAfterRule []byte) []byte {
	ruleName := lineAfterRule[len(rulePrefix):]
	if len(ruleName) == 0 {
		return ruleName
	}
	if ruleName[0] == ' ' {
		ruleName = bytes.TrimLeft(ruleName, " ")
	}
	if n := bytes.IndexAny(ruleName, " \t\r\n"); n >= 0 {
		ruleName = ruleName[0:n]
	}
	return ruleName
}
