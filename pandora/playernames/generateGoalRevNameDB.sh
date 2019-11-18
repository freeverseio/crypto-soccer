python3 parseMatthias.py
python3 parseEnglishSurnames.py
python3 generateSQL.py
echo "...purging Spanish names..."
patch patch goalRevNames purgeSpanishNames.patch
echo "...generated database: goalRev.db...DONE"
