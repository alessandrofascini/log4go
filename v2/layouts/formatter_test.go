package layouts_test

import (
	"fmt"
	"os"
	"reflect"
	"sort"
	"testing"
)

const (
	keyMask = "*****"
)

type AuthInfo struct {
	Login  string // Login user
	ACL    uint   // ACL bitmask
	APIKey string // API key
}

var authInfoFields []string

func init() {
	typ := reflect.TypeOf(AuthInfo{})
	authInfoFields = make([]string, typ.NumField())
	for i := 0; i < typ.NumField(); i++ {
		authInfoFields[i] = typ.Field(i).Name
	}
	sort.Strings(authInfoFields) // People are better with sorted data
}

// String implements Stringer interface
func (ai *AuthInfo) String() string {
	key := ai.APIKey
	if key != "" {
		key = keyMask
	}
	return fmt.Sprintf("Login:%s, ACL:%08b, APIKey: %s", ai.Login, ai.ACL, key)
}

func (ai *AuthInfo) Format(state fmt.State, verb rune) {
	switch verb {
	case 'q':
		if precision, ok := state.Precision(); ok {
			fmt.Println("precision:", precision)
		}
		if width, ok := state.Width(); ok {
			fmt.Println("width:", width)
		}
		if ok := state.Flag('?'); ok {
			fmt.Println("# flag present")
		}
		fallthrough
	case 's':
		val := ai.String()
		if verb == 'q' {
			val = fmt.Sprintf("%q", val)
		}
		fmt.Fprint(state, val)
	case 'v':
		if state.Flag('#') {
			// Emit type before
			fmt.Fprintf(state, "%T", ai)
		}
		fmt.Fprint(state, "{")
		val := reflect.ValueOf(*ai)
		for i, name := range authInfoFields {
			if state.Flag('#') || state.Flag('+') {
				fmt.Fprintf(state, "%s:", name)
			}
			fld := val.FieldByName(name)
			if name == "APIKey" && fld.Len() > 0 {
				fmt.Fprint(state, keyMask)
			} else {
				fmt.Fprint(state, fld)
			}
			if i < len(authInfoFields)-1 {
				fmt.Fprint(state, " ")
			}
		}
		fmt.Fprint(state, "}")
	}

}

func TestFormatter(t *testing.T) {
	ai := &AuthInfo{
		Login:  "daffy",
		ACL:    777,
		APIKey: "duck season",
	}
	fmt.Println(ai.String())
	fmt.Printf("ai %%s: %s\n", ai)
	fmt.Printf("ai %%q: %5.10q{helo}\n", ai)
	fmt.Printf("ai %%v: %#v\n", ai)
}

func TestGetPid(t *testing.T) {
	fmt.Println(os.Getpid())
	fmt.Println(os.Getppid())
}
