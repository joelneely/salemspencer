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
1 | 1 | 1 | 10.458µs | 9.846µs
2 | 2 | 1 | 39.472µs | 1.324µs
3 | 2 | 3 | 46.723µs | 4.37µs
4 | 3 | 2 | 52.377µs | 2.824µs
5 | 4 | 1 | 63.14µs | 8.124µs
6 | 4 | 4 | 75.248µs | 9.21µs
7 | 4 | 10 | 129.206µs | 51.309µs
8 | 4 | 25 | 161.9µs | 28.946µs
9 | 5 | 4 | 192.948µs | 27.688µs
10 | 5 | 24 | 252.95µs | 54.929µs
11 | 6 | 7 | 319.695µs | 61.868µs
12 | 6 | 25 | 431.707µs | 108.826µs
13 | 7 | 6 | 567.916µs | 131.51µs
14 | 8 | 1 | 742.112µs | 169.605µs
15 | 8 | 4 | 1.02709ms | 280.828µs
16 | 8 | 14 | 1.488549ms | 456.88µs
17 | 8 | 43 | 2.251725ms | 758.899µs
18 | 8 | 97 | 3.498938ms | 1.242416ms
19 | 8 | 220 | 5.294901ms | 1.789115ms
20 | 9 | 2 | 7.31615ms | 2.016582ms
21 | 9 | 18 | 10.362176ms | 3.041573ms
22 | 9 | 62 | 15.03003ms | 4.663711ms
23 | 9 | 232 | 21.901343ms | 6.855117ms
24 | 10 | 2 | 29.788852ms | 7.875433ms
25 | 10 | 33 | 40.461324ms | 10.665428ms
26 | 11 | 2 | 54.733126ms | 14.258ms
27 | 11 | 12 | 73.614609ms | 18.866431ms
28 | 11 | 36 | 103.358128ms | 29.728164ms
29 | 11 | 106 | 145.85737ms | 42.481511ms
30 | 12 | 1 | 196.780976ms | 50.906935ms
31 | 12 | 11 | 273.502778ms | 76.703212ms
32 | 13 | 2 | 359.195864ms | 85.675736ms
33 | 13 | 4 | 487.156511ms | 127.944688ms
34 | 13 | 14 | 679.408913ms | 192.217127ms
35 | 13 | 40 | 960.444445ms | 281.019276ms
36 | 14 | 2 | 1.286233047s | 325.773181ms
37 | 14 | 4 | 1.76699717s | 480.746023ms
38 | 14 | 86 | 2.471030608s | 704.016256ms
39 | 14 | 307 | 3.497364885s | 1.026316827s
40 | 15 | 20 | 4.685364008s | 1.187981077s
41 | 16 | 1 | 6.015981032s | 1.330598588s
42 | 16 | 4 | 7.964285081s | 1.94828555s
43 | 16 | 14 | 10.788987051s | 2.824683937s
44 | 16 | 41 | 14.890528236s | 4.101522858s
45 | 16 | 99 | 20.800035648s | 5.909489448s

### To-do


### Done

* I had failed to move all of the constant declarations from ssmain.go to ssdata/ssset.go, so needed to fix that.
* The first implementation of SSSet.data used slices (`[]uint8`), which prevented the use of a hashmap to store previously-found sets. I replacing that slice with an array, then replaced the slice holding maximal set with a map (eliminating the array search to prevent duplicates). These two changes reduced the processing time by about 1/3 (the previous total time for _N_=45 was 31.73s), with most of the gains coming from the first of those two changes.
