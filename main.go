package main

import (
    "fmt"
    "golang.org/x/image/font"
    "golang.org/x/image/font/basicfont"
    "golang.org/x/image/math/fixed"
    "image"
    "image/color"
    "image/png"
    "log"
    "os"
)

type HelixMeter struct {
    minValueLabel string
    minValue      float64
    maxValueLabel string
    maxValue      float64
    value         float64
    format        string
    transparent   bool
    paddingX      int
    paddingY      int
}

func (m *HelixMeter) SetMin(minValue float64) {
    m.minValue = minValue
}

func (m *HelixMeter) SetMax(maxValue float64) {
    m.maxValue = maxValue
}

func (m *HelixMeter) addLabel(img *image.RGBA, x, y int, label string, color color.RGBA) {

    point := fixed.Point26_6{
        fixed.Int26_6((m.paddingX + x) * 64),
        fixed.Int26_6((m.paddingY + y) * 64),
    }

    d := &font.Drawer{
        Dst:  img,
        Src:  image.NewUniform(color),
        Face: basicfont.Face7x13,
        Dot:  point,
    }
    d.DrawString(label)
}

// drawHLine draws a horizontal line
func (m *HelixMeter) drawHLine(img *image.RGBA, x1, y, x2 int, color color.RGBA) {
    for ; x1 < x2; x1++ {
        img.Set(m.paddingX + x1, m.paddingY + y, color)
    }
}

// drawVLine draws a vertical line
func (m *HelixMeter) drawVLine(img *image.RGBA, x, y1, y2 int, color color.RGBA) {
    for y := y1; y < y2; y++ {
        img.Set(m.paddingX + x, m.paddingY + y, color)
    }
}

// Rect draws a rectangle utilizing drawHLine() and drawVLine()
func (m *HelixMeter) drawRect(img *image.RGBA, x, y, width, height int, color color.RGBA) {
    for xn := x; xn < x + width; xn++ {
        m.drawVLine(img, xn, y, y + height, color)
    }
}

func (m *HelixMeter) drawRule(img *image.RGBA) {

    var minLabel string = m.minValueLabel
    if minLabel == "" {
        minLabel = fmt.Sprintf(m.format, m.minValue)
    }

    var maxLabel string = m.minValueLabel
    if maxLabel == "" {
        maxLabel = fmt.Sprintf(m.format, m.maxValue)
    }

    const xo0 = 0
    const xo1 = 0 + 8*5 + 1*4
    const xo2 = 0 + (8*5 + 1*4) + 5 + (8*4 + 1*3) + 5

    m.drawVLine(img, xo1 - 1, 0, 13, BlackColor)
    m.drawVLine(img, xo2, 0, 13, BlackColor)

    m.addLabel(img, xo0 + 12, 3 + 8, minLabel, GrayColor)
    m.addLabel(img, xo2 + 5, 3 + 8, maxLabel, GrayColor)

}

func (m *HelixMeter) drawBars(img *image.RGBA) {

    // Blue bars (37 px)
    const xo0 = 0
    for i := 0; i < 5; i++ {
        m.drawRect(img, xo0 + (8+1) * i, 37 - (3 + i), 8, 3 + i, BlueColor)
    }

    // Green bars
    const xo1 = 0 + 9*5 + 4
    for i := 0; i < 4; i++ {
        m.drawRect(img, xo1 + (8+1) * i, 37 - (8 + i), 8, 8 + i, GreenColor)
    }

    // Red bars
    const xo2 = 0 + (9*5 + 4) + (9*4 + 4)
    for i := 0; i < 5; i++ {
        m.drawRect(img, xo2 + (8+1) * i, 37 - (12 + i), 8, 12 + i, RedColor)
    }

}

