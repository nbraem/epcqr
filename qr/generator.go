package qr

import (
    "fmt"
    go_qr "github.com/piglig/go-qr"
)

// Generate creates a QR code with the given content and saves it to outputPath.
// format must be "png" or "svg".
func Generate(content string, format string, outputPath string) error {
    // Encode text with Medium error correction
    code, err := go_qr.EncodeText(content, go_qr.Medium)
    if err != nil {
        return err
    }

    // 10x scale, 4 blocks border
    config := go_qr.NewQrCodeImgConfig(10, 4)

    switch format {
    case "png":
        return code.PNG(config, outputPath)
    case "svg":
        // Valid SVG color strings: hex codes
        return code.SVG(config, outputPath, "#000000", "#FFFFFF")
    case "console":
        printQRCodeToConsole(code)
        return nil
    default:
        return fmt.Errorf("unsupported format: %s", format)
    }
}

func printQRCodeToConsole(code *go_qr.QrCode) {
    size := code.GetSize()
    border := 4 // Minimum quiet zone

    // ANSI colors
    reset := "\033[0m"
    white := "\033[47m  " + reset
    black := "\033[40m  " + reset

    // Print top border
    for i := 0; i < border; i++ {
        for j := 0; j < size+2*border; j++ {
            fmt.Print(white)
        }
        fmt.Println()
    }

    // Print content
    for y := 0; y < size; y++ {
        // Left border
        for i := 0; i < border; i++ {
            fmt.Print(white)
        }

        for x := 0; x < size; x++ {
            if code.GetModule(x, y) {
                fmt.Print(black)
            } else {
                fmt.Print(white)
            }
        }

        // Right border
        for i := 0; i < border; i++ {
            fmt.Print(white)
        }
        fmt.Println()
    }

    // Print bottom border
    for i := 0; i < border; i++ {
        for j := 0; j < size+2*border; j++ {
            fmt.Print(white)
        }
        fmt.Println()
    }
}
