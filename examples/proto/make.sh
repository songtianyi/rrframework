for f in `find . -name \*.proto`
do 
	protoc --go_out=. --proto_path=/usr/include/google --proto_path=. $f
done
