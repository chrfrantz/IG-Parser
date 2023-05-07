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
const RAW_STATEMENT = "Regional Managers, on behalf of the Secretary, may review, reward, or sanction approved certified production and handling operations and accredited certifying agents for compliance with the Act or regulations in this part, under the condition that Operations were non-compliant or violated organic farming provisions and Manager has concluded investigation."

// Encoded example statement, including properties and semantic annotations
const ANNOTATED_STATEMENT = "A,p(Regional) A[role=enforcer,type=animate](Managers), Cex(on behalf of the Secretary), D[stringency=permissive](may) I[act=performance]((review [AND] (reward [XOR] sanction))) Bdir,p(approved) Bdir1,p(certified) Bdir1[role=monitored,type=animate](production [operations]) and Bdir[role=monitored,type=animate](handling operations) and Bdir2,p(accredited) Bdir2[role=monitor,type=animate](certifying agents) Cex[ctx=purpose](for compliance with the (Act or [XOR] regulations in this part)) under the condition that Cac{Cac[state]{A[role=monitored,type=animate](Operations) I[act=violate]((were non-compliant [OR] violated)) Bdir[type=inanimate](organic farming provisions)} [AND] Cac[state]{A[role=enforcer,type=animate](Manager) I[act=terminate](has concluded) Bdir[type=activity](investigation)}}."

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
var HELP_CODED_STMT = "The <em>'Encoded Statement'</em> field in the IG Parser UI is where the actual encoding of the institutional statement in the IG Script syntax occurs." + LINEBREAK + LINEBREAK +
	"The basic structure of a statement is the component symbol (e.g., " + HTML_EM_START + "A" + HTML_EM_STOP + "), immediately followed by the coded text in parentheses, e.g., " + HTML_EM_START + "A(certifying agent)" + HTML_EM_STOP + "." + LINEBREAK +
	"Within the coded component, logical combinations of type " + HTML_EM_START + "[AND]" + HTML_EM_STOP + ", " + HTML_EM_START + "[OR]" + HTML_EM_STOP + ", and " + HTML_EM_START + "[XOR]" + HTML_EM_STOP + " are supported, e.g., " + HTML_EM_START + "A(Both (certifying agent [AND] inspector)) ..." + HTML_EM_STOP + ". " + LINEBREAK +
	"Note the parentheses that indicate the combination scope within the component. These need to be explicitly specified for every logical operator (i.e., " + HTML_EM_START + "A((first [AND] second))" + HTML_EM_STOP + "; " + HTML_EM_START + "A(first [AND] second)" + HTML_EM_STOP + " will lead to an error)." + LINEBREAK + LINEBREAK +
	"In addition, the notion of statement-level nesting is supported (i.e., the substitution of component content with entire statements), " +
	"e.g., " + HTML_EM_START + "Cac{A(certifier) I(observes) Bdir(violation)}" + HTML_EM_STOP + ", including the combination of nested statements, e.g., " + HTML_EM_START + "{Cac{A(certifier) I(observes) Bdir(violation)} [AND] Cac{A(certifier) I(sanctions) Bdir(violation)}}" + HTML_EM_STOP + " (note the outer braces)." + LINEBREAK +
	"Nesting is supported on all property types (as detailed below), Activation conditions (" + HTML_EM_START + tree.ACTIVATION_CONDITION + "{}" + HTML_EM_STOP + "), Execution constraints (" + HTML_EM_START + tree.EXECUTION_CONSTRAINT + "{}" + HTML_EM_STOP + "), " +
	" and the Or else component (" + HTML_EM_START + tree.OR_ELSE + "{}" + HTML_EM_STOP + ")." + LINEBREAK + LINEBREAK +
	"Additional features include the use of suffices to indicate exclusive linkages between properties and associated components (e.g., " + HTML_EM_START + "Bdir1,p(violating) Bdir1(citizens)" + HTML_EM_STOP + " as well as " + HTML_EM_START + "Bdir2,p(compliant) Bdir2(customers)" + HTML_EM_STOP + " indicating that the properties are exclusively associated with the given corresponding object, i.e., as \"violating citizens\" and \"compliant customers\", respectively). This principle applies to most component types and is described in the comprehensive syntax overview (linked at the top of the page). The ability to use suffices on properties is further supported, and likewise described in the comprehensive overview." + LINEBREAK + LINEBREAK +
	"The parser further supports the encoding of IG Logico annotations to capture semantic information associated with component values (e.g., " + HTML_EM_START + "A[type=animate](Officer)" + HTML_EM_STOP + "). " +
	"Such annotations can be combined with suffices indicating private component relationships (e.g., " + HTML_EM_START + "A1,p[prop=qualitative](personal) A1[type=animate](agent)" + HTML_EM_STOP + ")." + LINEBREAK + LINEBREAK +
	"Supported component symbols include (with indication of optional component-level nesting):" + LINEBREAK + LINEBREAK +
	"<table>" +
	"<tr><th>IG Script Symbol</th><th>Corresponding IG 2.0 Component</th></tr>" +
	"<tr><td>" + tree.ATTRIBUTES + "()" + "</td><td>" + tree.IGComponentSymbolNameMap[tree.ATTRIBUTES] + "</td></tr>" +
	"<tr><td>" + tree.ATTRIBUTES_PROPERTY + "()</td><td>" + tree.IGComponentSymbolNameMap[tree.ATTRIBUTES_PROPERTY] + "*" + "</td></tr>" +
	"<tr><td>" + tree.DEONTIC + "()</td><td>" + tree.IGComponentSymbolNameMap[tree.DEONTIC] + "</td></tr>" +
	"<tr><td>" + tree.AIM + "()</td><td>" + tree.IGComponentSymbolNameMap[tree.AIM] + "</td></tr>" +
	"<tr><td>" + tree.DIRECT_OBJECT + "()</td><td>" + tree.IGComponentSymbolNameMap[tree.DIRECT_OBJECT] + "*" + "</td></tr>" +
	"<tr><td>" + tree.DIRECT_OBJECT_PROPERTY + "()</td><td>" + tree.IGComponentSymbolNameMap[tree.DIRECT_OBJECT_PROPERTY] + "*" + "</td></tr>" +
	"<tr><td>" + tree.INDIRECT_OBJECT + "()</td><td>" + tree.IGComponentSymbolNameMap[tree.INDIRECT_OBJECT] + "*" + "</td></tr>" +
	"<tr><td>" + tree.INDIRECT_OBJECT_PROPERTY + "()</td><td>" + tree.IGComponentSymbolNameMap[tree.INDIRECT_OBJECT_PROPERTY] + "*" + "</td></tr>" +
	"<tr><td>" + tree.ACTIVATION_CONDITION + "()</td><td>" + tree.IGComponentSymbolNameMap[tree.ACTIVATION_CONDITION] + "*" + "</td></tr>" +
	"<tr><td>" + tree.EXECUTION_CONSTRAINT + "()</td><td>" + tree.IGComponentSymbolNameMap[tree.EXECUTION_CONSTRAINT] + "*" + "</td></tr>" +
	"<tr><td>" + tree.CONSTITUTED_ENTITY + "()</td><td>" + tree.IGComponentSymbolNameMap[tree.CONSTITUTED_ENTITY] + "</td></tr>" +
	"<tr><td>" + tree.CONSTITUTED_ENTITY_PROPERTY + "()</td><td>" + tree.IGComponentSymbolNameMap[tree.CONSTITUTED_ENTITY_PROPERTY] + "*" + "</td></tr>" +
	"<tr><td>" + tree.MODAL + "()</td><td>" + tree.IGComponentSymbolNameMap[tree.MODAL] + "</td></tr>" +
	"<tr><td>" + tree.CONSTITUTIVE_FUNCTION + "()</td><td>" + tree.IGComponentSymbolNameMap[tree.CONSTITUTIVE_FUNCTION] + "</td></tr>" +
	"<tr><td>" + tree.CONSTITUTING_PROPERTIES + "()</td><td>" + tree.IGComponentSymbolNameMap[tree.CONSTITUTING_PROPERTIES] + "*" + "</td></tr>" +
	"<tr><td>" + tree.CONSTITUTING_PROPERTIES_PROPERTY + "()</td><td>" + tree.IGComponentSymbolNameMap[tree.CONSTITUTING_PROPERTIES_PROPERTY] + "*" + "</td></tr>" +
	"<tr><td>" + tree.OR_ELSE + "{}</td><td>" + tree.IGComponentSymbolNameMap[tree.OR_ELSE] + "**" + "</td></tr>" +
	"</table>" + LINEBREAK +
	"* In addition to component annotation, these components support component-level nesting, with braces scoping the nested statements (e.g., " + HTML_EM_START + "Bdir{ ... }" + HTML_EM_STOP + ", " + HTML_EM_START + "Bdir,p{ ... }" + HTML_EM_STOP + ", etc.)." + LINEBREAK +
	"** The Or else component only allows component-level nesting (i.e., substitution by an entire statement)."

// Help for statement ID field
const HELP_STMT_ID = "This entry field should contain a statement ID (consisting of numbers and/or letters) that is the basis for generating substatement IDs."

// Help for parameter fields
const HELP_PARAMETERS = "This section includes specific customizations for the output generation, which affect the generated output. Where larger numbers of statements are encoded for analytical purposes, ensure the consistent parameterization for all generated statements."

// Help for output field
const HELP_OUTPUT_TYPE = "The application currently supports two output types, either Google Sheets output, which can be directly copied into any Google sheet in your browser, or CSV format, which can be used for further processing in Excel or by scripts. Note that the CSV variant uses the pipe symbol ('|') as delimiter/separator."

// Help for report error field
const HELP_REPORT = "Clicking on this link should open your mail client with a pre-populated mail." + LINEBREAK +
	"Alternatively, right-click on the link, copy the e-mail address, and send a mail manually. Ensure to provide the Request ID in the subject line or body of your mail."
