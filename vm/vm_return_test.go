package vm_test

import (
	"testing"
)

func TestReturn(t *testing.T) {
	expect(t, `out = func() { return 10; }()`, 10)
	expect(t, `out = func() { return 10; return 9; }()`, 10)
	expect(t, `out = func() { return 2 * 5; return 9 }()`, 10)
	expect(t, `out = func() { 9; return 2 * 5; return 9 }()`, 10)
	expect(t, `
	out = func() { 
		if (10 > 1) {
			if (10 > 1) {
				return 10;
	  		}

	  		return 1;
		}
	}()`, 10)

	expect(t, `f1 = func() { return 2 * 5; }; out = f1()`, 10)
}
