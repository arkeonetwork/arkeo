package keeper

import (
	"mercury/common"
	"mercury/common/cosmos"
	"mercury/x/mercury/configs"
	"mercury/x/mercury/types"

	. "gopkg.in/check.v1"
)

type ClaimContractIncomeSuite struct{}

var _ = Suite(&ClaimContractIncomeSuite{})

func (ClaimContractIncomeSuite) TestValidate(c *C) {
	var err error
	ctx, k := SetupKeeper(c)
	ctx = ctx.WithBlockHeight(20)

	s := newMsgServer(k)

	// setup
	pubkey := types.GetRandomPubKey()
	acc := types.GetRandomBech32Addr()
	chain := common.BTCChain
	client := types.GetRandomPubKey()

	contract := types.NewContract(pubkey, chain, client)
	contract.Duration = 100
	contract.Rate = 10
	contract.Height = 10
	contract.Type = types.ContractType_PayAsYouGo
	contract.Deposit = cosmos.NewInt(contract.Duration * contract.Rate)
	c.Assert(k.SetContract(ctx, contract), IsNil)

	// happy path
	msg := types.MsgClaimContractIncome{
		PubKey:  pubkey,
		Chain:   chain.String(),
		Client:  client,
		Creator: acc.String(),
		Nonce:   20,
		Height:  10,
	}
	c.Assert(s.ClaimContractIncomeValidate(ctx, &msg), IsNil)

	// check bad height
	msg.Height = contract.Height * 2
	err = s.ClaimContractIncomeValidate(ctx, &msg)
	c.Check(err, ErrIs, types.ErrClaimContractIncomeBadHeight)

	// check closed contract
	msg.Height = contract.Height
	ctx = ctx.WithBlockHeight(ctx.BlockHeight() + contract.Duration)
	err = s.ClaimContractIncomeValidate(ctx, &msg)
	c.Check(err, ErrIs, types.ErrClaimContractIncomeClosed)
}

func (ClaimContractIncomeSuite) TestHandlePayAsYouGo(c *C) {
	ctx, k := SetupKeeper(c)

	s := newMsgServer(k)

	// setup
	pubkey := types.GetRandomPubKey()
	acc, err := pubkey.GetMyAddress()
	c.Assert(err, IsNil)
	chain := common.BTCChain
	client := types.GetRandomPubKey()
	c.Assert(k.MintToModule(ctx, types.ModuleName, getCoin(common.Tokens(10*100*2))), IsNil)
	c.Assert(k.SendFromModuleToModule(ctx, types.ModuleName, types.ContractName, getCoins(10*100)), IsNil)

	contract := types.NewContract(pubkey, chain, client)
	contract.Duration = 100
	contract.Rate = 10
	contract.Type = types.ContractType_PayAsYouGo
	contract.Deposit = cosmos.NewInt(contract.Duration * contract.Rate)
	c.Assert(k.SetContract(ctx, contract), IsNil)

	// happy path
	msg := types.MsgClaimContractIncome{
		PubKey:  pubkey,
		Chain:   chain.String(),
		Creator: acc.String(),
		Client:  client,
		Nonce:   20,
		Height:  ctx.BlockHeight(),
	}
	c.Assert(s.ClaimContractIncomeHandle(ctx, &msg), IsNil)

	c.Check(k.GetBalance(ctx, acc).AmountOf(configs.Denom).Int64(), Equals, int64(180))
	c.Check(k.GetBalanceOfModule(ctx, types.ContractName, configs.Denom).Int64(), Equals, int64(800))
	c.Check(k.GetBalanceOfModule(ctx, types.ReserveName, configs.Denom).Int64(), Equals, int64(20))

	// repeat the same thing and ensure we don't pay providers twice
	c.Assert(s.ClaimContractIncomeHandle(ctx, &msg), IsNil)
	c.Check(k.GetBalance(ctx, acc).AmountOf(configs.Denom).Int64(), Equals, int64(180))
	c.Check(k.GetBalanceOfModule(ctx, types.ContractName, configs.Denom).Int64(), Equals, int64(800))
	c.Check(k.GetBalanceOfModule(ctx, types.ReserveName, configs.Denom).Int64(), Equals, int64(20))

	// increase the nonce and get slightly more funds for the provider
	msg.Nonce = 25
	c.Assert(s.ClaimContractIncomeHandle(ctx, &msg), IsNil)
	acct := k.GetBalance(ctx, acc).AmountOf(configs.Denom).Int64()
	c.Check(acct, Equals, int64(225))
	cname := k.GetBalanceOfModule(ctx, types.ContractName, configs.Denom).Int64()
	c.Check(cname, Equals, int64(750))
	rname := k.GetBalanceOfModule(ctx, types.ReserveName, configs.Denom).Int64()
	c.Check(rname, Equals, int64(25))
	c.Check(rname+cname+acct, Equals, contract.Rate*contract.Duration)

	// ensure provider cannot take more than what is deposited into the account, overspend the contract
	msg.Nonce = contract.Deposit.Int64() / contract.Rate * 1000000000000
	c.Assert(s.ClaimContractIncomeHandle(ctx, &msg), IsNil)
	acct = k.GetBalance(ctx, acc).AmountOf(configs.Denom).Int64()
	c.Check(acct, Equals, int64(900))
	cname = k.GetBalanceOfModule(ctx, types.ContractName, configs.Denom).Int64()
	c.Check(cname, Equals, int64(0))
	rname = k.GetBalanceOfModule(ctx, types.ReserveName, configs.Denom).Int64()
	c.Check(rname, Equals, int64(100))
	c.Check(rname+cname+acct, Equals, contract.Rate*contract.Duration)
}

