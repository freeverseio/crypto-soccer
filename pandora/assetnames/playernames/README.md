pip3 install web3
brew install postgresql
brew install openssl

# check brew info openssl and follow it:
# echo 'export PATH="/usr/local/opt/openssl@1.1/bin:$PATH"' >> ~/.bash_profile
# export LDFLAGS="-L/usr/local/opt/openssl@1.1/lib"
# export CPPFLAGS="-I/usr/local/opt/openssl@1.1/include"


pip3 install db_utils
sqlite3 goalRev.db; .output; .dump; .exit



