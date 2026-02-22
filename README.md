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

The last question is answered on [The On-Line Encyclopedia of Integer Sequences](https://oeis.org/), which is an absolute treasure for anyone who likes to play around with numbers (or who has real work to do with such information)! The counts for these maximally-sized subsets is given by [entry A262347](https://oeis.org/A262347). According to the history of that entry, progressively more values have been identified.

Contributor | Date | Maximum _N_
:---------: | :--: | :---------:
N. McNew | 2015-09-18 | 24
R. Israel | 2015-09-20 | 44
F. Cariboni | 2018-01-15 | 75
F. Cariboni | 2018-02-19 | 140

I was able to reproduce the first few entries of the sequence fairly quickly in JavaScript, but the time required for successive values of _N_ appears to grow exponentially. So, I became curious about the possibility of using a relatively current CPU and a performance-oriented language to match the previous work or even provide OEIS with a new record.

Hence, this little project, and the choice of Go.

I have already discovered a few minor optimizations for performance, but need to go further before it's practical to take on the current record of _N_ through 140.

### Usage

Requires Go 1.18+. The search limit _N_ is set by `LIMIT` in `ssdata/ssset.go` (currently 75).

```bash
# Run directly
go run ssmain.go

# Or build and run
go build -o salemspencer .
./salemspencer

# Run tests
go test ./...
```

Output is a Markdown table (like the one below) showing the maximal set size, count of distinct maximal sets, and cumulative and incremental timing for each _N_ from 1 to `LIMIT`.

### Salem-Spencer Search (revised Go implementation)

Timings below are from a MacBook Pro with an M2 Pro CPU.

N | Size | Count | Total time | Unit time
:-: | :-: | :-: | :-: | :-:
1 | 1 | 1 | 9.333µs | 8.833µs
2 | 2 | 1 | 55.5µs | 1.542µs
3 | 2 | 3 | 67.042µs | 1.959µs
4 | 3 | 2 | 79.167µs | 4.792µs
5 | 4 | 1 | 90.833µs | 3.791µs
6 | 4 | 4 | 108.167µs | 10.25µs
7 | 4 | 10 | 184.708µs | 68.125µs
8 | 4 | 25 | 221.875µs | 25.792µs
9 | 5 | 4 | 256.792µs | 26µs
10 | 5 | 24 | 304.583µs | 39.791µs
11 | 6 | 7 | 367.625µs | 54.792µs
12 | 6 | 25 | 465.958µs | 89.583µs
13 | 7 | 6 | 588.417µs | 113.584µs
14 | 8 | 1 | 746.792µs | 150.334µs
15 | 8 | 4 | 989.458µs | 234.041µs
16 | 8 | 14 | 1.406167ms | 408.042µs
17 | 8 | 43 | 2.072958ms | 654.125µs
18 | 8 | 97 | 3.150667ms | 1.068334ms
19 | 8 | 220 | 4.812708ms | 1.652291ms
20 | 9 | 2 | 6.854125ms | 2.02775ms
21 | 9 | 18 | 9.912667ms | 3.041875ms
22 | 9 | 62 | 14.175708ms | 4.235291ms
23 | 9 | 232 | 20.796875ms | 6.607833ms
24 | 10 | 2 | 28.503708ms | 7.69675ms
25 | 10 | 33 | 38.26275ms | 9.740083ms
26 | 11 | 2 | 48.814208ms | 10.539083ms
27 | 11 | 12 | 63.830333ms | 15.006458ms
28 | 11 | 36 | 84.233ms | 20.375625ms
29 | 11 | 106 | 112.438542ms | 28.1905ms
30 | 12 | 1 | 144.396208ms | 31.939958ms
31 | 12 | 11 | 191.835333ms | 47.423833ms
32 | 13 | 2 | 247.418667ms | 55.562084ms
33 | 13 | 4 | 330.800333ms | 83.3455ms
34 | 13 | 14 | 455.766042ms | 124.928042ms
35 | 13 | 40 | 640.284542ms | 184.4695ms
36 | 14 | 2 | 853.109167ms | 212.777834ms
37 | 14 | 4 | 1.168411292s | 315.242834ms
38 | 14 | 86 | 1.629601583s | 461.147208ms
39 | 14 | 307 | 2.313013125s | 683.358917ms
40 | 15 | 20 | 3.094501667s | 781.446542ms
41 | 16 | 1 | 3.974850292s | 880.305042ms
42 | 16 | 4 | 5.273787125s | 1.298902167s
43 | 16 | 14 | 7.118227917s | 1.844389334s
44 | 16 | 41 | 9.765020042s | 2.6467475s
45 | 16 | 99 | 13.58191425s | 3.816843333s
46 | 16 | 266 | 19.091391958s | 5.509432208s
47 | 16 | 674 | 27.000038958s | 7.908582166s
48 | 16 | 1505 | 38.290067458s | 11.289981416s
49 | 16 | 3510 | 54.438600375s | 16.148483792s
50 | 16 | 7726 | 1m17.731108625s | 23.292443958s
51 | 17 | 14 | 1m44.156129792s | 26.4249635s
52 | 17 | 50 | 2m21.839798583s | 37.683616s
53 | 17 | 156 | 3m15.33699425s | 53.497143833s
54 | 18 | 2 | 4m15.464001917s | 1m0.126960292s
55 | 18 | 8 | 5m39.762071917s | 1m24.298011125s
56 | 18 | 26 | 7m36.931934458s | 1m57.169809458s
57 | 18 | 56 | 10m19.995192125s | 2m43.063193417s
58 | 19 | 2 | 13m21.965277167s | 3m1.970022667s
59 | 19 | 4 | 17m34.352815792s | 4m12.387482042s
60 | 19 | 6 | 23m26.260105083s | 5m51.90724725s
61 | 19 | 14 | 31m36.885394583s | 8m10.624554541s
62 | 19 | 48 | 43m4.450671333s | 11m27.56464475s
63 | 20 | 2 | 55m39.590232125s | 12m35.1389575s
64 | 20 | 4 | 1h12m33.71035525s | 16m54.119635958s
65 | 20 | 8 | 1h36m30.406733875s | 23m56.69552575s

### To-do


### Done

* I had failed to move all of the constant declarations from ssmain.go to ssdata/ssset.go, so needed to fix that.
* The first implementation of SSSet.data used slices (`[]uint8`), which prevented the use of a hashmap to store previously-found sets. I replacing that slice with an array, then replaced the slice holding maximal set with a map (eliminating the array search to prevent duplicates). These two changes reduced the processing time by about 1/3 (the previous total time for _N_=45 was 31.73s), with most of the gains coming from the first of those two changes.
* Switched the five hot-path `SSSet` methods (`Equals`, `IsClosedAt`, `IsOpenAt`, `Move`, `MoveLR`) from value receivers to pointer receivers, eliminating 92-byte struct copies on every call. This reduced the unit time at _N_=45 from 5.909s to 3.884s — a **~34% speedup** — with no changes required at call sites.
