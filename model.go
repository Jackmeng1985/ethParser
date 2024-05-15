package ethParser

import (
	"strconv"
	"strings"
)

// Transaction represents an Ethereum transaction
type Transaction struct {
	Hash     string `json:"hash"`
	From     string `json:"from"`
	To       string `json:"to"`
	Value    string `json:"value"`
	Gas      string `json:"gas"`
	GasPrice string `json:"gasPrice"`
}

type Transactions []*Transaction

// Block represents Ethereum block data
type Block struct {
	Number       HexNumber `json:"number"`
	Hash         string    `json:"hash"`
	ParentHash   string    `json:"parentHash"`
	Timestamp    HexNumber `json:"timestamp"`
	Transactions Transactions
}

type HexNumber uint64

// MarshalJSON implements json.Marshaler.
func (b HexNumber) MarshalJSON() ([]byte, error) {
	buf := make([]byte, 3, 12)
	copy(buf, `"0x`)
	buf = strconv.AppendUint(buf, uint64(b), 16)
	buf = append(buf, '"')
	return buf, nil
}

// UnmarshalJSON implements json.Unmarshaler
func (b *HexNumber) UnmarshalJSON(input []byte) error {
	// Remove quotes and "0x" prefix from the JSON string
	str := strings.Trim(string(input), "\"")
	if strings.HasPrefix(str, "0x") {
		str = str[2:]
	}
	// Parse the hexadecimal string into uint64 value
	val, err := strconv.ParseUint(str, 16, 64)
	if err != nil {
		return err
	}
	// Assign the parsed value to the HexUint64 pointer
	*b = HexNumber(val)
	return nil
}
