package controller

// Contract|hash,signatures
func GenerateSnapshot(c Contract) string {
	serial := c.marshalTransaction() + "|"
	serial += c.hashTransaction()
	serial += "," + c.marshalSignatures()
	return serial
}
