package encrypt

import (
	"testing"
)

func Test_Encrypt(t *testing.T) {
	input := "abcdefghijklmnopqrstuvwxyz1234567890123"
	code, err := Encode(input)
	if err != nil {
		t.Fatalf("%v", err)
	}

	output, err := Decode(code)
	if err != nil {
		t.Fatalf("%v", err)
	}

	if output != input {
		t.Error("加密解密後輸出輸入必須相等！")
	}

}
