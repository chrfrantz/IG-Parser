package converter

/*
Struct for interacting with template via handler
 */
type ReturnStruct struct{
	// Indicates whether operation was successful
	Success bool;
	// Indicates whether an error has occurred
	Error bool;
	// Message shown to user
	Message string;
	// Original unparsed statement
	RawStmt string;
	// IG-Script annotated statement
	CodedStmt string;
	// Statement ID
	StmtId string;
	// Generated tabular output
	TabularOutput string
	// Transaction ID
	TransactionId string;
}
