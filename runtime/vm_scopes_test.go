package runtime_test

import "testing"

func TestVMScopes(t *testing.T) {
	// shadowed global variable
	expect(t, `
c := 5
if a := 3; a {
	c := 6
} else {
	c := 7
}
out = c
`, nil, 5)

	// shadowed local variable
	expect(t, `
func() {
	c := 5
	if a := 3; a {
		c := 6
	} else {
		c := 7
	}
	out = c
}()
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

	// shadowing inside for statement
	expect(t, `
a := 4
b := 5
for i:=0;i<3;i++ {
	b := 6
	for j:=0;j<2;j++ {
		b := 7
		a = i*j
	}
}
out = a`, nil, 2)

	// shadowing variable declared in init statement
	expect(t, `
if a := 5; a {
	a := 6
	out = a
}`, nil, 6)
	expect(t, `
a := 4
if a := 5; a {
	a := 6
	out = a
}`, nil, 6)
	expect(t, `
a := 4
if a := 0; a {
	a := 6
	out = a
} else {
	a := 7
	out = a
}`, nil, 7)
	expect(t, `
a := 4
if a := 0; a {
	out = a
} else {
	out = a
}`, nil, 0)

	// shadowing function level
	expect(t, `
a := 5
func() {
	a := 6
	a = 7
}()
out = a
`, nil, 5)
	expect(t, `
a := 5
func() {
	if a := 7; true {
		a = 8
	}
}()
out = a
`, nil, 5)
}
