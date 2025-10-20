## Transactions:
Jeder Kauf (buy) und Verkauf (sell) stellt eine Transaktion dar. Die Transaktionen werden in einer Liste in der Reihefolge ihres Auftretens gespeichert. Zur Zeit wird beim Einfügen einer Transaktion kein Duplikat-Test durchgeführt. Dieser könnte anhand einer Auftragsnummer gemacht werden.

## Unclosed transactions:
Sind Transaktionen die noch nicht abgerechnet sind. Bedeutet, das das Asset im Depot vorhanden ist. Die unclosed transactions können nur vom Typ "buy" sein, da bei Verkaufs-Transaktionen "sell" die Abrechnung (Gewinn / Verlust) ausgelöst wird. Eine oder mehrere unclosed transactions vom gleichen Asset bilden einen Depoteintrag.

## Realized gains:
Jede Abrechnung erzeugt einen "Realized Gains" (Gewinn / Verlust) Datensatz.
Die Gebühren beim Kauf als auch beim Verkauf werden vom Gewinn abgezogen.
Sollten für ein Asset mehrere Kauf oder Verkauf Transaktionen existieren, so werden alle möglichen Gebühren bei der ersten Abrechnung eingerechnet. In den nächsten Abrechnungen wird nur noch die jeweilige fehlende Gebühr eingerechnet.
Beispiel:
Zwei Käufe, ein Verkauf. Die Gebühr vom ersten Kauf und die Gebühr beim Verkauf sind in der ersten Abrechnung enthalten.
Bei der zweiten Abrechnung ist nur die Gebühr des zweiten Kaufs enthalten, weil die Verkaufsgebühr schon in der ersten Abrechnung enthalten ist.
Ein Kauf, zwei Verkäufe. Die Kauf-Gebühr und die Gebühr des ersten Verkaufs sind in der ersten Abrechnung enthalten. In der zweiten Abrechnung ist nur noch die Gebühr des zweiten Verkaufs enthalten.

Der Gewinn bzw. Verlust pro Abrechnung ist dadurch etwas ungenau. Aber in der Gesamt-Statistik ist die Gebühr richtig berechnet.
Ansonsten könnte man die Gebühr erst berechnen (und dann auch anteilsmässig) nachdem, alle Assets verkauft sind.

Eine andere Lösung wäre, die Gebühren bei der Abrechnung zu ignorieren und die Gebühren über die Transaktionen zu berechnen und sichtbar zu machen. Es gibt aber auch Broker, das ist die Gebühr schon im Preis des Assets enthalten. Die Gebührenausweisung ist dann nicht oder nur schwer möglich.

## Persistenz
Alle transactions, unclosed transactions und realized gains werden in der Db abgespeichert.

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
- Von der Anzahl der buy Assets (unclosed transactions) die sell Assets abziehen und die uclosed transaction mit ihrer neuen Anzahl der Assets speichern.

#### Anzahl der sell Assets ist größer
- Erste unclosed transaction behandeln wie in "Anzahl der Assets ist gleich".
- Mit den restlichen sell Assets wieder von vorne anfangen.
- **Erstellt für jede und jede angefangene Buy-Transaktion eine Abrechnung**

## Abrechnungen (Realized Gains) und offene Transaktionen (unclosed transactions) und Depotbestand berechnen
Vor Nutzung des Programms können, wenn vorhanden, bereits getätigten Transaktionen importiert werden. Sollten Transaktionen importiert worden sein, so können Gewinne / Verluste (Abrechnungen), offene Transaktionen und der Depotbestand mit "ComputeAllTransactions" berechnet werden. Die Abrechnungen und die offenen Transaktionen müssen danach persistiert werden, um bei einem Neustart, nicht die Berechnung der Abrechnungen und offenen Transaktionen wiederholen zu müssen. Der Depotbestand wird immer anhand der offenen Transaktionen berechnet. Wenn eine neue Sell-Transaktion hinzu kommt, wird die Abrechnung mit dieser und der passende(n) unclosed transaction(s) berechnet. Dann werden die unclosed transcations aktualisiert. Handelt es sich um eine Buy-Transaktion, so werden nur die unclosed transactions aktualisiert. Bei jeder hinzugefügten Transaktion wird der Depotbestand neu berechnet.
