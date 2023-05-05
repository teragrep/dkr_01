package main

import (
	"crypto/tls"
	"crypto/x509"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

type KubernetesAPI struct {
	clientKeyPath     string                  // used for tls.Certificate X509 cert
	clientCertPath    string                  // used for tls.Certificate X509 key
	serverCertPath    string                  // used for root CA cert
	clientAuthEnabled bool                    // uses x509 cert/key
	serverAuthEnabled bool                    // ca cert
	url               string                  // url should be in form of "https://localhost:10250"
	data              *map[string]interface{} // unmarshaled json response data
	metadata          map[string]string       // extracted metadata from json response
	containerStatus   map[string]interface{}  // extracted current container status json data
}

func (kapi *KubernetesAPI) FetchData() error {
	client := &http.Client{}
	tlsConfig := &tls.Config{}

	// client certs
	if kapi.clientAuthEnabled {
		if _, err := os.Stat(kapi.clientCertPath); err == nil {
			// client cert path exists
			if _, err2 := os.Stat(kapi.clientKeyPath); err2 == nil {
				// client key path exists
			} else {
				// client key path doesn't exist
				return errors.New("client key path doesn't exist: '" + kapi.clientKeyPath + "'")
			}
		} else {
			// client cert path doesn't exist
			return errors.New("client cert path doesn't exist: '" + kapi.clientCertPath + "'")
		}

		cert, err := tls.LoadX509KeyPair(kapi.clientCertPath, kapi.clientKeyPath)
		if err != nil {
			return err
		}
		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	// server (CA) certs
	if kapi.serverAuthEnabled {
		if _, err := os.Stat(kapi.serverCertPath); err == nil {
			caCert, err := os.ReadFile(kapi.serverCertPath)
			if err != nil {
				return err
			}
			caCertPool := x509.NewCertPool()
			caCertPool.AppendCertsFromPEM(caCert)

			tlsConfig.RootCAs = caCertPool

		} else {
			return errors.New("Server cert path pointed to a file that does not exist: '" + kapi.serverCertPath + "'")
		}

	}

	// set tls config
	client.Transport = &http.Transport{TLSClientConfig: tlsConfig}

	// get '/pods'
	resp, err := client.Get(fmt.Sprintf("%s/pods", kapi.url))
	if err != nil {
		return err
	}

	// read body of response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// get json
	var objMap map[string]interface{}
	if err := json.Unmarshal(body, &objMap); err != nil {
		return err
	}

	kapi.data = &objMap
	return nil
}

func (kapi *KubernetesAPI) GetContainerData(cntId string) error {
	// check if data has been fetched
	if kapi.data == nil {
		return errors.New("KubernetesAPI data was nil, please retrieve data with FetchData() first")
	}

	// .items[] > status > containerStatuses[] > containerID
	items, _ := (*kapi.data)["items"] // get "items" array
	item := items.([]interface{})[0]  // get first item of "items" array

	status := item.(map[string]interface{})["status"]                         // get "status" map
	containerStatuses := status.(map[string]interface{})["containerStatuses"] // get "containerStatuses" array

	// get status for container id
	statuses := containerStatuses.([]interface{})
	for i := 0; i < len(statuses); i++ {
		statusItem := statuses[i]
		if statusItem != nil {
			idString := statusItem.(map[string]interface{})["containerID"].(string)
			if idString == cntId {
				// matching container id
				kapi.containerStatus = statusItem.(map[string]interface{})
				break
			}
		}
	}

	// prepare metadata map
	kapi.metadata = make(map[string]string)

	// .items[] > metadata > name
	metadata := item.(map[string]interface{})["metadata"]
	name := metadata.(map[string]interface{})["name"]
	kapi.metadata["name"] = name.(string)
	// .items[] > metadata > namespace
	namespace := metadata.(map[string]interface{})["namespace"]
	kapi.metadata["namespace"] = namespace.(string)
	// .items[] > metadata > labels
	labels := metadata.(map[string]interface{})["labels"]
	jsonLabels, err := json.Marshal(&labels)
	if err != nil {
		jsonLabels = []byte("null")
	}
	kapi.metadata["labels"] = string(jsonLabels)
	// .items[] > metadata > uid
	uid := metadata.(map[string]interface{})["uid"]
	kapi.metadata["uid"] = uid.(string)
	// .items[] > metadata > creationTimestamp
	creationTimestamp := metadata.(map[string]interface{})["creationTimestamp"]
	if creationTimestamp != nil {
		kapi.metadata["creationTimestamp"] = creationTimestamp.(string)
	} else {
		kapi.metadata["creationTimestamp"] = "null"
	}

	//fmt.Println(kapi.metadata)
	//fmt.Println(kapi.containerStatus)

	return nil
}
