package ovh

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

func notificationEmailSortedIds(config *Config) ([]int64, error) {
	// Create Order
	log.Printf("[DEBUG] Will read notification emails ids")
	res := []int64{}

	endpoint := "/me/notification/email/history"
	if err := config.OVHClient.Get(endpoint, &res); err != nil {
		return nil, fmt.Errorf("calling Get %s:\n\t %q", endpoint, err)
	}

	sort.Slice(res, func(i, j int) bool { return res[i] < res[j] })

	return res, nil
}

func getNewNotificationEmail(matches []string, oldIds []int64, config *Config) (*NotificationEmail, error) {
	curIds, err := notificationEmailSortedIds(config)
	if err != nil {
		return nil, err
	}

	lastOldId := oldIds[len(oldIds)-1]
	for _, id := range curIds {
		// matching only new ids (NOTE; a set subtract would be a better impl)
		if id > lastOldId {
			log.Printf("[DEBUG] Will read notification email %d", id)
			email := &NotificationEmail{}
			endpoint := fmt.Sprintf("/me/notification/email/history/%d", id)
			if err := config.OVHClient.Get(endpoint, email); err != nil {
				return nil, fmt.Errorf("calling Get %s:\n\t %q", endpoint, err)
			}

			match := true
			for _, m := range matches {
				log.Printf("[DEBUG] test match %v", m)
				if !strings.Contains(email.Body, m) {
					match = false
				}
			}
			if match {
				return email, nil
			}
		}
	}

	log.Printf("[DEBUG] no new notification email")
	return nil, nil
}
