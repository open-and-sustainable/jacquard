# Synthetic interface examples

All source text, claims, and populations here are invented for Jacquard. The two
output files are hand-authored examples of intended reports, not engine-produced
results. The checker example uses illustrative `a…a` and `b…b` hash values; real
reports must hash the actual contract and language table.

Current integration tests validate these interface shapes and typed round trips.
They do not assert checker/audit computation. Engine acceptance must compare
computed results against independent synthetic expectations before release.

The funnel example includes a missing field, an unparsed measurement, zero,
different positive deciles, and repeated linked values. Item `f` contributes once
to each of deciles 1 and 10, despite its repeated value of 1. Every empty cell is
included in the example output.
