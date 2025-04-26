package shared

import (
	"IG-Parser/core/tree"
)

/*
Variables (de facto constants) for Web GUI.
*/

// Link in header to IG Script overview
const HEADER_SCRIPT_LINK = "Opens an overview of of IG Script syntax (opens new tab)"

// Link in header to IG 2.0 website
const HEADER_IG_LINK = "Opens Institutional Grammar 2.0 website (opens new tab)"

// Default example statement
const RAW_STATEMENT = "Once policy comes into force, relevant regulators must monitor and enforce compliance."

// Encoded example statement, including properties and semantic annotations
const ANNOTATED_STATEMENT = "Cac{Once E(policy) F(comes into force)} A,p(relevant) A(regulators) D(must) I(monitor [AND] enforce) Bdir(compliance)."

// Default example ID
const STATEMENT_ID = "123"

// Default dimensions for visual output
const HEIGHT = 2000
const WIDTH = 4000

// Minimum values
const MIN_HEIGHT = 100
const MIN_WIDTH = 100

// Help for raw statement field
const HELP_RAW_STMT = "This entry field is for optional use. You can paste the original statement here as a reference while encoding it in the 'Encoded Statement' field."

// Help for coded statement field

// Default line break accommodating both tooltips (use \n) and HTML (specifically, for help page - for which it is substituted locally)
const LINEBREAK = "\n"

// Emphasis of syntax examples
const HTML_EM_START = "<b>"
const HTML_EM_STOP = "</b>"

// UI information for Encoded Statement field, refering to Help page
const HELP_REF = "Click to open a separate help page explaining the IG Script syntax (opens new tab)."

