# Root Language Types and Orchestration

This document describes the root `lang` package types that serve as the central orchestration layer for all linguistic components in the Kismet language system.

## Overview

The root `lang` package provides the `Language` struct and related types that coordinate between specialized linguistic modules:
- **phoneme**: Fundamental sound units and pools
- **phonology**: Syllable structure and phonotactic rules  
- **orthography**: Writing system mappings and styles
- **morphology**: Word formation and lexicon generation
- **grammar**: Sentence structure and agreement systems

## Core Types

### LanguageID

A BCP-47 inspired identifier system for uniquely identifying languages:

```go
type LanguageID struct {
    Family   string // Language family (elv, hum, orc, dwf, etc.)
    Branch   string // Sub-family branch (wood, high, black, etc.)
    Language string // Specific language name (sindarin, quenya, etc.)
    Dialect  string // Regional or social dialect variant
    Script   string // Writing system identifier
    Region   string // Geographic or cultural region
}
```

**Format**: `[family]-[branch]-[language]-[dialect]-[script]-[region]`

**Examples**:
- `elv-wood-sindarin` - Wood Elf Sindarin
- `hum-germanic-old-english-latin` - Old English with Latin script
- `orc-black-speech` - Orcish Black Speech
- `dwf-mountain-khuzdul-cirth` - Dwarven Khuzdul with Cirth script

**Features**:
- Hierarchical organization from broad to specific
- Flexible parsing (1-6 components supported)
- String representation for easy serialization
- Validation and error handling

### LanguageType

Classification of languages based on origin and context:

```go
const (
    LanguageTypeUnknown     // Undetermined type
    LanguageTypeNatural     // Naturally evolved language
    LanguageTypeConstructed // Artificially constructed language
    LanguageTypeDivine      // Language of divine/supernatural origin
    LanguageTypeAncient     // Ancient or extinct language
    LanguageTypeModern      // Contemporary or actively used language
)
```

### LanguageComplexity

Overall complexity level of a language:

```go
const (
    LanguageComplexityUnknown     // Undetermined complexity
    LanguageComplexitySimple      // Simple grammar, small phoneme inventory
    LanguageComplexityModerate    // Moderate complexity, balanced features
    LanguageComplexityComplex     // Complex grammar, rich morphology
    LanguageComplexityVeryComplex // Highly complex, many grammatical categories
)
```

### Language

The central orchestrator that integrates all linguistic components:

```go
type Language struct {
    ID          LanguageID           // Unique identifier
    Name        string               // Human-readable name
    Type        LanguageType         // Classification
    Complexity  LanguageComplexity   // Complexity level
    Culture     string               // Associated culture
    Description string               // Description
    
    // Core linguistic components
    Phonology   *phonology.Phonology     // Sound system
    Orthography *orthography.WritingSystem // Writing system
    Morphology  *morphology.MorphemeList  // Word formation
    Grammar     *grammar.Grammar          // Sentence structure
    
    // Metadata and evolution tracking
    CreatedAt   time.Time    // Creation timestamp
    EvolvedAt   time.Time    // Last evolution timestamp
    ParentID    *LanguageID  // Parent language reference
    ChildIDs    []LanguageID // Child language references
    
    // Configuration
    Seed        int64        // Random seed for generation
    RNG         *rand.Rand   // Random number generator
}
```

## Key Features

### 1. Component Integration

The `Language` struct serves as a facade that coordinates all linguistic subsystems:

```go
lang := NewLanguage(id, "Sindarin", LanguageTypeConstructed, 12345)

// Set up phonology
lang.SetPhonology(phonology.NewPhonology(pool))

// Set up writing system  
lang.SetOrthography(orthography.NewWritingSystem())

// Set up morphology
lang.SetMorphology(morphology.NewMorphemeList())

// Set up grammar
lang.SetGrammar(grammar.NewGrammar())
```

### 2. Completeness Checking

Languages can verify they have all required components:

```go
if !lang.IsComplete() {
    missing := lang.GetMissingComponents()
    // missing = ["phonology", "orthography", "morphology", "grammar"]
}
```

### 3. Language Evolution

Support for language families and evolution:

