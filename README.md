# salemspencer

_Practicing Go by generating Salem-Spencer sequences_

A [recent article on QuantaMagazine](https://www.quantamagazine.org/to-win-this-numbers-game-learn-to-avoid-math-patterns-20200507/) described a "Salem-Spencer set" as a subset of the positive integers within some boundary (such as {1, 2, 3, 4, 5, 6, 7, 8, 9}) that contains no arithmetic progressions. For those who don't use mathematical jargon on a daily basis, that means that the set can't contain three evenly-spaced values, such as {1, 2, 3} or {1, 3, 5}. Within the first nine positive integers, {3, 5, 8, 9} meets this condition, as shown below.

Left | Step | Middle | Step | Right
:--: | :--: | :----: | :--: | :---:
_3_ | 2  | _5_ | 3 | _8_
_3_ | 2  | _5_ | 4 | _9_
_3_ | 5  | _8_ | 1 | _9_
_5_ | 3  | _8_ | 1 | _9_

For any three distinct values taken from that set, the difference from the first to the second is unequal to the difference from the second to the third.

On the other hand, {2, 3, 5, 8} fails, because it contains three evenly-spaced values (2, 5, and 8).

Left | Step | Middle | Step | Right
:--: | :--: | :----: | :--: | :---:
_2_ | 1  | _3_ | 2 | _5_
_2_ | 1  | _3_ | 5 | _8_
_2_ | **3**  | _5_ | **3** | _8_
_3_ | 2  | _5_ | 3 | _8_

The article poses a question: can one find a larger subset within {1, 2, ... 9} that still meets the condition? As programmers and mathematicians tend to do, we can immediately generalize that question. For a given value of _N_, can we find all such subsets of maximum size? What is the maximal size for a given _N_, and how many are there?

The last question is answered on [The On-Line Encyclopedia of Integer Sequences](https://oeis.org/), which is an absolute treasure for anyone who likes to play around with numbers (or who has real work to do with such information)! The counts for these maximally-sized subsets is given by (entry A262347)[https://oeis.org/A262347]. According to the history of that entry, progressively more values have been identified.

Contributor | Date | Maximum _N_
:---------: | :--: | :---------:
N. McNew | 2015-09-18 | 24
R. Israel | 2015-09-20 | 44
F. Cariboni | 2018-01-15 | 75
F. Cariboni | 2018-02-19 | 140

I was able to reproduce the first few entries of the sequence fairly quickly in JavaScript, but the time required for successive values of _N_ appears to grow exponentially. So, I became curious about the possibility of using a relatively current CPU and a performance-oriented language to match the previous work or even provide OEIS with a new record.

Hence, this little project, and the choice of Go.

I have already discovered a few minor optimizations for performance, but need to go further before it's practical to take on the current record of _N_ through 140.

### Salem-Spencer Search (first Go implementation)

N | Size | Count | Total Time
:-: | :-: | :-: | :-:
1 | 1 | 1 | 3.843µs
2 | 2 | 1 | 40.539µs
3 | 2 | 3 | 87.759µs
4 | 3 | 2 | 112.956µs
5 | 4 | 1 | 125.046µs
6 | 4 | 4 | 147.936µs
7 | 4 | 10 | 211.995µs
8 | 4 | 25 | 367.726µs
9 | 5 | 4 | 466.356µs
10 | 5 | 24 | 640.309µs
11 | 6 | 7 | 835.443µs
12 | 6 | 25 | 1.096034ms
13 | 7 | 6 | 1.339049ms
14 | 8 | 1 | 1.650173ms
15 | 8 | 4 | 2.086159ms
16 | 8 | 14 | 2.926158ms
17 | 8 | 43 | 4.029673ms
18 | 8 | 97 | 5.657035ms
19 | 8 | 220 | 8.452234ms
20 | 9 | 2 | 11.396026ms
21 | 9 | 18 | 16.009507ms
22 | 9 | 62 | 21.611098ms
23 | 9 | 232 | 29.854792ms
24 | 10 | 2 | 38.995551ms
25 | 10 | 33 | 53.477444ms
26 | 11 | 2 | 71.132202ms
27 | 11 | 12 | 97.051327ms
28 | 11 | 36 | 135.48256ms
29 | 11 | 106 | 194.488647ms
30 | 12 | 1 | 266.474042ms
31 | 12 | 11 | 369.880037ms
32 | 13 | 2 | 497.291588ms
33 | 13 | 4 | 684.022605ms
34 | 13 | 14 | 966.744526ms
35 | 13 | 40 | 1.386886481s
36 | 14 | 2 | 1.870248927s
37 | 14 | 4 | 2.586800647s
38 | 14 | 86 | 3.641474187s
39 | 14 | 307 | 5.185412762s
40 | 15 | 20 | 6.966690355s
41 | 16 | 1 | 8.975331207s
42 | 16 | 4 | 11.928156667s
43 | 16 | 14 | 16.282357294s
44 | 16 | 41 | 22.592319016s
45 | 16 | 99 | 31.730036466s

### To-do

* The current implementation of SSSet uses slices (`[]uint8`), which prevents the use of a hashmap to store previously-found sets. I previously made a quick-and-dirty stab at replacing the data slice with an array, but it was too dirty to work, so I backed it out. I believe that fixing this issue will improve performance.

### Done

* I had failed to move all of the constant declarations from ssmain.go to ssdata/ssset.go, so needed to fix that.
