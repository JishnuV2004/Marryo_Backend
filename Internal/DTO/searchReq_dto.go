package dto
// dto/search.go
type SearchRequest struct {
	Name string `query:"name"`
	Age  int    `query:"age"`
	Star string `query:"star"`
}
