package converter

import (
	"IG-Parser/tree"
)

/*
Variables (de facto constants) for Web GUI.
*/

// Default example statement
var RAW_STATEMENT = "Regional Managers, on behalf of the Secretary, may review, reward, or sanction approved certified production and handling operations and accredited certifying agents for compliance with the Act or regulations in this part, under the condition that Operations were non-compliant or violated organic farming provisions and Manager has concluded investigation."

//var ANNOTATED_STATEMENT = "A,p(Regional) A(Managers), Cex(on behalf of the Secretary), D(may) I((review [AND] (reward [XOR] sanction))) Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) Cex(for compliance with the (Act or [XOR] regulations in this part)) under the condition that Cac{A(Operations) I(were (non-compliant [OR] violated)) Bdir(organic farming provisions)}."
//"A,p(Regional) A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), D(may) I(inspect and), I((review [AND] (reward [XOR] sanction))) Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) Cex(for compliance with the (Act or [XOR] regulations in this part))"
// Statement with properties
//var ANNOTATED_STATEMENT_PRIVATE_PROPERTIES
var ANNOTATED_STATEMENT = "A,p(Regional) A[role=enforcer,type=animate](Managers), Cex(on behalf of the Secretary), D[stringency=permissive](may) I[act=performance]((review [AND] (reward [XOR] sanction))) Bdir,p(approved) Bdir1,p(certified) Bdir1[role=monitored,type=animate](production [operations]) and Bdir[role=monitored,type=animate](handling operations) and Bdir2,p(accredited) Bdir2[role=monitor,type=animate](certifying agents) Cex[ctx=purpose](for compliance with the (Act or [XOR] regulations in this part)) under the condition that {Cac[state]{A[role=monitored,type=animate](Operations) I[act=violate](were (non-compliant [OR] violated)) Bdir[type=inanimate](organic farming provisions)} [AND] Cac[state]{A[role=enforcer,type=animate](Manager) I[act=terminate](has concluded) Bdir[type=activity](investigation)}}."

// Default example ID
var STATEMENT_ID = "123"

// Default dimensions for visual output
var HEIGHT = 2000
var WIDTH = 4000

// Minimum values
const MIN_HEIGHT = 100
const MIN_WIDTH = 100

// Help for raw statement field
var HELP_RAW_STMT = "This entry field is for optional use. You can paste the original statement here as a reference while encoding it in the 'Encoded Statement' field."

// Help for coded statement field
const HTML_LINEBREAK = "\n"

