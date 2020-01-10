package access

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestConcurrent_Limit(t *testing.T) {
	tests := []struct {
		name         string
		Max          int64
		Current      int64
		want         int
		callExpected bool
	}{
		{
			name:         "unlimited-success",
			Max:          0, // Unlimited.
			want:         http.StatusOK,
			callExpected: true,
		},
		{
			name:         "limited-success",
			Max:          1,
			Current:      0, // Not limited
			want:         http.StatusOK,
			callExpected: true,
		},
		{
			name:         "limited-service-unavailable",
			Max:          1,
			Current:      1, // Limited.
			want:         http.StatusServiceUnavailable,
			callExpected: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Concurrent{
				Max:     tt.Max,
				Current: tt.Current,
			}
			wasCalled := false
			tester := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				wasCalled = true
			})
			h := c.Limit(tester)
			rw := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodGet, "/", nil)

			// Simulate HTTP request with the Concurrent.Limit handler.
			h.ServeHTTP(rw, req)

			result := rw.Result()
			if result.StatusCode != tt.want {
				t.Errorf("Concurrent.Limit() = %v, want %v", result, tt.want)
			}
			if wasCalled != tt.callExpected {
				t.Errorf("Concurrent.Limit() called unexpected got = %t, want %t", wasCalled, tt.callExpected)
			}
		})
	}
}
