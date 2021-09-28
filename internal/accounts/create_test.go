package accounts

import (
	"encoding/json"
	"testing"

	"github.com/onflow/flow-go-sdk"
	"github.com/onflow/flow-go-sdk/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"github.com/onflow/flow-cli/internal/command"
	"github.com/onflow/flow-cli/pkg/flowkit"
	"github.com/onflow/flow-cli/pkg/flowkit/output"
	"github.com/onflow/flow-cli/pkg/flowkit/services"
	"github.com/onflow/flow-cli/tests"
)

func setup() (*services.Services, *flowkit.State, *tests.TestGateway) {
	readerWriter := tests.ReaderWriter()
	state, err := flowkit.Init(readerWriter, crypto.ECDSA_P256, crypto.SHA3_256)
	if err != nil {
		panic(err)
	}

	gw := tests.DefaultMockGateway()
	s := services.NewServices(gw.Mock, state, output.NewStdoutLogger(output.NoneLog))
	return s, state, gw
}

func TestAccountsCreate(t *testing.T) {
	services, state, gw := setup()

	newAddress := flow.HexToAddress("192440c99cb17282")
	gw.GetAccount.Run(func(args mock.Arguments) {
		gw.GetAccount.Return(
			tests.NewAccountWithAddress(newAddress.String()), nil,
		)
	})

	gw.GetTransactionResult.Return(
		tests.NewAccountCreateResult(newAddress), nil,
	)

	res, err := create(nil, nil, command.GlobalFlags{}, services, state)

	assert.NoError(t, err)

	out := res.JSON()

	// cycle json
	bytes, err := json.Marshal(out)
	assert.NoError(t, err)
	account := struct {
		Address string `json:"address"`
	}{}
	err = json.Unmarshal(bytes, &account)
	assert.NoError(t, err)

	// confirm the address is the one we created in our mock gateway
	assert.Equal(t, account.Address, newAddress.Hex())
}
