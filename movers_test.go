package td

import (
	"testing"

	"github.com/zeebo/assert"
)

func TestMoversReq_Encode(t *testing.T) {
	tests := []struct {
		name      string
		req       *MoversReq
		wantQuery string
		wantErr   bool
	}{
		{
			name: "valid request with sort",
			req: &MoversReq{
				SymbolID:  SymbolIdDJI,
				Sort:      SortVolume,
				Frequency: Frequency5,
			},
			wantQuery: "frequency=5&sort=VOLUME",
			wantErr:   false,
		},
		{
			name: "valid request without sort",
			req: &MoversReq{
				SymbolID:  SymbolIdSPY,
				Frequency: Frequency10,
			},
			wantQuery: "frequency=10",
			wantErr:   false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.req == nil {
				// test nil pointer safely
				var req *MoversReq
				_, err := req.Encode()
				assert.Error(t, err)
				return
			}

			got, err := tt.req.Encode()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}

			assert.NoError(t, err)
			assert.Equal(t, tt.wantQuery, got)
		})
	}
}
