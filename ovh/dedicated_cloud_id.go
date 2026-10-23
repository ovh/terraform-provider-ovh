package ovh

import (
	"fmt"
	"strconv"
	"strings"
)

func splitDedicatedCloudUserId(id string) (serviceName string, userId int64, err error) {
	splits := strings.Split(id, "/")
	if len(splits) != 2 {
		return "", 0, fmt.Errorf("ID must be formatted like the following: <service_name>/<user_id>")
	}

	userId, convErr := strconv.ParseInt(splits[1], 10, 64)
	if convErr != nil {
		return "", 0, fmt.Errorf("ID must be formatted like the following: <service_name>/<user_id> where user_id is a number")
	}

	return splits[0], userId, nil
}

func splitDedicatedCloudUserObjectRightId(id string) (serviceName string, userId int64, objectRightId int64, err error) {
	splits := strings.Split(id, "/")
	if len(splits) != 3 {
		return "", 0, 0, fmt.Errorf("ID must be formatted like the following: <service_name>/<user_id>/<object_right_id>")
	}

	userId, convErr := strconv.ParseInt(splits[1], 10, 64)
	if convErr != nil {
		return "", 0, 0, fmt.Errorf("ID must be formatted like the following: <service_name>/<user_id>/<object_right_id> where user_id is a number")
	}

	objectRightId, convErr = strconv.ParseInt(splits[2], 10, 64)
	if convErr != nil {
		return "", 0, 0, fmt.Errorf("ID must be formatted like the following: <service_name>/<user_id>/<object_right_id> where object_right_id is a number")
	}

	return splits[0], userId, objectRightId, nil
}
