package ovh

import (
	"context"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
)

// dedicatedCloudUserTestPrefix marks users created by the dedicatedCloud user
// acceptance tests. It can't reuse test_prefix ("testacc-terraform") or the
// "tf-" idioms elsewhere: the dedicatedCloud user "name" is validated
// server-side as a shortname and rejects hyphens.
const dedicatedCloudUserTestPrefix = "tfacctest"

func init() {
	resource.AddTestSweepers("ovh_dedicated_cloud_user", &resource.Sweeper{
		Name: "ovh_dedicated_cloud_user",
		F:    testSweepDedicatedCloudUser,
	})
}

type dedicatedCloudUserSweepItem struct {
	UserId int64  `json:"userId"`
	Name   string `json:"name"`
}

// testSweepDedicatedCloudUser deletes leftover test users (and, as a
// consequence, any object rights granted to them) from OVH_DEDICATED_CLOUD_TEST.
// Best-effort: listing/deletion errors are logged and swallowed, never failed.
func testSweepDedicatedCloudUser(region string) error {
	client, err := sharedClientForRegion(region)
	if err != nil {
		return fmt.Errorf("error getting client: %s", err)
	}

	serviceName := os.Getenv("OVH_DEDICATED_CLOUD_TEST")
	if serviceName == "" {
		log.Print("[DEBUG] OVH_DEDICATED_CLOUD_TEST is not set. No dedicatedCloud user to sweep")
		return nil
	}

	// Deleting a user also removes every object right granted to it, so
	// sweeping users is sufficient - no separate object right sweep needed.
	usersEndpoint := fmt.Sprintf("/dedicatedCloud/%s/user", serviceName)

	var userIds []int64
	if err := client.Get(usersEndpoint, &userIds); err != nil {
		log.Printf("[DEBUG] error listing users to sweep (GET %s): %s", usersEndpoint, err)
		return nil
	}

	for _, userId := range userIds {
		userEndpoint := fmt.Sprintf("%s/%d", usersEndpoint, userId)

		var user dedicatedCloudUserSweepItem
		if err := client.Get(userEndpoint, &user); err != nil {
			log.Printf("[DEBUG] error reading user %d to sweep: %s", userId, err)
			continue
		}

		if !strings.HasPrefix(user.Name, dedicatedCloudUserTestPrefix) {
			continue
		}

		log.Printf("[DEBUG] sweeping dedicatedCloud user %d (%q) from %s", userId, user.Name, serviceName)

		var task DedicatedCloudTask
		if err := client.Delete(userEndpoint, &task); err != nil {
			log.Printf("[DEBUG] error deleting user %d: %s", userId, err)
			continue
		}

		if _, err := waitForDedicatedCloudTask(context.Background(), client, serviceName, task.TaskId); err != nil {
			log.Printf("[DEBUG] error waiting for user %d deletion: %s", userId, err)
		}
	}

	return nil
}
