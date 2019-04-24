package gpio

//Clock struct
type Clock struct {
	Pin
}

//PullUpDown (GPPUD)
func (c *Clock) PullUpDown() error {
	return nil
}
