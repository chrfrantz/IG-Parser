ace.define("ace/mode/igscriptnotation_highlight_rules",["require","exports","module","ace/lib/oop","ace/mode/text_highlight_rules"], function(require, exports, module){"use strict";
var oop = require("../lib/oop");
var TextHighlightRules = require("./text_highlight_rules").TextHighlightRules;
var IGScriptNotationHighlightRules = function () {
    this.$rules = {
        "start": [
            {
                token : "Attribute",
                // Symbol, followed by an opening bracket
                //regex : /A(?:\d*)(?=(,p)*\s*[\[|\(|\{])/
                regex : /A(?=[\[|\(|\{])/
            },{
                token : "Attribute_Property",
                // Symbol, followed by an opening bracket
                //regex : /A(?:\d*)(?=(,p)*\s*[\[|\(|\{])/
                regex : /A\d*,p(?=[\[|\(|\{])/
            },{
                token : "Deontic",
                regex : /D(?:\d*)(?=[\[|\(|\{])/
            },{
                token : "Aim",
                regex : /I(?:\d*)(?=[\[|\(|\{])/
            },{
                token : "Direct_Object",
                regex : /Bdir(?:\d*)(?=[\[|\(|\{])/
            },{
                token : "Direct_Object_Property",
                regex : /Bdir\d*,p(?=[\[|\(|\{])/
            },{
                token : "Indirect_Object",
                //regex : /Bind(?:\d*)(?=(,p)*\s*[\[|\(|\{])/
                regex : /Bind(?:\d*)(?=[\[|\(|\{])/
            },{
                token : "Indirect_Object_Property",
                regex : /Bind\d*,p(?=[\[|\(|\{])/
            },{
                token : "Activation_Condition",
                regex : /Cac(?:\d*)(?=[\[|\(|\{])/
            },{
                token : "Activation_Constraint",
                regex : /Cex(?:\d*)(?=[\[|\(|\{])/
            },{
                token : "Constituted_Entity",
                regex : /E(?:\d*)(?=(,p)*[\[|\(|\{])/
            },{
                token : "Modal",
                regex : /M(?:\d*)(?=[\[|\(|\{])/
            },{
                token : "Constitutive_Function",
                regex : /F(?:\d*)(?=[\[|\(|\{])/
            },{
                token : "Constituting_Properties",
                regex : /D(?:\d*)(?=[\[|\(|\{])/
            },{
                token : "Constituting_Properties_Properties",
                regex : /D(?:\d*)(?=(,p)[\[|\(|\{])/
            },{
                token : "Or_Else",
                // Symbol, followed by any amount of whitespace and an opening bracket
                regex : /O(?=[\{])/
            },{
                token : "Logical_Operator",
                // Symbol within square brackets, with any amount of whitespace within the brackets
                regex : /\[(AND|OR|XOR)\]/
            },{
                token : "Semantic_Annotation",
                // Other text within square brackets
                regex : /\[(:?[^\]]+)\]/
            }]
    };
};
oop.inherits(IGScriptNotationHighlightRules, TextHighlightRules);
exports.IGScriptNotationHighlightRules = IGScriptNotationHighlightRules;

});

ace.define("ace/mode/igscriptnotation",["require","exports","module","ace/lib/oop","ace/mode/text","ace/mode/igscriptnotation_highlight_rules","ace/range"], function(require, exports, module){"use strict";
var oop = require("../lib/oop");
var TextMode = require("./text").Mode;
var IGScriptNotationHighlightRules = require("./igscriptnotation_highlight_rules").IGScriptNotationHighlightRules;
var Mode = function () {
    this.HighlightRules = IGScriptNotationHighlightRules;
    this.$behaviour = this.$defaultBehaviour;
};
oop.inherits(Mode, TextMode);
(function () {
    this.$id = "ace/mode/igscriptnotation";
}).call(Mode.prototype);
exports.Mode = Mode;

});                (function() {
                    ace.require(["ace/mode/igscriptnotation"], function(m) {
                        if (typeof module == "object" && typeof exports == "object" && module) {
                            module.exports = m;
                        }
                    });
                })();
            