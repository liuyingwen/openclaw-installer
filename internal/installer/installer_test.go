package installer

import (
	"errors"
	"reflect"
	"testing"

	"github.com/liuyingwen/openclaw-installer/internal/config"
	"github.com/liuyingwen/openclaw-installer/internal/repair"
)

type stubRunner struct {
	executed []string
	failures map[string]error
}

func (r *stubRunner) Run(command string) error {
	r.executed = append(r.executed, command)
	if err := r.failures[command]; err != nil {
		return err
	}
	return nil
}

func TestDoctorModeReturnsPlanWithoutExecutingCommands(t *testing.T) {
	runner := &stubRunner{failures: map[string]error{}}
	workflow := Workflow{
		Runner: runner,
	}

	cfg := config.Config{
		Install: []config.CommandStep{{Name: "install", Run: "install-command"}},
		Verify:  []config.CommandStep{{Name: "verify", Run: "verify-command"}},
	}

	result, err := workflow.Doctor(cfg, repair.Plan{Commands: []string{"repair-command"}})
	if err != nil {
		t.Fatalf("expected doctor to succeed, got %v", err)
	}

	if len(runner.executed) != 0 {
		t.Fatalf("expected no executed commands, got %#v", runner.executed)
	}

	expected := []string{"repair-command", "install-command", "verify-command"}
	if !reflect.DeepEqual(result.PlannedCommands, expected) {
		t.Fatalf("expected %#v, got %#v", expected, result.PlannedCommands)
	}
}

func TestInstallExecutesRepairBeforeInstallAndVerify(t *testing.T) {
	runner := &stubRunner{failures: map[string]error{}}
	workflow := Workflow{
		Runner: runner,
	}

	cfg := config.Config{
		Install: []config.CommandStep{{Name: "install", Run: "install-command"}},
		Verify:  []config.CommandStep{{Name: "verify", Run: "verify-command"}},
	}

	_, err := workflow.Install(cfg, repair.Plan{Commands: []string{"repair-command"}})
	if err != nil {
		t.Fatalf("expected install to succeed, got %v", err)
	}

	expected := []string{"repair-command", "install-command", "verify-command"}
	if !reflect.DeepEqual(runner.executed, expected) {
		t.Fatalf("expected %#v, got %#v", expected, runner.executed)
	}
}

func TestInstallStopsOnVerificationFailure(t *testing.T) {
	runner := &stubRunner{
		failures: map[string]error{
			"verify-command": errors.New("binary missing"),
		},
	}
	workflow := Workflow{
		Runner: runner,
	}

	cfg := config.Config{
		Install: []config.CommandStep{{Name: "install", Run: "install-command"}},
		Verify:  []config.CommandStep{{Name: "verify", Run: "verify-command"}},
	}

	_, err := workflow.Install(cfg, repair.Plan{})
	if err == nil {
		t.Fatal("expected verification error")
	}
}
