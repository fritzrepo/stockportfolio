## Transaction:
Jeder Kauf (buy) und Verkauf (sell) stellt eine Transaktion dar. Die Transaktionen werden in einer Liste in der Reihefolge ihres Auftretens gespeichert.

## Unclosed transactions:
Sind Transaktionen die noch nicht abgerechnet sind. Bedeutet, das das Asset im Depot vorhanden ist. Die unclosed Transaktionen können nur vom Typ "buy" sein, da bei Verkaufs-Transaktionen "sell" die Abrechnung (Gewinn / Verlust) ausgelöst wird.
Jede Abrechnung erzeugt einen "Realized Gains" (Gewinn / Verlust) Datensatz.

## Sell Transaktionen lösen eine Abrechnung aus
Wenn die nächste Transaktion vom Typ "sell" ist, wird zu diesem Asset die erste vorhandene unclosed transaction gesucht.
Es gilt das FiFo-Prinzip (First in, first out).

Drei mögliche Abrechnungen gibt es dann:
1. Anzahl der Assets ist gleich
2. Anzahl der sell Assets ist kleiner
3. Anzahl der sell Assets ist größer

#### Anzahl der Assets ist gleich
- Gewinn / Verlust ausrechnen
- Buy-Transaktion aus "unclosed transactions" löschen
- Die original buy Transaktion auf "IsClosed" setzen

#### Anzahl der sell Assets ist kleiner
- Gewinn / Verlust ausrechnen
- Von der Anzahl der buy Assets (unclosed transactions) die sell Assets abziehen und die uclosed tranaction mit ihrer neuen Anzahl der Assets speichern.

#### Anzahl der sell Assets ist größer
- Erste unclosed transaction behandeln wie in "Anzahl der Assets ist gleich".
- Mit den restlichen sell Assets wieder von vorne anfangen.
- Erstellt für jede und jede angefangene Buy-Transaktion eine Abrechnung


Noch zu tun:
Test erstellen von einfachen realisierten Transanktionen und Depoteinträgen bis zu komplizierten Kauf- und Verkauftransaktionen.
