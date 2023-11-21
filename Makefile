SHELL:=/bin/bash

build:
	sh clean.sh
	npx tailwindcss -i assets/styles.css -o assets/styles.out.css
	templ generate

setup:
	go install
	npm install
