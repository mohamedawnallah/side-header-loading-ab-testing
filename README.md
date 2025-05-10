# Side-Header Loading A/B Testing Experiments
This repository serves as a knowledge base documenting my A/B testing experiments. Here, I explore and evaluate various hypotheses, recording methodologies, results, and insights gained throughout the process.

## Hypotheses Regarding `filters_db/header-index` and `block_headers_bin`

| # | Header in `filters_db` | Header in `block_headers_bin` | Description                                                  | Variables                       | Default Outcome (`H₀`)                                         | Accepted or Rejected |
|---|:----------------------:|:-----------------------------:|--------------------------------------------------------------|-------------------------------------------|----------------------------------------------------------------|:-------------------:|
| 1 | ✔️                     | ❌                            | Header exists in `filters_db/header-index` but not in `block_headers_bin` | – | Error creating chain service: unable to read block header: `EOF` | Accepted            |
| 2 | ❌                     | ✔️                            | Header not in `filters_db/header-index` but exists in `block_headers_bin` | – | Error creating chain service: target height not found in index  | Accepted            |
| 3 | ❌                     | ❌                            | Removing Tail Header from both stores                                | Not updating chain tip for both btc and regular header | Error creating chain service: target height not found in index                                                           | Accepted                     |
| 4 | ❌                     | ❌                            | Removing Tail Header from both stores                                | Updating chain tip for btc only and not regular header | Error creating chain service: target height not found in index                                                           | Accepted                     |
| 5 | ❌                     | ❌                            | Removing Tail Header from both stores                                | Updating chain tips for btc and  regular header | `OK`                                                           | Accepted                     |
| 6 | ❌                     | ❌                            | Removing Tail Header from both stores                                | Updating chain tips for btc and  regular header | `OK` and that Tail Filter header should be computed and indexed automatically                                                         | –                     |
| 7 | ❌                     | ❌                            | Removing Head Header from both stores                                | Not Removing Filter Header from store | `OK`                                                           | –                     |
| 8 | ❌                     | ❌                            | Removing Head Header from both stores                                | Removing Filter Header from store | `OK` and that Head Filter header should be computed and indexed automatically                                                         | –                     |
| 9 | ❌                     | ❌                            | Removing Middle Header from both stores                                | Not Removing Filter Header from store | `OK`                                                           | –                     |
| 10 | ❌                     | ❌                            | Removing Middle Header from both stores                                | Removing Filter Header from store | `OK`  and that Mid Filter header should be computed and indexed automatically                                                          | –                     |
| 11 | ✔️                     | ✔️                            | Header exists in both                                        | – | `OK`                                                           | Accepted            |

## Hypotheses Regarding `filters_db/regular/filter-store` and `regular_filter_headers_bin`


## Hypotheses Regarding Side Loading of Block AND Filter Headers

## Conclusions

1. The permitted operations for side-loading are limited to either no-operation (no-op) or extending. This enforces the **monotonicity/append only property**, meaning the state can only remain the same or grow, but never decrease or revert. This restriction is based on mutual trust, as established in the merge intervals algorithm and the process of overlapping checkpointing.
2. When side-loading block headers for a given range (e.g., blocks M to N), the corresponding filter headers for the same range (M to N) must also be provided. The **conjunction property** requires that both block headers and filter headers for the specified range are present and valid. The semantic mapping between block headers and filter headers is assumed to be correct and trusted; however, syntactic validation may still be performed during side-loading to ensure structural correctness.
3. The **idempotence property** ensures that performing the same side-loading operation multiple times has the same effect as performing it once. In other words, repeated application of a no-op or extending operation does not result in inconsistencies or unintended changes to the state.
4. The **atomicity property** guarantees that each side-loading operation is all-or-nothing: either the entire operation is applied successfully, or no part of it takes effect. This prevents partial application of data, which could result in inconsistencies or corruption.

## Rationale behind Conclusions

1. By explicitly disallowing replace, update, and delete operations, and limiting actions to no-op or extending, the system enforces **monotonicity/append only property**—the state can only stay the same or grow. This greatly simplifies reasoning about system behavior and prevents inconsistencies or malicious rollbacks during side-loading. Relying on mutual trust and well-defined procedures (such as merge intervals and overlapping checkpointing) enables safe coordination of updates, while minimizing the attack surface and operational complexity.
2. Requiring both block headers and filter headers together enforces the **conjunction property**, ensuring data consistency and preventing partial or mismatched data from being introduced during side-loading. This maintains the integrity of the chain state and ensures that all necessary components are available for downstream processing or validation.
3. Enforcing the **idempotence property** means that side-loading operations can be safely retried or repeated without concern for corrupting or altering the intended state. This increases robustness and fault tolerance in distributed or unreliable environments, as repeated operations have no adverse side effects.
4. Enforcing the **atomicity property** ensures that side-loading operations are applied in their entirety or not at all. For example, if an operation fails partway—such as only the block headers are written but not the filter headers—the system should roll back or recover to a consistent state. This guarantees that the system never ends up in a partially updated or inconsistent state, preserving data integrity and reliability even in the case of errors or failures during the operation.