func (ClaimContractIncomeSuite) TestHandleSubscription(c *C) {
	ctx, k := SetupKeeper(c)
	ctx = ctx.WithBlockHeight(20)

	s := newMsgServer(k)

	// setup
	pubkey := types.GetRandomPubKey()
	acc, err := pubkey.GetMyAddress()
	c.Assert(err, IsNil)
	chain := common.BTCChain
	client := types.GetRandomPubKey()
	c.Assert(k.MintToModule(ctx, types.ModuleName, getCoin(common.Tokens(10*100*2))), IsNil)
	c.Assert(k.SendFromModuleToModule(ctx, types.ModuleName, types.ContractName, getCoins(10*100)), IsNil)

	contract := types.NewContract(pubkey, chain, client)
	contract.Duration = 100
	contract.Height = 10
	contract.Rate = 10
	contract.Type = types.ContractType_Subscription
	contract.Deposit = cosmos.NewInt(contract.Duration * contract.Rate)
	c.Assert(k.SetContract(ctx, contract), IsNil)

	// happy path
	msg := types.MsgClaimContractIncome{
		PubKey:  pubkey,
		Chain:   chain.String(),
		Creator: acc.String(),
		Client:  client,
		Nonce:   20,
		Height:  ctx.BlockHeight(),
	}
	c.Assert(s.ClaimContractIncomeHandle(ctx, &msg), IsNil)

	c.Check(k.GetBalance(ctx, acc).AmountOf(configs.Denom).Int64(), Equals, int64(90))
	c.Check(k.GetBalanceOfModule(ctx, types.ContractName, configs.Denom).Int64(), Equals, int64(900))
	c.Check(k.GetBalanceOfModule(ctx, types.ReserveName, configs.Denom).Int64(), Equals, int64(10))

	// repeat the same thing and ensure we don't pay providers twice
	c.Assert(s.ClaimContractIncomeHandle(ctx, &msg), IsNil)
	c.Check(k.GetBalance(ctx, acc).AmountOf(configs.Denom).Int64(), Equals, int64(90))
	c.Check(k.GetBalanceOfModule(ctx, types.ContractName, configs.Denom).Int64(), Equals, int64(900))
	c.Check(k.GetBalanceOfModule(ctx, types.ReserveName, configs.Denom).Int64(), Equals, int64(10))

	// increase the nonce and get slightly more funds for the provider
	ctx = ctx.WithBlockHeight(30)
	c.Assert(s.ClaimContractIncomeHandle(ctx, &msg), IsNil)
	acct := k.GetBalance(ctx, acc).AmountOf(configs.Denom).Int64()
	c.Check(acct, Equals, int64(180))
	cname := k.GetBalanceOfModule(ctx, types.ContractName, configs.Denom).Int64()
	c.Check(cname, Equals, int64(800))
	rname := k.GetBalanceOfModule(ctx, types.ReserveName, configs.Denom).Int64()
	c.Check(rname, Equals, int64(20))
	c.Check(rname+cname+acct, Equals, contract.Rate*contract.Duration)

	// ensure provider cannot take more than what is deposited into the account, overspend the contract
	ctx = ctx.WithBlockHeight(30000000)
	c.Assert(s.ClaimContractIncomeHandle(ctx, &msg), IsNil)
	acct = k.GetBalance(ctx, acc).AmountOf(configs.Denom).Int64()
	c.Check(acct, Equals, int64(900))
	cname = k.GetBalanceOfModule(ctx, types.ContractName, configs.Denom).Int64()
	c.Check(cname, Equals, int64(0))
	rname = k.GetBalanceOfModule(ctx, types.ReserveName, configs.Denom).Int64()
	c.Check(rname, Equals, int64(100))
	c.Check(rname+cname+acct, Equals, contract.Rate*contract.Duration)
}
