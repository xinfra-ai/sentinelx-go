package sentinelx_test

import (
	"context"
	"testing"

	sentinelx "github.com/xinfra-ai/sentinelx-go"
)

func TestEnforceInadmissible(t *testing.T) {
	sx := sentinelx.New("sx_demo_live")

	_, err := sx.Enforce(context.Background(), "ai.agent.action.execute", map[string]any{
		"agent_id":               "agent-001",
		"action_type":            "file.write",
		"human_in_loop_required": true,
		"human_in_loop":          false,
		"action_within_scope":    true,
		"action_logged":          true,
	})

	if err == nil {
		t.Fatal("expected AdmissibilityError, got nil")
	}

	ae, ok := err.(*sentinelx.AdmissibilityError)
	if !ok {
		t.Fatalf("expected *AdmissibilityError, got %T", err)
	}

	if ae.Receipt.Verdict != "INADMISSIBLE" {
		t.Errorf("expected INADMISSIBLE, got %s", ae.Receipt.Verdict)
	}

	t.Logf("✅ INADMISSIBLE — constraint: %v, code: %v", ae.Receipt.Constraint, ae.Receipt.ViolationCode)
}
