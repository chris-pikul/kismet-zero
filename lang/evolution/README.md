# Language Evolution Package

The `evolution` package provides comprehensive mechanisms for simulating diachronic language change and cross-cultural linguistic influence. It enables languages to evolve over time through systematic sound shifts, morphological simplifications, and contact-induced changes like loanwords and blended phonologies.

## Overview

This package simulates how real languages evolve over time, either naturally within themselves or from other cultures' influence. Real-world examples include:

- **American English** derived from British English with slight evolution in spellings and sounds
- **Modern Turkish** adopting parts of French and Arabic into their dialect and grammar
- **Grimm's Law** in Germanic languages (systematic sound shifts)
- **Romance languages** evolving from Latin through natural changes and contact

## Key Features

### 1. Sound Change System
- **Systematic phonological changes** like Grimm's Law, Verner's Law analogues
- **Environment-based rules** (e.g., "stops become fricatives between vowels")
- **Era-specific evolution** (early, middle, late evolution periods)
- **Historically attested patterns** from real linguistic research

### 2. Morphological Evolution
- **Simplification changes** (reducing case systems, agreement complexity)
- **Regularization** (making irregular forms regular over time)
- **Innovation** (introducing new grammatical features)
- **Complexity tracking** and measurement

### 3. Orthographic Evolution
- **Writing system changes** over time and through contact
- **Script reforms**: spelling, character simplification, system-wide reforms
- **Borrowing**: graphemes, mapping patterns, writing styles from other systems
- **Phonological adaptation**: writing system updates to reflect sound changes
- **Complexity tracking** and measurement for writing systems
- **Writing System Family Trees**: genealogical relationships and evolution history

### 4. Contact-Induced Change
- **Multiple contact types**: trade, conquest, migration, cultural, religious, educational
- **Borrowing hierarchy**: lexical (most common) → phonological → orthographic → grammatical (least common)
- **Intensity-based influence** calculation
- **Cultural context** and duration factors

### 4.5. Cultural Influence Modeling ⭐ **NEW**
- **Sophisticated cultural compatibility** calculations
- **Fantasy-focused contact types**: trade, conquest, migration, religious, magical, ancient
- **Cultural identity factors**: prestige, innovation, preservation, magical tradition, religious influence
- **Power dynamics** and cultural resistance modeling
- **Selective borrowing patterns** with cultural adaptation
- **Fantasy elements**: magical influence, ancient civilization effects, supernatural contact

### 5. Language Family Tree
- **Genealogical relationships** between languages
- **Divergence tracking** and timing estimation
- **Evolution history** recording
- **Contact event** logging

## Quick Start

### Basic Language Evolution

```go
package main

import (
    "time"
    "github.com/chris-pikul/kismet-zero/lang/evolution"
)

func main() {
    // Create evolution engine
    config := evolution.DefaultEvolutionConfig(42)
    engine := evolution.NewEvolutionEngine(config)
    
    // Create a base language
    protoLang := createProtoLanguage()
    
    // Evolve the language over time
    evolvedLang, evolutionEvent, err := engine.EvolveLanguage(
        protoLang,
        "early_evolution",
        time.Hour*24*365*100, // 100 years
    )
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Language evolved: %s\n", evolvedLang.Name)
    fmt.Printf("Changes applied: %d\n", len(evolutionEvent.Changes))
}

### Orthographic Evolution

```go
// Apply orthographic changes to a writing system
orthoEngine := evolution.NewOrthographicEvolutionEngine(config)

// Apply natural orthographic evolution
changes := orthoEngine.ApplyOrthographicChanges(
    language.Orthography,
    "middle_evolution",
)

// Generate a script reform
reform := orthoEngine.GenerateOrthographicReform(
    language.Orthography,
    "spelling",
    "reform_era",
)

// Calculate writing system complexity
complexity := orthoEngine.CalculateOrthographicComplexity(language.Orthography)
```
```

### Creating Child Languages

