package money

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
)

func TestCAD_ParseCAD(t *testing.T) {
	tests := []struct {
		Input  string
		Output CAD
		Error  string
	}{
		{
			Input: "-$1234.56",
			Output: CAD{
				cents: -123456,
			},
			Error: "",
		},
		{
			Input: "$-1234.56",
			Output: CAD{
				cents: -123456,
			},
			Error: "",
		},
		{
			Input: "-$1,234.56",
			Output: CAD{
				cents: -123456,
			},
			Error: "",
		},
		{
			Input: "$-1,234.56",
			Output: CAD{
				cents: -123456,
			},
			Error: "",
		},
		{
			Input: "CAD -$1234.56",
			Output: CAD{
				cents: -123456,
			},
			Error: "",
		},
		{
			Input: "CAD $-1234.56",
			Output: CAD{
				cents: -123456,
			},
			Error: "",
		},
		{
			Input: "CAD S $-1234.56",
			Output: CAD{
				cents: 0,
			},
			Error: "strconv.ParseInt: parsing \"S-123456\": invalid syntax",
		},
		{
			Input: "CAD  $-1234.56A",
			Output: CAD{
				cents: 0,
			},
			Error: "strconv.ParseInt: parsing \"-123456A\": invalid syntax",
		},
	}

	for testNumber, test := range tests {
		cad, err := ParseCAD(test.Input)
		if cad.cents != test.Output.cents {
			t.Errorf("Test %d :  cents is %d  but was expecting %d", testNumber, cad.cents, test.Output.cents)
		}
		if test.Error != "" && ((err == nil) || (err != nil && err.Error() != test.Error)) {
			t.Errorf("Test %d :  error is %s  but was expecting %s", testNumber, err, test.Error)
		}
	}
}

func TestCAD_Abs(t *testing.T) {
	tests := []struct {
		Input    CAD
		Expected CAD
	}{
		{
			Input: CAD{
				cents: 78954,
			},
			Expected: CAD{
				cents: 78954,
			},
		},
		{
			Input: CAD{
				cents: -78954,
			},
			Expected: CAD{
				cents: 78954,
			},
		},
		{
			Input: CAD{
				cents: 0,
			},
			Expected: CAD{
				cents: 0,
			},
		},
	}

	for testNumber, test := range tests {
		cad := test.Input.Abs()
		if cad != test.Expected {
			t.Errorf("Test %d :  output is %d  but was expecting %d", testNumber, cad, test.Expected)
		}
	}
}

func TestCAD_Add(t *testing.T) {
	tests := []struct {
		Input    CAD
		Addition CAD
		Expected CAD
	}{
		{
			Input: CAD{
				cents: 1,
			},
			Addition: CAD{
				cents: 2,
			},
			Expected: CAD{
				cents: 3,
			},
		},
		{
			Input: CAD{
				cents: -7,
			},
			Addition: CAD{
				cents: 2,
			},
			Expected: CAD{
				cents: -5,
			},
		},
		{
			Input: CAD{
				cents: 0,
			},
			Addition: CAD{
				cents: 2,
			},
			Expected: CAD{
				cents: 2,
			},
		},
	}

	for testNumber, test := range tests {
		cad := test.Input.Add(test.Addition)
		if cad != test.Expected {
			t.Errorf("Test %d :  output is %d  but was expecting %d", testNumber, cad, test.Expected)
		}
	}
}

func TestCAD_AsCents(t *testing.T) {
	tests := []struct {
		Input    CAD
		Expected int64
	}{
		{
			Input: CAD{
				cents: 1,
			},
			Expected: 1,
		},
		{
			Input: CAD{
				cents: -7,
			},
			Expected: -7,
		},
		{
			Input: CAD{
				cents: 0,
			},
			Expected: 0,
		},
	}

	for testNumber, test := range tests {
		cents := test.Input.AsCents()
		if cents != test.Expected {
			t.Errorf("Test %d :  output is %d  but was expecting %d", testNumber, cents, test.Expected)
		}
	}
}

func TestCAD_CanonicalForm(t *testing.T) {
	type Canonical struct {
		Dollars int64
		Cents   int64
	}
	tests := []struct {
		Input    CAD
		Expected Canonical
	}{
		{
			Input: CAD{
				cents: 112,
			},

			Expected: Canonical{
				Dollars: 1,
				Cents:   12,
			},
		},
		{
			Input: CAD{
				cents: -7008,
			},
			Expected: Canonical{
				Dollars: -70,
				Cents:   -8,
			},
		},
		{
			Input: CAD{
				cents: 0,
			},
			Expected: Canonical{
				Dollars: 0,
				Cents:   0,
			},
		},
	}

	for testNumber, test := range tests {
		dollars, cents := test.Input.CanonicalForm()
		if cents != test.Expected.Cents || dollars != test.Expected.Dollars {
			t.Errorf("Test %d :  output is %d$ and %d cents  but was expecting %d$ and %d cents", testNumber, dollars, cents, test.Expected.Dollars, test.Expected.Cents)
		}
	}
}

