# Language Package `lang`

The primary purpose of the `lang` package is to provide fantasy-based linguistic
generation and evolution for creation of "conlangs" (constructed languages). Any
created language can be translated to, and from, the user's preferred language such
as English (default standard).

## Sub-Packages

### lang (core)
Orchestrates all linguistic components to create complete, culturally consistent languages. Provides high-level APIs for language creation, evolution, and text generation while maintaining consistency across all subsystems. Enables deterministic but diverse language generation for fantasy world-building.

### evolution
Simulates how languages change over time through natural evolution and cultural contact. Creates realistic language families with sound shifts, morphological changes, and dialect formation. Tracks genealogical relationships between languages and writing systems, enabling rich world-building with historical depth and linguistic consistency.

**Core Evolution:**
- Natural diachronic changes (sound shifts, morphological evolution)
- Contact-induced changes (borrowing, calquing, phonological adaptation)
- Language family tree tracking and genealogical relationships

**Dialect Formation:**
- Geographic dialects based on climate, terrain, population, urbanization
- Social dialects based on social class, education, urban/rural environment
- Automatic feature generation (phonological, lexical, grammatical)
- Mutual intelligibility tracking and dialect continua

**Cultural Influence:**
- Sophisticated cultural compatibility calculations
- Fantasy-focused contact types (trade, conquest, migration, religious, magical)
- Cultural identity factors (prestige, innovation, preservation, magical tradition)

**Writing Systems:**
- Orthographic evolution and script reforms
- Writing system family trees and genealogy
- Cultural influence through orthographic borrowing

**Configuration:**
- Reproducible evolution through seeded RNG
- Configurable change rates and influence weights
- Era-based evolution with customizable time periods

### phoneme
Provides the fundamental building blocks for language generation by defining consonants and vowels with their articulatory features. Creates phoneme pools for weighted random selection while maintaining consistency. Serves as the foundation for all higher-level linguistic modules, ensuring realistic sound systems for constructed languages.

**Sound Units:**
- Strongly typed consonant and vowel representations
- Articulatory features (manner, place, voicing, height, backness, rounding)
- Phoneme pools with weighted random selection
- Deduplication guarantees for consistency

**Foundation Layer:**
- Core types for phonology, morphology, and grammar
- Deterministic behavior through seeded randomness
- Composable design for higher-level modules
- Support for diverse linguistic systems

### phonology
Defines the rules and constraints for how sounds combine into syllables and words. Creates phonotactic patterns that determine valid sound combinations, syllable structures, and consonant clusters. Ensures generated languages follow realistic phonological patterns while allowing customization for different linguistic systems.

**Syllable Generation:**
- Syllable templates (CV, CVC, CCV, CCVC, etc.)
- Phonotactic rule enforcement and validation
- Weighted probabilities for syllable structures
- Customizable rules for different linguistic systems

**Sound Constraints:**
- Composable phonotactic rule units
- Cluster validation and filtering
- Deterministic generation under seeded randomness
- Bridge between phonemes and word-building

### morphology
Constructs meaningful words from phonological building blocks by defining morphemes, word formation rules, and lexicon generation. Creates culturally appropriate vocabulary and names while maintaining phonological consistency. Supports different morphological systems like agglutinative, fusional, and isolating languages.

**Word Building:**
- Morpheme types (roots, affixes, infixes, circumfixes)
- Derivational and inflectional word formation rules
- Cultural context for semantic domains
- Phonological well-formedness validation

**Lexicon Generation:**
- Culturally appropriate word inventories
- Name generation for people, places, and concepts
- Weighted selection for morpheme combinations
- Support for different language types and patterns

### orthography
Creates writing systems that map spoken sounds to written symbols, supporting alphabetic, syllabic, and logographic styles. Generates culturally appropriate writing systems that maintain consistency between spoken and written forms. Enables the creation of complete written languages for fantasy world-building.

**Writing Systems:**
- Multiple styles (alphabetic, syllabic, logographic)
- Grapheme-phoneme mapping consistency
- Cultural customization and preferences
- Deterministic generation for reproducibility

**Text Generation:**
- Consistent written representation of spoken language
- Support for names, words, and phrases
- Configurable writing system styles
- Bridge between phonology and written forms

