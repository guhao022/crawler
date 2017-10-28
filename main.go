package main

import (
	"fmt"
	"strings"
)

/*import (
	"github.com/num5/env"
)

func main() {
	_, err := env.Load()
	if err != nil {
		panic(err)
	}

	opts := NewOptions()

	Run(opts)
}*/

func main() {
	s := "//image.haha.mx/2017/10/27/middle/2605964_4ebd2e73d8c03b7b29c8d1440614b99a_1509098142.gid"
	ss := strings.HasSuffix(s, ".gif")
	fmt.Printf("%v",ss)
	/*s = strings.Replace(s, "", "|", -1)
	fmt.Println(s)*/
}


