package flagf

import (
	"fmt"
	"io"
	"strings"

	"github.com/google/shlex"
)

type FFlag struct {
	Values []string
	BFlags map[string]*bool
	SFlags map[string]*string

	FlagTypes map[string]string
}

func (ff *FFlag) Init(values []string, sflags map[string]*string, bflags map[string]*bool, types map[string]string) {
	ff.Values = values
	ff.BFlags = bflags
	ff.SFlags = sflags
	ff.FlagTypes = types
}

func (ff *FFlag) Bool(symbol string, dvalue bool, description string) *bool {
	ff.FlagTypes[symbol] = "bool"
	ff.BFlags[symbol] = &dvalue
	return ff.BFlags[symbol]
}

func (ff *FFlag) String(symbol string, dvalue string, description string) *string {
	ff.FlagTypes[symbol] = "string"
	ff.SFlags[symbol] = &dvalue
	return ff.SFlags[symbol]
}

func (ff *FFlag) stringToFFlags(str string) {
	r := strings.NewReader(str)
	l := shlex.NewLexer(r)

	var stopFlags = false
	var nextVal = false
	var currentFlag = ""
	var firstToken = true

	for {
		token, err := l.Next()
		if err == io.EOF {
			break
		}

		if nextVal {
			*ff.SFlags[currentFlag] = token
			nextVal = false
		} else if !stopFlags && isFlag(token) {
			runes := []rune(token)
			flag := string(runes[1:])

			if ff.FlagTypes[flag] == "string" {
				currentFlag = flag
				nextVal = true
			} else if ff.FlagTypes[flag] == "bool" {
				*ff.BFlags[flag] = true
			} else {
				if firstToken {
					panic(fmt.Sprintf("flagf provided but not defined: %s", token))
				}

				ff.Values = append(ff.Values, token)
				stopFlags = true
			}

			firstToken = false

		} else {
			ff.Values = append(ff.Values, token)
			stopFlags = true
		}
	}
}

func isFlag(str string) bool {
	return strings.HasPrefix(str, "-")
}

func (ff *FFlag) Parse(fakeArgs string) {
	ff.stringToFFlags(fakeArgs)
}

func (ff *FFlag) Args() []string {
	return ff.Values
}
