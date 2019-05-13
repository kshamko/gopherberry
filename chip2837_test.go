package gopherberry

import (
	"fmt"
	"testing"
)

func TestAddressOffset(t *testing.T) {

	c := newChip2837()

	x := c.(*Chip2837).addressOffset(0x7E200000)

	/*addr1 := 0x7E200000
	addr2 := 0x7E200004
	addr3 := 0x7E200008
	addr4 := 0x7E20000C

	offset := 0x0000000*/
	fmt.Println(0x200004 - 0x200000)
	fmt.Printf("%X", x)
	//assert.Equals(t, addr1, addr2)
	t.Error("xxxx")
}
