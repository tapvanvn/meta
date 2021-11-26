package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/tapvanvn/goutil"
)

const (
	ResponseNotFound         = "not_found"
	ResponseErrorType        = "wrong_type"
	ResponseErrorUnsupported = "unsupported"
)

var supported = map[string]bool{
	".jsonc": true,
	".json":  true,
}

func ResponseError(reason string) {
	fmt.Print(fmt.Sprintf("error.\"%s\"", reason))
	os.Exit(1)
}

func Response(result string) {
	fmt.Print(result)
	os.Exit(0)
}

//ars:
// #1: input_file
// #2: type
// #3: path
//example:
// #1: meta length master.nodes 	`return the length of master.nodes stack`
// #2: meta offset.0 master.nodes 	`return the first item on of master.nodes stack in string -> ""`
// #3: meta value master.nodes.0 	`return the value in string of the `
// #4: meta key master.nodes.0 		`return the key in string of the first element in master.nodes stack`

type pair struct {
	key   string
	value interface{}
}

func main() {
	numArg := len(os.Args)
	if numArg != 4 {
		log.Fatal("invalid syntax: meta <input_file> <info_type> <info_path>")
	}

	clusterInfoFile := os.Args[1]

	dot := strings.Index(clusterInfoFile, ".")
	ext := strings.ToLower(clusterInfoFile[dot:])

	if _, isSupport := supported[ext]; !isSupport {

		ResponseError(ResponseErrorUnsupported)
	}
	file, err := os.Open(clusterInfoFile)

	if err != nil {
		ResponseError("cluster info file is not existed!")
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		ResponseError("reading cluster info err!")
	}
	if ext == ".jsonc" {
		bytes = goutil.TripJSONComment(bytes)
	}
	meta := &map[string]interface{}{}
	err = json.Unmarshal(bytes, meta)
	if err != nil {
		ResponseError("parse cluster info err!")

	}
	parts := strings.Split(os.Args[3], ".")
	var iter *map[string]interface{} = meta
	offset := 0
	numPart := len(parts)

	for {
		if iter == nil {
			ResponseError("bad request")

		}
		label := parts[offset]
		var value interface{} = nil

		valueOffset, err := strconv.Atoi(label)
		if err == nil {
			count := 0

			for key, elementValue := range *iter {

				if count == valueOffset {

					value = map[string]interface{}{
						key: elementValue,
					}
				}
				count++
			}
		} else {
			if testValue, ok := (*iter)[label]; !ok {

				ResponseError(fmt.Sprintf("bad request at %s", label))

			} else {

				value = testValue
			}
		}
		if value == nil {

			Response(ResponseNotFound)
		}

		if offset == numPart-1 { //end of the map

			export(os.Args[2], value)
			break
		} else {
			offset++
			mval := value.(map[string]interface{})
			iter = &mval
		}
	}
	//log.Printf("%v", meta)
	//log.Printf("%v", parts)
}

//export the infomation
//example:
// #1: meta length master.nodes 	`return the length of master.nodes stack`
// #2: meta value master.nodes.0 	`return the value in string of the `
// #3: meta key master.nodes.0 		`return the key in string of the first element in master.nodes stack`
func export(infoType string, value interface{}) {

	switch infoType {
	case "length":
		exportLength(value)
	case "value":
		exportValue(value)
	case "key":
		exportKey(value)
	default:
		Response(ResponseNotFound)
	}
}

func exportKey(value interface{}) {

	switch value.(type) {

	case map[string]interface{}:

		for key, _ := range value.(map[string]interface{}) {
			Response(key)
		}
	default:
		ResponseError(ResponseErrorType)
	}
}
func exportValue(value interface{}) {

	switch value.(type) {
	case string:
		Response(value.(string))
	case map[string]interface{}:
		for _, val := range value.(map[string]interface{}) {
			Response(fmt.Sprint(val))
		}
	default:
		Response(fmt.Sprint(value))
	}
}
func exportLength(value interface{}) {

	switch value.(type) {
	case string:
		Response(fmt.Sprint(len(value.(string))))
	case map[string]interface{}:
		Response(fmt.Sprint(len(value.(map[string]interface{}))))
	default:
		ResponseError(ResponseErrorType)
	}
}
