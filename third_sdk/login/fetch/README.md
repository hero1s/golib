# fetch

a http request lib with go

### how to use

GET

```
url := "http://example.com"

result, err := fetch.Cmd(fetch.Request{
    Method: "GET",
    URL:    url,
})
fmt.Println(result, err)
```

DELETE

```
url := "http://example.com"

result, err := fetch.Cmd(fetch.Request{
    Method: "DELETE",
    URL:    url,
})
fmt.Println(result, err)
```

POST

```
url := "http://example.com"

data := map[string]string{
    "go":   "golang",
    "java": "javalang",
    "rust": "rustlang",
}
body, _ := json.Marshal(data)

header := http.Header{}
header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 8_0 like Mac OS X) AppleWebKit/600.1.3 (KHTML, like Gecko) Version/8.0 Mobile/12A4345d Safari/600.1.4")

result, err := fetch.Cmd(fetch.Request{
    Method: "POST",
    URL:    url,
    Body:   body,
    Header: header,
})
fmt.Println(result, err)
```

PUT

```
url := "http://example.com"

data := map[string]string{
    "go":   "golang",
    "java": "javalang",
    "rust": "rustlang",
}
body, _ := json.Marshal(data)

header := http.Header{}
header.Add("User-Agent", "Mozilla/5.0 (iPhone; CPU iPhone OS 8_0 like Mac OS X) AppleWebKit/600.1.3 (KHTML, like Gecko) Version/8.0 Mobile/12A4345d Safari/600.1.4")

result, err := fetch.Cmd(fetch.Request{
    Method: "PUT",
    URL:    url,
    Body:   body,
    Header: header,
})
fmt.Println(result, err)
```
