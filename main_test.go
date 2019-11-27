package main

import (
	"github.com/dm0275/go-utils/files"
	"reflect"
	"testing"
)

var awsAccounts = files.ReadFile("./awsSampleAccounts")

var defaultAccount = `[default]
aws_access_key_id = 12345678912345678912
aws_secret_access_key = 1234567890abcdefghigklmnopqrstuvwxyz1234
region = us-west-2
`

var accountA = `[accountA]
aws_access_key_id = 12345678912345678912
aws_secret_access_key = 1234567890abcdefghigklmnopqrstuvwxyz1234
region = us-west-2
`

func Test_getAccountMatches(t *testing.T) {
	type args struct {
		account           string
		credentialsString string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{name: "Default_Account", args: args{account: "default", credentialsString: awsAccounts,}, want: defaultAccount},
		{name: "Account_A", args: args{account: "accountA", credentialsString: awsAccounts,}, want: accountA},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAccountMatches(tt.args.account, tt.args.credentialsString); got != tt.want {
				t.Errorf("getAccountMatches() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_getAccountNames(t *testing.T) {
	type args struct {
		credentialsString string
	}
	tests := []struct {
		name string
		args args
		want []string
	}{
		{name: "Default_Account", args: args{credentialsString: awsAccounts,}, want: []string{"default", "accountA", "accountB"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getAccountNames(tt.args.credentialsString); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getAccountNames() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseAcctFields(t *testing.T) {
	type args struct {
		account           string
		credentialsString string
	}
	tests := []struct {
		name string
		args args
		want map[string]awsAccountFields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseAcctFields(tt.args.account, tt.args.credentialsString); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseAcctFields() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseAwsCredentials(t *testing.T) {
	type args struct {
		awsCredsLocation string
		awsAccount       string
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
		})
	}
}

func Test_setDefaultAccount(t *testing.T) {
	type args struct {
		awsAccount  string
		awsAccounts []map[string]awsAccountFields
	}
	tests := []struct {
		name string
		args args
		want []map[string]awsAccountFields
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := setDefaultAccount(tt.args.awsAccount, tt.args.awsAccounts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setDefaultAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}