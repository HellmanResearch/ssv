package eth1

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/pkg/errors"
	"go.uber.org/zap"

	"github.com/bloxapp/ssv/eth1/abiparser"
)

var (
	contractABI = `[{"inputs":[],"name":"ApprovalNotWithinTimeframe","type":"error"},{"inputs":[],"name":"CallerNotOwner","type":"error"},{"inputs":[],"name":"CallerNotWhitelisted","type":"error"},{"inputs":[],"name":"ClusterAlreadyEnabled","type":"error"},{"inputs":[],"name":"ClusterDoesNotExists","type":"error"},{"inputs":[],"name":"ClusterIsLiquidated","type":"error"},{"inputs":[],"name":"ClusterNotLiquidatable","type":"error"},{"inputs":[],"name":"ExceedValidatorLimit","type":"error"},{"inputs":[],"name":"FeeExceedsIncreaseLimit","type":"error"},{"inputs":[],"name":"FeeIncreaseNotAllowed","type":"error"},{"inputs":[],"name":"FeeTooLow","type":"error"},{"inputs":[],"name":"IncorrectClusterState","type":"error"},{"inputs":[],"name":"IncorrectValidatorState","type":"error"},{"inputs":[],"name":"InsufficientBalance","type":"error"},{"inputs":[],"name":"InvalidOperatorIdsLength","type":"error"},{"inputs":[],"name":"InvalidPublicKeyLength","type":"error"},{"inputs":[],"name":"NewBlockPeriodIsBelowMinimum","type":"error"},{"inputs":[],"name":"NoFeeDeclared","type":"error"},{"inputs":[],"name":"NotAuthorized","type":"error"},{"inputs":[],"name":"OperatorAlreadyExists","type":"error"},{"inputs":[],"name":"OperatorDoesNotExist","type":"error"},{"inputs":[],"name":"OperatorsListNotUnique","type":"error"},{"inputs":[],"name":"SameFeeChangeNotAllowed","type":"error"},{"inputs":[],"name":"TargetModuleDoesNotExist","type":"error"},{"inputs":[],"name":"TokenTransferFailed","type":"error"},{"inputs":[],"name":"UnsortedOperatorsList","type":"error"},{"inputs":[],"name":"ValidatorAlreadyExists","type":"error"},{"inputs":[],"name":"ValidatorDoesNotExist","type":"error"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"address","name":"previousAdmin","type":"address"},{"indexed":false,"internalType":"address","name":"newAdmin","type":"address"}],"name":"AdminChanged","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"beacon","type":"address"}],"name":"BeaconUpgraded","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":false,"internalType":"uint64[]","name":"operatorIds","type":"uint64[]"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"},{"components":[{"internalType":"uint32","name":"validatorCount","type":"uint32"},{"internalType":"uint64","name":"networkFeeIndex","type":"uint64"},{"internalType":"uint64","name":"index","type":"uint64"},{"internalType":"bool","name":"active","type":"bool"},{"internalType":"uint256","name":"balance","type":"uint256"}],"indexed":false,"internalType":"struct ISSVNetworkCore.Cluster","name":"cluster","type":"tuple"}],"name":"ClusterDeposited","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":false,"internalType":"uint64[]","name":"operatorIds","type":"uint64[]"},{"components":[{"internalType":"uint32","name":"validatorCount","type":"uint32"},{"internalType":"uint64","name":"networkFeeIndex","type":"uint64"},{"internalType":"uint64","name":"index","type":"uint64"},{"internalType":"bool","name":"active","type":"bool"},{"internalType":"uint256","name":"balance","type":"uint256"}],"indexed":false,"internalType":"struct ISSVNetworkCore.Cluster","name":"cluster","type":"tuple"}],"name":"ClusterLiquidated","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":false,"internalType":"uint64[]","name":"operatorIds","type":"uint64[]"},{"components":[{"internalType":"uint32","name":"validatorCount","type":"uint32"},{"internalType":"uint64","name":"networkFeeIndex","type":"uint64"},{"internalType":"uint64","name":"index","type":"uint64"},{"internalType":"bool","name":"active","type":"bool"},{"internalType":"uint256","name":"balance","type":"uint256"}],"indexed":false,"internalType":"struct ISSVNetworkCore.Cluster","name":"cluster","type":"tuple"}],"name":"ClusterReactivated","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":false,"internalType":"uint64[]","name":"operatorIds","type":"uint64[]"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"},{"components":[{"internalType":"uint32","name":"validatorCount","type":"uint32"},{"internalType":"uint64","name":"networkFeeIndex","type":"uint64"},{"internalType":"uint64","name":"index","type":"uint64"},{"internalType":"bool","name":"active","type":"bool"},{"internalType":"uint256","name":"balance","type":"uint256"}],"indexed":false,"internalType":"struct ISSVNetworkCore.Cluster","name":"cluster","type":"tuple"}],"name":"ClusterWithdrawn","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint64","name":"value","type":"uint64"}],"name":"DeclareOperatorFeePeriodUpdated","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint64","name":"value","type":"uint64"}],"name":"ExecuteOperatorFeePeriodUpdated","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":false,"internalType":"address","name":"recipientAddress","type":"address"}],"name":"FeeRecipientAddressUpdated","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint8","name":"version","type":"uint8"}],"name":"Initialized","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint64","name":"value","type":"uint64"}],"name":"LiquidationThresholdPeriodUpdated","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"MinimumLiquidationCollateralUpdated","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"},{"indexed":false,"internalType":"address","name":"recipient","type":"address"}],"name":"NetworkEarningsWithdrawn","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint256","name":"oldFee","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"newFee","type":"uint256"}],"name":"NetworkFeeUpdated","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint64","name":"operatorId","type":"uint64"},{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":false,"internalType":"bytes","name":"publicKey","type":"bytes"},{"indexed":false,"internalType":"uint256","name":"fee","type":"uint256"}],"name":"OperatorAdded","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"uint64","name":"operatorId","type":"uint64"}],"name":"OperatorFeeCancellationDeclared","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"uint64","name":"operatorId","type":"uint64"},{"indexed":false,"internalType":"uint256","name":"blockNumber","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"fee","type":"uint256"}],"name":"OperatorFeeDeclared","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"uint64","name":"operatorId","type":"uint64"},{"indexed":false,"internalType":"uint256","name":"blockNumber","type":"uint256"},{"indexed":false,"internalType":"uint256","name":"fee","type":"uint256"}],"name":"OperatorFeeExecuted","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"internalType":"uint64","name":"value","type":"uint64"}],"name":"OperatorFeeIncreaseLimitUpdated","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint64","name":"operatorId","type":"uint64"}],"name":"OperatorRemoved","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"uint64","name":"operatorId","type":"uint64"},{"indexed":false,"internalType":"address","name":"whitelisted","type":"address"}],"name":"OperatorWhitelistUpdated","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":true,"internalType":"uint64","name":"operatorId","type":"uint64"},{"indexed":false,"internalType":"uint256","name":"value","type":"uint256"}],"name":"OperatorWithdrawn","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previousOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnershipTransferStarted","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"previousOwner","type":"address"},{"indexed":true,"internalType":"address","name":"newOwner","type":"address"}],"name":"OwnershipTransferred","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"implementation","type":"address"}],"name":"Upgraded","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":false,"internalType":"uint64[]","name":"operatorIds","type":"uint64[]"},{"indexed":false,"internalType":"bytes","name":"publicKey","type":"bytes"},{"indexed":false,"internalType":"bytes","name":"shares","type":"bytes"},{"components":[{"internalType":"uint32","name":"validatorCount","type":"uint32"},{"internalType":"uint64","name":"networkFeeIndex","type":"uint64"},{"internalType":"uint64","name":"index","type":"uint64"},{"internalType":"bool","name":"active","type":"bool"},{"internalType":"uint256","name":"balance","type":"uint256"}],"indexed":false,"internalType":"struct ISSVNetworkCore.Cluster","name":"cluster","type":"tuple"}],"name":"ValidatorAdded","type":"event"},{"anonymous":false,"inputs":[{"indexed":true,"internalType":"address","name":"owner","type":"address"},{"indexed":false,"internalType":"uint64[]","name":"operatorIds","type":"uint64[]"},{"indexed":false,"internalType":"bytes","name":"publicKey","type":"bytes"},{"components":[{"internalType":"uint32","name":"validatorCount","type":"uint32"},{"internalType":"uint64","name":"networkFeeIndex","type":"uint64"},{"internalType":"uint64","name":"index","type":"uint64"},{"internalType":"bool","name":"active","type":"bool"},{"internalType":"uint256","name":"balance","type":"uint256"}],"indexed":false,"internalType":"struct ISSVNetworkCore.Cluster","name":"cluster","type":"tuple"}],"name":"ValidatorRemoved","type":"event"},{"stateMutability":"nonpayable","type":"fallback"},{"inputs":[],"name":"acceptOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint64","name":"operatorId","type":"uint64"}],"name":"cancelDeclaredOperatorFee","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint64","name":"operatorId","type":"uint64"},{"internalType":"uint256","name":"fee","type":"uint256"}],"name":"declareOperatorFee","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint64[]","name":"operatorIds","type":"uint64[]"},{"internalType":"uint256","name":"amount","type":"uint256"},{"components":[{"internalType":"uint32","name":"validatorCount","type":"uint32"},{"internalType":"uint64","name":"networkFeeIndex","type":"uint64"},{"internalType":"uint64","name":"index","type":"uint64"},{"internalType":"bool","name":"active","type":"bool"},{"internalType":"uint256","name":"balance","type":"uint256"}],"internalType":"struct ISSVNetworkCore.Cluster","name":"cluster","type":"tuple"}],"name":"deposit","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint64","name":"operatorId","type":"uint64"}],"name":"executeOperatorFee","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"userAddress","type":"address"}],"name":"getRegisterAuth","outputs":[{"components":[{"internalType":"bool","name":"registerOperator","type":"bool"},{"internalType":"bool","name":"registerValidator","type":"bool"}],"internalType":"struct Authorization","name":"","type":"tuple"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"contract IERC20","name":"token_","type":"address"},{"internalType":"contract IFnSSVOperators","name":"ssvOperators_","type":"address"},{"internalType":"contract IFnSSVClusters","name":"ssvClusters_","type":"address"},{"internalType":"contract IFnSSVDAO","name":"ssvDAO_","type":"address"},{"internalType":"contract IFnSSVViews","name":"ssvViews_","type":"address"},{"internalType":"uint64","name":"minimumBlocksBeforeLiquidation_","type":"uint64"},{"internalType":"uint256","name":"minimumLiquidationCollateral_","type":"uint256"},{"internalType":"uint32","name":"validatorsPerOperatorLimit_","type":"uint32"},{"internalType":"uint64","name":"declareOperatorFeePeriod_","type":"uint64"},{"internalType":"uint64","name":"executeOperatorFeePeriod_","type":"uint64"},{"internalType":"uint64","name":"operatorMaxFeeIncrease_","type":"uint64"}],"name":"initialize","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"owner","type":"address"},{"internalType":"uint64[]","name":"operatorIds","type":"uint64[]"},{"components":[{"internalType":"uint32","name":"validatorCount","type":"uint32"},{"internalType":"uint64","name":"networkFeeIndex","type":"uint64"},{"internalType":"uint64","name":"index","type":"uint64"},{"internalType":"bool","name":"active","type":"bool"},{"internalType":"uint256","name":"balance","type":"uint256"}],"internalType":"struct ISSVNetworkCore.Cluster","name":"cluster","type":"tuple"}],"name":"liquidate","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"owner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"pendingOwner","outputs":[{"internalType":"address","name":"","type":"address"}],"stateMutability":"view","type":"function"},{"inputs":[],"name":"proxiableUUID","outputs":[{"internalType":"bytes32","name":"","type":"bytes32"}],"stateMutability":"view","type":"function"},{"inputs":[{"internalType":"uint64[]","name":"operatorIds","type":"uint64[]"},{"internalType":"uint256","name":"amount","type":"uint256"},{"components":[{"internalType":"uint32","name":"validatorCount","type":"uint32"},{"internalType":"uint64","name":"networkFeeIndex","type":"uint64"},{"internalType":"uint64","name":"index","type":"uint64"},{"internalType":"bool","name":"active","type":"bool"},{"internalType":"uint256","name":"balance","type":"uint256"}],"internalType":"struct ISSVNetworkCore.Cluster","name":"cluster","type":"tuple"}],"name":"reactivate","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint64","name":"operatorId","type":"uint64"},{"internalType":"uint256","name":"fee","type":"uint256"}],"name":"reduceOperatorFee","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes","name":"publicKey","type":"bytes"},{"internalType":"uint256","name":"fee","type":"uint256"}],"name":"registerOperator","outputs":[{"internalType":"uint64","name":"id","type":"uint64"}],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes","name":"publicKey","type":"bytes"},{"internalType":"uint64[]","name":"operatorIds","type":"uint64[]"},{"internalType":"bytes","name":"sharesData","type":"bytes"},{"internalType":"uint256","name":"amount","type":"uint256"},{"components":[{"internalType":"uint32","name":"validatorCount","type":"uint32"},{"internalType":"uint64","name":"networkFeeIndex","type":"uint64"},{"internalType":"uint64","name":"index","type":"uint64"},{"internalType":"bool","name":"active","type":"bool"},{"internalType":"uint256","name":"balance","type":"uint256"}],"internalType":"struct ISSVNetworkCore.Cluster","name":"cluster","type":"tuple"}],"name":"registerValidator","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint64","name":"operatorId","type":"uint64"}],"name":"removeOperator","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"bytes","name":"publicKey","type":"bytes"},{"internalType":"uint64[]","name":"operatorIds","type":"uint64[]"},{"components":[{"internalType":"uint32","name":"validatorCount","type":"uint32"},{"internalType":"uint64","name":"networkFeeIndex","type":"uint64"},{"internalType":"uint64","name":"index","type":"uint64"},{"internalType":"bool","name":"active","type":"bool"},{"internalType":"uint256","name":"balance","type":"uint256"}],"internalType":"struct ISSVNetworkCore.Cluster","name":"cluster","type":"tuple"}],"name":"removeValidator","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[],"name":"renounceOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"recipientAddress","type":"address"}],"name":"setFeeRecipientAddress","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint64","name":"operatorId","type":"uint64"},{"internalType":"address","name":"whitelisted","type":"address"}],"name":"setOperatorWhitelist","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"userAddress","type":"address"},{"components":[{"internalType":"bool","name":"registerOperator","type":"bool"},{"internalType":"bool","name":"registerValidator","type":"bool"}],"internalType":"struct Authorization","name":"auth","type":"tuple"}],"name":"setRegisterAuth","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"newOwner","type":"address"}],"name":"transferOwnership","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint64","name":"timeInSeconds","type":"uint64"}],"name":"updateDeclareOperatorFeePeriod","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint64","name":"timeInSeconds","type":"uint64"}],"name":"updateExecuteOperatorFeePeriod","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint64","name":"blocks","type":"uint64"}],"name":"updateLiquidationThresholdPeriod","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"updateMinimumLiquidationCollateral","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"fee","type":"uint256"}],"name":"updateNetworkFee","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint64","name":"percentage","type":"uint64"}],"name":"updateOperatorFeeIncreaseLimit","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"enum SSVModules","name":"moduleId","type":"uint8"},{"internalType":"address","name":"moduleAddress","type":"address"}],"name":"upgradeModule","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"newImplementation","type":"address"}],"name":"upgradeTo","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"address","name":"newImplementation","type":"address"},{"internalType":"bytes","name":"data","type":"bytes"}],"name":"upgradeToAndCall","outputs":[],"stateMutability":"payable","type":"function"},{"inputs":[{"internalType":"uint64[]","name":"operatorIds","type":"uint64[]"},{"internalType":"uint256","name":"amount","type":"uint256"},{"components":[{"internalType":"uint32","name":"validatorCount","type":"uint32"},{"internalType":"uint64","name":"networkFeeIndex","type":"uint64"},{"internalType":"uint64","name":"index","type":"uint64"},{"internalType":"bool","name":"active","type":"bool"},{"internalType":"uint256","name":"balance","type":"uint256"}],"internalType":"struct ISSVNetworkCore.Cluster","name":"cluster","type":"tuple"}],"name":"withdraw","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"withdrawNetworkEarnings","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint64","name":"operatorId","type":"uint64"},{"internalType":"uint256","name":"amount","type":"uint256"}],"name":"withdrawOperatorEarnings","outputs":[],"stateMutability":"nonpayable","type":"function"},{"inputs":[{"internalType":"uint64","name":"operatorId","type":"uint64"}],"name":"withdrawOperatorEarnings","outputs":[],"stateMutability":"nonpayable","type":"function"}]`
)

