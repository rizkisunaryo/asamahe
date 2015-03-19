package sunhttp

import (
	"bytes"
	"encoding/json"
	//	"fmt"
	"io/ioutil"
	"net/http"
)

func Post(url string, b []byte, v interface{}) error {
	req, err3 := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err3 != nil {
		return err3
	}
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return err2
	}

	defer resp.Body.Close()
	return json.Unmarshal([]byte(body), &v)
}

func PostStruct(url string, i interface{}, o interface{}) error {
	iStr, err1 := json.Marshal(i)
	if err1 != nil {
		return err1
	}
	req, err2 := http.NewRequest("POST", url, bytes.NewBuffer(iStr))
	if err2 != nil {
		return err2
	}
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err3 := client.Do(req)
	if err3 != nil {
		return err3
	}
	body, err4 := ioutil.ReadAll(resp.Body)
	if err4 != nil {
		return err4
	}

	defer resp.Body.Close()
	return json.Unmarshal([]byte(body), &o)
}

func PostNoResponse(url string, b []byte) error {
	req, err1 := http.NewRequest("POST", url, bytes.NewBuffer(b))
	if err1 != nil {
		return err1
	}
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err2 := client.Do(req)
	//	body, _ := ioutil.ReadAll(resp.Body)
	//	fmt.Printf(string(body))
	defer resp.Body.Close()

	return err2
}

func GetResponse(url string, v interface{}) error {
	resp, err1 := http.Get(url)
	if err1 != nil {
		return err1
	}
	body, err2 := ioutil.ReadAll(resp.Body)
	if err2 != nil {
		return err2
	}

	defer resp.Body.Close()
	return json.Unmarshal([]byte(body), &v)
}

func Delete(url string) error {
	req, err1 := http.NewRequest("DELETE", url, nil)
	if err1 != nil {
		return err1
	}
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err2 := client.Do(req)
	defer resp.Body.Close()

	return err2
}

func DeleteStruct(url string, b []byte) error {
	req, err1 := http.NewRequest("DELETE", url, bytes.NewBuffer(b))
	if err1 != nil {
		return err1
	}
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err2 := client.Do(req)
	defer resp.Body.Close()

	return err2
}
