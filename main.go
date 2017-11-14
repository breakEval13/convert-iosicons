package main


import (
	//"encoding/json"
	"fmt"
	//"os"
	//"io/ioutil"
	"github.com/nfnt/resize"
	"github.com/mkideal/cli"
	"os"
	"image/jpeg"
	"path/filepath"
	// "flag"
	"strings"
	"image"
	"image/png"
	"image/gif"
	"errors"
	//"flag"
)

func main() {
	if err := cli.Root(root,
		cli.Tree(help),
		cli.Tree(child),
	).Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var help = cli.HelpCommand("display help information")


// root command
type rootT struct {
	cli.Helper
	Path string `cli:"path" usage:"press  path"`
}

var root = &cli.Command{
	// Argv is a factory function of argument object
	// ctx.Argv() is if Command.Argv == nil or Command.Argv() is nil
	Argv: func() interface{} { return new(rootT) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*rootT)
		ctx.String("Test Path  %s\n", argv.Path)
		start(argv.Path)
		return nil
	},
}

// child command
type childT struct {
	cli.Helper
	Path string `cli:"path" usage:"press  path"`
}

var child = &cli.Command{
	Name: "order",
	Desc: "测试",
	Argv: func() interface{} { return new(childT) },
	Fn: func(ctx *cli.Context) error {
		argv := ctx.Argv().(*childT)

		ctx.String("Order Test Path %s\n", argv.Path)

		return nil
	},
}

func start(path string) {
	if path==""{
		fmt.Println("参数不存在！")
		return
	}

	f1, err := os.Open(path)
	if err != nil {
		panic(err)
	}
	defer f1.Close()

	m1, err := jpeg.Decode(f1)
	if err != nil {
		panic(err)
	}
	bounds := m1.Bounds()

	if bounds.Dx()>1024 && bounds.Dy()>1024  {
		fmt.Println("原图大于1024px")
		return
	}
	fmt.Println("源图X:",bounds.Dx())
	fmt.Println("源图Y:",bounds.Dy())


	var height int = bounds.Dy();
	var width int = bounds.Dx();


	var height2x int = (height*2/3)
	var width2x int=((width*2/3))

	var height1x int=((height/3))
	var width1x int=((width/3))
	SaveImage("@1x.jpg",ImageResize(m1,width1x,height1x))
	SaveImage("@2x.jpg",ImageResize(m1,width2x,height2x))
	SaveImage("@3x.jpg",ImageResize(m1,width,height))

	fmt.Println("处理好的图片已经保存在当前文件目录下！")

}

// 将图片保存到指定的路径
func SaveImage(p string, src image.Image) error {

	f, err := os.OpenFile(p, os.O_SYNC|os.O_RDWR|os.O_CREATE, 0666)

	if err != nil {
		return err
	}
	defer f.Close()
	ext := filepath.Ext(p)

	if strings.EqualFold(ext, ".jpg") || strings.EqualFold(ext, ".jpeg") {

		err = jpeg.Encode(f, src, &jpeg.Options{Quality: 100})

	} else if strings.EqualFold(ext, ".png") {
		err = png.Encode(f, src)
	} else if strings.EqualFold(ext, ".gif") {
		err = gif.Encode(f, src, &gif.Options{NumColors: 256})
	}
	return err
}


func ImageCopy(src image.Image, x, y, w, h int) (image.Image, error) {

	var subImg image.Image

	if rgbImg, ok := src.(*image.YCbCr); ok {
		subImg = rgbImg.SubImage(image.Rect(x, y, x+w, y+h)).(*image.YCbCr) //图片裁剪x0 y0 x1 y1
	} else if rgbImg, ok := src.(*image.RGBA); ok {
		subImg = rgbImg.SubImage(image.Rect(x, y, x+w, y+h)).(*image.RGBA) //图片裁剪x0 y0 x1 y1
	} else if rgbImg, ok := src.(*image.NRGBA); ok {
		subImg = rgbImg.SubImage(image.Rect(x, y, x+w, y+h)).(*image.NRGBA) //图片裁剪x0 y0 x1 y1
	} else {

		return subImg, errors.New("图片解码失败")
	}

	return subImg, nil
}

func ImageResize(src image.Image, w, h int) image.Image {
	return resize.Resize(uint(w), uint(h), src, resize.Lanczos3)
}