package consts

const (
	KeyPurposeSigning        = "SIGNING"
	KeyPurposeSigningPrivate = "PRIVATE_SIGNING"
	KeyPurposeAuth           = "AUTHENTICATION"

	KeyTypeSecp256r12018 = "EcdsaSecp256r1VerificationKey2019"
	KeyTypeRSA2018       = "RsaVerificationKey2018"
)

var SupportedKeys = []string{KeyTypeSecp256r12018, KeyTypeRSA2018}
