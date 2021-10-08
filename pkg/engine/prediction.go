package engine

import (
	"encoding/json"
	"strconv"

	"github.com/JustHumanz/Go-Simp/pkg/config"
	"github.com/JustHumanz/Go-Simp/pkg/database"
	"github.com/JustHumanz/Go-Simp/pkg/network"
)

func Prediction(vtuber database.Member, state string, tar int) (int, error) {

	type PromeMetric struct {
		Status string `json:"status"`
		Data   struct {
			Resulttype string `json:"resultType"`
			Result     []struct {
				Metric interface{}   `json:"metric"`
				Value  []interface{} `json:"value"`
			} `json:"result"`
		} `json:"data"`
	}

	var Data PromeMetric

	RawQuery := "/api/v1/query?query=%0Apredict_linear(get_subscriber{state%3D\"" + state + "\"%2C+vtuber%3D\"" + vtuber.Name + "\"}+[356d]%2C86400*" + strconv.Itoa(tar) + ")"

	RawData, err := network.Curl(config.GoSimpConf.PrometheusURL+RawQuery, nil)
	if err != nil {
		return 0, err
	}
	err = json.Unmarshal(RawData, &Data)
	if err != nil {
		return 0, err
	}
	val := Data.Data.Result[0].Value[1].(string)
	valInt, err := strconv.ParseFloat(val, 64)
	if err != nil {
		return 0, err
	}
	return int(valInt), nil

}
