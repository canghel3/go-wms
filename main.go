package main

import (
	"fmt"
	"os"
)

func main() {
	//renderGdal("/home/cristian/Documents/quarticle/nasa_lights_compr.tif")
	//err := renderGodal("/home/cristian/Documents/quarticle/nasa_lights_compr.tif")
	//if err != nil {
	//	panic(err)
	//}

	fmt.Println("inp: ", os.Getenv("MAPNIK_INPUT_PLUGINS"))
	nik()
}
