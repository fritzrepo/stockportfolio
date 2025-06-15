## Transaction:
Jeder Kauf (buy) und Verkauf (sell) stellt eine Transaktion dar. Die Transaktionen werden in einer Liste in der Reihefolge ihres Auftretens gespeichert.

## Unclosed transactions:
Sind Transaktionen die noch nicht abgerechnet sind. Bedeutet, das das Asset im Depot vorhanden ist. Die unclosed Transaktionen können nur vom Typ "buy" sein, da bei Verkaufs-Transaktionen "sell" die Abrechnung (Gewinn / Verlust) ausgelöst wird. Mehrere unclosed transaction vom gleichen Asset bilden einen Depoteintrag.

## Realized Gains
Jede Abrechnung erzeugt einen "Realized Gains" (Gewinn / Verlust) Datensatz.

#### Sell Transaktionen lösen eine Abrechnung aus
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
- **Erstellt für jede und jede angefangene Buy-Transaktion eine Abrechnung**

## Abrechnungen (Realized Gains) und offene Transaktionen (unclosed transactions) und Depotbestand berechnen
Vor Nutzung des Programms können, wenn vorhanden, bereits getätigten Transaktionen importiert werden. Sollten Transaktionen importiert worden sein, so können Gewinne / Verluste (Abrechnungen), offene Transaktionen und der Depotbestand mit "ComputeAllTransactions" berechnet werden. Die Abrechnungen und die offenen Transaktionen müssen danach persistiert werden, um bei einem Neustart, nicht die Berechnung der Abrechnungen und offenen Transaktionen wiederholen zu müssen. Der Depotbestand wird immer anhand der offenen Transaktionen berechnet. Wenn eine neue Sell-Transaktion hinzu kommt, wird die Abrechnung mit dieser und der passende(n) unclosed transaction(s) berechnet. Handelt es sich um eine Buy-Transaktion, so werden die unclosed transactions aktualisiert. Bei jeder hinzugefügten Transaktion wird der Depotbestand aktualisiert.
