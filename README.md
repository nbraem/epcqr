# epcqr

epcqr is a Go-based tool for generating SEPA Credit Transfer (EPC069-12) QR codes, commonly known as "GiroCode" or "EPC QR Code". It supports generating QR codes in the terminal (console), as well as exporting them to PNG or SVG files. It can run as a CLI tool or as a web server.

## Features

- **CLI Support**: Generate QR codes directly from the command line.
- **Multiple Formats**: Output to Console (ANSI), PNG, or SVG.
- **Web Server Mode**: Run a lightweight HTTP server to generate QR codes on the fly via URL parameters.
- **Standard Compliant**: meaningful implementation of the EPC069-12 standard.

## Installation

Ensure you have Go installed (1.13+).

```bash
git clone https://github.com/yourusername/epcqr.git
cd epcqr
go build
```

This will create an `epcqr` binary in your directory.

## Usage

### CLI Mode

The default mode is CLI. You must provide at least an IBAN and a Beneficiary Name. The default output format is `console`.

**Basic Usage (Console Output):**

```bash
./epcqr -iban "BE12345678901234" -name "John Doe" -amount 50.00
```

**Save to PNG:**

```bash
./epcqr -iban "BE12345678901234" -name "John Doe" -amount 50.00 -format png -out my_qr.png
```

**Save to SVG:**

```bash
./epcqr -iban "BE12345678901234" -name "John Doe" -amount 50.00 -format svg -out my_qr.svg
```

**All Available Flags:**

| Flag      | Description                             | Required? | Default   |
| --------- | --------------------------------------- | --------- | --------- |
| `-iban`   | Beneficiary IBAN                        | Yes (CLI) |           |
| `-name`   | Beneficiary Name                        | Yes (CLI) |           |
| `-amount` | Amount in EUR                           | No        | 0.0       |
| `-bic`    | Beneficiary BIC                         | No        |           |
| `-ref`    | Remittance Reference (ISO 11649)        | No        |           |
| `-text`   | Remittance Text                         | No        |           |
| `-info`   | Beneficiary to Originator Info          | No        |           |
| `-format` | Output format (`console`, `png`, `svg`) | No        | `console` |
| `-out`    | Output filename                         | No        | `qr.png`  |
| `-mode`   | Operation mode (`cli`, `server`)        | No        | `cli`     |
| `-port`   | Server port                             | No        | `8080`    |

### Server Mode

You can run `epcqr` as a web server to generate QR codes dynamically.

```bash
./epcqr -mode server -port 8080
```

Once running, access the generator via your browser or HTTP client:

```
http://localhost:8080/?iban=BE123456789&name=JohnDoe&amount=10.50&format=png
```

**Query Parameters:**

- `iban` (required)
- `name` (required)
- `amount`
- `bic`
- `ref`
- `text`
- `info`
- `format` (`png` or `svg`)

## Dependencies

- [github.com/piglig/go-qr](https://github.com/piglig/go-qr) for QR code generation.
