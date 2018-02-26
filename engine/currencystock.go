package engine

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

var stockNames []string

func init() {
	stockNames = []string{"EUR/USD", "USD/JPY", "GBP/USD", "USD/CHF", "USD/CAD", "AUD/USD"}
}

// GetStocks gets the current stock rates based on the chart above
func GetStocks() map[string]float64 {
	result := make(map[string]float64)
	for _, v := range stockNames {
		cur := strings.Split(v, "/")[1]
		base := strings.Split(v, "/")[0]
		resp, err := http.Get("https://api.fixer.io/latest?symbols=" + cur + "&base=" + base)
		if err != nil {
			fmt.Println("error getting stocks: " + err.Error())
			// we encountered an error lets not use these stock values
			return nil
		}
		jsonData, _ := ioutil.ReadAll(resp.Body)
		jsonObject := make(map[string]interface{})
		_ = json.Unmarshal(jsonData, &jsonObject)
		rates := jsonObject["rates"].(map[string]interface{})
		result[v] = rates[strings.ToUpper(cur)].(float64)
	}
	return result
}
