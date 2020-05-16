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
