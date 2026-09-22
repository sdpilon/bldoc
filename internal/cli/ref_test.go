package cli

import "testing"

func TestParseTargetRef_WholeFile(t *testing.T) {
	ref, err := parseTargetRef("README")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ref.Target != "README" || ref.Field != "" {
		t.Fatalf("got %+v", ref)
	}
}

func TestParseTargetRef_Field(t *testing.T) {
	ref, err := parseTargetRef("README@version")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if ref.Target != "README" || ref.Field != "version" {
		t.Fatalf("got %+v", ref)
	}
}

func TestParseTargetRef_Malformed(t *testing.T) {
	_, err := parseTargetRef("README@version@extra")
	if err == nil {
		t.Fatal("expected an error for a target-ref with more than one '@'")
	}
}
