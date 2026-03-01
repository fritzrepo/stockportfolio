# Comdirect API
## Anmeldevorgang
2.1 Access token holen:<br>
Credentials nicht als json, sondern als form-url-encoded senden.
Session_id und request_id erzeugen. Wird in die headers des nächsten Requests eingetragen.

2.2 Session Identifier (Session-Objekt) holen

    Die Client Request-Id wird vom Client erstellt.
    Beispiel für eine Client Request-Id:
    {
    "clientRequestId": {
        "sessionId": "550e8400e29b11d4a716446655440000",
        "requestId": "140113250"
        }
    }

    Diese muss bei jedem Request im Header stehen.
    HTTP-Header „x-http-request-info“

2.3 Session Identifier durch TAN challenge validieren (2 Factor Auth)
Im Response Header "x-once-authentication-info" steckt die eigentliche Antwort.
Die TAN als Photo.
Aus Response Header "x-once-authentication-info" die Challange-ID holen

Wichtig!
Zwischenschritt: Warten auf die Freigabe per Photo-App

2.4 Aktiviere die Session TAN
 - TAN vom PhotoTAN-Generator eintragen.
Wenn das TAN-Verfahren photoTAN-Push genutzt wird, erfolgt die Freigabe der TAN in der comdirect photoTAN-App. Eine TAN-Eingabe im Header der Schnittstelle ist damit nicht mehr erforderlich. Der Parameter wird deshalb in diesem Falle nicht im Header benötigt.

In Request Header "x-once-authentication-info" die Challenge-ID eintragen.

2.5 Zugriffsrechte innerhalb der Session erweitern
    Bislang stand der Scope auf 2Factor. Wir erweitern den Scope, damit wir auf die eigentliche API zugreifen können. Wir erhalten einen neuen Access-Token als auch einen Refresh-Token.

Access-Token ist 10 Minuten gültig und kann vor Ablauf mit dem refresh-token erneuert werden.

Fehlermeldungen sind im Header-Feld „x-http-response-info“ untergebracht.


| transactionType | Übersetzung |
|---|---|
| Sparplan | Saving Plan |
| Wertpapier | Securities |
| Geldanlage | Investment Saving |
| Bankgebühren | Bank fees |
| Sonstiges | Miscellaneous |
| Bar | Cash |
| Zinsen / Dividenden | Interest / Dividends |
| Devisen | Currency Exchange |
| Storno | Cancellation |
| Scheck | Cheque |
| Lastschrift | Direct Debit |
| Überweisung | Transfer |
| Kartenverfügung | Card transaction |
| Sorten (Kasse) | Foreign Currency exchange |
| Geldautomat | ATM Withdrawal |
| Geldanlage | Savings |
| Dauerauftrag | Standing Order |


## Implementierte Requests
- Anmelden bzw. Ssssion Token verwalten
- Session schließen
- Alle Konten laden
- Einzelnes Konto laden
- Letzte Umsätze laden
- Alle Depots laden (nötig für die DepotId, die für die anderen Depotabfragen benötigt wird.)
- Depot-Positionen laden
- Einzelne Position laden
- Depot Transaktionen laden
- Instruments (Info über Wertpapier laden)
- Orders