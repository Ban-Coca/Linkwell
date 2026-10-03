package shortener

import "testing"

func TestGenerateCodeFormat(t *testing.T) {
	for i := 0; i < 100; i++ {
		code, err := GenerateCode()
		if err != nil {
			t.Fatalf("GenerateCode() returned an error: %v", err)
		}
		if len(code) != codeLength {
			t.Fatalf("generated code length = %d, want %d", len(code), codeLength)
		}
		for _, character := range code {
			if !containsCharacter(alphabet, character) {
				t.Fatalf("generated code contains invalid character %q", character)
			}
		}
	}
}

func containsCharacter(value string, want rune) bool {
	for _, character := range value {
		if character == want {
			return true
		}
	}
	return false
}
