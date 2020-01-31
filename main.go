package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/dm0275/go-utils/errors"
	"github.com/dm0275/go-utils/files"
	"github.com/dm0275/go-utils/slices"
	"os"
	"regexp"
	"strings"
)

type awsAccountFields struct {
	region string
	accessKeyID string
	secretAccessKey string
}

/*
Return the profile names from the AWS credentials file
 */
func getProfileNames(awsCredentialsString string) []string {
	awsProfiles := []string{}
	if awsCredentialsString == "" {
		return awsProfiles
	}
	re := regexp.MustCompile(`\[.*\]`)
	matches := re.FindAllString(string(awsCredentialsString), -1)
	for _, acct := range matches {
		awsProfiles = append(awsProfiles, string(acct[1:len(acct)-1]))
	}
	return awsProfiles
}

/*
Parse profile name from an AWS Credentials file
 */
func parseAwsProfileName(profile string, awsCredentialsString string) string {
	awsProfile := ""
	if len(profile) > 0 {
		re := regexp.MustCompile(`(?m)^.*\[` + profile + `]*\]\n(([a-z].+(\n|$))*)`)
		match := re.FindAllString(string(awsCredentialsString), -1)
		if match[0] != "" {
			awsProfile = match[0]
		}
	}
	return awsProfile
}

/*
Parse profile data from an AWS Credentials file
 */
func parseAwsProfileData(profile string, awsCredentialsString string) map[string]awsAccountFields {
	awsProfiles := make(map[string]awsAccountFields)

	region := ""
	accessKeyID := ""
	secretAccessKey := ""

	profileScanner := bufio.NewScanner(strings.NewReader(parseAwsProfileName(profile, awsCredentialsString)))
	for profileScanner.Scan() {
		if !strings.Contains(profileScanner.Text(), "[") {
			field := strings.Split(strings.TrimSpace(profileScanner.Text()), "=")
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
	awsProfiles[profile] = awsAccountFields{region, accessKeyID, secretAccessKey}
	return awsProfiles
}

/*
Given a [] of awsProfiles, this func sets the default profile
 */
func setDefaultAccount(profile string, awsProfiles []map[string]awsAccountFields) []map[string]awsAccountFields {
	defaultAcct := awsAccountFields{}
	for _, acctMap := range awsProfiles {
		for key, acct := range acctMap {
			if key == profile {
				defaultAcct = awsAccountFields{region:acct.region, accessKeyID:acct.accessKeyID, secretAccessKey:acct.secretAccessKey}
			}
		}
	}

	for _, acctMap := range awsProfiles {
		for key := range acctMap {
			if key == "default" {
				acctMap["default"] = defaultAcct
			}
		}
	}
	return awsProfiles
}

/*
updateAwsCredentialsFile overwrites the AWS credentials file with the updated default profile
 */
func updateAwsCredentialsFile(awsCredentialsFile string, awsProfiles string) bool {
	updatedProfiles := files.OverwriteFile(awsCredentialsFile, awsProfiles, 0644)
	errors.CheckError(updatedProfiles)

	if updatedProfiles == nil {
		return true
	}
	return false
}

func parseAwsCredentials(awsCredsLocation string, awsAccount string) {
	awsCredsFile := files.ReadFile(awsCredsLocation)
	accountNames := getProfileNames(awsCredsFile)

	if !slices.StringInSlice(awsAccount, accountNames) {
		fmt.Println("Invalid AWS profile account: " + awsAccount + "\nThe only valid options are [" + strings.Join(accountNames, ", ") + "]")
		os.Exit(1)
	}
	
	awsAccounts := make([]map[string]awsAccountFields, 0)
	for _, acct := range accountNames {
		awsAccounts = append(awsAccounts, parseAwsProfileData(acct, awsCredsFile))
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

	if updateAwsCredentialsFile(awsCredsLocation, stringBuilder) {
		fmt.Println("Successfully set \"" + awsAccount + "\" as the default AWS profile.")
	} else {
		fmt.Println("There was an issue updating your AWS credentials file.")
		os.Exit(1)
	}
}

func main() {
	awsAccount := flag.String("awsAccount", "","Pass the awsAccount that you want to set as default. " +
		"You can also pass this parameter as a command line argument. Ex. awsProfileSwitcher <account>")
	awsCredentialsFile := flag.String("awsCredentialsFile",
		os.Getenv("HOME")+"/.aws/credentials", "The full path to your AWS credentials file.")
	cmdArgs := os.Args[1:]

	flag.Parse()
	if flag.NArg() > 0 {
		*awsAccount = cmdArgs[0]
	}

	if *awsAccount == "" {
		fmt.Println("No aws profile was passed, aws credentials was not modified")
		flag.Usage()
		os.Exit(1)
	} else {
		parseAwsCredentials(*awsCredentialsFile, *awsAccount)
	}
}