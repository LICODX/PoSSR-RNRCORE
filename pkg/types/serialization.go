package types

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
)

// SerializeTransaction creates a canonical byte representation for signing
func SerializeTransaction(tx Transaction) []byte {
	var buf bytes.Buffer
	buf.Write(tx.Sender[:])
	buf.Write(tx.Receiver[:])
	binary.Write(&buf, binary.LittleEndian, tx.Amount)
	binary.Write(&buf, binary.LittleEndian, tx.Nonce)
	buf.Write(tx.Payload)
	return buf.Bytes()
}

// HashTransaction calculates transaction hash
func HashTransaction(tx Transaction) [32]byte {
	return sha256.Sum256(SerializeTransaction(tx))
}

// SerializeBlockHeader creates canonical header bytes
func SerializeBlockHeader(h BlockHeader) []byte {
	var buf bytes.Buffer
	binary.Write(&buf, binary.LittleEndian, h.Version)
	buf.Write(h.PrevBlockHash[:])
	buf.Write(h.MerkleRoot[:])
	binary.Write(&buf, binary.LittleEndian, h.Timestamp)
	binary.Write(&buf, binary.LittleEndian, h.Height)
	for _, node := range h.WinningNodes {
		buf.Write(node[:])
	}
	buf.Write(h.VRFSeed[:])
	return buf.Bytes()
}

// HashBlockHeader calculates block header hash
func HashBlockHeader(h BlockHeader) [32]byte {
	return sha256.Sum256(SerializeBlockHeader(h))
}
