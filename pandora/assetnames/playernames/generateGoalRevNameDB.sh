rm -rf tmp
mkdir -p tmp
ipython parseMatthias.py
echo "...purging Spanish names..."
patch tmp/goalRevNames purgeSpanishNames.patch
python3 parseSpanishNames.py
patch tmp/goalRevNames purgeSpanishNames2.patch
python3 parseEnglishSurnames.py
python3 parseSpanishSurnames.py
python3 parseBelgNethPolItalySurnames.py
ipython3 generateSQL.py
echo "...generated database: goalRev.db...DONE"
echo "run this command to update the GO code:"
echo "cp goalRev.db ~/gits/crypto-soccer/go/names/sql/names.db"