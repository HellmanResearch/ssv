package storage_test

import (
	"bytes"
	"testing"

	spectypes "github.com/bloxapp/ssv-spec/types"

	"github.com/ethereum/go-ethereum/common"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"github.com/bloxapp/ssv/logging"
	"github.com/bloxapp/ssv/registry/storage"
	ssvstorage "github.com/bloxapp/ssv/storage"
	"github.com/bloxapp/ssv/storage/basedb"
	"github.com/bloxapp/ssv/utils/blskeygen"
	"github.com/bloxapp/ssv/utils/rsaencryption"
)

func TestStorage_SaveAndGetOperatorData(t *testing.T) {
	logger := logging.TestLogger(t)
	storageCollection, done := newOperatorStorageForTest(logger)
	require.NotNil(t, storageCollection)
	defer done()

	_, pk := blskeygen.GenBLSKeyPair()

	operatorData := storage.OperatorData{
		PublicKey:    pk.Serialize(),
		OwnerAddress: common.Address{},
		ID:           1,
	}

	t.Run("get non-existing operator", func(t *testing.T) {
		nonExistingOperator, found, err := storageCollection.GetOperatorData(1)
		require.NoError(t, err)
		require.Nil(t, nonExistingOperator)
		require.False(t, found)
	})

	t.Run("get non-existing operator by public key", func(t *testing.T) {
		nonExistingOperator, found, err := storageCollection.GetOperatorDataByPubKey(logger, []byte("dummyPK"))
		require.NoError(t, err)
		require.Nil(t, nonExistingOperator)
		require.False(t, found)
	})

	t.Run("create and get operator", func(t *testing.T) {
		_, err := storageCollection.SaveOperatorData(logger, &operatorData)
		require.NoError(t, err)
		operatorDataFromDB, found, err := storageCollection.GetOperatorData(operatorData.ID)
		require.NoError(t, err)
		require.True(t, found)
		require.Equal(t, operatorData.ID, operatorDataFromDB.ID)
		require.True(t, bytes.Equal(operatorData.PublicKey, operatorDataFromDB.PublicKey))
		operatorDataFromDBCmp, found, err := storageCollection.GetOperatorDataByPubKey(logger, operatorData.PublicKey)
		require.NoError(t, err)
		require.True(t, found)
		require.Equal(t, operatorDataFromDB.ID, operatorDataFromDBCmp.ID)
		require.True(t, bytes.Equal(operatorDataFromDB.PublicKey, operatorDataFromDBCmp.PublicKey))
	})

	t.Run("create existing operator", func(t *testing.T) {
		od := storage.OperatorData{
			PublicKey:    []byte("010101010101"),
			OwnerAddress: common.Address{},
			ID:           1,
		}
		_, err := storageCollection.SaveOperatorData(logger, &od)
		require.NoError(t, err)
		odDup := storage.OperatorData{
			PublicKey:    []byte("010101010101"),
			OwnerAddress: common.Address{},
			ID:           1,
		}
		_, err = storageCollection.SaveOperatorData(logger, &odDup)
		require.NoError(t, err)
		_, found, err := storageCollection.GetOperatorData(od.ID)
		require.NoError(t, err)
		require.True(t, found)
	})

	t.Run("create and get multiple operators", func(t *testing.T) {
		ods := []storage.OperatorData{
			{
				PublicKey:    []byte("01010101"),
				OwnerAddress: common.Address{},
				ID:           10,
			}, {
				PublicKey:    []byte("02020202"),
				OwnerAddress: common.Address{},
				ID:           11,
			}, {
				PublicKey:    []byte("03030303"),
				OwnerAddress: common.Address{},
				ID:           12,
			},
		}
		for _, od := range ods {
			odCopy := od
			_, err := storageCollection.SaveOperatorData(logger, &odCopy)
			require.NoError(t, err)
		}

		for _, od := range ods {
			operatorDataFromDB, found, err := storageCollection.GetOperatorData(od.ID)
			require.NoError(t, err)
			require.True(t, found)
			require.Equal(t, od.ID, operatorDataFromDB.ID)
			require.Equal(t, od.PublicKey, operatorDataFromDB.PublicKey)
		}
	})
}

func TestStorage_ListOperators(t *testing.T) {
	logger := logging.TestLogger(t)
	storageCollection, done := newOperatorStorageForTest(logger)
	require.NotNil(t, storageCollection)
	defer done()

	n := 5
	for i := 0; i < n; i++ {
		pk, _, err := rsaencryption.GenerateKeys()
		require.NoError(t, err)
		operator := storage.OperatorData{
			PublicKey: pk,
			ID:        spectypes.OperatorID(i),
		}
		_, err = storageCollection.SaveOperatorData(logger, &operator)
		require.NoError(t, err)
	}

	t.Run("successfully list operators", func(t *testing.T) {
		operators, err := storageCollection.ListOperators(logger, 0, 0)
		require.NoError(t, err)
		require.Equal(t, n, len(operators))
	})

	t.Run("successfully list operators in range", func(t *testing.T) {
		operators, err := storageCollection.ListOperators(logger, 1, 2)
		require.NoError(t, err)
		require.Equal(t, 2, len(operators))
	})
}

func newOperatorStorageForTest(logger *zap.Logger) (storage.Operators, func()) {
	db, err := ssvstorage.GetStorageFactory(logger, basedb.Options{
		Type: "badger-memory",
		Path: "",
	})
	if err != nil {
		return nil, func() {}
	}
	s := storage.NewOperatorsStorage(db, []byte("test"))
	return s, func() {
		db.Close(logger)
	}
}
