package runtime_test

import (
	"testing"
)

func TestFor(t *testing.T) {
	expect(t, `
	out = 0
	for {
		out++
		if out == 5 {
			break
		}
	}`, nil, 5)

	expect(t, `
	out = 0
	for {
		out++
		if out == 5 {
			break
		}
	}`, nil, 5)

	expect(t, `
	out = 0
	a := 0
	for {
		a++
		if a == 3 { continue }
		if a == 5 { break }
		out += a
	}`, nil, 7) // 1 + 2 + 4

	expect(t, `
	out = 0
	a := 0
	for {
		a++
		if a == 3 { continue }
		out += a
		if a == 5 { break }
	}`, nil, 12) // 1 + 2 + 4 + 5

	expect(t, `
	out = 0
	for true {
		out++
		if out == 5 {
			break
		}
	}`, nil, 5)

	expect(t, `
	a := 0
	for true {
		a++
		if a == 5 {
			break
		}
	}
	out = a`, nil, 5)

	expect(t, `
	out = 0
	a := 0
	for true {
		a++
		if a == 3 { continue }
		if a == 5 { break }
		out += a
	}`, nil, 7) // 1 + 2 + 4

	expect(t, `
	out = 0
	a := 0
	for true {
		a++
		if a == 3 { continue }
		out += a
		if a == 5 { break }
	}`, nil, 12) // 1 + 2 + 4 + 5

	expect(t, `
	out = 0
	func() {
		for true {
			out++
			if out == 5 {
				return
			}
		}
	}()`, nil, 5)

	expect(t, `
	out = 0
	for a:=1; a<=10; a++ {
		out += a
	}`, nil, 55)

	expect(t, `
	out = 0
	for a:=1; a<=3; a++ {
		for b:=3; b<=6; b++ {
			out += b
		}
	}`, nil, 54)

	expect(t, `
	out = 0
	func() {
		for {
			out++
			if out == 5 {
				break
			}
		}
	}()`, nil, 5)

	expect(t, `
	out = 0
	func() {
		for true {
			out++
			if out == 5 {
				break
			}
		}
	}()`, nil, 5)

	expect(t, `
	out = func() {
		a := 0
		for {
			a++
			if a == 5 {
				break
			}
		}
		return a
	}()`, nil, 5)

	expect(t, `
	out = func() {
		a := 0
		for true {
			a++
			if a== 5 {
				break
			}
		}
		return a
	}()`, nil, 5)

	expect(t, `
	out = func() {
		a := 0
		func() {
			for {
				a++
				if a == 5 {
					break
				}
			}
		}()
		return a
	}()`, nil, 5)

	expect(t, `
	out = func() {
		a := 0
		func() {
			for true {
				a++
				if a == 5 {
					break
				}
			}
		}()
		return a
	}()`, nil, 5)

	expect(t, `
	out = func() {
		sum := 0
		for a:=1; a<=10; a++ {
			sum += a
		}
		return sum
	}()`, nil, 55)

	expect(t, `
	out = func() {
		sum := 0
		for a:=1; a<=4; a++ {
			for b:=3; b<=5; b++ {
				sum += b
			}
		}
		return sum
	}()`, nil, 48) // (3+4+5) * 4

	expect(t, `
	a := 1
	for ; a<=10; a++ {
		if a == 5 {
			break
		}
	}
	out = a`, nil, 5)

	expect(t, `
	out = 0
	for a:=1; a<=10; a++ {
		if a == 3 {
			continue
		}
		out += a
		if a == 5 {
			break
		}
	}`, nil, 12) // 1 + 2 + 4 + 5

	expect(t, `
	out = 0
	for a:=1; a<=10; {
		if a == 3 {
			a++
			continue
		}
		out += a
		if a == 5 {
			break
		}
		a++
	}`, nil, 12) // 1 + 2 + 4 + 5
}
