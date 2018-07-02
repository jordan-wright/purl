package main

import (
	"net/url"
	"testing"
)

func TestPrintHostname(t *testing.T) {
	expected := "www.example.com"
	u := URL{}
	u.Host = expected
	got := u.String()
	if got != expected {
		t.Fatalf("Unexpected value returned. Expected %s Got %s", expected, got)
	}
}

func TestPrintScheme(t *testing.T) {
	expected := "https"
	u := URL{}
	u.Scheme = expected
	got := u.String()
	if got != expected {
		t.Fatalf("Unexpected value returned. Expected %s Got %s", expected, got)
	}
}

func TestPrintHostAndScheme(t *testing.T) {
	expected := "https://www.example.com"
	u := URL{}
	u.Scheme = "https"
	u.Host = "www.example.com"
	got := u.String()
	if got != expected {
		t.Fatalf("Unexpected value returned. Expected %s Got %s", expected, got)
	}
}

func TestPrintHostAndUser(t *testing.T) {
	expected := "foo:bar@www.example.com"
	u := URL{}
	ui := url.UserPassword("foo", "bar")
	u.User = ui
	u.Host = "www.example.com"
	got := u.String()
	if got != expected {
		t.Fatalf("Unexpected value returned. Expected %s Got %s", expected, got)
	}
}

func TestPrintUser(t *testing.T) {
	expected := "foo:bar"
	u := URL{}
	ui := url.UserPassword("foo", "bar")
	u.User = ui
	got := u.String()
	if got != expected {
		t.Fatalf("Unexpected value returned. Expected %s Got %s", expected, got)
	}
}
