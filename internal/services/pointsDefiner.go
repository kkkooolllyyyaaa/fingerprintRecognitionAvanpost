package services

func DefineDots(bin [][]bool) (map[int]int, [][]int) {

	hashesMp := make(map[int]int, 0)
	pointI := make([]int, 0)
	pointJ := make([]int, 0)

	for i := 1; i < len(bin)-1; i++ {
		for j := 1; j < len(bin[0])-1; j++ {
			if bin[i][j] {
				local := locality(bin, i, j)
				k := defineK(local)
				if inHashes(k) {
					hashesMp[k]++
					pointI = append(pointI, i)
					pointJ = append(pointJ, j)
				}
			}
		}
	}

	return hashesMp, [][]int{pointI, pointJ}
}

func inHashes(k int) bool {
	hashes := [...]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 13, 16, 20, 21, 26, 29, 32, 33,
		35, 37, 39, 40, 41, 45, 49, 50, 51, 53, 55, 57, 58, 59, 61, 64, 69, 71, 74,
		76, 77, 78, 81, 82, 83, 85, 87, 88, 89, 90, 91, 92, 93, 94, 96, 97, 101, 103,
		109, 113, 114, 115, 117, 119, 125, 128, 132, 133, 134, 135, 138, 140, 141, 142,
		144, 148, 149, 154, 156, 157, 160, 161, 162, 163, 164, 165, 166, 167, 168, 169,
		170, 172, 173, 174, 175, 176, 177, 178, 179, 180, 181, 182, 183, 185, 186, 187,
		188, 189, 190, 191, 192, 196, 197, 202, 204, 205, 206, 207, 213, 218, 221, 225,
		226, 227, 228, 229, 230, 231, 236, 237, 238, 239, 245}
	for _, hash := range hashes {
		if hash == k {
			return true
		}
	}
	return false
}

func defineK(local [][]bool) (k int) {
	mask := [9]int{128, 64, 32, 16, 0, 8, 4, 2, 1}
	k = 0
	mdx := 0
	for i, _ := range local {
		for j := 0; j < 3; j++ {
			k += Btoi(local[i][j]) * mask[mdx]
			mdx++
		}
	}
	return
}

func locality(m [][]bool, i int, j int) (res [][]bool) {
	res = make([][]bool, 3)
	for k := range res {
		res[k] = m[i+k-1][j-1 : j+2]
	}
	return
}
