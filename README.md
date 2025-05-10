# Side-Header Loading A/B Testing Experiments
This repository serves as a knowledge base documenting my A/B testing experiments. Here, I explore and evaluate various hypotheses, recording methodologies, results, and insights gained throughout the process.

## Hypotheses Across `filters_db/header-index` and `block_headers_bin`

| # | Header in `filters_db` | Header in `block_headers_bin` | Description                                                  | Default Outcome (`H₀`)                |
|---|:----------------------:|:-----------------------------:|--------------------------------------------------------------|----------------------------------------|
| 1 | ✔️                     | ❌                            | Header exists in `filters_db/header-index` but not in `block_headers_bin` | `EOF` or `NotFound` in headers         |
| 2 | ❌                     | ✔️                            | Header not in `filters_db/header-index` but exists in `block_headers_bin` | `EOF` or `NotFound` in `filters/header-index`         |
| 3 | ❌                     | ❌                            | Header not in either store                                   | `OK`                                   |
| 4 | ✔️                     | ✔️                            | Header exists in both                                        | `OK`                                   |
