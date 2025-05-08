# IG Parser
Parser for IG 2.0 Statements based on the IG Script Notation. 

Contact: Christopher Frantz (christopher.frantz@ntnu.no)

Institutional Grammar 2.0 Website: https://newinstitutionalgrammar.org

Deployed IG Parser: 
* Tabular output: https://ig-parser.newinstitutionalgrammar.org 
* Visual output: https://ig-parser.newinstitutionalgrammar.org/visual

Note: Either version allows interactive switching to the respective other while preserving encoded statement information.

See [Revision history](changelog.md) for a detailed overview of changes.

See [Contributors](contributors.md) for an overview of contributions to the project. We explicitly encourage external contributions. Please feel free to get in touch if you plan to contribute to the repository. Please create an issue to report bugs, or to propose features (alternatively also per mail). 

## Overview

IG Parser is a parser for IG Script, a formal notation for the representation institutional statements (e.g., policy statements) used in the Institutional Grammar 2.0. The parser can be used locally, as well as via a web interface that produces tabular output of parsed statements (currently supporting Google Sheets format). In the following, you will find a brief introduction to the user interface, followed by a comprehensive introduction to the syntactic principles and essential features of IG Script. This includes a set of examples showcasing all features and various levels of complexity, while highlighting typical mistakes in the encoding. As a final aspect, the deployment instructions for IG Parser are provided. 

