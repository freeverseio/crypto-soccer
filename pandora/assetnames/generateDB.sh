#./playernames/generateGoalRevNameDB.sh
cp playernames/goalRev.db teamnames/goalRevFull.db
(cd teamnames && ipython generateSQL.py)
cp teamnames/goalRevFull.db .
echo "Now run this command:"
echo "cp goalRevFull.db ~/gits/crypto-soccer/go/names/sql/names.db"
