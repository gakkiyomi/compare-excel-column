package main

import "github.com/gakkiyomi/compare-excel-column/handler"

func main() {

	r := handler.SetupRouter()
	r.Run("0.0.0.0:8085")
}
