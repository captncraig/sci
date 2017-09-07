package resources

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
)

//The SCI0 map file format is pretty simple: It consists of 6-byte entries, terminated by the sequence 0xffff ffff ffff.
// The first 2 bytes, interpreted as little endian 16 bit integer, encode resource type (high 5 bits) and number (low 11 bits).
//The next 4 bytes are a 32 bit LE integer that contains the resource file number in the high 6 bits,
//and the absolute offset within the file in the low 26 bits.
//SCI0 performs a linear search to find the resource; however, multiple entries may match the search, since resources may be present
//more than once (the inverse mapping is not injective).
func ReadMap(dir string) error {
	dat, err := ioutil.ReadFile(filepath.Join(dir, "RESOURCE.MAP"))
	if err != nil {
		return err
	}
	typeCounts := map[byte]int{}
	for i := 0; i < len(dat)-5; i += 6 {
		fmt.Println(dat[i : i+6])
		rType := dat[i+1] >> 3
		fmt.Println(rType)
		typeCounts[rType]++
	}
	fmt.Println(typeCounts)
	return nil
}