func TestCAD_Mul(t *testing.T) {
	tests := []struct {
		Input    CAD
		Scalar   int64
		Expected CAD
	}{
		{
			Input: CAD{
				cents: 1,
			},
			Scalar: 2,
			Expected: CAD{
				cents: 2,
			},
		},
		{
			Input: CAD{
				cents: -7,
			},
			Scalar: 2,
			Expected: CAD{
				cents: -14,
			},
		},
		{
			Input: CAD{
				cents: 0,
			},
			Scalar: 2,
			Expected: CAD{
				cents: 0,
			},
		},
	}

	for testNumber, test := range tests {
		cad := test.Input.Mul(test.Scalar)
		if cad != test.Expected {
			t.Errorf("Test %d :  output is %d  but was expecting %d", testNumber, cad, test.Expected)
		}
	}
}

func TestCAD_Sub(t *testing.T) {
	tests := []struct {
		Input    CAD
		Sub      CAD
		Expected CAD
	}{
		{
			Input: CAD{
				cents: 5,
			},
			Sub: CAD{
				cents: 2,
			},
			Expected: CAD{
				cents: 3,
			},
		},
		{
			Input: CAD{
				cents: -7,
			},
			Sub: CAD{
				cents: 2,
			},
			Expected: CAD{
				cents: -9,
			},
		},
		{
			Input: CAD{
				cents: 0,
			},
			Sub: CAD{
				cents: 2,
			},
			Expected: CAD{
				cents: -2,
			},
		},
	}

	for testNumber, test := range tests {
		cad := test.Input.Sub(test.Sub)
		if cad != test.Expected {
			t.Errorf("Test %d :  output is %d  but was expecting %d", testNumber, cad, test.Expected)
		}
	}
}

func TestCAD_GoString(t *testing.T) {
	tests := []struct {
		Input    CAD
		Expected string
	}{
		{
			Input: CAD{
				cents: 1,
			},
			Expected: "main.Cents(1)",
		},
		{
			Input: CAD{
				cents: -7,
			},
			Expected: "main.Cents(-7)",
		},
		{
			Input: CAD{
				cents: 0,
			},
			Expected: "main.Cents(0)",
		},
		{
			Input: CAD{
				cents: 1785,
			},
			Expected: "main.Cents(1785)",
		},
	}

	for testNumber, test := range tests {
		var b bytes.Buffer
		fmt.Fprintf(&b, "%#v", test.Input)
		out := b.String()
		if out != test.Expected {
			t.Errorf("Test %d :  output is %s  but was expecting %s", testNumber, out, test.Expected)
		}
	}
}

func TestCAD_Stringer(t *testing.T) {
	tests := []struct {
		Input    CAD
		Expected string
	}{
		{
			Input: CAD{
				cents: 1,
			},
			Expected: "$0.01\n",
		},
		{
			Input: CAD{
				cents: -7,
			},
			Expected: "-$0.07\n",
		},
		{
			Input: CAD{
				cents: 0,
			},
			Expected: "$0.00\n",
		},
		{
			Input: CAD{
				cents: 1785,
			},
			Expected: "$17.85\n",
		},
	}

	for testNumber, test := range tests {
		var b bytes.Buffer
		fmt.Fprintln(&b, test.Input)
		out := b.String()
		if out != test.Expected {
			t.Errorf("Test %d :  output is %q  but was expecting %s", testNumber, out, test.Expected)
		}
		b.Reset()
	}
}

func TestCAD_MarshalJson(t *testing.T) {
	tests := []struct {
		Input    CAD
		Expected string
	}{
		{
			Input: CAD{
				cents: 1,
			},
			Expected: "$0.01\n",
		},
		{
			Input: CAD{
				cents: -7,
			},
			Expected: "-$0.07\n",
		},
		{
			Input: CAD{
				cents: 0,
			},
			Expected: "$0.00\n",
		},
		{
			Input: CAD{
				cents: 1785,
			},
			Expected: "$17.85\n",
		},
	}

	for testNumber, test := range tests {
		var b bytes.Buffer
		fmt.Fprintln(&b, test.Input)
		out := b.String()
		if out != test.Expected {
			t.Errorf("Test %d :  output is %q  but was expecting %s", testNumber, out, test.Expected)
		}
		b.Reset()
	}
}

func TestCAD_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		Input    []byte
		Expected CAD
	}{
		{
			Expected: CAD{
				cents: 1,
			},
			Input: []byte("1"),
		},
		{
			Expected: CAD{
				cents: -7,
			},
			Input: []byte("-7"),
		},
		{
			Expected: CAD{
				cents: 0,
			},
			Input: []byte("0"),
		},
		{
			Expected: CAD{
				cents: 1785,
			},
			Input: []byte("1785"),
		},
		{
			Expected: CAD{
				cents: 1785,
			},
			Input: []byte("\"$17.85\""),
		},
	}

	for testNumber, test := range tests {
		out := &CAD{
			cents: -100,
		}
		err := json.Unmarshal(test.Input, out)
		if err != nil {
			t.Error(err)
			continue
		}

		if out.cents != test.Expected.cents {
			t.Errorf("Test %d :  output is %s  but was expecting %s", testNumber, out, test.Expected)
		}
	}
}