### grammar
Organizes words into grammatically correct sentences by defining syntactic structures, word order patterns, and agreement systems. Creates culturally appropriate sentence structures while maintaining grammatical consistency. Supports different language types through flexible rule-based systems for complete language generation.

**Sentence Structure:**
- Configurable word order patterns (SVO, SOV, VSO, VOS, OVS, OSV)
- Template-based sentence generation (declarative, interrogative, imperative)
- Complex sentences with subordinate clauses and coordination
- Phrase construction (noun, verb, prepositional, adverbial)

**Agreement Systems:**
- Case marking (nominative, accusative, genitive, dative)
- Number agreement (singular, dual, plural)
- Gender agreement (masculine, feminine, neuter, animate, inanimate)
- Tense and aspect marking (past, present, future, perfective, imperfective)

**Syntactic Rules:**
- Composable rule system for grammatical phenomena
- Weighted rule application for probabilistic generation
- Validation of sentence structure and grammatical well-formedness
- Integration with morphology, phonology, and orthography

### interlingua
Provides a semantic bridge between generated languages and user languages through semantic graph representations. Enables bi-directional translation with meaning preservation, allowing users to understand generated content in their preferred language while maintaining the linguistic richness of the constructed language.

**Semantic Pivot:**
- Semantic graph representation of sentences (Events, Entities, Roles)
- TAM features (tense, aspect, mood, polarity, evidentiality, voice)
- Nominal features (person, number, gender, case, definiteness, classifiers)
- Discourse features (topic, focus, givenness, coreference)

**Translation Pipeline:**
- Conlang generation → Interlingua emission with construction tracking
- Interlingua → User language realization with loss-aware translation
- User language → Interlingua analysis for round-trip consistency
- Deterministic behavior through seeded randomness

**Integration:**
- Hooks into morphology, grammar, and phonology systems
- Maintains construction IDs and lemma hints for reversibility
- Supports multiple target languages through analyzer/realizer pairs
- Preserves linguistic features even when targets cannot realize them

## High-Level APIs

The `lang` package provides ergonomic, RNG-driven APIs that make language generation simple for simulation engines. These functions handle the complexity internally while providing consistent interfaces for world-building tasks.

### Language Creation & Management

**`CreateRandomLanguage(family, culture, seed)`** - Generates complete languages with cultural preferences in a single call. Creates realistic phoneme inventories, writing systems, and grammatical structures appropriate to the specified culture type.

**`CulturalLanguageConfig(culture)`** - Returns pre-configured settings for different fantasy cultures (elvish, dwarven, orcish, human, ancient). Each culture has appropriate complexity targets, phoneme preferences, and evolution rates.

**`DefaultLanguageConfig()`** - Provides sensible defaults for language generation when custom configuration isn't needed.

### Evolution & World-Building

**`CreateLanguageFamily(protoLanguage, numBranches, evolutionTime, seed)`** - Generates realistic language families from a proto-language. Creates branching relationships that simulate historical language development and dialect formation.

**`EvolveLanguage(language, timePeriod, culturalEvents, seed)`** - Applies time-based changes to languages, simulating how they change over centuries. Supports different evolution periods and cultural event influences.

**`SimulateLanguageContact(language1, language2, contactType, intensity, duration, seed)`** - Models how languages influence each other through trade, conquest, migration, or cultural exchange. Returns both languages modified by contact.

### Content Generation

**`GenerateNames(language, nameType, count, culturalContext, seed)`** - Creates culturally appropriate names for people, places, deities, and artifacts. Names follow cultural patterns and maintain linguistic consistency.

**`GenerateText(language, textType, length, culturalContext, seed)`** - Produces sample texts in generated languages, including lore, poetry, dialogue, and inscriptions. Content reflects cultural themes and linguistic patterns.

**`CreateCulturalLanguageSet(culture, numLanguages, diversity, seed)`** - Generates complete linguistic ecosystems for fantasy cultures, including main languages and regional variants with proper family relationships.

### Utility & Validation

**`GetLanguageInfo(language)`** - Extracts human-readable statistics and cultural information from languages. Provides phoneme counts, complexity scores, and evolution history for world-building reference.

**`ValidateLanguage(language)`** - Ensures language consistency and completeness. Reports issues, warnings, and suggestions for creating playable linguistic systems.

**`ExportLanguage(language, format)`** - Exports language data for external use, saving/loading languages or sharing them between simulation systems.

## Evolution Package Features