# IG Parser Revisions

* Version 0.6   - Fixed inference of logical operators for component combinations when combined with nested components. Added support for tilde ('~') in input. Added convenience function for node substitution.
* Version 0.5   - Added Original Statement inclusion option for output. Fixed bug related to inclusion of IG Script output. 
* Version 0.4.1 - Added validation button to test raw statement input for imbalanced parentheses (and copies content into 'Encoded Statement' field following successful valiation).
* Version 0.4   - Added contributor information and licensing. Reorganized dependencies and code base, including file commenting. Bugfix: Prevent crash when parsing pair combinations without logical operators.
* Version 0.3   - Added option to include IG Script source in generated tabular output.
* Version 0.2.1 - Bugfix: Prevent parsing of duplicate input (e.g., duplicate components with same values/statements).
* Version 0.2   - Added support for component pair combinations in input; enabled interactive switching between visual and tabular output variants; added saving of user settings in browser storage.
* Version 0.1   - First publicly deployed version of the parser with visual and tabular output, support for container-based and local deployment.
