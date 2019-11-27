package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/dm0275/Utils"
	"os"
	"regexp"
	"strings"
)

type awsAccountFields struct {
	region string
	accessKeyID string
	secretAccessKey string
}

func getAccountNames(credentialsString string) []string {
	awsAccounts := []string{}
	if credentialsString == "" {
		return awsAccounts
	}
	re := regexp.MustCompile(`\[.*\]`)
	matches := re.FindAllString(string(credentialsString), -1)
	for _, acct := range matches {
		awsAccounts = append(awsAccounts, string(acct[1:len(acct)-1]))
	}
	return awsAccounts
}

func getAccountMatches(account string, credentialsString string) string {
	accountMatches := ""
	if len(account) > 0 {
		re := regexp.MustCompile(`(?m)^.*\[` + account + `]*\]\n(([a-z].+(\n|$))*)`)
		match := re.FindAllString(string(credentialsString), -1)
		if match[0] != "" {
			accountMatches = match[0]
		}
	}
	return accountMatches
}

func parseAcctFields(account string, credentialsString string) map[string]awsAccountFields {
	accounts := make(map[string]awsAccountFields)

	region := ""
	accessKeyID := ""
	secretAccessKey := ""

	accountScanner := bufio.NewScanner(strings.NewReader(getAccountMatches(account, credentialsString)))
	for accountScanner.Scan() {
		if !strings.Contains(accountScanner.Text(), "[") {
			field := strings.Split(strings.TrimSpace(accountScanner.Text()), "=")
			if len(field) > 1 {
				switch acctField := strings.TrimSpace(field[0]); acctField {
				case "aws_access_key_id":
					accessKeyID = strings.TrimSpace(field[1])
				case "aws_secret_access_key":
					secretAccessKey = strings.TrimSpace(field[1])
				case "region":
					region = strings.TrimSpace(field[1])
				}
			}
		}
	}
	accounts[account] = awsAccountFields{region, accessKeyID, secretAccessKey}
	return accounts
}

func setDefaultAccount(awsAccount string, awsAccounts []map[string]awsAccountFields) []map[string]awsAccountFields {
	defaultAcct := awsAccountFields{}
	for _, acctMap := range awsAccounts {
		for key, acct := range acctMap {
			if key == awsAccount {
				defaultAcct = awsAccountFields{region:acct.region, accessKeyID:acct.accessKeyID, secretAccessKey:acct.secretAccessKey}
			}
		}
	}

	for _, acctMap := range awsAccounts {
		for key := range acctMap {
			if key == "default" {
				acctMap["default"] = defaultAcct
			}
		}
	}
	return awsAccounts
}

func parseAwsCredentials(awsCredsLocation string, awsAccount string) {
	awsCredsFile := utils.ReadFile(awsCredsLocation)
	accountNames := getAccountNames(awsCredsFile)

	if utils.StringInSlice(awsAccount, accountNames) {

	} else {
		fmt.Println("Invalid AWS profile account: " + awsAccount + "\nThe only valid options are [" + strings.Join(accountNames, ", ") + "]")
		os.Exit(1)
	}
	
	awsAccounts := make([]map[string]awsAccountFields, 0)
	for _, acct := range accountNames {
		awsAccounts = append(awsAccounts, parseAcctFields(acct, awsCredsFile))
	}

	updatedAccounts := setDefaultAccount(awsAccount, awsAccounts)

	stringBuilder := ""
	for index, acctMap := range updatedAccounts {
		for key, acct := range acctMap {
			stringBuilder = stringBuilder + "[" + key + "]\n" +
				"aws_access_key_id = " + acct.accessKeyID + "\n" +
				"aws_secret_access_key = " + acct.secretAccessKey + "\n"
				if acct.region != "" {
					stringBuilder = stringBuilder + "region = " + acct.region
				}
			if index == len(updatedAccounts)-1 {
				stringBuilder += "\n"
			} else {
				stringBuilder += "\n\n"
			}
		}
	}
	fmt.Println(stringBuilder)
}

func main() {
	//argsWithProg := os.Args
	//argsWithoutProg := os.Args[1:]
	//arg := os.Args[1]

	//fmt.Println(argsWithProg)
	//fmt.Println(argsWithoutProg)
	//fmt.Println(arg)

	awsAccount := flag.String("awsAccount", "","Pass the awsAccount that you want to set as default")
	awsCredentialsFile := flag.String("awsCredentialsFile",
		os.Getenv("HOME")+"/.aws/credentials2", "The full path to your AWS credentials file.")
	flag.Parse()

	if *awsAccount == "" {
		fmt.Println("No aws profile was passed, aws credentials was not modified")
	} else {
		parseAwsCredentials(*awsCredentialsFile, *awsAccount)
	}
}