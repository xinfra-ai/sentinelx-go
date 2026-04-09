# sentinelx-go

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
