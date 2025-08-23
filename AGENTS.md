# Kismet Zero - Agents Information

Kismet Zero is a high-fantasy procedural world generation system implemented in
Go. Its core aim is to algorithmically generate deeply detailed fictional
worlds—including geography, cultures, languages, history, and
mythologies—initially for immersive reading, with future expansion toward full
RPG game mechanics. The system is modular, combining deterministic pseudo-random
algorithms for structure and LLMs for rich narrative generation. Agents should
treat Kismet Zero as a framework for creating unique, generative lore
ecosystems.

# Culture Features

## Linguistics

The linguistic module `lang` includes a global phoneme pool (consonants and
vowels with IPA-derived categorization and weights) from which individual
languages derive their phonologies. A phonology is defined as a curated,
de-duplicated subset of phonemes sampled from the pool, with local weights
adjusted to reflect language-specific rarity or commonness. This enables
procedural languages to diverge significantly from global averages while
maintaining internal consistency, supporting realistic linguistic diversity
across generated cultures.