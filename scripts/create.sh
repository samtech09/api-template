#!/bin/bash
IFS=$'\n'

displayhelp(){
	echo Creates new project from github.com/samtech09/api-template
	echo syntax:
	echo "    bash setup.sh <name\/of-your\/package>"
	echo
	echo "Example: following will create new package 'github.com/someuser/apiproject'"
	echo
	echo "    bash setup.sh 'github.com\/someuser\/apiproject'"
	echo 
}


if [ $# != 1 ]
then
  displayhelp
  exit
fi

#clone repo
git clone https://github.com/samtech09/api-template.git

# find given file type in current folder and all sub folders
for f in $(find api-template/ -iname "*.go")
do 
	#echo Processing $f
	# find and replace text
	sed -i "s/github.com\/samtech09\/api-template/${1}/g" "$f"
done

#mod file
sed -i "s/github.com\/samtech09\/api-template/${1}/g" "api-template/go.mod"

#replace in Makefile
pname=`basename ${1}`
sed -i "s/api-template/${pname}/g" "api-template/Makefile"

#delete .git folder & rename
rm -rf api-template/.git
rm -rf api-template/scripts
mv api-template "${pname}"

echo "Done. Project '${pname}' is ready for development."
echo