// Version enum to support more than one abi format
type Version int64

// Version types
const (
	V1 Version = iota
)

func (v Version) String() string {
	switch v {
	case V1:
		return "v1"
	}
	return "v1"
}

// AbiParser serves as a parsing client for events from contract
type AbiParser struct {
	Version AbiVersion
}

// NewParser return parser client based on the contract version
func NewParser(logger *zap.Logger, version Version) AbiParser {
	var parserVersion AbiVersion
	switch version {
	case V1:
		parserVersion = &abiparser.AbiV1{}
	}
	return AbiParser{Version: parserVersion}
}

// ParseOperatorAddedEvent parses an OperatorAddedEvent
func (ap AbiParser) ParseOperatorAddedEvent(log types.Log, contractAbi abi.ABI) (*abiparser.OperatorAddedEvent, error) {
	return ap.Version.ParseOperatorAddedEvent(log, contractAbi)
}

// ParseOperatorRemovedEvent parses an OperatorRemovedEvent
func (ap AbiParser) ParseOperatorRemovedEvent(log types.Log, contractAbi abi.ABI) (*abiparser.OperatorRemovedEvent, error) {
	return ap.Version.ParseOperatorRemovedEvent(log, contractAbi)
}

// ParseValidatorAddedEvent parses ValidatorAddedEvent
func (ap AbiParser) ParseValidatorAddedEvent(log types.Log, contractAbi abi.ABI) (*abiparser.ValidatorAddedEvent, error) {
	return ap.Version.ParseValidatorAddedEvent(log, contractAbi)
}