// Content of help page
var HELP_CODED_STMT = "The <em>'Encoded Statement'</em> field in the IG Parser UI is where the actual encoding of the institutional statement in the IG Script syntax occurs." +
	LINEBREAK + LINEBREAK +
	"The IG Script syntax knows four distinct cases to facilitate the encoding:" +
	LINEBREAK + LINEBREAK +
	HTML_EM_START + "1.) Basic component coding as well as component combinations: " + HTML_EM_STOP + LINEBREAK +
	"The basic structure of a statement is the component symbol (e.g., " + HTML_EM_START + "A" + HTML_EM_STOP + " -- see the list of all component symbols supported by the IG Parser at the bottom), immediately followed by the coded text in parentheses, e.g., " + HTML_EM_START + "A(certifying agent)" + HTML_EM_STOP + "." +
	LINEBREAK +
	"Within the coded component, logical combinations of type " + HTML_EM_START + "[AND]" + HTML_EM_STOP + ", " + HTML_EM_START + "[OR]" + HTML_EM_STOP + ", and " + HTML_EM_START + "[XOR]" + HTML_EM_STOP + " are supported, e.g., " + HTML_EM_START + "A(Both (certifying agent [AND] inspector)) ..." + HTML_EM_STOP + ". " +
	"(In this example, only 'certifying agent' and 'inspector' are part of the combination; 'Both' is not part of it.)" +
	LINEBREAK +
	"Note the use of parentheses to indicate scope within the component text. As indicated above, parentheses are relevant if the scope of combinations lies within the component (e.g., " + HTML_EM_START + "A(Both (certifying agent [AND] inspector)) ... " + HTML_EM_STOP + ") " +
	"but also for the indication of precedence in the case of nested combinations of multiple values linked by different logical operators (e.g., " + HTML_EM_START + "A(farmer [XOR] (certifying agent [AND] inspector)) ... " + HTML_EM_STOP + "). " +
	"In all other cases the parser assumes that all component content is part of the combination (i.e., if a logical operator is present), " +
	"(i.e., " + HTML_EM_START + "A(first [AND] second)" + HTML_EM_STOP + " is the same as " + HTML_EM_START + "A((first [AND] second))" + HTML_EM_STOP + ")." +
	LINEBREAK + LINEBREAK +
	"Example 1: " + HTML_EM_START + "A(Concert visitors) D(must) I(present) Bdir(tickets) to Bind(agent) at the Cex(venue)." + HTML_EM_STOP + LINEBREAK +
	" Note: Text outside any encoded component (here: " + HTML_EM_START + "to" + HTML_EM_STOP + " and " + HTML_EM_START + "at the" + HTML_EM_STOP + ") is ignored during parsing." + LINEBREAK + LINEBREAK +
	"Example 2: " + HTML_EM_START + "A(Both (concert visitors [AND] reporters)) D(must) I(present) Bdir,p(corresponding) Bdir(tickets) to Bind(agent [XOR] security personnel) at the Cex(venue)." + HTML_EM_STOP + LINEBREAK +
	" Note: This example displays variably scoped component combinations as well as the coding of properties (in this case for the direct object; see supported property symbols for other components (e.g., A,p) in the table at the bottom)." +
	LINEBREAK + LINEBREAK +
	HTML_EM_START + "2.) Nested components (\"component-level nesting\" in IG): " + HTML_EM_STOP + LINEBREAK +
	"Component-level nesting (i.e., the substitution of component content with entire statements) applies when a single component is further decomposed into individual elements, " +
	"e.g., an activation condition that captures a distinctive event such as " + HTML_EM_START + "Cac{A(certifier) I(observes) Bdir(violation)}" + HTML_EM_STOP + " (read: \"if certifier observes violation, ...\"), " +
	"including the combination of components as in the previous case (e.g., " + HTML_EM_START + "Cac{A(certifier) I((observes [AND] reports)) Bdir(violation)}" + HTML_EM_STOP + ")" +
	LINEBREAK + LINEBREAK +
	"Example: " + HTML_EM_START + "A(Agent) D(must) I(reject) Bdir(admission) Cac{A(concert visitor) I(refuses to present) Bdir(ticket)}." + HTML_EM_STOP + LINEBREAK + " Note: Nesting is supported on all property types (as detailed below), Activation conditions (" + HTML_EM_START + tree.ACTIVATION_CONDITION + "{}" + HTML_EM_STOP + "), Execution constraints (" + HTML_EM_START + tree.EXECUTION_CONSTRAINT + "{}" + HTML_EM_STOP + "), " +
	" and the Or else component (" + HTML_EM_START + tree.OR_ELSE + "{}" + HTML_EM_STOP + ")." +
	LINEBREAK + LINEBREAK +
	HTML_EM_START + "3.) Nested statement combinations: " + HTML_EM_STOP + LINEBREAK +
	"Nested statement combinations occur if multiple distinct nested components are logically linked. For instance, if two functionally distinct activation conditions apply, " +
	"such as \"certifier observes violation\" OR \"certified agent requests revocation\". Using braces (i.e., { and }), such combinations can be encoded as follows " +
	HTML_EM_START + "Cac{Cac{A(certifier) I(observes) Bdir(violation)} [OR] Cac{A,p(certified) A(agent) I(revokes) Bdir(revocation)}}" + HTML_EM_STOP +
	". Central aspect here is to indicate the distinctive component type preceding the braces (activation conditions should only be combined with activation conditions, for instance)." +
	LINEBREAK + LINEBREAK +
	"Example: " + HTML_EM_START + "A(Agent) D(must) I(reject) Bdir(admission) Cac{Cac{A(visitor) I(refuses to present) Bdir(ticket)} [OR] Cac{A(organizer) I(cancels) Bdir(event)}}." + HTML_EM_STOP + LINEBREAK +
	"Note: This applies to all components that support component-level nesting (see table at the bottom). Nested statement combinations can, similar to nested components, contain component combinations." +
	LINEBREAK + LINEBREAK +
	HTML_EM_START + "4.) Component pair combinations: " + HTML_EM_STOP + LINEBREAK +
	"Component pair combinations are similar to the previous case, but instead of applying to completely distinctive expressions, they apply in cases where some but not all parts " +
	"of the statement (pairs of components, hence \"component pairs\") are different, such as \"certifiers must review certification procedure and monitor compliance\". " +
	"This would be encoded as " + HTML_EM_START + "A(Certifiers) D(must) {I(review) Bdir(certification procedures) [AND] I(monitor) Bdir(compliance)}" + HTML_EM_STOP + ". " +
	"Note the use of braces to signal the component pairs that are distinct (here \"review certification procedures\" and \"monitor compliance\", both of which consist of a distinct aim and direct object), " +
	"but are, in this instance, executed by the same actor (here: \"certifiers\")." +
	LINEBREAK + LINEBREAK +
	"Example 1: " + HTML_EM_START + "A(Agent) D(must) {I(reject) Bdir(admission) [AND] I(report) Bdir(occurrence)}" + HTML_EM_STOP + LINEBREAK +
	"Note: Component pairs can consist of any type and number of components (e.g., " +
	HTML_EM_START + "{A(actor1) D(must) I(perform action 1) [XOR] A(actor2) D(may) I(perform action2)} Cac(Under any circumstance)" + HTML_EM_STOP + "), and also applies to nested components " +
	"(e.g., in an activation condition, such as " + HTML_EM_START + "Cac{A(actor) {I(action1) Bdir(object1) [XOR] I(action2) Bdir(object2)}}" + HTML_EM_STOP + "). " +
	"Component pairs can further embed any form of the syntactic cases introduced above (component combinations, nested statements and nested statement combinations)." + LINEBREAK + LINEBREAK +
	"Example 2: " + HTML_EM_START + "A(Agent) D(must) {I(reject) Bdir(admission) [AND] {I(report) Bdir(occurrence) [OR] I(consult) Bdir(supervisor)}}" + HTML_EM_STOP + LINEBREAK +
	"Note: This example shows the combination of multiple component pairs with explicit linkage via logical operators. On a given nesting level (e.g., on top-level statement, within nested component), " +
	"only one component pair expression (including nested linkages as shown above) is necessary to capture any number of component pair alternatives. The parser will offer a corresponding indication if multiple " +
	"separate component pairs are identified in the coded statement." +
	LINEBREAK + LINEBREAK +
	HTML_EM_START + "Additional features (Suffixes, Semantic Annotations)" + HTML_EM_STOP + LINEBREAK +
	"IG Script supports additional features specifically aimed at handling property associations and facilitating semantic annotations: " +
	LINEBREAK + LINEBREAK +
	HTML_EM_START + "Suffixes" + HTML_EM_STOP + LINEBREAK +
	"This includes the use of " + HTML_EM_START + "suffixes" + HTML_EM_STOP + " to indicate exclusive linkages between properties and associated components (e.g., " + HTML_EM_START + "Bdir1,p(violating) Bdir1(citizens)" + HTML_EM_STOP +
	" as well as " + HTML_EM_START + "Bdir2,p(compliant) Bdir2(customers)" + HTML_EM_STOP + " indicating that the properties are exclusively associated with the given corresponding object, " +
	"i.e., as \"violating citizens\" and \"compliant customers\", respectively). This principle applies to most component types and is described at greater detail in the comprehensive syntax overview (linked at the top of the page). " +
	LINEBREAK + LINEBREAK +
	HTML_EM_START + "Semantic annotations" + HTML_EM_STOP + LINEBREAK +
	"The parser further supports the encoding of " + HTML_EM_START + "semantic annotations" + HTML_EM_STOP + ", reflecting IG Logico's focus on capturing semantic information associated " +
	"with component values (e.g., " + HTML_EM_START + "A[type=animate](Officer)" + HTML_EM_STOP + "). " +
	"Such annotations apply to any component and can be combined with suffixes indicating private component relationships (e.g., " + HTML_EM_START + "A1,p[prop=qualitative](personal) A1[type=animate](agent)" + HTML_EM_STOP + ")." +
	"They can further be used to annotate nested components (e.g., " + HTML_EM_START + "Cac[event=violation]{ A(actor) I(violates) ... }" + HTML_EM_STOP + "), as well as combinations thereof " +
	"(e.g., " + HTML_EM_START + "Cac[state=condition]{Cac[event=violation]{ A(actor) I(violates) ... } [OR] Cac[event=non-compliance]{ A(actor) I(does not comply) ... }}" + HTML_EM_STOP + ")." +
	LINEBREAK +
	"Annotations can finally apply on entire statements to facilitate statement-level annotations. The syntax for this simply relies on the presence of an annotation not directly associated with a component and further " +
	"supports the presence of multiple annotations (which are concatenated in the generated output). Their positioning in the statement is arbitrary." +
	LINEBREAK + LINEBREAK +
	"Example (for multiple statement-level annotations): " + HTML_EM_START + "A(actor) [statement-level annotation] I(aim) Bdir(direct object) [yet another statement level annotation]" + HTML_EM_STOP +
	LINEBREAK + LINEBREAK +
	"Supported " + HTML_EM_START + "IG Script symbols" + HTML_EM_STOP + " for the encoding of components include (with indication of support for component-level nesting where applicable):" +
	LINEBREAK + LINEBREAK +
	"<table>" +
	"<tr><th>IG Script Symbol</th><th>Corresponding IG 2.0 Component</th></tr>" +
	"<tr><td>" + HTML_EM_START + tree.ATTRIBUTES + HTML_EM_STOP + "</td><td>" + tree.IGComponentSymbolNameMap[tree.ATTRIBUTES] + "</td></tr>" +
	"<tr><td>" + HTML_EM_START + tree.ATTRIBUTES_PROPERTY + HTML_EM_STOP + "</td><td>" + tree.IGComponentSymbolNameMap[tree.ATTRIBUTES_PROPERTY] + "*" + "</td></tr>" +
	"<tr><td>" + HTML_EM_START + tree.DEONTIC + HTML_EM_STOP + "</td><td>" + tree.IGComponentSymbolNameMap[tree.DEONTIC] + "</td></tr>" +
	"<tr><td>" + HTML_EM_START + tree.AIM + HTML_EM_STOP + "</td><td>" + tree.IGComponentSymbolNameMap[tree.AIM] + "</td></tr>" +
	"<tr><td>" + HTML_EM_START + tree.DIRECT_OBJECT + HTML_EM_STOP + "</td><td>" + tree.IGComponentSymbolNameMap[tree.DIRECT_OBJECT] + "*" + "</td></tr>" +
	"<tr><td>" + HTML_EM_START + tree.DIRECT_OBJECT_PROPERTY + HTML_EM_STOP + "</td><td>" + tree.IGComponentSymbolNameMap[tree.DIRECT_OBJECT_PROPERTY] + "*" + "</td></tr>" +
	"<tr><td>" + HTML_EM_START + tree.INDIRECT_OBJECT + HTML_EM_STOP + "</td><td>" + tree.IGComponentSymbolNameMap[tree.INDIRECT_OBJECT] + "*" + "</td></tr>" +
	"<tr><td>" + HTML_EM_START + tree.INDIRECT_OBJECT_PROPERTY + HTML_EM_STOP + "</td><td>" + tree.IGComponentSymbolNameMap[tree.INDIRECT_OBJECT_PROPERTY] + "*" + "</td></tr>" +
	"<tr><td>" + HTML_EM_START + tree.ACTIVATION_CONDITION + HTML_EM_STOP + "</td><td>" + tree.IGComponentSymbolNameMap[tree.ACTIVATION_CONDITION] + "*" + "</td></tr>" +
	"<tr><td>" + HTML_EM_START + tree.EXECUTION_CONSTRAINT + HTML_EM_STOP + "</td><td>" + tree.IGComponentSymbolNameMap[tree.EXECUTION_CONSTRAINT] + "*" + "</td></tr>" +
	"<tr><td>" + HTML_EM_START + tree.CONSTITUTED_ENTITY + HTML_EM_STOP + "</td><td>" + tree.IGComponentSymbolNameMap[tree.CONSTITUTED_ENTITY] + "</td></tr>" +
	"<tr><td>" + HTML_EM_START + tree.CONSTITUTED_ENTITY_PROPERTY + HTML_EM_STOP + "</td><td>" + tree.IGComponentSymbolNameMap[tree.CONSTITUTED_ENTITY_PROPERTY] + "*" + "</td></tr>" +
	"<tr><td>" + HTML_EM_START + tree.MODAL + HTML_EM_STOP + "</td><td>" + tree.IGComponentSymbolNameMap[tree.MODAL] + "</td></tr>" +
	"<tr><td>" + HTML_EM_START + tree.CONSTITUTIVE_FUNCTION + HTML_EM_STOP + "</td><td>" + tree.IGComponentSymbolNameMap[tree.CONSTITUTIVE_FUNCTION] + "</td></tr>" +
	"<tr><td>" + HTML_EM_START + tree.CONSTITUTING_PROPERTIES + HTML_EM_STOP + "</td><td>" + tree.IGComponentSymbolNameMap[tree.CONSTITUTING_PROPERTIES] + "*" + "</td></tr>" +
	"<tr><td>" + HTML_EM_START + tree.CONSTITUTING_PROPERTIES_PROPERTY + HTML_EM_STOP + "</td><td>" + tree.IGComponentSymbolNameMap[tree.CONSTITUTING_PROPERTIES_PROPERTY] + "*" + "</td></tr>" +
	"<tr><td>" + HTML_EM_START + tree.OR_ELSE + HTML_EM_STOP + "</td><td>" + tree.IGComponentSymbolNameMap[tree.OR_ELSE] + "**" + "</td></tr>" +
	"</table>" + LINEBREAK +
	"* In addition to component annotation, these components support component-level nesting, with braces scoping the nested statements (e.g., " + HTML_EM_START + "Bdir{ ... }" + HTML_EM_STOP + ", " + HTML_EM_START + "Bdir,p{ ... }" + HTML_EM_STOP + ", etc.)." + LINEBREAK +
	"** The Or else component only allows component-level nesting (i.e., substitution by an entire statement)."

