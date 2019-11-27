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
aws_access_key_id = abcdefghigklmnopqrst
aws_secret_access_key = abcdefghigklmnopqrstuvwxyz12345678901234
region = us-west-2
`

var defaultAccounts = []map[string]awsAccountFields {
	{"default": {region: "us-west-2", accessKeyID: "12345678912345678912", secretAccessKey: "1234567890abcdefghigklmnopqrstuvwxyz1234"}},
	{"accountA": {region: "us-west-2", accessKeyID: "abcdefghigklmnopqrst", secretAccessKey: "abcdefghigklmnopqrstuvwxyz12345678901234"}},
	{"accountB": {region: "us-west-2", accessKeyID: "ABCD1234567890abcdef", secretAccessKey: "ABCD1234567890abcdefghigklmnopqrstuvwxyz"}},
}

var updatedAccountsA = []map[string]awsAccountFields {
	{"default": {region: "us-west-2", accessKeyID: "abcdefghigklmnopqrst", secretAccessKey: "abcdefghigklmnopqrstuvwxyz12345678901234"}},
	{"accountA": {region: "us-west-2", accessKeyID: "abcdefghigklmnopqrst", secretAccessKey: "abcdefghigklmnopqrstuvwxyz12345678901234"}},
	{"accountB": {region: "us-west-2", accessKeyID: "ABCD1234567890abcdef", secretAccessKey: "ABCD1234567890abcdefghigklmnopqrstuvwxyz"}},
}

var updatedAccountsB = []map[string]awsAccountFields {
	{"default": {region: "us-west-2", accessKeyID: "ABCD1234567890abcdef", secretAccessKey: "ABCD1234567890abcdefghigklmnopqrstuvwxyz"}},
	{"accountA": {region: "us-west-2", accessKeyID: "abcdefghigklmnopqrst", secretAccessKey: "abcdefghigklmnopqrstuvwxyz12345678901234"}},
	{"accountB": {region: "us-west-2", accessKeyID: "ABCD1234567890abcdef", secretAccessKey: "ABCD1234567890abcdefghigklmnopqrstuvwxyz"}},
}

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
		{name: "Get_All_Accounts", args: args{credentialsString: awsAccounts,}, want: []string{"default", "accountA", "accountB"}},
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
		{name: "Verify_Default_Account", args: args{account: "default", credentialsString: awsAccounts,}, want: defaultAccounts[0]},
		{name: "Verify_Account_A", args: args{account: "accountA", credentialsString: awsAccounts,}, want: defaultAccounts[1]},
		{name: "Verify_Account_B", args: args{account: "accountB", credentialsString: awsAccounts,}, want: defaultAccounts[2]},
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
		{name: "Verify_Account_B", args: args{awsAccount: "accountA", awsAccounts: defaultAccounts,}, want: updatedAccountsA},
		{name: "Verify_Account_B", args: args{awsAccount: "accountB", awsAccounts: defaultAccounts,}, want: updatedAccountsB},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := setDefaultAccount(tt.args.awsAccount, tt.args.awsAccounts); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("setDefaultAccount() = %v, want %v", got, tt.want)
			}
		})
	}
}