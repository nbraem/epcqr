package main

import (
    "epcqr/epc"
    "epcqr/qr"
    "flag"
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strconv"
    "strings"

    ibanpkg "github.com/jbub/banking/iban"
)

func main() {
    // Common flags
    iban := flag.String("iban", "", "IBAN (required)")
    name := flag.String("name", "", "Beneficiary Name (required)")
    amount := flag.Float64("amount", 0.0, "Amount in EUR")
    bic := flag.String("bic", "", "BIC (optional)")
    ref := flag.String("ref", "", "Remittance Reference (ISO 11649)")
    text := flag.String("text", "", "Remittance Text")
    info := flag.String("info", "", "Beneficiary to Originator Information")
    force := flag.Bool("force", false, "Force generation even if IBAN is invalid")

    // CLI specific
    out := flag.String("out", "qr.png", "Output file path")
    format := flag.String("format", "console", "Output format (png, svg, console)")

    // Mode
    mode := flag.String("mode", "cli", "Mode of operation: 'cli' or 'server'")
    port := flag.String("port", "8080", "Port to listen on (server mode)")

    flag.Parse()

    if *mode == "server" {
        startServer(*port)
    } else {
        if *iban == "" || *name == "" {
            fmt.Println("Error: -iban and -name are required for CLI mode")
            flag.PrintDefaults()
            os.Exit(1)
        }

        *iban = strings.ReplaceAll(*iban, " ", "")

        if !*force {
            _, err := ibanpkg.Parse(*iban)
            if err != nil {
                log.Fatalf("Invalid IBAN: %v", err)
            }
        }

        data := &epc.Data{
            Name:           *name,
            IBAN:           *iban,
            BIC:            *bic,
            Amount:         *amount,
            Reference:      *ref,
            RemittanceText: *text,
            Information:    *info,
        }

        content, err := data.GenerateString()
        if err != nil {
            log.Fatalf("Error generating EPC data: %v", err)
        }

        err = qr.Generate(content, *format, *out)
        if err != nil {
            log.Fatalf("Error generating QR code: %v", err)
        }

        if *format != "console" {
            fmt.Printf("QR code saved to %s\n", *out)
        }
    }
}

func startServer(port string) {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        query := r.URL.Query()

        iban := query.Get("iban")
        name := query.Get("name")
        amountStr := query.Get("amount")
        bic := query.Get("bic")
        ref := query.Get("ref")
        text := query.Get("text")
        info := query.Get("info")
        format := query.Get("format")
		force := query.Get("force") == "1" || strings.ToLower(query.Get("force")) == "true"

        if format == "" {
            format = "png"
        }
        if format != "png" && format != "svg" {
            http.Error(w, "Invalid format. Supported: png, svg", http.StatusBadRequest)
            return
        }

        if iban == "" || name == "" {
            http.Error(w, "Missing required parameters: iban, name", http.StatusBadRequest)
            return
        }

        iban = strings.ReplaceAll(iban, " ", "")

        if !force {
            _, err := ibanpkg.Parse(iban)
            if err != nil {
                http.Error(w, fmt.Sprintf("Invalid IBAN: %v", err), http.StatusBadRequest)
                return
            }
        }

        amount := 0.0
        if amountStr != "" {
            var err error
            amount, err = strconv.ParseFloat(amountStr, 64)
            if err != nil {
                http.Error(w, "Invalid amount", http.StatusBadRequest)
                return
            }
        }

        data := &epc.Data{
            Name:           name,
            IBAN:           iban,
            BIC:            bic,
            Amount:         amount,
            Reference:      ref,
            RemittanceText: text,
            Information:    info,
        }

        content, err := data.GenerateString()
        if err != nil {
            http.Error(w, fmt.Sprintf("Error building EPC data: %v", err), http.StatusBadRequest)
            return
        }

        // Create temp file
        tmpFile, err := ioutil.TempFile("", "epcqr-*."+format)
        if err != nil {
            http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            log.Printf("Temp file error: %v", err)
            return
        }
        defer os.Remove(tmpFile.Name()) // Clean up
        tmpFile.Close()                 // Close so generator can write to it (by path)

        err = qr.Generate(content, format, tmpFile.Name())
        if err != nil {
            http.Error(w, "Failed to generate QR", http.StatusInternalServerError)
            log.Printf("Generation error: %v", err)
            return
        }

        // Read back
        imgData, err := ioutil.ReadFile(tmpFile.Name())
        if err != nil {
            http.Error(w, "Failed to read generated image", http.StatusInternalServerError)
            return
        }

        // Serve
        if format == "png" {
            w.Header().Set("Content-Type", "image/png")
        } else {
            w.Header().Set("Content-Type", "image/svg+xml")
        }
        w.Write(imgData)
    })

    fmt.Printf("Server listening on port %s...\n", port)
    // Example usage hint
    fmt.Printf("Try: http://localhost:%s/?iban=BE123456789&name=JohnDoe&amount=10.50\n", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}
