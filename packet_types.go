package packet

type ( 
	Boolean bool

	Byte int8
	UnsignedByte uint8

	Short int16
	UnsignedShort uint16

	Int int32
	UnsignedInt uint32

	Float float32
	Double float64

	String string

	Chat string
	Identifier string
	VarInt int32
	VarLong int64
	
	Position struct {
		x,y,z int32
	}

	Angle byte
	UUID [16]byte

}