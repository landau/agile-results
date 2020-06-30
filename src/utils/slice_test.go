package utils

import "testing"

func TestIndexOf(t *testing.T) {
	items := []int{0, 1, 2}

	// predicate := func(i) bool {
	// 	return items
	// }

	type args struct {
		predicate func(i int) bool
	}

	tests := []struct {
		name string
		args args
		want int
	}{
		{
			"Returns 0 for an item found at index 0",
			args{predicate: func(i int) bool { return i == items[0] }},
			0,
		},
		{
			"Returns 1 for an item found at index 1",
			args{predicate: func(i int) bool { return i == items[1] }},
			1,
		},
		{
			"Returns 2 for an item found at index 2",
			args{predicate: func(i int) bool { return i == items[2] }},
			2,
		},
		{
			"Returns -1 if an item is not found",
			args{predicate: func(i int) bool { return i == 500 }},
			-1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IndexOf(len(items), tt.args.predicate); got != tt.want {
				t.Errorf("IndexOf() = %v, want %v", got, tt.want)
			}
		})
	}
}
