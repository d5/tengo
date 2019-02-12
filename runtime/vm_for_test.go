package runtime_test

import (
	"testing"
)

func TestFor(t *testing.T) {
	expect(t, `
	for {
		out++
		if out == 5 {
			break
		}
	}`, 5)

	expect(t, `
	for {
		out++
		if out == 5 {
			break
		}
	}`, 5)

	expect(t, `
	a := 0
	for {
		a++
		if a == 3 { continue }
		if a == 5 { break }
		out += a
	}`, 7) // 1 + 2 + 4

	expect(t, `
	a := 0
	for {
		a++
		if a == 3 { continue }
		out += a
		if a == 5 { break }
	}`, 12) // 1 + 2 + 4 + 5

	expect(t, `
	for true {
		out++
		if out == 5 {
			break
		}
	}`, 5)

	expect(t, `
	a := 0
	for true {
		a++
		if a == 5 {
			break
		}
	}
	out = a`, 5)

	expect(t, `
	a := 0
	for true {
		a++
		if a == 3 { continue }
		if a == 5 { break }
		out += a
	}`, 7) // 1 + 2 + 4

	expect(t, `
	a := 0
	for true {
		a++
		if a == 3 { continue }
		out += a
		if a == 5 { break }
	}`, 12) // 1 + 2 + 4 + 5

	expect(t, `
	func() {
		for true {
			out++
			if out == 5 {
				return
			}
		}
	}()`, 5)

	expect(t, `
	for a:=1; a<=10; a++ {
		out += a
	}`, 55)

	expect(t, `
	for a:=1; a<=3; a++ {
		for b:=3; b<=6; b++ {
			out += b
		}
	}`, 54)

	expect(t, `
	func() {
		for {
			out++
			if out == 5 {
				break
			}
		}
	}()`, 5)

	expect(t, `
	func() {
		for true {
			out++
			if out == 5 {
				break
			}
		}
	}()`, 5)

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
	}()`, 5)

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
	}()`, 5)

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
	}()`, 5)

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
	}()`, 5)

	expect(t, `
	out = func() {
		sum := 0
		for a:=1; a<=10; a++ {
			sum += a
		}
		return sum
	}()`, 55)

	expect(t, `
	out = func() {
		sum := 0
		for a:=1; a<=4; a++ {
			for b:=3; b<=5; b++ {
				sum += b
			}
		}
		return sum
	}()`, 48) // (3+4+5) * 4

	expect(t, `
	a := 1
	for ; a<=10; a++ {
		if a == 5 {
			break
		}
	}
	out = a`, 5)

	expect(t, `
	for a:=1; a<=10; a++ {
		if a == 3 {
			continue
		}
		out += a
		if a == 5 {
			break
		}
	}`, 12) // 1 + 2 + 4 + 5

	expect(t, `
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
	}`, 12) // 1 + 2 + 4 + 5
}
