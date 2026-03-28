package spec

// Spec represents a behavioral test specification for an AI agent.
type Spec struct {
	Name   string `yaml:"name"`
	Input  Input  `yaml:"input"`
	Expect Expect `yaml:"expect"`
}

// Input defines the input message to the agent.
type Input struct {
	Message string `yaml:"message"`
}

// Expect defines the expected agent behavior.
type Expect struct {
	Tools       []ToolExpect `yaml:"tools"`
	Response    *Response    `yaml:"response,omitempty"`
	Constraints []Constraint `yaml:"constraints,omitempty"`
}

// ToolExpect defines an expected tool call.
type ToolExpect struct {
	Name string                 `yaml:"name"`
	Args map[string]interface{} `yaml:"args,omitempty"`
}

// Response defines expected response content.
type Response struct {
	Contains string `yaml:"contains,omitempty"`
	Exact    string `yaml:"exact,omitempty"`
	Regex    string `yaml:"regex,omitempty"`
}

// Constraint defines a behavioral constraint.
type Constraint struct {
	NoTool   string `yaml:"no_tool,omitempty"`
	MaxTools int    `yaml:"max_tools,omitempty"`
	Ordered  bool   `yaml:"ordered,omitempty"`
}
