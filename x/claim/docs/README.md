# Claims module

## Abstract

The Arkeo claims module creates functionality to allow users to recieve airdroppped tokens. Native arkeo addresses
that are to recieve the airdrop will have a small amount of arkeo in their accounts on genesis. This will be enough to pay for the gas fees of claiming their initial airdrop.

Users of chains with differing address structures or derivation paths, specifically ethereum and thorchain, will be able to claim on arkeo using a signed message that transfers their airdrop from the designated Ethereum or Thorchain address to their Arkeo address. These users will either need to either use
a faucet to recieve a small amount of arkeo to pay for the gas fees of claiming their initial airdrop or we will need to determine another mechanism to
make this process as easy as possible.

Arkeo airdrop amounts 'expire' if not claimed. Users have three months (`DurationUntilDecay`) to claim their full airdrop amount.
After three months, the reward amount available will decline over 3 months (`DurationOfDecay`) in real time, until it hits `0%` at 6 months from launch (`DurationUntilDecay + DurationOfDecay`).

After 6 months from launch, all unclaimed tokens get sent to the community pool.

## Contents

1. **[Concept](01_concepts.md)**
2. **[State](02_state.md)**
3. **[Events](03_events.md)**
4. **[Hooks](04_hooks.md)**
5. **[Params](05_params.md)**

## Genesis State

## Actions

There are 3 types of actions, each of which release another 1/3 of the airdrop allocation.
The 3 actions are as follows:

```golang
	ACTION_CLAIM          Action = 0
	ACTION_VOTE           Action = 1
	ACTION_DELEGATE_STAKE Action = 2
```

The vote and delegate stake actions are monitored by registering claim **hooks** to the governance, and staking modules.
This means that when you perform an action, the claims module will immediately unlock those coins if they are applicable.
These actions can be performed in any order.

The code is structured by separating out a segment of the tokens as "claimable", indexed by each action type.
So if Alice delegates tokens, the claims module will move 1/3 of the claimables associated with staking to her liquid balance.
If she delegates again, there will not be additional tokens given, as the relevant action has already been performed.
Every action must be performed to claim the full amount.

## ClaimRecords

A claim record is a struct that contains data about the claims process of each airdrop recipient.

It contains the chain, an address, the initial claimable airdrop amount, and an array of bools representing
whether each action has been completed. The position in the array refers to enum number of the action.

So for example, `[true, true, false]` means that `ACTION_CLAIM` and `ACTION_VOTE` are completed.

```golang
// A Claim Records is the metadata of claim data per address
type ClaimRecord struct {
	Chain Chain `protobuf:"varint,1,opt,name=chain,proto3,enum=arkeonetwork.arkeo.claim.Chain" json:"chain,omitempty"`
	// arkeo address of claim user
	Address string `protobuf:"bytes,2,opt,name=address,proto3" json:"address,omitempty" yaml:"address"`
	// total initial claimable amount for the user
	InitialClaimableAmount types.Coin `protobuf:"bytes,3,opt,name=initial_claimable_amount,json=initialClaimableAmount,proto3" json:"initial_claimable_amount" yaml:"initial_claimable_amount"`
	// true if action is completed
	// index of bool in array refers to action enum #
	ActionCompleted []bool `protobuf:"varint,4,rep,packed,name=action_completed,json=actionCompleted,proto3" json:"action_completed,omitempty" yaml:"action_completed"`
}
```
