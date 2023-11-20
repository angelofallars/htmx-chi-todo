#!/usr/bin/env sh

deleteRenamedTemplFiles() {
    shopt -s globstar
    for builtFile in **/*_templ.go; do
        originalFile=$(echo "$builtFile" | sed -e "s/_templ\.go/.templ/")

        if [ ! -f "$originalFile" ]; then
            echo $originalFile deleted, deleting output $builtFile
            rm $builtFile
        fi
    done

}

deleteRenamedTemplFiles
