# Election Simulator

This tool simulates elections using different methods. Right now, it uses
- [Plurality](https://en.wikipedia.org/wiki/Plurality_voting)
- [Instant-Runoff](https://en.wikipedia.org/wiki/Instant-runoff_voting)
- [Schulze method](https://en.wikipedia.org/wiki/Schulze_method) (Condorcet)

It generates a CSV of election results via each method.

## Usage
```
Usage of election-simulator:
  -candidates int
        Sets the number of candidates per race. (default 3)
  -output string
        The output file for election results. (default "results.csv")
  -total-ballots int
        Sets the total number of ballots per election. (default 10000)
  -total-elections int
        Sets the total number of elections to run. (default 10000)
```