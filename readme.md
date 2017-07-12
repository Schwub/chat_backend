Dies ist eine Anleitung um die Chat-Applikation mit dem Go-Backend auf Windows zu starten.

# MongoDB
Für das Go-Backend muss eine MongoDB-Datenbank auf Port 27017 laufen.
MongoDB muss auf dem Rechner installiert sein.
Die Datenbank wird mit folgendem Befehl gestartet:

'''
(Pfad_zu_MongoDB)\Server\3.x\bin\mongod.exe --dbpath (Pfad_zu_diesem_Verzeichnis)\db
'''

# Go-Backend
In diesem Verzeichnis befindet sich eine für Windows kompilierte Version des Backends.
Diese kann einfach über die Eingabeaufforderung gestartet werden.

# Angular-Frontend
Um das Angular-Frontend starten zu können, muss man in das Unterverzeichnis *web* wechseln und dort den Befehl:

'''
ng serve
'''

eingeben.
