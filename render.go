package main

import (
	"bytes"
	"fmt"
	"github.com/canghel3/go-wms/render"
	"github.com/omniscale/go-mapnik"
	"image"
	"image/color"
	"image/png"
	"log"
	"os"
	"time"

	"github.com/Ginger955/gdal"
	"github.com/airbusgeo/godal"
)

func renderGdal(input string) {
	s := time.Now()
	// Define the bounding box (xmin, ymin, xmax, ymax)
	bbox := [4]float64{2320628.1787, 5690583.8818, 2417550.3306, 5779556.5827} // Example coordinates

	// Define the output image width and height (zoom level can affect these)
	outputWidth := 1000
	outputHeight := 1000

	// Open the .tif file using GDAL
	dataset, err := gdal.Open(input, gdal.ReadOnly)
	if err != nil {
		log.Fatalf("Failed to open TIFF file: %v", err)
	}
	defer dataset.Close()

	// Perform warp operation to crop and resample the image according to bbox, width, and height
	options := []string{
		"-of", "MEM",
		"-t_srs", "EPSG:3857",
		"-s_srs", "EPSG:3857",
		"-ts", fmt.Sprintf("%d", outputWidth), fmt.Sprintf("%d", outputHeight),
		"-te", fmt.Sprintf("%f", bbox[0]), fmt.Sprintf("%f", bbox[1]), fmt.Sprintf("%f", bbox[2]), fmt.Sprintf("%f", bbox[3]),
		"-ot", "Float32",
	}

	warpedDataset, err := gdal.Warp("", nil, []gdal.Dataset{dataset}, options)
	if err != nil {
		log.Fatalf("Failed to warp the dataset: %v", err)
	}
	defer warpedDataset.Close()

	// Read the raster band from the warped dataset (assume grayscale for simplicity)
	band := warpedDataset.RasterBand(1)

	// Read the resampled raster data into a buffer
	data := make([]float32, outputWidth*outputHeight)
	err = band.IO(gdal.Read, 0, 0, outputWidth, outputHeight, data, outputWidth, outputHeight, 0, 0)
	if err != nil {
		log.Fatalf("Failed to read raster data: %v", err)
	}

	// Create a new grayscale image
	img := image.NewGray(image.Rect(0, 0, outputWidth, outputHeight))

	// Convert the TIFF data to a grayscale image
	for y := 0; y < outputHeight; y++ {
		for x := 0; x < outputWidth; x++ {
			value := data[y*outputWidth+x]
			grayValue := uint8(value * 255)
			img.SetGray(x, y, color.Gray{Y: grayValue})
		}
	}

	// Create a PNG file to save the output
	outputFile, err := os.Create("output.png")
	if err != nil {
		log.Fatalf("Failed to create PNG file: %v", err)
	}
	defer outputFile.Close()

	// Encode the image to PNG format
	err = png.Encode(outputFile, img)
	if err != nil {
		log.Fatalf("Failed to encode PNG: %v", err)
	}

	log.Println("PNG file successfully created with bounding box and dimensions!")
	log.Println("elapsed: ", time.Since(s))
}

func renderGodal(input string, bbox [4]float64, width, height int) ([]byte, error) {
	s := time.Now()
	// Open the .tif file using godal
	ds, err := godal.Open(input)
	if err != nil {
		log.Fatalf("Failed to open TIFF file: %v", err)
	}
	defer ds.Close()

	switch len(ds.Bands()) {
	case 1:
		img, err := handleSingleBand(ds, bbox, width, height)
		if err != nil {
			return nil, err
		}

		var content bytes.Buffer
		err = png.Encode(&content, img)
		if err != nil {
			log.Fatalf("Failed to encode PNG: %v", err)
		}

		return content.Bytes(), nil
	}

	log.Println("elapsed: ", time.Since(s))
	return nil, nil
}

func handleSingleBand(ds *godal.Dataset, bbox [4]float64, width, height int) (image.Image, error) {
	//bbox := [4]float64{1878516.407136492, 5635549.221409476, 2035059.441064533, 5792092.255337515}
	//width := 256
	//// Calculate the aspect ratio of the bounding box
	//aspectRatio := (bbox[2] - bbox[0]) / (bbox[3] - bbox[1])
	//// Calculate the height based on the aspect ratio
	//height := int(float64(width)/aspectRatio) + 1

	switches := []string{
		"-te", fmt.Sprintf("%f", bbox[0]), fmt.Sprintf("%f", bbox[1]), fmt.Sprintf("%f", bbox[2]), fmt.Sprintf("%f", bbox[3]),
		"-te_srs", "EPSG:3857",
		"-ts", fmt.Sprintf("%d", width), fmt.Sprintf("%d", height),
		"-s_srs", "EPSG:3857",
		"-t_srs", "EPSG:3857",
		"-of", "MEM",
	}

	ds, err := ds.Warp("", switches)
	if err != nil {
		log.Printf("Failed to translate: %v\n", err)
		return nil, err
	}
	defer ds.Close()

	// Get raster band (assuming the first band, for simplicity)
	band := ds.Bands()[0]

	// Read the raster data into a float64 slice
	data := make([]float64, width*height)
	err = band.Read(0, 0, data, width, height)
	if err != nil {
		log.Printf("Failed to read raster data: %v", err)
		return nil, err
	}

	//gs := render.Gray(data, width, height, true)
	gs := render.RGB(data, width, height)
	render, err := gs.Render()
	if err != nil {
		log.Printf("Failed to render: %v", err)
		return nil, err
	}

	return render, nil

	//// Create the output PNG file
	//outputFile, err := os.Create("output.png")
	//if err != nil {
	//	log.Printf("Failed to create PNG file: %v", err)
	//	return nil, err
	//}
	//defer outputFile.Close()

	//// Encode the image to PNG format
	//err = png.Encode(outputFile, render)
	//if err != nil {
	//	log.Printf("Failed to encode PNG: %v", err)
	//	return nil, err
	//}
	//
	//log.Println("Styled PNG file successfully created!")
}

func nik() {
	mapnik.LogSeverity(mapnik.Debug)
	m := mapnik.New()

	err := m.Load("/home/cristian/GolandProjects/go-wms/mapnik_sample.xml")
	if err != nil {
		log.Fatalf("Failed to load mapnik: %v", err)
	}

	bbox := [4]float64{2353018.3695, 5733388.6176, 2367713.3882, 5748064.5270}
	m.ZoomTo(bbox[0], bbox[1], bbox[2], bbox[3])
	//m.SetMaxExtent(bbox[0], bbox[1], bbox[2], bbox[3])
	m.Resize(256, 256)
	m.SetSRS("epsg:3857")
	fmt.Println(m.SRS())
	nrgba, err := m.RenderImage(mapnik.RenderOpts{
		Scale:       1,
		ScaleFactor: 1,
		Format:      "PNG",
	})
	if err != nil {
		log.Fatalf("Failed to render image: %v", err)
	}

	encode, err := mapnik.Encode(nrgba, "png")
	if err != nil {
		log.Fatalf("Failed to encode image: %v", err)
	}

	// Create the output PNG file
	outputFile, err := os.Create("output2.png")
	if err != nil {
		log.Fatalf("Failed to create PNG file: %v", err)
	}
	defer outputFile.Close()

	img, err := png.Decode(bytes.NewReader(encode))
	if err != nil {
		log.Fatalf("Failed to decode PNG: %v", err)
	}
	err = png.Encode(outputFile, img)
	if err != nil {
		log.Fatalf("Failed to encode PNG: %v", err)
	}

	log.Println("DONE!")
}
