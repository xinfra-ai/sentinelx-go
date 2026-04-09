# sentinelx-go

Whether it's an AI agent, a bank transfer, or a power grid — execution is permanent.

`sentinelx-go` enforces at the commit boundary. Before the action executes. Server-side. Unbypassable. Cryptographic receipt on every decision.

[![license](https://img.shields.io/badge/license-Apache--2.0-blue)](LICENSE)
[![go version](https://img.shields.io/badge/go-1.21+-00ADD8)](https://golang.org)

## Install

```bash
go get github.com/xinfra-ai/sentinelx-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    sentinelx "github.com/xinfra-ai/sentinelx-go"
)

func main() {
    sx := sentinelx.New("YOUR_API_KEY")

    receipt, err := sx.Enforce(context.Background(), "ai.agent.action.execute", map[string]any{
        "agent_id":               "agent-001",
        "action_type":            "file.write",
        "human_in_loop_required": true,
        "human_in_loop":          true,
        "action_within_scope":    true,
        "action_logged":          true,
    })
    if err != nil {
        if ae, ok := err.(*sentinelx.AdmissibilityError); ok {
            fmt.Println("INADMISSIBLE:", ae.Receipt.Summary)
            fmt.Println("Constraint:", *ae.Receipt.Constraint)
            fmt.Println("Receipt hash:", ae.Receipt.ReceiptHash)
            return
        }
        log.Fatal(err)
    }
    fmt.Println("ADMISSIBLE | Receipt:", receipt.ReceiptHash)
}
```

## SCADA Example

```go
receipt, err := sx.Enforce(context.Background(), "scada.setpoint.change", map[string]any{
    "device_id":                  "rtu-456",
    "parameter":                  "voltage_setpoint",
    "operator_authorized":        true,
    "change_ticket_linked":       true,
    "change_logged":              true,
    "two_person_auth":            true,
    "rollback_procedure_defined": true,
    "action_logged":              true,
})
```

## Observe Mode

```go
// Always returns receipt. Never errors on INADMISSIBLE.
// Useful for logging pipelines and observe mode.
receipt, err := sx.Evaluate(context.Background(), action, context)
fmt.Println(receipt.Verdict) // "ADMISSIBLE" or "INADMISSIBLE"
```

## How It Works

SentinelX sits at the commit boundary between your system and execution. Before any irreversible action fires, the enforcement engine evaluates it against invariant constraints and returns a deterministic verdict with a provenance receipt.

- **ADMISSIBLE** → receipt returned, action may proceed
- **INADMISSIBLE** → `AdmissibilityError` returned, nothing executes, receipt sealed

The enforcement decision is made server-side. It cannot be bypassed client-side.

## Domain Coverage

| Domain | Example Actions |
|--------|----------------|
| AI/ML Agents | `ai.agent.action.execute`, `ml.model.deploy.production` |
| Financial | `wire.transfer.execute`, `algo.trade.execute` |
| OT/SCADA | `scada.setpoint.change`, `breaker.open.execute` |
| Grid/Energy | `load.transfer.execute`, `der.curtailment.execute.batch` |
| Cyber/RMM | `rmm.script.execute`, `rmm.privilege.escalate` |
| Healthcare | `medication.order.execute`, `patient.record.modify` |
| Mobility | `driver.payout.execute`, `surge.pricing.apply` |

## Get an API Key

```bash
curl -X POST https://enforce.sentinelx.ai/generate-key
```

Or visit [sentinelx.ai](https://sentinelx.ai).

## Links

- [sentinelx.ai](https://sentinelx.ai)
- [enforce.sentinelx.ai](https://enforce.sentinelx.ai)
- [@sentinelx/sdk on npm](https://npmjs.com/package/@sentinelx/sdk)

## License

Apache-2.0# sentinelx-go

Go client for the [SentinelX Enforcement API](https://sentinelx.ai).

Pre-execution enforcement at the commit boundary. Deterministic. Server-side. No bypass.

## Install

```bash
go get github.com/xinfra-ai/sentinelx-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    sentinelx "github.com/xinfra-ai/sentinelx-go"
)

func main() {
    sx := sentinelx.New("YOUR_API_KEY")

    receipt, err := sx.Enforce(context.Background(), "ai.agent.action.execute", map[string]any{
        "agent_id":               "agent-001",
        "action_type":            "file.write",
        "human_in_loop_required": true,
        "human_in_loop":          false,
        "action_within_scope":    true,
        "action_logged":          true,
    })

    if err != nil {
        if ae, ok := err.(*sentinelx.AdmissibilityError); ok {
            fmt.Println("INADMISSIBLE:", ae.Receipt.Summary)
            fmt.Println("Receipt hash:", ae.Receipt.ReceiptHash)
            return
        }
        log.Fatal(err)
    }

    fmt.Println("ADMISSIBLE")
    fmt.Println("Receipt hash:", receipt.ReceiptHash)
}
```

## Get an API Key

```bash
curl -X POST https://enforce.sentinelx.ai/generate-key
```

Or visit [sentinelx.ai](https://sentinelx.ai).

## License

Apache-2.0
