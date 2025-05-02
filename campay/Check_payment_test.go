package campay

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCheckPaymentStatus(t *testing.T) {
	mockResponse := `{
		"reference": "12345",
		"external_reference": "ext-12345",
		"status": "success",
		"amount": "1000",
		"currency": "USD",
		"operator": "OperatorX",
		"code": "200",
		"operator_reference": "op-12345",
		"description": "Payment successful",
		"external_user": "user-123",
		"reason": "",
		"phone_number": "1234567890"
	}`

	server := MakeTestServer(http.StatusOK, []byte(mockResponse))

	mockClient := &Requests{
		baseUrl: server.URL + "/",
		apikey:  "mockapikey",
		client:  http.DefaultClient,
	}

	expectedStatus := Status{
		Reference:          "12345",
		Ext_ref:            "ext-12345",
		Status:             "success",
		Amount:             "1000",
		Currency:           "USD",
		Operator:           "OperatorX",
		Code:               "200",
		Operator_Reference: "op-12345",
		Description:        "Payment successful",
		Exterbal_User:      "user-123",
		Reason:             "",
		Phone_Number:       "1234567890",
	}

	result := mockClient.CheckPaymentStatus("12345")

	if result.Ext_ref != expectedStatus.Ext_ref {
		t.Errorf("Expected %+v, got %+v", expectedStatus.Ext_ref, result.Ext_ref)
	}
}

// MakeTestServer creates an api server for testing
func MakeTestServer(responseCode int, body []byte) *httptest.Server {
	return httptest.NewServer(http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		res.WriteHeader(responseCode)
		_, err := res.Write(body)
		if err != nil {
			panic(err)
		}
	}))
}
