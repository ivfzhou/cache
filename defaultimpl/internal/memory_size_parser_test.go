package internal_test

import (
	"testing"

	"gitee.com/ivfzhou/cache/defaultimpl/internal"
)

func TestSzeParse(t *testing.T) {
	size, err := internal.ParseMemorySize("1KB")
	if err != nil {
		t.Fatal(err)
	}
	if size != 1024 {
		t.Fatal("size != 1024")
	}

	size, err = internal.ParseMemorySize("1mB")
	if err != nil {
		t.Fatal(err)
	}
	if size != 1024*1024 {
		t.Fatal("size != 1024*124")
	}

	size, err = internal.ParseMemorySize("0mb")
	if err != nil {
		t.Fatal(err)
	}
	if size != 0 {
		t.Fatal("size != 0")
	}

	size, err = internal.ParseMemorySize("-1mb")
	if err == nil {
		t.Fatal("err is nil")
	}

	size, err = internal.ParseMemorySize("mb")
	if err == nil {
		t.Fatal("err is nil")
	}
}
