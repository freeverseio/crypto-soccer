# moves plots to folders and, for plots1, it joins them
# 
outDir1=tmp1
outDir2=tmp2

rm -rf $outDir1
mkdir $outDir1

for a in *hist.png; do montage $(basename $a _hist.png).png $a -tile 2x1 -geometry +0+0 ${outDir1}/$(basename $a _hist.png).png; done

rm -rf $outDir2
mkdir $outDir2
mkdir ${outDir2}/win
mkdir ${outDir2}/tie

cp *Win*png ${outDir2}/win
cp *Tie*png ${outDir2}/tie

