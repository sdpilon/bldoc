package manifest

import "testing"

func TestSplitRecordField(t *testing.T) {
	cases := []struct {
		name       string
		field      string
		wantRecord string
		wantField  string
		wantOK     bool
	}{
		{"record only", "bldoc", "bldoc", "", true},
		{"record and field", "bldoc.description", "bldoc", "description", true},
		{"empty", "", "", "", false},
		{"over-qualified", "bldoc.description.extra", "", "", false},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			record, field, ok := SplitRecordField(c.field)
			if ok != c.wantOK {
				t.Fatalf("ok = %v, want %v", ok, c.wantOK)
			}
			if !ok {
				return
			}
			if record != c.wantRecord || field != c.wantField {
				t.Fatalf("got record=%q field=%q, want record=%q field=%q", record, field, c.wantRecord, c.wantField)
			}
		})
	}
}
