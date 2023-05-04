package main

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
)

type KubernetesAPI struct {
	clientKeyPath     string
	clientCertPath    string
	serverCertPath    string
	clientAuthEnabled bool
	serverAuthEnabled bool
	// url should be in form of "https://localhost:10250"
	url             string
	data            *map[string]interface{}
	metadata        map[string]string
	containerStatus map[string]interface{}
}

func (kapi *KubernetesAPI) FetchData() error {
	resp, err := http.Get(kapi.url + "/pods")
	if err != nil {
		return err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

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
