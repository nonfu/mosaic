package common

import (
    "fmt"
    "image"
    "image/color"
    "io/ioutil"
    "math"
    "os"
    "sync"
)

var TILESDB map[string][3]float64

func CloneTilesDB() DB {
    store := make(map[string][3]float64)
    for k, v := range TILESDB {
        store[k] = v
    }
    db := DB{
        store: store,
        mutex: &sync.Mutex{},
    }
    return db
}

func averageColor(img image.Image) [3]float64 {
    bounds := img.Bounds()
    r, g, b := 0.0, 0.0, 0.0
    for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
        for x := bounds.Min.X; x < bounds.Max.X; x++ {
            r1, g1, b1, _ := img.At(x, y).RGBA()
            r, g, b = r+float64(r1), g+float64(g1), b+float64(b1)
        }
    }
    totalPixels := float64(bounds.Max.X * bounds.Max.Y)
    return [3]float64{r / totalPixels, g / totalPixels, b / totalPixels}
}

func Resize(in image.Image, newWidth int) image.NRGBA {
    bounds := in.Bounds()
    width := bounds.Max.X - bounds.Min.X
    ratio := width / newWidth
    out := image.NewNRGBA(image.Rect(bounds.Min.X/ratio, bounds.Min.X/ratio, bounds.Max.X/ratio, bounds.Max.Y/ratio))
    for y, j := bounds.Min.Y, bounds.Min.Y; y < bounds.Max.Y; y, j = y+ratio, j+1 {
        for x, i := bounds.Min.X, bounds.Min.X; x < bounds.Max.X; x, i = x+ratio, i+1 {
            r, g, b, a := in.At(x, y).RGBA()
            out.SetNRGBA(i, j, color.NRGBA{uint8(r), uint8(g), uint8(b), uint8(a)})
        }
    }
    return *out
}

func TilesDB() map[string][3]float64 {
    fmt.Println("开始构建嵌入图片数据库...")
    db := make(map[string][3]float64)
    files, _ := ioutil.ReadDir("tiles")
    for _, f := range files {
        name := "tiles/" + f.Name()
        file, err := os.Open(name)
        if err == nil {
            img, _, err := image.Decode(file)
            if err == nil {
                db[name] = averageColor(img)
            } else {
                fmt.Println("构建嵌入图片数据库出错：", err, name)
            }
        } else {
            fmt.Println("构建嵌入图片数据库出错：", err, "无法打开文件", name)
        }
        file.Close()
    }
    fmt.Println("完成嵌入图片数据库构建")
    return db
}

func (db *DB) Nearest(target [3]float64) string {
    var filename string
    db.mutex.Lock()
    smallest := 1000000.0
    for k, v := range db.store {
        dist := distance(target, v)
        if dist < smallest {
            filename, smallest = k, dist
        }
    }
    delete(db.store, filename)
    db.mutex.Unlock()
    return filename
}

func distance(p1 [3]float64, p2 [3]float64) float64 {
    return math.Sqrt(sq(p2[0]-p1[0]) + sq(p2[1]-p1[1]) + sq(p2[2]-p1[2]))
}

func sq(n float64) float64 {
    return n * n
}

type DB struct {
    mutex *sync.Mutex
    store map[string][3]float64
}
