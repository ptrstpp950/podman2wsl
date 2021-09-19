package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_convertArgument_Should_Escape_Windows_Paths(t *testing.T) {
	cases := []struct{ test, expected, caseName string }{
		{"c\\a\\", "c\\a\\", "Not a path shouldn't be converted"},
		{"c:\\a\\b\\c", "/mnt/c/a/b/c", "Convert to mnt"},
		{"C:\\a\\B\\c", "/mnt/c/a/B/c", "Driver letter must be small in mnt"},
		{"C:\\", "/mnt/c/", "Root should be converted"},
		{"X:\\a\\", "/mnt/x/a/", "Other driver letter then C should be supported"},
	}
	for _, c := range cases {
		actual := convertArgument(c.test)
		assert.Equal(t, c.expected, actual, c.caseName)
	}
}

func Test_convertArguments_Should_Ignore_Host_Arg(t *testing.T) {
	args := []string{"--host", "127.0.0.1", "last"}
	result := convertArguments(args)
	expected := []string{"last"}
	assert.Equal(t, result, expected)
}

func Test_convertArguments_Should_Escape_Collon_In_Env(t *testing.T) {
	expected := []string{"something", "-e", "AAA=x\\;y", "-v", "aaaa;bbb"}
	args := []string{"something", "-e", "AAA=x;y", "-v", "aaaa;bbb"}
	result := convertArguments(args)
	assert.Equal(t, result, expected)
}

func Test_convertArguments_Should_Escape_Dollar_In_Last_Arg_In_exec(t *testing.T) {
	expected := []string{"exec", "aaa$bbb", "aaa\\$bbb"}
	args := []string{"exec", "aaa$bbb", "aaa$bbb"}
	result := convertArguments(args)
	assert.Equal(t, result, expected)
}

func Test_convertArguments_Should_Use_Full_Filter_Name_In_events(t *testing.T) {
	expected := []string{"events", "--filter", "xxx"}
	args := []string{"events", "-f", "xxx"}
	result := convertArguments(args)
	assert.Equal(t, result, expected)
}

func Test_convertArguments_Should_Replace_JSON_Format(t *testing.T) {
	cases := []struct {
		test, expected []string
		caseName       string
	}{
		{
			[]string{"ps", "{{json .}}"},
			[]string{"ps", "{\"Command\":\"\\\"{{.Command}}\\\"\",\"CreatedAt\":\"{{.CreatedAt}}\",\"ID\":\"{{.ID}}\",\"Image\":\"{{.Image}}\",\"Labels\":\"{{.Labels}}\",\"LocalVolumes\":\"0\",\"Mounts\":\"{{.Mounts}}\",\"Names\":\"{{.Names}}\",\"Networks\":\"{{.Networks}}\",\"Ports\":\"{{.Ports}}\",\"RunningFor\":\"{{.RunningFor}}\",\"Size\":\"Unknown (TODO)\",\"State\":\"{{.State}}\",\"Status\":\"{{.Status}}\"}"},
			"ps JSON format conversion failed",
		},
		{
			[]string{"other", "{{json .}}"},
			[]string{"other", "table {{json .}}"},
			"JSON format conversion failed",
		},
	}
	for _, c := range cases {
		result := convertArguments(c.test)
		assert.Equal(t, c.expected, result)
	}
}

func Test_versionHack_When_Not_Version_Should_Return_False(t *testing.T) {
	args := []string{"podman", "not_version", "--format", "{{.Server.Os}}"}
	assert.False(t, versionHack(args))
}

func Test_versionHack_When_System_Os_Return_Linux(t *testing.T) {
	args := []string{"podman", "version", "--format", "{{.Server.Os}}"}
	assert.True(t, versionHack(args))
}

func Test_versionHack_When_Client_Version_Return_Static_Text(t *testing.T) {
	args := []string{"podman", "version", "--format", "{{.Client.Version}}"}
	assert.True(t, versionHack(args))
}
