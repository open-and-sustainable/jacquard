# Synthetic interface examples

All source text, claims, and populations here are invented for Jacquard. The two
output files are hand-authored schema examples. The checker example uses
illustrative `a…a` and `b…b` hash values; document-checking
reports must hash the actual contract and language table.

Integration tests validate these interface shapes and typed round trips.
The example output values describe the fixtures; evidence checks compare findings
with the source document, and audits compare counts with the item population.

The funnel example includes a missing field, an unparsed measurement, zero,
different positive deciles, and repeated linked values. Item `f` contributes once
to each of deciles 1 and 10, despite its repeated value of 1. Every empty cell is
included in the example output.
