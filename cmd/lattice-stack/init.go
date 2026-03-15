package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func initStack(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("stack name cannot be empty")
	}

	// Create project directory.
	if err := os.MkdirAll(name, 0o755); err != nil {
		return fmt.Errorf("create directory: %w", err)
	}

	// Write main.go
	mainGo := fmt.Sprintf(`package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/latticehq/latticesdk/stack"
)

func main() {
	ctx := context.Background()

	s, err := stack.New(stack.Config{
		RuntimeURL:   os.Getenv("LATTICE_RUNTIME_URL"),
		APIKey:       os.Getenv("LATTICE_API_KEY"),
		StackName:    %q,
		StackVersion: "0.1.0",
	})
	if err != nil {
		log.Fatalf("failed to create stack: %%v", err)
	}

	me, err := s.Identity.GetCurrentUser(ctx)
	if err != nil {
		log.Fatalf("failed to authenticate: %%v", err)
	}
	fmt.Printf("%%s stack running as %%s\n", %q, me.Username)

	// Add your stack logic here.
	// Available services: s.Identity, s.Authz, s.Audit, s.Budget, s.Coordination, s.Lifecycle

	if err := s.Run(ctx); err != nil {
		log.Fatalf("stack exited: %%v", err)
	}
}
`, name, name)

	if err := os.WriteFile(filepath.Join(name, "main.go"), []byte(mainGo), 0o644); err != nil {
		return fmt.Errorf("write main.go: %w", err)
	}

	// Write go.mod
	goMod := fmt.Sprintf(`module %s

go 1.22.0

require github.com/latticehq/latticesdk v0.1.0
`, name)

	if err := os.WriteFile(filepath.Join(name, "go.mod"), []byte(goMod), 0o644); err != nil {
		return fmt.Errorf("write go.mod: %w", err)
	}

	// Write .env.example
	envExample := `# Lattice Runtime Configuration
LATTICE_RUNTIME_URL=http://localhost:3000
LATTICE_API_KEY=your-api-key-here
`
	if err := os.WriteFile(filepath.Join(name, ".env.example"), []byte(envExample), 0o644); err != nil {
		return fmt.Errorf("write .env.example: %w", err)
	}

	// Write .gitignore
	gitignore := `.env
*.exe
*.test
*.out
`
	if err := os.WriteFile(filepath.Join(name, ".gitignore"), []byte(gitignore), 0o644); err != nil {
		return fmt.Errorf("write .gitignore: %w", err)
	}

	fmt.Printf("Created Department Stack: %s/\n", name)
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  cd %s\n", name)
	fmt.Println("  cp .env.example .env")
	fmt.Println("  # Edit .env with your Runtime URL and API key")
	fmt.Println("  go run .")

	return nil
}