```go
// Clone a language to create a dialect
dialect := original.Clone(newID, newSeed)
dialect.ParentID = &original.ID
original.ChildIDs = append(original.ChildIDs, dialect.ID)
```

### 4. Deterministic Generation

Each language has a seeded RNG for reproducible generation:

```go
lang := NewLanguage(id, "Quenya", LanguageTypeConstructed, 12345)
// lang.RNG is initialized with seed 12345
// All generation will be deterministic with this seed
```

## Usage Patterns

### Creating a New Language

```go
// 1. Define the language ID
id := LanguageID{
    Family:   "elv",
    Branch:   "high", 
    Language: "quenya",
}

// 2. Create the language
lang := NewLanguage(id, "Quenya", LanguageTypeConstructed, 12345)

// 3. Configure properties
lang.SetComplexity(LanguageComplexityComplex)
lang.SetCulture("elvish")
lang.SetDescription("The ancient language of the High Elves")

// 4. Add linguistic components
lang.SetPhonology(createPhonology())
lang.SetOrthography(createOrthography())
lang.SetMorphology(createMorphology())
lang.SetGrammar(createGrammar())
```

### Language Identification

```go
// Parse from string
id, err := ParseLanguageID("elv-wood-sindarin")
if err != nil {
    log.Fatal(err)
}

// Convert back to string
idStr := id.String() // "elv-wood-sindarin"

// Access individual components
family := id.Family   // "elv"
branch := id.Branch   // "wood"
lang := id.Language   // "sindarin"
```

### Language Families

```go
// Create parent language
parent := NewLanguage(LanguageID{Family: "hum"}, "Proto-Germanic", LanguageTypeAncient, 1000)

// Create child languages
oldEnglish := parent.Clone(LanguageID{Family: "hum", Language: "old-english"}, 2000)
german := parent.Clone(LanguageID{Family: "hum", Language: "german"}, 3000)

// Track relationships
parent.ChildIDs = []LanguageID{oldEnglish.ID, german.ID}
```

## Design Principles

### 1. Composition Over Inheritance

The `Language` struct composes linguistic components rather than inheriting from them, allowing flexible configuration and easy testing.

### 2. Dependency Injection

All linguistic components are injected via setter methods, enabling:
- Easy mocking for testing
- Flexible component swapping
- Clear dependency relationships

### 3. Immutable Core

The core `LanguageID` and enums are immutable, while the `Language` struct provides controlled mutation through setter methods.

### 4. Deterministic Behavior

Seeded RNGs ensure reproducible generation, essential for:
- Consistent world generation
- Debugging and testing
- Version control of generated content

### 5. Extensibility

The design supports future extensions:
- Additional linguistic components
- More complex evolution models
- Cultural and historical linking
- Performance optimizations

## Integration Points

### With Other Packages

- **phoneme**: Provides sound units for validation
- **phonology**: Ensures syllable well-formedness
- **orthography**: Enables written representation
- **morphology**: Supplies word-building blocks
- **grammar**: Provides sentence structure

### With External Systems

- **Culture System**: Links languages to cultural entities
- **World History**: Tracks language evolution over time
- **Name Generation**: Provides consistent naming across cultures
- **Text Generation**: Enables procedural content creation

## Future Enhancements

### Planned Features

1. **Language Contact**: Modeling influence between neighboring languages
2. **Diachronic Change**: Systematic sound shifts and grammatical evolution
3. **Loanword Integration**: Borrowing and adaptation mechanisms
4. **Register Variation**: Formal/informal language variants
5. **Performance Optimization**: Caching and lazy loading of components

### Extension Points

1. **Custom Phonotactics**: Language-specific sound combination rules
2. **Morphological Templates**: Configurable word formation patterns
3. **Grammatical Features**: Extensible agreement and inflection systems
4. **Writing System Styles**: Customizable orthographic representations
5. **Evolution Algorithms**: Pluggable language change models

## Conclusion

The root `lang` types provide a solid foundation for orchestrating complex linguistic systems while maintaining flexibility and extensibility. The BCP-47 inspired identification system enables clear language classification, while the `Language` struct serves as a central coordinator for all linguistic components.

This design supports the Kismet system's goals of:
- Procedural language generation
- Cultural and historical consistency
- Deterministic but diverse output
- Modular and maintainable architecture
- Future expansion and enhancement
