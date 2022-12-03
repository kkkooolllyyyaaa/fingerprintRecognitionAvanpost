package services

func Skeleton(bin [][]bool) {
	table := [256]int{1, 0, 0, 0, 0, 0, 0, 1, 0, 0, 1, 0, 0, 0, 1, 1, 0, 0, 1, 1, 0, 4,
		1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 3, 0, 1, 1, 1, 0, 4, 1, 1, 0, 0, 2, 2, 2, 0, 2, 2,
		0, 0, 1, 1, 0, 2, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 4, 4, 1, 1, 1, 0, 0, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1,
		1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 5, 5, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 2, 0, 0, 2, 2, 0, 0, 1, 1, 1, 0, 1, 1, 0, 0, 1, 1, 0, 0, 1, 1, 0,
		0, 0, 0, 0, 0, 0, 0, 0, 0, 5, 5, 0, 0, 0, 5, 0, 0, 0, 0, 2, 0, 2, 2, 4, 4, 1, 1, 4, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0,
		1, 1, 1, 1, 4, 0, 1, 1, 1, 4, 1, 1, 1, 3, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 4, 4, 1, 1, 1, 4,
		1, 1, 1, 0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}

	mask := [9]int{128, 64, 32, 16, 0, 8, 4, 2, 1}

	for i := 0; i < 100; i++ {
		for i := 1; i < len(bin)-1; i++ {
			for j := 1; j < len(bin[i])-1; j++ {
				if bin[i][j] {
					local := locality(bin, i, j)
					k := defineK(local, mask)
					action := table[k]
					switch action {
					case 1:
						bin[i][j] = false
					case 2:
						bin[i][j] = false
						bin[i-1][j] = true
					case 3:
						bin[i][j] = false
						bin[i][j+1] = true
					case 4:
						bin[i][j] = false
						bin[i+1][j] = true
					case 5:
						bin[i][j] = false
						bin[i][j-1] = true
					}
				}
			}
		}
	}

}

func defineK(local [][]bool, mask [9]int) (k int) {
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

func Btoi(val bool) int {
	if val {
		return 1
	}
	return 0
}

func locality(m [][]bool, i int, j int) (res [][]bool) {
	res = make([][]bool, 3)
	for k := range res {
		res[k] = m[i+k-1][j-1 : j+2]
	}
	return
}

//func Skeleton(bin [][]bool) {
//	wasChanged := true
//
//	for wasChanged {
//		st1 := iterate(1, bin)
//		st2 := iterate(2, bin)
//
//		wasChanged = st1 == 0 && st2 == 0
//	}
//}
//
//func iterate(stepNum int8, matrix [][]bool) int {
//	step := make([][]int, 0)
//	for i := 1; i < len(matrix)-1; i++ {
//		for j := 1; j < len(matrix[i])-1; j++ {
//			if matrix[i][j] {
//				local := locality(matrix, i, j)
//				prm1, prm2 := params(local)
//				p2, p4, p6, p8 := local[0][1], local[1][2], local[2][0], local[1][0]
//				if prm2 >= 2 && prm2 <= 6 && prm1 == 1 {
//					if stepNum == 1 {
//						if (p2 || p4 || p6) && (p4 || p6 || p8) {
//							step = append(step, []int{i, j})
//						}
//					} else {
//						if (p2 || p4 || p8) && (p2 || p6 || p8) {
//							step = append(step, []int{i, j})
//						}
//					}
//				}
//			}
//		}
//	}
//
//	for _, pair := range step {
//		i := pair[0]
//		j := pair[1]
//		matrix[i][j] = false
//	}
//
//	return len(step)
//}
//

//
//func params(m [][]bool) (int, int) {
//	i := 1
//	j := 1
//	p2, p3, p4, p5, p6, p7, p8, p9 := m[i-1][j], m[i-1][j+1],
//		m[i][j+1], m[i+1][j+1], m[i+1][j], m[i+1][j-1], m[i][j-1], m[i-1][j-1]
//
//	//fmt.Printf("p2: %v \np3: %v \np4: %v \np5: %v \np6: %v \np7: %v \np8: %v \np9: %v \n\n",
//	//	p2, p3, p4, p5, p6, p7, p8, p9)
//
//	return Btoi(!p2 && p3) + Btoi(!p3 && p4) + Btoi(!p4 && p5) + Btoi(!p5 && p6) +
//			Btoi(!p6 && p7) + Btoi(!p7 && p8) + Btoi(!p8 && p9) + Btoi(!p9 && p2),
//		Btoi(p2) + Btoi(p3) + Btoi(p4) + Btoi(p5) + Btoi(p6) + Btoi(p7) + Btoi(p8) + Btoi(p9)
//}
//
