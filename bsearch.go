package main

import "sort"

//half it, check values -> look half it again -> ...
//After finding the the values, look up the index in keepSort
func bsearch(lsorted []uint32, hsorted []uint32, symbolInterval uint32) uint8 {

	//Binary search
	start := 0
	end := len(lsorted) - 1
	match := 0
	for start <= end {
		middle := (start + end) / 2
		if symbolInterval >= lsorted[middle] && symbolInterval <= hsorted[middle] {
			match = middle
		}
		if symbolInterval > hsorted[middle] {
			start = middle + 1
		}
		if symbolInterval < lsorted[middle] {
			end = middle - 1
		}
	}
	/*
			//Used for later lookup(to get the index back)
		hmap := keep(low, high)
		lsorted, hsorted := sorter(low, high)
		hvalue := hsorted[match]
		//Look up this value in the saved map(either is fine)
		index := hmap[hvalue]
	*/
	return uint8(match)
}

/*Get 2 unordered slices, make map where key = index, value = sorted values
//Slices by default passed by ref.
//map[intervalLOW] = index
//Call this and sort once*/
func keep(low []uint32, high []uint32) map[uint32]int {
	hmap := make(map[uint32]int)
	//Store the original values in a map for easier lookup by value
	for i := 0; i < 256; i++ {
		hmap[high[i]] = i
	}
	return hmap
}
func sorter(low []uint32, high []uint32) ([]uint32, []uint32) {
	lsorted := low
	hsorted := high
	//Sort slices, hsorted and lsorted are sorted copies of low and high
	sort.Slice(lsorted, func(i, j int) bool { return lsorted[i] < lsorted[j] })
	sort.Slice(hsorted, func(i, j int) bool { return hsorted[i] < hsorted[j] })
	return lsorted, hsorted
}
