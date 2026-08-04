package kpi

import "testing"

func TestPctReturnsNilWhenDenominatorIsZero(t *testing.T) {
	if got := pct(5, 0); got != nil {
		t.Fatalf("expected nil, got %v", *got)
	}
}

func TestPctReturnsPercentage(t *testing.T) {
	got := pct(25, 50)
	if got == nil {
		t.Fatal("expected percentage value, got nil")
	}
	if *got != 50 {
		t.Fatalf("expected 50, got %v", *got)
	}
}

func TestNormalizeRoleCode(t *testing.T) {
	if got := normalizeRoleCode("sales"); got != "sales_rep" {
		t.Fatalf("expected sales_rep, got %s", got)
	}
	if got := normalizeRoleCode("sales_manager"); got != "sales_manager" {
		t.Fatalf("expected sales_manager, got %s", got)
	}
}
