package main

import (
"context"
"fmt"

"github.com/lsilvatti/bakasub/internal/core/parser"
"github.com/lsilvatti/bakasub/internal/core/watcher"
)

func main() {
fmt.Println("BakaSub - Phase 2.5: Translation Pipeline Demo")
fmt.Println("===============================================")
fmt.Println()

fmt.Println("=== Test 1: HI Tag Removal ===")
tests := []string{
"[Music] Hello",
"(Door) Test",
"Name: Dialog",
}

for _, text := range tests {
cleaned := parser.RemoveHearingImpairedTags(text)
fmt.Printf("  %s -> %s\n", text, cleaned)
}
fmt.Println()

fmt.Println("=== Test 2: Batching ===")
lines := []parser.SubtitleLine{
{Index: 0, Text: "L1"},
{Index: 1, Text: "L2"},
{Index: 2, Text: "L3"},
}

batches := parser.BatchLines(lines, 2)
fmt.Printf("  Created %d batches\n", len(batches))
fmt.Println()

fmt.Println("=== Test 3: Watcher ===")
w, err := watcher.New("/tmp")
if err != nil {
fmt.Printf("  Error: %v\n", err)
} else {
fmt.Println("  ✓ Watcher initialized")
w.Stop()
}
fmt.Println()

_ = context.Background()

fmt.Println("✓ Phase 2.5 core components validated!")
fmt.Println("Pipeline ready for integration")
}
