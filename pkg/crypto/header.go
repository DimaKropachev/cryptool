package crypto

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
)

const MagicNum = "CRPT"

type Header struct {
	MagicNum  string
	AlgID     uint16
	BlockSize uint64
	SaltSize  uint32
	Salt      []byte
	NonceSize uint32
}

func NewHeader(algID, blockSize, saltSize, nonceSize int, salt []byte) *Header {
	header := &Header{
		MagicNum:  MagicNum,
		AlgID:     uint16(algID),
		BlockSize: uint64(blockSize),
		SaltSize:  uint32(saltSize),
		Salt:      salt,
		NonceSize: uint32(nonceSize),
	}

	return header
}

func DecryptHeader(r io.Reader) (*Header, error) {
	header := Header{}

	// Decrypt MagicNum
	magicNum := make([]byte, 4)
	_, err := r.Read(magicNum)
	if err != nil {
		return nil, err
	}
	if MagicNum != string(magicNum) {
		return nil, fmt.Errorf("")
	}
	header.MagicNum = string(magicNum)

	// Decrypt algID
	if err := binary.Read(r, binary.LittleEndian, &header.AlgID); err != nil {
		return nil, err
	}

	// Decrypt BlockSize
	if err := binary.Read(r, binary.LittleEndian, &header.BlockSize); err != nil {
		return nil, err
	}

	// Drcrypt SaltSize
	if err := binary.Read(r, binary.LittleEndian, &header.SaltSize); err != nil {
		return nil, err
	}

	// Decrypt Salt
	salt := make([]byte, header.SaltSize)
	if _, err := r.Read(salt); err != nil {
		return nil, err
	}
	header.Salt = salt

	// Decrypt NonceSize
	if err := binary.Read(r, binary.LittleEndian, &header.NonceSize); err != nil {
		return nil, err
	}

	return &header, nil
}

func EncryptHeader(header *Header) ([]byte, error) {
	result := bytes.NewBuffer([]byte{})

	if _, err := result.Write([]byte(MagicNum)); err != nil {
		return nil, err
	}

	if err := binary.Write(result, binary.LittleEndian, header.AlgID); err != nil {
		return nil, err
	}

	if err := binary.Write(result, binary.LittleEndian, header.BlockSize); err != nil {
		return nil, err
	}

	if err := binary.Write(result, binary.LittleEndian, header.SaltSize); err != nil {
		return nil, err
	}

	if _, err := result.Write(header.Salt); err != nil {
		return nil, err
	}

	if err := binary.Write(result, binary.LittleEndian, header.NonceSize); err != nil {
		return nil, err
	}

	return result.Bytes(), nil
}
