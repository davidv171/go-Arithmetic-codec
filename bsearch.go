package main

/*Get 2 unordered slices, make map where key = index, value = sorted values
//Slices by default passed by ref.
//map[intervalLOW] = index
//Call this and sort once*/
func keep(high []uint32) map[uint32]int {
	hmap := make(map[uint32]int)
	//Store the original values in a map for easier lookup by value
	for i := 0; i < 256; i++ {
		hmap[high[i]] = i
	}
	return hmap
}
