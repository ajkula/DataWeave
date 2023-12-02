#!/bin/bash

# Attendre que SQL Server soit prêt
echo "Attente du démarrage de SQL Server..."
sleep 30s

# Exécuter le script SQL pour configurer le compte 'sa'
echo "Exécution du script SQL..."
/opt/mssql-tools/bin/sqlcmd -S localhost -U sa -P "StrongP@ssw0rd!" -d master -i /var/opt/mssql/scripts/mssql_init.sql

echo "Script d'initialisation terminé."
