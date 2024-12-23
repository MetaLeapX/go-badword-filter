package helper

import (
	"math/rand"
)

func CreateString(replaceChars string, stringLength int) string {
	lenReplaceString := len(replaceChars)
	if lenReplaceString == 0 {
		return ""
	}
	runeString := []rune(replaceChars)
	chars := make([]rune, stringLength)
	if lenReplaceString > 1 {
		for ii := 0; ii < stringLength; ii++ {
			chars[ii] = runeString[rand.Intn(lenReplaceString)]
		}
	} else {
		for i := range chars {
			chars[i] = runeString[0]
		}
	}
	return string(chars)
}
