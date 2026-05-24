package main

import (
	"fmt"
	"log"

	"github.com/rpstvs/json-parser-go/internal/dora"
)

const testJSONObject = `{
    "item1": ["aryitem1", "aryitem2", {"some": {"thing": "coolObj"}}],
    "item2": "simplestringvalue"
}`

func main() {
	c, err := dora.NewFromString(testJSONObject)
	if err != nil {
		log.Printf("error creating client %v", err)
	}
	result, err := c.GetString("$.item1[2].some.thing")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(result)
}