// ParseValidatorRemovedEvent parses ValidatorRemovedEvent
func (ap AbiParser) ParseValidatorRemovedEvent(log types.Log, contractAbi abi.ABI) (*abiparser.ValidatorRemovedEvent, error) {
	return ap.Version.ParseValidatorRemovedEvent(log, contractAbi)
}

// ParseClusterLiquidatedEvent parses ClusterLiquidatedEvent
func (ap AbiParser) ParseClusterLiquidatedEvent(log types.Log, contractAbi abi.ABI) (*abiparser.ClusterLiquidatedEvent, error) {
	return ap.Version.ParseClusterLiquidatedEvent(log, contractAbi)
}

// ParseClusterReactivatedEvent parses ClusterReactivatedEvent
func (ap AbiParser) ParseClusterReactivatedEvent(log types.Log, contractAbi abi.ABI) (*abiparser.ClusterReactivatedEvent, error) {
	return ap.Version.ParseClusterReactivatedEvent(log, contractAbi)
}

// ParseFeeRecipientAddressUpdatedEvent parses FeeRecipientAddressUpdatedEvent
func (ap AbiParser) ParseFeeRecipientAddressUpdatedEvent(log types.Log, contractAbi abi.ABI) (*abiparser.FeeRecipientAddressUpdatedEvent, error) {
	return ap.Version.ParseFeeRecipientAddressUpdatedEvent(log, contractAbi)
}

