package main

import "testing"

func TestConvertArgument(t *testing.T) {
	cases := []struct{ test, expected, caseName string }{
		{"c\\a\\", "c\\a\\", "Not a path shouldn't be converted"},
		{"c:\\a\\b\\c", "/mnt/c/a/b/c", "Convert to mnt"},
		{"C:\\a\\B\\c", "/mnt/c/a/B/c", "Driver letter must be small in mnt"},
		{"C:\\", "/mnt/c/", "Root should be converted"},
		{"X:\\a\\", "/mnt/x/a/", "Other driver letter then C should be supported"},
		{"string ; with ;;colons", "string \\; with \\;\\;colons", "All colons should be escaped"},
	}
	//expected := []string{"c\\a\\b", "/mnt/c/a/b/c", "/mnt/D/bbbbb", "-'/mnt/c/q'", ""}
	//test := []string{"c\\a\\b", "c:\\a\\b\\c", "D:\\bbbbb", "-'c:\\q'", `"C:\projects\aspnet-core-test"`}
	/*actual := convertPath(test)
	for index := 0; index < len(test); index++ {
		if actual[index] != expected[index] {
			t.Error("Test failed:" + actual[index] + "!=" + expected[index])
		}
	}*/
	for _, c := range cases {
		actual := convertArgument(c.test)
		if actual != c.expected {
			t.Error("Test " + c.caseName + " failed:" + actual + "!=" + c.expected)
		}
	}
}
