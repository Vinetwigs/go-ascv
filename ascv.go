package ascv

import (
	"bytes"
	"encoding/binary"
	"errors"
	"os"
)

const Magic = "AAVF"

// Header represents the header of the .ascv file
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

// Frame represents a frame of ASCII art
type Frame struct {
	Size    int
	Content []byte
}

// EncodeRLE compresses data using Run-Length Encoding (RLE) encoding
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

// DecodeRLE decompresses data using Run-Length Encoding (RLE)
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

// WriteASCV writes a sequence of frames to an .ascv file
func WriteASCV(filename string, header Header, frames []Frame) error {
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Writing the header
	err = binary.Write(file, binary.LittleEndian, header)
	if err != nil {
		return err
	}

	// Writing frames
	for _, frame := range frames {
		// Writing frame size (in VLQ format)
		vlq := encodeVLQ(uint32(frame.Size))
		_, err = file.Write(vlq)
		if err != nil {
			return err
		}
		// Writing frame contents
		_, err = file.Write(frame.Content)
		if err != nil {
			return err
		}
	}

	return nil
}

// ReadASCV reads an .ascv file and returns the header and frames
func ReadASCV(filename string) (Header, []Frame, error) {
	var header Header

	file, err := os.Open(filename)
	if err != nil {
		return header, nil, err
	}
	defer file.Close()

	// Header reading
	err = binary.Read(file, binary.LittleEndian, &header)
	if err != nil {
		return header, nil, err
	}

	if string(header.Magic[:]) != Magic {
		return header, nil, errors.New("invalid file format")
	}

	var frames []Frame
	for {
		// Reading the frame size (in VLQ format)
		size, err := decodeVLQ(file)
		if err != nil {
			break // EOF
		}

		// Reading frame content
		content := make([]byte, size)
		_, err = file.Read(content)
		if err != nil {
			return header, nil, err
		}
		frames = append(frames, Frame{Size: int(size), Content: content})
	}

	return header, frames, nil
}

// encodeVLQ encodes a uint32 in VLQ (Variable-Length Quantity) format
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

// decodeVLQ decodes a Variable-Length Quantity (VLQ) value from a file
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
