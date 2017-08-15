package hashid

import (
	hashids "github.com/speps/go-hashids"
	"github.com/tbbr/tbbr-api/config"
)

// Generate takes id of type uint and returns a string hash
func Generate(id uint) string {
	hd := hashids.NewData()
	hd.Salt = config.HashID.Salt
	hd.MinLength = config.HashID.MinLength
	h := hashids.NewWithData(hd)

	a := []int{0}
	a[0] = int(id)

	// Encode
	e, _ := h.Encode(a)
	return e
}
