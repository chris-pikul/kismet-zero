package interlingua

// InterlinguaServices provides access to interlingua realizers and analyzers for a language.
type InterlinguaServices struct {
	Realizers map[string]Realizer // "en","fr",...
	Analyzers map[string]Analyzer
}

// NewInterlinguaServices creates a new interlingua services container.
func NewInterlinguaServices() *InterlinguaServices {
	return &InterlinguaServices{
		Realizers: make(map[string]Realizer),
		Analyzers: make(map[string]Analyzer),
	}
}

// AddRealizer adds a realizer for a specific language code.
func (is *InterlinguaServices) AddRealizer(langCode string, realizer Realizer) {
	is.Realizers[langCode] = realizer
}

// AddAnalyzer adds an analyzer for a specific language code.
func (is *InterlinguaServices) AddAnalyzer(langCode string, analyzer Analyzer) {
	is.Analyzers[langCode] = analyzer
}

// GetRealizer returns a realizer for the specified language code.
func (is *InterlinguaServices) GetRealizer(langCode string) (Realizer, bool) {
	realizer, exists := is.Realizers[langCode]
	return realizer, exists
}

// GetAnalyzer returns an analyzer for the specified language code.
func (is *InterlinguaServices) GetAnalyzer(langCode string) (Analyzer, bool) {
	analyzer, exists := is.Analyzers[langCode]
	return analyzer, exists
}

// HasRealizer checks if a realizer exists for the specified language code.
func (is *InterlinguaServices) HasRealizer(langCode string) bool {
	_, exists := is.Realizers[langCode]
	return exists
}

// HasAnalyzer checks if an analyzer exists for the specified language code.
func (is *InterlinguaServices) HasAnalyzer(langCode string) bool {
	_, exists := is.Analyzers[langCode]
	return exists
}

// SupportedLanguages returns a list of language codes that have realizers.
func (is *InterlinguaServices) SupportedLanguages() []string {
	var languages []string
	for langCode := range is.Realizers {
		languages = append(languages, langCode)
	}
	return languages
}
