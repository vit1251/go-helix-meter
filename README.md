# Helix Meter

Source code example:

```golang
meter := NewHelixMeter()
meter.SetMin(28)
meter.SetMax(100)
meter.SetValue(7)
meter.SetFormat("%.0f")
err := meter.Render("output.png")
if err != nil {
    panic(err)
}
```

## Origianl Helix meter

![Helix Original Meter](/docs/Original.png)

## Library output meter

![Library example output](/docs/output.png)

![Library example output 1](/docs/output1.png)
