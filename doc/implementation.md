## Data types
Im Modul "Depot" werden für alle Geld-Beträge float64 verwendet. Eigentlich wäre int auf Basis von Cent besser. Würde Rundungsfehler vermeiden und schneller in der Verarbeitung. Allerdings gibt es auch mittlerweile die Möglichkeit Teile einer Aktie zu kaufen. Bei Fonds sowieso. Dann benötigt man wieder float. Oder auch hier eine andere Basis ( mal 100) verwenden?
