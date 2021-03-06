package oci

import (
	"context"
	"fmt"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/oracle/oci-go-sdk/common"
	oci_database "github.com/oracle/oci-go-sdk/database"
)

func sweepDatabaseDbSystemResource(compartment string) error {
	databaseClient := GetTestClients(&schema.ResourceData{}).databaseClient
	dbSystemIds, err := getDbSystemIds(compartment)
	if err != nil {
		return err
	}
	for _, dbSystemId := range dbSystemIds {
		if ok := SweeperDefaultResourceId[dbSystemId]; !ok {
			terminateDbSystemRequest := oci_database.TerminateDbSystemRequest{}

			terminateDbSystemRequest.DbSystemId = &dbSystemId

			terminateDbSystemRequest.RequestMetadata.RetryPolicy = getRetryPolicy(true, "database")
			_, error := databaseClient.TerminateDbSystem(context.Background(), terminateDbSystemRequest)
			if error != nil {
				fmt.Printf("Error deleting DbSystem %s %s, It is possible that the resource is already deleted. Please verify manually \n", dbSystemId, error)
				continue
			}
			waitTillCondition(testAccProvider, &dbSystemId, dbSystemSweepWaitCondition, time.Duration(3*time.Minute),
				dbSystemSweepResponseFetchOperation, "database", true)
		}
	}
	return nil
}

func getDbSystemIds(compartment string) ([]string, error) {
	ids := getResourceIdsToSweep(compartment, "DbSystemId")
	if ids != nil {
		return ids, nil
	}
	var resourceIds []string
	compartmentId := compartment
	databaseClient := GetTestClients(&schema.ResourceData{}).databaseClient

	listDbSystemsRequest := oci_database.ListDbSystemsRequest{}
	listDbSystemsRequest.CompartmentId = &compartmentId
	listDbSystemsRequest.LifecycleState = oci_database.DbSystemSummaryLifecycleStateAvailable

	// Terminate the newest dbSystem first to make sure any standby databases created by Data Guard Assocuations are deleted first
	listDbSystemsRequest.SortBy = oci_database.ListDbSystemsSortByTimecreated
	listDbSystemsRequest.SortOrder = oci_database.ListDbSystemsSortOrderDesc

	listDbSystemsResponse, err := databaseClient.ListDbSystems(context.Background(), listDbSystemsRequest)

	if err != nil {
		return resourceIds, fmt.Errorf("Error getting DbSystem list for compartment id : %s , %s \n", compartmentId, err)
	}
	for _, dbSystem := range listDbSystemsResponse.Items {
		id := *dbSystem.Id
		resourceIds = append(resourceIds, id)
		addResourceIdToSweeperResourceIdMap(compartmentId, "DbSystemId", id)
	}
	return resourceIds, nil
}

func dbSystemSweepWaitCondition(response common.OCIOperationResponse) bool {
	// Only stop if the resource is available beyond 3 mins. As there could be an issue for the sweeper to delete the resource and manual intervention required.
	if dbSystemResponse, ok := response.Response.(oci_database.GetDbSystemResponse); ok {
		return dbSystemResponse.LifecycleState != oci_database.DbSystemLifecycleStateTerminated
	}
	return false
}

func dbSystemSweepResponseFetchOperation(client *OracleClients, resourceId *string, retryPolicy *common.RetryPolicy) error {
	_, err := client.databaseClient.GetDbSystem(context.Background(), oci_database.GetDbSystemRequest{
		DbSystemId: resourceId,
		RequestMetadata: common.RequestMetadata{
			RetryPolicy: retryPolicy,
		},
	})
	return err
}
