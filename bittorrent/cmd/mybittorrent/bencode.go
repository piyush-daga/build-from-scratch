package main

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

// Returns match, rest of the string and any errors
// 5:hello5:world -> hello, 5:world, nil
func decodeBencodedStrings(bencodedString string) (interface{}, string, error) {
	before, remaining, found := strings.Cut(bencodedString, ":")
	if !found {
		return "", bencodedString, fmt.Errorf("colon not found in given bencoded string")
	}
	digits, err := strconv.Atoi(before)

	if err != nil {
		return "", "", err
	}

	return remaining[:digits], remaining[digits:], nil
}

// i52e5:hello -> 52, 5:hello
func decodeBencodedInts(bencodedString string) (interface{}, string, error) {
	index := strings.Index(bencodedString, "e")

	// Not Found
	if index == -1 {
		return -1, bencodedString, fmt.Errorf("given bencodedstring does not contain an integer")
	}

	result, err := strconv.ParseInt(bencodedString[1:index], 10, 64)
	if err != nil {
		return -1, bencodedString, err
	}
	return result, bencodedString[index+1:], nil
}

func isString(bencodedString string) bool {
	return unicode.IsDigit(rune(bencodedString[0]))
}

func isInteger(bencodedString string) bool {
	return bencodedString[0] == 'i'
}

func isList(bencodedString string) bool {
	return bencodedString[0] == 'l'
}

func isDict(bencodedString string) bool {
	return bencodedString[0] == 'd'
}

// le -> []
// llee -> [[]]
// lli52eee -> [[52]]
// li52ee5:hello -> [52], 5:hello
func decodeBencodedList(ben string) (interface{}, string, error) {
	var match []interface{}

	// Exit condition - This condition will only arise when a list is ending
	if ben[0] == 'e' {
		return nil, ben[1:], nil
	} else if isList(ben) {
		res, remaining, err := decodeBencodedList(ben[1:])
		if err != nil {
			return nil, remaining, err
		}
		// Initialize an array
		if match == nil {
			match = make([]interface{}, 0)
		} else {
			match = append(match, res)
		}

		return match, remaining, nil
	} else if isInteger(ben) {
		return decodeBencodedInts(ben)
	} else if isString(ben) {
		return decodeBencodedStrings(ben)
	}

	return nil, "", fmt.Errorf("oof, should not have reached this state in decoding lists")
}

// d5:hello3:heyei52e -> {"hello": hey}, i52e, nil
// d5:hellod2:ok5:helloee -> {"hello": {"ok": "hello"}}
func decodeBencodedDict(ben string) (interface{}, string, error) {
	var match map[string]interface{}
	// Any e encountered, means a dict has ended
	if ben[0] == 'e' {
		return nil, ben[1:], nil
	}

	if isDict(ben) {
		var (
			k         interface{}
			v         interface{}
			remaining string
		)

		if len(ben) == 2 {
			return make(map[string]interface{}, 0), "", nil
		}

		remaining = ben[1:]

		for remaining[0] != 'e' {
			// Find the key and the value and then proceed ahead
			k, remaining, _ = decodeBencodedDict(remaining)
			// Check if remaining exists or not
			if remaining != "" {
				v, remaining, _ = decodeBencodedDict(remaining)
			} else {
				v, remaining = nil, ""
			}

			if match == nil {
				match = make(map[string]interface{}, 0)
			}

			if v != nil {
				match[k.(string)] = v
			}

			// Update ben's value
			// ben = remaining
		}

		// For the remaining data
		return match, remaining[1:], nil
	} else if isInteger(ben) {
		return decodeBencodedInts(ben)
	} else if isString(ben) {
		return decodeBencodedStrings(ben)
	} else if isList(ben) {
		return decodeBencodedList(ben)
	}

	return nil, "", fmt.Errorf("oof, should not have reached this state in decoding dicts")
}

// This needs to be overhauled to support recursion
func DecodeBencode(bencodedString string) (interface{}, error) {
	var (
		rest  string
		match interface{}
		err   error
	)

	if isInteger(bencodedString) {
		match, rest, err = decodeBencodedInts(bencodedString)
	} else if isString(bencodedString) {
		match, rest, err = decodeBencodedStrings(bencodedString)
	} else if isList(bencodedString) {
		match, rest, err = decodeBencodedList(bencodedString)
	} else if isDict(bencodedString) {
		match, rest, err = decodeBencodedDict(bencodedString)
	}

	if rest == "" {
		return match, err
	}
	return DecodeBencode(rest)
}
