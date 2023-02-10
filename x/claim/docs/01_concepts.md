<!--
order: 1
-->

# Concepts

Arkeo will facilitate an airdrop to multiple communities and user groups as outlined [here](https://www.notion.so/shapeshift/Arkeo-Airdrop-Spec-7d4abec5aa9444399b51a5ddb99a3a54)

Users are required to take multiple actions in order to recieve the full amount of their airdrop.

- 1/3 is sent to users when they claim their airdop
- 1/3 is send to users after they delegate arkeo tokens to a validator
- 1/3 is sent to users after they have voted in governance

Users of chains with differing address structures or derivation paths, specifically ethereum and thorchain, will be able to claim on arkeo using a signed message that transfers their airdrop from the designated Ethereum or Thorchain address to their Arkeo address.

Addresses eligible for native claims on Arkeo, will have a small amount of Arkeo in their accounts on genesis. This will be enough to pay for the gas fees of claiming their initial airdrop.

To incentivize users to claim in a timely manner, the amount of claimable airdrop reduces over time. Users can claim the full airdrop amount for three months (`DurationUntilDecay`).
After three months, the claimable amount linearly decays until 6 months after launch. (At which point none of it is claimable) This is controlled by the parameter `DurationOfDecay` in the code, which is set to 3 months. (6 months - 3 months).

After 6 months from launch, all unclaimed airdrop tokens are sent back to the community pool.
