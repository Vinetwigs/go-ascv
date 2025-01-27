package ascv

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
)

const Magic = "AAVF"

// Header rappresenta l'header del file .ascv
type Header struct {
	Magic       [4]byte
	Version     uint8
	Width       uint16
	Height      uint16
	FPS         uint8
	Frames      uint32
	Compression uint8
	Charset     uint8
	Reserved    [16]byte
}

// Frame rappresenta un frame di ASCII art
type Frame struct {
	Size    int
	Content []byte
}

// EncodeRLE comprime i dati utilizzando la codifica Run-Length Encoding (RLE)
func EncodeRLE(data []byte) []byte {
	var buffer bytes.Buffer
	n := len(data)
	for i := 0; i < n; {
		count := 1
		for i+count < n && data[i] == data[i+count] && count < 255 {
			count++
		}
		buffer.WriteByte(byte(count))
		buffer.WriteByte(data[i])
		i += count
	}
	return buffer.Bytes()
}

// DecodeRLE decomprime i dati utilizzando la codifica Run-Length Encoding (RLE)
func DecodeRLE(data []byte) ([]byte, error) {
	var buffer bytes.Buffer
	n := len(data)
	if n%2 != 0 {
		return nil, errors.New("invalid RLE data")
	}
	for i := 0; i < n; i += 2 {
		count := int(data[i])
		char := data[i+1]
		for j := 0; j < count; j++ {
			buffer.WriteByte(char)
		}
	}
	return buffer.Bytes(), nil
}

// WriteASCV scrive una sequenza di frame in un file .ascv
func WriteASCV(filename string, header Header, frames []Frame) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Scrittura dell'header
	err = binary.Write(file, binary.LittleEndian, header)
	if err != nil {
		return err
	}

	// Scrittura dei frame
	for _, frame := range frames {
		// Scrittura della dimensione del frame (in formato VLQ)
		vlq := encodeVLQ(uint32(frame.Size))
		_, err = file.Write(vlq)
		if err != nil {
			return err
		}
		// Scrittura del contenuto del frame
		_, err = file.Write(frame.Content)
		if err != nil {
			return err
		}
	}

	return nil
}

// ReadASCV legge un file .ascv e restituisce l'header e i frame
func ReadASCV(filename string) (Header, []Frame, error) {
	var header Header

	file, err := os.Open(filename)
	if err != nil {
		return header, nil, err
	}
	defer file.Close()

	// Lettura dell'header
	err = binary.Read(file, binary.LittleEndian, &header)
	if err != nil {
		return header, nil, err
	}

	if string(header.Magic[:]) != Magic {
		return header, nil, errors.New("invalid file format")
	}

	var frames []Frame
	for {
		// Lettura della dimensione del frame (in formato VLQ)
		size, err := decodeVLQ(file)
		if err != nil {
			break // Fine del file
		}

		// Lettura del contenuto del frame
		content := make([]byte, size)
		_, err = file.Read(content)
		if err != nil {
			return header, nil, err
		}
		frames = append(frames, Frame{Size: int(size), Content: content})
	}

	return header, frames, nil
}

// encodeVLQ codifica un uint32 nel formato VLQ (Variable-Length Quantity)
func encodeVLQ(value uint32) []byte {
	var buffer []byte
	for {
		byteValue := value & 0x7F
		value >>= 7
		if value > 0 {
			buffer = append(buffer, byte(byteValue|0x80))
		} else {
			buffer = append(buffer, byte(byteValue))
			break
		}
	}
	return buffer
}

// decodeVLQ decodifica un valore VLQ (Variable-Length Quantity) da un file
func decodeVLQ(file *os.File) (uint32, error) {
	var value uint32
	var shift uint32
	for {
		byteValue := make([]byte, 1)
		_, err := file.Read(byteValue)
		if err != nil {
			return 0, err
		}
		value |= uint32(byteValue[0]&0x7F) << shift
		if byteValue[0]&0x80 == 0 {
			break
		}
		shift += 7
	}
	return value, nil
}
