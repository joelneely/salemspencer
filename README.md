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

### Salem-Spencer Search (revised Go implementation)

N | Size | Count | Total time | Unit time
:-: | :-: | :-: | :-: | :-:
1 | 1 | 1 | 2.834µs | 2.584µs
2 | 2 | 1 | 27.459µs | 750ns
3 | 2 | 3 | 31.917µs | 833ns
4 | 3 | 2 | 37.209µs | 2µs
5 | 4 | 1 | 42.25µs | 1.791µs
6 | 4 | 4 | 50.292µs | 4.708µs
7 | 4 | 10 | 90.584µs | 37.084µs
8 | 4 | 25 | 106.209µs | 12.084µs
9 | 5 | 4 | 122.209µs | 12.709µs
10 | 5 | 24 | 145.167µs | 19.958µs
11 | 6 | 7 | 175.5µs | 27.458µs
12 | 6 | 25 | 224.125µs | 45.541µs
13 | 7 | 6 | 285.209µs | 58.042µs
14 | 8 | 1 | 359.959µs | 71.792µs
15 | 8 | 4 | 482.792µs | 119.792µs
16 | 8 | 14 | 686.542µs | 200.792µs
17 | 8 | 43 | 1.022375ms | 332.5µs
18 | 8 | 97 | 1.553667ms | 527µs
19 | 8 | 220 | 2.717084ms | 1.158959ms
20 | 9 | 2 | 4.122209ms | 1.385ms
21 | 9 | 18 | 6.279917ms | 2.149917ms
22 | 9 | 62 | 9.503292ms | 3.212417ms
23 | 9 | 232 | 14.357334ms | 4.845875ms
24 | 10 | 2 | 19.891709ms | 5.5105ms
25 | 10 | 33 | 27.901584ms | 7.994625ms
26 | 11 | 2 | 36.912375ms | 8.992458ms
27 | 11 | 12 | 50.114709ms | 13.186625ms
28 | 11 | 36 | 68.608917ms | 18.476ms
29 | 11 | 106 | 95.581792ms | 26.953208ms
30 | 12 | 1 | 127.97775ms | 32.38075ms
31 | 12 | 11 | 175.86825ms | 47.873666ms
32 | 13 | 2 | 230.473584ms | 54.576584ms
33 | 13 | 4 | 313.852625ms | 83.337791ms
34 | 13 | 14 | 437.50725ms | 123.611583ms
35 | 13 | 40 | 622.563709ms | 185.009875ms
36 | 14 | 2 | 834.685834ms | 212.079375ms
37 | 14 | 4 | 1.146833542s | 312.104ms
38 | 14 | 86 | 1.608748334s | 461.867334ms
39 | 14 | 307 | 2.284896375s | 676.104791ms
40 | 15 | 20 | 3.057892417s | 772.9505ms
41 | 16 | 1 | 3.925507417s | 867.564917ms
42 | 16 | 4 | 5.197943375s | 1.272393041s
43 | 16 | 14 | 7.049222334s | 1.8512345s
44 | 16 | 41 | 9.735946209s | 2.686677167s
45 | 16 | 99 | 13.619907375s | 3.883912333s
46 | 16 | 266 | 19.226968709s | 5.607017709s
47 | 16 | 674 | 27.289871042s | 8.062857792s
48 | 16 | 1505 | 38.793948917s | 11.504033167s
49 | 16 | 3510 | 55.22905425s | 16.435039458s
50 | 16 | 7726 | 1m18.69316725s | 23.464065958s
51 | 17 | 14 | 1m45.21406925s | 26.520860708s
52 | 17 | 50 | 2m23.066534042s | 37.852414333s
53 | 17 | 156 | 3m16.4607845s | 53.394198s
54 | 18 | 2 | 4m15.993352542s | 59.532519625s
55 | 18 | 8 | 5m39.200345625s | 1m23.206942875s
56 | 18 | 26 | 7m37.0556005s | 1m57.855204666s

### To-do


### Done

* I had failed to move all of the constant declarations from ssmain.go to ssdata/ssset.go, so needed to fix that.
* The first implementation of SSSet.data used slices (`[]uint8`), which prevented the use of a hashmap to store previously-found sets. I replacing that slice with an array, then replaced the slice holding maximal set with a map (eliminating the array search to prevent duplicates). These two changes reduced the processing time by about 1/3 (the previous total time for _N_=45 was 31.73s), with most of the gains coming from the first of those two changes.
* Switched the five hot-path `SSSet` methods (`Equals`, `IsClosedAt`, `IsOpenAt`, `Move`, `MoveLR`) from value receivers to pointer receivers, eliminating 92-byte struct copies on every call. This reduced the unit time at _N_=45 from 5.909s to 3.884s — a **~34% speedup** — with no changes required at call sites.