var HELP_CODED_STMT = "This entry field should be used to encode your institutional statement using the IG-Script notation." + HTML_LINEBREAK +
	"The basic structure of a statement is the component symbol (e.g., 'A'), immediately followed by the coded text in parentheses, e.g., 'A(certifying agent)'." + HTML_LINEBREAK +
	"Within the coded component, logical combinations of type [AND], [OR], and [XOR] are supported, e.g., 'A(Both (certifying agent [AND] inspector)) ...'. " + HTML_LINEBREAK +
	"Note the parentheses indicating the combination scope within the component; these need to be explicitly specified for every logical operator (i.e., 'A((first [AND] second))'; 'A(first [AND] second)' will lead to an error)." + HTML_LINEBREAK + HTML_LINEBREAK +
	"In addition, the notion of statement-level nesting is supported (i.e., the substitution of component content with entire statements), " +
	"e.g., 'Cac{A(certifier) I(observes) Bdir(violation)}', including the combination of nested statements, e.g., '{Cac{A(certifier) I(observes) Bdir(violation)} [AND] Cac{A(certifier) I(sanctions) Bdir(violation)}}' (note the outer braces)." + HTML_LINEBREAK +
	"Nesting is supported on all property types (as detailed below), Activation conditions (" + tree.ACTIVATION_CONDITION + "{}), Execution constraints (" + tree.EXECUTION_CONSTRAINT + "{}), " + //HTML_LINEBREAK +
	" and the Or else component (" + tree.OR_ELSE + "{})." + HTML_LINEBREAK +
	"Additional features include the use of suffices to indicate private linkages between properties and associated components (e.g., 'Bdir,p1(violating) Bdir1(citizens) as well as Bdir,p2(compliant) Bdir2(customers)')." + HTML_LINEBREAK +
	"The parser further supports the encoding of IG Logico annotations to capture semantic information associated with component values (e.g., 'A[type=animate](Officer)'). " +
	"Such annotations can be combined with suffices indicating private component relationships (e.g., 'A,p1[prop=qualitative](personal) A1[type=animate](agent)')." + HTML_LINEBREAK + HTML_LINEBREAK +
	"Supported component symbols include (with indication of optional component-level nesting):" + HTML_LINEBREAK +
	tree.ATTRIBUTES + "() --> " + tree.IGComponentSymbolNameMap[tree.ATTRIBUTES] + HTML_LINEBREAK +
	tree.ATTRIBUTES_PROPERTY + "() --> " + tree.IGComponentSymbolNameMap[tree.ATTRIBUTES_PROPERTY] + "*" + HTML_LINEBREAK +
	tree.DEONTIC + "() --> " + tree.IGComponentSymbolNameMap[tree.DEONTIC] + HTML_LINEBREAK +
	tree.AIM + "() --> " + tree.IGComponentSymbolNameMap[tree.AIM] + HTML_LINEBREAK +
	tree.DIRECT_OBJECT + "() --> " + tree.IGComponentSymbolNameMap[tree.DIRECT_OBJECT] + "*" + HTML_LINEBREAK +
	tree.DIRECT_OBJECT_PROPERTY + "() --> " + tree.IGComponentSymbolNameMap[tree.DIRECT_OBJECT_PROPERTY] + "*" + HTML_LINEBREAK +
	tree.INDIRECT_OBJECT + "() --> " + tree.IGComponentSymbolNameMap[tree.INDIRECT_OBJECT] + "*" + HTML_LINEBREAK +
	tree.INDIRECT_OBJECT_PROPERTY + "() --> " + tree.IGComponentSymbolNameMap[tree.INDIRECT_OBJECT_PROPERTY] + "*" + HTML_LINEBREAK +
	tree.ACTIVATION_CONDITION + "() --> " + tree.IGComponentSymbolNameMap[tree.ACTIVATION_CONDITION] + "*" + HTML_LINEBREAK +
	tree.EXECUTION_CONSTRAINT + "() --> " + tree.IGComponentSymbolNameMap[tree.EXECUTION_CONSTRAINT] + "*" + HTML_LINEBREAK +
	tree.CONSTITUTED_ENTITY + "() --> " + tree.IGComponentSymbolNameMap[tree.CONSTITUTED_ENTITY] + HTML_LINEBREAK +
	tree.CONSTITUTED_ENTITY_PROPERTY + "() --> " + tree.IGComponentSymbolNameMap[tree.CONSTITUTED_ENTITY_PROPERTY] + "*" + HTML_LINEBREAK +
	tree.MODAL + "() --> " + tree.IGComponentSymbolNameMap[tree.MODAL] + HTML_LINEBREAK +
	tree.CONSTITUTIVE_FUNCTION + "() --> " + tree.IGComponentSymbolNameMap[tree.CONSTITUTIVE_FUNCTION] + HTML_LINEBREAK +
	tree.CONSTITUTING_PROPERTIES + "() --> " + tree.IGComponentSymbolNameMap[tree.CONSTITUTING_PROPERTIES] + "*" + HTML_LINEBREAK +
	tree.CONSTITUTING_PROPERTIES_PROPERTY + "() --> " + tree.IGComponentSymbolNameMap[tree.CONSTITUTING_PROPERTIES_PROPERTY] + "*" + HTML_LINEBREAK +
	tree.OR_ELSE + "{} --> " + tree.IGComponentSymbolNameMap[tree.OR_ELSE] + "**" + HTML_LINEBREAK +
	"* In addition to component annotation, these components support component-level nesting, with braces scoping the nested statements (e.g., Bdir,p{ ... })." + HTML_LINEBREAK +
	"** The Or else component only allows component-level nesting (i.e., substitution by an entire statement)."

// Help for statement ID field
var HELP_STMT_ID = "This entry field should contain a statement ID (consisting of numbers and/or letters) that is the basis for generating substatement IDs."

// Help for parameter fields
var HELP_PARAMETERS = "This section includes specific customizations for the output generation, which affect the generated output. Where larger numbers of statements are encoded for analytical purposes, ensure the consistent parameterization for all generated statements."

// Help for report error field
var HELP_REPORT = "Clicking on this link should open your mail client with a pre-populated mail." + HTML_LINEBREAK +
	"Alternatively, right-click on the link, copy the e-mail address, and send a mail manually. Ensure to provide the Request ID in the subject line or body of your mail."
