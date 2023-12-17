package integration

import (
	"FxService/controllers"
	"FxService/model/dto/request"
	"FxService/model/dto/response"
	"FxService/service/dal"
	"bytes"
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/mongodb"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"
)

func TestMongoDBIntegration(t *testing.T) {
	ctx := context.Background()

	mongodbContainer, err := mongodb.RunContainer(ctx, testcontainers.WithImage("mongo:latest"))
	if err != nil {
		panic(err)
	}

	// Clean up the container
	defer func() {
		if err := mongodbContainer.Terminate(ctx); err != nil {
			panic(err)
		}
	}()

	endpoint, err := mongodbContainer.ConnectionString(ctx)
	if err != nil {
		panic(err)
	}

	envErr := os.Setenv("MONGODB_URI", endpoint)
	if envErr != nil {
		panic(err)
	}
	dal.Init(os.Getenv("MONGODB_URI"))
	router := gin.Default()
	controllers.AddRoutes(router)
	ts := httptest.NewServer(router)

	id := createForexData(ts, t)
	getForexDataById(ts, t, id)
	updateForexDataById(ts, t, id)
	createForexData(ts, t)
	getForexData(ts, t)
	getConversion(ts, t)
	deleteForexDataById(ts, t, id)
}

func createForexData(ts *httptest.Server, t *testing.T) string {
	var effectiveDate = time.Now()
	requestDataJSON, _ := json.Marshal(request.CreateForexDataRequest{
		TenantId:                     1,
		BankId:                       1,
		BaseCurrency:                 "USD",
		TargetCurrency:               "EUR",
		Tier:                         "1",
		DirectIndirectFlag:           "Y",
		Multiplier:                   1,
		BuyRate:                      1,
		SellRate:                     2,
		TolerancePercentage:          1,
		EffectiveDate:                &effectiveDate,
		ExpirationDate:               nil,
		ContractRequirementThreshold: "",
	})
	resp, err := http.Post(ts.URL+"/api/forexrates", "application/json", bytes.NewBuffer(requestDataJSON))
	if err != nil {
		t.Fatalf("Failed to send HTTP request: %v", err)
	}
	body, err := io.ReadAll(resp.Body)
	var data response.ResponseWithSimpleData[response.ForexDataResponse]
	if err := json.Unmarshal(body, &data); err != nil {

	}

	assert.Equal(t, 1, data.Data.TenantId)
	assert.Equal(t, 200, resp.StatusCode)

	return data.Data.Id
}

func getForexDataById(ts *httptest.Server, t *testing.T, id string) {
	println(ts.URL + "/api/forexrates/" + id)
	resp, err := http.Get(ts.URL + "/api/forexrates/" + id)
	if err != nil {
		t.Fatalf("Failed to send HTTP request: %v", err)
	}
	body, err := io.ReadAll(resp.Body)
	var data response.ResponseWithSimpleData[response.ForexDataResponse]
	if err := json.Unmarshal(body, &data); err != nil {
	}
	assert.Equal(t, id, data.Data.Id)
	assert.Equal(t, 200, resp.StatusCode)
}

func updateForexDataById(ts *httptest.Server, t *testing.T, id string) {
	var effectiveDate = time.Now()
	requestDataJSON, _ := json.Marshal(request.UpdateForexDataRequest{
		DirectIndirectFlag:           "N",
		Multiplier:                   1,
		BuyRate:                      1,
		SellRate:                     2,
		TolerancePercentage:          1,
		EffectiveDate:                &effectiveDate,
		ExpirationDate:               nil,
		ContractRequirementThreshold: "",
	})
	req, err := http.NewRequest("PUT", ts.URL+"/api/forexrates/"+id, bytes.NewBuffer(requestDataJSON))
	if err != nil {
		t.Fatalf("Failed to create HTTP put request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send PUT request:%v", err)
	}
	body, err := io.ReadAll(resp.Body)
	var data response.ResponseWithSimpleData[response.ForexDataResponse]
	if err := json.Unmarshal(body, &data); err != nil {

	}

	assert.Equal(t, "N", data.Data.DirectIndirectFlag)
	assert.Equal(t, 200, resp.StatusCode)
}

func getForexData(ts *httptest.Server, t *testing.T) {
	resp, err := http.Get(ts.URL + "/api/forexrates?tenantId=1&bankId=1&baseCurrency=USD&targetCurrency=EUR")
	if err != nil {
		t.Fatalf("Failed to send HTTP request: %v", err)
	}
	body, err := io.ReadAll(resp.Body)
	var data response.ResponseWithArrayData[response.ForexDataResponse]
	if err := json.Unmarshal(body, &data); err != nil {
	}
	assert.Equal(t, 2, len(*data.Data))
	assert.Equal(t, 200, resp.StatusCode)
}

func getConversion(ts *httptest.Server, t *testing.T) {
	resp, err := http.Get(ts.URL + "/api/convert?tenantId=1&bankId=1&baseCurrency=USD&targetCurrency=EUR&tier=1&amount=1000")
	if err != nil {
		t.Fatalf("Failed to send HTTP request: %v", err)
	}
	body, err := io.ReadAll(resp.Body)
	var data response.ResponseWithSimpleData[response.ConversionResponse]
	if err := json.Unmarshal(body, &data); err != nil {
	}
	assert.Equal(t, 1000.00, data.Data.ConvertedAmount)
	assert.Equal(t, 200, resp.StatusCode)
}

func deleteForexDataById(ts *httptest.Server, t *testing.T, id string) {
	req, err := http.NewRequest("DELETE", ts.URL+"/api/forexrates/"+id, nil)
	if err != nil {
		t.Fatalf("Failed to create HTTP Delete request: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("Failed to send Delete request:%v", err)
	}

	body, err := io.ReadAll(resp.Body)
	var data response.ResponseWithSimpleData[response.ForexDataResponse]
	if err := json.Unmarshal(body, &data); err != nil {

	}

	assert.Equal(t, 200, resp.StatusCode)
	assert.Equal(t, response.Success, data.Status)
}
