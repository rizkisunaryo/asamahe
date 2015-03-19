package asamahetil

import (
	//	"fmt"
	"sunaryo/util"
	"sunaryo/util/sunhttp"
	"time"
)

func SaveViewActivity(viewer string, activity string, long float64, lat float64) error {
	url := "http://localhost:9200/asamahe/viewact"
	if viewer == "" {
		viewer = "n a"
	}
	b := viewActReqJsonTpl(viewer, activity, long, lat)
	err := sunhttp.PostNoResponse(url, b)

	//	fmt.Printf(string(b))
	return err
}

func viewActReqJsonTpl(viewer string, activity string, long float64, lat float64) []byte {
	b := []byte(`{"Viewer":"`)
	b = append(b, viewer...)
	b = append(b, []byte(`","Activity":"`)...)
	b = append(b, activity...)
	b = append(b, []byte(`","Time":"`)...)
	tm := time.Now().Local()
	b = append(b, tm.Format("2006/01/02 15:04:05.000")...)
	if long != 0.0 || lat != 0.0 {
		b = append(b, []byte(`","Geo":{"Loc":[`)...)
		b = append(b, util.FloatToString(long)...)
		b = append(b, ","...)
		b = append(b, util.FloatToString(lat)...)
		b = append(b, "]}}"...)
	} else {
		b = append(b, []byte(`"}`)...)
	}

	return b
}
