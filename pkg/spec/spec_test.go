package spec

import "testing"

func TestFieldType(t *testing.T) {
	tests := []struct {
		name       string
		fieldType  FieldType
		wantString string
		wantDigits int
		wantVar    bool
	}{
		{"Fixed", FieldTypeFixed, "Fixed", 0, false},
		{"L", FieldTypeL, "L", 1, true},
		{"LL", FieldTypeLL, "LL", 2, true},
		{"LLL", FieldTypeLLL, "LLL", 3, true},
		{"Bitmap", FieldTypeBitmap, "Bitmap", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.fieldType.String(); got != tt.wantString {
				t.Errorf("String() = %v, want %v", got, tt.wantString)
			}

			if got := tt.fieldType.LengthIndicatorDigits(); got != tt.wantDigits {
				t.Errorf("LengthIndicatorDigits() = %v, want %v", got, tt.wantDigits)
			}

			if got := tt.fieldType.IsVariable(); got != tt.wantVar {
				t.Errorf("IsVariable() = %v, want %v", got, tt.wantVar)
			}
		})
	}
}

func TestDataType(t *testing.T) {
	tests := []struct {
		name     string
		dataType DataType
		want     string
	}{
		{"Numeric", DataTypeNumeric, "Numeric"},
		{"Alpha", DataTypeAlpha, "Alpha"},
		{"Alphanumeric", DataTypeAlphanumeric, "Alphanumeric"},
		{"AlphaNumericSpecial", DataTypeAlphaNumericSpecial, "AlphaNumericSpecial"},
		{"Binary", DataTypeBinary, "Binary"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.dataType.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestEncodingType(t *testing.T) {
	tests := []struct {
		name     string
		encoding EncodingType
		want     string
	}{
		{"ASCII", EncodingASCII, "ASCII"},
		{"EBCDIC", EncodingEBCDIC, "EBCDIC"},
		{"BCD", EncodingBCD, "BCD"},
		{"Binary", EncodingBinary, "Binary"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.encoding.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestPaddingType(t *testing.T) {
	tests := []struct {
		name    string
		padding PaddingType
		want    string
	}{
		{"None", PaddingNone, "None"},
		{"Left", PaddingLeft, "Left"},
		{"Right", PaddingRight, "Right"},
		{"Center", PaddingCenter, "Center"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.padding.String(); got != tt.want {
				t.Errorf("String() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestFieldSpec(t *testing.T) {
	spec := &FieldSpec{
		Number:      2,
		Name:        "PrimaryAccountNumber",
		Aliases:     []string{"PAN", "CardNumber"},
		Type:        FieldTypeLL,
		MaxLength:   19,
		DataType:    DataTypeNumeric,
		Encoding:    EncodingASCII,
		Description: "Primary Account Number",
	}

	if spec.Number != 2 {
		t.Errorf("Number = %v, want 2", spec.Number)
	}

	if spec.Name != "PrimaryAccountNumber" {
		t.Errorf("Name = %v, want PrimaryAccountNumber", spec.Name)
	}

	if len(spec.Aliases) != 2 {
		t.Errorf("len(Aliases) = %v, want 2", len(spec.Aliases))
	}

	if spec.Type != FieldTypeLL {
		t.Errorf("Type = %v, want FieldTypeLL", spec.Type)
	}
}

func TestSpecDefaults(t *testing.T) {
	defaults := FieldDefaults{
		Encoding: EncodingASCII,
		Padding:  PaddingLeft,
		PadChar:  '0',
	}

	if defaults.Encoding != EncodingASCII {
		t.Errorf("Encoding = %v, want ASCII", defaults.Encoding)
	}

	if defaults.Padding != PaddingLeft {
		t.Errorf("Padding = %v, want Left", defaults.Padding)
	}

	if defaults.PadChar != '0' {
		t.Errorf("PadChar = %v, want '0'", defaults.PadChar)
	}
}

func TestSpec(t *testing.T) {
	spec := &Spec{
		Name:    "ISO 8583 v1987",
		Version: "1.0",
		Defaults: FieldDefaults{
			Encoding: EncodingASCII,
			Padding:  PaddingLeft,
			PadChar:  '0',
		},
		Fields: map[int]*FieldSpec{
			0: {
				Number:   0,
				Name:     "MessageTypeIndicator",
				Aliases:  []string{"MTI"},
				Type:     FieldTypeFixed,
				Length:   4,
				DataType: DataTypeNumeric,
			},
			2: {
				Number:    2,
				Name:      "PrimaryAccountNumber",
				Aliases:   []string{"PAN"},
				Type:      FieldTypeLL,
				MaxLength: 19,
				DataType:  DataTypeNumeric,
			},
		},
	}

	if spec.Name != "ISO 8583 v1987" {
		t.Errorf("Name = %v, want 'ISO 8583 v1987'", spec.Name)
	}

	if len(spec.Fields) != 2 {
		t.Errorf("len(Fields) = %v, want 2", len(spec.Fields))
	}

	if spec.Fields[0].Name != "MessageTypeIndicator" {
		t.Errorf("Fields[0].Name = %v, want MessageTypeIndicator", spec.Fields[0].Name)
	}
}
