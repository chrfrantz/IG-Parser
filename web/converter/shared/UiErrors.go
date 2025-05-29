package shared

import "strconv"

/*
Lists constants used for direct interaction with users.
*/
const ERROR_INPUT_STATEMENT_ID = "Error: The Statement ID is missing. Please review corresponding field."
const ERROR_INPUT_NO_STATEMENT = "Error: The 'Encoded Statement' field does not contain IG Script-encoded content."
const ERROR_INPUT_IGNORED_ELEMENTS = "Error: Please review the 'Encoded Statement' for the following element(s) that could not be parsed: "
const WARNING_INPUT_NON_PARSED_ELEMENTS = "Warning: The input text might have contained IG Script text fragments that have not been parsed (e.g., annotation parts, nested statements). If you believe the following text, or parts of it, should have been parsed, please review your coding accordingly (else ignore this message): "

// Made the following ones variables to allow flexible concatenation of variables.
var ERROR_INPUT_WIDTH = "The input value for output canvas width (in px) is invalid and has been reset (Minimum value: " + strconv.Itoa(MIN_WIDTH) + ")."
var ERROR_INPUT_HEIGHT = "The input value for output canvas height (in px) is invalid and has been reset (Minimum value: " + strconv.Itoa(MIN_HEIGHT) + ")."
