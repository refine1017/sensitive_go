package sensitive_go

import (
	"bufio"
	"io"
	"os"
	"strings"
	"unicode"
)

var sensitive_replacement = byte('*')

var sensitive_words = map[rune][][]rune{}

func InitSensitiveWords(filename string) error {
	swords, err := readSensitiveWords(filename)
	if err == nil {
		sensitive_words = swords
	}
	return err
}

func readSensitiveWords(filename string) (map[rune][][]rune, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	sensitive_words := make(map[rune][][]rune)
	reader := bufio.NewReader(f)
	var line string
	for {
		lstr, isPrefix, err := reader.ReadLine()
		if err != nil {
			if err != io.EOF {
				return nil, err
			} else {
				break
			}
		}
		line += string(lstr)
		if isPrefix {
		} else {
			strs := strings.Split(line, "|")
			str := strs[0]
			if len(str) <= 0 {
				continue
			}
			firstWord := []rune(str)[0]
			wordArr := make([][]rune, 0, 10)
			for i := 1; i < len(strs); i++ {
				wordArr = append(wordArr, []rune(strs[i]))
			}
			sensitive_words[firstWord] = wordArr
			line = ""
		}
	}
	return sensitive_words, nil
}

// 替换文本中的屏蔽字
func TransformSensitiveWords(words string) string {
	lenWords := len(words)
	if lenWords == 0 {
		return words
	}

	inputRunes := []rune(words)
	outBytes := make([]byte, 0, lenWords*3)
IR_S:
	for i := 0; i < len(inputRunes); i++ {
		curRune := inputRunes[i]

		// 忽略字符
		if isSensitiveFilterCharacter(unicode.ToLower(curRune)) {
			outBytes = append(outBytes, []byte(string(curRune))...)
			continue
		}

		// 非屏蔽字符
		v, found := sensitive_words[unicode.ToLower(curRune)]
		if !found {
			outBytes = append(outBytes, []byte(string(curRune))...)
			continue
		}

		if len(v) == 0 {
			outBytes = append(outBytes, sensitive_replacement)
			continue IR_S
		}

		for _, sws := range v {
			len_sws := len(sws)

			if (len(inputRunes) - 1 - i) < len_sws {
				continue
			}

			sw_indexes := make([]int, len_sws)
			for j, ri := 0, i+1; j < len_sws && ri < len(inputRunes); ri++ {
				oRune := unicode.ToLower(inputRunes[ri])
				if unicode.ToLower(sws[j]) == oRune {
					sw_indexes[j] = ri
					if j == (len_sws - 1) {
						outBytes = append(outBytes, sensitive_replacement)
						for k, l := i+1, 0; k <= ri; k++ {
							if k == sw_indexes[l] {
								outBytes = append(outBytes, sensitive_replacement)
								l++
							} else {
								outBytes = append(outBytes, []byte(string(inputRunes[k]))...)
							}
						}
						i = ri
						continue IR_S
					}
					j++
				} else {
					if isSensitiveFilterCharacter(oRune) {
						continue
					} else {
						break
					}
				}
			}
		}
		outBytes = append(outBytes, []byte(string(curRune))...)
	}
	return string(outBytes)
}

// 检查文本中是否有屏蔽字
func CheckSensitiveWords(words string) bool {
	if len(words) == 0 {
		return false
	}

	inputRunes := []rune(words)

	for i := 0; i < len(inputRunes); i++ {
		curRune := inputRunes[i]

		// 忽略字符
		if isSensitiveFilterCharacter(unicode.ToLower(curRune)) {
			continue
		}

		// 非屏蔽字符
		v, found := sensitive_words[unicode.ToLower(curRune)]
		if !found {
			continue
		}

		if len(v) == 0 {
			return true
		}

		for _, sws := range v {
			len_sws := len(sws)

			if (len(inputRunes) - 1 - i) < len_sws {
				continue
			}

			for j, ri := 0, i+1; j < len_sws && ri < len(inputRunes); ri++ {
				oRune := unicode.ToLower(inputRunes[ri])
				if unicode.ToLower(sws[j]) == oRune {
					if j == (len_sws - 1) {
						return true
					}
					j++
				} else {
					if isSensitiveFilterCharacter(oRune) {
						continue
					} else {
						break
					}
				}
			}
		}
	}

	return false
}

// 是否为可忽略的字符
func isSensitiveFilterCharacter(r rune) bool {
	i := int64(r)
	// ASCII部分
	if i < 128 {
		if (65 <= i && i <= 90) || (97 <= i && i <= 122) {
			//去除大小写字母
		} else if '[' == i || i == ']' {
			//去除有特殊含义字符(聊天中的表情)
		} else {
			return true
		}
	}
	// 中文标点符号
	if (0x2000 <= i && i <= 0x33ff) || (0xff01 <= i && i <= 0xff5e) {
		return true
	}
	return false
}
