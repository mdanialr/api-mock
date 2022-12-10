package api

import (
	"errors"
	"fmt"
	"go/types"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/spf13/viper"
)

// NewMainHandler create new instance of MainHandler with separate config. Pointer viper
// param is important to make config live reload work.
func NewMainHandler(conf *viper.Viper) *MainHandler {
	return &MainHandler{
		conf: conf,
	}
}

type MainHandler struct {
	conf *viper.Viper
}

// Entry handle all incoming request including all http methods that supported by echo.
func (h *MainHandler) Entry(c echo.Context) error {
	// return not found if there is no matching endpoint in the config
	if err := h.checkEndpoint(c.Request().RequestURI); err != nil {
		msg := map[string]string{
			"message": "oops!! there is nothing here ¯\\_(ツ)_/¯",
		}
		return c.JSON(http.StatusNotFound, msg)
	}

	endpoint := fmt.Sprintf("endpoint.%s", c.Request().RequestURI)
	statusCode := h.getStatusCode(endpoint)
	sample, err := h.getSample(endpoint, statusCode)
	if err != nil {
		msg := map[string]string{
			"message": "there is no sample in the config, make sure you add it and follow the convention",
		}
		return c.JSON(http.StatusBadRequest, msg)
	}

	return c.JSONBlob(statusCode, []byte(sample))
}

// checkEndpoint check whether the given endpoint is exists or not in config.
// return error if not found.
func (h *MainHandler) checkEndpoint(endpoint string) error {
	path := h.conf.GetStringMap("endpoint")

	if path[endpoint] == nil {
		return errors.New("")
	}

	return nil
}

// getStatusCode get which status code will be used in the config. Default is 200.
func (h *MainHandler) getStatusCode(endpoint string) int {
	status := 200

	if code := h.conf.GetInt(fmt.Sprintf("%s.status", endpoint)); code > status {
		status = code
	}

	return status
}

// getSample retrieve sample data from the given endpoint and status code.
// return error if the sample was not found.
func (h *MainHandler) getSample(endpoint string, status int) (string, error) {
	var dataLen int
	var samples []string
	chosenData := 1 // default index to choose is the first sample

	conf := h.conf.GetStringMap(fmt.Sprintf("%s.%d", endpoint, status))

	if whichData, ok := conf["data"].(int); ok {
		chosenData = whichData
	}
	if allData, ok := conf["sample"].([]interface{}); ok { // can not directly cast to []string
		for _, val := range allData { // manually convert to []string
			if v, ok := val.(string); ok {
				samples = append(samples, v)
			}
		}

		dataLen = len(allData)
	}
	// in case config only has one sample
	if allData, ok := conf["sample"].(string); ok {
		samples = append(samples, allData)

		dataLen = 1 // because there is only one data
	}

	// if samples was not found, give error back
	if len(samples) < 1 {
		return "", types.Error{}
	}
	// if the given chosen data is exceeding the number of samples then use the last sample instead
	if chosenData > dataLen {
		return samples[dataLen-1], nil
	}

	return samples[chosenData-1], nil
}