```go
// Create a new language that evolved from a parent
germanicID := lang.LanguageID{
    Family:   "indoeuropean",
    Branch:   "germanic",
    Language: "proto_germanic",
}

childLang, evolutionEvent, err := engine.CreateChildLanguage(
    evolvedLang,
    germanicID,
    "Proto-Germanic",
    "middle_evolution",
)
```

### Simulating Language Contact

```go
// Simulate trade contact between languages
contactEvolved, evolutionEvent, err := engine.EvolveLanguageFromContact(
    targetLang,
    sourceLang,
    evolution.ContactTypeTrade,
    0.7, // High intensity
    time.Hour*24*365*200, // 200 years
)
```

## Configuration

The `EvolutionConfig` struct controls various aspects of language evolution:

```go
config := evolution.EvolutionConfig{
    NaturalChangeRate:     0.3,  // Probability of natural changes per era
    SoundShiftProbability: 0.4,  // Likelihood of sound changes
    MorphologyChangeRate:  0.2,  // Rate of morphological evolution
    OrthographyChangeRate: 0.15, // Rate of orthographic evolution
    ScriptReformRate:      0.1,  // Rate of script reforms
    ContactInfluenceRate:  0.25, // Base rate of contact influence
    BorrowingThreshold:    0.3,  // Minimum intensity for borrowing
    AdaptationStrength:    0.6,  // How strongly borrowed features adapt
    CulturalInfluenceWeight: 0.4, // Weight of cultural factors
    EraDuration:           time.Hour * 24 * 365 * 100, // 100 years
    Seed:                  42,    // RNG seed for reproducibility
}
```

## Evolution Types

### Natural Evolution
- **Internal changes** that occur over time
- **Simplification** of complex grammatical systems
- **Regularization** of irregular forms
- **Innovation** of new features

### Contact Evolution
- **Lexical borrowing**: words and vocabulary
- **Phonological borrowing**: sound patterns and phonemes
- **Grammatical influence**: syntax and morphology
- **Cultural adaptation**: how borrowed features integrate

### Cultural Influence Modeling ⭐ **NEW**
- **Sophisticated cultural compatibility** calculations
- **Fantasy-focused contact types**: trade, conquest, migration, religious, magical, ancient
- **Cultural identity factors**: prestige, innovation, preservation, magical tradition, religious influence
- **Power dynamics** and cultural resistance modeling
- **Selective borrowing patterns** with cultural adaptation
- **Fantasy elements**: magical influence, ancient civilization effects, supernatural contact

### Sound Changes
- **Grimm's Law**: voiced stops → voiceless stops
- **Lenition**: stops → fricatives between vowels
- **Palatalization**: velar → palatal before front vowels
- **Custom rules**: user-defined phonological changes

## Family Tree Management

The package maintains a complete genealogical tree of languages:

```go
familyTree := engine.GetFamilyTree()

// Get ancestors of a language
ancestors := familyTree.GetAncestors("germanic-proto_germanic")

// Get descendants
descendants := familyTree.GetDescendants("indoeuropean-proto")

// Find common ancestor
commonAncestor, err := familyTree.GetCommonAncestor("germanic-old_english", "germanic-old_high_german")

// Print the tree
fmt.Println(familyTree.PrintFamilyTree())
```

## Writing System Family Trees

The package also maintains genealogical relationships between writing systems:

```go
// Create a writing system family tree
tree := evolution.NewWritingSystemFamilyTree(config)

// Add writing systems with genealogical relationships
tree.AddWritingSystem(protoWriting, nil, "origin")
tree.AddWritingSystem(childWriting, protoWriting, "natural_evolution")
tree.AddWritingSystem(siblingWriting, protoWriting, "style_evolution")

// Navigate family relationships
ancestors := tree.GetAncestors(childWriting.Name)
descendants := tree.GetDescendants(protoWriting.Name)
siblings := tree.GetSiblings(childWriting.Name)

// Calculate similarity and divergence times
similarity := tree.CalculateWritingSystemSimilarity(childWriting.Name, siblingWriting.Name)
divergenceTime := tree.EstimateDivergenceTime(childWriting.Name, siblingWriting.Name)

// Track evolution and borrowing history
tree.AddEvolutionEvent(writingSystem.Name, evolutionChange)
tree.AddBorrowingEvent(writingSystem.Name, borrowingEvent)

// Generate new writing system lineages
newLineage := tree.GenerateWritingSystemLineage(
    parentWriting,
    "Child-Name",
    orthography.WritingStyleSyllabic,
    evolutionEngine,
    "modern_era",
)
```

