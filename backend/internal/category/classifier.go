package category

import "strings"

type categoryRule struct {
	name     string
	keywords []string
}

// Ordered rules make classification deterministic across runs.
var categoryRules = []categoryRule{
	{name: "Security", keywords: []string{"security", "cybersecurity", "infosec", "application security"}},
	{name: "Design", keywords: []string{"designer", "ux", "ui", "product design", "visual design"}},
	{name: "Data", keywords: []string{"data", "analytics", "scientist", "machine learning", "ml", "ai"}},
	{name: "Product", keywords: []string{"product manager", "product owner", "product"}},
	{name: "Engineering", keywords: []string{"engineer", "developer", "software", "backend", "frontend", "full stack", "devops", "platform", "sre", "mobile", "ios", "android"}},
	{name: "Marketing", keywords: []string{"marketing", "growth", "seo", "content", "brand"}},
	{name: "Sales", keywords: []string{"sales", "account executive", "business development"}},
	{name: "People", keywords: []string{"recruiter", "talent", "people operations", "human resources", "hr"}},
	{name: "Finance", keywords: []string{"finance", "accounting", "financial"}},
	{name: "Operations", keywords: []string{"operations", "program manager", "technical program manager"}},
	{name: "Customer Success", keywords: []string{"customer success", "support", "customer support", "success manager"}},
	{name: "Legal", keywords: []string{"legal", "counsel", "paralegal", "compliance"}},
}

var normalizedAliases = map[string]string{
	"engineering": "Engineering",
	"eng":         "Engineering",
	"developer":   "Engineering",
	"software":    "Engineering",
	"technology":  "Engineering",

	"data":      "Data",
	"analytics": "Data",

	"product": "Product",

	"design": "Design",

	"marketing": "Marketing",
	"growth":    "Marketing",

	"sales":                "Sales",
	"business development": "Sales",

	"hr":              "People",
	"people":          "People",
	"people ops":      "People",
	"human resources": "People",
	"recruiting":      "People",

	"finance":    "Finance",
	"accounting": "Finance",

	"security": "Security",

	"operations": "Operations",
	"ops":        "Operations",

	"customer success": "Customer Success",
	"support":          "Customer Success",

	"legal":      "Legal",
	"compliance": "Legal",
}

func inferFromText(text string) string {
	text = strings.ToLower(strings.TrimSpace(text))
	if text == "" {
		return "Other"
	}

	for _, rule := range categoryRules {
		for _, keyword := range rule.keywords {
			if strings.Contains(text, keyword) {
				return rule.name
			}
		}
	}

	return "Other"
}

// Normalize first tries to map source-provided categories to canonical values.
// If that fails, it falls back to text inference.
func Normalize(rawCategory string, fallbackText string) string {
	normalized := strings.ToLower(strings.TrimSpace(rawCategory))
	if normalized != "" {
		if mapped, ok := normalizedAliases[normalized]; ok {
			return mapped
		}

		if inferred := inferFromText(normalized); inferred != "Other" {
			return inferred
		}
	}

	return inferFromText(fallbackText)
}
