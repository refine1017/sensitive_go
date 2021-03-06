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

func TestInitSensitiveWords(t *testing.T) {
	if e := InitSensitiveWords("no.txt"); e == nil {
		t.Errorf("InitSensitiveWords no.txt, e == nil")
	}
}

func Test_TransformSensitiveWords(t *testing.T) {
	if TransformSensitiveWords("傻逼啊啊啊啊") != "**啊啊啊啊" {
		t.Errorf("TransformSensitiveWords 傻逼啊啊啊啊 != **啊啊啊啊")
	}

	if TransformSensitiveWords("碡") != "*" {
		t.Errorf("TransformSensitiveWords 碡 != *")
	}

	if TransformSensitiveWords("ml，问题") != "**，问题" {
		t.Errorf("TransformSensitiveWords ml != **，问题")
	}

	if TransformSensitiveWords("12345") != "12345" {
		t.Errorf("TransformSensitiveWords 12345 != 12345")
	}

	if TransformSensitiveWords("测试") != "测试" {
		t.Errorf("CheckSensitiveWords 测试 != 测试")
	}

	if TransformSensitiveWords("") != "" {
		t.Errorf("CheckSensitiveWords '' != ''")
	}

	if TransformSensitiveWords("傻  B") != "*  *" {
		t.Errorf("TransformSensitiveWords 傻  B != *  *")
	}
}

func Test_CheckSensitiveWords(t *testing.T) {

	if CheckSensitiveWords("傻逼啊啊啊啊") == false {
		t.Errorf("CheckSensitiveWords 傻逼啊啊啊啊 = false")
	}

	if CheckSensitiveWords("碡") == false {
		t.Errorf("CheckSensitiveWords 碡 = false")
	}

	if CheckSensitiveWords("ml1，问题") == false {
		t.Errorf("CheckSensitiveWords ml1，问题 = false")
	}

	if CheckSensitiveWords("傻  B") == false {
		t.Errorf("CheckSensitiveWords 傻  B = false")
	}

	if CheckSensitiveWords("222") == true {
		t.Errorf("CheckSensitiveWords 222 = true")
	}

	if CheckSensitiveWords("测试") == true {
		t.Errorf("CheckSensitiveWords 测试 = true")
	}

	if CheckSensitiveWords("") == true {
		t.Errorf("CheckSensitiveWords '' = true")
	}
}
