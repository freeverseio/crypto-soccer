new_header="pragma solidity >=0.5.12 <0.6.2;"
for f in contracts/*.sol; do \
    sed -i.bak "1 s/^.*$/$new_header/" $f;
done
rm contracts/*bak
