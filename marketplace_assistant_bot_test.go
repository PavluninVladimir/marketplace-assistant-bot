package main

import (
	"reflect"
	"testing"
)

func TestCreateInlineKeyboardButtonsBot(t *testing.T) {
	matrix := make([][]InlineKeyboardButton, 4)
	matrix[0] = make([]InlineKeyboardButton, 5)
	matrix[1] = make([]InlineKeyboardButton, 2)
	matrix[2] = make([]InlineKeyboardButton, 1)
	matrix[3] = make([]InlineKeyboardButton, 1)
	matrix[0][0] = InlineKeyboardButton{Text: "1:1", CallbackData: "/setting"}
	matrix[0][1] = InlineKeyboardButton{Text: "1:2", CallbackData: "/setting"}
	matrix[0][2] = InlineKeyboardButton{Text: "1:3", CallbackData: "/setting"}
	matrix[0][3] = InlineKeyboardButton{Text: "1:4", CallbackData: "/setting"}
	matrix[0][4] = InlineKeyboardButton{Text: "1:5", CallbackData: "/setting"}
	matrix[1][0] = InlineKeyboardButton{Text: "2:1", CallbackData: "/setting"}
	matrix[1][1] = InlineKeyboardButton{Text: "2:2", CallbackData: "/setting"}
	matrix[2][0] = InlineKeyboardButton{Text: "3:2", CallbackData: "/setting"}
	matrix[3][0] = InlineKeyboardButton{Text: "5:2", CallbackData: "/setting"}
	type args struct {
		b []ButtonBot[InlineKeyboardButton]
	}
	tests := []struct {
		name string
		args args
		want [][]InlineKeyboardButton
	}{
		{
			name: "Create struct button for bot telegram",
			args: args{
				b: []ButtonBot[InlineKeyboardButton]{
					{Row: 1, Col: 1, Button: InlineKeyboardButton{Text: "1:1", CallbackData: "/setting"}},
					{Row: 1, Col: 2, Button: InlineKeyboardButton{Text: "1:2", CallbackData: "/setting"}},
					{Row: 1, Col: 3, Button: InlineKeyboardButton{Text: "1:3", CallbackData: "/setting"}},
					{Row: 1, Col: 4, Button: InlineKeyboardButton{Text: "1:4", CallbackData: "/setting"}},
					{Row: 1, Col: 5, Button: InlineKeyboardButton{Text: "1:5", CallbackData: "/setting"}},
					{Row: 2, Col: 1, Button: InlineKeyboardButton{Text: "2:1", CallbackData: "/setting"}},
					{Row: 2, Col: 2, Button: InlineKeyboardButton{Text: "2:2", CallbackData: "/setting"}},
					{Row: 5, Col: 2, Button: InlineKeyboardButton{Text: "5:2", CallbackData: "/setting"}},
					{Row: 3, Col: 2, Button: InlineKeyboardButton{Text: "3:2", CallbackData: "/setting"}},
				},
			},
			want: matrix,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CreateButtonsBot[InlineKeyboardButton](tt.args.b); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("CreateInlineKeyboardButtonsBot() = %v, want %v", got, tt.want)
			}
		})
	}
}
