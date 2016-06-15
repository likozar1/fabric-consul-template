package main

import(
	 "fmt"
	 "net/http"
	 "io/ioutil"
	 "strings"
	 "crypto/tls"
)

//REST request call
func requestREST(REST string, url string, data []byte)([]byte, string){

	//checks config file for ssl
	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	if c.Ssl["ssl"]{
		tr = &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: false}}
	}
	
    client := &http.Client{Transport: tr}

	request, err := http.NewRequest(REST, url, nil)
    
    if(data != nil){
    	payload := strings.NewReader(string(data))
    	request, err = http.NewRequest(REST, url, payload)	
    }

	//request.SetBasicAuth(USERNAME, PASSWORD)

	//send request
	response, err := client.Do(request)
	if err != nil {
		return nil, fmt.Sprintf("ERROR-%s", err)
	} else {
		//read response
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, fmt.Sprintf("ERROR-%s", err)
		}
		return contents, response.Status
	}
    return nil, fmt.Sprintf("ERROR-%s", err)
}

//Send GET request
func getRequest(url string) ([]byte, string) {
	return requestREST(GET,url,nil)
}

//Send DELETE request
func deleteRequest(url string)([]byte, string){
	return requestREST(DELETE,url,nil)
}

//Send PUT request
func putRequest(url string, data []byte)([]byte, string){
	return requestREST(PUT,url,data)
}

//Send POST request
func postRequest(url string, data []byte) ([]byte, string){
	return requestREST(POST,url,data)
}