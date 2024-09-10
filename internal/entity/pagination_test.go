package entity

import "testing"

func TestPagination_ToSQL(t *testing.T) {
	type fields struct {
		Limit  int32
		Offset int32
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "limit + offset",
			fields: fields{
				Limit:  1,
				Offset: 2,
			},
			want: " LIMIT 1 OFFSET 2",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := Pagination{
				Limit:  tt.fields.Limit,
				Offset: tt.fields.Offset,
			}
			if got := p.ToSQL(); got != tt.want {
				t.Errorf("ToSQL() = %v, want %v", got, tt.want)
			}
		})
	}
}
