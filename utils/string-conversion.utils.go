package utils

import "strings"

func SplitNameStr(nameStr string) (string, string) {
	// if claims.name is not empty, then split the name into first name and last name, if it is 1, dont fill, if the lenght is 2, then split it, if it is 3 then choose the first two as first name and the last one as last name, and if it is 4 then choose the first two as first name and the last two as last name. otherwise, dont fill
	firstName := ""
	lastName := ""
	if nameStr != "" {
		// split the string into parts based on space
		nameParts := strings.Split(nameStr, " ")
		if len(nameParts) == 1 {
			// do nothing
		} else if len(nameParts) == 2 {
			firstName = nameParts[0]
			lastName = nameParts[1]
		} else if len(nameParts) == 3 {
			firstName = nameParts[0] + " " + nameParts[1]
			lastName = nameParts[2]
		} else if len(nameParts) == 4 {
			firstName = nameParts[0] + " " + nameParts[1]
			lastName = nameParts[2] + " " + nameParts[3]
		}
	}

	return firstName, lastName
}
