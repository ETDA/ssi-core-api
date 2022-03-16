package middlewares

//func TestValidateSignatureRegisterOperation(t *testing.T) {
//	str := utils.JSONToString(&services.TXBroadcastPayload{
//		Message:   "eyJjdXJyZW50X2tleSI6ICItLS0tLUJFR0lOIFBVQkxJQyBLRVktLS0tLVxuTUZrd0V3WUhLb1pJemowQ0FRWUlLb1pJemowREFRY0RRZ0FFdGVZQVFsbW42dk00Q2lFKzdRMnRBQm5TNHBuWFxuaHVxNEhNR0dFRDdpWGliOXpCMXI5YUE2VHZUbVc2VEk4ajBNMnYyaEN1R1JqdTVPWkl6MENWcHRRUT09XG4tLS0tLUVORCBQVUJMSUMgS0VZLS0tLS0iLCAibmV4dF9rZXlfaGFzaCI6ICJhYThmY2VhNzAyOWMzYmU3OWJlMWUyM2JjOTc1OGQwNTAyODlmZjRhNTNjMmQxNmNjNWI0YjI4ZTliYWE0MTQ2IiwgIm9wZXJhdGlvbiI6ICJESURfUkVHSVNURVIifQ==",
//		Signature: "MEQCIAcWV7PaCMc02m05+iFjLeKkocNow/kLSEKXu5DMdc2XAiAacL0h6lSwaQn7jP4IhiFHAGUu4NMHjjJo0vkgOLCIeA==",
//		Version:   consts.BlockChainVersion100,
//		CreatedAt: utils.GetCurrentDateTime(),
//	})
//
//	assert.Nil(t, ValidateSignatureMessage(core.NewABCIContext(nil), utils.StringToBytes(str)))
//
//}
