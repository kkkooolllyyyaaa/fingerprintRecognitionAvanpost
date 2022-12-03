package services

type Bitset struct {
	width  int
	height int
	bin    [][]bool
}

func Skeleton(bin [][]bool) {
	matrix := bin
	wasChanged := true

	for wasChanged {
		st1 := iterate(1, matrix)
		st2 := iterate(2, matrix)

		wasChanged = st1 == 0 && st2 == 0
	}
}

func iterate(stepNum int8, matrix [][]bool) int {
	step := make([][]int, 0)
	for i := 1; i < len(matrix)-1; i++ {
		for j := 1; j < len(matrix)-1; j++ {
			if matrix[i][j] {
				local := locality(matrix, i, j)
				prm1, prm2 := params(local)
				p2, p4, p6, p8 := local[0][1], local[1][2], local[2][0], local[1][0]
				if prm2 >= 2 && prm2 <= 6 && prm1 == 1 {
					if stepNum == 1 {
						if (p2 || p4 || p6) && (p4 || p6 || p8) {
							step = append(step, []int{i, j})
						}
					} else {
						if (p2 || p4 || p8) && (p2 || p6 || p8) {
							step = append(step, []int{i, j})
						}
					}
				}
			}
		}
	}

	for _, pair := range step {
		i := pair[0]
		j := pair[1]
		matrix[i][j] = false
	}

	return len(step)
}

func locality(m [][]bool, i int, j int) (res [][]bool) {
	res = make([][]bool, 3)
	for k := range res {
		res[k] = m[i+k-1][j-1 : j+2]
	}
	return
}

func params(m [][]bool) (int, int) {
	i := 1
	j := 1
	p2, p3, p4, p5, p6, p7, p8, p9 := m[i-1][j], m[i-1][j+1],
		m[i][j+1], m[i+1][j+1], m[i+1][j], m[i+1][j-1], m[i][j-1], m[i-1][j-1]

	//fmt.Printf("p2: %v \np3: %v \np4: %v \np5: %v \np6: %v \np7: %v \np8: %v \np9: %v \n\n",
	//	p2, p3, p4, p5, p6, p7, p8, p9)

	return Btoi(!p2 && p3) + Btoi(!p3 && p4) + Btoi(!p4 && p5) + Btoi(!p5 && p6) +
			Btoi(!p6 && p7) + Btoi(!p7 && p8) + Btoi(!p8 && p9) + Btoi(!p9 && p2),
		Btoi(p2) + Btoi(p3) + Btoi(p4) + Btoi(p5) + Btoi(p6) + Btoi(p7) + Btoi(p8) + Btoi(p9)
}

func Btoi(val bool) int {
	if val {
		return 1
	}
	return 0
}
