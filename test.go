package main

import (
    "fmt"
    "strings"
)

func convertYouTubeLink(link string) string {
    if idx := strings.Index(link, "&"); idx != -1 {
        link = link[:idx]
    }

   	link = strings.Replace(link, "watch?v=", "embed/", 1)

    return link
}

func main() {
    original := "https://www.youtube.com/watch?v=SQLDFttWurk&ab_channel=WBRussia"
    converted := convertYouTubeLink(original)
    fmt.Println("Преобразованная ссылка:", converted)
}
