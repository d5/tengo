package runtime_test

import "testing"

func TestVMScopes(t *testing.T) {
	// 'c' in outer scope
	expect(t, `
c := 5
if a := 3; a {
	c := 6 // shadowed
} else {
	c := 7 // shadowed
}
out = c
`, nil, 5)

	// 'b' is declared in 2 separate blocks
	expect(t, `
c := 5
if a := 3; a {
	b := 8
	c = b
} else {
	b := 9
	c = b
}
out = c
`, nil, 8)

	//
	expect(t, `
a := 5
for i:=0;i<3;i++ {
	b := 6
	for j:=0;j<2;j++ {
		b := 7
		a = i*j
	}
}
out = a`, nil, 2)
}
