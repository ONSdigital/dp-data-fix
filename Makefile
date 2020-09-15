.PHONY: install
install:
	go install github.com/ONSdigital/dp-data-fix

.PHONY: debug
debug:
	go build -o dp-data-fix
	./dp-data-fix pdfs -m=/Users/dave/Desktop/zebedee-data/content/zebedee/master