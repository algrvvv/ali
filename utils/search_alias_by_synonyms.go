package utils

import "slices"

func SearchSynonyms(aliases map[string]AliasEntry, search string) *AliasEntry {
	// сначала ищем только по основному алиасу, а не по синонимам
	for alias, entry := range aliases {
		if alias == search {
			return &entry
		}
	}

	// а только потом проверяем по синонимам
	for _, entry := range aliases {
		if slices.Contains(entry.Aliases, search) {
			return &entry
		}
	}

	return nil
}
