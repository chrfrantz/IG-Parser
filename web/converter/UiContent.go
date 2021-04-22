package converter

import "IG-Parser/tree"

/*
This file holds variables (de facto constants) for Web GUI.
 */

// Default example statement
var ANNOTATED_STATEMENT = "A,p(Regional) A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), D(may) I(inspect and), I((review [AND] (reward [XOR] sanction))) Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) Cex(for compliance with the (Act or [XOR] regulations in this part)) Cac{A(Operations) I((were non-compliant [OR] violated)) Bdir(organic farming provisions)}."

		//"A,p(Regional) A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), D(may) I(inspect and), I((review [AND] (reward [XOR] sanction))) Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) Cex(for compliance with the (Act or [XOR] regulations in this part))"
// Default example ID
var STATEMENT_ID = "650"
// Help for raw statement field
var HELP_RAW_STMT = "This entry field is for optional use. You can paste the original statement here to maintain a reference, while reconstructing it in the 'Annotated Statement' field."
// Help for coded statement field
const HTML_LINEBREAK = "\n"
var HELP_CODED_STMT = "This entry field should be used to annotate your institutional statement using the IG-Script notation." + HTML_LINEBREAK +
						"The basic structure of a statement is the component symbol (e.g., 'A'), immediately followed by the coded text in parentheses, e.g., 'A(certifying agent)'." + HTML_LINEBREAK +
						"Within the coded component, logical combinations of type [AND], [OR], and [XOR] are supported, e.g., 'A(Both (certifying agent [AND] inspector)) ...'. " + HTML_LINEBREAK +
						"Note the parentheses indicating the combination scope within the component; these need to be explicitly specified for every logical operator (i.e., 'A((first [AND] second))'; 'A(first [AND] second)' will lead to an error)." + HTML_LINEBREAK + HTML_LINEBREAK +
						"In addition, the notion of statement-level nesting is supported (i.e., the substitution of component content with entire statements), " +
						"e.g., 'Cac{A(certifier) I(observes) Bdir(violation)}'" + HTML_LINEBREAK +
	                    "Nesting is supported on all property types (e.g., Attribute Property (" + tree.ATTRIBUTES_PROPERTY + "{})), as well as Activation condition (" + tree.ACTIVATION_CONDITION + "{})" + HTML_LINEBREAK +
	                    ", Execution constraints (" + tree.EXECUTION_CONSTRAINT +"{}) and Or else (" + tree.OR_ELSE + "{})" + HTML_LINEBREAK +
						"Supported component symbols include:" + HTML_LINEBREAK +
						tree.ATTRIBUTES + "() --> " + tree.IGComponentSymbolNameMap[tree.ATTRIBUTES] + HTML_LINEBREAK +
						tree.ATTRIBUTES_PROPERTY + "() --> " + tree.IGComponentSymbolNameMap[tree.ATTRIBUTES_PROPERTY] + HTML_LINEBREAK +
						tree.DEONTIC + "() --> " + tree.IGComponentSymbolNameMap[tree.DEONTIC] + HTML_LINEBREAK +
						tree.AIM + "() --> " + tree.IGComponentSymbolNameMap[tree.AIM] + HTML_LINEBREAK +
						tree.DIRECT_OBJECT + "() --> " + tree.IGComponentSymbolNameMap[tree.DIRECT_OBJECT] + HTML_LINEBREAK +
						tree.DIRECT_OBJECT_PROPERTY + "() --> " + tree.IGComponentSymbolNameMap[tree.DIRECT_OBJECT_PROPERTY] + HTML_LINEBREAK +
						tree.INDIRECT_OBJECT + "() --> " + tree.IGComponentSymbolNameMap[tree.INDIRECT_OBJECT] + HTML_LINEBREAK +
						tree.INDIRECT_OBJECT_PROPERTY + "() --> " + tree.IGComponentSymbolNameMap[tree.INDIRECT_OBJECT_PROPERTY] + HTML_LINEBREAK +
						tree.ACTIVATION_CONDITION + "() --> " + tree.IGComponentSymbolNameMap[tree.ACTIVATION_CONDITION] + HTML_LINEBREAK +
						tree.EXECUTION_CONSTRAINT + "() --> " + tree.IGComponentSymbolNameMap[tree.EXECUTION_CONSTRAINT] + HTML_LINEBREAK +
						tree.CONSTITUTED_ENTITY + "() --> " + tree.IGComponentSymbolNameMap[tree.CONSTITUTED_ENTITY] + HTML_LINEBREAK +
						tree.CONSTITUTED_ENTITY_PROPERTY + "() --> " + tree.IGComponentSymbolNameMap[tree.CONSTITUTED_ENTITY_PROPERTY] + HTML_LINEBREAK +
						tree.MODAL + "() --> " + tree.IGComponentSymbolNameMap[tree.MODAL] + HTML_LINEBREAK +
						tree.CONSTITUTIVE_FUNCTION + "() --> " + tree.IGComponentSymbolNameMap[tree.CONSTITUTIVE_FUNCTION] + HTML_LINEBREAK +
						tree.CONSTITUTING_PROPERTIES + "() --> " + tree.IGComponentSymbolNameMap[tree.CONSTITUTING_PROPERTIES] + HTML_LINEBREAK +
						tree.CONSTITUTING_PROPERTIES_PROPERTY + "() --> " + tree.IGComponentSymbolNameMap[tree.CONSTITUTING_PROPERTIES_PROPERTY] + HTML_LINEBREAK +
						tree.OR_ELSE + "{} --> " + tree.IGComponentSymbolNameMap[tree.OR_ELSE] + HTML_LINEBREAK

// Help for statement ID field
var HELP_STMT_ID = "This entry field should contain a numeric ID that is the basis for generating substatement IDs."

// Help for report error field
var HELP_REPORT = "Clicking on this link should open your mail client with a pre-populated mail." + HTML_LINEBREAK +
	"Alternatively, right-click on the link, copy the e-mail address, and send a mail manually. Ensure to provide the Request ID in the subject line or body of your mail."