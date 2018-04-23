# bfactor
Factoring numbers with Bill

So this repository started when one day my friend Bill asked me what the smallest number that had 36 factors was. He came up with an algorithm, which was to find the prime factorization of 36 (or whatever), then using those as the exponents (plus one) of the numbers used to create the number with 36 (or whatever) factors, distributing them in "reverse" order so the largest exponents are used on the smallest prime (2 has the most, then 3, etc).

This instantly came up with 1,260 as the answer, which was correct. But does this algorithm *always* use the correct algorithm? To find out, I coded up code to test all the possibilities. It was a tricky combinatorics problem, but I figured out a way to do it. The way I do it is by using bitflags, one for each of the factors of 36 (or whatever), and, in the outer loop, it counts down from the number with all the bits set down to 0, and that represents the factors that will go in the first group, and in the first inner loop, it uses the outer loop as a mask to find all combinations of the bits not used in the outer loop, and that represents the factors that go in the second group, and in the next inner loop, it combines both of the outer loops to make a mask so it finds all combinations of bits not used by either of the outer two, and that represents the factors that go in the third group, and so on -- with no limit to the number of inner loops. In this way, the system can find all combinations of factors with no fixed number of groups.

As it turns out, this method found repetitions and also combinations where not all the factors were used at all. So, added some code to combine all the bitflags from all the loops to verify that all the factors were being used. Didn't do anything about the repetitions, which comes from the fact that factors are repeated and the combinatorics code assumes they are all unique (e.g. 2, 2, 3, 3 for 36, the 2s and 3s are not unique, so we repeat them in the calculations). The code probably could be optimized to remove these duplicate calculations, which are harmless but waste a bit of computing power.

Anyway, running Bill's algorithm against my "brute force" algorithm to test all combinations, we found Bills algorithm was almost always right -- but not always. The first number for which it was wrong was 8 factors -- I've been saying "or whatever" because we made the 36 a parameter that could be changed, and then looped through to find all the numbers that were the minimum for N factors where N goes in a sequence. For 8 factors, Bill's algorithm said 30, but mine said 24, and, 24 does indeed have 8 factors (1, 2, 3, 4, 6, 8, 12, and 24), so 24 is lower and is in fact the correct answer. (We ultimately verified this by Google search and the encyclopedia of integer sequences).

The code here provides the following functions:

findMinWithFactors() finds the smallest number with N factors, e.g. findMinWithFactors(36) finds the smallest number with 36 factors (which is 24), using my combinatorics algorithm.
billAlgorithm() finds the smallest using Bill's algorithm.
findBillAlgoDifferent() finds where the two algorithms are different.
findSeries() generates a series of the smallest number with N factors for sequential N, and
findHighlyCompositeNumbers(), finds a list of highly composite numbers, which I haven't mentioned thus far, but is what got me and Bill talking about "36 factors" in the first place

In addition to all this, this code has some handy helper functions that you can use anywhere:

calcPrimes(), which calculates primes (using the sieve of Eratosthenes algorithm), up to some limit (the other functions require a list of primes as input)
calcFactors(), which finds the prime factors of some number
countFactors(), which counts the number of factors a number has given its prime factorization
showFactorization(), which displays the prime factorization to the user

The code does have some issues, notably that, even when using unsigned 64-bit integers, it frequently overflows. To understand why this is the case, consider the prime number 67. The only numbers that have 67 factors are prime numbers to the 68th power. The smallest of these is 2^68. However, a 64-bit unsigned int can only go up to 2^65 - 1. Therefore it is easy to see, that as soon as you ask for the minimum number with N factors, as soon as N gets up to 67 this program will fail to come up with the answers. This can be remedied by using Go's math/big package for infinite-precision arithmetic. Bill wrote his program in Python, though, which automatically bumped his variables up to infinite precision ints, so was able to calculate all of these without modifying his code. However, since we are searching for the minimums, and as long as we're interested in minimums that aren't huge numbers (N highly composite instead of prime), this code will work ok even with this issue. But it's something to be aware of before you put this code into some production system or something.

