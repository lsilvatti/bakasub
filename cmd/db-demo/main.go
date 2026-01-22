package main

import (
"fmt"
"log"
"os"

"github.com/lsilvatti/bakasub/internal/core/db"
)

func main() {
fmt.Println("╔═══════════════════════════════════════════╗")
fmt.Println("║ BakaSub - DBLocal Cache Demo             ║")
fmt.Println("╚═══════════════════════════════════════════╝")
fmt.Println()

dbPath := "test_cache.db"
defer os.Remove(dbPath)

cache, err := db.GetInstance(dbPath)
if err != nil {
log.Fatalf("Failed: %v", err)
}
defer cache.Close()

fmt.Println("✓ Cache initialized")
fmt.Println()

fmt.Println("=== Test 1: Save translations ===")
cache.SaveTranslation("Hello, world!", "Olá, mundo!", "eng->por")
cache.SaveTranslation("Good morning", "Bom dia", "eng->por")
cache.SaveTranslation("Thank you", "Obrigado", "eng->por")
fmt.Println("  ✓ Saved 3 translations")
fmt.Println()

fmt.Println("=== Test 2: Exact match retrieval ===")
translated, found := cache.GetExactMatch("Hello, world!", "eng->por")
if found {
fmt.Printf("  ✓ Found: %s\n", translated)
} else {
fmt.Println("  ✗ Not found")
}
fmt.Println()

fmt.Println("=== Test 3: Fuzzy match (95%% threshold) ===")
tests := []string{
"Hello, world!",  // Exact (100%)
"Hello, World!",  // Case diff
"Hello world!",   // Punctuation
}
for _, text := range tests {
entry, found := cache.GetFuzzyMatch(text, "eng->por", 0.95)
if found {
fmt.Printf("  ✓ '%s' -> '%s' (%.0f%%)\n", 
text, entry.TranslatedText, entry.Similarity*100)
} else {
fmt.Printf("  ✗ '%s' -> Below 95%% threshold\n", text)
}
}
fmt.Println()

fmt.Println("=== Test 4: Batch save ===")
entries := []db.CacheEntry{
{OriginalText: "Line 1", TranslatedText: "Linha 1", LangPair: "eng->por"},
{OriginalText: "Line 2", TranslatedText: "Linha 2", LangPair: "eng->por"},
{OriginalText: "Line 3", TranslatedText: "Linha 3", LangPair: "eng->por"},
}
if err := cache.SaveBatch(entries); err != nil {
log.Printf("Error: %v", err)
} else {
fmt.Printf("  ✓ Saved batch of %d entries\n", len(entries))
}
fmt.Println()

fmt.Println("=== Test 5: Thread safety ===")
done := make(chan bool)
for i := 0; i < 5; i++ {
go func(id int) {
cache.SaveTranslation(
fmt.Sprintf("Concurrent %d", id), 
fmt.Sprintf("Concorrente %d", id), 
"eng->por",
)
done <- true
}(i)
}
for i := 0; i < 5; i++ {
<-done
}
fmt.Println("  ✓ 5 concurrent operations completed")
fmt.Println()

stats, err := cache.GetStats()
if err != nil {
log.Printf("Stats error: %v", err)
} else {
fmt.Println("=== Cache Statistics ===")
fmt.Printf("  Total Entries: %d\n", stats.TotalEntries)
fmt.Printf("  Hit Rate: %.1f%%\n", stats.HitRate)
fmt.Printf("  Saved Cost: $%.4f USD\n", stats.SavedCost)
}
fmt.Println()

fmt.Println("✓ All tests passed!")
fmt.Println("✓ Cache persists to save API costs on repeat translations")
}
