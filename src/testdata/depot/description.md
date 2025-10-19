## Test decription

#### Test 1
- Eine Abrechnung Verkauf gleich groß Kauf.
- Ein Asset im Depot.

#### Test 2
- Ein Asset wird gekauft. Dann ein anderes. Vom ersten Asset wird die Hälfte verkauft.

#### Test 3 
- Checken ob der Durchschnittspreis bei einem Asset das durch zwei Käufe im Depot vorhanden ist, korrekt berechnet wird.

#### Test 4
Teste, ob die sell Transaktion Gebühr wenn es zwei buy Transaktionen gibt, nur einmal vom Gewinn abgezogen wird.
Kauf 1: Anzahl 30, Kurs 100 Gebühr 6.
Kauf 2: Anzahl 20, Kurs 100 Gebühr 5.
Verkauf: Anzahl 50, Kurs 110 Gebühr 8.
Rechnung:
(Verkaufteil1 ) 30 * 110 = 3300 - (Ankauf 30 * 100) -6 -8 =  286 (Gewinn abzüglich Gebühren)
(Verkaufteil2 ) 20 * 110 = 2200 - (Ankauf 20 * 100) -5 =  195 (Gewinn abzüglich Gebühren)

#### Test 5
Es gibt zwei sell Transaktionen zu einer buy Transaktion. Testen ob die buy transaktions Gebühr nur einmal vom Gewinn abgezogen wird.
Kauf 1: Anzahl 50, Kurs 100 Gebühr 8.
Verkauf 1: Anzahl 20, Kurs 110 Gebühr 5.
Verkauf 2: Anzahl 30, Kurs 110 Gebühr 6.
Rechnung:
(Verkauf 1 ) 20 * 110 = 2200 - (Ankauf 20 * 100) -5 -8 =  187 (Gewinn abzüglich Gebühren)
(Verkauf 2 ) 30 * 110 = 3300 - (Ankauf 30 * 100) -6 =  294 (Gewinn abzüglich Gebühren)

#### Test 6
- Eine Abrechnung Verkauf gleich groß Kauf.
- Eine Abrechnung Verkauf größer als erster Kauf, aber kleiner als zweiter Kauf. Es verbleiben also Aktien im Depot.
- Zwei Assets im Depot, drei Abrechnungen.

#### Test 7
Asset_1 wird 4 mal gekauft und einmal verkauft. Der Verkauf beinhaltet 3 Käufe. Vom dritten Kauf bleibt die Hälfte übrig. Somit ist die Hälfte vom dritten Kauf und alle vom vierten Kauf im Depot.
Zusätzlich wurde ein zweites Asset zweimal gekauft und komplett verkauft.
