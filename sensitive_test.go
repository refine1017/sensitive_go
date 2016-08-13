package sensitive_go

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
)

func substr(s string, pos, length int) string {
	runes := []rune(s)
	l := pos + length
	if l > len(runes) {
		l = len(runes)
	}
	return string(runes[pos:l])
}

func init() {
	_, fulleFilename, _, _ := runtime.Caller(0)
	path := substr(fulleFilename, 0, strings.LastIndex(fulleFilename, "/"))
	path = fmt.Sprintf("%v%v", path, "/library.txt")
	InitSensitiveWords(path)
}

func Test_isSensitiveFilterCharacter(t *testing.T) {
	var r = rune(123)
	if isSensitiveFilterCharacter(r) != true {
		t.Errorf("isSensitiveFilterCharacter rune(123) != true")
	}

	r = rune('B')
	if isSensitiveFilterCharacter(r) == true {
		t.Errorf("isSensitiveFilterCharacter rune(59) = true")
	}
}

func Test_TransformSensitiveWords(t *testing.T) {
	if TransformSensitiveWords("傻逼啊啊啊啊") != "**啊啊啊啊" {
		t.Errorf("TransformSensitiveWords 傻逼啊啊啊啊 != **啊啊啊啊")
	}

	if TransformSensitiveWords("碡") != "*" {
		t.Errorf("TransformSensitiveWords 碡 != *")
	}

	if TransformSensitiveWords("ml") != "**" {
		t.Errorf("TransformSensitiveWords ml != **")
	}

	if TransformSensitiveWords("12345") != "12345" {
		t.Errorf("TransformSensitiveWords 12345 != 12345")
	}
}

func Test_CheckSensitiveWords(t *testing.T) {

	if CheckSensitiveWords("傻逼啊啊啊啊") == false {
		t.Errorf("CheckSensitiveWords 傻逼啊啊啊啊 = false")
	}

	if CheckSensitiveWords("碡") == false {
		t.Errorf("CheckSensitiveWords 碡 = false")
	}

	if CheckSensitiveWords("ml1") == false {
		t.Errorf("CheckSensitiveWords ml1 = false")
	}

	if CheckSensitiveWords("222") == true {
		t.Errorf("CheckSensitiveWords 222 = true")
	}
}
