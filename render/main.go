package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"os"
	"regexp"
	"strconv"

	"github.com/mattn/go-isatty"
)

var palette = color.Palette{
	color.White,
	color.Black,
}

func main() {
	rxReader := regexp.MustCompile(`([0-9]+)[,\]]`)
	var outfn string
	var scale int
	var framedelay int
	flag.StringVar(&outfn, "o", "", "output filename")
	flag.IntVar(&scale, "s", 16, "scale (pixels/item)")
	flag.IntVar(&framedelay, "d", 10, "frame delay in 1/00ths of a sec")
	flag.Parse()

	if scale < 1 || scale > 1024 {
		fmt.Fprintln(os.Stderr, "invalid scale:", scale, "not between 1 and 1024")
		os.Exit(1)
	}

	if framedelay < 0 {
		fmt.Fprintln(os.Stderr, "invalid frame delay:", framedelay, "not greater than or equal to 0")
		os.Exit(1)
	}

	var err error
	var infile io.ReadCloser = os.Stdin

	if len(flag.Args()) > 0 {
		infile, err = os.Open(flag.Arg(0))
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to open", flag.Arg(0), ":", err.Error())
			os.Exit(1)
		}
		defer infile.Close()
	} else if isatty.IsTerminal(os.Stdin.Fd()) {
		fmt.Fprintln(os.Stderr, "send input via stdin or filename")
	}
	var data [][]int

	scn := bufio.NewScanner(infile)
	var maxv, minv int
	for scn.Scan() {
		parsed := rxReader.FindAllStringSubmatch(scn.Text(), -1)
		values := make([]int, len(parsed))
		for i, x := range parsed {
			values[i], _ = strconv.Atoi(x[1])
			maxv = max(maxv, values[i])
			minv = min(minv, values[i])
		}
		data = append(data, values)
	}

	g := &gif.GIF{}
	for _, values := range data {
		frame := image.NewPaletted(
			image.Rectangle{
				image.Point{},
				image.Point{(2 + len(values)) * scale, (2 + maxv - minv) * scale},
			},
			palette,
		)
		g.Image = append(g.Image, frame)
		g.Delay = append(g.Delay, framedelay)
		for x, v := range values {
			x1 := (x + 1) * scale
			x2 := x1 + scale
			y1 := ((2 + maxv - minv) - (v + 1)) * scale
			y2 := ((2 + maxv - minv) - 1) * scale
			for x := x1; x < x2; x++ {
				for y := y1; y < y2; y++ {
					frame.SetColorIndex(x, y, 1)
				}
			}
		}
	}

	if len(g.Image) == 0 {
		return
	}

	g.Delay[len(g.Delay)-1] = 100

	var outfile io.WriteCloser = os.Stdout
	if len(outfn) > 0 {
		outfile, err = os.OpenFile(outfn, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
		if err != nil {
			fmt.Fprintln(os.Stderr, "failed to open", outfn, ":", err.Error())
			os.Exit(1)
		}
		defer outfile.Close()
	}

	if err := gif.EncodeAll(outfile, g); err != nil {
		fmt.Fprintln(os.Stderr, "failed to write GIF:", err.Error())
		os.Exit(1)
	}
}
