package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
)

type projectName struct {
	Project []string `json:"project"`
}
type helpNameStruct struct {
	Name string `json:"name"`
}
type responseStruct struct {
	Languages []string
	Auth      bool
}

func getLanguages(w http.ResponseWriter, auth string, limit string, r *http.Request) {
	stopID := getNumber()
	auththen2 := ""
	if auth != "" {
		auththen2 = "?private_token=" + auth
	}
	defer r.Body.Close()
	proj := &projectName{}
	err := json.NewDecoder(r.Body).Decode(proj)
	if err != nil {
		fmt.Println("Error decoding line 20")
	}
	var fulldata = make(map[string]int)
	tallA := make([]int, len(proj.Project))
	fmt.Println(proj.Project)
	lmt := 0
	if limit == "" {
		lmt = 5
	} else {
		lmt, err = strconv.Atoi(limit)
		if err != nil {
			fmt.Println("Error converting limit to an integer line 26")
		}
	}
	fmt.Println(len(proj.Project))
	if len(proj.Project) == 0 {
		for i := 1; i <= stopID; i++ {
			currentNum := strconv.Itoa(i)
			resp, ok := http.Get("https://git.gvk.idi.ntnu.no/api/v4/projects/" + currentNum + "/languages" + auththen2)
			if ok != nil {
				fmt.Println("Error getting z loop")
			}
			defer resp.Body.Close()
			var data map[string]float32
			json.NewDecoder(resp.Body).Decode(&data)
			for index, element := range data {
				fmt.Println(element)
				if index != "message" {
					fulldata[index] = fulldata[index] + 1
					fmt.Println("Eureka: ")
					fmt.Println(fulldata)
				}
			}
		}
	} else {
		for i := 0; i <= stopID; i++ {
			strint := strconv.Itoa(i)
			urlstring := apiURL + strint + auththen2
			fmt.Println(urlstring)
			resp, ok := http.Get(urlstring)
			if ok != nil {
				fmt.Println("Error looping http.Get")
			}
			defer resp.Body.Close()
			helpStruct := &helpNameStruct{}
			ok = json.NewDecoder(resp.Body).Decode(helpStruct)
			if ok != nil {
				fmt.Println(ok)
				fmt.Println("Error decoding in loop in cycle: ")
				fmt.Println(i)
			}
			for j := 0; j < len(proj.Project); j++ {
				if helpStruct.Name == proj.Project[j] {
					tallA[j] = i
				}
			}
		}
		for z := 0; z < len(tallA); z++ {
			currentNum := strconv.Itoa(tallA[z])
			resp, ok := http.Get("https://git.gvk.idi.ntnu.no/api/v4/projects/" + currentNum + "/languages" + auththen2)
			if ok != nil {
				fmt.Println("Error getting z loop")
			}
			defer resp.Body.Close()
			var data map[string]float32
			json.NewDecoder(resp.Body).Decode(&data)
			for index, element := range data {
				fmt.Println(element)
				if index != "message" {
					fulldata[index] = fulldata[index] + 1
				}
			}
		}
	}
	fmt.Println(fulldata)
	if lmt > len(fulldata) {
		lmt = len(fulldata)
	}
	fmt.Println("LMT: ")
	fmt.Println(lmt)
	intArray := make([]int, len(fulldata))
	intArraySorted := make([]int, len(fulldata))
	stringArray := make([]string, len(fulldata))
	stringArrayResponse := make([]string, lmt)
	for index, element := range fulldata {
		if index != "" {
			fmt.Println(element)
			output1 := append(intArray, element)
			output2 := append(intArraySorted, element)
			output3 := append(stringArray, index)
			intArray = output1
			intArraySorted = output2
			stringArray = output3
		}
	}
	sort.Ints(intArraySorted)
	loopn := 1
	stoploop := 0
	fmt.Println("len(intarray): ")
	fmt.Println(len(intArray))
	fmt.Println(lmt)
	fmt.Println("len(stringarrayresponse): ")
	fmt.Println(len(stringArrayResponse))
	for i := 0; i < len(stringArrayResponse); i++ {
		stoploop = 0
		for j := 0; j < len(intArray); j++ {
			if intArraySorted[len(intArray)-loopn] == intArray[j] && stoploop == 0 {
				intArray[j] = 0
				stringArrayResponse[i] = stringArray[j]
				loopn++
				stoploop = 1
			}
		}
	}
	fmt.Println(stringArrayResponse)
	answerpointer := &responseStruct{}
	answerpointer.Languages = stringArrayResponse
	if auth == "" {
		answerpointer.Auth = false
	} else {
		answerpointer.Auth = true
	}
	err = json.NewEncoder(w).Encode(answerpointer) //encodes it to the w parameter
	if err != nil {
		fmt.Println("Error Encoding")
		fmt.Fprintf(w, "500"+http.StatusText(500))
		return
	}

}
