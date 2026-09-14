package cli

import "testing"

func TestHelp_Root(t *testing.T) {
	_, _, err := execute(t, "--help")
	if err != nil {
		t.Fatalf("expected --help to exit cleanly, got %v", err)
	}
}

func TestHelp_EverySubcommand(t *testing.T) {
	for _, name := range []string{"new", "add-dep", "rm-dep", "make", "rm", "rename", "show", "list"} {
		t.Run(name, func(t *testing.T) {
			_, _, err := execute(t, name, "--help")
			if err != nil {
				t.Fatalf("expected %q --help to exit cleanly, got %v", name, err)
			}
		})
	}
}

func TestUnknownCommand_Rejected(t *testing.T) {
	_, _, err := execute(t, "frobnicate")
	if err == nil {
		t.Fatal("expected an unknown-command error")
	}
}
