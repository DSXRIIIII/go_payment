/*
order service

Testing DefaultAPIService

*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech);

package openapi

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
	openapiclient "github.com/GIT_USER_ID/GIT_REPO_ID"
)

func Test_openapi_DefaultAPIService(t *testing.T) {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)

	t.Run("Test DefaultAPIService CustomerCustomerIDOrdersOrderIDGet", func(t *testing.T) {

		t.Skip("skip test")  // remove to run test

		var customerID string
		var orderID string

		resp, httpRes, err := apiClient.DefaultAPI.CustomerCustomerIDOrdersOrderIDGet(context.Background(), customerID, orderID).Execute()

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

	t.Run("Test DefaultAPIService CustomerCustomerIDOrdersPost", func(t *testing.T) {

		t.Skip("skip test")  // remove to run test

		var customerID string

		resp, httpRes, err := apiClient.DefaultAPI.CustomerCustomerIDOrdersPost(context.Background(), customerID).Execute()

		require.Nil(t, err)
		require.NotNil(t, resp)
		assert.Equal(t, 200, httpRes.StatusCode)

	})

}