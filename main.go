package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

	"github.com/savioxavier/termlink"
)

var api string = "https://api.github.com/"

type userResponse struct {
	Username string `json:"login"`
	AvatarURL string `json:"avatar_url"`
	ID int64 `json:"id"`
	Type string `json:"type"`
	Followers string `json:"followers_url"`
	
}

type reposResponse struct {
	URL string `json:"url"`
	Name string `json:"name"`
	Private bool `json:"private"`
	Description string `json:"description"`
	Language string `json:"language"`
	GitURL string `json:"git_url"`
	ForksCount int64 `json:"forks_count"`
	Owner reposOwner `json:"owner"`
	
}

type reposOwner struct {
	userResponse

}


func fetchUser(username string) {
	resp, err := http.Get(fmt.Sprint(api + "users/", username))
	if err != nil {
		fmt.Println(err)
	}

	var user userResponse
	bytes, err :=  io.ReadAll(resp.Body)
	jsonData := json.Unmarshal(bytes, &user)
	if jsonData != nil {
		fmt.Println(jsonData)
	}
	if err != nil {
		fmt.Println(err)
		return
	}

	defer resp.Body.Close()

	followers := []string{}

	respFollowers, errFollowers := http.Get(user.Followers)
	if errFollowers != nil{
		fmt.Println(errFollowers)
	}

	bytesFollowers, _ := io.ReadAll(respFollowers.Body)
	data := json.Unmarshal(bytesFollowers, &followers)

	if data != nil {
		fmt.Println(data)
	}

	defer respFollowers.Body.Close()

	countFollowers := len(followers)


	fmt.Println("Username:", user.Username)
	fmt.Println("Avatar URL: ", termlink.Link("Click here!", user.AvatarURL))
	fmt.Println("ID: ",user.ID)
	fmt.Println("Followers: ",countFollowers)
	fmt.Println("Type:", user.Type)

}

func fetchRepos(username string) {

	var repos []reposResponse

	resp, _ := http.Get(fmt.Sprint(api + "users/", username + "/repos"))
	bytes, _ := io.ReadAll(resp.Body)
	jsonData := json.Unmarshal(bytes, &repos,)

	if jsonData != nil {
		fmt.Println(jsonData)
	}

	for _, repo := range repos {
		fmt.Println(termlink.Link(repo.Name, repo.URL), "-", repo.Language, "-", repo.Description)
	}

	defer resp.Body.Close()


}

func main() {
	var username string
	flag.StringVar(&username, "user", "", "usernamevar")

	for {

		flag.Parse()

		fmt.Println("Username:", username)

		fmt.Println("-----------------------------------")
		fmt.Println("Github CLI made by darknessmalware.")
		fmt.Println("-----------------------------------")

		fmt.Println("Options: ")

		fmt.Println("[1] - Fetch User`s Profile Information")
		fmt.Println("[2] - Fetch User`s Repositorys")
		fmt.Println("[0] - Exit")

		var choice int

		fmt.Print("Option: ")
		fmt.Scan(&choice)
		
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()

		switch choice {
		case 1:
			fetchUser(username)
		
		case 2:
			fetchRepos(username)

		case 0:
			os.Exit(1)

		}
	}


	

}