// AbiVersion serves as the parser client interface
type AbiVersion interface {
	ParseOperatorAddedEvent(log types.Log, contractAbi abi.ABI) (*abiparser.OperatorAddedEvent, error)
	ParseOperatorRemovedEvent(log types.Log, contractAbi abi.ABI) (*abiparser.OperatorRemovedEvent, error)
	ParseValidatorAddedEvent(log types.Log, contractAbi abi.ABI) (*abiparser.ValidatorAddedEvent, error)
	ParseValidatorRemovedEvent(log types.Log, contractAbi abi.ABI) (*abiparser.ValidatorRemovedEvent, error)
	ParseClusterLiquidatedEvent(log types.Log, contractAbi abi.ABI) (*abiparser.ClusterLiquidatedEvent, error)
	ParseClusterReactivatedEvent(log types.Log, contractAbi abi.ABI) (*abiparser.ClusterReactivatedEvent, error)
	ParseFeeRecipientAddressUpdatedEvent(log types.Log, contractAbi abi.ABI) (*abiparser.FeeRecipientAddressUpdatedEvent, error)
}

// LoadABI enables to load a custom abi json
func LoadABI(logger *zap.Logger, abiFilePath string) error {
	jsonFile, err := os.Open(filepath.Clean(abiFilePath))
	if err != nil {
		return errors.Wrap(err, "failed to open abi")
	}

	raw, err := io.ReadAll(jsonFile)
	if err != nil {
		return errors.Wrap(err, "failed to read abi")
	}
	if err := jsonFile.Close(); err != nil {
		logger.Warn("failed to close abi json", zap.Error(err))
	}
	s := string(raw)

	if err := jsonFile.Close(); err != nil {
		logger.Warn("failed to close abi json", zap.Error(err))
	}

	// assert valid JSON
	var obj []interface{}
	err = json.Unmarshal(raw, &obj)
	if err != nil {
		return errors.Wrap(err, "abi is not a valid json")
	}
	contractABI = s
	return nil
}

// ContractABI abi of the ssv-network contract
func ContractABI(abiVersion Version) string {
	switch abiVersion {
	case V1:
		return contractABI
	default:
		return contractABI
	}
}
