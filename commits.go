package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
)

const apiURL = "https://git.gvk.idi.ntnu.no/api/v4/projects/"

type getAmnt []struct {
	ShortID string `json:"short_id"`
}
type getName struct {
	PathWithNamespace string `json:"path"`
}

//asd
type RepositoryStruct struct {
	Repository string
	Commits    int
}

type commitStruct struct {
	Repos []RepositoryStruct
	Auth  bool
}

func getCommits(w http.ResponseWriter, lmt string, vrf string) {
	booly := false
	auththen := ""
	if vrf != "" {
		booly = true
		auththen = "private_token=" + vrf
	}
	fmt.Println(auththen)
	stopID := (getNumber())
	if lmt == "" {
		lmt = "5"
	}
	limit, err := strconv.Atoi(lmt)
	if err != nil {
		fmt.Println("Error converting lmt to an integer")
	}
	tallA := make([]int, stopID)
	tallB := make([]int, stopID)
	tallC := make([]int, limit)
	idA := make([]string, limit)
	finito := make([]RepositoryStruct, limit)
	helpStruct := &getAmnt{}

	for i := 0; i < stopID; i++ {
		strint := strconv.Itoa(i)
		pageNum := 0
		urlstring := apiURL + strint + "/repository/commits?all=true&per_page=100&" + auththen
		fmt.Println(urlstring)
		resp, ok := http.Get(urlstring)
		if ok != nil {
			fmt.Println("Error looping http.Get")
		}
		defer resp.Body.Close()
		ok = json.NewDecoder(resp.Body).Decode(helpStruct)
		if ok != nil {
			fmt.Println(ok)
			fmt.Println("Error decoding in loop in cycle: ")
			fmt.Println(i)
			*helpStruct = getAmnt{}
		}
		if len(*helpStruct) != 100 {
			tallA[i] = len(*helpStruct)
			tallB[i] = len(*helpStruct)
			fmt.Println(tallA)
		} else if len(*helpStruct) == 100 {
			tooLong := 0
			for tooLong%100 == 0 || tooLong == 0 {
				pageNum++
				pagenumber := strconv.Itoa(pageNum)
				urlstring = apiURL + strint + "/repository/commits?all=true&per_page=100&page=" + pagenumber + "&" + auththen
				resp, ok := http.Get(urlstring)
				if ok != nil {
					fmt.Println("Error looping http.Get")
				}
				defer resp.Body.Close()
				ok = json.NewDecoder(resp.Body).Decode(helpStruct)
				if ok != nil {
					fmt.Println("Printing 0, no access on line: ")
					fmt.Println(i)
					*helpStruct = getAmnt{}
				}
				tooLong += len(*helpStruct)
			}
			tallA[i] = tooLong
			tallB[i] = tooLong
			fmt.Println(tallA)
		}
	}
	sort.Ints(tallB)
	fmt.Println(tallB)
	number := 0
	for z := stopID - 1; z >= stopID-limit; z-- {
		tallC[number] = tallB[z]
		number++
	}
	fmt.Println(tallC)
	ehhn := 0
	arraynumber := 0
	helpName := &getName{}
	for k := 0; k < limit; k++ {
		for i := 0; i < stopID; i++ {
			if ehhn != limit {
				if tallC[ehhn] == tallA[i] {
					ehhn++
					strint := strconv.Itoa(i)
					urlstring := apiURL + strint + "?" + auththen
					resp, ok := http.Get(urlstring)
					fmt.Println(urlstring)
					if ok != nil {
						fmt.Println(ok)
					}
					defer resp.Body.Close()
					ok = json.NewDecoder(resp.Body).Decode(helpName)
					if ok != nil {
						fmt.Println(ok)
					}
					idA[arraynumber] = (*helpName).PathWithNamespace
					arraynumber++
				}
			}
		}
	}

	fmt.Println(idA)
	fmt.Println(tallC)
	for x := 0; x < limit; x++ {
		finito[x] = RepositoryStruct{idA[x], tallC[x]}
	}
	RepositoryAnswer := &commitStruct{}
	RepositoryAnswer.Repos = finito
	RepositoryAnswer.Auth = booly
	fmt.Println(RepositoryAnswer)
	fmt.Println(finito)
	err = json.NewEncoder(w).Encode(RepositoryAnswer) //encodes it to the w parameter
	if err != nil {
		fmt.Println("Error Encoding")
		fmt.Fprintf(w, "500"+http.StatusText(500))
		return
	}

}

type getID []struct {
	ID int `json:"id"`
}

func getNumber() int {
	resp, err := http.Get(apiURL)
	if err != nil { //If error, print 520 (unknown error)
		fmt.Println("Error with getting nURL")
		fmt.Println(err)
	}
	defer resp.Body.Close() //dereferencing & closing response
	topid := &getID{}
	err = json.NewDecoder(resp.Body).Decode(topid)
	if err != nil {
		fmt.Println(err)
		fmt.Println(resp.Body)
		fmt.Println("Error with topid Decoding")
	}

	n := *topid
	return n[0].ID

}
