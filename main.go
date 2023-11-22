package main

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

func removeWhiteSpacesExceptQuotes(input string) string {
	// This pattern captures content within double quotes and matches all other whitespaces
	customRegex := regexp.MustCompile(`"([^"]*)"|\s+`)

	// Replace based on the custom regular expression
	result := customRegex.ReplaceAllStringFunc(input, func(match string) string {
		// If the match is within double quotes, preserve it
		if match[0] == '"' {
			return match
		}
		// Otherwise, remove whitespaces
		return ""
	})

	return result
}

func checkValueToken(value string) bool {
		if value[0] == '"' && value[len(value) - 1] == '"' {
			return true;
		} else if value == "null" || value == "true" || value == "false" {
			return true;	
		} 


		_, errInt := strconv.ParseInt(value, 10, 64);
		_, errFloat := strconv.ParseFloat(value, 64);

		return errInt == nil || errFloat == nil;

}

func checkValue(value string) bool {
	if value[0] == '{' {
		return checkValidJSON(value);
	} else if value[0] == '[' {
		value = value[1 : len(value) - 1];
		valueArray := splitKeyValuePairs(value);
		for _, val := range valueArray {
			if !checkValueToken(val) {
				return false;
			}
		}
		return true;		
	} else {
		return checkValueToken(value);
	}

}

func checkKey(key string) bool {
	if key[0] != '"' || key[len(key) - 1] != '"' {
		return false;
	}
	return true;

}

func checkKeyValue(input string) bool {

	keyValueArray := strings.SplitN(input, ":", 2)
	key := keyValueArray[0];
	value := keyValueArray[1];

	if(checkKey(key) && checkValue(value)) {
		return true;
	}
	return false;
		

}

func splitKeyValuePairs(input string) []string {
	// Define a regular expression pattern to match commas outside of curly braces and square brackets
	regexPattern := regexp.MustCompile(`(?:[^{},\[\]]|\[[^[\]]*]|\{[^{}]*\})*`)

	// Find all matches in the input string
	matches := regexPattern.FindAllString(input, -1)

	// If there are no matches (no commas), return the entire string
	if len(matches) == 0 {
		return []string{input}
	}

	var result []string

	// Process matches and add to the result
	for _, match := range matches {
		// Remove trailing commas
		trimmedMatch := strings.TrimRight(match, ",")
		// If the match is not empty, add it to the result
		if trimmedMatch != "" {
			result = append(result, trimmedMatch)
		}
	}

	return result
}

func checkValidJSON(fileString string) bool {
	if(fileString[0] != '{' || fileString[len(fileString) - 1] != '}') {
		return false;
	} else {
		fileString = fileString[1 : len(fileString) - 1];

		if len(fileString) == 0 {
			return true;
		}

		keyValuePairArray := splitKeyValuePairs(fileString);

		for _, pair := range keyValuePairArray {
			if !checkKeyValue(pair) {
				return false;
			}
		}
		return true;
	}

}


func main() {
	fileName := "tests/step4/valid.json";

	fileContents, err := os.ReadFile(fileName);

	if err != nil {
		fmt.Println(err);
	}

	str := string(fileContents)

	fileString := removeWhiteSpacesExceptQuotes(str);
	
	if(len(fileString) == 0) {
		fmt.Println("Invalid");
		return;
	}

	if checkValidJSON(fileString) {
		fmt.Println("Valid");
	} else {
		fmt.Println("Invalid");
	}

}