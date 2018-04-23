package main

import (
	"fmt"
	"strconv"
)

// Sieve of Erastrathenes algorithm -- fastest way to get primes, but
// you have to specify what number to go up to in advance
func calcPrimes(upto int) []int {
	sieve := make([]bool, upto)
	result := make([]int, 0)

	p := 2
	z := p
	for p < upto {
		result = append(result, p)
		z = p
		for z < upto {
			sieve[z] = true
			z += p
		}
		p++
		if p == upto {
			return result
		}
		for sieve[p] {
			p++
			if p == upto {
				return result
			}
		}
	}
	return result // we will never reach this line
}

// One fuction for factoring a number that gives back both an ordered
// list of primes and exponents (factors and order) and an expanded
// list without exponents (sometimes more useful)
func calcFactors(primes []int, number int) (map[int]int, []int, []int) {
	factors := make(map[int]int)
	order := make([]int, 0)
	expanded := make([]int, 0)
	idx := 0
	for number != 1 {
		if (number % primes[idx]) == 0 {
			_, ok := factors[primes[idx]]
			if ok {
				factors[primes[idx]]++
			} else {
				factors[primes[idx]] = 1
				order = append(order, primes[idx])
			}
			expanded = append(expanded, primes[idx])
			number = number / primes[idx]
		} else {
			idx++
		}
	}
	return factors, order, expanded
}

// crude brute force integer exponentiation so we don't have to convert to float and back to use math.Exp()
// we're using small enough ints that this should be ok
func intexp(aa int, bb int) uint64 {
	var abig uint64
	var rr uint64
	var prev uint64
	rr = 1
	prev = 1
	abig = uint64(aa)
	for ii := 0; ii < bb; ii++ {
		prev = rr
		rr *= abig
		if rr < prev {
			// This should be an error-prone way of detecting overflows
			return 0 // too big for us to care about
		}
	}
	return rr
}

// finds the smallest using Bill's algorithm.
func billAlgorithm(factors []int, primes []int) uint64 {
	var total uint64
	total = 1
	currentPrimeIdx := 0
	for idx := len(factors) - 1; idx >= 0; idx-- {
		total *= intexp(primes[currentPrimeIdx], factors[idx]-1)
		currentPrimeIdx++
	}
	return total
}

func showIdxes(idxes []int, currentIdx int, comboResult int) {
	for ii := 0; ii <= currentIdx; ii++ {
		fmt.Print("    ", idxes[ii])
	}
	fmt.Println(" -> ", comboResult)
}

func usable(idxes []int, currentIdx int, bitflags int) bool {
	acc := 0
	for jj := 0; jj <= currentIdx; jj++ {
		acc |= idxes[jj]
	}
	if acc == bitflags {
		return true
	}
	return false
}

func doCrazyCalculation(idxes []int, currentIdx int, factors []int, primes []int) uint64 {
	// uncomment the prints in this function to see the expansion of all the exponents multiplied together to get a number
	var total uint64
	total = 1
	primenum := 0
	for ii := 0; ii <= currentIdx; ii++ {
		thisNum := 1
		flags := idxes[ii]
		bit := 1
		number := 0
		for flags != 0 {
			if (flags & bit) != 0 {
				thisNum *= factors[number]
				flags -= bit
			}
			bit *= 2
			number++
		}
		total *= intexp(primes[primenum], thisNum-1)
		// fmt.Print(" * ",primes[primenum], "^", thisNum - 1)
		primenum++
	}
	// fmt.Println(" = ", total)
	return total
}

func findAlls(factors []int, primes []int) uint64 {
	var bitflags int
	bitflags = int(intexp(2, len(factors))) - 1
	idxes := make([]int, len(factors)+1)
	masks := make([]int, len(factors)+1)
	idxes[0] = bitflags + 1
	masks[0] = bitflags
	currentIdx := 0
	keepgoing := true
	iterz := 65536 // Sanity loop breaker if things run forever
	var minimum uint64
	minimum = 0 // flag to indicate not set yet
	for keepgoing {
		idxes[currentIdx]--
		for (idxes[currentIdx] & masks[currentIdx]) != idxes[currentIdx] {
			idxes[currentIdx]--
		}
		if idxes[currentIdx] == 0 {
			// we hit the end -- back up to previous level
			currentIdx--
			if usable(idxes, currentIdx, bitflags) {
				comboResult := doCrazyCalculation(idxes, currentIdx, factors, primes)
				// showIdxes(idxes, currentIdx, comboResult)
				if minimum == 0 {
					minimum = comboResult
				} else {
					if comboResult < minimum {
						if comboResult > 0 {
							minimum = comboResult
						}
					}
				}
			}
			// no previous level? exit the whole function!
			if currentIdx < 0 {
				return minimum
			}
		} else {
			// else -- advance to NEXT level
			currentIdx++
			idxes[currentIdx] = bitflags + 1
			// sanity check
			// if (masks[currentIdx-1] - idxes[currentIdx-1]) != (masks[currentIdx-1] & (bitflags - idxes[currentIdx-1])) {
			//	fmt.Println("idxes", idxes)
			//	fmt.Println("masks", masks)
			//	fmt.Println("currentIdx", currentIdx)
			//	panic("Mask subtract failure")
			// }
			masks[currentIdx] = masks[currentIdx-1] - idxes[currentIdx-1]
		}

		// this is to break the loop if it becomes endless -- which should never happen
		iterz--
		if iterz == 0 {
			keepgoing = false
		}
	}
	return minimum
}

