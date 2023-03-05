# Regression Testing Framework

This is a port of the `thornode` regression testing framework: https://gitlab.com/thorchain/thornode/-/tree/develop/test/regression.

## Design

The high level structure of these tests is straightforward - each test case is defined in a YAML file consisting of start state and a set of interleaved transactions and checks. These test cases are organized hierarchically into directories as suites for testing specific features and boundary conditions. In order to avoid raciness in the test cases, we prevent blocks from being created and expose a special operation type to trigger the creation of blocks.

### Directory Structure

Test cases should be organized into directories as “suites” consisting of test cases for specific conditions or features.

```none
suites/
  initialize.yaml
  free-query.yaml
  paid-subscription.yaml
  micropayments.yaml
```

### Test Structure

The simplest of test structures may look something like:

```yaml
# yaml state deep merged with genesis
---
# create-block
---
# check: endpoint + jq conditions to assert
---
# transaction: ...
---
# create-block
---
# check: endpoint + jq conditions to assert
---
# ...
```

### Dynamic Values

In order to preserve the human-readability of test cases, the harness will populate embedded variables at runtime for addresses and transaction IDs. These values will be expressed as Go template functions and can be used in the test cases like:

- `{{ addr_dog }}` (the arkeo address for the "dog" mnemonic)
- `{{ native_txid 1 }}` (the txid of the first native transaction)
- `{{ native_txid -1 }}` (the txid of the most recent native transaction)
- `{{ template "default-state.yaml" }}` (default go template, embed the contents of the template `default-state.yaml`)

Addresses will be from the following mnemonics and each will be referenced by the corresponding animal:

```none
dog dog dog dog dog dog dog dog dog dog dog dog dog dog dog dog dog dog dog dog dog dog dog fossil
cat cat cat cat cat cat cat cat cat cat cat cat cat cat cat cat cat cat cat cat cat cat cat crawl
fox fox fox fox fox fox fox fox fox fox fox fox fox fox fox fox fox fox fox fox fox fox fox filter
pig pig pig pig pig pig pig pig pig pig pig pig pig pig pig pig pig pig pig pig pig pig pig quick
```

The `dog` mnemonic is a special case and will be used as the mnemonic for the default simulation validator, and as the mimir admin.

### State Definition

The state definition can be any valid data to deep merge with the default generated genesis file before initializing the simulation validator. There can be multiple `state` operations at the beginning of a test case:

```yaml
type: state
data:
  app_state: ...
```

### Check Definition

The check definition contains an endpoint, optional query parameters, and a set of `jq` query assertions against the endpoint response:

```yaml
type: check
endpoint:
params: {}
asserts:
```

### Transaction Definitions

There are multiple types of transactions that may be defined, which map to the protobuf types and should be self-explanatory by example:

```yaml
type: tx-send
from_address: {{ addr_fox }}
to_address: {{ addr_cat }}
amount:
  - amount: "200000000"
    asset: "rune"
```

### Create Blocks Definition

In order to allow the defining of test cases that are sensitive to the timing of blocks and placement of transactions therein, we expose the ability to explicitly trigger the creation of blocks during the test simulation:

```yaml
type: create-blocks
count: 1
```

If a specific transaction should cause the process to exit, an optional `exit` parameter will verify `arkeod` exits with the provided code.

## Tips for Writing Tests

The simplest way to approach test creation is to define state changes and transactions and keep the following operation at the end of the test file:

```yaml
type: check
endpoint: <endpoint>
```

This assertion will always fail and optinoally print the endpoint response to the console for inspecting the current state of a given endpoint after the test run up that point. Remember that the Cosmos and Arkeo APIs are available on port `1317` and the Tendermint APIs are available on port `26657`:

- https://v1.cosmos.network/rpc/v0.45.1
- https://docs.tendermint.com/master/rpc/#/

Pass the `RUN` environment variable the name of your test to avoid running all suites (it will also match a regex):

```bash
RUN=my-test make test-regression
```

If stuck set `DEBUG=1` to output the entire log output from the `arkeo` process and pause execution at the end of the test to inspect endpoints:

```bash
DEBUG=1 RUN=my-test make test-regression
```

Setting `EXPORT=1` will force overwrite the exported genesis after the test:

```bash
EXPORT=1 make test-regression
```

### Conventions

...

### Coverage

We leverage functionality in Golang 1.20 to track code coverage on the `arkeod` binary during live execution. Every run of the regression tests will generate a coverage percentage with archived, versioned, and generated code filtered - the value will be output to the console at the end of the test run. Coverage data is cleared after each run and a convenience target exists to open the coverage data from the last test run in the browser.

```bash
make test-regression-coverage
```

### Flakiness

Since block creation acquires a lock in process that will prevent query handling, all checks between blocks must complete within the block time - this block time defaults to `1s`. Additionally there is some raciness between the return of the application `EndBlock` and the time at which Tendermint, Cosmos, and Arkeo endpoints will execute against the new blocks data - we have a default sleep after the return of `EndBlock` set to `200ms`.

In order to avoid raciness more conveniently while running on resource constrained hardware, all time values above can optionally be increased by an integer factor defined in the `TIME_FACTOR` environment variable. If you find tests are hitting timeouts or returning inconsistent data, simply increase this factor (this will slow down the test run):

```bash
TIME_FACTOR=2 make test-regression
```
