package nbt

func Parse(r io.Reader) (*Tag, error) {
	buf := bufio.NewBuffer(r)
	for {

	}
}

func ReadTagType(buf *bufio.Reader) (TagType, error) {
	b, err := buf.ReadByte()
	if err != nil {
		return TagEnd, err
	}
	return TagType(b), nil
}

func ReadName(buf *bufio.Reader) (string, error) {
	len := binary.BigEndian.ReadUint16(buf)
	
	nameBytes := make([]byte, len)
	_, err := buf.Read(nameBytes)
	if err != nil {
		return "", err
	}

	return string(nameBytes), nil
}

func 