// finds the smallest number with N factors, e.g. findMinWithFactors(36)
// finds the smallest number with 36 factors (which is 24), using my
// combinatorics algorithm.
func findMinWithFactors(numFactors int) {
	primes := calcPrimes(100)
	factors, order, expanded := calcFactors(primes, numFactors)
	fmt.Println("factors", factors)
	fmt.Println("order", order)
	fmt.Println("expanded", expanded)
	billMin := billAlgorithm(expanded, primes)
	fmt.Println("billMin", billMin)
	minimum := findAlls(expanded, primes)
	fmt.Println("minimum", minimum)
}

// finds where the two algorithms are different.
func findBillAlgoDifferent() {
	primes := calcPrimes(1000)
	for numFactors := 1; numFactors < 1000; numFactors++ {
		_, _, expanded := calcFactors(primes, numFactors)
		billMin := billAlgorithm(expanded, primes)
		minimum := findAlls(expanded, primes)
		if billMin != minimum {
			if minimum > 0 {
				// fmt.Println("Bill min is different", numFactors)
				// fmt.Println("    Bill min:", billMin)
				// fmt.Println("    minimum:", minimum)
				fmt.Println("(", minimum, ",", billMin, ")")
			}
		}
	}
}

// generates a series of the smallest number with N factors for sequential N
func findSeries(upto int) {
	primes := calcPrimes(upto)
	for numFactors := 1; numFactors < upto; numFactors++ {
		_, _, expanded := calcFactors(primes, numFactors)
		minimum := findAlls(expanded, primes)
		fmt.Println(numFactors, minimum)
		// fmt.Println(minimum)
	}
}

// counts the number of factors a number has given its prime factorization
func countFactors(primes []int, factors map[int]int, order []int) int {
	total := 1
	for ii := 0; ii < len(order); ii++ {
		total *= factors[order[ii]] + 1
	}
	return total
}

func intToStr(nn int) string {
	return strconv.FormatInt(int64(nn), 10)
}

// displays the prime factorization to the user
func showFactorization(primes []int, factors map[int]int, order []int) {
	if len(order) == 0 {
		return
	}
	result := ""
	for ii := 0; ii < len(order); ii++ {
		result += " * " + intToStr(order[ii]) + "^" + intToStr(factors[order[ii]])
	}
	result = result[3:]
	fmt.Println(result)
}

// finds a list of highly composite numbers, which is what got me and
// Bill talking about "36 factors" in the first place
func findHighlyCompositeNumbers(upto int) {
	primes := calcPrimes(upto)
	maxSoFar := 0
	hasBeenUsed := make(map[int]bool)
	for number := 1; number < upto; number++ {
		factors, order, _ := calcFactors(primes, number)
		numFactors := countFactors(primes, factors, order)
		if numFactors >= maxSoFar {
			foundNew := false
			firstAppearance := ""
			for ii := 0; ii < len(order); ii++ {
				kk := order[ii]
				_, ok := hasBeenUsed[kk]
				if !ok {
					firstAppearance += ", first appearance of " + intToStr(order[ii])
					foundNew = true
					hasBeenUsed[kk] = true
				}
			}
			if foundNew {
				if numFactors == maxSoFar {
					fmt.Println("        ", number, "(", numFactors, "factors", firstAppearance, ")")
				}
			} else {
				if numFactors == maxSoFar {
					fmt.Println("        ", number, "(", numFactors, "factors )")
				}
			}
		}
		if numFactors == maxSoFar {
		}
		if numFactors > maxSoFar {
			maxSoFar = numFactors
			fmt.Println(number, "(", numFactors, "factors )")
			showFactorization(primes, factors, order)
		}
	}
}

func factorNumber(primes []int, number int) {
	fmt.Println("Factoring: ", number)
	factors, order, _ := calcFactors(primes, number)
	showFactorization(primes, factors, order)
	numFactors := countFactors(primes, factors, order)
	fmt.Println(numFactors, "factors")
}

func main() {
	// uncomment one of these to run
	findMinWithFactors(36)
	// findBillAlgoDifferent()
	// findSeries(1000)
	// findHighlyCompositeNumbers(10000)
	// some simple factorization
	primes := calcPrimes(100)
	factorNumber(primes, 24)
	factorNumber(primes, 60)
	factorNumber(primes, 1440)
	factorNumber(primes, 86400)
}
