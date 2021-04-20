# IG-Parser
Parser for IG 2.0 Statements based on the IG Script Notation. 

## Overview

IG-Parser is a parser for IG Script, a formal notation for the representation institutional statements 
(e.g., policy statements) used in the Institutional Grammar 2.0. 
The parser can be used by direct invocation, as well as via a web interface 
that produces tabular output of parsed statements (currently supporting Google Sheets format). 
Deployment instructions are shown at the bottom. 

## IG Script

IG Script is a notation introduced in the context of the Institutional Grammar 2.0 (IG 2.0) that aims at a deep 
structural representation of legal statements alongside selected levels of expressiveness. While IG 2.0 highlights the 
conceptual background, the objective of IG Script is to provide an accessible, but formal approach to provide a 
format-independent representation of institutional statements. While the parser currently only supports export in 
tabular format, future refinements will include other formats (e.g.,  XML, JSON, YAML).

### Principles of IG Script Syntax

IG Script centers around a set of fundamental primitives that can be combined to parse statements comprehensively, including:

* Component annotations
* Component combinations
* Nested statements
* Nested statement combinations
* Semantic Annotations (not yet supported)

#### Component Annotations

*Component Annotations* provide the basic building block for any statement annotation. A component is represented as 

* `componentSymbol(naturalLanguageText)`
  
`componentSymbol` is one of the supported symbols for the different component types (see below). 
An example for an annotated Attributes component is `A(Farmer)`.

`naturalLanguageText` is the human-readable text annotated as the entity the annotation describes. The text is
open-ended and can include special symbols, including parentheses (e.g., `A(Farmer (e.g., organic farmer))`). 
Exceptions to this rule are discussed in the context of combinations).

The scope of a component is specified by opening and closing parentheses.

All components of a statement are annotated correspondingly, without concern for order, or repetition. Multiple 
component annotations of the same kind. Multiple Attributes, 
for example (e.g., `A(Farmer) D(must) I(comply) A(Certifier)`, are effectively interpreted as a combination (i.e., 
`(Farmer [AND] Certifier)`) in the parsing process.  

Any symbols outside the annotations and combinations are ignored. The parser is
thus robust against non-annotated text.

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

`(componentAnnotation1 [logicalOperator] componentAnnotation2)`

Example: 
* `(A(actor) [XOR] A(owner))`
* `(I(report) [XOR] I(review))`

Component combinations can be nested arbitrarily deep, i.e., any component can be a combination itself, for example:

`(componentAnnotation1 [logicalOperator] (componentAnnotation2 [logicalOperator] componentAnnotation3))`

Example:
* `(A(certifier) [AND] (A(owner) [XOR] A(inspector)))` 
* `((A(operator) [OR] A(certifier)) [AND] (A(owner) [XOR] A(inspector)))`

Supported logical operators:

* `[AND]` - conjunction (i.e., "and")
* `[OR]` - inclusive disjunction (i.e., "and/or")
* `[XOR]` - exclusive disjunction (i.e., "either or")

#### Nested Statements

Selected components can be substituted by statements entirely. 
For example, the activation condition could consist of a statement 
on its own. The syntax is as follows:

`componentSymbol{ componentSymbol(naturalLanguageText) ... }`

Essentially, a fully annotated statement (e.g., `A(), I(), Cex()`) is 
framed by the statement annotation (e.g., `Cac{ A(), I(), Cex() }`).

Nesting can occur to arbitrary depth, i.e., a nested statement 
can contain another nested statement, e.g., `Cac{ A(), I(), Cac{ A(), I(), Cac() } }` 

#### Nested Statement Combinations

Nested statements can, analogous to components, be combined to an 
arbitrary depth and using the same logical operators as for component
combinations. The only constraint is that only components of the same 
kind can be combined.  

The following example would be accepted:

`{ Cac{ A(), I(), Cex() } [AND] Cac{ A(), I(), Cex() }}`

* Note the outer braces surrounding the individual nested statements.

The following example will result in an error due to missing outer braces:

`Cac{ A(), I(), Cex() } [AND] Cac{ A(), I(), Cex() }`

Non-nested and nested components can be used in the same statement 
(e.g., `... Cac( text ), Cac{ A(text) I(text) Cac(text) } ...` )

### General Comments

* Examples
  * `A(Operator) D(must) I(comply) Bdir(with regulations)`
  * `A(Operator) D(must) I((comply [AND] respond)) Bdir(with/to regulations)`
  * `A(National Organic Program's Program Manager), Cex(on behalf of the Secretary), 
D(may) I(inspect and), I(sustain (review [AND] (revise [AND] resubmit))) 
Bdir(approved (certified production and [AND] handling operations and [AND] accredited certifying agents)) 
Cex(for compliance with the (Act or [XOR] regulations in this part)) 
Cac{A(Programme Manager) I((suspects [AND] assesses)) Bdir(violations)}`

* Note the explicit specification of the combination using parentheses (i.e., `A(( actor1 [XOR] actor2 ))`); 
  unscoped combinations lead to errors (i.e., `A( actor1 [XOR] actor2 )`). 
  The same applies for braces used for statement-level combinations, i.e., 
  `{Cac{A(actor1) I(complies) Cac(at all time)} [OR] Cac{A(actor1) I(complies) Cac(at all time)}}`.

## Deployment

* Prerequisites:
  * [Docker](https://docs.docker.com/engine/install/)
  * [docker-compose](https://docs.docker.com/compose/install/) (optional if environment, volume and port parameterization is done manually)
  * Quick installation of docker under Ubuntu LTS: `sudo apt install docker.io`

* Deployment Guidelines
  * Clone this repository
  * Make `deploy.sh` executable (`chmod 740 deploy.sh`)
  * Run `deploy.sh` with superuser permissions (`sudo deploy.sh`)
    * This script deletes old versions of IG-Parser, before pulling the latest version and deploying it.
    * For manual start, run `sudo docker-compose up -d`. Run `sudo docker-compose down` to stop execution.
  
* Service Configuration
  * By default, the web service listens on port 4040, and logging is enabled in the subfolder `./logs`. 
  * The service automatically restarts if it crashes. 
  * Adjust the docker-compose.yml file to modify any of these characteristics.
