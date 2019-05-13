package api

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/pborman/uuid"
	"log"
)

var db = map[string]*ContentDeliveryNetworkV1{}

func ContentDeliveryNetworkCreateV1(w http.ResponseWriter, r *http.Request) {
	if AuthenticateRequest(r, w) != nil {
		return
	}
	err := validateExpectedHeaders(r.Header,"X-Request-Id", "User-Agent")
	if err != nil {
		sendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}
	cdn := &ContentDeliveryNetworkV1{}
	err = readRequest(r, cdn)
	if err != nil {
		sendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	if err := validateMandatoryFields(cdn); err != nil {
		sendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}

	populateComputePropertiesCDN(cdn)
	db[cdn.Id] = cdn
	log.Printf("POST [%+v\n]", cdn)
	sendResponse(http.StatusCreated, w, cdn)
}

func ContentDeliveryNetworkGetV1(w http.ResponseWriter, r *http.Request) {
	if AuthenticateRequest(r, w) != nil {
		return
	}
	err := validateExpectedHeaders(r.Header, "User-Agent")
	if err != nil {
		sendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}
	cdn, err := retrieveCdn(r)
	log.Printf("GET [%+v\n]", cdn)
	if err != nil {
		sendErrorResponse(http.StatusNotFound, err.Error(), w)
		return
	}
	sendResponse(http.StatusOK, w, cdn)
}

func ContentDeliveryNetworkUpdateV1(w http.ResponseWriter, r *http.Request) {
	if AuthenticateRequest(r, w) != nil {
		return
	}
	err := validateExpectedHeaders(r.Header, "User-Agent")
	if err != nil {
		sendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}
	cdn, err := retrieveCdn(r)
	if err != nil {
		sendErrorResponse(http.StatusNotFound, err.Error(), w)
		return
	}
	newCDN := &ContentDeliveryNetworkV1{}
	err = readRequest(r, newCDN)
	if err != nil {
		sendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}
	log.Printf("UPDATE [%+v\n]", newCDN)
	updateCDN(cdn, newCDN)
	sendResponse(http.StatusOK, w, newCDN)
}

func ContentDeliveryNetworkDeleteV1(w http.ResponseWriter, r *http.Request) {
	if AuthenticateRequest(r, w) != nil {
		return
	}
	err := validateExpectedHeaders(r.Header, "User-Agent")
	if err != nil {
		sendErrorResponse(http.StatusBadRequest, err.Error(), w)
		return
	}
	cdn, err := retrieveCdn(r)
	if err != nil {
		sendErrorResponse(http.StatusNotFound, err.Error(), w)
		return
	}
	delete(db, cdn.Id)
	log.Printf("DELETE [%s]", cdn.Id)
	updateResponseHeaders(http.StatusNoContent, w)
}

func populateComputePropertiesCDN(cdn *ContentDeliveryNetworkV1) {
	cdn.Id = uuid.New()
	cdn.ObjectNestedSchemeProperty[0].Name = "autogenerated name"
}

func validateMandatoryFields(cdn *ContentDeliveryNetworkV1) error {
	if cdn.Label == "" {
		return fmt.Errorf("mandatory 'label' field not populated")
	}
	if len(cdn.Ips) <= 0 {
		return fmt.Errorf("mandatory 'ips' list field not populated")
	}
	if len(cdn.Hostnames) <= 0 {
		return fmt.Errorf("mandatory 'hostnames' field not populated")
	}
	if cdn.Label == "" {
		return fmt.Errorf("mandatory label field not populated")
	}
	if len(cdn.ObjectNestedSchemeProperty) == 1 {
		if len(cdn.ObjectNestedSchemeProperty[0].ObjectProperty) != 1 {
			return fmt.Errorf("mandatory object_nested_scheme_property.object_property field not populated")
		}
	}
	return nil
}

func updateCDN(dbCDN, updatedCDN *ContentDeliveryNetworkV1) {
	dbCDN.Label = updatedCDN.Label
	dbCDN.Ips = updatedCDN.Ips
	dbCDN.Hostnames = updatedCDN.Hostnames
	dbCDN.ExampleInt = updatedCDN.ExampleInt
	dbCDN.ExampleNumber = updatedCDN.ExampleNumber
	dbCDN.ExampleBoolean = updatedCDN.ExampleBoolean
	dbCDN.ObjectProperty = updatedCDN.ObjectProperty
	dbCDN.ArrayOfObjectsExample = updatedCDN.ArrayOfObjectsExample
	db[dbCDN.Id] = dbCDN
}

func validateExpectedHeaders(headers map[string][]string, expectedHeaders ...string) error {
	for _, header := range expectedHeaders {
		err := validateHeader(headers, header)
		if err != nil {
			return err
		}
	}
	return nil
}

func validateHeader(headers map[string][]string, expectedHeader string) error {
	value, exists := headers[expectedHeader]
	log.Printf("validating header %s = %s", expectedHeader, value)
	if !exists {
		return fmt.Errorf("expected header [%s] is missing", expectedHeader)
	}
	return nil
}

func retrieveCdn(r *http.Request) (*ContentDeliveryNetworkV1, error) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/cdns/")
	if id == "" {
		return nil, fmt.Errorf("cdn id path param not provided")
	}
	cdn := db[id]
	if cdn == nil {
		return nil, fmt.Errorf("cdn id '%s' not found", id)
	}
	return cdn, nil
}