## Evolution Events

Every change is recorded as a `LinguisticChange` with detailed metadata:

```go
type LinguisticChange struct {
    ID          string          // Unique identifier
    Type        ChangeType      // Sound shift, morphological, lexical, etc.
    Direction   ChangeDirection // Additive, subtractive, modifying, blending
    Description string          // Human-readable description
    Details     string          // Technical details
    Timestamp   time.Time       // When the change occurred
    Era         string          // Evolution era
    Trigger     string          // What caused the change
    Intensity   float32         // Strength of the change (0.0-1.0)
}
```

## Use Cases

### World-Building
- **Generate language families** with realistic evolution
- **Model cultural interactions** through linguistic borrowing
- **Create historical depth** with evolution timelines
- **Maintain consistency** across generated content

### Historical Linguistics
- **Test linguistic theories** through simulation
- **Model sound change** patterns and timing
- **Study contact effects** between languages
- **Reconstruct proto-languages** through reverse evolution

### Educational Tools
- **Demonstrate language change** concepts
- **Interactive linguistics** lessons
- **Historical language** reconstruction exercises
- **Cultural influence** modeling

## Advanced Usage

### Custom Sound Change Rules

```go
// Create a custom sound change rule
customRule := evolution.SoundChangeRule{
    ID:          "custom_palatalization",
    Name:        "Custom Palatalization",
    Description: "Custom palatalization rule for specific context",
    FromPhonemes: []phoneme.Phoneme{/* ... */},
    ToPhonemes:   []phoneme.Phoneme{/* ... */},
    Environment:  "before high front vowels",
    Probability:  0.6,
    Era:         "custom_era",
}

// Add to sound change engine
soundEngine := evolution.NewSoundChangeEngine(config)
soundEngine.AddRule(customRule)
```

### Complex Evolution Scenarios

```go
// Simulate a complex evolution scenario
// 1. Natural evolution over time
evolvedLang, _, err := engine.EvolveLanguage(baseLang, "early", time.Hour*24*365*500)

// 2. Contact with another language
contactEvolved, _, err := engine.EvolveLanguageFromContact(
    evolvedLang, otherLang, evolution.ContactTypeConquest, 0.8, time.Hour*24*365*100)

// 3. Create a divergent branch
branchLang, _, err := engine.CreateChildLanguage(
    contactEvolved, branchID, "Branch Language", "late")

// 4. Continue evolution
finalLang, _, err := engine.EvolveLanguage(branchLang, "modern", time.Hour*24*365*200)
```

## Performance Considerations

- **RNG seeding** ensures reproducible evolution
- **Efficient tree traversal** for large language families
- **Configurable change rates** to balance realism vs. performance
- **Lazy evaluation** of complex linguistic changes

## Future Enhancements

- **Orthographic evolution** (writing system changes)
  - Writing system family trees and genealogy
  - Cultural influence modeling through orthographic borrowing
  - Advanced script reform simulation
- **Semantic drift** (meaning changes over time)
- **Dialect formation** (regional variation)
- **Sociolinguistic factors** (social class, age, gender)
- **Machine learning** integration for more realistic patterns

## Contributing

When contributing to the evolution package:

1. **Follow Go best practices** and the project's coding standards
2. **Add comprehensive tests** for new features
3. **Document linguistic concepts** with academic references
4. **Maintain backward compatibility** for existing evolution scenarios
5. **Consider performance implications** for large-scale simulations

## References

- **Historical Linguistics**: Campbell, L. (2013). Historical Linguistics: An Introduction
- **Sound Change**: Blevins, J. (2004). Evolutionary Phonology
- **Language Contact**: Thomason, S. G. (2001). Language Contact: An Introduction
- **Indo-European**: Fortson, B. W. (2010). Indo-European Language and Culture: An Introduction
