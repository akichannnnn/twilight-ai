package responses

import (
	"crypto/rand"
	"fmt"
)

type streamingToolCall struct {
	id       string
	name     string
	args     string
	finished bool
}

func generateID() string {
	b := make([]byte, 12)
	rand.Read(b)
	return fmt.Sprintf("call_%x", b)
}
