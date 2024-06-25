package test

import (
	"fmt"
	"net/url"
	"testing"
)

func TestGatLab(t *testing.T) {
	originalString := "/ok/tt.txt"

	// 使用QueryEscape进行编码
	encodedString := url.QueryEscape(originalString)

	fmt.Printf("Original String: %s\n", originalString)
	fmt.Printf("encodedString String: %s\n", encodedString)
}
