package luhn

import "testing"

func TestIsValid(t *testing.T) {
	tests := []struct {
		name   string
		number string
		want   bool
	}{
		{
			name:   "valid luhn number",
			number: "79927398713",
			want:   true,
		},
		{
			name:   "valid luhn number 2",
			number: "12345678903",
			want:   true,
		},
		{
			name:   "valid luhn number 3",
			number: "4561261212345467",
			want:   true,
		},
		{
			name:   "invalid luhn number",
			number: "79927398710",
			want:   false,
		},
		{
			name:   "invalid luhn number 2",
			number: "12345678901",
			want:   false,
		},
		{
			name:   "empty string",
			number: "",
			want:   false,
		},
		{
			name:   "contains letters",
			number: "1234abc",
			want:   false,
		},
		{
			name:   "single digit valid",
			number: "0",
			want:   true,
		},
		{
			name:   "single digit invalid",
			number: "1",
			want:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsValid(tt.number)
			if got != tt.want {
				t.Errorf("IsValid(%q) = %v, want %v", tt.number, got, tt.want)
			}
		})
	}
}

