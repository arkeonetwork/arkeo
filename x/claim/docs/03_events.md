<!--
order: 3
-->

# Events

## External module hooks

`claim` module emits one of the following events upon claiming:

| Type  | Attribute Key | Attribute Value |
| ----- | ------------- | --------------- |
| claim | sender        | {receiver}      |
| claim | amount        | {claim_amount}  |

| Type           | Attribute Key | Attribute Value |
| -------------- | ------------- | --------------- |
| claim_from_eth | sender        | {receiver}      |
| claim_from_eth | amount        | {claim_amount}  |