func (m *HelixMeter) drawMarker(img *image.RGBA) {

    var valueLavel string = fmt.Sprintf(m.format, m.value)
    var value = m.value

    var wide = m.maxValue - m.minValue
    var pos float64

    const xo0_0 = 0
    const xo0_1 = 8*5 + 1*4
    const xo1_0 = (8*5 + 1*4) + 5
    const xo1_1 = (8*5 + 1*4) + 5 + (8*4 + 1*3)
    const xo2_0 = (8*5 + 1*4) + 5 + (8*4 + 1*3) + 5
    const xo2_1 = (8*5 + 1*4) + 5 + (8*4 + 1*3) + 5 + (8*5 + 1*4)

    if value < m.minValue {
        // | ------- | ------- | ------- |
        //      ^ value
        log.Printf("Case 1")
        var minValue = m.minValue - wide
        if value < minValue {
            value = minValue
        }
        var width float64 = 8*5 + 1*4
        pos = width * (value) / wide
        pos += xo0_0
    }

    if m.minValue < value && value < m.maxValue {
        // | ------- | ------- | ------- |
        //                ^ value
        var width float64 = 4*8 + 3*1
        log.Printf("Case 2: wide = %f width = %f", wide, width)
        pos = width * (value - m.minValue) / wide
        pos += xo1_0
    }

    if m.maxValue < value {
        // | ------- | ------- | ------- |
        //                        ^ value
        log.Printf("Case 3")
        var maxValue = m.maxValue + wide
        if value > maxValue {
            value = maxValue
        }
        var width float64 = 8*5 + 1*4
        pos = width * (value - m.maxValue) / wide
        pos += xo2_0
    }

    log.Printf("marker x position is %.2f", pos)

    var x = int(pos)

    m.drawRect(img, x, 17, 2, 34, OrangeColor)
    m.addLabel(img, x + 2 + 5, 17 + 34, valueLavel, DarkColor)

}

func (m *HelixMeter) Render(path string) error {
    const sizeX = 133
    const sizeY = 54
    size := image.Rect(0, 0, sizeX + 2*m.paddingX, sizeY + 2*m.paddingY)
    img := image.NewRGBA(size)

    if !m.transparent {
        m.drawRect(img, 0, 0, sizeX, sizeY, WhiteColor)
    }

    m.drawRule(img)
    m.drawBars(img)
    m.drawMarker(img)

    f, err := os.Create(path)
    if err != nil {
        panic(err)
    }
    defer f.Close()
    return png.Encode(f, img)
}

func (m *HelixMeter) SetValue(value float64) {
    m.value = value
}

func (m *HelixMeter) SetFormat(format string) {
    m.format = format
}

func (m *HelixMeter) SetPadding(paddingX int, paddingY int) {
    m.paddingX = paddingX
    m.paddingY = paddingY
}

func NewHelixMeter() *HelixMeter {
    return &HelixMeter{
        format: "%.2f",
        transparent: false,
    }
}

var (
    RedColor    = color.RGBA{0xFF,0x00,0x0C,0xFF}
    BlueColor   = color.RGBA{0x00,0x52,0x9F,0xFF}
    GreenColor  = color.RGBA{0x6F,0xE0,0x5C,0xFF}
    OrangeColor = color.RGBA{0xFF,0x79,0x1E,0xFF}
    GrayColor   = color.RGBA{0x65,0x65,0x65,0xFF}
    DarkColor   = color.RGBA{0x3A,0x3A,0x3A,0xFF}
    BlackColor  = color.RGBA{0x00,0x00,0x00,0xFF}
    WhiteColor  = color.RGBA{0xFF,0xFF,0xFF,0xFF}
)

func main() {

    meter := NewHelixMeter()
    meter.SetMin(28)
    meter.SetMax(100)
    meter.SetValue(7)
    meter.SetFormat("%.0f")
    err := meter.Render("output.png")
    if err != nil {
        panic(err)
    }

    meter1 := NewHelixMeter()
    meter1.SetMin(50)
    meter1.SetMax(100)
    meter1.SetValue(117)
    meter1.SetFormat("%.0f")
    meter1.SetPadding(20, 10)
    err1 := meter1.Render("output1.png")
    if err1 != nil {
        panic(err1)
    }

}
