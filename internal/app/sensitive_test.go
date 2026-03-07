package app

import (
	"bytes"
	"errors"
	"io"
	"strconv"
	"strings"
	"testing"
)

// --- Test structs ---

type testStrJSONTag struct {
	FieldName string `json:"json_tag_value" sensitive:"true" prompt:"Prompt Tag Value"`
}

type testIntJSONTag struct {
	FieldName int `json:"json_tag_value" sensitive:"true" prompt:"Prompt Tag Value"`
}

type testStrNoJSONTag struct {
	FieldName string `sensitive:"true" prompt:"Prompt Tag Value"`
}

type testStrNoPromptTag struct {
	FieldName string `json:"json_tag_value" sensitive:"true"`
}

type testSliceJSONTag struct {
	FieldName []string `json:"json_tag_value" sensitive:"true"`
}

// --- Env var tests ---

func TestPromptSensitiveFields_EnvJSONTagExactMatch(t *testing.T) {
	t.Setenv("STC_json_tag_value", "val123")
	s := new(testStrJSONTag)
	if err := promptSensitiveFields(s, strings.NewReader(""), io.Discard); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.FieldName != "val123" {
		t.Errorf("expected %q, got %q", "val123", s.FieldName)
	}
}

func TestPromptSensitiveFields_EnvJSONTagCaseInsensitiveMatch(t *testing.T) {
	t.Setenv("STC_JSON_TAG_VALUE", "val123")
	s := new(testStrJSONTag)
	if err := promptSensitiveFields(s, strings.NewReader(""), io.Discard); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.FieldName != "val123" {
		t.Errorf("expected %q, got %q", "val123", s.FieldName)
	}
}

func TestPromptSensitiveFields_EnvEmptyValueReturnsError(t *testing.T) {
	t.Setenv("STC_json_tag_value", "")
	s := new(testStrJSONTag)
	err := promptSensitiveFields(s, strings.NewReader(""), io.Discard)
	if !errors.Is(err, errEmptyEnvVar) {
		t.Fatalf("expected errEmptyEnvVar, got %v", err)
	}
}

func TestPromptSensitiveFields_EnvInvalidValueForIntReturnsError(t *testing.T) {
	t.Setenv("STC_json_tag_value", "val123")
	s := new(testIntJSONTag)
	err := promptSensitiveFields(s, strings.NewReader(""), io.Discard)
	if !errors.Is(err, strconv.ErrSyntax) {
		t.Fatalf("expected strconv.ErrSyntax, got %v", err)
	}
}

func TestPromptSensitiveFields_EnvFieldNameExactMatch(t *testing.T) {
	t.Setenv("STC_FieldName", "val123")
	s := new(testStrNoJSONTag)
	if err := promptSensitiveFields(s, strings.NewReader(""), io.Discard); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.FieldName != "val123" {
		t.Errorf("expected %q, got %q", "val123", s.FieldName)
	}
}

func TestPromptSensitiveFields_EnvFieldNameCaseInsensitiveMatch(t *testing.T) {
	t.Setenv("STC_FIELDNAME", "val123")
	s := new(testStrNoJSONTag)
	if err := promptSensitiveFields(s, strings.NewReader(""), io.Discard); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if s.FieldName != "val123" {
		t.Errorf("expected %q, got %q", "val123", s.FieldName)
	}
}

// --- Prompt tests ---

func TestPromptSensitiveFields_PromptUsesPromptTag(t *testing.T) {
	var out bytes.Buffer
	s := new(testStrJSONTag)
	if err := promptSensitiveFields(s, strings.NewReader("val123\n"), &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "Prompt Tag Value") {
		t.Errorf("expected prompt to contain %q, got %q", "Prompt Tag Value", out.String())
	}
	if s.FieldName != "val123" {
		t.Errorf("expected field value %q, got %q", "val123", s.FieldName)
	}
}

func TestPromptSensitiveFields_PromptUsesFieldNameWhenNoPromptTag(t *testing.T) {
	var out bytes.Buffer
	s := new(testStrNoPromptTag)
	if err := promptSensitiveFields(s, strings.NewReader("val123\n"), &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "FieldName") {
		t.Errorf("expected prompt to contain %q, got %q", "FieldName", out.String())
	}
	if s.FieldName != "val123" {
		t.Errorf("expected field value %q, got %q", "val123", s.FieldName)
	}
}

func TestPromptSensitiveFields_PromptNoJSONTagUsesPromptTag(t *testing.T) {
	var out bytes.Buffer
	s := new(testStrNoJSONTag)
	if err := promptSensitiveFields(s, strings.NewReader("val123\n"), &out); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !strings.Contains(out.String(), "Prompt Tag Value") {
		t.Errorf("expected prompt to contain %q, got %q", "Prompt Tag Value", out.String())
	}
	if s.FieldName != "val123" {
		t.Errorf("expected field value %q, got %q", "val123", s.FieldName)
	}
}

func TestPromptSensitiveFields_PromptInvalidValueForIntReturnsError(t *testing.T) {
	s := new(testIntJSONTag)
	err := promptSensitiveFields(s, strings.NewReader("val123\n"), io.Discard)
	if !errors.Is(err, strconv.ErrSyntax) {
		t.Fatalf("expected strconv.ErrSyntax, got %v", err)
	}
}

// --- Other error cases ---

func TestPromptSensitiveFields_NotPointerReturnsError(t *testing.T) {
	s := testStrJSONTag{}
	err := promptSensitiveFields(s, strings.NewReader(""), io.Discard)
	if !errors.Is(err, errNotPointer) {
		t.Fatalf("expected errNotPointer, got %v", err)
	}
}

func TestPromptSensitiveFields_UnsupportedTypeReturnsError(t *testing.T) {
	t.Setenv("STC_json_tag_value", "val123")
	s := new(testSliceJSONTag)
	err := promptSensitiveFields(s, strings.NewReader(""), io.Discard)
	if !errors.Is(err, errUnsupportedType) {
		t.Fatalf("expected errUnsupportedType, got %v", err)
	}
}

func TestPromptSensitiveFields_EmptyPromptInputReturnsError(t *testing.T) {
	s := new(testStrJSONTag)
	err := promptSensitiveFields(s, strings.NewReader("\n"), io.Discard)
	if !errors.Is(err, errEmptyInput) {
		t.Fatalf("expected errEmptyInput, got %v", err)
	}
}
