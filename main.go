package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

func main() {

	server := http.Server{
		Addr: "localhost:9000",
	}
	i := 0
	http.HandleFunc("/geoserver/wms", func(w http.ResponseWriter, r *http.Request) {
		request := r.URL.Query().Get("request")
		if request == "GetCapabilities" {
			content, err := os.ReadFile("./capabilities.xml")
			if err != nil {
				log.Printf("failed to read capabilities: %v", err)
				panic(err)
			}

			w.Header().Set("Content-Type", "application/xml")
			_, err = w.Write(content)
			if err != nil {
				log.Printf("failed to write response: %v", err)
				panic(err)
			}
		} else if request == "GetMap" {
			log.Println("map request: ", i)
			i++
			width := r.URL.Query().Get("width")
			height := r.URL.Query().Get("height")
			bbox := r.URL.Query().Get("bbox")

			wh, err := strconv.ParseInt(width, 10, 32)
			if err != nil {
				panic(err)
			}

			hh, err := strconv.ParseInt(height, 10, 32)
			if err != nil {
				panic(err)
			}

			var flBbox [4]float64
			spl := strings.Split(bbox, ",")
			for i, s := range spl {
				fl, err := strconv.ParseFloat(s, 64)
				if err != nil {
					panic(err)
				}
				flBbox[i] = fl
			}
			s := time.Now()
			data, err := renderGodal("/home/cristian/Documents/quarticle/nasa_lights_compr.tif", flBbox, int(wh), int(hh))
			if err != nil {
				panic(err)
			}

			fmt.Println("elapsed:", time.Since(s))
			s = time.Now()

			_, err = w.Write(data)
			if err != nil {
				panic(err)
			}
		}
	})

	err := server.ListenAndServe()
	if err != nil {
		panic(err)
	}
	//renderGdal("/home/cristian/Documents/quarticle/nasa_lights_compr.tif")

	//TODO: implement map rendering with mapnik and simple gdal reading and image formatting
	//fmt.Println("inp: ", os.Getenv("MAPNIK_INPUT_PLUGINS"))
	//nik()
	//fmt.Println("elapsed 2:", time.Since(s))

}
