package services

func Skeleton(bin [][]bool) {
	temp := bin
	count := 1

	//count = step(1, temp, bin)
	//temp = bin
	//count += step(2, temp, bin)
	//temp = bin
	for count > 0 {
		count = step(1, temp, bin)
		temp = bin
		count += step(2, temp, bin)
		temp = bin
	}
}

func step(stepNum int, temp [][]bool, s [][]bool) (count int) {
	count = 0
	for i := 1; i < len(temp)-1; i++ {
		for j := 1; j < len(temp[0])-1; j++ {
			if def(i, j, temp, stepNum == 2) {
				if s[i][j] {
					count++
					s[i][j] = false
				}
			}
		}
	}

	return
}

func def(x int, y int, s [][]bool, even bool) bool {
	p2 := s[x][y-1]
	p4 := s[x+1][y]
	p6 := s[x][y+1]
	p8 := s[x-1][y]

	bp1 := neighbors(x, y, s)
	if bp1 >= 2 && bp1 <= 6 {
		if transition(x, y, s) == 1 {
			if even {
				if !((p2 && p4) && p8) {
					if !((p2 && p6) && p8) {
						return true
					}
				}
			} else {
				if !((p2 && p4) && p6) {

					if !((p4 && p6) && p8) {
						return true
					}
				}
			}
		}
	}
	return false
}

func transition(x int, y int, s [][]bool) (count int) {
	p2 := s[x][y-1]
	p3 := s[x+1][y-1]
	p4 := s[x+1][y]
	p5 := s[x+1][y+1]
	p6 := s[x][y+1]
	p7 := s[x-1][y+1]
	p8 := s[x-1][y]
	p9 := s[x-1][y-1]

	count = Btoi(!p2 && p3) + Btoi(!p3 && p4) +
		Btoi(!p4 && p5) + Btoi(!p5 && p6) +
		Btoi(!p6 && p7) + Btoi(!p7 && p8) +
		Btoi(!p8 && p9) + Btoi(!p9 && p2)

	return
}

func neighbors(x int, y int, s [][]bool) (count int) {
	count = 0
	if s[x-1][y+1] {
		count++
	}
	if s[x-1][y+1] {
		count++
	}
	if s[x-1][y-1] {
		count++
	}
	if s[x][y+1] {
		count++
	}
	if s[x][y-1] {
		count++
	}
	if s[x+1][y] {
		count++
	}
	if s[x+1][y+1] {
		count++
	}
	if s[x+1][y-1] {
		count++
	}

	return
}

func Btoi(val bool) int {
	if val {
		return 1
	}
	return 0
}

func HashCode(bin [][]bool, i, j int) int {
	Btoi := func(bl bool) int {
		if bl {
			return 1
		}
		return 0
	}

	hash := 0
	hash += 128 * Btoi(bin[i-1][j-1])
	hash += 64 * Btoi(bin[i-1][j])
	hash += 32 * Btoi(bin[i-1][j+1])

	hash += 16 * Btoi(bin[i][j-1])
	hash += 8 * Btoi(bin[i][j+1])

	hash += 4 * Btoi(bin[i+1][j-1])
	hash += 2 * Btoi(bin[i+1][j])
	hash += 1 * Btoi(bin[i+1][j+1])
	return hash
}
