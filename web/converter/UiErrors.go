package converter

import "strconv"

/*
Lists constants used for direct interaction with users.
*/
const ERROR_INPUT_STATEMENT_ID = "The Statement ID is missing. Please review corresponding field."
const ERROR_INPUT_NO_STATEMENT = "The 'Encoded Statement' field does not contain IG Script-encoded content."

// Made the following ones variables to allow flexible concatenation of variables.
var ERROR_INPUT_WIDTH = "The input value for output canvas width (in px) is invalid and has been reset (Minimum value: " + strconv.Itoa(MIN_WIDTH) + ")."
var ERROR_INPUT_HEIGHT = "The input value for output canvas height (in px) is invalid and has been reset (Minimum value: " + strconv.Itoa(MIN_HEIGHT) + ")."
