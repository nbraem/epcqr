package epc

import (
    "fmt"
    "strings"
)

// Data holds the information required for a SEPA Credit Transfer QR Code (EPC069-12).
type Data struct {
    Name           string
    IBAN           string
    BIC            string
    Amount         float64
    Reference      string // Creditor Reference (ISO 11649)
    RemittanceText string // Unstructured Remittance Information
    Information    string // Beneficiary to Originator Information
}

// GenerateString constructs the content string for the QR code according to EPC069-12 standards.
// Format:
// BCD
// 002
// 1
// SCT
// BIC (optional)
// Name
// IBAN
// Amount (EUR+Amount)
// Purpose (optional)
// Reference (optional)
// RemittanceText (optional)
// Information (optional)
func (d *Data) GenerateString() (string, error) {
    if d.IBAN == "" {
        return "", fmt.Errorf("IBAN is required")
    }
    if d.Name == "" {
        return "", fmt.Errorf("beneficiary name is required")
    }

    var sb strings.Builder

    // 1. Service Tag
    sb.WriteString("BCD\n")
    // 2. Version
    sb.WriteString("002\n")
    // 3. Character Set (1 = UTF-8)
    sb.WriteString("1\n")
    // 4. Identification (SCT = SEPA Credit Transfer)
    sb.WriteString("SCT\n")
    // 5. BIC
    sb.WriteString(strings.TrimSpace(d.BIC) + "\n")
    // 6. Beneficiary Name (max 70 chars)
    name := strings.TrimSpace(d.Name)
    if len(name) > 70 {
        name = name[:70]
    }
    sb.WriteString(name + "\n")
    // 7. IBAN
    sb.WriteString(strings.ReplaceAll(d.IBAN, " ", "") + "\n")
    // 8. Amount (EUR<Amount>)
    sb.WriteString(fmt.Sprintf("EUR%.2f\n", d.Amount))
    // 9. Purpose Code (optional, empty here for simplicity)
    sb.WriteString("\n")
    // 10. Remittance Reference (Structured)
    // Must be ISO 11649 (RF...) if present.
    sb.WriteString(strings.TrimSpace(d.Reference) + "\n")
    // 11. Remittance Text (Unstructured)
    sb.WriteString(strings.TrimSpace(d.RemittanceText) + "\n")
    // 12. Information
    sb.WriteString(strings.TrimSpace(d.Information) + "\n") // End of data, standard says it doesn't need trailing newline necessarily but usually handled by splitter

    return sb.String(), nil
}
