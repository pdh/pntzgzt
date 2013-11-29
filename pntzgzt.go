package main

import (
        "bytes"
        "encoding/json"
        "flag"
        "fmt"
        "io/ioutil"
        "net/http"
        "os"
)

var GIST_URL = "https://api.github.com/gists"

type Config struct {
        Username string
        Password string
}

type FileMeta struct {
        Content string `json:"content"`
}

type Gist struct {
        Description string              `json:"description"`
        Public      bool                `json:"public"`
        Files       map[string]FileMeta `json:"files"`
}

func PostGist(config Config, gist []byte) (*http.Response, error) {
        client := &http.Client{}
        req, err := http.NewRequest("POST", GIST_URL, bytes.NewReader(gist))
        if err != nil {
                fmt.Printf("Error creating request: %v\n")
                os.Exit(1)
        }
        req.SetBasicAuth(config.Password, "x-oauth-basic")
        return client.Do(req)
}

func ParseGistURL(body []byte) string {
        var response interface{}
        err = json.Unmarshal(body, &response)
        if err != nil {
                fmt.Printf("Failed to parse URL: %v\n", err)
                os.Exit(1)
        }
        m := response.(map[string]interface{})
        return m["html_url"]
}

func main() {
        file, err := ioutil.ReadFile("./config.json")
        if err != nil {
                fmt.Printf("Error reading config: %v\n", err)
                os.Exit(1)
        }
        var conf Config
        json.Unmarshal(file, &conf)

        if len(os.Args) < 2 {
                fmt.Println("usage: pntzgzt <filename> (options)")
                os.Exit(1)
        }
        gistFileName := os.Args[1]

        content, err := ioutil.ReadFile(gistFileName)
        if err != nil {
                fmt.Printf("Error reading input file: %v\n", err)
                os.Exit(1)
        }

        description := flag.String("desc", gistFileName, "")
        is_public := flag.Bool("public", false, "")

        file_meta := FileMeta{string(content)}
        file_map := make(map[string]FileMeta)
        file_map[gistFileName] = file_meta
        gist := Gist{*description, *is_public, file_map}

        gist_json, err := json.Marshal(gist)
        if err != nil {
                fmt.Printf("Unable to create json: %v\n", err)
        }

        res, err := PostGist(conf, gist_json)
        if err != nil {
                fmt.Println("fuk")
        }
        defer res.Body.Close()
        body, err := ioutil.ReadAll(res.Body)

        url, err := ParseGistURL(body)
        fmt.Printf("%v\n", url)
}
