package wallet

import (
	"github.com/ToirovSadi/wallet/pkg/types"
	"reflect"
	"testing"
)

func TestService_Export(t *testing.T) {
	type fields struct {
		nextAccountID int64
		accounts      []*types.Account
		payments      []*types.Payment
		favorites     []*types.Favorite
	}
	type args struct {
		dir string
	}
	var tests = []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				nextAccountID: 4,
				accounts:      accounts,
				payments:      payments,
				favorites:     favorites,
			},
			args: args{
				dir: ".",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				nextAccountID: tt.fields.nextAccountID,
				accounts:      tt.fields.accounts,
				payments:      tt.fields.payments,
				favorites:     tt.fields.favorites,
			}
			if err := s.Export(tt.args.dir); (err != nil) != tt.wantErr {
				t.Errorf("Export() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_ExportToFile(t *testing.T) {
	type fields struct {
		nextAccountID int64
		accounts      []*types.Account
		payments      []*types.Payment
		favorites     []*types.Favorite
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				nextAccountID: 4,
				accounts:      accounts,
				payments:      payments,
				favorites:     favorites,
			},
			args: args{
				path: "./accounts.dump",
			},
			wantErr: false,
		},
		{
			name: "test2",
			fields: fields{
				nextAccountID: 4,
				accounts:      accounts,
				payments:      payments,
				favorites:     favorites,
			},
			args: args{
				path: "./nowhere/nothing/accounts.dump",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				nextAccountID: tt.fields.nextAccountID,
				accounts:      tt.fields.accounts,
				payments:      tt.fields.payments,
				favorites:     tt.fields.favorites,
			}
			if err := s.ExportToFile(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("ExportToFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestService_ExportAccountHistory(t *testing.T) {
	type fields struct {
		nextAccountID int64
		accounts      []*types.Account
		payments      []*types.Payment
		favorites     []*types.Favorite
	}
	type args struct {
		accountID int64
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []*types.Payment
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				nextAccountID: 4,
				accounts:      accounts,
				payments:      payments,
				favorites:     favorites,
			},
			args: args{
				accountID: 1,
			},
			want: []*types.Payment{
				payments[0],
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				nextAccountID: tt.fields.nextAccountID,
				accounts:      tt.fields.accounts,
				payments:      tt.fields.payments,
				favorites:     tt.fields.favorites,
			}
			got, err := s.ExportAccountHistory(tt.args.accountID)
			if (err != nil) != tt.wantErr {
				t.Errorf("ExportAccountHistory() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ExportAccountHistory() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestService_HistoryToFiles(t *testing.T) {
	type fields struct {
		nextAccountID int64
		accounts      []*types.Account
		payments      []*types.Payment
		favorites     []*types.Favorite
	}
	type args struct {
		payments []*types.Payment
		dir      string
		records  int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "test1",
			fields: fields{
				nextAccountID: 4,
				accounts:      accounts,
				payments:      payments,
				favorites:     favorites,
			},
			args: args{
				payments: payments,
				dir:      ".",
				records:  2,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Service{
				nextAccountID: tt.fields.nextAccountID,
				accounts:      tt.fields.accounts,
				payments:      tt.fields.payments,
				favorites:     tt.fields.favorites,
			}
			if err := s.HistoryToFiles(tt.args.payments, tt.args.dir, tt.args.records); (err != nil) != tt.wantErr {
				t.Errorf("HistoryToFiles() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
