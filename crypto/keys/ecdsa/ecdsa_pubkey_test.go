package ecdsa

import (
	"encoding/hex"
	"fmt"
	"testing"

	proto "github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"

	cryptotypes "github.com/cosmos/cosmos-sdk/crypto/types"
)

type EcdsaSuite struct {
	suite.Suite
	pk cryptotypes.PubKey
	sk cryptotypes.PrivKey
}

func TestEcdsaSuite(t *testing.T) {
	suite.Run(t, new(EcdsaSuite))
}

func (suite *EcdsaSuite) SetupSuite() {
	sk, err := GenSecp256r1()
	suite.Require().NoError(err)
	suite.sk = sk
	suite.pk = sk.PubKey()
}

func (suite *EcdsaSuite) TestString() {
	assert := suite.Assert()
	require := suite.Require()

	pkStr := suite.pk.String()
	prefix := "secp256r1{"
	require.Len(pkStr, len(prefix)+PubKeySize*2+1) // prefix + hex_len + "}"
	assert.Equal(prefix, pkStr[:len(prefix)])
	assert.EqualValues('}', pkStr[len(pkStr)-1])

	bz, err := hex.DecodeString(pkStr[len(prefix) : len(pkStr)-1])
	require.NoError(err)
	assert.EqualValues(suite.pk.Bytes(), bz)
}

func (suite *EcdsaSuite) TestEquals() {
	require := suite.Require()

	skOther, err := GenSecp256r1()
	require.NoError(err)
	pkOther := skOther.PubKey()
	pkOther2 := &ecdsaPK{skOther.(ecdsaSK).PublicKey, nil}

	require.False(suite.pk.Equals(pkOther))
	require.True(pkOther.Equals(pkOther2))
	require.True(pkOther2.Equals(pkOther), "Equals must be reflexive")
}

func (suite *EcdsaSuite) TestMarshalAmino() {
	require := suite.Require()
	type AminoPubKey interface {
		cryptotypes.PubKey
		MarshalAmino() ([]byte, error)
	}

	pk := suite.pk.(AminoPubKey)
	bz, err := pk.MarshalAmino()
	require.NoError(err)

	var pk2 = new(ecdsaPK)
	require.NoError(pk2.UnmarshalAmino(bz))
	require.True(pk2.Equals(suite.pk))
}

func (suite *EcdsaSuite) TestMarshalProto() {
	require := suite.Require()
	// registry := types.NewInterfaceRegistry()
	// cdc := codec.NewProtoCodec(registry)

	bz, err := proto.Marshal(suite.pk)
	require.NoError(err)
	fmt.Println("bytes:", bz)

	// bz, err := cdc.MarshalInterface(suite)
	// require.NoError(t, err)

	// var animal testdata.Animal

	// // empty registry should fail
	// err = cdc.UnmarshalInterface(bz, &animal)
	// require.Error(t, err)

	// // wrong type registration should fail
	// registry.RegisterImplementations((*testdata.Animal)(nil), &testdata.Dog{})
	// err = cdc.UnmarshalInterface(bz, &animal)
	// require.Error(t, err)

	// // should pass
	// registry = NewTestInterfaceRegistry()
	// cdc = codec.NewProtoCodec(registry)
	// err = cdc.UnmarshalInterface(bz, &animal)
	// require.NoError(t, err)
	// require.Equal(t, kitty, animal)

	// // nil should fail
	// registry = NewTestInterfaceRegistry()
	// err = cdc.UnmarshalInterface(bz, nil)
	// require.Error(t, err)

}