// Help for statement ID field
const HELP_STMT_ID = "This entry field should contain a statement ID (consisting of numbers and/or letters) that is the basis for generating substatement IDs."

// Help for parameter fields
const HELP_PARAMETERS = "This section includes specific customizations for the output generation, which affect the generated output. Where larger numbers of statements are encoded for analytical purposes, ensure the consistent parameterization for all generated statements."

// Help for Original Statement output inclusion
const HELP_ORIGINAL_STATEMENT_OUTPUT = "Indicates whether the Original Statement is included in the output by introducing an additional column following the Statement ID. Choices include the exclusion (no additional column), the inclusion for the first atomic statement only (i.e., first row following the header row), or the inclusion for all atomic statements (i.e., each row)."

// Help for IG Script output inclusion
const HELP_IG_SCRIPT_OUTPUT = "Indicates whether the IG Script-encoded statement is included in the output by introducing an additional column following the Statement ID (or the Original Statement if activated). Choices include the exclusion (no additional column), the inclusion for the first atomic statement only (i.e., first row following the header row), or the inclusion for all atomic statements (i.e., each row)."

// Help for output field
const HELP_OUTPUT_TYPE = "The application currently supports two output types, either Google Sheets output, which can be directly copied into any Google sheet in your browser, or CSV format, which can be used for further processing in Excel or by scripts. Note that the CSV variant uses the pipe symbol ('|') as delimiter/separator."

// Help for report error field
const HELP_REPORT = "Clicking on this link should open your mail client with a pre-populated mail." + LINEBREAK +
	"Alternatively, right-click on the link, copy the e-mail address, and send a mail manually. Ensure to provide the Request ID in the subject line or body of your mail."
