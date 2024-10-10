package main

import (
	"fmt"
	"time"
)

func main() {
	//renderGdal("/home/cristian/Documents/quarticle/nasa_lights_compr.tif")
	s := time.Now()
	err := renderGodal("/home/cristian/Documents/quarticle/nasa_lights_compr.tif")
	if err != nil {
		panic(err)
	}

	fmt.Println("elapsed:", time.Since(s))
	s = time.Now()

	//TODO: implement map rendering with mapnik and simple gdal reading and image formatting
	//fmt.Println("inp: ", os.Getenv("MAPNIK_INPUT_PLUGINS"))
	nik()
	fmt.Println("elapsed 2:", time.Since(s))

}