The conceptual background of the Institutional Grammar 2.0 is provided in the corresponding [article](https://doi.org/10.1111/padm.12719) and [book](https://newinstitutionalgrammar.org), augmented with supplementary operational [coding guidelines](https://arxiv.org/abs/2008.08937).

## User Interface Guide

The user interface consists of various entry fields, followed by parameters (specific to each version of the parser) and an output section that will contain the generated output. The following subsections highlight the key features of each element.

### General Entry Fields

The initial entry field holds the 'Original Statement'. The purpose of this field is to keep track of the original statement during coding, but also to include it in the output if you choose to do so (see 'Parameters' section in the UI; discussed below). Prior to coding, it also allows you to vet the statement for obvious challenges in coding (e.g., imbalanced parentheses, non-supported symbol), which you can do by clicking on 'Validate 'Original Statement' input'. This will allow you copy the content to the 'Encoded Statement' field (it warns you if encoded code exists to prevent accidental overwriting of the existing coding).

The 'Encoded Statement' area provides the actual coding editor. It features both an 'Advanced Mode' that includes additional input options (see below) as well as color-coding of the encoded statement. Alternatively, you can use the standard mode that does not provide any advanced features beyond basic bracket matching, and operates by directly encoding in the IG Script syntax (described later). This version is more useful if using mobile devices or facing accessibility issues based on the color-coding.

### Advanced Editor
The advanced user interface has four main input options available to annotate statements:
* Firstly the option of simply writing out the symbols manually in the encoded statement field to annotate the text. This means manually typing each symbol and creating the brackets for each symbol. 
* The second option is to highlight the section or span you want to annotate with a symbol and then click on the button for that specific symbol below the editor. This will then encapsulate that section with the specified symbol.
* Thirdly there is a "Selection Mode" toggle above the editor, which reverses this interaction flow. So instead of first highlighting and then selecting a symbol, you first select a symbol and can then highlight text to annotate the selection with the symbol. The symbol stays selected until either a new symbol is selected or the selection mode is toggled off. The selection mode is only available on devices using a mouse as an input device. 
* Finally, there is the "edit mode". By clicking the "Text mode" button the editor switches to edit mode which disables typing in the editor in exchange for allowing keybinds to be used for annotation. This mode works by first selecting the text you want to annotate, and then clicking a button corresponding with the symbol you want to annotate the selection with. For example for a direct object the key is "r". To add a property instead for symbols which support it the same key is used, in addition to the "shift" key. So, in this instance the keybinding for a direct object property is "shift+r".

In addition to these features there are toggles for nesting, which enables the correct brackets for nested symbols, semantic annotation which adds the semantic annotation brackets to symbols on creation. There are also buttons for undo, redo and for displaying a quick guide of the supported symbols and keybinds in the editor. 

Another feature of the website is keyboard interaction. Where each element can be navigated using the "tab" key or "shift+tab" to navigate in reverse. Additionally, each button can be pressed using the "enter" key and the parameter toggles can be clicked using the "space" key. Further for selecting or highlighting text in the editor the combination of "ctrl+shift" and the arrow keys can be used. This makes the website accessible using a keyboard as the only input.

### Toggling between variants (editors, parser versions)

You can interactively toggle between basic and the advanced mode by clicking on 'Toggle advanced editor features'. You can further interactively switch between the tabular output mode as well as the visual output mode of the parser. In all those cases no coding information is lost.

### Output-specific parameters

Below the editor area you will find output-specific parameters, all of which have a simple help (by hovering over their label; some allow for opening of external pages for more extensive guidance).

#### Tabular output

* For the tabular parser version, the parameters include
  * the 'Statement ID' (which is the ID based on which individual statement IDs are generated),
  * options to generate output based on the different levels of expressiveness (IG Core, IG Extended, IG Logico),
  * the selective inclusion of a header row in the generated output,
  * the option to include the original statement (from the field 'Original Statement') as well as the IG-Script-encoded statement (Field 'Encoded Statement') in the output (either only in the first content line, or for all generated output lines),
  * the selection of the output format, which currently includes Google Sheets-parseable output (you can paste it directly into Google Sheets spreadsheets), or as CSV (which can be used in statistical programming tools, or in Excel, etc.)
 
By clicking on 'Generate tabular output', the input is parsed and output generated, which can be copied into the clipboard for transferral into a tool of your choice.

##### Note on Google Sheets 

When using Google Sheets output, pasting this into Google Sheets, and exporting the result to Excel, please note that seemingly empty cells actually contain the function `=IFERROR(@__xludf.DUMMYFUNCTION("""COMPUTED_VALUE""");" ")`, which is automatically generated when exporting xlsx files from Google Sheets.  To remove these contents (in order to arrive at empty cells), please use the Find & Replace feature in Excel (`Ctrl + F`): Put the function in the `Search` field, leave the `Replace` field empty, and click `Replace all`. Afterward, all seemingly empty cells should be actually empty. Please note that this aspect is specific to Google Sheets output; it does not apply to CSV output.

#### Visual output

* For the visual parser version, the parameters include
  * the inclusion of 'IG Logico' annotations in the output
  * the inclusion of the Degree of Variability (a metric introduced as part of the IG 2.0 - see the conceptual guidance),
  * the inclusion of component properties as tree nodes attached to their parent components (as opposed to just labels),
  * the display of a fully decomposed binary tree structure (this is useful for "debugging" your interpretation of the tree structure)
  * the choice to print activation conditions on top of the tree (to make statements more readable by reflecting the logical precedence of activation conditions over the rest of the statement),
  * the specification of the canvas dimensions (this is useful to customize the size to your needs and graph scale -- a future feature will be to automate the scaling)
 
By clicking on 'Generate visual output', the input information is parsed and the statement tree structure displayed.

### Usage considerations

* To support efficient coding, specifically for complex statements it is often useful to encode and evaluate those in visual mode, before generating the tabular output for downstream processing. Use the interactive switching features for this purpose.
* As indicated before, when hovering of a label on the UI, a short description is displayed highlighting its key function (or variably opens a new tab, where more support is provided, e.g., for the 'Encoded Statement').
* Note that the UI stores the last entry in *your* browser. When reopening your browser, it should hence display the statement you have been working on (not the history). This feature is useful if internet connection is lost, or if you work with interruptions (so you can continue at a later stage), or if your browser crashes. No information is stored on the server side, but only in your specific browser (not across browsers). To delete this information, you can delete your browser cache.

In the following, you will find an overview of the actual IG Script syntax used for the encoding of institutional statements entered into the 'Encoded Statement' field.

## IG Script

IG Script is a notation introduced in the context of the [Institutional Grammar 2.0](https://newinstitutionalgrammar.org) (IG 2.0) that aims at a deep structural representation of legal statements alongside selected levels of expressiveness. While IG 2.0 highlights the conceptual background, the objective of IG Script is to provide an accessible, but formal approach to provide a format-independent representation of institutional statements of any type (e.g., regulative, constitutive, hybrid). While the parser currently supports exemplary export formats (e.g., tabular format and visual output), the tool is open to be extended to support other output formats (e.g.,  XML, JSON, YAML). The introduction below focuses on the operational coding. Syntactic and semantic foundations are provided [elsewhere](https://github.com/InstitutionalGrammar/IG-2.0-Resources).

### Principles of IG Script Syntax

IG Script centers around a set of fundamental primitives that can be combined to parse statements comprehensively, including:

* Basic component coding
* Component combinations
* Nested statements
* Nested statement combinations
* Component pair combinations
* Object-Property relationships
* Semantic Annotations

#### Component Coding

*Component Coding* provides the basic building block for any statement encoding. A component is represented as 

* `componentSymbol(naturalLanguageText)`
  
`componentSymbol` is one of the supported symbols for the different component types (see below). 
An example for an annotated Attributes component is `A(Farmer)`.

`naturalLanguageText` is the human-readable text annotated as the entity the annotation describes. The text is
open-ended and can include special symbols, including parentheses (e.g., `A(Farmer (e.g., organic farmer))`). 
Exceptions to this rule are discussed in the context of combinations).

The scope of a component is specified by opening and closing parentheses.

All components of a statement are annotated correspondingly, without concern for order, or repetition. 
The parser further tolerates multiple component annotations of the same kind. Multiple Attributes, 
for example (e.g., `A(Farmer) D(must) I(comply) A(Certifier)`), are effectively interpreted as a combination (i.e., 
`A(Farmer [AND] Certifier) D(must) I(comply)`) in the parsing process.  

Any symbols outside the encoded components and combinations of components are ignored by the parser.

The parser supports a fixed set of component type symbols that uniquely identify a given component type.

_Supported Component Type Symbols:_

* `A` - Attributes
* `A,p` - Attributes Property
* `D` - Deontic
* `I` - Aim
* `Bdir` - Direct Object
* `Bdir,p` - Direct Object Property
* `Bind` - Indirect Object
* `Bind,p` - Indirect Object Property
* `Cac` - Activation Condition
* `Cex` - Execution Constraint
* `E` - Constituted Entity
* `E,p` - Constituted Entity Property
* `M` - Modal
* `F` - Constitutive Function 
* `P` - Constituting Properties
* `P,p` - Constituting Properties Property
* `O` - Or else

#### Component Combinations

Statements often express alternatives or combinations of activities that need to be 
administered, thus reflecting logical combinations of entities, such as actors, actions, conditions, etc. 
To combine individual components logically, IG Script supports the notion of *Component 
Combinations*. As indicated above, separate components of the same type are interpreted as 
AND-combined. To explicitly specify the nature of the logical relationship (e.g., conjunction, 
inclusive/exclusive disjunction), combinations need to be explicitly specified in the following format:

`componentSymbol(componentValue1 [logicalOperator] componentValue2)`

Note: Per default, combinations are scoped to cover both component values (i.e., `componentValue1` and `componentValue2`). In practice, it may be relevant to scope combinations more narrowly within a statement to capture the statement semantics, e.g., `A((production [AND] handling) responsible)` reflecting `production responsible and handling responsible`.

Examples: 
* `A(actor [XOR] owner)` (Alternative: `A((actor [XOR] owner))`)
* `A(involved (driver [OR] passenger))` (showcasing partial scoping of combination -- resulting in `involved driver or involved passenger`) 
* `I(report [XOR] review)`

Component combinations can be nested arbitrarily deep, i.e., any component can be a combination itself, for example:

`componentSymbol(componentValue1 [logicalOperator] (componentValue2 [logicalOperator] componentValue3))`

Example:
* `A(certifier [AND] (owner [XOR] inspector))` 
* `A((operator [OR] certifier) [AND] (owner [XOR] inspector))`

Where components are linked by the same logical operator, the indication of precedence is optional. 

Example:
* `A(certifier [AND] owner [AND] inspector)`

Supported logical operators:

* `[AND]` - conjunction (i.e., "and")
* `[OR]` - inclusive disjunction (i.e., "and/or")
* `[XOR]` - exclusive disjunction (i.e., "either or")
* `[NOT]` - negation (i.e., "not") -- Note: not yet fully supported by parser

Invalid operators (e.g., `[AN]`) will be ignored in the parsing process.

#### Nested Statements

Selected components can be substituted by statements entirely. 
For example, the activation condition could consist of a statement 
on its own. The syntax is as follows:

`componentSymbol{ componentSymbol(naturalLanguageText) ... }`

As before, components are annotated using parentheses, and are now augmented 
with braces that delineate the nested statements.

Essentially, a fully annotated statement (e.g., `A(), I(), Cex()`) is 
framed by the statement annotation (e.g., `Cac{ A(), I(), Cex() }`).

Nesting can occur to arbitrary depth, i.e., a nested statement 
can contain another nested statement, e.g., `Cac{ A(), I(), Cac{ A(), I(), Cac() } }`, 
and can further include combinations as introduced in the following.

It is important to note that component-level nesting is limited to specific 
component types and properties.

Components for which nested statements are supported:

* `A,p` - Attributes Property
* `Bdir` - Direct Object
* `Bdir,p` - Direct Object Property
* `Bind` - Indirect Object
* `Bind,p` - Indirect Object Property
* `Cac` - Activation Condition
* `Cex` - Execution Constraint
* `E,p` - Constituted Entity Property
* `P` - Constituting Properties
* `P,p` - Constituting Properties Property
* `O` - Or else

#### Nested Statement Combinations

Nested statements can, analogous to components, be combined to an arbitrary depth and using the same logical operators as for component combinations. Two constraints are to be considered: 

* Combinations of statements need to be surrounded by braces (i.e., `{` and `}`).
* *Only components of the same kind* can be combined, with the component being indicated left to the surrounding braces (e.g., `Cac{ ... }`).

The following example correctly reflects the combination of two nested AND-combined 
activation conditions:

`Cac{ Cac{ A(), I(), Cex() } [AND] Cac{ A(), I(), Cex() } }`

This example, in contrast, will fail (note the differing component types `Cac` and `Cex`): 

`Cac{ Cac{ A(), I(), Cex() } [AND] Cex{ A(), I(), Cex() } }`

Another important aspect are the outer braces surrounding the nested statement 
combination, i.e., form a `componentSymbol{ ... [AND] ... }` pattern (where logical operators can vary, of course).

Unlike the first example, the following one will result in an error due to missing outer braces:

`Cac{ A(), I(), Cex() } [AND] Cac{ A(), I(), Cex() }`

Nesting can be of multi levels, i.e., similar to component combinations, braces can be used to signal precedence when linking multiple component-level nested statements. 

Example: `Cac{ Cac{ A(), I(), Cex() } [AND] { Cac{ A(), I(), Cex() } [XOR] Cac{ A(), I(), Cex() } } }`

Note: the inner brace indicating precedence (the expression `... { Cac{ A(), I(), Cex() } [XOR] Cac{ A(), I(), Cex() } }` in the previous example) does not require the leading component symbol.

Non-nested and nested components can be used in the same statement (e.g., `... Cac( text ), Cac{ A(text) I(text) Cac(text) } ...` ). Those are implicitly AND-combined.

#### Component Pair Combinations

Where a range of components together form an alternative in a given statement (require linkage to another range of components by logical operators), IG Script supports the ability to indicate such so-called *component pair combinations* or *component tuple combinations*. 

For instance the statement `A(actor) D(must) {I(perform action) on Bdir(object1) Cex(in a particular way) [XOR] I(prevent action) on Bdir(object2) by Cex(some specific means)}` draws on the same actor, but -- in contrast to combinations of nested statements -- we see *combinations of pairs of different components* in this statement, effective rendering those as two distinct statements. Braces (without indication of the component symbol as done for nested statement combinations) are used to indicate the scope of a given component pair (here `I(perform action) on Bdir(object1) Cex(in a particular way)` as first component pair, and `I(prevent action) on Bdir(object2) by Cex(some specific means)` as the second one, both of which are combined by `[XOR]`; but both statements have the same attribute `actor` and deontic `must`). Operationally, the parser expands this into two distinct (but logically linked) statements: 
* `A(actor) D(must) I(perform action) on Bdir(object1) Cex(in a particular way)`
* `[XOR]`
* `A(actor) D(must) I(prevent action) on Bdir(object2) by Cex(some specific means)`

Note further that either pair or tuple can consist of an arbitrary number of components and can be imbalanced and use different components on either side. For example the statement `A(actor) D(must) {I(perform action) [XOR] I(prevent action) on Bdir(object2) and affects Bind(object3) by Cex(some specific means)}`, with a single component on the left side (`I(perform action)`) and multiple on the right side (`I(prevent action) on Bdir(object2) and affects Bind(object3) by Cex(some specific means)`) is equally valid input and decomposes into the statements
* `A(actor) D(must) I(perform action)`
* `[XOR]`
* `A(actor) D(must) I(prevent action) on Bdir(object2) and affects Bind(object3) by Cex(some specific means)`

The use of component pairs can occur on any level of nesting, including top-level statements (as shown in the first example), within nested statements, statement combinations, and embed basic component combinations (e.g., `A(actor) D(must) {I(perform action) on Bdir((objectA [AND] objectB)) Cex(in a particular way) [XOR] I(prevent action) on Bdir(objectC) by Cex(some specific means)}`). 

Furthermore, an arbitrary number of component pairs/tuples can be combined by using the syntax (similar to nested statement combinations) to indicate precedence amongst multiple component pairs with varying logical operators.

Example: `A(actor1) {I(action1) Bdir(directobject1) [XOR] {I(action2) Bdir(directobject2) Bind(indirectobject2) [AND] I(action3) Bdir(directobject3) Cex(constraint3)}} Cac(condition1)`

In this example the center part reflects the combination of different component pairs using the following logical pattern `{ ... [XOR] { ... [AND] ... }}`, all of which share the same attribute (`actor1`) and activation condition (`condition1`). 

#### Object-Property Relationships

Entities such as Attributes, Direct Object, Indirect Object, Constituted Entity and Constitutive Properties often carry private properties specific to a particular instance of that component (e.g., where multiple components of the same type exist).

An example is `Bdir,p(shared) Bdir1,p(private) Bdir1(object1) Bdir(object2)`, where both Direct Objects (Bdir) have a shared property (`Bdir,p(shared)`), but only one has an additional private property (`Bdir1,p(private)`) that is exclusively linked to `object1` (`Bdir1(object1)`).

In IG Script this is reflected based on suffices associated with the privately related components, where both need to carry the same suffix (i.e., `1` to signal direct linkage between `Bdir1,p` and `Bdir1` in the above example).

The basic syntax (without annotations -- see below) is `componentSymbolSuffix(component content)`, where the component symbol (`componentSymbol`) reflects the entity or property of concern, and the suffix (`Suffix`) is the identifier of the private linkage between particular instances of the related components (i.e, the suffix `1` identifies the relationship between `Bdir1,p` and `Bdir1`). The syntax further supports suffix information on properties (e.g., `Bdir1,p1(content)`, `Bdir1,p2(content2)`) to reflect dependency structures embedded within given components or their properties (here: `content` as the first property, and `content2` as the second property of `Bdir1` -- where of analytical relevance).

The coding of component-property relationships ensures that the specific intra-statement relationships are correctly captured and accessible to downstream analysis.

Suffixes can be attached to any component type, but private property linkages (i.e., linkages between particular types of components/properties) are currently supported for the following component-property pairs:

* `A` and `A,p`
* `Bdir` and `Bdir,p`
* `Bind` and `Bind,p`
* `E` and `E,p`
* `P` and `P,p`
* `I` and `Cex`

Note that the extended syntax that supports Object-Property Relationships is further augmented with the ability to capture IG Logico's Semantic Annotations as discussed in the following.

#### Semantic Annotations

In addition to the parsing of component annotations and combinations of various kinds, the parser further supports semantic annotations of components according to the taxonomies outlined in the [Institutional Grammar 2.0 Codebook](https://arxiv.org/abs/2008.08937).

The syntax (including support for suffices introduced above) is `componentSymbolSuffix[semanticAnnotation](component content)`, i.e., any component can be augmented with `[semantic annotation content]`, e.g., `Cac[context=state](Upon certification)`. 

This also applies to nested components, e.g., `Cac[condition=violation]{A[entity=actor,animate](actor) I[act=violate](violates) Bdir[entity=target,inanimate](something)}`, as well as for compound components, e.g., `Bdir[type=target](leftObject [XOR] rightObject)`, in which case the annotation `type=target` is attached to both `leftObject` and `rightObject` in the generated output. 

In addition to component-level annotations, the parser also supports statement-level annotations. The syntax for those is simply any bracket-based expression outside parentheses or braces, e.g., `A(actor) I(aim) Cac(condition) [Statement-level Annotation]`, where `[Statement-level Annotation]` is the corresponding annotation. Where multiple statement-level annotations are presented, they are concatenated in the generated output, e.g., `[annotation1][annotation2]`. Positioning of annotations in the encoding is arbitrary (`A(actor) [an annotation here] I(aim) Cac(condition) [and another one here]`), as long as it not embedded in nested components.

To include annotations in the output, the Parameter `Include IG Logico annotations` needs to be activated in the Parameters section of the tabular version of the parser.

Annotations are then generated in the corresponding columns (for statement-level annotations as well as annotations on nested components (e.g., `Cac[Nested Component Annotation]{ ... }`) in the `Statement Annotation` column), and for components, annotations are provided in the corresponding annotation column following the component column (e.g., `A (Annotation)` for Attributes). 

### Examples

In the following, you will find selected examples that highlight the practical use of the features introduced above. These can be tested and explored using the parser.

* Simple regulative statement: 
  * `A(Operator) D(must) I(comply) Bdir(with regulations)`
* Component combination: 
  * `A(Operator) D(must) I(comply with [AND] respond to) Bdir(regulations)`
* Component-level nesting: 
  * `A(Farmer) D(must) I(comply) with Bdir(provisions) Cac{A(Farmer) I(has lodged) for Bdir(application) Bdir,p(certification)}`
* Component-level nesting over multiple levels: 
  * `A(Farmer) D(must) I(comply) with Bdir(Organic Farming provisions) Cac{A(Farmer) I(has lodged) for Bdir(application) Bdir,p(certification) Cex(successfully) with the Bind(Organic Farming Program) Cac{E(Organic Program) F(covers) P,p(relevant) P(region)}}`
* Component-level nesting on various components (e.g., Bdir and Cac) embedded in combinations of nested statements (Cac): 
  * `A(Program Manager) D(may) I(administer) Bdir(sanctions) Cac{Cac{A(Program Manager) I(suspects) Bdir{A(farmer) I(violates [OR] does not comply) with Bdir(regulations)}} [OR] Cac{A(Program Manager) I(has witnessed) Bdir,p(farmer's) Bdir(non-compliance) Cex(in the past)}}`
* Complex statement; showcasing combined use of various features (e.g., component-level combinations, nested statement combinations (activation conditions, Or else)): 
  * `A,p(National Organic Program's) A(Program Manager), Cex(on behalf of the Secretary), D(must) I(inform), I(review [AND] (recertify [XOR] sanction)) Bdir(approved (certified production [AND] handling operations [AND] accredited certifying agents)) Cex(according to the guidelines provided in the (Act [XOR] regulations in this part)) if Cac{Cac{A(Program Manager) I(suspects [OR] establishes) Bdir(violations)} [AND] Cac{E(Program Manager) F(is authorized) for the P,p(relevant) P(region)}}, or else O{O{A,p(Manager's) A(supervisor) D(may) I(suspend [XOR] revoke) Bdir,p(Program Manager's) Bdir(authority)} [XOR] O{A(regional board) D(may) I(warn [OR] fine) Bdir,p(violating) Bdir(Program Manager)}}`
* Object-Property Relationships; showcasing linkage of private nodes with specific component instances (especially where implicit component-level combinations occur), alongside a shared property that apply to both objects (note that for this example, the AND-linkage between the different direct objects is implicit):
  * `A,p(Certified) A(agent) D(must) I(request) Bdir,p(independently) Bdir1,p(authorized) Bdir1(report) and Bdir2,p(audited) Bdir2(financial documents).`
* Semantic annotations; showcasing the basic use of annotations on arbitrary components:
  * `A,p[property=qualitative](Certified) A[role=responsible](organic farmers) D[stringency=obligation](must) I[type=act](submit) Bdir[type=inanimate](report) to Bind[role=authority](Organic Program Representative) Cac[context=tim](at the end of each year).`
* Semantic annotation on statement-level (this can of course be combined with component-level annotations, but separated here for clarity):
  * `A(responsible actor) [statement annotation0] I(action) Cac(condition) [statement annotation1]`

 ### Common issues
  * Parentheses/braces need to match. The parser calls out if a mismatch exists. Parentheses are applied when components are specified (or combinations thereof), and braces are used to signal component-level nesting (i.e., statements embedded within components) and statement combinations (i.e., combinations of component-level nested statements).
    * *One important special case to bear in mind* is that the input itself (i.e., the content to be encoded) may contain parentheses (e.g., `... under the following conditions 1) the actor must not ..., 2) the actor must not ...`. In such cases, the parentheses must be removed as part of the preprocessing of the text input, since the parser cannot differentiate between parentheses used for the encoding and the ones contained in the content.
  * Ensure whitespace between symbols and logical operators to ensure correct parsing. However, excess words between parentheses and logical operator are permissible. (Note: The parser has builtin tolerance toward various issues, but may not handle this correctly under all circumstances.)
    * This example works: `A(actor) I(act)`
    * This one is incorrect: `A(actor)I(act)` due to missing whitespace.
    * This example works: `(A(actor1) [AND] A(actor2))`
    * This one is incorrect: `(A(actor1)[AND] A(actor2))` due to missing whitespace between Attributes component and operator.
  * Suffices to indicate object-property relationships need to be specified immediately following the component identifier, not following the property indicator.
    * This example works: `A1(actor1) A1,p(approved)`
    * This one does not: `A1(actor1) A,p1(approved)` -- here the suffix 1 would not be linked to the Attributes (A1), but qualify the property as the first property (as opposed to second (i.e., A,p2), etc.
  * In nested statement combinations, only components of the same kind can be logically linked (e.g., `Cac{Cac{ ... } [AND] Cac{ ... }}` will parse; `Cac{Cac{ ... } [AND] Bdir{ ... }}` will not parse).
  * Component pairs have a similar syntax as nested statement combinations (e.g., `{I() Bdir() [AND ] I() Bdir()}`), but *no presence of leading component symbol* since they, unlike nested statement combinations, allow for the *combination of pairs of different components* (e.g., action-object pairs), not just single components of the same kind (see `Cac{ ... }` syntax in the previous comment on nested statement combinations). This encoding leads to distinctively different outcomes in the statement structure (observable in tabular and visual output): For nested statement combinations the individual nested components generate individual nested statements linked to the same main statement (i.e., everything is still retained as a single statement), component pair combination encoding leads to the generation (extrapolation) of entirely separate but logically-linked statements in which the shared parts of the encoding are replicated across those additional statements. The following statements may be useful to observe the difference in the generated output structure (best done using the visual version of the parser):
  * Nested statement combination -- Combinations of a specific nested component type
    * `A(actor1) I(action1) Bdir{Bdir{A(actor2) I(action2)} [XOR] Bdir{A(actor2) I(action2)}} Cac(condition1)`   
  * Component pair combination -- Combination of multiple different component types (nested or non-nested)
    * `A(actor1) {I(action1) Bdir(object1) [XOR] I(action2) Bdir(object2)} Cac(condition1)` 

## Deployment

This section is particularly focused on the setup of IG Parser, not the practical use discussed above. IG Parser can both be run on a local machine, or deployed on a server. The corresponding instructions are provided in the following, alongside links to the prerequisites. 

Note that the server-deployed version is more reliable when considering production use.

### Local deployment

The purpose of building a local executable is to run IG Parser on a local machine (primarily for personal use on your own machine).

* Prerequisites:
  * Install [Go (Programming Language)](https://go.dev/dl/)
  * Open a console (e.g., Linux terminal, Windows command line (`cmd`) or PowerShell (`powershell`))
  * Clone this repository into a dedicated folder on your local machine
  * Navigate to the repository folder
  * Compile IG Parser in the corresponding console
    * Under Windows, execute `go build -o ig-parser.exe ./web`
      * This creates the executable `ig-parser.exe` in the repository folder
    * Under Linux, execute `go build -o ig-parser ./web`
      * This creates the executable `ig-parser` in the repository folder
  * Run the created executable
    * Under Windows, run `ig-parser` (or `ig-parser.exe`) either via command line or by doubleclicking
    * Under Linux (or Windows PowerShell), run `./ig-parser`
  * Once started, it should automatically open your browser and navigate to http://localhost:8080/visual. Alternatively, use your browser to manually navigate to one of the URLs listed in the console output. By default, this is the URL http://localhost:8080 (and http://localhost:8080/visual respectively)
  * Press `Ctrl` + `C` in the console window to terminate the execution (or simply close the console window)

### Server deployment

The purpose of deploying IG Parser on a server is to provide a deployment that allows remote use on the local network or the internet, as well as for production-level deployment (see comments at the bottom).

* Prerequisites:
  * Install [Docker](https://docs.docker.com/engine/install/)
    * Quick installation of docker under Ubuntu LTS: `sudo apt install docker.io`
  * Install [Docker Compose](https://docs.docker.com/compose/install/)
    * Quick installation of docker under Ubuntu LTS: `sudo apt install docker-compose`
  * If not already installed (check with `git version`), install [Git](https://git-scm.com/) (optional if IG Parser sources are downloaded as zip file)
    * Quick installation of git under Ubuntu LTS: `sudo apt install git`

* Deployment Guidelines
  * Clone (or download and unzip) this repository into dedicated local folder
  * Navigate into cloned (or unzipped) folder
  * Make `deploy.sh` executable (`chmod 740 deploy.sh`)
  * Run `deploy.sh` with superuser permissions (`sudo ./deploy.sh`)
    * This script automatically deploys the latest version of IG Parser by undeploying old versions, before pulling the latest version, building and deploying it.
    * For manual start, run `sudo docker-compose up -d`. Run `sudo docker-compose down` to stop the execution.
  * Open browser and enter the server address and port 4040 (e.g., http://server-ip:4040)
  
* Service Configuration & Additional Considerations
  * By default, the Docker-deployed web service listens on port 4040, and logging is enabled in the subfolder `./logs`.
  * The service automatically restarts if it crashes or if the docker daemon restarts. 
  * Adjust the docker-compose.yml file to modify any of these characteristics.
  * The service is exposed as http service by default. For production-level deployment, consider using an environment that provides additional security features (e.g., SSL, DDoS protection, etc.), as is the case for the deployed version linked at the top of this page.
